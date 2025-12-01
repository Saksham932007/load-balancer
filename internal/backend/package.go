package backend

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

// Backend represents a single backend server with its reverse proxy.
type Backend struct {
	URL          *url.URL
	ReverseProxy *httputil.ReverseProxy
	Alive        bool
	mux          sync.RWMutex
}

// NewBackend creates a new Backend instance with a configured reverse proxy.
func NewBackend(urlStr string) (*Backend, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(u)

	return &Backend{
		URL:          u,
		ReverseProxy: proxy,
		Alive:        true, // Assume alive initially
	}, nil
}

// IsAlive returns the alive status of the backend in a thread-safe manner.
func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Alive
}

// SetAlive sets the alive status of the backend in a thread-safe manner.
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Alive = alive
}
