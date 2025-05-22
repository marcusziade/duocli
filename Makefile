# DuoCLI Makefile

.PHONY: build run clean test install dev

# Build the application
build:
	go build -o duocli .

# Run the application in interactive mode
run: build
	./duocli

# Clean build artifacts
clean:
	rm -f duocli duocli.db

# Run tests
test:
	go test ./...

# Install dependencies
deps:
	go mod tidy
	go mod download

# Development mode with hot reload (requires air)
dev:
	air

# Install the application globally
install: build
	sudo cp duocli /usr/local/bin/

# Create a release build
release:
	CGO_ENABLED=1 go build -ldflags="-w -s" -o duocli .

# Show help
help:
	@echo "Available commands:"
	@echo "  build    - Build the DuoCLI application"
	@echo "  run      - Build and run the application"
	@echo "  clean    - Remove build artifacts and database"
	@echo "  test     - Run tests"
	@echo "  deps     - Install/update dependencies"
	@echo "  dev      - Run in development mode (requires air)"
	@echo "  install  - Install globally to /usr/local/bin"
	@echo "  release  - Create optimized release build"
	@echo "  help     - Show this help message"

# Default target
all: build