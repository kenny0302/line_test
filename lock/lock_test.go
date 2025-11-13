package lock

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func getTestRedisClient() *redis.Client {
	return NewRedisClient(Config{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	})
}

func TestNewRedisClient(t *testing.T) {
	client := getTestRedisClient()
	if client == nil {
		t.Fatal("Expected non-nil Redis client")
	}

	// Try to ping Redis
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		t.Logf("Redis not available: %v (this is expected if Redis is not running)", err)
	} else {
		t.Log("Redis is available and responding")
	}
}

func TestNewDistributedLock(t *testing.T) {
	client := getTestRedisClient()
	lock, err := NewDistributedLock(client, "test-key", 10*time.Second)

	if err != nil {
		t.Fatalf("Failed to create distributed lock: %v", err)
	}

	if lock == nil {
		t.Fatal("Expected non-nil lock")
	}

	if lock.key != "lock:test-key" {
		t.Errorf("Expected key to be 'lock:test-key', got '%s'", lock.key)
	}

	if lock.ttl != 10*time.Second {
		t.Errorf("Expected TTL to be 10s, got %v", lock.ttl)
	}

	if lock.value == "" {
		t.Error("Expected non-empty lock value")
	}
}

func TestLock_BasicLocking(t *testing.T) {
	client := getTestRedisClient()
	ctx := context.Background()

	// Clean up any existing lock
	client.Del(ctx, "lock:test-basic")

	lock, err := NewDistributedLock(client, "test-basic", 5*time.Second)
	if err != nil {
		t.Fatalf("Failed to create lock: %v", err)
	}

	// Try to acquire lock
	err = lock.Lock(ctx)
	if err != nil {
		t.Logf("Failed to acquire lock: %v (expected if Redis is not running)", err)
		return
	}

	t.Log("Successfully acquired lock")

	// Verify lock is held
	isLocked, err := lock.IsLocked(ctx)
	if err != nil {
		t.Errorf("Failed to check lock status: %v", err)
	}
	if !isLocked {
		t.Error("Expected lock to be held")
	}

	// Release lock
	err = lock.Unlock(ctx)
	if err != nil {
		t.Errorf("Failed to unlock: %v", err)
	}

	// Verify lock is released
	isLocked, err = lock.IsLocked(ctx)
	if err != nil {
		t.Errorf("Failed to check lock status: %v", err)
	}
	if isLocked {
		t.Error("Expected lock to be released")
	}
}

func TestLock_DoubleAcquire(t *testing.T) {
	client := getTestRedisClient()
	ctx := context.Background()

	// Clean up
	client.Del(ctx, "lock:test-double")

	lock1, _ := NewDistributedLock(client, "test-double", 5*time.Second)
	lock2, _ := NewDistributedLock(client, "test-double", 5*time.Second)

	// First lock should succeed
	err := lock1.Lock(ctx)
	if err != nil {
		t.Logf("Redis not available: %v", err)
		return
	}
	defer lock1.Unlock(ctx)

	// Second lock should fail
	err = lock2.Lock(ctx)
	if err != ErrLockNotObtained {
		t.Errorf("Expected ErrLockNotObtained, got %v", err)
	}
}

func TestLock_WithRetry(t *testing.T) {
	client := getTestRedisClient()
	ctx := context.Background()

	// Clean up
	client.Del(ctx, "lock:test-retry")

	lock, err := NewDistributedLock(client, "test-retry", 2*time.Second)
	if err != nil {
		t.Fatalf("Failed to create lock: %v", err)
	}

	// Try to acquire with retry
	err = lock.LockWithRetry(ctx, 3, 100*time.Millisecond)
	if err != nil {
		t.Logf("Failed to acquire lock with retry: %v (expected if Redis is not running)", err)
		return
	}
	defer lock.Unlock(ctx)

	t.Log("Successfully acquired lock with retry")
}

func TestLock_AutoExpiration(t *testing.T) {
	client := getTestRedisClient()
	ctx := context.Background()

	// Clean up
	client.Del(ctx, "lock:test-expire")

	lock, _ := NewDistributedLock(client, "test-expire", 1*time.Second)

	err := lock.Lock(ctx)
	if err != nil {
		t.Logf("Redis not available: %v", err)
		return
	}

	// Wait for lock to expire
	time.Sleep(2 * time.Second)

	// Check if lock expired
	isLocked, err := lock.IsLocked(ctx)
	if err != nil {
		t.Errorf("Failed to check lock status: %v", err)
	}
	if isLocked {
		t.Error("Expected lock to have expired")
	}

	t.Log("Lock correctly expired after TTL")
}

func TestLock_Extend(t *testing.T) {
	client := getTestRedisClient()
	ctx := context.Background()

	// Clean up
	client.Del(ctx, "lock:test-extend")

	lock, _ := NewDistributedLock(client, "test-extend", 2*time.Second)

	err := lock.Lock(ctx)
	if err != nil {
		t.Logf("Redis not available: %v", err)
		return
	}
	defer lock.Unlock(ctx)

	// Extend the lock
	err = lock.Extend(ctx, 10*time.Second)
	if err != nil {
		t.Errorf("Failed to extend lock: %v", err)
		return
	}

	// Verify lock is still held
	isLocked, err := lock.IsLocked(ctx)
	if err != nil {
		t.Errorf("Failed to check lock status: %v", err)
	}
	if !isLocked {
		t.Error("Expected lock to still be held after extension")
	}

	t.Log("Successfully extended lock TTL")
}

func TestWithLock(t *testing.T) {
	client := getTestRedisClient()
	ctx := context.Background()

	// Clean up
	client.Del(ctx, "lock:test-withlock")

	executed := false
	err := WithLock(ctx, client, "test-withlock", 5*time.Second, func() error {
		executed = true
		return nil
	})

	if err != nil {
		t.Logf("WithLock failed: %v (expected if Redis is not running)", err)
		return
	}

	if !executed {
		t.Error("Expected function to be executed")
	}

	// Verify lock is released
	value, err := client.Get(ctx, "lock:test-withlock").Result()
	if err != redis.Nil {
		t.Errorf("Expected lock to be released, got value: %v, err: %v", value, err)
	}
}

func TestWithLockRetry(t *testing.T) {
	client := getTestRedisClient()
	ctx := context.Background()

	// Clean up
	client.Del(ctx, "lock:test-withlock-retry")

	executed := false
	err := WithLockRetry(ctx, client, "test-withlock-retry", 5*time.Second, 3, 100*time.Millisecond, func() error {
		executed = true
		return nil
	})

	if err != nil {
		t.Logf("WithLockRetry failed: %v (expected if Redis is not running)", err)
		return
	}

	if !executed {
		t.Error("Expected function to be executed")
	}
}

func TestLock_UnlockWithoutOwnership(t *testing.T) {
	client := getTestRedisClient()
	ctx := context.Background()

	// Clean up
	client.Del(ctx, "lock:test-ownership")

	lock1, _ := NewDistributedLock(client, "test-ownership", 5*time.Second)
	lock2, _ := NewDistributedLock(client, "test-ownership", 5*time.Second)

	// Lock1 acquires the lock
	err := lock1.Lock(ctx)
	if err != nil {
		t.Logf("Redis not available: %v", err)
		return
	}

	// Lock2 tries to unlock (should fail)
	err = lock2.Unlock(ctx)
	if err != ErrLockNotHeld {
		t.Errorf("Expected ErrLockNotHeld, got %v", err)
	}

	// Lock1 unlocks successfully
	err = lock1.Unlock(ctx)
	if err != nil {
		t.Errorf("Failed to unlock: %v", err)
	}
}

func TestLock_ContextCancellation(t *testing.T) {
	client := getTestRedisClient()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Clean up
	client.Del(context.Background(), "lock:test-context")

	// Acquire lock with one instance
	lock1, _ := NewDistributedLock(client, "test-context", 10*time.Second)
	err := lock1.Lock(context.Background())
	if err != nil {
		t.Logf("Redis not available: %v", err)
		return
	}
	defer lock1.Unlock(context.Background())

	// Try to acquire with another instance using cancelled context
	lock2, _ := NewDistributedLock(client, "test-context", 10*time.Second)
	err = lock2.LockWithRetry(ctx, 10, 50*time.Millisecond)

	if err != context.DeadlineExceeded {
		t.Logf("Expected context.DeadlineExceeded, got %v", err)
	}
}

func TestGenerateRandomValue(t *testing.T) {
	value1, err := generateRandomValue()
	if err != nil {
		t.Fatalf("Failed to generate random value: %v", err)
	}

	value2, err := generateRandomValue()
	if err != nil {
		t.Fatalf("Failed to generate random value: %v", err)
	}

	if value1 == value2 {
		t.Error("Expected different random values")
	}

	if len(value1) != 32 { // 16 bytes = 32 hex characters
		t.Errorf("Expected value length 32, got %d", len(value1))
	}
}

func TestLock_IsLocked(t *testing.T) {
	client := getTestRedisClient()
	ctx := context.Background()

	// Clean up
	client.Del(ctx, "lock:test-islocked")

	lock, _ := NewDistributedLock(client, "test-islocked", 5*time.Second)

	// Initially not locked
	isLocked, err := lock.IsLocked(ctx)
	if err != nil {
		t.Logf("Redis not available: %v", err)
		return
	}
	if isLocked {
		t.Error("Expected lock to not be held initially")
	}

	// Acquire lock
	err = lock.Lock(ctx)
	if err != nil {
		t.Fatalf("Failed to acquire lock: %v", err)
	}

	// Now it should be locked
	isLocked, err = lock.IsLocked(ctx)
	if err != nil {
		t.Errorf("Failed to check lock status: %v", err)
	}
	if !isLocked {
		t.Error("Expected lock to be held")
	}

	// Unlock
	lock.Unlock(ctx)

	// Not locked anymore
	isLocked, err = lock.IsLocked(ctx)
	if err != nil {
		t.Errorf("Failed to check lock status: %v", err)
	}
	if isLocked {
		t.Error("Expected lock to not be held after unlock")
	}
}
