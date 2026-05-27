.PHONY: build test lint fmt tidy clean run

# Application name
APP_NAME := api
CMD_PATH := cmd/api

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOFMT := $(GOCMD) fmt
GOMOD := $(GOCMD) mod

# Build the application
build:
	$(GOBUILD) -o bin/$(APP_NAME) ./$(CMD_PATH)

# Run the application
run:
	$(GOCMD) run ./$(CMD_PATH)

# Run tests
test:
	$(GOTEST) -race -cover ./...

# Run tests with verbose output
test-v:
	$(GOTEST) -v -race -cover ./...

# Format code
fmt:
	$(GOFMT) ./...

# Tidy dependencies
tidy:
	$(GOMOD) tidy

# Lint code
lint:
	golangci-lint run ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Build for Linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(APP_NAME) ./$(CMD_PATH)

# Run all quality checks
check: fmt tidy lint test

# Docker build
docker-build:
	docker build -t $(APP_NAME):latest .

# Help
help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  test       - Run tests with coverage"
	@echo "  test-v     - Run tests with verbose output"
	@echo "  fmt        - Format code"
	@echo "  tidy       - Tidy dependencies"
	@echo "  lint       - Run linter"
	@echo "  clean      - Clean build artifacts"
	@echo "  check      - Run all quality checks"
	@echo "  help       - Show this help"
