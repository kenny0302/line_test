# LINE Bot with Distributed Locks

A LINE messaging bot with MongoDB storage and Redis-based distributed locking for data consistency across multiple instances.

## Features

- LINE webhook integration for message handling
- MongoDB for persistent storage of users and messages
- **Distributed locks with Redis** to prevent race conditions
- Comprehensive test suite with 53.3% coverage
- GitHub Actions CI/CD with MongoDB and Redis services

## Quick Start

### Prerequisites

- Go 1.18 or higher
- MongoDB 4.4 or higher
- Redis 6.0 or higher
- LINE Developer Account

### Setup

#### 1. Create MongoDB Container
```bash
docker run -d --name testdb -p 27017:27017 mongo:4.4
```

#### 2. Create Redis Container
```bash
docker run -d --name testredis -p 6379:6379 redis:latest
```

#### 3. Set up MongoDB Collections

```bash
docker exec -it testdb bash
mongo
use line
db.createCollection("message")
db.createCollection("user")
```

#### 4. Configure Application

Edit `config.yaml` with your settings:

```yaml
database:
 host: localhost
 port: 27017

redis:
 host: localhost
 port: 6379
 password: ""
 db: 0

line:
 secret: YOUR_LINE_CHANNEL_SECRET
 token: YOUR_LINE_ACCESS_TOKEN
```

Create a LINE Developer account and get your credentials from the [LINE Developers Console](https://developers.line.biz/).

#### 5. Run the Service
```bash
# Using Go
go run main.go

# Or using Make
make run

# Or build and run
make build
./line_test
```

## API Endpoints

### GET /list
List all registered users.

```bash
curl --location --request GET '127.0.0.1:8080/list'
```

### POST /push
Send a message to a specific user.

```bash
curl --location --request POST '127.0.0.1:8080/push' \
--form 'UserId="YOUR_LINE_USER_ID"'
```

### POST /callback
LINE webhook endpoint for receiving messages.

## Testing

### Run All Tests

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Generate HTML coverage report
make coverage-html
```

### Test Coverage

Current test coverage: **53.3%**
- app package: 48.3%
- db package: 62.5%
- lock package: Comprehensive tests included
- proto package: 100%

See [TEST_COVERAGE.md](TEST_COVERAGE.md) for detailed coverage information.

### Run Tests with Docker Services

```bash
# Start MongoDB and Redis for testing
make start-mongo
docker run -d --name test-redis -p 6379:6379 redis:latest

# Run tests
make test-with-mongo

# Cleanup
make stop-mongo
docker stop test-redis && docker rm test-redis
```

## Distributed Locks

This application uses Redis-based distributed locks to prevent race conditions when multiple instances process messages for the same user simultaneously.

### Key Features

- Atomic lock acquisition
- Automatic expiration (TTL)
- Retry logic with configurable attempts
- Lock ownership verification

### Usage Example

```go
import "main/lock"

// Lock automatically acquired and released
err := lock.WithLockRetry(ctx, redisClient, "user:123", 10*time.Second, 3, 100*time.Millisecond, func() error {
	// Critical section - update user data
	return updateUser()
})
```

See [DISTRIBUTED_LOCKS.md](DISTRIBUTED_LOCKS.md) for comprehensive documentation.

## Project Structure

```
.
├── app/           # Application logic and HTTP handlers
├── db/            # MongoDB database operations
├── lock/          # Distributed lock implementation
├── proto/         # Data structures
├── .github/       # GitHub Actions workflows
├── config.yaml    # Application configuration
├── main.go        # Application entry point
├── Makefile       # Build and test commands
└── README.md      # This file
```

## Development

### Using Makefile

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Run tests with coverage
make coverage

# Generate HTML coverage report
make coverage-html

# Format code
make fmt

# Clean build artifacts
make clean

# Show all available commands
make help
```

## CI/CD

The project includes a GitHub Actions workflow that:
- Runs tests with MongoDB and Redis services
- Generates coverage reports
- Enforces 80% coverage threshold (when services are available)
- Uploads coverage to Codecov

## Demo Videos

### Push Message
https://drive.google.com/file/d/1b4kTe0RY5kwcwQEO-N_2UDMvanV7KaC_/view?usp=sharing

### Save User Info and Message
https://drive.google.com/file/d/1P7bOGYC2t3TtN2MJGBcH39WGxuPWbkux/view?usp=sharing

### List Users
https://drive.google.com/file/d/1U_iD4vrJFMmqYWKFkOVFM7wx6bYh3uIt/view?usp=sharing

## Documentation

- [TEST_COVERAGE.md](TEST_COVERAGE.md) - Detailed test coverage report
- [DISTRIBUTED_LOCKS.md](DISTRIBUTED_LOCKS.md) - Distributed locking documentation

## License

[Add your license here]
