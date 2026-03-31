package pool

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

type Backend struct {
	URL    *url.URL
	Alive  bool
	mux    sync.RWMutex
	Client *http.Client
}

type ServerPool struct {
	backends []*Backend
	current  uint64
}

func (p *ServerPool) nextIndex() int {
	return int(atomic.AddUint64(&p.current, uint64(1)) % uint64(len(p.backends)))
}

func NewBackend(rawURL string) (*Backend, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid backend URL: %w", err)
	}
	b := &Backend{
		URL:    u,
		Alive:  true,
		Client: &http.Client{Timeout: 5 * time.Second},
	}
	return b, nil
}

func (p *ServerPool) AddBackend(b *Backend) {
	p.backends = append(p.backends, b)
}

func (p *ServerPool) Next() *Backend {
	if len(p.backends) == 0 {
		return nil
	}

	next := p.nextIndex()
	loop := len(p.backends)

	for i := 0; i < loop; i++ {
		idx := (next + i) % loop
		if p.backends[idx].IsAlive() {
			if i != 0 {
				atomic.StoreUint64(&p.current, uint64(idx))
			}
			return p.backends[idx]
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
