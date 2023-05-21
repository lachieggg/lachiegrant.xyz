.PHONY: run clean docker test testgo htmlcoverage coverage help

build:
	go build -o bin/app src/*.go

run: clean build
	./bin/app

webpack:
	npm run build

.PHONY: docker
docker: webpack
	docker-compose up --build

.PHONY: clean
clean:
	rm -f bin/app

.PHONY: test
test:
	find . -type f -name "*.go" -exec dirname {} \; | sort -u | xargs -I {} go test -v {}

.PHONY: testgo
testgo:
	go test $(go list ./...) -coverageprofile=coverage.out

.PHONY: htmlcoverage
htmlcoverage:
	go tool cover --html=coverage.out

.PHONY: coverage
coverage:
	go test ./... -covermode=atomic -coverprofile=coverage.out -coverpkg=./...
	go tool cover -func=coverage.out

.PHONY: help
help:
	@echo "Available targets:"
	@grep '^[^#[:space:]].*:' Makefile | awk -F':' '{print $$1}' | grep -v ".PHONY"
