CLI_BINARY_NAME=SecureSyncDrive
GO_CMD=go
GO_BUILD_FLAGS=-o $(CLI_BINARY_NAME) ./cmd/cli/

TEST_ARGS := -cover ./...

test:
	go test $(TEST_ARGS)

build:
	$(GO_CMD) build $(GO_BUILD_FLAGS)

clean:
	rm -f $(CLI_BINARY_NAME)

.PHONY: build clean

