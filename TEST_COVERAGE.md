# Test Coverage Report

## Overview

This document provides an overview of the test coverage for the LINE Bot project.

## Current Coverage

- **Overall Coverage**: 53.3%
- **app package**: 48.3%
- **db package**: 62.5%
- **proto package**: 100% (no testable statements, only struct definitions)

## Test Structure

### Package: app
Location: `app/app_test.go`

Tests include:
- Router setup and configuration
- HTTP endpoint handlers (GET /list, POST /callback, POST /push)
- Request/response validation
- Error handling scenarios
- Configuration loading
- Multiple request scenarios

### Package: db
Location: `db/db_test.go`

Tests include:
- MongoDB connection establishment
- User data retrieval with various filters
- User data insertion
- Message data insertion
- Connection context management
- Error handling for invalid hosts/ports

### Package: proto
Location: `proto/proto_test.go`

Tests include:
- Struct initialization and field validation
- JSON marshaling/unmarshaling
- Data structure integrity

## Test Execution

Run all tests:
```bash
go test ./... -v
```

Run tests with coverage:
```bash
go test ./app ./proto ./db -coverprofile=coverage.out -covermode=atomic
```

View coverage report:
```bash
go tool cover -html=coverage.out
```

View coverage summary:
```bash
go tool cover -func=coverage.out
```

## Coverage Details by Function

### app package
| Function | Coverage |
|----------|----------|
| SetupRouter | 91.7% |
| GetBotData | 26.7% |
| Set | 35.7% |
| ListUser | 83.3% |
| PushMessage | 83.3% |

### db package
| Function | Coverage |
|----------|----------|
| Connect | 55.6% |
| GetUser | 47.4% |
| SetUser | 80.0% |
| SetMessage | 80.0% |

## Limitations and Challenges

### External Dependencies
The project has significant dependencies on external services that affect test coverage:

1. **MongoDB**: Many functions require a running MongoDB instance. Without it:
   - Connection tests timeout after 30 seconds
   - Data operations cannot be fully tested
   - Integration tests are limited

2. **LINE Bot API**:
   - Webhook signature validation requires valid LINE credentials
   - Message sending requires authenticated API access
   - Profile fetching needs valid tokens

### Current Test Approach
Tests are designed to:
- Validate function structure and interfaces
- Test error handling paths
- Verify request/response formats
- Handle both success and failure scenarios gracefully

### Test Execution Time
Due to MongoDB connection timeouts (30 seconds each), full test suite execution takes approximately 7-8 minutes.

## Recommendations for Improving Coverage

### To reach 80%+ coverage:

1. **Set up MongoDB for Testing**
   ```bash
   # Using Docker
   docker run -d --name test-mongo -p 27017:27017 mongo:4.4

   # Run tests
   go test ./... -coverprofile=coverage.out
   ```

2. **Implement Dependency Injection**
   - Refactor code to accept interfaces instead of concrete types
   - Create mock implementations for MongoDB and LINE Bot SDK
   - Use libraries like `testify/mock` or `gomock`

3. **Add Integration Tests**
   - Create test fixtures with sample data
   - Use testcontainers-go for ephemeral MongoDB instances
   - Implement proper test setup/teardown

4. **Optimize Test Timeouts**
   - Configure shorter connection timeouts for tests
   - Use context with shorter deadlines
   - Skip long-running tests in CI environments

## CI/CD Integration

Example GitHub Actions workflow:
```yaml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    services:
      mongodb:
        image: mongo:4.4
        ports:
          - 27017:27017

    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: '1.18'

      - name: Run tests
        run: go test ./... -coverprofile=coverage.out -covermode=atomic

      - name: Upload coverage
        uses: codecov/codecov-action@v2
        with:
          files: ./coverage.out
```

## Test Files

- `app/app_test.go` - Application logic tests
- `db/db_test.go` - Database operations tests
- `proto/proto_test.go` - Data structure tests
- `main.go` - Entry point (minimal, not directly tested)

## Conclusion

The current test suite provides a solid foundation with 53.3% coverage. The primary limitation is the lack of available external services (MongoDB, LINE API) during test execution. With proper test infrastructure (mocking or actual services), coverage can be increased to 80%+.
