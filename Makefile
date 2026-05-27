.PHONY: build test lint fmt tidy clean run check
.PHONY: docker-build docker-up docker-down docker-logs docker-logs-traefik docker-restart docker-test

# Application name
APP_NAME := api
CMD_PATH := cmd/api

# Docker
DOCKER_COMPOSE := docker compose
DOCKER_COMPOSE_FILE := -f docker-compose.yml
DOCKER_COMPOSE_PROD := -f docker-compose.yml -f docker-compose.prod.yml

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

# ============================================
# Docker targets
# ============================================

# Build Docker images
docker-build:
	docker build -t $(APP_NAME):latest .

# Start all services (development)
docker-up:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_FILE) up -d --build

# Start all services (production)
docker-up-prod:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_PROD) up -d --build

# Stop all services
docker-down:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_FILE) down

# Stop all services (production)
docker-down-prod:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_PROD) down

# View logs - all services
docker-logs:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_FILE) logs -f

# View logs - Traefik only
docker-logs-traefik:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_FILE) logs -f traefik

# View logs - API only
docker-logs-api:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_FILE) logs -f api

# Restart all services
docker-restart: docker-down docker-up

# Test Traefik routing
docker-test:
	@echo "Testing Traefik routing..."
	@echo "=== Health Check ==="
	@curl -s http://localhost/health | jq .
	@echo "=== API Users List ==="
	@curl -s http://localhost/api/v1/users | jq .
	@echo "=== Traefik Dashboard ==="
	@curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" http://localhost:8080/api/overview

# Docker system prune (clean up)
docker-prune:
	docker system prune -f

# Help
help:
	@echo "Available targets:"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application"
	@echo "  test            - Run tests with coverage"
	@echo "  test-v          - Run tests with verbose output"
	@echo "  fmt             - Format code"
	@echo "  tidy            - Tidy dependencies"
	@echo "  lint            - Run linter"
	@echo "  clean           - Clean build artifacts"
	@echo "  check           - Run all quality checks"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-up       - Start services (dev)"
	@echo "  docker-up-prod  - Start services (prod)"
	@echo "  docker-down     - Stop services"
	@echo "  docker-logs     - View all logs"
	@echo "  docker-test     - Test routing"
	@echo "  help            - Show this help"
