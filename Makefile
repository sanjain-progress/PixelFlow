export PATH := $(shell go env GOPATH)/bin:$(PATH)

.PHONY: help up down restart logs ps test clean build rebuild health db-shell kafka-shell

# Default target
help:
	@echo "PixelFlow - Distributed Image Processing Application"
	@echo ""
	@echo "Available targets:"
	@echo "  make up          - Start all services in detached mode"
	@echo "  make down        - Stop all services"
	@echo "  make restart     - Restart all services"
	@echo "  make build       - Build all Docker images"
	@echo "  make rebuild     - Rebuild all images without cache"
	@echo "  make logs        - Follow logs from all services"
	@echo "  make ps          - Show running containers"
	@echo "  make test        - Run end-to-end tests"
	@echo "  make clean       - Stop services and remove volumes (CAUTION: deletes data)"
	@echo "  make health      - Check health of all services"
	@echo ""
	@echo "Service-specific logs:"
	@echo "  make logs-auth   - View Auth service logs"
	@echo "  make logs-api    - View API service logs"
	@echo "  make logs-worker - View Worker service logs"
	@echo ""
	@echo "Database access:"
	@echo "  make db-postgres - Access PostgreSQL shell"
	@echo "  make db-mongo    - Access MongoDB shell"
	@echo ""
	@echo "Kafka debugging:"
	@echo "  make kafka-topics   - List Kafka topics"
	@echo "  make kafka-consumer - Consume messages from image-tasks topic"
	@echo "  make kafka-groups   - List consumer groups"

# Start all services
up:
	@echo "🚀 Starting all PixelFlow services..."
	docker-compose up -d
	@echo "✅ Services started! Run 'make ps' to verify."
	@echo "💡 Run 'make test' to verify end-to-end workflow."

# Stop all services
down:
	@echo "🛑 Stopping all services..."
	docker-compose down
	@echo "✅ Services stopped."

# Restart all services
restart:
	@echo "🔄 Restarting all services..."
	docker-compose restart
	@echo "✅ Services restarted."

# Build all images
build:
	@echo "🔨 Building all Docker images..."
	docker-compose build
	@echo "✅ Build complete!"

# Rebuild without cache
rebuild:
	@echo "🔨 Rebuilding all images (no cache)..."
	docker-compose build --no-cache
	@echo "✅ Rebuild complete!"

# View logs from all services
logs:
	docker-compose logs -f

# View logs from specific services
logs-auth:
	docker-compose logs -f auth-service

logs-api:
	docker-compose logs -f api-service

logs-worker:
	docker-compose logs -f worker-service

logs-kafka:
	docker-compose logs -f kafka

logs-mongo:
	docker-compose logs -f mongo

logs-postgres:
	docker-compose logs -f postgres-auth

# Show container status
ps:
	docker-compose ps

# Run end-to-end tests
test:
	@echo "🧪 Running E2E tests..."
	@chmod +x ./test_e2e.sh
	@./test_e2e.sh

# Clean up (WARNING: Deletes volumes/data)
clean:
	@echo "⚠️  WARNING: This will delete all data (volumes). Press Ctrl+C to cancel..."
	@sleep 3
	docker-compose down -v
	@echo "✅ Cleanup complete. All data removed."

# Health checks
health:
	@echo "🏥 Checking service health..."
	@echo ""
	@echo "API Service:"
	@curl -s http://localhost:8080/health || echo "❌ API not responding"
	@echo ""
	@echo ""
	@echo "Container Status:"
	@docker-compose ps

# PostgreSQL shell
db-postgres:
	@echo "📊 Connecting to PostgreSQL..."
	docker exec -it pixelflow-postgres-auth psql -U postgres -d auth_db

# MongoDB shell
db-mongo:
	@echo "📊 Connecting to MongoDB..."
	docker exec -it pixelflow-mongo mongosh pixelflow

# Kafka topics
kafka-topics:
	@echo "📋 Listing Kafka topics..."
	docker exec pixelflow-kafka kafka-topics --list --bootstrap-server localhost:9092

# Kafka consumer (read messages)
kafka-consumer:
	@echo "📨 Consuming from image-tasks topic (Ctrl+C to stop)..."
	docker exec pixelflow-kafka kafka-console-consumer \
		--topic image-tasks \
		--from-beginning \
		--bootstrap-server localhost:9092

# Kafka consumer groups
kafka-groups:
	@echo "👥 Listing Kafka consumer groups..."
	docker exec pixelflow-kafka kafka-consumer-groups \
		--list --bootstrap-server localhost:9092
	@echo ""
	@echo "Worker group details:"
	@docker exec pixelflow-kafka kafka-consumer-groups \
		--describe --group worker-group-1 \
		--bootstrap-server localhost:9092

# Quick start (build + up + test)
quickstart: build up
	@echo "⏳ Waiting 10 seconds for services to initialize..."
	@sleep 10
	@make test

# Development helpers
dev-auth:
	cd apps/auth && go run cmd/main.go

dev-api:
	cd apps/api && go run cmd/main.go

dev-worker:
	cd apps/worker && go run cmd/main.go

# Frontend helpers
install-frontend:
	cd apps/frontend && npm install

dev-frontend:
	cd apps/frontend && npm start

build-frontend:
	cd apps/frontend && npm run build
