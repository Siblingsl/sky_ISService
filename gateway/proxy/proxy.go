package proxy

import (
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
)

// WeightedNode 代表一个服务节点
type WeightedNode struct {
	addr   string // 服务器地址
	weight int    // 配置的权重
}

// Proxy 代理结构体，包含服务节点映射
type Proxy struct {
	services map[string][]*WeightedNode // 路由与服务节点的映射
	mu       sync.Mutex                 // 保护并发访问
}

// NewProxy 构造函数，初始化代理
func NewProxy() *Proxy {
	p := &Proxy{}
	p.initServices()
	return p
}

// 设置 黑白名单、访问策略、流量策略、开启 CORS、同步 swagger3.0
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
	p.services["default"] = []*WeightedNode{
		{addr: "127.0.0.1:8085", weight: 10},
	}
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
	}
	return p.getServiceForPath("default")
}

//func (p *Proxy) findService(r *http.Request) *WeightedNode {
//	if r.URL.Path == "/auth" {
//		return p.getServiceForPath("auth")
//	} else if r.URL.Path == "/system" {
//		return p.getServiceForPath("system")
//	}
//	return p.getServiceForPath("default")
//}

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

// 处理请求并转发到对应的服务节点
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
