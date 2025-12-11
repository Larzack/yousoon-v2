FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy shared module first
COPY shared/ ./shared/

# Copy service module
COPY notification-service/go.mod ./notification-service/
COPY notification-service/go.sum* ./notification-service/
COPY notification-service/ ./notification-service/
WORKDIR /app/notification-service
RUN go mod tidy && go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /main ./cmd/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates, timezone data and wget for healthcheck
RUN apk --no-cache add ca-certificates tzdata wget

# Copy binary from builder
COPY --from=builder /main .

# Expose port
EXPOSE 4000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:4000/health || exit 1

# Run
CMD ["./main"]
