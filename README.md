# Go Load Balancer

A simple HTTP load balancer written in Go to explore backend selection strategies and core service design principles.

## Features

- Round Robin load balancing
- Weighted Round Robin
- IP Hashing (deterministic routing)
- Pluggable selector interface
- Reverse proxy request forwarding (`httputil.ReverseProxy`)
- Dead backend skipping
- Unit tests for core selection logic

## Architecture

The project is structured with clear separation of concerns:

- `pool` — backend definitions and state (alive/dead)
- `algorithms` — selection strategies (round robin, weighted, IP hash)
- `handler` — request lifecycle and routing
- `cmd/lb` — application entrypoint

### Selector Interface

```go
type Selector interface {
    Next(r *http.Request, p *pool.ServerPool) *pool.Backend
}
```

This allows routing strategies to be swapped without modifying request handling logic.

## Example

Start the backend servers:

```bash
go run ./testservers/backend1
go run ./testservers/backend2
```

Start the load balancer:

```bash
go run ./cmd/lb
```

Send requests:

```bash
curl localhost:8080
```

Responses will be routed according to the selected algorithm.

## Design Notes

- Uses Go’s `httputil.ReverseProxy` for HTTP forwarding
- Selection logic is isolated from networking concerns
- Weighted round robin is implemented via a precomputed index space
- IP hashing uses FNV hashing on client IP (port stripped)
- Emphasis on clarity and separation of concerns over production completeness

## Testing

Run all tests:

```bash
go test ./...
```

Tests focus on:
- selection correctness
- deterministic behavior
- handling of dead backends

## Scope

This project focuses on:

> how traffic is routed to backend services

It intentionally does **not** include:

- health checks
- retries or circuit breakers
- metrics / observability
- dynamic configuration
- distributed coordination

## Future Extensions (Optional)

- Least connections strategy
- Health checks
- Consistent hashing

## Notes

This project is intended as a learning exercise and is not production-ready.
