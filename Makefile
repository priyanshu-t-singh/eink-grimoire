BINARY_NAME=le-grimoire
BUILD_DIR=bin
MAIN_FILE=main.go

DOCKER_USERNAME=priyanshu9943
DOCKER_IMAGE=le-grimoire
DOCKER_TAG=latest

.PHONY: build run clean dev docker-build docker-run build-linux build-windows build-mac build-all

dev:
	go run $(MAIN_FILE)

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	rm -rf $(BUILD_DIR)

fmt:
	@echo "Formatting code..."
	go fmt ./...

vet:
	@echo "Vetting code..."
	go vet ./...

lint: ## Run golangci-lint (requires installation)
	@echo "Linting..."
	golangci-lint run


docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	@echo "Running Docker container..."
	docker run \
		-p 8080:8080 \
		--name $(BINARY_NAME) \
		--rm \
		-v $(PWD)/.db:/app/.db \
		$(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(DOCKER_TAG)

build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)

build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)

build-mac:
	@echo "Building for macOS (Intel)..."
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	@echo "Building for macOS (Apple Silicon)..."
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)

build-all: build-linux build-windows build-mac
