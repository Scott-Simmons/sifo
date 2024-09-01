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

TEST_ARGS := -v -cover -coverprofile=coverage.out ./...
COVERAGE_ARGS := -html=coverage.out -o coverage.html
LINT_ARGS := -w

all: clean lint build test

test:
	go test $(TEST_ARGS)
	go tool cover $(COVERAGE_ARGS)
	go tool cover -func=coverage.out

build:
	$(GO_CMD) build $(GO_BUILD_FLAGS)

clean:
	rm -f $(CLI_BINARY_NAME)
	go mod tidy

lint: 
	gofmt $(LINT_ARGS) .

.PHONY: clean lint build test

