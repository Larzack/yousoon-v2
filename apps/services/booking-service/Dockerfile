# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# ===========================================
# STEP 1: Copy go.mod/go.sum for dependency caching
# ===========================================
# Copy shared module dependencies first
COPY shared/go.mod shared/go.sum* ./shared/

# Copy service dependencies
COPY booking-service/go.mod booking-service/go.sum* ./booking-service/

# Download dependencies (cached if go.mod/go.sum unchanged)
WORKDIR /app/booking-service
RUN go mod download

# ===========================================
# STEP 2: Copy source code and build
# ===========================================
WORKDIR /app
COPY shared/ ./shared/
COPY booking-service/ ./booking-service/

WORKDIR /app/booking-service

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /booking-service \
    ./cmd/main.go

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app

COPY --from=builder /booking-service .

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 4000
EXPOSE 9093

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:4000/health || exit 1

ENTRYPOINT ["./booking-service"]
