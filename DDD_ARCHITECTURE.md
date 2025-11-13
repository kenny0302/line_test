# Domain-Driven Design (DDD) Architecture

This document describes the Domain-Driven Design architecture refactoring for the LINE Bot application.

## Architecture Overview

The application is structured into four main layers following DDD principles:

```
├── cmd/server/             # Application entry point
├── internal/
│   ├── domain/            # Domain Layer (Business Logic)
│   ├── application/       # Application Layer (Use Cases)
│   ├── infrastructure/    # Infrastructure Layer (External Services)
│   └── presentation/      # Presentation Layer (HTTP)
└── pkg/                   # Shared packages
```

## Layers

### 1. Domain Layer (`internal/domain/`)

The core business logic with no external dependencies.

#### User Aggregate
- **Entities**: `User` - Represents a LINE bot user
- **Value Objects**:
  - `UserID` - Unique identifier with validation
  - `DisplayName` - User's display name (max 100 chars)
  - `PictureURL` - Profile picture URL
  - `StatusMessage` - User's status
  - `Language` - Language preference
- **Repository Interface**: Defines persistence contract

#### Message Aggregate
- **Entities**: `Message` - Represents a chat message
- **Value Objects**:
  - `MessageID` - Unique identifier
  - `Content` - Message text (max 5000 chars)
  - `UserReference` - Reference to user
- **Repository Interface**: Defines persistence contract

#### Shared
- **Lock Service Interface**: Distributed locking contract

### 2. Application Layer (`internal/application/`)

Coordinates business workflows and implements use cases.

#### Use Cases
- `RegisterUserUseCase` - Register/update user with distributed locking
- `SaveMessageUseCase` - Save user messages
- `ListUsersUseCase` - Retrieve all users

#### Services
- `LINEBotService` - Coordinates webhook processing, user registration, and messaging

#### DTOs (Data Transfer Objects)
- `UserDTO` - User data transfer
- `MessageDTO` - Message data transfer
- Request/Response models

### 3. Infrastructure Layer (`internal/infrastructure/`)

Implements interfaces defined in domain layer.

#### Persistence
- `UserRepository` - MongoDB implementation of user.Repository
- `MessageRepository` - MongoDB implementation of message.Repository
- Connection management

#### Lock
- `RedisLockService` - Distributed lock implementation

#### Messaging
- `LINEClient` - LINE Bot API integration

#### Config
- Configuration management

### 4. Presentation Layer (`internal/presentation/`)

HTTP handlers and routing.

#### HTTP Handlers
- `UserHandler` - User-related endpoints
- `WebhookHandler` - LINE webhook processing

#### Middleware
- Error handling
- Logging
- Authentication

## Key DDD Concepts Implemented

### Entities
Objects with identity that can change over time:
- `User` - Has lifecycle methods (UpdateProfile, UpdateStatusMessage, etc.)
- `Message` - Immutable once created

### Value Objects
Immutable objects defined by their attributes:
- `UserID`, `DisplayName`, `PictureURL`
- `MessageID`, `Content`
- Encapsulate validation logic
- Provide type safety

### Aggregates
Consistency boundaries:
- **User Aggregate**: User + related operations
- **Message Aggregate**: Message + user reference

### Repositories
Abstract persistence:
- Domain defines interfaces
- Infrastructure provides implementations
- Allows switching storage without changing domain

### Use Cases
Application-specific business rules:
- Each use case represents a single business operation
- Coordinates between domain entities and repositories
- Enforces business rules and constraints

## Benefits of This Architecture

### 1. Separation of Concerns
- Business logic isolated in domain layer
- Infrastructure details separated
- Easy to understand and maintain

### 2. Testability
- Domain layer has no dependencies - easy to unit test
- Application layer can be tested with mocks
- Infrastructure can be integration tested

### 3. Flexibility
- Can swap MongoDB for PostgreSQL without changing domain
- Can add new use cases without modifying existing code
- Can change presentation layer (REST to GraphQL) easily

### 4. Type Safety
- Value objects prevent invalid data
- Compiler catches many errors
- Self-documenting code

### 5. Business Logic Clarity
- Domain entities contain business rules
- Use cases make workflows explicit
- Easy for non-technical stakeholders to understand

## Data Flow

```
HTTP Request
    ↓
Handler (Presentation)
    ↓
Use Case (Application)
    ↓
Domain Entity + Repository
    ↓
Repository Implementation (Infrastructure)
    ↓
MongoDB/Redis
```

## Example: Register User Flow

1. **HTTP Request** arrives at `/callback`
2. **WebhookHandler** (Presentation) parses request
3. **LINEBotService** (Application) coordinates workflow
4. **RegisterUserUseCase** (Application) executes business logic:
   - Creates value objects (UserID, DisplayName, etc.)
   - Acquires distributed lock
   - Checks if user exists
   - Creates or updates User entity
   - Saves via repository
5. **UserRepository** (Infrastructure) persists to MongoDB
6. **Response** returns to client

## Migration Path

### Current Structure → DDD Structure

```
Old:                          New:
app/app.go              →    presentation/http/handler/
db/db.go                →    infrastructure/persistence/mongodb/
proto/proto.go          →    domain/{user,message}/
lock/lock.go            →    infrastructure/lock/
```

## Testing Strategy

### Unit Tests
- Domain entities and value objects
- Use cases with mocked repositories
- Pure business logic

### Integration Tests
- Repository implementations with test database
- HTTP handlers with test server
- End-to-end workflows

### Contract Tests
- Verify infrastructure matches domain interfaces
- Ensure DTOs map correctly to entities

## Future Enhancements

### Domain Events
Add event sourcing for audit trail:
```go
type UserRegistered struct {
    UserID    string
    Timestamp time.Time
}
```

### CQRS (Command Query Responsibility Segregation)
Separate read and write models:
- Commands: RegisterUser, SaveMessage
- Queries: ListUsers, GetUserMessages

### Domain Services
Extract complex business logic:
```go
type UserRegistrationService struct {
    // Complex registration logic
}
```

## Best Practices

### 1. Keep Domain Pure
- No external dependencies in domain layer
- All business rules in entities/value objects
- Use interfaces for external concerns

### 2. Thin Controllers
- Handlers only parse requests and format responses
- Delegate to use cases
- No business logic

### 3. Fat Models
- Entities contain behavior, not just data
- Methods enforce invariants
- Self-validating value objects

### 4. Use Cases are Workflows
- One use case = one business operation
- Coordinate between entities and services
- Handle transactions and locking

### 5. Repositories are Collections
- Think of them as in-memory collections
- Hide persistence details
- Return domain entities, not DTOs

## References

- [Domain-Driven Design by Eric Evans](https://www.domainlanguage.com/ddd/)
- [Clean Architecture by Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)

## Implementation Status

### Completed ✅
- Domain layer (entities, value objects, repository interfaces)
- Application layer (use cases, services, DTOs)
- Infrastructure layer (MongoDB repositories - partial)

### In Progress 🚧
- Infrastructure layer (complete implementations)
- Presentation layer (HTTP handlers)
- Migration from old structure

### Pending 📋
- Comprehensive tests for new structure
- Documentation updates
- Performance optimization
- Monitoring and observability
