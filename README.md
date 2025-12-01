# Load Balancer (Golang)

Overview
-
This repository contains a high-performance, thread-safe HTTP load balancer written in Go. The design follows a modular architecture separating concerns into Listener, Pool, Strategy, and Health components:

- Listener: Accepts incoming HTTP requests and forwards them to the ServerPool according to a Strategy.
- Pool: Maintains a set of backend servers (peers) and exposes selection methods.
- Strategy: Implements load-balancing algorithms (e.g., Round Robin). The Strategy Pattern allows swapping algorithms easily.
- Health: Actively checks backend liveness and updates the Pool to avoid routing to dead backends.

Concurrency & Safety
-
Request handling is done in Goroutines. Shared state (peer alive flags, current index for round-robin) uses `sync.RWMutex` and atomic operations to ensure thread safety.

Project Layout
-
```
/cmd/lb           # entrypoint (main)
/config           # configuration files (yaml)
/internal
  /backend        # backend server wrapper and reverse proxy
  /strategy       # load balancing algorithms (round robin)
  /health         # active health checking routines

README.md
go.mod
```

Execution
-
This project will be runnable as a single binary. The `cmd/lb` package will contain the `main` that loads configuration, starts health checks, and launches an HTTP listener.

Phases
-
- Phase 0: Scaffolding (repo init, README, config, .gitignore, directories)
- Phase 1: Basic reverse proxy and forwarding
- Phase 2: Round robin & concurrency safety
- Phase 3: Active health checks
- Phase 4: Logging, Docker, and Compose for local dev

See commit history for step-by-step incremental implementation.
