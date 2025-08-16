# Makefile for pomodoro - a minimalist macOS CLI Pomodoro timer

# Variables
BINARY_NAME=pomodoro
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE=$(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildDate=$(BUILD_DATE) -s -w"

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)
GOPATH=$(shell go env GOPATH)

# Build targets
.PHONY: all build clean test coverage fmt lint vet install uninstall help

all: test build

build: 
	@echo "Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o $(GOBIN)/$(BINARY_NAME) .

# Build for all supported platforms
build-all: clean
	@echo "Building for all platforms..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(GOBIN)/$(BINARY_NAME)_darwin_amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(GOBIN)/$(BINARY_NAME)_darwin_arm64 .
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(GOBIN)/$(BINARY_NAME)_linux_amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(GOBIN)/$(BINARY_NAME)_linux_arm64 .

# Install the binary in $GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME) to $(GOPATH)/bin..."
	@cp $(GOBIN)/$(BINARY_NAME) $(GOPATH)/bin/

# Uninstall the binary from $GOPATH/bin
uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(GOPATH)/bin..."
	@rm -f $(GOPATH)/bin/$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run golint
lint:
	@echo "Running linter..."
	@if command -v golint >/dev/null 2>&1; then \
		golint ./...; \
	else \
		echo "golint not installed. Run: go install golang.org/x/lint/golint@latest"; \
	fi

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(GOBIN)
	@rm -f coverage.out coverage.html
	@rm -f $(BINARY_NAME)

# Help command
help:
	@echo "pomodoro Makefile help:"
	@echo "make              - Run tests and build binary"
	@echo "make build        - Build for current platform"
	@echo "make build-all    - Build for all supported platforms"
	@echo "make install      - Install binary to GOPATH/bin"
	@echo "make uninstall    - Remove binary from GOPATH/bin" 
	@echo "make test         - Run tests"
	@echo "make coverage     - Generate test coverage report"
	@echo "make fmt          - Format code"
	@echo "make lint         - Run linter"
	@echo "make vet          - Run go vet"
	@echo "make clean        - Remove build artifacts"
