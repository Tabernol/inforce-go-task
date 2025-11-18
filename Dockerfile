# ---- builder ----
FROM golang:1.23.10 AS builder

# Set build args (overrideable)
ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64

WORKDIR /src

# Copy go.mod first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy rest
COPY . .

# Build static binary
ENV CGO_ENABLED=${CGO_ENABLED}
ENV GOOS=${GOOS}
ENV GOARCH=${GOARCH}

# Adjust the package path if your main is in ./cmd (we assume ./cmd)
RUN go build -ldflags="-s -w" -o /out/rarible-service ./cmd

# ---- runtime ----
FROM alpine:3.18 AS runtime

# Add CA certs and curl for healthcheck
RUN apk add --no-cache ca-certificates curl && update-ca-certificates

# Create non-root user (UID/GID 65532 is common for nonroot)
RUN addgroup -g 65532 nonroot && adduser -u 65532 -G nonroot -D -H -s /sbin/nologin nonroot

WORKDIR /app

# Copy binary from builder
COPY --from=builder /out/rarible-service /usr/local/bin/rarible-service

# Ensure binary is executable
RUN chmod +x /usr/local/bin/rarible-service

# Use non-root user
USER nonroot

EXPOSE 8080

ENV SERVER_PORT=8080

# Healthcheck using curl against /health
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD curl -f http://127.0.0.1:8080/health || exit 1

ENTRYPOINT ["/usr/local/bin/rarible-service"]