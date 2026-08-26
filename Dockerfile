# Stage 1: Build the binary with CGO enabled
FROM golang:1.27-alpine AS builder

# Install build tools required by go-sqlite3 (CGo)
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
    -ldflags="-w -s -X 'le-grimoire/internal/constants.Host=${HOST}'" \
    -o server main.go

# Stage 2: Minimal runtime image with Chromium & dependencies
FROM alpine:latest

# Install Chromium, fonts, and CA certificates for chromedp
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont

# Ensure chromedp detects Chromium path in Alpine
ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/lib/chromium/

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
