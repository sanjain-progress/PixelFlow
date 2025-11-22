export PATH := $(shell go env GOPATH)/bin:$(PATH)

.PHONY: all build test clean proto

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

# Generate Protobuf code
proto:
	@echo "Generating Protobuf code..."
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/auth.proto
	@echo "Moving generated files to pkg/pb..."
	mkdir -p pkg/pb
	mv proto/*.go pkg/pb/ 2>/dev/null || true
	@echo "Done!"
