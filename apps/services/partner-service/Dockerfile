# Partner Service Dockerfile
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# ===========================================
# STEP 1: Copy go.mod/go.sum for dependency caching
# ===========================================
# Copy shared module dependencies first
COPY shared/go.mod shared/go.sum* ./shared/

# Copy service dependencies
COPY partner-service/go.mod partner-service/go.sum* ./partner-service/

# Download dependencies (cached if go.mod/go.sum unchanged)
WORKDIR /app/partner-service
RUN go mod download

# ===========================================
# STEP 2: Copy source code and build
# ===========================================
WORKDIR /app
COPY shared/ ./shared/
COPY partner-service/ ./partner-service/

WORKDIR /app/partner-service

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /partner-service \
    ./cmd/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata wget
RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /partner-service .

USER appuser

EXPOSE 4000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:4000/health || exit 1

CMD ["./partner-service"]
