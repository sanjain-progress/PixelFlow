.PHONY: all build test clean

# Default target
all: build

# Build all services
build:
	@echo "Building Auth Service..."
	cd apps/auth && go build -o ../../bin/auth ./cmd
	@echo "Building API Service..."
	cd apps/api && go build -o ../../bin/api ./cmd
	@echo "Building Worker Service..."
	cd apps/worker && go build -o ../../bin/worker ./cmd
	@echo "Build complete!"

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
