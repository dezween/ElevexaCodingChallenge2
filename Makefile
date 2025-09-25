# Makefile for development tasks
# Usage: make <target>

.PHONY: all test testrace build run clean coverage fmt vet

all: test

# Run unit and integration tests (without race detector)
test:
	@echo "Running tests..."
	go test ./... -v

# Build the binary
build:
	@echo "Building binary..."
	go build -v -o kyber-server main.go

# Run the built binary (build first if needed)
run: build
	@echo "Running kyber-server..."
	./kyber-server

# Format code
fmt:
	gofmt -w .

# Vet code
vet:
	go vet ./...

