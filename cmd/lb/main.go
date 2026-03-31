package main

import (
	"log"
	"net/http"

	"github.com/avaswani-build/load-balancer/internal/pool"
)

func lbHandler(sp *pool.ServerPool) http.HandlerFunc {
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

func main() {
	sp := pool.ServerPool{}
	b, err := pool.NewBackend("http://localhost:9001")
	if err != nil {
		panic(err)
	}
	sp.AddBackend(b)

	b, err = pool.NewBackend("http://localhost:9002")
	if err != nil {
		panic(err)
	}
	sp.AddBackend(b)

	//Always use NewServeMux for production
	http.DefaultServeMux.HandleFunc("/", lbHandler(&sp))

	http.ListenAndServe(":8080", http.DefaultServeMux)
	log.Println("Load Balancer running on :8080")
}
