.PHONY: build clean run test test-coverage lint swagger docker-build docker-run docker-compose-up docker-compose-down help migrate migrate-up migrate-down

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet
GOLINT=golangci-lint

# Application parameters
APP_NAME=web3-edu-core
MAIN_PATH=./cmd/api
BINARY_NAME=api-server
BUILD_DIR=./bin

# Docker parameters
DOCKER_IMAGE_NAME=web3-edu-api
DOCKER_IMAGE_TAG=latest

# Default target
.DEFAULT_GOAL := help

# Help target
help:
	@echo "Available commands:"
	@echo "  make build              - Build the application"
	@echo "  make clean              - Remove build artifacts"
	@echo "  make run                - Run the application"
	@echo "  make test               - Run tests"
	@echo "  make test-coverage      - Run tests with coverage"
	@echo "  make lint               - Run linters"
	@echo "  make swagger            - Generate Swagger documentation"
	@echo "  make docker-build       - Build Docker image"
	@echo "  make docker-run         - Run Docker container"
	@echo "  make docker-compose-up  - Start all services with Docker Compose"
	@echo "  make docker-compose-down - Stop all services with Docker Compose"
	@echo "  make migrate            - Run all migrations"
	@echo "  make migrate-up         - Run migrations up"
	@echo "  make migrate-down       - Run migrations down"
	@echo "  make deps               - Download dependencies"
	@echo "  make tidy               - Tidy up the go.mod file"

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete. Binary: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	$(GORUN) $(MAIN_PATH)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linters
lint:
	@echo "Running linters..."
	$(GOFMT) ./...
	$(GOVET) ./...
	@if command -v $(GOLINT) > /dev/null; then \
		$(GOLINT) run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g internal/api/docs.go -o docs; \
	else \
		echo "swag not installed. Run: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8003:8003 --name $(APP_NAME) $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

# Start all services with Docker Compose
docker-compose-up:
	@echo "Starting all services with Docker Compose..."
	docker-compose up -d

# Stop all services with Docker Compose
docker-compose-down:
	@echo "Stopping all services with Docker Compose..."
	docker-compose down

# Run database migrations
migrate:
	@echo "Running all migrations..."
	@for file in migrations/*.up.sql; do \
		echo "Running migration: $$file"; \
		psql -h localhost -p 5432 -U postgres -d web3_edu_db -f $$file; \
	done

# Run migrations up
migrate-up:
	@echo "Running migrations up..."
	@if [ -z "$(version)" ]; then \
		echo "Please specify a version. Example: make migrate-up version=20230101000001"; \
	else \
		echo "Running migration up to version $(version)"; \
		psql -h localhost -p 5432 -U postgres -d web3_edu_db -f migrations/$(version)_*.up.sql; \
	fi

# Run migrations down
migrate-down:
	@echo "Running migrations down..."
	@if [ -z "$(version)" ]; then \
		echo "Please specify a version. Example: make migrate-down version=20230101000001"; \
	else \
		echo "Running migration down to version $(version)"; \
		psql -h localhost -p 5432 -U postgres -d web3_edu_db -f migrations/$(version)_*.down.sql; \
	fi

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOGET) -u ./...
	$(GOGET) -u github.com/swaggo/swag/cmd/swag
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint

# Tidy up the go.mod file
tidy:
	@echo "Tidying up go.mod..."
	$(GOMOD) tidy
