package main

import (
	"fmt"
	"net/http"
	"github.com/avaswani-build/load-balancer/internal/pool"
)

func main() {
	sp := pool.ServerPool{}
	b, err := pool.NewBackend("http://localhost:9001")
	sp.AddBackend(b)
	b, err = pool.NewBackend("http://localhost:9002")
	sp.AddBackend(b)

	func lb(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	}

	http.lb("/", pickServer)

	http.ListenAndServe(":8080", nil)
}