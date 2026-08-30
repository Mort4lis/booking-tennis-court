BINARY_NAME := bot
CMD_PATH    := ./cmd/bot
BIN_DIR     := bin
GO          ?= go

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: help build run lint fmt test test-unit test-integration tidy generate \
	migrate migrate-down-step migrate-down \
	docker-build compose-up compose-down

## help: show this help
help:
	@grep -hE '^##' $(MAKEFILE_LIST) | sed -e 's/## //' | awk -F ': ' '{printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

## build: compile the binary into bin/
build:
	$(GO) build -trimpath -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_PATH)

## run: compile and run the bot
run: build
	./$(BIN_DIR)/$(BINARY_NAME) serve

## lint: run golangci-lint
lint:
	golangci-lint -v run ./...

## fmt: format the code and run go vet
fmt:
	$(GO) fmt ./...
	$(GO) vet ./...

## test: run unit tests, then integration tests (requires Docker)
test: test-unit test-integration

## test-unit: run unit tests with the race detector
test-unit:
	$(GO) test -v -race -count=1 ./...

## test-integration: run integration tests (requires Docker)
test-integration:
	$(GO) test -v -race -count=1 -tags=integration ./test/...

## tidy: tidy go.mod / go.sum
tidy:
	$(GO) mod tidy

## generate: regenerate mocks with mockery v3 (config in .mockery.yml)
generate:
	$(GO) tool mockery

## migrate: apply all pending migrations
migrate:
	$(GO) run $(CMD_PATH) migrate up

## migrate-down-step: roll back one migration
migrate-down-step:
	$(GO) run $(CMD_PATH) migrate down --step 1

## migrate-down: roll back all migrations
migrate-down:
	$(GO) run $(CMD_PATH) migrate down --all

## docker-build: build the Docker image
docker-build:
	docker build --build-arg VERSION=$(VERSION) -t $(BINARY_NAME):$(VERSION) .

## compose-up: start local dependencies (PostgreSQL, Jaeger)
compose-up:
	docker compose up -d

## compose-down: stop local dependencies
compose-down:
	docker compose down
