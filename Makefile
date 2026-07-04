APP_NAME := distributed-rate-limiter
BUILD_DIR := bin

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  run      - Run the server"
	@echo "  build    - Build the binary"
	@echo "  test     - Run all tests"
	@echo "  fmt      - Format Go code"
	@echo "  vet      - Run go vet"
	@echo "  lint     - Run golangci-lint"
	@echo "  tidy     - Tidy Go modules"
	@echo "  clean    - Remove build artifacts"

.PHONY: run
run:
	go run ./cmd/server --addr=:8080

.PHONY: node1
node1:
	go run ./cmd/server --addr=:8081

.PHONY: node2
node2:
	go run ./cmd/server --addr=:8082

.PHONY: node3
node3:
	go run ./cmd/server --addr=:8083

.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

.PHONY: test
test:
	go test -v ./...

.PHONY: test-race
test-race:
	go test -race ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)