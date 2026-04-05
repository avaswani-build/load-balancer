package handler

import (
	"log"
	"net/http"

	"github.com/avaswani-build/load-balancer/internal/algorithms"
	"github.com/avaswani-build/load-balancer/internal/pool"
)

func LB(sp *pool.ServerPool) http.HandlerFunc {
	ws := algorithms.NewWeightedSelection(sp.Backends)
	return func(w http.ResponseWriter, r *http.Request) {
		peer := algorithms.WeightedNext(ws, sp)
		if peer == nil {
			http.Error(w, "No backend available", http.StatusServiceUnavailable)
			return
		}

		log.Printf("Request forwarded to port : %s\n", peer.URL)
		peer.Proxy.ServeHTTP(w, r)
	}
}
