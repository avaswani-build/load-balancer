package pool

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

type Backend struct {
	URL   *url.URL
	Alive bool
	mux   sync.RWMutex
	Proxy *httputil.ReverseProxy
}

type ServerPool struct {
	Backends []*Backend
	Current  uint64
}

func (p *ServerPool) nextIndex() int {
	return int(atomic.AddUint64(&p.Current, uint64(1)) % uint64(len(p.Backends)))
}

func NewBackend(rawURL string) (*Backend, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid backend URL: %w", err)
	}
	b := &Backend{
		URL:   u,
		Alive: true,
		Proxy: httputil.NewSingleHostReverseProxy(u),
	}
	return b, nil
}

func (p *ServerPool) AddBackend(b *Backend) {
	p.Backends = append(p.Backends, b)
}

func (p *ServerPool) Next() *Backend {
	if len(p.Backends) == 0 {
		return nil
	}

	next := p.nextIndex()
	loop := len(p.Backends)

	for i := 0; i < loop; i++ {
		idx := (next + i) % loop
		if p.Backends[idx].IsAlive() {
			if i != 0 {
				atomic.StoreUint64(&p.Current, uint64(idx))
			}
			return p.Backends[idx]
		}
	}
	return nil
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
