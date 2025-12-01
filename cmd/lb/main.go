package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Saksham932007/load-balancer/internal/backend"
	"github.com/Saksham932007/load-balancer/internal/health"
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
	logger     *slog.Logger
)

func main() {
	// Initialize structured JSON logger
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Initialize backends from URL list
	var backends []*backend.Backend
	for _, urlStr := range backendURLs {
		b, err := backend.NewBackend(urlStr)
		if err != nil {
			logger.Error("Failed to create backend", "url", urlStr, "error", err)
			os.Exit(1)
		}
		backends = append(backends, b)
		logger.Info("Added backend", "url", urlStr)
	}

	// Initialize server pool
	serverPool = strategy.NewServerPool(backends)

	// Start health checking in background
	health.StartHealthCheck(backends)

	server := http.Server{
		Addr:    listenAddr,
		Handler: http.HandlerFunc(handleRequest),
	}

	logger.Info("Starting load balancer", "address", listenAddr, "backends", len(backends))
	if err := server.ListenAndServe(); err != nil {
		logger.Error("Server failed", "error", err)
		os.Exit(1)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	clientIP := r.RemoteAddr

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
			logger.Warn("No backends available", "client_ip", clientIP)
			http.Error(w, "No backends available", http.StatusServiceUnavailable)
			return
		}

		// Use a custom response writer to detect failures
		recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		peer.ReverseProxy.ServeHTTP(recorder, r)

		latency := time.Since(start)

		// If backend responded successfully (not 503), we're done
		if recorder.statusCode != http.StatusServiceUnavailable {
			logger.Info("Request completed",
				"client_ip", clientIP,
				"method", r.Method,
				"path", r.URL.Path,
				"backend", peer.URL.String(),
				"status", recorder.statusCode,
				"latency_ms", latency.Milliseconds(),
			)
			return
		}

		// Backend failed, try next one
		attempts++
		logger.Warn("Backend failed, retrying",
			"backend", peer.URL.String(),
			"attempt", attempts,
			"max_attempts", maxAttempts,
			"client_ip", clientIP,
		)
	}

	// All backends failed
	logger.Error("All backends unavailable", "client_ip", clientIP)
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
