# Makefile for Prometheus CLI

# Variables
BINARY_NAME=prom-cli
GO=go
GOFMT=gofmt
GOTEST=$(GO) test
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOMOD=$(GO) mod
GOGET=$(GO) get
BIN_DIR=bin

# Version information
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
REVISION ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
BUILD_USER ?= $(shell whoami)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# LDFLAGS for injecting version information
LDFLAGS = \
	-X 'github.com/prometheus/common/version.Version=$(VERSION)' \
	-X 'github.com/prometheus/common/version.Revision=$(REVISION)' \
	-X 'github.com/prometheus/common/version.Branch=$(BRANCH)' \
	-X 'github.com/prometheus/common/version.BuildUser=$(BUILD_USER)' \
	-X 'github.com/prometheus/common/version.BuildDate=$(BUILD_DATE)'

# Build targets for different platforms
BINARY_UNIX=$(BIN_DIR)/$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BIN_DIR)/$(BINARY_NAME).exe
BINARY_MACOS=$(BIN_DIR)/$(BINARY_NAME)_macos

# Main targets
.PHONY: all build clean test fmt vet run deps lint help

all: test build

build:
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME) -v ./cmd/prom-cli

clean:
	$(GOCLEAN)
	rm -rf $(BIN_DIR)

test:
	$(GOTEST) -v ./...

fmt:
	$(GOFMT) -w .

vet:
	$(GO) vet ./...

run: build
	$(BIN_DIR)/$(BINARY_NAME)

deps:
	$(GOMOD) download

# Linting target
lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint is not installed. Please install it from https://golangci-lint.run/usage/install/"; exit 1; }
	@echo "Running golangci-lint..."
	@golangci-lint run --verbose

# Cross-compilation targets
.PHONY: build-all build-linux build-windows build-macos

build-all: build-linux build-windows build-macos

build-linux:
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_UNIX) -v ./cmd/prom-cli

build-windows:
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_WINDOWS) -v ./cmd/prom-cli

build-macos:
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BINARY_MACOS) -v ./cmd/prom-cli

# Help target
help:
	@echo "Available targets:"
	@echo "  all        - Run tests and build the binary"
	@echo "  build      - Build the binary for the current platform"
	@echo "  clean      - Remove binaries and temporary files"
	@echo "  test       - Run all tests"
	@echo "  fmt        - Format the code"
	@echo "  vet        - Run go vet"
	@echo "  run        - Build and run the binary"
	@echo "  deps       - Download dependencies"
	@echo "  lint       - Run golangci-lint to check code quality"
	@echo "  build-all  - Build binaries for Linux, Windows, and macOS"
	@echo "  build-linux   - Build binary for Linux"
	@echo "  build-windows - Build binary for Windows"
	@echo "  build-macos   - Build binary for macOS"