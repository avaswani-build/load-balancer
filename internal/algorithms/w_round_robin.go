package algorithms

import (
	"net/http"
	"sync/atomic"

	"github.com/avaswani-build/load-balancer/internal/pool"
)

type WeightedRoundRobin struct {
	Indices []int
	Current uint64
}

func NewWeightedRoundRobin(backends []*pool.Backend) *WeightedRoundRobin {
	indices := make([]int, 0)

	for idx, backend := range backends {
		if backend == nil {
			continue
		}

		for j := 0; j < backend.Weight; j++ {
			indices = append(indices, idx)
		}
	}

	return &WeightedRoundRobin{
		Indices: indices,
		Current: 0,
	}
}

func (wrr *WeightedRoundRobin) weightedNextIndex() int {
	if len(wrr.Indices) == 0 {
		return 0
	}

	pos := atomic.AddUint64(&wrr.Current, 1) - 1
	return int(pos % uint64(len(wrr.Indices)))
}

func (wrr *WeightedRoundRobin) Next(_ *http.Request, p *pool.ServerPool) *pool.Backend {
	if len(p.Backends) == 0 || len(wrr.Indices) == 0 {
		return nil
	}

	start := wrr.weightedNextIndex()
	count := len(wrr.Indices)

	for i := 0; i < count; i++ {
		pos := (start + i) % count
		backendIdx := wrr.Indices[pos]

		if p.Backends[backendIdx].IsAlive() {
			if i != 0 {
				atomic.StoreUint64(&wrr.Current, uint64(pos))
			}
			return p.Backends[backendIdx]
		}
	}

	return nil
}
