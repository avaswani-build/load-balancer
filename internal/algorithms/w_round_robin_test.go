package algorithms

import (
	"testing"

	"github.com/avaswani-build/load-balancer/internal/pool"
)

func TestWeightedRoundRobinNext_CycleBackendsRatio1(t *testing.T) {
	b1 := &pool.Backend{
		Alive:  true,
		Weight: 3,
	}
	b2 := &pool.Backend{
		Alive:  true,
		Weight: 2,
	}

	sp := &pool.ServerPool{
		Backends: []*pool.Backend{b1, b2},
		Current:  0,
	}
	wrr := NewWeightedRoundRobin(sp.Backends)

	cycle1 := wrr.Next(nil, sp)
	cycle2 := wrr.Next(nil, sp)
	cycle3 := wrr.Next(nil, sp)
	cycle4 := wrr.Next(nil, sp)
	cycle5 := wrr.Next(nil, sp)
	cycle6 := wrr.Next(nil, sp)

	if cycle1 != b1 {
		t.Fatalf("first cycle: got %p, want %p", cycle1, b1)
	}
	if cycle2 != b1 {
		t.Fatalf("second cycle: got %p, want %p", cycle2, b1)
	}
	if cycle3 != b1 {
		t.Fatalf("third cycle: got %p, want %p", cycle3, b1)
	}
	if cycle4 != b2 {
		t.Fatalf("fourth cycle: got %p, want %p", cycle4, b2)
	}
	if cycle5 != b2 {
		t.Fatalf("fifth cycle: got %p, want %p", cycle5, b2)
	}
	if cycle6 != b1 {
		t.Fatalf("sixth cycle: got %p, want %p", cycle6, b1)
	}
}

func TestWeightedRoundRobinNext_CycleBackendsRatio2(t *testing.T) {
	b1 := &pool.Backend{
		Alive:  true,
		Weight: 1,
	}
	b2 := &pool.Backend{
		Alive:  true,
		Weight: 2,
	}
	b3 := &pool.Backend{
		Alive:  true,
		Weight: 3,
	}

	sp := &pool.ServerPool{
		Backends: []*pool.Backend{b1, b2, b3},
		Current:  0,
	}
	wrr := NewWeightedRoundRobin(sp.Backends)

	cycle1 := wrr.Next(nil, sp)
	cycle2 := wrr.Next(nil, sp)
	cycle3 := wrr.Next(nil, sp)
	cycle4 := wrr.Next(nil, sp)
	cycle5 := wrr.Next(nil, sp)
	cycle6 := wrr.Next(nil, sp)
	cycle7 := wrr.Next(nil, sp)

	if cycle1 != b1 {
		t.Fatalf("first cycle: got %p, want %p", cycle1, b1)
	}
	if cycle2 != b2 {
		t.Fatalf("second cycle: got %p, want %p", cycle2, b2)
	}
	if cycle3 != b2 {
		t.Fatalf("third cycle: got %p, want %p", cycle3, b2)
	}
	if cycle4 != b3 {
		t.Fatalf("fourth cycle: got %p, want %p", cycle4, b3)
	}
	if cycle5 != b3 {
		t.Fatalf("fifth cycle: got %p, want %p", cycle5, b3)
	}
	if cycle6 != b3 {
		t.Fatalf("sixth cycle: got %p, want %p", cycle6, b3)
	}
	if cycle7 != b1 {
		t.Fatalf("seventh cycle: got %p, want %p", cycle7, b1)
	}
}
func TestWeightedRoundRobinNext_SkipDeadBackendsRatio1(t *testing.T) {
	b1 := &pool.Backend{
		Alive:  true,
		Weight: 2,
	}
	b2 := &pool.Backend{
		Alive:  false,
		Weight: 3,
	}
	b3 := &pool.Backend{
		Alive:  true,
		Weight: 1,
	}

	sp := &pool.ServerPool{
		Backends: []*pool.Backend{b1, b2, b3},
		Current:  0,
	}
	wrr := NewWeightedRoundRobin(sp.Backends)

	cycle1 := wrr.Next(nil, sp)
	cycle2 := wrr.Next(nil, sp)
	cycle3 := wrr.Next(nil, sp)
	cycle4 := wrr.Next(nil, sp)
	cycle5 := wrr.Next(nil, sp)
	cycle6 := wrr.Next(nil, sp)

	if cycle1 != b1 {
		t.Fatalf("first cycle: got %p, want %p", cycle1, b1)
	}
	if cycle2 != b1 {
		t.Fatalf("second cycle: got %p, want %p", cycle2, b1)
	}
	if cycle3 != b3 {
		t.Fatalf("third cycle: got %p, want %p", cycle3, b3)
	}
	if cycle4 != b1 {
		t.Fatalf("fourth cycle: got %p, want %p", cycle4, b1)
	}
	if cycle5 != b1 {
		t.Fatalf("fifth cycle: got %p, want %p", cycle5, b1)
	}
	if cycle6 != b3 {
		t.Fatalf("sixth cycle: got %p, want %p", cycle6, b3)
	}
}

func TestWeightedRoundRobinNext_SkipDeadBackendsRatio2(t *testing.T) {
	b1 := &pool.Backend{
		Alive:  false,
		Weight: 4,
	}
	b2 := &pool.Backend{
		Alive:  false,
		Weight: 2,
	}
	b3 := &pool.Backend{
		Alive:  true,
		Weight: 2,
	}

	sp := &pool.ServerPool{
		Backends: []*pool.Backend{b1, b2, b3},
		Current:  0,
	}
	wrr := NewWeightedRoundRobin(sp.Backends)

	cycle1 := wrr.Next(nil, sp)
	cycle2 := wrr.Next(nil, sp)
	cycle3 := wrr.Next(nil, sp)
	cycle4 := wrr.Next(nil, sp)

	if cycle1 != b3 {
		t.Fatalf("first cycle: got %p, want %p", cycle1, b3)
	}
	if cycle2 != b3 {
		t.Fatalf("second cycle: got %p, want %p", cycle2, b3)
	}
	if cycle3 != b3 {
		t.Fatalf("third cycle: got %p, want %p", cycle3, b3)
	}
	if cycle4 != b3 {
		t.Fatalf("fourth cycle: got %p, want %p", cycle4, b3)
	}
}

func TestNewWeightedRoundRobin_BuildsIndices(t *testing.T) {
	b1 := &pool.Backend{
		Alive:  true,
		Weight: 3,
	}
	b2 := &pool.Backend{
		Alive:  true,
		Weight: 1,
	}
	b3 := &pool.Backend{
		Alive:  true,
		Weight: 2,
	}

	wrr := NewWeightedRoundRobin([]*pool.Backend{b1, b2, b3})

	want := []int{0, 0, 0, 1, 2, 2}

	if len(wrr.Indices) != len(want) {
		t.Fatalf("indices length: got %d, want %d", len(wrr.Indices), len(want))
	}

	for i := range want {
		if wrr.Indices[i] != want[i] {
			t.Fatalf("index %d: got %d, want %d", i, wrr.Indices[i], want[i])
		}
	}
}
