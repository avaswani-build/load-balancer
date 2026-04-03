package main

import (
	"log"
	"net/http"

	"github.com/avaswani-build/load-balancer/internal/handler"
	"github.com/avaswani-build/load-balancer/internal/pool"
)

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
	http.DefaultServeMux.HandleFunc("/", handler.LB(&sp))

	http.ListenAndServe(":8080", http.DefaultServeMux)
	log.Println("Load Balancer running on :8080")
}
