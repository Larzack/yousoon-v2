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
COPY engagement-service/go.mod engagement-service/go.sum* ./engagement-service/

# Download dependencies (cached if go.mod/go.sum unchanged)
WORKDIR /app/engagement-service
RUN go mod download

# ===========================================
# STEP 2: Copy source code and build
# ===========================================
WORKDIR /app
COPY shared/ ./shared/
COPY engagement-service/ ./engagement-service/

WORKDIR /app/engagement-service

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /engagement-service \
    ./cmd/main.go

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app

COPY --from=builder /engagement-service .

RUN adduser -D -g '' appuser
USER appuser

EXPOSE 4000
EXPOSE 9094

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:4000/health || exit 1

ENTRYPOINT ["./engagement-service"]
