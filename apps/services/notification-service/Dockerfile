# Notification Service Dockerfile
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# ===========================================
# STEP 1: Copy go.mod/go.sum for dependency caching
# ===========================================
# Copy shared module dependencies first
COPY shared/go.mod shared/go.sum* ./shared/

# Copy service dependencies
COPY notification-service/go.mod notification-service/go.sum* ./notification-service/

# Download dependencies (cached if go.mod/go.sum unchanged)
WORKDIR /app/notification-service
RUN go mod download

# ===========================================
# STEP 2: Copy source code and build
# ===========================================
WORKDIR /app
COPY shared/ ./shared/
COPY notification-service/ ./notification-service/

WORKDIR /app/notification-service

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /notification-service \
    ./cmd/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates, timezone data and wget for healthcheck
RUN apk --no-cache add ca-certificates tzdata wget

# Create non-root user
RUN adduser -D -g '' appuser

# Copy binary from builder
COPY --from=builder /notification-service .

USER appuser

# Expose port
EXPOSE 4000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:4000/health || exit 1

# Run
CMD ["./notification-service"]
