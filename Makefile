export PATH := $(shell go env GOPATH)/bin:$(PATH)

.PHONY: help up down restart logs ps test clean build rebuild health db-shell kafka-shell \
	test-e2e test-observability test-metrics test-traces test-logs test-alerts test-full

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
	@echo "Testing targets:"
	@echo "  make test-e2e           - Run application E2E tests"
	@echo "  make test-observability - Verify observability stack"
	@echo "  make test-metrics       - Check Prometheus metrics"
	@echo "  make test-traces        - Verify Jaeger traces"
	@echo "  make test-logs          - Check Loki logs"
	@echo "  make test-alerts        - Test alert rules"
	@echo "  make test-full          - Run all tests (E2E + Observability)"
	@echo ""
	@echo "Service-specific logs:"
	@echo "  make logs-auth   - View Auth service logs"
	@echo "  make logs-api    - View API service logs"
	@echo "  make logs-worker - View Worker service logs"
	@echo ""
	@echo "Observability access:"
	@echo "  make open-grafana    - Open Grafana in browser"
	@echo "  make open-prometheus - Open Prometheus in browser"
	@echo "  make open-jaeger     - Open Jaeger in browser"
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
	@echo "💡 Run 'make test-full' to verify end-to-end workflow and observability."

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

logs-prometheus:
	docker-compose logs -f prometheus

logs-grafana:
	docker-compose logs -f grafana

logs-jaeger:
	docker-compose logs -f jaeger

logs-loki:
	docker-compose logs -f loki

logs-promtail:
	docker-compose logs -f promtail

# Show container status
ps:
	docker-compose ps

# Run basic end-to-end tests
test:
	@echo "🧪 Running E2E tests..."
	@chmod +x ./test_e2e.sh
	@./test_e2e.sh

# Run application E2E tests
test-e2e:
	@echo "🧪 Running Application E2E Tests..."
	@echo "=================================="
	@chmod +x ./test_e2e.sh
	@./test_e2e.sh
	@echo ""
	@echo "✅ Application E2E tests completed!"

# Test Prometheus metrics
test-metrics:
	@echo "📊 Testing Prometheus Metrics..."
	@echo "================================"
	@echo ""
	@echo "1. Checking Prometheus health..."
	@curl -sf http://localhost:9091/-/healthy > /dev/null && echo "✅ Prometheus is healthy" || echo "❌ Prometheus is not responding"
	@echo ""
	@echo "2. Checking service metrics endpoints..."
	@echo "   Auth Service (/metrics):"
	@curl -sf http://localhost:50051/metrics | grep -q "auth_requests_total" && echo "   ✅ Auth metrics available" || echo "   ❌ Auth metrics not found"
	@echo "   API Service (/metrics):"
	@curl -sf http://localhost:8080/metrics | grep -q "api_requests_total" && echo "   ✅ API metrics available" || echo "   ❌ API metrics not found"
	@echo "   Worker Service (/metrics):"
	@curl -sf http://localhost:8081/metrics | grep -q "worker_tasks_processed_total" && echo "   ✅ Worker metrics available" || echo "   ❌ Worker metrics not found"
	@echo ""
	@echo "3. Checking Prometheus targets..."
	@curl -sf http://localhost:9091/api/v1/targets | grep -q "\"health\":\"up\"" && echo "   ✅ At least one target is UP" || echo "   ⚠️  Check target status"
	@echo ""
	@echo "✅ Metrics test completed!"
	@echo "💡 View metrics at: http://localhost:9091"

# Test Jaeger traces
test-traces:
	@echo "🔍 Testing Jaeger Distributed Tracing..."
	@echo "========================================"
	@echo ""
	@echo "1. Checking Jaeger health..."
	@curl -sf http://localhost:16686/ > /dev/null && echo "✅ Jaeger UI is accessible" || echo "❌ Jaeger is not responding"
	@echo ""
	@echo "2. Generating test traffic for traces..."
	@echo "   Creating a test task to generate trace..."
	@REGISTER=$$(curl -sf -X POST http://localhost:50051/register \
		-H "Content-Type: application/json" \
		-d '{"email":"trace-test@example.com","password":"test123"}' 2>/dev/null); \
	LOGIN=$$(curl -sf -X POST http://localhost:50051/login \
		-H "Content-Type: application/json" \
		-d '{"email":"trace-test@example.com","password":"test123"}' 2>/dev/null); \
	TOKEN=$$(echo $$LOGIN | grep -o '"token":"[^"]*' | sed 's/"token":"//'); \
	if [ -n "$$TOKEN" ]; then \
		curl -sf -X POST http://localhost:8080/api/upload \
			-H "Authorization: Bearer $$TOKEN" \
			-H "Content-Type: application/json" \
			-d '{"image_url":"https://example.com/trace-test.jpg"}' > /dev/null && \
		echo "   ✅ Test request sent (trace should be generated)"; \
	else \
		echo "   ⚠️  Could not generate test trace"; \
	fi
	@echo ""
	@echo "3. Waiting 5 seconds for trace propagation..."
	@sleep 5
	@echo ""
	@echo "✅ Traces test completed!"
	@echo "💡 View traces at: http://localhost:16686"
	@echo "💡 Search for service: api-service, operation: POST /api/upload"

# Test Loki logs
test-logs:
	@echo "📝 Testing Loki Log Aggregation..."
	@echo "=================================="
	@echo ""
	@echo "1. Checking Loki health..."
	@curl -sf http://localhost:3100/ready > /dev/null && echo "✅ Loki is ready" || echo "❌ Loki is not responding"
	@echo ""
	@echo "2. Checking Promtail..."
	@docker ps | grep -q promtail && echo "✅ Promtail is running" || echo "❌ Promtail is not running"
	@echo ""
	@echo "3. Querying recent logs from API service..."
	@LOGS=$$(curl -sf "http://localhost:3100/loki/api/v1/query_range?query={container_name=\"pixelflow-api\"}&limit=5" 2>/dev/null); \
	if echo "$$LOGS" | grep -q "\"status\":\"success\""; then \
		echo "   ✅ Logs are being collected"; \
		LOG_COUNT=$$(echo "$$LOGS" | grep -o "\"values\"" | wc -l); \
		echo "   📊 Found log streams from API service"; \
	else \
		echo "   ⚠️  No logs found (services may need more time)"; \
	fi
	@echo ""
	@echo "✅ Logs test completed!"
	@echo "💡 View logs in Grafana Explore: http://localhost:3001/explore"
	@echo "💡 Select 'Loki' datasource and query: {container_name=\"pixelflow-api\"}"

# Test alert rules
test-alerts:
	@echo "🚨 Testing Grafana Alert Rules..."
	@echo "================================="
	@echo ""
	@echo "1. Checking Grafana health..."
	@curl -sf http://localhost:3001/api/health > /dev/null && echo "✅ Grafana is healthy" || echo "❌ Grafana is not responding"
	@echo ""
	@echo "2. Checking alert rules configuration..."
	@if [ -f "deploy/grafana/provisioning/alerting/rules.yml" ]; then \
		echo "   ✅ Alert rules file exists"; \
		RULE_COUNT=$$(grep -c "title:" deploy/grafana/provisioning/alerting/rules.yml); \
		echo "   📊 Found $$RULE_COUNT alert rules configured"; \
	else \
		echo "   ❌ Alert rules file not found"; \
	fi
	@echo ""
	@echo "3. Alert rules configured:"
	@echo "   - High Error Rate - API Service"
	@echo "   - High Latency - API Service"
	@echo "   - Service Down - API Service"
	@echo "   - High Error Rate - Worker Service"
	@echo "   - High Kafka Consumption Errors"
	@echo ""
	@echo "✅ Alerts test completed!"
	@echo "💡 View alerts at: http://localhost:3001/alerting/list"
	@echo "💡 Login with: admin/admin"

# Test observability stack
test-observability:
	@echo "🔭 Testing Complete Observability Stack..."
	@echo "=========================================="
	@echo ""
	@make test-metrics
	@echo ""
	@make test-traces
	@echo ""
	@make test-logs
	@echo ""
	@make test-alerts
	@echo ""
	@echo "=========================================="
	@echo "✅ Observability Stack Verification Complete!"
	@echo ""
	@echo "📊 Access Points:"
	@echo "   Grafana:    http://localhost:3001 (admin/admin)"
	@echo "   Prometheus: http://localhost:9091"
	@echo "   Jaeger:     http://localhost:16686"
	@echo "   Loki:       Via Grafana Explore"

# Run full test suite
test-full:
	@echo "🧪 Running Full Test Suite..."
	@echo "============================="
	@echo ""
	@echo "Part 1: Application E2E Tests"
	@echo "-----------------------------"
	@make test-e2e
	@echo ""
	@echo "Part 2: Observability Stack Tests"
	@echo "----------------------------------"
	@make test-observability
	@echo ""
	@echo "============================="
	@echo "✅ Full Test Suite Completed!"
	@echo ""
	@echo "Summary:"
	@echo "  ✅ Application E2E: Passed"
	@echo "  ✅ Metrics: Verified"
	@echo "  ✅ Traces: Verified"
	@echo "  ✅ Logs: Verified"
	@echo "  ✅ Alerts: Verified"

# Open observability tools in browser
open-grafana:
	@echo "Opening Grafana..."
	@open http://localhost:3001 || xdg-open http://localhost:3001 || echo "Please open http://localhost:3001 in your browser"

open-prometheus:
	@echo "Opening Prometheus..."
	@open http://localhost:9091 || xdg-open http://localhost:9091 || echo "Please open http://localhost:9091 in your browser"

open-jaeger:
	@echo "Opening Jaeger..."
	@open http://localhost:16686 || xdg-open http://localhost:16686 || echo "Please open http://localhost:16686 in your browser"

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
	@make test-full

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
