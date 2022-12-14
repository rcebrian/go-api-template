SHELL = /bin/bash
MAKEFLAGS += --silent

.PHONY: build
build: api-clean api
	go build -o ./build/users-api-server ./cmd/users-api-server/main.go

build-clean:
	rm -rf ./build/*

.PHONY: api
api:
	scripts/api/codegen.sh

api-clean:
	find internal/platform/server/openapi -type f -not -name '*_service.go'  -and ! -name '*_service_test.go' \
		-and ! -name '*mapper*.go' -delete

lint:
	goimports -w .
	golangci-lint run

clean: build-clean api-clean

all: clean api build

test:
	go test ./...

test-coverage:
	go test ./... -cover

test-coverage-reporter:
	scripts/test/reporter.sh
