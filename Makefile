# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Testing commands
test:
	go test ./... -v

test-short:
	go test ./... -v -short

test-coverage:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race:
	go test ./... -v -race

test-benchmark:
	go test ./... -v -bench=.

test-models:
	go test ./internal/models/... -v

test-repository:
	go test ./internal/repository/... -v

test-service:
	go test ./internal/service/... -v

test-config:
	go test ./internal/config/... -v

test-logger:
	go test ./internal/logger/... -v

test-migrations:
	go test ./internal/migrations/... -v

test-unit:
	go test ./internal/models/... ./internal/config/... ./internal/logger/... -v

test-integration:
	go test ./internal/repository/... ./internal/service/... ./internal/migrations/... -v

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Docker commands
docker-build:
	docker build -t tickets-api deployments/

docker-run:
	cd deployments && docker compose up -d

docker-stop:
	cd deployments && docker compose down

docker-logs:
	cd deployments && docker compose logs -f

docker-clean:
	cd deployments && docker compose down -v
	docker system prune -f

# Development with Docker
dev-docker:
	cd deployments && docker compose up --build

# Start with automatic database reset
start-fresh:
	chmod +x deployments/start-with-reset.sh
	./deployments/start-with-reset.sh

# Database commands
db-migrate:
	./bin/migrate -command up

db-connect:
	cd deployments && docker-compose exec postgres psql -U postgres -d tickets_db

db-reset:
	cd deployments && docker-compose down
	cd deployments && docker volume rm tickets_postgres_data 2>/dev/null || true
	cd deployments && docker-compose up -d postgres
	@echo "Database reset completed. Waiting for PostgreSQL to be ready..."
	@echo "You can now start your application with: make docker-run"

db-reset-full:
	cd deployments && docker-compose down -v
	cd deployments && docker-compose up -d postgres
	@echo "Full database reset completed. Waiting for PostgreSQL to be ready..."
	@echo "You can now start your application with: make docker-run"

# Protocol Buffer commands
proto:
	@echo "Cleaning up existing generated files..."
	rm -f proto/*.pb.go
	rm -f api/*.pb.go
	@echo "Generating Protocol Buffer code..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/tickets.proto
	mkdir -p api
	mv proto/*.pb.go api/
	@echo "Protocol Buffer code generated successfully in api/ directory"

proto-clean:
	rm -f proto/*.pb.go
	rm -f api/*.pb.go
	rmdir api 2>/dev/null || true

# Migration commands
migrate-build:
	go build -o bin/migrate cmd/migrate/main.go

migrate-up:
	./bin/migrate -command up

migrate-down:
	./bin/migrate -command down -steps 1

migrate-status:
	./bin/migrate -command status

migrate-create:
	./bin/migrate -command create -name $(name)

migrate-reset:
	./bin/migrate -command down -steps 999
	./bin/migrate -command up

# Show help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo ""
	@echo "Docker commands:"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Start services with Docker Compose"
	@echo "  docker-stop   - Stop services"
	@echo "  docker-logs   - View logs"
	@echo "  docker-clean  - Clean Docker resources"
	@echo "  dev-docker    - Development with Docker"
	@echo "  start-fresh   - Start application with automatic database reset"
	@echo ""
	@echo "Database commands:"
	@echo "  db-migrate    - Run database migrations"
	@echo "  db-connect    - Connect to database"
	@echo "  db-reset      - Reset database to initial state"
	@echo "  db-reset-full - Full database reset (removes all volumes)"
	@echo ""
	@echo "Migration commands:"
	@echo "  migrate-build  - Build migration tool"
	@echo "  migrate-up     - Apply all pending migrations"
	@echo "  migrate-down   - Rollback last migration"
	@echo "  migrate-status - Show migration status"
	@echo "  migrate-create - Create new migration (name=<migration_name>)"
	@echo "  migrate-reset  - Reset all migrations"
	@echo ""
	@echo "Protocol Buffer commands:"
	@echo "  proto          - Generate Go code from .proto files"
	@echo "  proto-clean    - Remove generated .pb.go files"
	@echo ""
	@echo "Testing commands:"
	@echo "  test           - Run all tests with verbose output"
	@echo "  test-short     - Run tests with short flag"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-race      - Run tests with race detection"
	@echo "  test-benchmark - Run benchmark tests"
	@echo "  test-models    - Run model tests only"
	@echo "  test-repository - Run repository tests only"
	@echo "  test-service   - Run service tests only"
	@echo "  test-config    - Run config tests only"
	@echo "  test-logger    - Run logger tests only"
	@echo "  test-migrations - Run migration tests only"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  help          - Show this help"

.PHONY: build run test test-short test-coverage test-race test-benchmark test-models test-repository test-service test-config test-logger test-migrations test-unit test-integration clean deps lint fmt docker-build docker-run docker-stop docker-logs docker-clean dev-docker db-migrate db-connect proto proto-clean migrate-build migrate-up migrate-down migrate-status migrate-create migrate-reset help