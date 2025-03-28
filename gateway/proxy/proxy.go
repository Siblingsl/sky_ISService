package proxy

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sky_ISService/config"
	"strings"
	"sync"
	"time"
)

// WeightedNode 代表一个服务节点
type WeightedNode struct {
	addr   string // 服务器地址
	weight int    // 配置的权重
}

// Proxy 代理结构体，包含服务节点映射和配置项
type Proxy struct {
	services                map[string][]*WeightedNode // 服务节点映射
	mu                      sync.Mutex                 // 保护并发访问
	blacklist               map[string]bool            // 黑名单
	whitelist               map[string]bool            // 白名单
	restrictedRoutes        map[string][]string        // 受限路径映射（路径 -> 允许的 IP 列表）
	rateLimit               sync.Map                   // 记录 IP 访问次数
	EnableBlacklist         bool                       // 是否启用黑名单检查
	EnableWhitelist         bool                       // 是否启用白名单检查
	EnableRestrictedRoutes  bool                       // 是否启用受限路径检查
	EnableRateLimiting      bool                       // 是否启用速率限制
	RateLimitRequestsPerSec int                        // 每秒允许的请求次数
}

// NewProxy 构造函数，初始化代理
func NewProxy() *Proxy {
	p := &Proxy{
		services:                make(map[string][]*WeightedNode),
		blacklist:               make(map[string]bool),
		whitelist:               make(map[string]bool),
		restrictedRoutes:        make(map[string][]string),
		EnableBlacklist:         true,
		EnableWhitelist:         false,
		EnableRestrictedRoutes:  true,
		EnableRateLimiting:      true,
		RateLimitRequestsPerSec: 5, // 默认每秒 5 次请求
	}
	p.initServices()
	return p
}

// 初始化服务节点
func (p *Proxy) initServices() {
	p.services = make(map[string][]*WeightedNode)
	p.services["security"] = []*WeightedNode{
		{addr: fmt.Sprintf("%s:%s", config.GetConfig().Security.Addr, config.GetConfig().Security.Port), weight: config.GetConfig().Security.Weight1},
		{addr: fmt.Sprintf("%s:%s", config.GetConfig().Security.Addr, config.GetConfig().Security.Port1), weight: config.GetConfig().Security.Weight2},
	}
	p.services["system"] = []*WeightedNode{
		{addr: fmt.Sprintf("%s:%s", config.GetConfig().System.Addr, config.GetConfig().System.Port), weight: config.GetConfig().System.Weight1},
		{addr: fmt.Sprintf("%s:%s", config.GetConfig().System.Addr, config.GetConfig().System.Port), weight: config.GetConfig().System.Weight2},
	}
	//p.services["order"] = []*WeightedNode{
	//	{addr: "0.0.0.0:8085", weight: 10},
	//	{addr: "0.0.0.0:8086", weight: 10},
	//}
	p.services["default"] = []*WeightedNode{
		{addr: config.GetConfig().Default.Addr, weight: config.GetConfig().Default.Weight},
	}

	// 添加黑名单示例
	p.blacklist["192.168.1.100"] = true

	// 添加白名单示例（若为空，则默认所有 IP 可访问 -- 127.0.0.1 或者 IPV6 的 localhost == ::1 也会直接默认通过）
	p.whitelist["::1"] = true
	p.whitelist["127.0.0.1"] = true
	p.whitelist["192.168.10.6"] = true

	// 限制 `/admin` 只允许 `192.168.1.50` 访问
	p.restrictedRoutes["/admin"] = []string{"192.168.1.50"}
}

// NewHttpReverseProxy 创建一个新的反向代理
func (p *Proxy) NewHttpReverseProxy(target *url.URL) *httputil.ReverseProxy {
	return httputil.NewSingleHostReverseProxy(target)
}

// 根据请求路径选择服务节点
func (p *Proxy) findService(r *http.Request) *WeightedNode {
	path := r.URL.Path
	serviceMap := map[string]string{
		"/security": "security",
		"/system":   "system",
		"/order":    "order",
	}

	for prefix, service := range serviceMap {
		if strings.HasPrefix(path, prefix) {
			return p.getServiceForPath(service)
		}
	}
	return p.getServiceForPath("default")
}

// 获取指定路径的服务节点
func (p *Proxy) getServiceForPath(path string) *WeightedNode {
	p.mu.Lock()
	defer p.mu.Unlock()
	if nodes, ok := p.services[path]; ok && len(nodes) > 0 {
		return p.getTargetService(nodes)
	}
	return nil
}

// 负载均衡逻辑，使用权重选择服务节点
func (p *Proxy) getTargetService(nodes []*WeightedNode) *WeightedNode {
	totalWeight := 0
	for _, node := range nodes {
		totalWeight += node.weight
	}

	randWeight := rand.Intn(totalWeight)
	for _, node := range nodes {
		if randWeight < node.weight {
			return node
		}
		randWeight -= node.weight
	}
	return nodes[len(nodes)-1] // 兜底，防止异常
}

// 获取真实ip
func getClientIP(r *http.Request) string {
	// 检查 X-Forwarded-For 头
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0]) // 取第一个 IP
	}

	// 检查 X-Real-IP 头
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// 默认回退到 RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println("解析 RemoteAddr 失败:", err)
		return ""
	}
	return ip
}

// 处理请求并转发
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientIP := getClientIP(r)

	// 先添加 CORS 头部（即使后面拦截请求也要带 CORS 头）
	p.handleCORS(w)

	// 处理 OPTIONS 预检请求（直接返回 200）
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 黑名单检查
	if p.EnableBlacklist && p.blacklist[clientIP] {
		http.Error(w, "403 Forbidden - 黑名单", http.StatusForbidden)
		return
	}

	// 白名单检查
	if p.EnableWhitelist && len(p.whitelist) > 0 && !p.whitelist[clientIP] {
		http.Error(w, "403 Forbidden - 不在白名单", http.StatusForbidden)
		return
	}

	// 访问策略检查
	if p.EnableRestrictedRoutes {
		if allowedIPs, ok := p.restrictedRoutes[r.URL.Path]; ok {
			allowed := false
			for _, ip := range allowedIPs {
				if ip == clientIP {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(w, "403 Forbidden - 受限路径", http.StatusForbidden)
				return
			}
		}
	}

	// 流量限制检查
	if p.EnableRateLimiting && !p.rateLimitCheck(clientIP) {
		http.Error(w, "429 请求速度过快", http.StatusTooManyRequests)
		return
	}

	// 选择目标服务
	route := p.findService(r)
	if route == nil {
		http.NotFound(w, r)
		return
	}

	// 代理请求
	target := &url.URL{
		Scheme: "http",
		Host:   route.addr,
	}
	proxy := p.NewHttpReverseProxy(target)
	proxy.ServeHTTP(w, r)
}

// CORS 处理
func (p *Proxy) handleCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true") // 允许携带认证信息（如果需要）
	w.WriteHeader(http.StatusOK)
}

// 访问限制（流量控制）
func (p *Proxy) rateLimitCheck(ip string) bool {
	val, _ := p.rateLimit.LoadOrStore(ip, &sync.Mutex{})
	mu := val.(*sync.Mutex)

	mu.Lock()
	defer mu.Unlock()

	// 每秒请求次数
	windowKey := fmt.Sprintf("%s-%d", ip, time.Now().Unix())
	if _, exists := p.rateLimit.Load(windowKey); exists {
		return false
	}
	p.rateLimit.Store(windowKey, struct{}{})

	// 1秒后清理
	go func() {
		time.Sleep(1 * time.Second)
		p.rateLimit.Delete(windowKey)
	}()

	return true
}
