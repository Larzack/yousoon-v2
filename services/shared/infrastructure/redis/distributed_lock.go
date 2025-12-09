package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

// DistributedLock provides a distributed lock implementation using Redis.
type DistributedLock struct {
	client   *Client
	key      string
	token    string
	ttl      time.Duration
	acquired bool
}

// NewDistributedLock creates a new distributed lock.
func NewDistributedLock(client *Client, key string, ttl time.Duration) *DistributedLock {
	return &DistributedLock{
		client: client,
		key:    fmt.Sprintf("lock:%s", key),
		token:  uuid.New().String(),
		ttl:    ttl,
	}
}

// Acquire attempts to acquire the lock.
// Returns true if the lock was acquired, false otherwise.
func (l *DistributedLock) Acquire(ctx context.Context) (bool, error) {
	acquired, err := l.client.SetNX(ctx, l.key, l.token, l.ttl)
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	l.acquired = acquired
	return acquired, nil
}

// AcquireWithRetry attempts to acquire the lock with retries.
func (l *DistributedLock) AcquireWithRetry(ctx context.Context, maxRetries int, retryInterval time.Duration) (bool, error) {
	for i := 0; i < maxRetries; i++ {
		acquired, err := l.Acquire(ctx)
		if err != nil {
			return false, err
		}

		if acquired {
			return true, nil
		}

		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(retryInterval):
			continue
		}
	}

	return false, nil
}

// Release releases the lock.
// Uses Lua script to ensure only the owner can release the lock.
func (l *DistributedLock) Release(ctx context.Context) error {
	if !l.acquired {
		return nil
	}

	// Lua script to release lock only if we own it
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	result, err := l.client.Client().Eval(ctx, script, []string{l.key}, l.token).Result()
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	if result.(int64) == 0 {
		return ErrLockNotOwned
	}

	l.acquired = false
	return nil
}

// Extend extends the lock TTL.
// Uses Lua script to ensure only the owner can extend the lock.
func (l *DistributedLock) Extend(ctx context.Context, ttl time.Duration) error {
	if !l.acquired {
		return ErrLockNotAcquired
	}

	// Lua script to extend lock only if we own it
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("pexpire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`

	result, err := l.client.Client().Eval(ctx, script, []string{l.key}, l.token, int64(ttl/time.Millisecond)).Result()
	if err != nil {
		return fmt.Errorf("failed to extend lock: %w", err)
	}

	if result.(int64) == 0 {
		l.acquired = false
		return ErrLockNotOwned
	}

	l.ttl = ttl
	return nil
}

// IsAcquired returns whether the lock is currently held.
func (l *DistributedLock) IsAcquired() bool {
	return l.acquired
}

// Key returns the lock key.
func (l *DistributedLock) Key() string {
	return l.key
}

// WithLock executes a function with a distributed lock.
// The lock is automatically released after the function returns.
func WithLock(ctx context.Context, client *Client, key string, ttl time.Duration, fn func() error) error {
	lock := NewDistributedLock(client, key, ttl)

	acquired, err := lock.Acquire(ctx)
	if err != nil {
		return err
	}

	if !acquired {
		return ErrLockNotAcquired
	}

	defer lock.Release(ctx)

	return fn()
}

// WithLockRetry executes a function with a distributed lock, retrying if the lock is not available.
func WithLockRetry(ctx context.Context, client *Client, key string, ttl time.Duration, maxRetries int, retryInterval time.Duration, fn func() error) error {
	lock := NewDistributedLock(client, key, ttl)

	acquired, err := lock.AcquireWithRetry(ctx, maxRetries, retryInterval)
	if err != nil {
		return err
	}

	if !acquired {
		return ErrLockNotAcquired
	}

	defer lock.Release(ctx)

	return fn()
}

// Semaphore provides a distributed semaphore using Redis.
type Semaphore struct {
	client   *Client
	key      string
	limit    int64
	ttl      time.Duration
	acquired bool
	token    string
}

// NewSemaphore creates a new distributed semaphore.
func NewSemaphore(client *Client, key string, limit int64, ttl time.Duration) *Semaphore {
	return &Semaphore{
		client: client,
		key:    fmt.Sprintf("semaphore:%s", key),
		limit:  limit,
		ttl:    ttl,
		token:  uuid.New().String(),
	}
}

// Acquire attempts to acquire a slot in the semaphore.
func (s *Semaphore) Acquire(ctx context.Context) (bool, error) {
	now := time.Now().UnixNano()
	expires := now + int64(s.ttl)

	// Remove expired entries
	_, err := s.client.Client().ZRemRangeByScore(ctx, s.key, "0", fmt.Sprintf("%d", now)).Result()
	if err != nil {
		return false, fmt.Errorf("failed to clean semaphore: %w", err)
	}

	// Check current count
	count, err := s.client.Client().ZCard(ctx, s.key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to get semaphore count: %w", err)
	}

	if count >= s.limit {
		return false, nil
	}

	// Try to add our entry
	added, err := s.client.Client().ZAdd(ctx, s.key, goredis.Z{
		Score:  float64(expires),
		Member: s.token,
	}).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire semaphore: %w", err)
	}

	if added == 0 {
		return false, nil
	}

	// Verify we're within the limit
	rank, err := s.client.Client().ZRank(ctx, s.key, s.token).Result()
	if err != nil {
		s.client.Client().ZRem(ctx, s.key, s.token)
		return false, fmt.Errorf("failed to verify semaphore position: %w", err)
	}

	if rank >= s.limit {
		s.client.Client().ZRem(ctx, s.key, s.token)
		return false, nil
	}

	s.acquired = true
	return true, nil
}

// Release releases the semaphore slot.
func (s *Semaphore) Release(ctx context.Context) error {
	if !s.acquired {
		return nil
	}

	_, err := s.client.Client().ZRem(ctx, s.key, s.token).Result()
	if err != nil {
		return fmt.Errorf("failed to release semaphore: %w", err)
	}

	s.acquired = false
	return nil
}

// Errors
var (
	ErrLockNotAcquired = fmt.Errorf("lock not acquired")
	ErrLockNotOwned    = fmt.Errorf("lock not owned by this instance")
)
