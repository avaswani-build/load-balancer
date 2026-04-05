package algorithms

import (
	"sync/atomic"

	"github.com/avaswani-build/load-balancer/internal/pool"
)

func nextIndex(p *pool.ServerPool) int {
	return int(atomic.AddUint64(&p.Current, 1) % uint64(len(p.Backends)))
}

func Next(p *pool.ServerPool) *pool.Backend {
	if len(p.Backends) == 0 {
		return nil
	}

	start := nextIndex(p)
	count := len(p.Backends)

	for i := 0; i < count; i++ {
		idx := (start + i) % count
		if p.Backends[idx].IsAlive() {
			if i != 0 {
				atomic.StoreUint64(&p.Current, uint64(idx))
			}
			return p.Backends[idx]
		}
	}

	return nil
}
