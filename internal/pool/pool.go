package pool

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL    *url.URL
	Alive  bool
	mux    sync.RWMutex
	Proxy  *httputil.ReverseProxy
	Weight int
}

type ServerPool struct {
	Backends []*Backend
	Current  uint64
}

func NewBackend(rawURL string, weight int) (*Backend, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid backend URL: %w", err)
	}
	b := &Backend{
		URL:    u,
		Alive:  true,
		Proxy:  httputil.NewSingleHostReverseProxy(u),
		Weight: weight,
	}
	return b, nil
}

func (p *ServerPool) AddBackend(b *Backend) {
	p.Backends = append(p.Backends, b)
}

func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Alive
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Alive = alive
}
