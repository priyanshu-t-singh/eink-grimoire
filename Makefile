BINARY_NAME=le-grimoire
BUILD_DIR=bin
RELEASE_DIR=release
MAIN_FILE=main.go
DOCKER_USERNAME=priyanshu9943
DOCKER_PLATFORM ?= linux/amd64,linux/arm64
VERSION ?= 1.0.0
DOCKER_IMAGE=$(BINARY_NAME)
OUTPUT_TAR      ?= $(DOCKER_IMAGE)-$(VERSION).tar
FULL_IMAGE_NAME := $(DOCKER_USERNAME)/$(DOCKER_IMAGE)

LDFLAGS := -w -s -X 'le-grimoire/internal/constants.Version=$(VERSION)'

# $(1)=GOOS $(2)=GOARCH $(3)=output suffix (e.g. .exe or empty)
define cross_build
	@echo "Building for $(1)/$(2)..."
	GOOS=$(1) GOARCH=$(2) go build \
		-ldflags="$(LDFLAGS)" \
		-o $(RELEASE_DIR)/$(VERSION)/$(BINARY_NAME)-$(VERSION)-$(1)-$(2)$(3) $(MAIN_FILE)
endef

.PHONY: build run clean dev docker-build docker-run docker-save build-linux build-windows build-mac build-all register-device fmt vet lint

dev:
	go run $(MAIN_FILE)

build:
	go build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)

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
	docker buildx build \
		--platform $(DOCKER_PLATFORM) \
		-t $(FULL_IMAGE_NAME):$(VERSION) \
		-t $(FULL_IMAGE_NAME):latest \
		--push .

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
	$(call cross_build,linux,amd64)
	$(call cross_build,linux,arm64)

build-windows:
	$(call cross_build,windows,amd64,.exe)
	$(call cross_build,windows,arm64,.exe)

build-mac:
	$(call cross_build,darwin,amd64)
	$(call cross_build,darwin,arm64)

build-all: build-linux build-windows build-mac

register-device:
	go run cmd/register-device/main.go -id "$(ID)" -key "$(KEY)" -db ".db/database.sqlite"
