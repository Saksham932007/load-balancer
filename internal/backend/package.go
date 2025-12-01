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
