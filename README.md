# Go Load Balancer

A learning project to build an HTTP load balancer in Go from scratch.

## Goals
- Learn Go networking, concurrency, and service design.
- Implement round-robin, weighted, and consistent hashing algorithms.
- Add health checks, retries, circuit breaker, and metrics.
- Deploy to Kubernetes.

## Current Status
- [x] Project scaffold
- [x] Basic HTTP server responding on `:8080`
- [ ] Two test backend servers
- [ ] Round-robin backend selection
- [ ] Request forwarding
- [ ] Health checks
- [ ] Metrics endpoint

## Running Locally
```bash
# Start the load balancer
go run ./cmd/lb

# (Later) Start the test backends
go run ./testservers/backend1
go run ./testservers/backend2