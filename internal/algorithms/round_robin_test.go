package algorithms

import (
	"testing"

	"github.com/avaswani-build/load-balancer/internal/pool"
)

func TestRoundRobinNext_AlternateAliveBackends(t *testing.T) {
	rr := NewRoundRobin()
	b1 := &pool.Backend{
		Alive: true,
	}
	b2 := &pool.Backend{
		Alive: true,
	}

	sp := &pool.ServerPool{
		Backends: []*pool.Backend{b1, b2},
		Current:  0,
	}

	cycle1 := rr.Next(nil, sp)
	cycle2 := rr.Next(nil, sp)
	cycle3 := rr.Next(nil, sp)
	cycle4 := rr.Next(nil, sp)

	if cycle1 != b1 {
		t.Fatalf("first cycle: got %p, want %p", cycle1, b1)
	}
	if cycle2 != b2 {
		t.Fatalf("second cycle: got %p, want %p", cycle2, b2)
	}
	if cycle3 != b1 {
		t.Fatalf("third cycle: got %p, want %p", cycle3, b1)
	}
	if cycle4 != b2 {
		t.Fatalf("fourth cycle: got %p, want %p", cycle4, b2)
	}
}

func TestRoundRobinNext_SkipDeadBackend(t *testing.T) {
	rr := NewRoundRobin()

	dead := &pool.Backend{
		Alive: false,
	}
	alive := &pool.Backend{
		Alive: true,
	}

	sp := &pool.ServerPool{
		Backends: []*pool.Backend{dead, alive},
		Current:  0,
	}

	cycle1 := rr.Next(nil, sp)
	cycle2 := rr.Next(nil, sp)
	cycle3 := rr.Next(nil, sp)

	if cycle1 != alive {
		t.Fatalf("first cycle: got %p, want alive backend %p", cycle1, alive)
	}
	if cycle2 != alive {
		t.Fatalf("second cycle: got %p, want alive backend %p", cycle2, alive)
	}
	if cycle3 != alive {
		t.Fatalf("third cycle: got %p, want alive backend %p", cycle3, alive)
	}
}

func TestRoundRobinNext_ReturnNilNoBackends(t *testing.T) {
	rr := NewRoundRobin()

	sp := &pool.ServerPool{
		Backends: nil,
		Current:  0,
	}

	got := rr.Next(nil, sp)
	if got != nil {
		t.Fatalf("got %p, want nil", got)
	}
}

func TestRoundRobinNext_ReturnNilAllBackendsDead(t *testing.T) {
	rr := NewRoundRobin()

	b1 := &pool.Backend{
		Alive: false,
	}
	b2 := &pool.Backend{
		Alive: false,
	}

	sp := &pool.ServerPool{
		Backends: []*pool.Backend{b1, b2},
		Current:  0,
	}

	got := rr.Next(nil, sp)
	if got != nil {
		t.Fatalf("got %p, want nil", got)
	}
}
