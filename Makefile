.PHONY: help up down local-certs certbot test test-go test-cov test-js test-nginx reload logs bans update-bots
.DEFAULT_GOAL := help

up: build ## Build and start services in background
	ALPINE_REPO=http://dl-cdn.alpinelinux.org/alpine/edge/community docker-compose up --build -d

down: ## Stop all services
	docker-compose down

build: ## Install dependencies and build frontend
	npm install
	npm run build

reset: build down up ## Stop, rebuild, and restart services

reload: ## Rebuild frontend and backend, then restart app
	@out=$$(npm run build 2>&1) || { echo "$$out"; exit 1; }
	@docker exec app go build -o /server_bin/app ./src
	@docker-compose restart app > /dev/null 2>&1
	@echo "Reloaded successfully."

test: ## Run all tests
	@echo "--- 🔧 Running Go Tests 🔧 ---"
	@go test -v ./src/... && GO_OK=1 || GO_OK=0; \
	echo "\n--- 🔧 Running JS Tests 🔧 ---"; \
	npm test && JS_OK=1 || JS_OK=0; \
	echo "\n--- 🔧 Linting Nginx Config 🔧 ---"; \
	$(MAKE) test-nginx && NGINX_OK=1 || NGINX_OK=0; \
	echo "-------------------------------"; \
	echo "         TEST SUMMARY          "; \
	echo "-------------------------------"; \
	if [ $$GO_OK -eq 1 ]; then echo "✅ Go tests OK"; else echo "❌ Go tests failed 😵"; fi; \
	if [ $$JS_OK -eq 1 ]; then echo "✅ JS tests OK"; else echo "❌ JS tests failed 😵"; fi; \
	if [ $$NGINX_OK -eq 1 ]; then echo "✅ Nginx tests OK"; else echo "❌ Nginx tests failed 😵"; fi; \
	echo "-------------------------------"; \
	if [ $$GO_OK -eq 0 ] || [ $$JS_OK -eq 0 ] || [ $$NGINX_OK -eq 0 ]; then exit 1; fi

test-go: ## Run backend unit tests
	go test -v ./src/...

test-cov: ## Show backend coverage in browser
	go test -coverprofile=coverage.out ./src/...
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html || xdg-open coverage.html || start coverage.html

test-js: ## Run frontend javascript tests
	npm test

test-nginx: ## Validate Nginx configuration syntax natively in the running webserver container
	@docker exec webserver /bin/sh -c "envsubst '\$$ALLOWED_HOSTS' < /etc/nginx/nginx.conf.template > /tmp/nginx-test.conf && nginx -t -c /tmp/nginx-test.conf"

certbot: ## Run certbot script for production SSL certificates
	./scripts/certbot.sh

local-certs: ## Generate self-signed SSL certificates for local development
	mkdir -p ./tls/nginx/keys ./tls/nginx/certs
	openssl req -x509 -nodes -days 365 -newkey rsa:4096 \
		-keyout ./tls/nginx/keys/privkey.pem \
		-out ./tls/nginx/certs/fullchain.pem

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

logs: ## View live tail of all Docker container logs
	docker-compose logs -f

bans: ## View the currently active Fail2Ban IP hitlist
	docker exec -it fail2ban fail2ban-client status nginx-bot-ban

kill: ## Kill all running containers
	docker kill $(docker ps -q)

prune: ## Remove all unused Docker resources (hard reset)
	docker system prune -af