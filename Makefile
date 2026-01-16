.PHONY: run build clean test

# Variables
BINARY_NAME=server
MAIN_PATH=./cmd/server

# Run the application
run:
	go run $(MAIN_PATH)/main.go

# Build the application
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)/main.go

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

# Run tests
test:
	go test -v ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Development with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air
