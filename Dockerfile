# =========================
# Build stage
# =========================
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/bin/api ./cmd/api

# =========================
# Runtime stage
# =========================
FROM alpine:3.21

RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    wget

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /app/bin/api /usr/local/bin/api
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/docs /app/docs

USER app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:8080/api/v1/health || exit 1

ENTRYPOINT ["api"]