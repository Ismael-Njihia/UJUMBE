.PHONY: help build run test clean docker-up docker-down migrate

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Go application
	@echo "Building application..."
	cd backend && go build -o ../bin/ujumbe ./cmd/api

run: ## Run the application locally
	@echo "Running application..."
	cd backend && go run ./cmd/api/main.go ./cmd/api/handlers.go

test: ## Run tests
	@echo "Running tests..."
	cd backend && go test ./... -v

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf frontend/dist/
	rm -rf frontend/node_modules/

docker-up: ## Start all services with Docker Compose
	@echo "Starting services..."
	docker-compose up -d

docker-down: ## Stop all services
	@echo "Stopping services..."
	docker-compose down

docker-logs: ## View logs from all services
	docker-compose logs -f

docker-build: ## Build Docker images
	docker-compose build

frontend-install: ## Install frontend dependencies
	cd frontend && npm install

frontend-dev: ## Run frontend development server
	cd frontend && npm run dev

frontend-build: ## Build frontend for production
	cd frontend && npm run build

deps: ## Install Go dependencies
	cd backend && go mod download
	cd backend && go mod tidy

migrate: ## Run database migrations (placeholder)
	@echo "Running migrations..."
	# Add migration commands here

lint: ## Run linters
	@echo "Running linters..."
	cd backend && go fmt ./...
	cd backend && go vet ./...

dev: ## Run both backend and frontend in development mode
	@echo "Starting development environment..."
	@make -j2 run frontend-dev

all: deps build ## Install dependencies and build

.DEFAULT_GOAL := help
