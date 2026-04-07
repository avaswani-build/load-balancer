package algorithms

import (
	"hash/fnv"
	"net"
	"net/http"

	"github.com/avaswani-build/load-balancer/internal/pool"
)

type IPHash struct{}

func NewIPHash() *IPHash {
	return &IPHash{}
}

func (iph *IPHash) Next(r *http.Request, p *pool.ServerPool) *pool.Backend {
	if len(p.Backends) == 0 || r == nil {
		return nil
	}

	ip := clientIP(r)
	if ip == "" {
		return nil
	}

	start := hashIndex(ip, len(p.Backends))
	count := len(p.Backends)

	for i := 0; i < count; i++ {
		idx := (start + i) % count
		if p.Backends[idx].IsAlive() {
			return p.Backends[idx]
		}
	}

	return nil
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}

func hashIndex(s string, n int) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return int(h.Sum32() % uint32(n))
}
