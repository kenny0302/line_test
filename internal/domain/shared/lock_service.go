package shared

import (
	"context"
	"time"
)

// LockService defines the interface for distributed locking
type LockService interface {
	// AcquireLock attempts to acquire a lock for a given resource
	AcquireLock(ctx context.Context, resource string, ttl time.Duration) (Lock, error)

	// AcquireLockWithRetry attempts to acquire a lock with retry logic
	AcquireLockWithRetry(ctx context.Context, resource string, ttl time.Duration, maxRetries int, retryDelay time.Duration) (Lock, error)
}

// Lock represents an acquired distributed lock
type Lock interface {
	// Release releases the lock
	Release(ctx context.Context) error

	// Extend extends the TTL of the lock
	Extend(ctx context.Context, ttl time.Duration) error

	// IsHeld checks if the lock is still held
	IsHeld(ctx context.Context) (bool, error)
}
