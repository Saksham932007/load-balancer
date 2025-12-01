package strategy

import (
	"sync"
	"testing"

	"github.com/Saksham932007/load-balancer/internal/backend"
)

// TestNextIndexThreadSafety verifies that the atomic counter is thread-safe
func TestNextIndexThreadSafety(t *testing.T) {
	// Create a pool with 3 backends
	backends := []*backend.Backend{
		{URL: nil, ReverseProxy: nil},
		{URL: nil, ReverseProxy: nil},
		{URL: nil, ReverseProxy: nil},
	}
	pool := NewServerPool(backends)

	const goroutines = 100
	const iterations = 1000
	var wg sync.WaitGroup

	// Launch multiple goroutines calling NextIndex concurrently
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				idx := pool.NextIndex()
				// Verify index is always within bounds
				if idx < 0 || idx >= len(backends) {
					t.Errorf("Index out of bounds: %d", idx)
				}
			}
		}()
	}

	wg.Wait()

	// Total calls should be goroutines * iterations
	expectedCalls := uint64(goroutines * iterations)
	if pool.current != expectedCalls {
		t.Errorf("Expected %d calls, got %d", expectedCalls, pool.current)
	}
}

// TestGetNextPeerRoundRobin verifies round-robin behavior
func TestGetNextPeerRoundRobin(t *testing.T) {
	backends := []*backend.Backend{
		{URL: nil, ReverseProxy: nil},
		{URL: nil, ReverseProxy: nil},
		{URL: nil, ReverseProxy: nil},
	}
	pool := NewServerPool(backends)

	// Call GetNextPeer multiple times and verify rotation
	seen := make(map[*backend.Backend]int)
	for i := 0; i < 9; i++ {
		peer := pool.GetNextPeer()
		if peer == nil {
			t.Fatal("GetNextPeer returned nil")
		}
		seen[peer]++
	}

	// Each backend should be selected 3 times (9 calls / 3 backends)
	for _, count := range seen {
		if count != 3 {
			t.Errorf("Expected each backend to be selected 3 times, got %d", count)
		}
	}
}

// TestGetNextPeerEmptyPool verifies behavior with no backends
func TestGetNextPeerEmptyPool(t *testing.T) {
	pool := NewServerPool([]*backend.Backend{})
	peer := pool.GetNextPeer()
	if peer != nil {
		t.Error("Expected nil for empty pool, got a backend")
	}
}
