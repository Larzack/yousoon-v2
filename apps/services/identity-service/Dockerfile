# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy shared module first (for caching)
COPY services/shared/go.mod services/shared/go.sum ./services/shared/
WORKDIR /app/services/shared
RUN go mod download

# Copy identity service module
WORKDIR /app
COPY services/identity-service/go.mod services/identity-service/go.sum ./services/identity-service/
WORKDIR /app/services/identity-service
RUN go mod download

# Copy source code
WORKDIR /app
COPY services/shared ./services/shared
COPY services/identity-service ./services/identity-service

# Build the binary
WORKDIR /app/services/identity-service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /identity-service ./cmd

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /identity-service .

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the binary
ENTRYPOINT ["./identity-service"]
