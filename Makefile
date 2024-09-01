SHELL := /bin/bash
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

all: clean lint build test

test:
	@mkdir -p $(COVERAGE_DIR)
	go test $(TEST_ARGS)
	go tool cover $(COVERAGE_ARGS)
	go tool cover -func=$(COVERAGE_OUT)

build:
	$(GO_CMD) build $(GO_BUILD_FLAGS)

clean:
	rm -f $(CLI_BINARY_NAME)
	go mod tidy

lint: 
	gofmt $(LINT_ARGS) .

.PHONY: clean lint build test

