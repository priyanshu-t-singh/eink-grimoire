# Stage 1: Build binary using Debian/glibc
FROM golang:1.27-bookworm AS builder

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build
COPY . .

ARG HOST=0.0.0.0

# Add -ldflags to your build step
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-w -s -X 'le-grimoire/internal/constants.Host=${HOST}'" \
    -o server main.go

# Stage 2: Runtime image
FROM chromedp/headless-shell:latest

# Default path for this specific container image
ENV CHROME_PATH=/headless-shell/headless-shell

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080

ENTRYPOINT ["/app/server"]
