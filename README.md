# Load Balancer (Golang)

A high-performance, thread-safe HTTP load balancer written in Go with active health checking, round-robin load balancing, and structured JSON logging.

## Overview

This repository contains a production-ready load balancer with the following features:
- **Round-robin load balancing** with thread-safe atomic operations
- **Active health checking** to detect and skip unhealthy backends
- **Automatic retry logic** when backends fail
- **Structured JSON logging** with request latency and client IP tracking
- **Goroutine-based concurrency** for high-throughput request handling
- **Docker support** with multi-stage builds and distroless runtime

## Architecture

The design follows a modular architecture separating concerns into:

- **Listener**: Accepts incoming HTTP requests and forwards them to the ServerPool according to a Strategy.
- **Pool**: Maintains a set of backend servers (peers) and exposes selection methods.
- **Strategy**: Implements load-balancing algorithms (e.g., Round Robin). The Strategy Pattern allows swapping algorithms easily.
- **Health**: Actively checks backend liveness and updates the Pool to avoid routing to dead backends.

### Concurrency & Safety

Request handling is done in Goroutines. Shared state (peer alive flags, current index for round-robin) uses `sync.RWMutex` and atomic operations to ensure thread safety.

## Project Layout

```
/cmd/lb           # entrypoint (main)
/config           # configuration files (yaml)
/internal
  /backend        # backend server wrapper and reverse proxy
  /strategy       # load balancing algorithms (round robin)
  /health         # active health checking routines
/test-data        # sample data for backend servers
Dockerfile        # multi-stage Docker build
docker-compose.yml # local development setup
README.md
go.mod
```

## How to Run

### Option 1: Using Docker Compose (Recommended)

This will start the load balancer and 3 Python backend servers:

```bash
docker-compose up --build
```

Access the load balancer at `http://localhost:8080`

### Option 2: Build and Run Locally

**Prerequisites**: Go 1.20+

1. Build the binary:
```bash
go build -o load-balancer ./cmd/lb
```

2. Start some backend servers (in separate terminals):
```bash
# Terminal 1
python3 -m http.server 8001

# Terminal 2
python3 -m http.server 8002

# Terminal 3
python3 -m http.server 8003
```

3. Run the load balancer:
```bash
./load-balancer
```

4. Test it:
```bash
curl http://localhost:8080
```

### Option 3: Run Tests

```bash
go test ./...
```

## Configuration

Edit `config/config.yaml` to customize:
- Backend URLs
- Health check interval
- Timeout settings
- Listen address

Or modify the `backendURLs` slice in `cmd/lb/main.go` for quick changes.

## Features

### Round-Robin Load Balancing
Uses atomic operations (`atomic.AddUint64`) to safely distribute requests across backends in a round-robin fashion, even under high concurrent load.

### Active Health Checks
Every 10 seconds, the load balancer pings each backend. Unhealthy backends are automatically removed from rotation until they recover.

### Retry Logic
If a backend returns a 503, the load balancer automatically retries with the next available backend (up to N attempts, where N = number of backends).

### Structured Logging
All logs are in JSON format with structured fields:
```json
{
  "time": "2025-12-01T10:30:45Z",
  "level": "INFO",
  "msg": "Request completed",
  "client_ip": "127.0.0.1:54321",
  "method": "GET",
  "path": "/",
  "backend": "http://localhost:8001",
  "status": 200,
  "latency_ms": 12
}
```

## Development Phases

This project was built incrementally with 23+ commits following Conventional Commits:

- **Phase 0**: Scaffolding (repo init, README, config, .gitignore, directories)
- **Phase 1**: Basic reverse proxy and forwarding
- **Phase 2**: Round robin & concurrency safety
- **Phase 3**: Active health checks
- **Phase 4**: Logging, Docker, and Compose for local dev

See commit history for step-by-step incremental implementation.

## License

MIT

