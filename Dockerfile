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
    -ldflags="-w -s -X 'le-grimoire/internal/constants.Host=${HOST}'" \
    -o server main.go

# Stage 2: Runtime image
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080

ENTRYPOINT ["/app/server"]
