# Makefile for development tasks
# Usage: make <target>

.PHONY: all test testrace build run clean coverage fmt vet

all: test

# Run unit and integration tests (without race detector)
test:
	@echo "Running tests..."
	go test ./... -v

# Build the binary (removes old binary first)
build:
ifeq ($(OS),Windows_NT)
	@echo "Building binary for Windows..."
	@if exist kyber-server.exe del kyber-server.exe
	go build -v -o kyber-server.exe main.go
else
	@echo "Building binary for Unix..."
	@if [ -f kyber-server ]; then rm kyber-server; fi
	go build -v -o kyber-server main.go
endif

# Run the built binary (build first if needed)
run:
ifeq ($(OS),Windows_NT)
	@echo "Running kyber-server.exe..."
	kyber-server.exe
else
	@echo "Running kyber-server..."
	./kyber-server
endif

# Format code
fmt:
	gofmt -w .

# Vet code
vet:
	go vet ./...

# Clean up binaries
clean:
ifeq ($(OS),Windows_NT)
	@if exist kyber-server.exe del kyber-server.exe
else
	@if [ -f kyber-server ]; then rm kyber-server; fi
endif
