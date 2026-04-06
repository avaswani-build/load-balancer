package algorithms

import (
	"net/http"

	"github.com/avaswani-build/load-balancer/internal/pool"
)

type Selector interface {
	Next(r *http.Request, p *pool.ServerPool) *pool.Backend
}
