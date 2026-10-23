# syntax=docker/dockerfile:1

# ==========================================
# 1. Build Stage
# ==========================================
ARG GO_VERSION=alpine
FROM golang:${GO_VERSION} AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache ca-certificates tzdata git

# Cache Go modules layer
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Build statically compiled binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /bin/party-service \
    ./cmd/service

# ==========================================
# 2. Production Runtime Stage
# ==========================================
FROM alpine:3.21 AS runner

# Install certificates, timezone data, and wget for healthchecks
RUN apk --no-cache add ca-certificates tzdata wget \
    && addgroup -S appgroup -g 10001 \
    && adduser -S appuser -u 10001 -G appgroup

WORKDIR /app

# Copy binary from build stage
COPY --from=builder /bin/party-service /bin/party-service

# Run as non-root user
USER appuser:appgroup

# Expose HTTP port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=15s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the compiled binary
ENTRYPOINT ["/bin/party-service"]
