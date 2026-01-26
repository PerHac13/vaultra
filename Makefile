.PHONY: help build run test lint clean docker-build docker-run fmt

PROJECT_NAME := vaultra
MAIN_PATH := cmd/vaultra
BINARY_NAME := vaultra
VERSION := v0.1.0-dev
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LD_FLAGS := -ldflags="-X 'github.com/PerHac13/$(PROJECT_NAME)/cmd/$(PROJECT_NAME).Version=$(VERSION)' \
	-X 'github.com/PerHac13/$(PROJECT_NAME)/cmd/$(PROJECT_NAME).BuildTime=$(BUILD_TIME)' \
	-X 'github.com/PerHac13/$(PROJECT_NAME)/cmd/$(PROJECT_NAME).GitCommit=$(GIT_COMMIT)'"

help: ## Show this help message
	@echo "$(PROJECT_NAME) - Makefile commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  make %-20s %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "-> Building $(BINARY_NAME)..."
	go build $(LD_FLAGS) -o ./bin/$(BINARY_NAME) ./$(MAIN_PATH)
	@echo "-> Binary created at ./bin/$(BINARY_NAME)"

run: build ## Build and run the binary
	@echo "-> Running $(BINARY_NAME)..."
	./bin/$(BINARY_NAME) --help

test: ## Run all tests
	@echo "-> Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "-> Coverage report generated: coverage.out"

test-unit: ## Run unit tests only
	@echo "-> Running unit tests..."
	go test -v -race ./test/unit/...

test-integration: ## Run integration tests (requires Docker)
	@echo "-> Running integration tests..."
	go test -v -race -tags integration ./test/integration/...

lint: ## Run linter
	@echo "-> Linting code..."
	golangci-lint run ./...

fmt: ## Format code
	@echo "-> Formatting code..."
	go fmt ./...
	goimports -w .

mod-tidy: ## Tidy go.mod
	@echo "-> Tidying modules..."
	go mod tidy

mod-download: ## Download dependencies
	@echo "->  Downloading dependencies..."
	go mod download

clean: ## Clean build artifacts
	@echo "-> Cleaning..."
	rm -rf ./bin
	rm -f coverage.out
	go clean

docker-build: ## Build Docker image
	@echo "-> Building Docker image..."
	docker build -f build/Dockerfile -t $(PROJECT_NAME):latest .
	docker tag $(PROJECT_NAME):latest $(PROJECT_NAME):$(VERSION)
	@echo "-> Docker image built: $(PROJECT_NAME):latest"

docker-run: docker-build ## Run Docker container
	@echo "-> Running Docker container..."
	docker run --rm -it $(PROJECT_NAME):latest --help

docker-compose-up: ## Start services with docker-compose
	@echo "-> Starting docker-compose services..."
	docker-compose -f deployments/docker-compose.yaml up -d
	@echo "-> Services started"

docker-compose-down: ## Stop docker-compose services
	@echo "-> Stopping docker-compose services..."
	docker-compose -f deployments/docker-compose.yaml down
	@echo "-> Services stopped"

setup-dev: mod-download ## Setup development environment
	@echo "->  Setting up development environment..."
	@command -v golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@command -v goimports >/dev/null 2>&1 || go install golang.org/x/tools/cmd/goimports@latest
	@echo "-> Development environment ready"

ci: mod-tidy fmt lint test ## Run all CI checks (format, lint, test)
	@echo "-> All CI checks passed"

.DEFAULT_GOAL := help
