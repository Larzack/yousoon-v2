# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy shared module first
COPY shared/ ./shared/

# Copy service module
COPY booking-service/go.mod ./booking-service/
COPY booking-service/go.sum* ./booking-service/
COPY booking-service/ ./booking-service/
WORKDIR /app/booking-service
RUN go mod tidy && go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /booking-service ./cmd/main.go

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
