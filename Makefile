.DEFAULT_GOAL := help
.PHONY: help build up down

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Install dependencies and build frontend
	npm install
	npm run build

up: build ## Build and start services in background
	docker-compose up --build -d

down: ## Stop all services
	docker-compose down
