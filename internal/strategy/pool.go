package strategy

import (
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
