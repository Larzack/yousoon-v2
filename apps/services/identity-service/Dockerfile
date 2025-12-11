# Identity Service Dockerfile
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy shared module first
COPY shared/ ./shared/

# Copy service module
COPY identity-service/go.mod ./identity-service/
COPY identity-service/go.sum* ./identity-service/
COPY identity-service/ ./identity-service/
WORKDIR /app/identity-service
RUN go mod tidy && go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /identity-service \
    ./cmd/main.go

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata wget

RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /identity-service .

USER appuser

EXPOSE 4000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:4000/health || exit 1

CMD ["./identity-service"]
