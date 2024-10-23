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
	docker-compose up
	./bin/${BINARY_NAME}

.PHONY: build_start
build_start: build start

.PHONY: start_watch
start_watch: 
	# docker-compose up
	$$(go env GOPATH)/bin/air -c .air.toml
	
.PHONY: test
test:
	go test -tags=unit ./... -count=1

.PHONY: test_integration
test_integration:
	docker-compose up --no-recreate -d pybe-db configs-api transaction-kit
	go test -tags=integration ./... -count=1
	
.PHONY: stop
stop:
	docker-compose down
	if [ -d bin ] ; then rm -rf bin ; fi
