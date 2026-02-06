.PHONY: build install clean test lint fmt run help

# Build variables
BINARY_NAME=fhir-cli
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/jbogarin/fhir-cli/cmd.Version=$(VERSION) -X github.com/jbogarin/fhir-cli/cmd.BuildDate=$(BUILD_DATE)"

# Go commands
GO=go
GOTEST=$(GO) test
GOBUILD=$(GO) build
GOINSTALL=$(GO) install
GOFMT=gofmt
GOLINT=golangci-lint

# Default target
all: build

## build: Build the binary
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

## install: Install the binary to GOPATH/bin
install:
	$(GOINSTALL) $(LDFLAGS) .

## clean: Remove build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

## test: Run tests
test:
	$(GOTEST) -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

## lint: Run linters
lint:
	$(GOLINT) run

## fmt: Format code
fmt:
	$(GOFMT) -s -w .

## tidy: Tidy go modules
tidy:
	$(GO) mod tidy

## deps: Download dependencies
deps:
	$(GO) mod download

## run: Run the CLI
run:
	$(GO) run . $(ARGS)

## build-all: Build for all platforms
build-all: clean
	mkdir -p dist
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 .
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe .

## help: Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
