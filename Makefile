BINARY_NAME := sifo
GO_CMD := go
BUILD_DIR := build
DIST_DIR := dist
PACKAGE := ./cmd/cli
INSTALL_PREFIX ?= /usr/local
BIN_DIR := $(INSTALL_PREFIX)/bin
VERSION_FILE := VERSION
VERSION := $(shell cat $(VERSION_FILE))

COVERAGE_DIR := coverage
COVERAGE_OUT := $(COVERAGE_DIR)/coverage.out
COVERAGE_HTML := $(COVERAGE_DIR)/coverage.html
COVERAGE_ARGS := -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
TEST_ARGS := -v -cover -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
LINT_ARGS := -w

GIT_TAG := $(shell git describe --tags --always --dirty)
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

BUILD_BINARY_LOC ?= $(BUILD_DIR)/$(GOOS)-$(GOARCH)/$(BINARY_NAME)
install: build
	@mkdir -p $(BIN_DIR); \
	@if [ -f $(BUILD_BINARY_LOC) ]; then \
		@cp $(BUILD_BINARY_LOC) $(BIN_DIR)/$(BINARY_NAME); \
		@echo Installed $(BINARY_NAME) to $(BIN_DIR); \
	else \
		echo Error: $(BUILD_BINARY_LOC) does not exist.; \
		exit 1; \
	fi


uninstall:
	@if [ -f $(BIN_DIR)/$(BINARY_NAME) ]; then \
		@rm -rf $(BIN_DIR)/$(BINARY_NAME); \
		@echo Uninstalled $(BINARY_NAME) from $(BIN_DIR); \
	else \
		echo Error: $(BIN_DIR)/$(BINARY_NAME) does not exist.; \
		exit 1; \
	fi

test:
	@mkdir -p $(COVERAGE_DIR)
	go test $(TEST_ARGS)
	go tool cover $(COVERAGE_ARGS)
	go tool cover -func=$(COVERAGE_OUT)

SUBDIR_NAME ?= $(GOOS)-$(GOARCH)
build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) SUBDIR_NAME=$(SUBDIR_NAME) $(GO_CMD) build $(GO_LD_FLAGS) -o $(BUILD_DIR)/$(SUBDIR_NAME)/$(BINARY_NAME) $(PACKAGE)

build-all:
	@for target in $(TARGETS); do \
		GOOS=$${target%/*} GOARCH=$${target#*/} SUBDIR_NAME=$${target%/*}-$${target#*/} $(MAKE) build || exit 1; \
	done
	

clean:
	rm -rf $(BUILD_DIR) $(DIST_DIR) $(COVERAGE_DIR)
	go mod tidy

lint: 
	gofmt $(LINT_ARGS) .

RELEASE_NAME ?= Release $(VERSION)
new-tag:
	git tag -a $(VERSION) -m $(RELEASE_NAME)
	git push origin $(VERSION)

new-release: new-tag
	gh release create $(VERSION) $(BUILD_BINARY_LOC) --title $(RELEASE_NAME)

.PHONY: clean lint build test build-all uninstall install new-tag new-release

