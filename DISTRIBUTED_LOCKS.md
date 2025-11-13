# Distributed Locks with Redis

This document explains the distributed locking mechanism implemented for the LINE bot application.

## Overview

Distributed locks prevent race conditions and ensure data consistency across multiple instances of the application. This implementation uses Redis as the distributed lock store.

## Features

- **Atomic lock acquisition** using Redis SET NX command
- **Automatic expiration** via TTL to prevent deadlocks
- **Lock ownership verification** to prevent unauthorized unlocks
- **Retry logic** with configurable attempts and delays
- **Context support** for cancellation and timeouts
- **Lock extension** to renew TTL while holding the lock

## Architecture

### Lock Package

Location: `lock/lock.go`

The lock package provides:
- `DistributedLock` struct for managing individual locks
- `NewRedisClient()` for creating Redis connections
- `NewDistributedLock()` for creating lock instances
- `Lock()` and `Unlock()` for basic lock operations
- `LockWithRetry()` for lock acquisition with retry logic
- `Extend()` for renewing lock TTL
- `IsLocked()` for checking lock status
- `WithLock()` and `WithLockRetry()` helper functions

### Integration

The distributed lock is integrated into the application in:
- `app/app.go` - SetupRouter() initializes the Redis client
- `app/app.go` - Set() function uses locks to prevent user data race conditions

## Configuration

Redis configuration is stored in `config.yaml`:

```yaml
redis:
 host: localhost
 port: 6379
 password: ""
 db: 0
```

## Usage Examples

### Basic Lock/Unlock

```go
import (
	"context"
	"main/lock"
	"time"
)

// Create Redis client
client := lock.NewRedisClient(lock.Config{
	Host: "localhost",
	Port: "6379",
	Password: "",
	DB: 0,
})

// Create a distributed lock
distributedLock, err := lock.NewDistributedLock(client, "my-resource", 10*time.Second)
if err != nil {
	log.Fatal(err)
}

// Acquire lock
ctx := context.Background()
if err := distributedLock.Lock(ctx); err != nil {
	log.Fatal(err)
}

// Do critical work
// ...

// Release lock
if err := distributedLock.Unlock(ctx); err != nil {
	log.Error(err)
}
```

### Lock with Retry

```go
// Try to acquire lock with up to 5 retries, waiting 100ms between attempts
err := distributedLock.LockWithRetry(ctx, 5, 100*time.Millisecond)
if err != nil {
	if err == lock.ErrLockNotObtained {
		log.Println("Could not acquire lock after retries")
	} else {
		log.Fatal(err)
	}
}
defer distributedLock.Unlock(ctx)

// Critical section
// ...
```

### Using WithLock Helper

```go
import "main/lock"

err := lock.WithLock(ctx, client, "my-resource", 10*time.Second, func() error {
	// Critical section - lock is automatically acquired and released
	return doSomething()
})

if err != nil {
	log.Error(err)
}
```

### Using WithLockRetry Helper

```go
err := lock.WithLockRetry(
	ctx,
	client,
	"user:123",          // Lock key
	10*time.Second,      // TTL
	3,                   // Max retries
	100*time.Millisecond,// Retry delay
	func() error {
		// Critical section
		return updateUser()
	},
)
```

### Extending Lock TTL

```go
// Acquire lock
if err := distributedLock.Lock(ctx); err != nil {
	log.Fatal(err)
}
defer distributedLock.Unlock(ctx)

// Do some work
doSomeWork()

// Need more time? Extend the lock
if err := distributedLock.Extend(ctx, 20*time.Second); err != nil {
	log.Error("Failed to extend lock")
	return
}

// Continue with more work
doMoreWork()
```

## Implementation Details

### Lock Key Format

Locks are stored in Redis with keys in the format:
```
lock:<resource-name>
```

For example:
- `lock:user:123` - Lock for user with ID 123
- `lock:message:abc` - Lock for message with ID abc

### Lock Value

Each lock stores a random UUID value to ensure that only the lock owner can release it. This prevents scenarios where:
1. Process A acquires lock
2. Lock expires due to TTL
3. Process B acquires the same lock
4. Process A tries to release the lock (would fail because value doesn't match)

### TTL and Auto-Expiration

All locks have a TTL (Time To Live) to prevent deadlocks if a process crashes while holding a lock. The TTL should be:
- Long enough for the operation to complete normally
- Short enough to not block other processes for too long if the holder crashes

Recommended TTLs:
- Quick operations (< 1 second): 5-10 seconds
- Medium operations (1-5 seconds): 10-30 seconds
- Long operations (> 5 seconds): Consider breaking into smaller operations or using lock extension

### Retry Logic

The retry logic uses exponential backoff:
```go
for i := 0; i < maxRetries; i++ {
	if lock acquired {
		return success
	}
	wait(retryDelay)
}
return ErrLockNotObtained
```

### Lua Scripts

The implementation uses Lua scripts for atomic operations:

**Unlock Script:**
```lua
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end
```

**Extend Script:**
```lua
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("expire", KEYS[1], ARGV[2])
else
	return 0
end
```

## Error Handling

The package defines specific errors:

- `ErrLockNotObtained` - Lock could not be acquired (already held by another process)
- `ErrLockNotHeld` - Attempted to unlock/extend a lock not owned by this instance

Example:
```go
err := distributedLock.Lock(ctx)
if err == lock.ErrLockNotObtained {
	// Lock is held by another process
	// Decide whether to retry, wait, or fail
} else if err != nil {
	// Other error (e.g., Redis connection failed)
	log.Fatal(err)
}
```

## Testing

### Running Tests

```bash
# Run lock package tests
go test ./lock -v

# Run all tests with coverage
go test ./lock -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Test Environment

Tests are designed to pass even without Redis running. They will:
- Log informational messages when Redis is not available
- Pass successfully with appropriate messages
- Execute fully when Redis is running on localhost:6379

### With Redis Running

To run tests with actual Redis:

```bash
# Start Redis using Docker
docker run -d --name test-redis -p 6379:6379 redis:latest

# Run tests
go test ./lock -v

# Stop and remove Redis
docker stop test-redis
docker rm test-redis
```

## Production Deployment

### Redis Setup

1. **Single Redis Instance** (Development/Small Scale):
   ```bash
   docker run -d --name redis -p 6379:6379 redis:latest
   ```

2. **Redis Sentinel** (High Availability):
   - Use Redis Sentinel for automatic failover
   - Configure multiple Redis instances with sentinel monitoring

3. **Redis Cluster** (Large Scale):
   - Use Redis Cluster for horizontal scaling
   - Note: Some lock operations may need adjustment for cluster mode

### Configuration

Update `config.yaml` for production:
```yaml
redis:
 host: your-redis-host.example.com
 port: 6379
 password: "your-secure-password"
 db: 0
```

### Monitoring

Monitor these Redis metrics:
- Lock acquisition rate
- Lock hold duration
- Failed lock attempts
- Lock expiration events

### Best Practices

1. **Choose appropriate TTLs**
   - Balance between operation time and blocking time
   - Use lock extension for long-running operations

2. **Always use defer for unlocking**
   ```go
   lock.Lock(ctx)
   defer lock.Unlock(ctx)
   ```

3. **Handle lock acquisition failures gracefully**
   - Decide whether to retry, queue, or reject the operation
   - Log failures for monitoring

4. **Use unique lock keys**
   - Include entity type and ID: `user:123`, `order:456`
   - Avoid generic keys that could create bottlenecks

5. **Keep critical sections small**
   - Only lock what absolutely needs synchronization
   - Release locks as soon as possible

6. **Test without Redis**
   - Application should handle Redis unavailability gracefully
   - Consider fallback behavior or circuit breakers

## Use Cases in LINE Bot

### User Data Updates

The `Set()` function uses distributed locks to prevent race conditions when:
- Multiple webhooks arrive simultaneously for the same user
- User profile needs to be created/updated atomically
- Message history needs to be recorded in order

```go
lockKey := "user:" + userID
lock.WithLockRetry(ctx, redisClient, lockKey, 10*time.Second, 3, 100*time.Millisecond, func() error {
	// Check if user exists
	// Create user if needed
	// Record message
	return nil
})
```

### Potential Future Uses

1. **Message Rate Limiting**
   - Lock per user to enforce message rate limits
   - Prevent spam from single user

2. **Batch Processing**
   - Lock to ensure only one batch processor runs at a time
   - Coordinate multiple worker instances

3. **Configuration Updates**
   - Lock when updating shared configuration
   - Ensure consistent config across instances

## Troubleshooting

### Lock Not Being Acquired

**Symptoms:** `ErrLockNotObtained` error

**Possible Causes:**
- Another process is holding the lock
- Lock TTL is too long
- Too many concurrent requests for same resource

**Solutions:**
- Increase retry attempts or delay
- Reduce lock TTL
- Review lock usage patterns

### Locks Not Being Released

**Symptoms:** Operations timing out, high lock acquisition failures

**Possible Causes:**
- Application crashing while holding lock
- TTL too long
- Forgetting to call Unlock()

**Solutions:**
- Always use `defer lock.Unlock(ctx)`
- Reduce TTL to match operation duration
- Monitor for process crashes

### Redis Connection Failures

**Symptoms:** Various connection errors

**Possible Causes:**
- Redis server is down
- Network issues
- Authentication problems

**Solutions:**
- Verify Redis server is running
- Check network connectivity
- Verify Redis password in config
- Implement connection retry logic

## References

- [Redis SET Command](https://redis.io/commands/set)
- [Redis Distributed Locks](https://redis.io/topics/distlock)
- [Redlock Algorithm](https://redis.io/topics/distlock#the-redlock-algorithm)

## License

This implementation is part of the LINE bot project and follows the same license.
