package proxy

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

// WeightedNode 代表一个服务节点
type WeightedNode struct {
	addr   string // 服务器地址
	weight int    // 配置的权重
}

// Proxy 代理结构体，包含服务节点映射
type Proxy struct {
	services         map[string][]*WeightedNode // 服务节点映射
	mu               sync.Mutex                 // 保护并发访问
	blacklist        map[string]bool            // 黑名单
	whitelist        map[string]bool            // 白名单
	restrictedRoutes map[string][]string        // 受限路径映射（路径 -> 允许的 IP 列表）
	rateLimit        sync.Map                   // 记录 IP 访问次数
}

// NewProxy 构造函数，初始化代理
func NewProxy() *Proxy {
	p := &Proxy{
		services:         make(map[string][]*WeightedNode),
		blacklist:        make(map[string]bool),
		whitelist:        make(map[string]bool),
		restrictedRoutes: make(map[string][]string),
	}
	p.initServices()
	return p
}

// 初始化服务节点
func (p *Proxy) initServices() {
	p.services = make(map[string][]*WeightedNode)
	p.services["auth"] = []*WeightedNode{
		{addr: "127.0.0.1:8081", weight: 10},
		{addr: "127.0.0.1:8082", weight: 10},
	}
	p.services["system"] = []*WeightedNode{
		{addr: "127.0.0.1:8083", weight: 10},
		{addr: "127.0.0.1:8084", weight: 10},
	}
	p.services["order"] = []*WeightedNode{
		{addr: "127.0.0.1:8085", weight: 10},
		{addr: "127.0.0.1:8086", weight: 10},
	}
	p.services["default"] = []*WeightedNode{
		{addr: "127.0.0.1:8099", weight: 10},
	}

	// 添加黑名单示例
	p.blacklist["192.168.1.100"] = true

	// 添加白名单示例（若为空，则默认所有 IP 可访问 -- 127.0.0.1 或者 IPV6 的 localhost == ::1 也会直接默认通过）
	p.whitelist["::1"] = true
	p.whitelist["127.0.0.1"] = true

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
	if strings.HasPrefix(path, "/auth") {
		return p.getServiceForPath("auth")
	} else if strings.HasPrefix(path, "/system") {
		return p.getServiceForPath("system")
	} else if strings.HasPrefix(path, "/order") {
		return p.getServiceForPath("order")
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
	var currentWeight int
	for _, node := range nodes {
		currentWeight += node.weight
		if randWeight < currentWeight {
			return node
		}
	}
	return nil
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

	// 黑名单检查
	if p.blacklist[clientIP] {
		http.Error(w, "403 Forbidden - 黑名单", http.StatusForbidden)
		return
	}

	// 白名单检查（如果有白名单限制）
	if len(p.whitelist) > 0 && !p.whitelist[clientIP] {
		http.Error(w, "403 Forbidden - 不在白名单", http.StatusForbidden)
		return
	}

	// 访问策略检查
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

	// 流量策略：限制每个 IP 每秒 5 次请求
	if !p.rateLimitCheck(clientIP) {
		http.Error(w, "429 请求速度过快", http.StatusTooManyRequests)
		return
	}

	// 处理 CORS
	if r.Method == "OPTIONS" {
		p.handleCORS(w, r)
		return
	}

	// 选择目标服务
	route := p.findService(r)
	if route == nil {
		http.NotFound(w, r)
		return
	}
	target := &url.URL{
		Scheme: "http",
		Host:   route.addr,
	}
	proxy := p.NewHttpReverseProxy(target)
	proxy.ServeHTTP(w, r)
}

// CORS 处理
func (p *Proxy) handleCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(http.StatusOK)
}

// 访问限制（流量控制）
func (p *Proxy) rateLimitCheck(ip string) bool {
	val, _ := p.rateLimit.LoadOrStore(ip, &sync.Mutex{})
	mu := val.(*sync.Mutex)

	mu.Lock()
	defer mu.Unlock()

	// 记录访问时间
	key := fmt.Sprintf("%s-%d", ip, time.Now().Unix())
	if _, exists := p.rateLimit.Load(key); exists {
		return false
	}
	p.rateLimit.Store(key, struct{}{})

	// 1秒后清理
	go func() {
		time.Sleep(1 * time.Second)
		p.rateLimit.Delete(key)
	}()

	return true
}
