.PHONY: build test lint fmt tidy clean run check
.PHONY: docker-build docker-up docker-down docker-logs docker-logs-traefik docker-restart docker-test
.PHONY: docker-up-prod docker-down-prod docker-up-multi docker-down-multi
.PHONY: docker-logs-api docker-logs-auth docker-logs-prom docker-logs-grafana
.PHONY: docker-up-monitoring docker-down-monitoring monitoring-test docker-prune

# Application name
APP_NAME := api
AUTH_NAME := auth
CMD_PATH := cmd/api

# Docker
DOCKER_COMPOSE := docker compose
DOCKER_COMPOSE_FILE := -f docker-compose.yml
DOCKER_COMPOSE_PROD := -f docker-compose.yml -f docker-compose.prod.yml
DOCKER_COMPOSE_MULTI := -f docker-compose.yml -f docker-compose.multi.yml
DOCKER_COMPOSE_MONITORING := -f docker-compose.monitoring.yml

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOFMT := $(GOCMD) fmt
GOMOD := $(GOCMD) mod

# ============================================
# Go targets
# ============================================

# Build the application
build:
	$(GOBUILD) -o bin/$(APP_NAME) ./$(CMD_PATH)

# Build auth service
build-auth:
	$(GOBUILD) -o bin/$(AUTH_NAME) ./cmd/auth

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
# Docker targets - Basic
# ============================================

# Build Docker images
docker-build:
	docker build -t $(APP_NAME):latest .
	docker build -t $(AUTH_NAME):latest . -f Dockerfile.auth

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

# View logs - Auth only
docker-logs-auth:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_FILE) logs -f auth

# Restart all services
docker-restart: docker-down docker-up

# Test Traefik routing
docker-test:
	@echo "Testing Traefik routing..."
	@echo "=== API Health ==="
	@curl -s http://localhost/health | jq .
	@echo "=== Auth Health ==="
	@curl -s http://localhost/auth/health | jq .
	@echo "=== API Users ==="
	@curl -s http://localhost/api/v1/users | jq .
	@echo "=== Traefik Dashboard ==="
	@curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" http://localhost:8080/api/overview

# Docker system prune
docker-prune:
	docker system prune -f

# ============================================
# Docker targets - Multi-service
# ============================================

# Start multi-service (API + Auth)
docker-up-multi:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_MULTI) up -d --build

# Stop multi-service
docker-down-multi:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_MULTI) down

# ============================================
# Docker targets - Monitoring
# ============================================

# Start monitoring stack (Prometheus + Grafana)
docker-up-monitoring:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_MONITORING) up -d

# Stop monitoring stack
docker-down-monitoring:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_MONITORING) down

# View Prometheus logs
docker-logs-prom:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_MONITORING) logs -f prometheus

# View Grafana logs
docker-logs-grafana:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_MONITORING) logs -f grafana

# Test monitoring
monitoring-test:
	@echo "Testing Monitoring Stack..."
	@echo "=== Prometheus ==="
	@curl -s http://localhost:9090/-/healthy | jq .
	@echo "=== Grafana ==="
	@curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" http://localhost:3000/login
	@echo "=== Traefik Metrics ==="
	@curl -s http://localhost:8080/metrics | head -20

# ============================================
# Help
# ============================================
help:
	@echo "Available targets:"
	@echo ""
	@echo "Go targets:"
	@echo "  build           - Build the application"
	@echo "  build-auth      - Build auth service"
	@echo "  run             - Run the application"
	@echo "  test            - Run tests with coverage"
	@echo "  fmt             - Format code"
	@echo "  tidy            - Tidy dependencies"
	@echo "  lint            - Run linter"
	@echo "  clean           - Clean build artifacts"
	@echo "  check           - Run all quality checks"
	@echo ""
	@echo "Docker (Basic):"
	@echo "  docker-build    - Build Docker images"
	@echo "  docker-up       - Start services (dev)"
	@echo "  docker-up-prod  - Start services (prod)"
	@echo "  docker-down     - Stop services"
	@echo "  docker-logs     - View all logs"
	@echo "  docker-test     - Test routing"
	@echo ""
	@echo "Docker (Multi-service):"
	@echo "  docker-up-multi    - Start API + Auth services"
	@echo "  docker-down-multi  - Stop multi-service"
	@echo "  docker-logs-api   - API logs"
	@echo "  docker-logs-auth   - Auth logs"
	@echo ""
	@echo "Docker (Monitoring):"
	@echo "  docker-up-monitoring  - Start Prometheus + Grafana"
	@echo "  docker-down-monitoring - Stop monitoring"
	@echo "  monitoring-test      - Test monitoring stack"
	@echo ""
	@echo "Other:"
	@echo "  docker-prune    - Clean up Docker"
	@echo "  help            - Show this help"
