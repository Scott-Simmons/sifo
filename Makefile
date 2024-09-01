BINARY_NAME := sifo
GO_CMD := go
BUILD_DIR := build
DIST_DIR := dist
PACKAGE := ./cmd/cli

TEST_ARGS := -v -cover -coverprofile=coverage.out ./...
COVERAGE_DIR := coverage
COVERAGE_OUT := $(COVERAGE_DIR)/coverage.out
COVERAGE_HTML := $(COVERAGE_DIR)/coverage.html
COVERAGE_ARGS := -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
LINT_ARGS := -w

VERSION := $(shell git describe --tags --always --dirty)
COMMIT_HASH := $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
TARGETS := linux/amd64 darwin/amd64 windows/amd64
VERSION_FLAGS := -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.commitHash=$(COMMIT_HASH)

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

ifdef DEBUG
LD_FLAGS :=
else
LD_FLAGS := -s -w 
endif

GO_LD_FLAGS := -ldflags "$(LD_FLAGS) $(VERSION_FLAGS)"

.DEFAULT_GOAL := build

all: build

install: build
	@cp $(BUILD_DIR) /

test:
	go test $(TEST_ARGS)
	go tool cover $(COVERAGE_ARGS)
	go tool cover -func=$(COVERAGE_OUT)

SUBDIR_NAME ?= $(GOOS)-$(GOARCH)
build:
	echo $(BINARY_NAME)
	GOOS=$(GOOS) GOARCH=$(GOARCH) SUBDIR_NAME=$(SUBDIR_NAME) $(GO_CMD) build $(GO_LD_FLAGS) -o $(BUILD_DIR)/$(SUBDIR_NAME)/$(BINARY_NAME) $(PACKAGE)

build-all:
	@for target in $(TARGETS); do \
		GOOS=$${target%/*} GOARCH=$${target#*/} SUBDIR_NAME=$${target%/*}-$${target#*/} $(MAKE) build; \
	done
	

clean:
	rm -rf $(BUILD_DIR) $(DIST_DIR) $(COVERAGE_DIR)
	go mod tidy

lint: 
	gofmt $(LINT_ARGS) .

.PHONY: clean lint build test

