package strategy

import (
	"sync/atomic"

	"github.com/Saksham932007/load-balancer/internal/backend"
)

// ServerPool holds multiple backend servers and manages load balancing.
type ServerPool struct {
	backends []*backend.Backend
	current  uint64
}

// NewServerPool creates a new ServerPool with the given backends.
func NewServerPool(backends []*backend.Backend) *ServerPool {
	return &ServerPool{
		backends: backends,
		current:  0,
	}
}

// NextIndex atomically increments and returns the next index for round-robin.
// This method is thread-safe and prevents race conditions.
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

// GetNextPeer returns the next available backend using round-robin algorithm.
// It skips backends that are marked as not alive.
func (s *ServerPool) GetNextPeer() *backend.Backend {
	if len(s.backends) == 0 {
		return nil
	}

	// Try each backend once (starting from next index)
	for i := 0; i < len(s.backends); i++ {
		nextIdx := s.NextIndex()
		backend := s.backends[nextIdx]

		// Skip dead backends
		if backend.IsAlive() {
			return backend
		}
	}

	// All backends are dead
	return nil
}
