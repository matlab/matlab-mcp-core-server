# Copyright 2025 The MathWorks, Inc.

# Set shell based on OS
ifeq ($(OS),Windows_NT)
	SHELL = powershell.exe
else
	SHELL = sh
endif

# Race detector flag
# Note: Disabled on Windows because CI agents don't have gcc available (required for -race)
ifeq ($(OS),Windows_NT)
    RACE_FLAG =
else
    RACE_FLAG = -race
endif

SEMANTIC_VERSION=v0.2.0
COMMIT_HASH := $(shell git rev-parse HEAD)

# Append Git commit hash to version unless building a release
ifeq ($(RELEASE),true)
	VERSION := $(SEMANTIC_VERSION)
else
	VERSION := $(SEMANTIC_VERSION).$(COMMIT_HASH)
endif

ifeq ($(OS),Windows_NT)
    RM_DIR = if (Test-Path "$(1)") { Remove-Item -Recurse -Force "$(1)" }
	PATHSEP = ;
	BIN_PATH = $(CURDIR)/.bin/win64
else
    RM_DIR = rm -rf $(1)
	PATHSEP = :
	BIN_PATH = $(CURDIR)/.bin/glnxa64
endif

# Capture CLI Environment variables
CLI_MATLAB_MCP_CORE_SERVER_BUILD_DIR := $(MATLAB_MCP_CORE_SERVER_BUILD_DIR)
CLI_MCP_MATLAB_PATH := $(MCP_MATLAB_PATH)

# Include .env file if it exists
ifneq (,$(wildcard .env))
    include .env
endif

# Set MATLAB_MCP_CORE_SERVER_BUILD_DIR with precendence CLI > .env > default
ifdef CLI_MATLAB_MCP_CORE_SERVER_BUILD_DIR
	MATLAB_MCP_CORE_SERVER_BUILD_DIR = $(CLI_MATLAB_MCP_CORE_SERVER_BUILD_DIR)
endif
ifndef MATLAB_MCP_CORE_SERVER_BUILD_DIR
	MATLAB_MCP_CORE_SERVER_BUILD_DIR = $(CURDIR)/.bin
endif
export MATLAB_MCP_CORE_SERVER_BUILD_DIR

# Set MCP_MATLAB_PATH with precendence CLI > .env > default (empty)
ifdef CLI_MCP_MATLAB_PATH
	MCP_MATLAB_PATH = $(CLI_MCP_MATLAB_PATH)
endif
export MCP_MATLAB_PATH

# Variables for MCP Inspector
export HOST = localhost
export PATH := $(BIN_PATH)$(PATHSEP)$(PATH)

# Go build flags
LDFLAGS := -ldflags "-X 'github.com/matlab/matlab-mcp-core-server/internal/adaptors/application/config.version=$(VERSION)'"

all: install wire mockery lint unit-tests build

version:
	@echo $(VERSION)

mcp-inspector: build
	npx @modelcontextprotocol/inspector matlab-mcp-core-server

# File checks

install:
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/vektra/mockery/v3@latest
	go install gotest.tools/gotestsum@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2

wire:
	wire github.com/matlab/matlab-mcp-core-server/internal/wire

mockery:
	@$(call RM_DIR,./mocks)
	@$(call RM_DIR,./tests/mocks)
	mockery

lint:
	golangci-lint run ./...

fix-lint:
	golangci-lint run ./... --fix

# Building

build: build-for-windows build-for-glnxa64 build-for-maci64 build-for-maca64

build-for-windows:
ifeq ($(OS),Windows_NT)
	$$env:GOOS='windows'; $$env:GOARCH='amd64'; go build $(LDFLAGS) -o ./.bin/win64/matlab-mcp-core-server.exe ./cmd/matlab-mcp-core-server
else
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o ./.bin/win64/matlab-mcp-core-server.exe ./cmd/matlab-mcp-core-server
endif

build-for-glnxa64:
ifeq ($(OS),Windows_NT)
	$$env:GOOS='linux'; $$env:GOARCH='amd64'; go build $(LDFLAGS) -o ./.bin/glnxa64/matlab-mcp-core-server ./cmd/matlab-mcp-core-server
else
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ./.bin/glnxa64/matlab-mcp-core-server ./cmd/matlab-mcp-core-server
endif

build-for-maci64:
ifeq ($(OS),Windows_NT)
	$$env:GOOS='darwin'; $$env:GOARCH='amd64'; go build $(LDFLAGS) -o ./.bin/maci64/matlab-mcp-core-server ./cmd/matlab-mcp-core-server
else
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o ./.bin/maci64/matlab-mcp-core-server ./cmd/matlab-mcp-core-server
endif

build-for-maca64:
ifeq ($(OS),Windows_NT)
	$$env:GOOS='darwin'; $$env:GOARCH='arm64'; go build $(LDFLAGS) -o ./.bin/maca64/matlab-mcp-core-server ./cmd/matlab-mcp-core-server
else
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o ./.bin/maca64/matlab-mcp-core-server ./cmd/matlab-mcp-core-server
endif

# Testing

unit-tests:
	gotestsum --packages="./internal/... ./tests/testutils/..." -- -race -coverprofile cover.out

system-tests:
	gotestsum --packages="./tests/system/..." -- -race -count=1 -timeout 30m
	@$(MAKE) --no-print-directory check-matlab-leaks

ci-unit-tests:
	go test $(RACE_FLAG) -json -count=1 -coverprofile cover.out ./internal/... ./tests/testutils/...

ci-system-tests:
	go test $(RACE_FLAG) -timeout 120m -json -count=1 ./tests/system/...
	@$(MAKE) --no-print-directory check-matlab-leaks

# Check for leaked MATLAB processes after system tests
# Tests should clean up all MATLAB sessions they create
check-matlab-leaks:
	@echo "Waiting for processes to settle..."
ifeq ($(OS),Windows_NT)
	@powershell -Command "Start-Sleep -Seconds 5"
	@echo "Checking for leaked MATLAB processes..."
	@powershell -Command "$$procs = Get-Process -Name MATLAB -ErrorAction SilentlyContinue | Where-Object { $$_.CommandLine -like '*matlab-mcp-core-server*' }; if ($$procs) { Write-Host 'WARNING: Found leaked MATLAB processes:'; $$procs | Format-Table Id,ProcessName,StartTime; exit 1 } else { Write-Host 'No leaked MATLAB processes found.' }"
else
	@sleep 5
	@echo "Checking for leaked MATLAB processes..."
	@leaked=$$(pgrep -a -f 'matlab.*matlab-mcp-core-server' | grep -v 'make\|grep' || true); \
	if [ -n "$$leaked" ]; then \
		echo "WARNING: Found leaked MATLAB processes:"; \
		echo "$$leaked"; \
		exit 1; \
	else \
		echo "No leaked MATLAB processes found."; \
	fi
endif
