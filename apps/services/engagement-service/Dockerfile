# Build stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY shared/ ./shared/
COPY engagement-service/go.mod ./engagement-service/
COPY engagement-service/go.sum* ./engagement-service/
COPY engagement-service/ ./engagement-service/
WORKDIR /app/engagement-service
RUN go mod tidy && go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /engagement-service ./cmd/main.go

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app

COPY --from=builder /engagement-service .

RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8084
EXPOSE 9094

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:4000/health || exit 1

ENTRYPOINT ["./engagement-service"]
