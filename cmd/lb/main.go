package main

import (
	"fmt"
	"net/http"

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

	lb := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
	}

	http.DefaultServeMux.HandleFunc("/", lb)

	http.ListenAndServe(":8080", http.DefaultServeMux)
}
