# Multi-stage build for smaller final image

# Stage 1: Builder
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-X 'github.com/PerHac13/vaultra/cmd/vaultra.Version=v0.1.0' \
              -X 'github.com/PerHac13/vaultra/cmd/vaultra.BuildTime=$(date -u '+%Y-%m-%d_%H:%M:%S')'" \
    -o vaultra ./cmd/vaultra

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /root

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata postgresql-client mysql-client mongodb-tools

# Copy binary from builder
COPY --from=builder /app/vaultra .

# Copy config templates
COPY configs/ ./configs/

# Create a non-root user
RUN addgroup -g 1000 vaultra && \
    adduser -D -u 1000 -G vaultra vaultra && \
    chown -R vaultra:vaultra /root

USER vaultra

ENTRYPOINT ["./vaultra"]
CMD ["--help"]
