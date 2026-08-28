# Stage 1: Build
FROM golang:1.27-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .

ARG HOST=0.0.0.0

# Add -ldflags to your build step
RUN CGO_ENABLED=1 GOOS=linux go build \
    -trimpath \
    -ldflags="-w -s -X 'le-grimoire/internal/constants.Host=${HOST}'" \
    -o server main.go

# Build the register-device binary
RUN CGO_ENABLED=1 GOOS=linux go build \
    -trimpath \
    -o register-device cmd/register-device/main.go

# Stage 2: Runtime image
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Database mount here
RUN mkdir -p /app/.db && \
    chown -R appuser:appgroup /app

# Copy both compiled binaries from the builder stage
COPY --from=builder /app/server /usr/local/bin/server
COPY --from=builder /app/register-device /usr/local/bin/register-device

ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8080

EXPOSE 8080

# Run as non-root user
USER appuser

ENTRYPOINT ["server"]
