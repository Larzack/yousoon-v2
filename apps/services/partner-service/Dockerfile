# Partner Service Dockerfile
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy shared module first
COPY shared/ ./shared/

# Copy service module
COPY partner-service/go.mod ./partner-service/
COPY partner-service/go.sum* ./partner-service/
COPY partner-service/ ./partner-service/
WORKDIR /app/partner-service
RUN go mod tidy && go mod download

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
