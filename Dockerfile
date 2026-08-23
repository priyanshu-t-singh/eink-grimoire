# Stage 1: Build the binary
FROM golang:1.27-alpine AS builder
WORKDIR /app

# Download dependencies first (cache them)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and compile
COPY . .

# CGO_ENABLED=0 ensures a static binary with no external dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Stage 2: Create the final minimal image
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .

# Document that this server listens on port 8080
EXPOSE 8080

CMD ["./server"]
