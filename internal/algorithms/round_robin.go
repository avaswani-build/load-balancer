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

	next := nextIndex(p)
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
