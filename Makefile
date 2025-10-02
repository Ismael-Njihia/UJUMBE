.PHONY: help build run stop clean test migrate-up migrate-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build Docker images
	docker-compose build

run: ## Start all services
	docker-compose up -d
	@echo "Services started!"
	@echo "API: http://localhost:8080"
	@echo "Health check: http://localhost:8080/health"

stop: ## Stop all services
	docker-compose down

clean: ## Remove all containers, volumes, and images
	docker-compose down -v --rmi all

logs: ## View logs from all services
	docker-compose logs -f

logs-backend: ## View backend logs
	docker-compose logs -f backend

logs-worker: ## View worker logs
	docker-compose logs -f worker

logs-kafka: ## View Kafka logs
	docker-compose logs -f kafka

test: ## Run tests
	cd backend && go test ./...

migrate-up: ## Run database migrations
	docker exec -it ujumbe-postgres psql -U postgres -d ujumbe -f /docker-entrypoint-initdb.d/001_init_schema.sql

dev: ## Run backend in development mode
	cd backend && go run cmd/server/main.go

worker-dev: ## Run worker in development mode
	cd backend && go run cmd/worker/main.go

install-deps: ## Install Go dependencies
	cd backend && go mod download

tidy: ## Tidy Go modules
	cd backend && go mod tidy
