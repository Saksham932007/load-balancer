package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	listenAddr = ":8080"
)

func main() {
	server := http.Server{
		Addr:    listenAddr,
		Handler: http.HandlerFunc(handleRequest),
	}

	log.Printf("Starting load balancer on %s", listenAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Placeholder handler - will forward to backend in next commit
	fmt.Fprintf(w, "Load balancer is running\n")
}
