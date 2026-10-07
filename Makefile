BINARY_NAME=hann
BIN_DIR=bin
VERSION?=0.1.0
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
BUILD_DATE?=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

MODULE=github.com/frhndevweb/hann-cli
LDFLAGS=-ldflags "-s -w \
	-X $(MODULE)/internal/version.Version=$(VERSION) \
	-X $(MODULE)/internal/version.Commit=$(COMMIT) \
	-X $(MODULE)/internal/version.BuildDate=$(BUILD_DATE)"

.PHONY: all build test install clean lint fmt check

all: build

build:
	@mkdir -p $(BIN_DIR)
	go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) .

test:
	go test -v ./...

install: build
	@mkdir -p $(HOME)/.local/bin
	@cp $(BIN_DIR)/$(BINARY_NAME) $(HOME)/.local/bin/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(HOME)/.local/bin/$(BINARY_NAME)"

clean:
	rm -rf $(BIN_DIR)

lint:
	go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	fi

fmt:
	go fmt ./...

check: fmt lint test
