BINARY_NAME=le-grimoire
BUILD_DIR=bin
MAIN_FILE=main.go

DOCKER_USERNAME=priyanshu9943
DOCKER_IMAGE=le-grimoire
DOCKER_TAG=latest

OUTPUT_TAR      ?= $(DOCKER_IMAGE)-$(DOCKER_TAG).tar
FULL_IMAGE_NAME := $(DOCKER_USERNAME)/$(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: build run clean dev docker-build docker-run docker-save build-linux build-windows build-mac build-all

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
	if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found. Please install it first."; \
		echo "You can install it by running: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi
	@echo "Linting..."
	golangci-lint run


docker-build:
	@echo "Building Docker image..."
	docker build -t $(FULL_IMAGE_NAME) .

docker-run:
	@echo "Running Docker container..."
	docker run \
		-p 8080:8080 \
		--name $(BINARY_NAME) \
		--rm \
		-v $(PWD)/.db:/app/.db \
		$(FULL_IMAGE_NAME)

docker-save:
	@echo "Saving $(FULL_IMAGE_NAME) to $(OUTPUT_TAR)..."
	docker save -o $(OUTPUT_TAR) $(FULL_IMAGE_NAME)
	@echo "Saved successfully: $(OUTPUT_TAR)"

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
