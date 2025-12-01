package main

import (
	"log"
	"net/http"

	"github.com/Saksham932007/load-balancer/internal/backend"
)

const (
	listenAddr     = ":8080"
	hardcodedBackend = "http://localhost:8001"
)

var backendInstance *backend.Backend

func main() {
	// Initialize hardcoded backend
	var err error
	backendInstance, err = backend.NewBackend(hardcodedBackend)
	if err != nil {
		log.Fatalf("Failed to create backend: %v", err)
	}

	server := http.Server{
		Addr:    listenAddr,
		Handler: http.HandlerFunc(handleRequest),
	}

	log.Printf("Starting load balancer on %s", listenAddr)
	log.Printf("Forwarding to backend: %s", hardcodedBackend)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Forward request to the hardcoded backend
	backendInstance.ReverseProxy.ServeHTTP(w, r)
}
