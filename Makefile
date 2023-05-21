.PHONY: run clean docker test testgo htmlcoverage coverage help

build:
	go build -o bin/app src/*.go

run: clean build
	./bin/app

.PHONY: kill
kill:
	docker kill $(docker ps -q)	

.PHONY: compile
compile:
	docker exec -it app go build -o /app/bin/app /app/src/

.PHONY: remove
remove:
	docker volume prune -f && docker rm -f webserver && docker rm -f app && docker ps -a

.PHONY: docker
export BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
docker:
	npm run build
	docker-compose up --build

.PHONY: daemon
export BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
daemon:
	npm run build
	docker-compose up --build -d

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

.PHONY: tidyhtml
tidyhtml:
	find . -name "*.html" -exec tidy -m {} \;