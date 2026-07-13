.PHONY: build run test clean docker docker-up docker-down mod tidy

APP_NAME := gin-layout
BUILD_DIR := ./bin
GO := go

build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

run: build
	$(BUILD_DIR)/$(APP_NAME) -c etc/config.toml

dev:
	$(GO) run ./cmd/server -c etc/config.toml

test:
	$(GO) test -v -race -cover ./...

test-cover:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	$(GO) clean

mod:
	$(GO) mod download
	$(GO) mod tidy

docker:
	docker build -t $(APP_NAME):latest -f deploy/Dockerfile .

docker-up:
	docker-compose -f deploy/docker-compose.yml up -d

docker-down:
	docker-compose -f deploy/docker-compose.yml down

docker-build-up:
	docker build -t $(APP_NAME):latest -f deploy/Dockerfile .
	docker-compose -f deploy/docker-compose.yml up -d

lint:
	golangci-lint run ./...

vet:
	$(GO) vet ./...

swag:
	swag init -g cmd/server/main.go -o docs

help:
	@echo "Available targets:"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Run without building"
	@echo "  test          - Run tests"
	@echo "  test-cover    - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  mod           - Download and tidy modules"
	@echo "  docker        - Build Docker image"
	@echo "  docker-up     - Start Docker containers"
	@echo "  docker-down   - Stop Docker containers"
	@echo "  lint          - Run linter"
	@echo "  vet           - Run go vet"
	@echo "  swag          - Generate Swagger docs"
