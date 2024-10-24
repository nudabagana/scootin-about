include .env
export

BINARY_NAME=scootin-service
.DEFAULT_GOAL := build

.PHONY: build
build:
	go mod tidy
	go build -o bin/${BINARY_NAME}

.PHONY: start
start:
	docker-compose up -d
	./bin/${BINARY_NAME}

.PHONY: build_start
build_start: build start

.PHONY: start_watch
start_watch: 
	docker-compose up -d
	$$(go env GOPATH)/bin/air -c .air.toml
	
.PHONY: test_e2e
test_e2e:
	docker-compose up -d
	go test ./... -run TestE2E
	
.PHONY: stop
stop:
	docker-compose down
	if [ -d bin ] ; then rm -rf bin ; fi
