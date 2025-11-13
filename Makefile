.PHONY: test coverage coverage-html clean build run

# Build the application
build:
	go build -o line_test .

# Run the application
run: build
	./line_test

# Run all tests
test:
	go test ./... -v

# Run tests with coverage
coverage:
	go test ./app ./proto ./db ./lock -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

# Generate HTML coverage report
coverage-html: coverage
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detector
test-race:
	go test ./... -race -v

# Run tests with MongoDB (requires MongoDB running on localhost:27017)
test-with-mongo:
	@echo "Ensure MongoDB is running on localhost:27017"
	@go test ./... -v -coverprofile=coverage.out -covermode=atomic
	@go tool cover -func=coverage.out | tail -1

# Clean build artifacts and test files
clean:
	rm -f line_test
	rm -f coverage.out coverage.html coverage_detail.txt
	rm -f app/config.yaml config.yaml

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

# Run all checks (format, lint, test)
check: fmt test

# Start MongoDB in Docker for testing
start-mongo:
	docker run -d --name line-test-mongo -p 27017:27017 mongo:4.4

# Stop and remove MongoDB container
stop-mongo:
	docker stop line-test-mongo || true
	docker rm line-test-mongo || true

# Start Redis in Docker for testing
start-redis:
	docker run -d --name line-test-redis -p 6379:6379 redis:latest

# Stop and remove Redis container
stop-redis:
	docker stop line-test-redis || true
	docker rm line-test-redis || true

# Start all services (MongoDB + Redis)
start-services: start-mongo start-redis
	@echo "Waiting for services to be ready..."
	@sleep 5

# Stop all services
stop-services: stop-mongo stop-redis

# Full test with all services in Docker
test-full: start-services
	@echo "Running tests with MongoDB and Redis..."
	@$(MAKE) test-with-mongo
	@$(MAKE) stop-services

# Show help
help:
	@echo "Available targets:"
	@echo "  build            - Build the application"
	@echo "  run              - Run the application"
	@echo "  test             - Run all tests"
	@echo "  coverage         - Run tests with coverage"
	@echo "  coverage-html    - Generate HTML coverage report"
	@echo "  test-race        - Run tests with race detector"
	@echo "  test-with-mongo  - Run tests with MongoDB (requires MongoDB running)"
	@echo "  clean            - Clean build artifacts"
	@echo "  deps             - Install dependencies"
	@echo "  fmt              - Format code"
	@echo "  lint             - Run linter"
	@echo "  check            - Run format and tests"
	@echo "  start-mongo      - Start MongoDB in Docker"
	@echo "  stop-mongo       - Stop MongoDB Docker container"
	@echo "  start-redis      - Start Redis in Docker"
	@echo "  stop-redis       - Stop Redis Docker container"
	@echo "  start-services   - Start MongoDB and Redis in Docker"
	@echo "  stop-services    - Stop MongoDB and Redis containers"
	@echo "  test-full        - Run full test suite with all services in Docker"
	@echo "  help             - Show this help message"
