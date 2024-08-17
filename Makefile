CLI_BINARY_NAME=SecureSyncDrive
GO_CMD=go
GO_BUILD_FLAGS=-o $(CLI_BINARY_NAME) ./cmd/cli/

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

