package backend

import (
	"net/http/httputil"
	"net/url"
)

// Backend represents a single backend server with its reverse proxy.
type Backend struct {
	URL          *url.URL
	ReverseProxy *httputil.ReverseProxy
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
	}, nil
}
