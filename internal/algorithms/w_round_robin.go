package algorithms

import (
	"sync/atomic"

	"github.com/avaswani-build/load-balancer/internal/pool"
)

type WeightedSelection struct {
	Indices []int
	Current uint64
}

func NewWeightedSelection(backends []*pool.Backend) *WeightedSelection {
	indices := make([]int, 0)

	for idx, backend := range backends {
		if backend == nil {
			continue
		}

		for j := 0; j < backend.Weight; j++ {
			indices = append(indices, idx)
		}
	}

	return &WeightedSelection{
		Indices: indices,
		Current: 0,
	}
}

func weightedNextIndex(ws *WeightedSelection) int {
	if len(ws.Indices) == 0 {
		return 0
	}

	pos := atomic.AddUint64(&ws.Current, 1) - 1
	return int(pos % uint64(len(ws.Indices)))
}

func WeightedNext(ws *WeightedSelection, p *pool.ServerPool) *pool.Backend {
	if len(p.Backends) == 0 || len(ws.Indices) == 0 {
		return nil
	}

	start := weightedNextIndex(ws)
	count := len(ws.Indices)

	for i := 0; i < count; i++ {
		pos := (start + i) % count
		backendIdx := ws.Indices[pos]

		if p.Backends[backendIdx].IsAlive() {
			if i != 0 {
				atomic.StoreUint64(&ws.Current, uint64(pos))
			}
			return p.Backends[backendIdx]
		}
	}

	return nil
}
