package lock

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	// ErrLockNotObtained is returned when a lock cannot be obtained
	ErrLockNotObtained = errors.New("lock not obtained")
	// ErrLockNotHeld is returned when trying to unlock a lock that is not held
	ErrLockNotHeld = errors.New("lock not held")
)

// DistributedLock represents a distributed lock using Redis
type DistributedLock struct {
	client *redis.Client
	key    string
	value  string
	ttl    time.Duration
}

// Config holds configuration for Redis connection
type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// NewRedisClient creates a new Redis client
func NewRedisClient(config Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password,
		DB:       config.DB,
	})
}

// NewDistributedLock creates a new distributed lock
func NewDistributedLock(client *redis.Client, key string, ttl time.Duration) (*DistributedLock, error) {
	value, err := generateRandomValue()
	if err != nil {
		return nil, err
	}

	return &DistributedLock{
		client: client,
		key:    "lock:" + key,
		value:  value,
		ttl:    ttl,
	}, nil
}

// Lock attempts to acquire the lock
func (l *DistributedLock) Lock(ctx context.Context) error {
	// Use SET NX (set if not exists) with expiration
	ok, err := l.client.SetNX(ctx, l.key, l.value, l.ttl).Result()
	if err != nil {
		return err
	}

	if !ok {
		return ErrLockNotObtained
	}

	return nil
}

// LockWithRetry attempts to acquire the lock with retry logic
func (l *DistributedLock) LockWithRetry(ctx context.Context, maxRetries int, retryDelay time.Duration) error {
	for i := 0; i < maxRetries; i++ {
		err := l.Lock(ctx)
		if err == nil {
			return nil
		}

		if err != ErrLockNotObtained {
			return err
		}

		// Wait before retrying
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(retryDelay):
			continue
		}
	}

	return ErrLockNotObtained
}

// Unlock releases the lock
func (l *DistributedLock) Unlock(ctx context.Context) error {
	// Lua script to ensure we only delete the lock if we own it
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	result, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Result()
	if err != nil {
		return err
	}

	if result == int64(0) {
		return ErrLockNotHeld
	}

	return nil
}

// Extend extends the TTL of the lock
func (l *DistributedLock) Extend(ctx context.Context, ttl time.Duration) error {
	// Lua script to extend TTL only if we own the lock
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("expire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`

	result, err := l.client.Eval(ctx, script, []string{l.key}, l.value, int(ttl.Seconds())).Result()
	if err != nil {
		return err
	}

	if result == int64(0) {
		return ErrLockNotHeld
	}

	l.ttl = ttl
	return nil
}

// IsLocked checks if the lock is currently held by this instance
func (l *DistributedLock) IsLocked(ctx context.Context) (bool, error) {
	value, err := l.client.Get(ctx, l.key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return value == l.value, nil
}

// generateRandomValue generates a random value for the lock
func generateRandomValue() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// WithLock executes a function while holding a distributed lock
func WithLock(ctx context.Context, client *redis.Client, key string, ttl time.Duration, fn func() error) error {
	lock, err := NewDistributedLock(client, key, ttl)
	if err != nil {
		return err
	}

	if err := lock.Lock(ctx); err != nil {
		return err
	}
	defer lock.Unlock(ctx)

	return fn()
}

// WithLockRetry executes a function while holding a distributed lock with retry logic
func WithLockRetry(ctx context.Context, client *redis.Client, key string, ttl time.Duration, maxRetries int, retryDelay time.Duration, fn func() error) error {
	lock, err := NewDistributedLock(client, key, ttl)
	if err != nil {
		return err
	}

	if err := lock.LockWithRetry(ctx, maxRetries, retryDelay); err != nil {
		return err
	}
	defer lock.Unlock(ctx)

	return fn()
}
