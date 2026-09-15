package l7

import (
	"math/rand"
	"net"
	"sync"
	"time"
)

// DNSCache
type DNSCache struct {
	mu    sync.Mutex
	cache map[string][]net.IP
	ttl   time.Duration
	last  map[string]time.Time
}

func NewDNSCache() *DNSCache {
	return &DNSCache{
		cache: make(map[string][]net.IP),
		last:  make(map[string]time.Time),
		ttl:   60 * time.Second,
	}
}

func (d *DNSCache) Resolve(host string) []net.IP {
	d.mu.Lock()
	defer d.mu.Unlock()

	if ips, ok := d.cache[host]; ok {
		if time.Since(d.last[host]) < d.ttl {
			return ips
		}
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return nil
	}
	d.cache[host] = ips
	d.last[host] = time.Now()
	return ips
}

func (d *DNSCache) PickRandom(host string) net.IP {
	ips := d.Resolve(host)
	if len(ips) == 0 {
		return nil
	}
	return ips[rand.Intn(len(ips))]
}