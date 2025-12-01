package health

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Saksham932007/load-balancer/internal/backend"
)

const (
	healthCheckInterval = 10 * time.Second
	healthCheckTimeout  = 5 * time.Second
)

// CheckBackend performs a health check on a single backend.
func CheckBackend(b *backend.Backend) {
	client := &http.Client{
		Timeout: healthCheckTimeout,
	}

	resp, err := client.Get(b.URL.String())
	if err != nil {
		slog.Debug("Health check failed", "backend", b.URL.String(), "error", err)
		b.SetAlive(false)
		return
	}
	defer resp.Body.Close()

	// Consider 2xx and 3xx status codes as healthy
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		if !b.IsAlive() {
			slog.Info("Backend recovered", "backend", b.URL.String())
		}
		b.SetAlive(true)
	} else {
		slog.Warn("Backend unhealthy", "backend", b.URL.String(), "status", resp.StatusCode)
		b.SetAlive(false)
	}
}

// StartHealthCheck launches a background goroutine that periodically checks backend health.
func StartHealthCheck(backends []*backend.Backend) {
	ticker := time.NewTicker(healthCheckInterval)
	go func() {
		for range ticker.C {
			slog.Debug("Running periodic health checks")
			for _, b := range backends {
				go CheckBackend(b)
			}
		}
	}()

	// Run initial health check immediately
	slog.Info("Running initial health checks")
	for _, b := range backends {
		go CheckBackend(b)
	}
}
