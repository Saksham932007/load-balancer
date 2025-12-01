# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies (if go.sum exists)
RUN go mod download || true

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o load-balancer ./cmd/lb

# Runtime stage (distroless)
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copy binary from builder
COPY --from=builder /app/load-balancer .

# Use non-root user
USER nonroot:nonroot

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/load-balancer"]
