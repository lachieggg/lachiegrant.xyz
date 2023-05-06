.PHONY: run clean docker test testgo htmlcoverage coverage help

build:
	go build -o bin/app src/*.go

run:
	go run src/main.go

clean:
	rm -f bin/app

docker:
	docker-compose up --build -d

test:
	find . -type f -name "*.go" -exec dirname {} \; | sort -u | xargs -I {} go test -v {}

testgo:
	go test $(go list ./...) -coverageprofile=coverage.out

htmlcoverage:
	go tool cover --html=coverage.out

coverage:
	go test ./... -covermode=atomic -coverprofile=coverage.out -coverpkg=./...
	go tool cover -func=coverage.out

help:
	@echo "Available targets:"
	@grep '^[^#[:space:]].*:' Makefile | awk -F':' '{print $$1}'
