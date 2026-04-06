package main

import (
	"log"
	"net/http"

	"github.com/avaswani-build/load-balancer/internal/algorithms"
	"github.com/avaswani-build/load-balancer/internal/handler"
	"github.com/avaswani-build/load-balancer/internal/pool"
)

func main() {
	sp := pool.ServerPool{}
	b, err := pool.NewBackend("http://localhost:9001", 3)
	if err != nil {
		panic(err)
	}
	sp.AddBackend(b)

	b, err = pool.NewBackend("http://localhost:9002", 1)
	if err != nil {
		panic(err)
	}
	sp.AddBackend(b)

	selector := algorithms.NewWeightedRoundRobin(sp.Backends)
	http.DefaultServeMux.HandleFunc("/", handler.LB(&sp, selector))

	http.ListenAndServe(":8080", http.DefaultServeMux)
	log.Println("Load Balancer running on :8080")
}
