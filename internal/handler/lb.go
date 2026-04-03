package handler

import (
	"log"
	"net/http"

	"github.com/avaswani-build/load-balancer/internal/pool"
)

func LB(sp *pool.ServerPool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		peer := sp.Next()
		if peer == nil {
			http.Error(w, "No backend available", http.StatusServiceUnavailable)
		}

		targetURL := peer.URL.String() + r.URL.Path

		log.Printf("Request forwarded to port : %s\n", targetURL)
		peer.Proxy.ServeHTTP(w, r)
	}
}
