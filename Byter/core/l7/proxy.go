package l7

import (
	"net/http"
	"net/url"
	"sync"
)

// ProxyPool
type ProxyPool struct {
	proxies []string
	mu      sync.Mutex
	idx     int
}

func NewProxyPool(list []string) *ProxyPool {
	return &ProxyPool{proxies: list}
}

func (p *ProxyPool) Next() string {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.proxies) == 0 {
		return ""
	}
	proxy := p.proxies[p.idx%len(p.proxies)]
	p.idx++
	return proxy
}

func (p *ProxyPool) ApplyToTransport(tr *http.Transport) {
	proxy := p.Next()
	if proxy == "" {
		return
	}
	u, err := url.Parse(proxy)
	if err != nil {
		return
	}
	tr.Proxy = http.ProxyURL(u)
}