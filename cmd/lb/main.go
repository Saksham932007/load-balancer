package main

import (
	"log"
	"net/http"

	"github.com/Saksham932007/load-balancer/internal/backend"
	"github.com/Saksham932007/load-balancer/internal/strategy"
)

const (
	listenAddr = ":8080"
)

var (
	backendURLs = []string{
		"http://localhost:8001",
		"http://localhost:8002",
		"http://localhost:8003",
	}
	serverPool *strategy.ServerPool
)

func main() {
	// Initialize backends from URL list
	var backends []*backend.Backend
	for _, urlStr := range backendURLs {
		b, err := backend.NewBackend(urlStr)
		if err != nil {
			log.Fatalf("Failed to create backend %s: %v", urlStr, err)
		}
		backends = append(backends, b)
		log.Printf("Added backend: %s", urlStr)
	}

	// Initialize server pool
	serverPool = strategy.NewServerPool(backends)

	server := http.Server{
		Addr:    listenAddr,
		Handler: http.HandlerFunc(handleRequest),
	}

	log.Printf("Starting load balancer on %s", listenAddr)
	log.Printf("Load balancing across %d backends", len(backends))
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Add X-Forwarded-For header to track original client IP
	if clientIP := r.Header.Get("X-Real-IP"); clientIP == "" {
		r.Header.Set("X-Forwarded-For", r.RemoteAddr)
	}

	// Retry logic: attempt to forward to available backends
	attempts := 0
	maxAttempts := len(backendURLs)

	for attempts < maxAttempts {
		peer := serverPool.GetNextPeer()
		if peer == nil {
			http.Error(w, "No backends available", http.StatusServiceUnavailable)
			return
		}

		// Use a custom response writer to detect failures
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		peer.ReverseProxy.ServeHTTP(recorder, r)

		// If backend responded successfully (not 503), we're done
		if recorder.statusCode != http.StatusServiceUnavailable {
			return
		}

		// Backend failed, try next one
		attempts++
		log.Printf("Backend %s failed (503), trying next peer (attempt %d/%d)", peer.URL, attempts, maxAttempts)
	}

	// All backends failed
	http.Error(w, "All backends unavailable", http.StatusServiceUnavailable)
}

// responseRecorder wraps http.ResponseWriter to capture the status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
