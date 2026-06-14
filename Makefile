.PHONY: help up down build clean logs shell test dev godot-dev

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Main commands
up: ## Start the backend server
	@echo "Starting Spheres backend server..."
	docker-compose up --build backend

up-detached: ## Start the backend server in detached mode
	@echo "Starting Spheres backend server in detached mode..."
	docker-compose up --build -d backend

dev: ## Start all services including Godot development container
	@echo "Starting all development services..."
	docker-compose --profile dev up --build

down: ## Stop all services
	@echo "Stopping all services..."
	docker-compose down

# Build commands
build: ## Build the Docker images
	@echo "Building Docker images..."
	docker-compose build

rebuild: ## Rebuild the Docker images without cache
	@echo "Rebuilding Docker images without cache..."
	docker-compose build --no-cache

# Development commands
logs: ## Show logs from all services
	docker-compose logs -f

logs-backend: ## Show logs from backend service only
	docker-compose logs -f backend

shell: ## Open a shell in the backend container
	docker-compose exec backend sh

godot-dev: ## Start Godot development container
	@echo "Starting Godot development environment..."
	docker-compose --profile dev up godot-dev

godot-shell: ## Open a shell in the Godot development container
	docker-compose --profile dev exec godot-dev bash

# Testing and utilities
test: ## Run tests (if available)
	@echo "Running tests..."
	docker-compose exec backend go test ./...

# Database commands
sqlc-generate: ## Generate Go code from SQL queries using sqlc
	@echo "Generating Go code from SQL..."
	cd backend && sqlc generate -f ./server/cmd/internal/server/db/config/sqlc.yml

clean: ## Remove all containers, networks, and volumes
	@echo "Cleaning up Docker resources..."
	docker-compose down -v --remove-orphans
	docker system prune -f

clean-all: ## Remove all containers, networks, volumes, and images
	@echo "Cleaning up all Docker resources..."
	docker-compose down -v --remove-orphans
	docker system prune -a -f

# Local development (without Docker)
run-local: ## Run the Go server locally (requires Go to be installed)
	@echo "Running server locally..."
	cd backend/server && go run cmd/main.go

install-deps: ## Install Go dependencies locally
	@echo "Installing Go dependencies..."
	cd backend/server && go mod download

# Health check
health: ## Check if the backend service is healthy
	@echo "Checking backend health..."
	curl -f http://localhost:8080/health || echo "Backend is not responding"

# Status
status: ## Show status of all services
	docker-compose ps
