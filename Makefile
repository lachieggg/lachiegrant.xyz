.PHONY: help up down local-certs certbot test reload
.DEFAULT_GOAL := help

up: build ## Build and start services in background
	ALPINE_REPO=http://dl-cdn.alpinelinux.org/alpine/edge/community docker-compose up --build -d

down: ## Stop all services
	docker-compose down

build: ## Install dependencies and build frontend
	npm install
	npm run build

reset: build down up ## Stop, rebuild, and restart services

reload: ## Restart only the app container to reload environment variables
	docker-compose restart app

test: ## Run backend unit tests
	go test -v ./src/...

certbot: ## Run certbot script for production SSL certificates
	./scripts/certbot.sh

local-certs: ## Generate self-signed SSL certificates for local development
	mkdir -p ./tls/nginx/keys ./tls/nginx/certs
	openssl req -x509 -nodes -days 365 -newkey rsa:4096 \
		-keyout ./tls/nginx/keys/privkey.pem \
		-out ./tls/nginx/certs/fullchain.pem

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)