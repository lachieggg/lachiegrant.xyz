.DEFAULT_GOAL := help
.PHONY: help build run clean docker daemon kill remove test coverage

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install dependencies (npm + Go)
	npm install
	go mod download

build: ## Build frontend and Go application
	npm run build
	go build -o bin/app src/*.go

run: clean build ## Clean, build, and run locally
	./bin/app

clean: ## Remove compiled binaries and frontend build
	rm -f bin/app
	rm -rf dist node_modules

docker: ## Build and run with docker-compose (foreground)
	docker-compose up --build

daemon: install build ## Install deps, build frontend, run docker-compose in background
	docker-compose up --build -d

daemon-stop: ## Stop background docker-compose
	docker-compose down

logs: ## View docker-compose logs
	docker-compose logs -f

kill: ## Kill all running Docker containers
	docker kill $$(docker ps -q) || true

remove: kill ## Remove Docker volumes and containers
	docker volume prune -f && docker rm -f webserver app || true

test: ## Run Go tests with verbose output
	go test -v ./...

coverage: ## Generate test coverage report
	go test ./... -covermode=atomic -coverprofile=coverage.out -coverpkg=./...
	go tool cover -func=coverage.out

coverage-html: coverage ## Generate HTML coverage report
	go tool cover -html=coverage.out

certbot: ## Run certbot script for TLS certificates
	./scripts/certbot.sh
