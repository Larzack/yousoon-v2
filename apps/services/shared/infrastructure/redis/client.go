// Package redis provides Redis infrastructure components.
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Config holds Redis connection configuration.
type Config struct {
	// Address is the Redis server address.
	Address string
	// Password is the Redis password.
	Password string
	// DB is the Redis database number.
	DB int
	// PoolSize is the maximum number of connections.
	PoolSize int
	// MinIdleConns is the minimum number of idle connections.
	MinIdleConns int
	// MaxRetries is the maximum number of retries.
	MaxRetries int
	// DialTimeout is the timeout for establishing connections.
	DialTimeout time.Duration
	// ReadTimeout is the timeout for read operations.
	ReadTimeout time.Duration
	// WriteTimeout is the timeout for write operations.
	WriteTimeout time.Duration
}

// DefaultConfig returns a default Redis configuration.
func DefaultConfig() Config {
	return Config{
		Address:      "localhost:6379",
		Password:     "",
		DB:           0,
		PoolSize:     100,
		MinIdleConns: 10,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

// Client wraps the Redis client with additional functionality.
type Client struct {
	client *redis.Client
	config Config
}

// NewClient creates a new Redis client with the given configuration.
func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Address,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	// Verify connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{
		client: client,
		config: cfg,
	}, nil
}

// Client returns the underlying Redis client.
func (c *Client) Client() *redis.Client {
	return c.client
}

// Close closes the Redis connection.
func (c *Client) Close() error {
	return c.client.Close()
}

// Ping verifies the connection to Redis.
func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// HealthCheck performs a health check on the Redis connection.
func (c *Client) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := c.Ping(ctx); err != nil {
		return fmt.Errorf("Redis health check failed: %w", err)
	}
	return nil
}

// Set stores a value with the given key and expiration.
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value by key.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	result, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrKeyNotFound
	}
	if err != nil {
		return "", fmt.Errorf("failed to get key: %w", err)
	}
	return result, nil
}

// GetBytes retrieves a value as bytes by key.
func (c *Client) GetBytes(ctx context.Context, key string) ([]byte, error) {
	result, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, ErrKeyNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}
	return result, nil
}

// Delete removes a key.
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists.
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check key existence: %w", err)
	}
	return result > 0, nil
}

// Expire sets the expiration for a key.
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// TTL returns the time-to-live for a key.
func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

// Incr increments a key's value by 1.
func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// IncrBy increments a key's value by the given amount.
func (c *Client) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.IncrBy(ctx, key, value).Result()
}

// Decr decrements a key's value by 1.
func (c *Client) Decr(ctx context.Context, key string) (int64, error) {
	return c.client.Decr(ctx, key).Result()
}

// SetNX sets a value only if the key does not exist.
func (c *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.client.SetNX(ctx, key, value, expiration).Result()
}

// HSet sets a hash field value.
func (c *Client) HSet(ctx context.Context, key string, values ...interface{}) error {
	return c.client.HSet(ctx, key, values...).Err()
}

// HGet gets a hash field value.
func (c *Client) HGet(ctx context.Context, key, field string) (string, error) {
	result, err := c.client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", ErrKeyNotFound
	}
	if err != nil {
		return "", fmt.Errorf("failed to get hash field: %w", err)
	}
	return result, nil
}

// HGetAll gets all hash field values.
func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, key).Result()
}

// HDel deletes hash fields.
func (c *Client) HDel(ctx context.Context, key string, fields ...string) error {
	return c.client.HDel(ctx, key, fields...).Err()
}

// LPush prepends values to a list.
func (c *Client) LPush(ctx context.Context, key string, values ...interface{}) error {
	return c.client.LPush(ctx, key, values...).Err()
}

// RPush appends values to a list.
func (c *Client) RPush(ctx context.Context, key string, values ...interface{}) error {
	return c.client.RPush(ctx, key, values...).Err()
}

// LRange gets a range of elements from a list.
func (c *Client) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.LRange(ctx, key, start, stop).Result()
}

// LLen gets the length of a list.
func (c *Client) LLen(ctx context.Context, key string) (int64, error) {
	return c.client.LLen(ctx, key).Result()
}

// SAdd adds members to a set.
func (c *Client) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return c.client.SAdd(ctx, key, members...).Err()
}

// SMembers gets all members of a set.
func (c *Client) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.client.SMembers(ctx, key).Result()
}

// SIsMember checks if a value is a member of a set.
func (c *Client) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return c.client.SIsMember(ctx, key, member).Result()
}

// SRem removes members from a set.
func (c *Client) SRem(ctx context.Context, key string, members ...interface{}) error {
	return c.client.SRem(ctx, key, members...).Err()
}

// ZAdd adds members to a sorted set.
func (c *Client) ZAdd(ctx context.Context, key string, members ...redis.Z) error {
	return c.client.ZAdd(ctx, key, members...).Err()
}

// ZRange gets a range from a sorted set.
func (c *Client) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores gets a range from a sorted set with scores.
func (c *Client) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	return c.client.ZRangeWithScores(ctx, key, start, stop).Result()
}

// Keys returns all keys matching a pattern.
func (c *Client) Keys(ctx context.Context, pattern string) ([]string, error) {
	return c.client.Keys(ctx, pattern).Result()
}

// Scan iterates over keys matching a pattern.
func (c *Client) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return c.client.Scan(ctx, cursor, match, count).Result()
}

// Pipeline returns a pipeline for batching commands.
func (c *Client) Pipeline() redis.Pipeliner {
	return c.client.Pipeline()
}

// TxPipeline returns a transactional pipeline.
func (c *Client) TxPipeline() redis.Pipeliner {
	return c.client.TxPipeline()
}

// Publish publishes a message to a channel.
func (c *Client) Publish(ctx context.Context, channel string, message interface{}) error {
	return c.client.Publish(ctx, channel, message).Err()
}

// Subscribe subscribes to channels.
func (c *Client) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.client.Subscribe(ctx, channels...)
}

// ErrKeyNotFound is returned when a key is not found in Redis.
var ErrKeyNotFound = fmt.Errorf("key not found")
