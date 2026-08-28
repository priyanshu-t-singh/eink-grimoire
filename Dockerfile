# Stage 1: Build
FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS builder

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Copy source code and build
COPY . .

ARG TARGETOS TARGETARCH
ARG HOST=0.0.0.0

# Build the binaries with CGO enabled for cross-compilation
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build \
      -trimpath \
      -ldflags="-w -s -X 'le-grimoire/internal/constants.Host=${HOST}'" \
      -o server main.go && \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build \
      -trimpath \
      -o register-device cmd/register-device/main.go

# Stage 2: Runtime image
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Create directory with correct user permissions upfront
RUN mkdir -p /app/.db && \
    chown -R appuser:appgroup /app

# Copy both compiled binaries from the builder stage to system PATH
COPY --from=builder /app/server /usr/local/bin/server
COPY --from=builder /app/register-device /usr/local/bin/register-device

ENV SERVER_HOST=0.0.0.0
ENV SERVER_PORT=8080

EXPOSE 8080

# Run as non-root user
USER appuser

ENTRYPOINT ["server"]
