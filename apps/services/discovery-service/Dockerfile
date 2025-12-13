# Discovery Service Dockerfile
# Multi-stage build for minimal image size

# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# ===========================================
# STEP 1: Copy go.mod/go.sum for dependency caching
# ===========================================
# Copy shared module dependencies first
COPY shared/go.mod shared/go.sum* ./shared/

# Copy service dependencies
COPY discovery-service/go.mod discovery-service/go.sum* ./discovery-service/

# Download dependencies (cached if go.mod/go.sum unchanged)
WORKDIR /app/discovery-service
RUN go mod download

# ===========================================
# STEP 2: Copy source code and build
# ===========================================
WORKDIR /app
COPY shared/ ./shared/
COPY discovery-service/ ./discovery-service/

WORKDIR /app/discovery-service

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION:-dev}" \
    -o /discovery-service \
    ./cmd/main.go

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /discovery-service .

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Set timezone
ENV TZ=Europe/Paris

# Use non-root user
USER appuser

# Expose ports
EXPOSE 4000 9090 9100

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:4000/health || exit 1

# Run the application
ENTRYPOINT ["/app/discovery-service"]
