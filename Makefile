.DEFAULT_GOAL := help
.PHONY: help build up down certs certbot

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Install dependencies and build frontend/backend
	npm install
	npm run build
	mkdir -p bin
	go build -o bin/app ./src

up: build ## Build and start services in background
	BUILD_DATE=$$(date -u +'%Y-%m-%dT%H:%M:%SZ') docker-compose up --build -d

down: ## Stop all services
	docker-compose down

reset: build down up

certs: ## Generate self-signed SSL certificates for local development
	mkdir -p ./tls/nginx/keys ./tls/nginx/certs
	openssl req -x509 -nodes -days 365 -newkey rsa:4096 \
		-keyout ./tls/nginx/keys/privkey.pem \
		-out ./tls/nginx/certs/fullchain.pem

certbot: ## Run certbot script for production SSL certificates
	./scripts/certbot.sh
