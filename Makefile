#! /usr/bin/make -f

PACKAGES=$(shell go list ./...)
VERSION := $(shell echo $(shell git describe --tags 2> /dev/null || echo "dev-$(shell git describe --always)") | sed 's/^v//')
SIMAPP = ./app
DOCKER := $(shell which docker)
COVER_FILE := coverage.txt
COVER_HTML_FILE := cover.html

## govet: Run go vet.
govet:
	@echo Running go vet...
	@go vet ./...

FIND_ARGS := -name '*.go' -type f -not -name '*.pb.go' -not -name '*.pb.gw.go'

## format: Run gofmt and goimports.
format:
	@echo Formatting...
	@find . $(FIND_ARGS) | xargs gofumpt -w .
	@find . $(FIND_ARGS) | xargs goimports -w -local github.com/ignite/cli-plugin-airdrop

## lint: Run Golang CI Lint.
lint:
	@echo Running gocilint...
	@golangci-lint run --out-format=tab --issues-exit-code=0

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)", or just run 'make' for install"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: lint format govet help

## test-unit: Run the unit tests.
test-unit:
	@echo Running unit tests...
	@VERSION=$(VERSION) go test -mod=readonly -v -timeout 30m $(PACKAGES)

## test-race: Run the unit tests checking for race conditions
test-race:
	@echo Running unit tests with race condition reporting...
	@VERSION=$(VERSION) go test -mod=readonly -v -race -timeout 30m  $(PACKAGES)

## test-cover: Run the unit tests and create a coverage html report
test-cover:
	@echo Running unit tests and creating coverage report...
	@VERSION=$(VERSION) go test -mod=readonly -v -timeout 30m -coverprofile=$(COVER_FILE) -coverpkg=./... -covermode=atomic $(PACKAGES)
	@go tool cover -html=$(COVER_FILE) -o $(COVER_HTML_FILE)
	@rm $(COVER_FILE)

## bench: Run the unit tests with benchmarking enabled
bench:
	@echo Running unit tests with benchmarking...
	@VERSION=$(VERSION) go test -mod=readonly -v -timeout 30m -bench=. $(PACKAGES)

## test: Run unit and integration tests.
test: govet test-unit

.PHONY: test test-unit test-race test-cover bench

.DEFAULT_GOAL := install
