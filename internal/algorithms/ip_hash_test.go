package algorithms

import (
	"net/http"
	"testing"

	"github.com/avaswani-build/load-balancer/internal/pool"
)

func TestHashIndex_EnsureDeterministic(t *testing.T) {
	first := hashIndex("127.0.0.1", 3)

	for i := 0; i < 20; i++ {
		got := hashIndex("127.0.0.1", 3)
		if got != first {
			t.Fatalf("iteration %d: got %d, want %d", i, got, first)
		}
	}
}

func TestIPHash_SameIPDifferentBackends(t *testing.T) {
	iph := NewIPHash()

	b1 := &pool.Backend{Alive: true}
	b2 := &pool.Backend{Alive: true}
	b3 := &pool.Backend{Alive: true}

	sp := &pool.ServerPool{
		Backends: []*pool.Backend{b1, b2, b3},
	}

	r1 := &http.Request{RemoteAddr: "127.0.0.1:5000"}
	r2 := &http.Request{RemoteAddr: "127.0.0.1:6000"}

	first := iph.Next(r1, sp)
	if first == nil {
		t.Fatal("got nil backend")
	}

	for i := 0; i < 20; i++ {
		got := iph.Next(r2, sp)
		if got != first {
			t.Fatalf("iteration %d: got %p, want %p", i, got, first)
		}
	}
}

func TestIPHash_DeadFallback(t *testing.T) {
	iph := NewIPHash()

	b1 := &pool.Backend{Alive: true}
	b2 := &pool.Backend{Alive: true}
	b3 := &pool.Backend{Alive: true}

	sp := &pool.ServerPool{
		Backends: []*pool.Backend{b1, b2, b3},
	}

	req := &http.Request{RemoteAddr: "10.0.0.42:5000"}
	start := hashIndex(clientIP(req), len(sp.Backends))

	sp.Backends[start].Alive = false
	want := sp.Backends[(start+1)%len(sp.Backends)]

	got := iph.Next(req, sp)
	if got != want {
		t.Fatalf("got %p, want %p", got, want)
	}
}

func TestIPHash_Nil(t *testing.T) {
	iph := NewIPHash()

	t.Run("no_backends", func(t *testing.T) {
		sp := &pool.ServerPool{}
		r := &http.Request{RemoteAddr: "127.0.0.1:5000"}

		got := iph.Next(r, sp)
		if got != nil {
			t.Fatalf("got %p, want nil", got)
		}
	})

	t.Run("nil_request", func(t *testing.T) {
		b1 := &pool.Backend{Alive: true}
		sp := &pool.ServerPool{
			Backends: []*pool.Backend{b1},
		}

		got := iph.Next(nil, sp)
		if got != nil {
			t.Fatalf("got %p, want nil", got)
		}
	})

	t.Run("all_dead", func(t *testing.T) {
		b1 := &pool.Backend{Alive: false}
		b2 := &pool.Backend{Alive: false}

		sp := &pool.ServerPool{
			Backends: []*pool.Backend{b1, b2},
		}
		r := &http.Request{RemoteAddr: "127.0.0.1:5000"}

		got := iph.Next(r, sp)
		if got != nil {
			t.Fatalf("got %p, want nil", got)
		}
	})
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}

	var digits [20]byte
	i := len(digits)

	for n > 0 {
		i--
		digits[i] = byte('0' + n%10)
		n /= 10
	}

	return string(digits[i:])
}
