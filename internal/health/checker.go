package health

import (
	"log"
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
		log.Printf("Health check failed for %s: %v", b.URL, err)
		b.SetAlive(false)
		return
	}
	defer resp.Body.Close()

	// Consider 2xx and 3xx status codes as healthy
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		if !b.IsAlive() {
			log.Printf("Backend %s is now alive", b.URL)
		}
		b.SetAlive(true)
	} else {
		log.Printf("Backend %s returned status %d", b.URL, resp.StatusCode)
		b.SetAlive(false)
	}
}

// StartHealthCheck launches a background goroutine that periodically checks backend health.
func StartHealthCheck(backends []*backend.Backend) {
	ticker := time.NewTicker(healthCheckInterval)
	go func() {
		for range ticker.C {
			log.Println("Running health checks...")
			for _, b := range backends {
				go CheckBackend(b)
			}
		}
	}()

	// Run initial health check immediately
	log.Println("Running initial health checks...")
	for _, b := range backends {
		go CheckBackend(b)
	}
}
