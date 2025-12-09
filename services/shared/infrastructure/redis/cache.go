package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Cache provides a type-safe caching layer on top of Redis.
type Cache struct {
	client *Client
	prefix string
}

// NewCache creates a new cache with the given prefix.
func NewCache(client *Client, prefix string) *Cache {
	return &Cache{
		client: client,
		prefix: prefix,
	}
}

// key builds a cache key with the prefix.
func (c *Cache) key(k string) string {
	if c.prefix == "" {
		return k
	}
	return fmt.Sprintf("%s:%s", c.prefix, k)
}

// Get retrieves a value from the cache and unmarshals it into the target.
func (c *Cache) Get(ctx context.Context, key string, target interface{}) error {
	data, err := c.client.GetBytes(ctx, c.key(key))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal cached value: %w", err)
	}

	return nil
}

// Set stores a value in the cache with the given expiration.
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.Set(ctx, c.key(key), data, expiration)
}

// Delete removes a value from the cache.
func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.client.Delete(ctx, c.key(key))
}

// DeletePattern removes all values matching a pattern.
func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := c.client.Keys(ctx, c.key(pattern))
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	if len(keys) > 0 {
		return c.client.Delete(ctx, keys...)
	}
	return nil
}

// Exists checks if a key exists in the cache.
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	return c.client.Exists(ctx, c.key(key))
}

// GetOrSet retrieves a value from the cache, or sets it using the provided function.
func (c *Cache) GetOrSet(ctx context.Context, key string, target interface{}, ttl time.Duration, fn func() (interface{}, error)) error {
	// Try to get from cache
	err := c.Get(ctx, key, target)
	if err == nil {
		return nil
	}

	// If not found, call the function
	if err == ErrKeyNotFound {
		value, fnErr := fn()
		if fnErr != nil {
			return fnErr
		}

		// Set in cache
		if setErr := c.Set(ctx, key, value, ttl); setErr != nil {
			return setErr
		}

		// Marshal and unmarshal to populate target
		data, marshalErr := json.Marshal(value)
		if marshalErr != nil {
			return fmt.Errorf("failed to marshal value: %w", marshalErr)
		}

		return json.Unmarshal(data, target)
	}

	return err
}

// TTL returns the remaining time-to-live for a key.
func (c *Cache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, c.key(key))
}

// Refresh extends the TTL of a key.
func (c *Cache) Refresh(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, c.key(key), expiration)
}

// Invalidate removes multiple keys from the cache.
func (c *Cache) Invalidate(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	prefixedKeys := make([]string, len(keys))
	for i, k := range keys {
		prefixedKeys[i] = c.key(k)
	}

	return c.client.Delete(ctx, prefixedKeys...)
}

// CacheKey helpers for building consistent cache keys.

// UserCacheKey builds a cache key for a user.
func UserCacheKey(userID string) string {
	return fmt.Sprintf("user:%s", userID)
}

// OfferCacheKey builds a cache key for an offer.
func OfferCacheKey(offerID string) string {
	return fmt.Sprintf("offer:%s", offerID)
}

// PartnerCacheKey builds a cache key for a partner.
func PartnerCacheKey(partnerID string) string {
	return fmt.Sprintf("partner:%s", partnerID)
}

// EstablishmentCacheKey builds a cache key for an establishment.
func EstablishmentCacheKey(establishmentID string) string {
	return fmt.Sprintf("establishment:%s", establishmentID)
}

// SessionCacheKey builds a cache key for a session.
func SessionCacheKey(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}

// RefreshTokenCacheKey builds a cache key for a refresh token.
func RefreshTokenCacheKey(tokenID string) string {
	return fmt.Sprintf("refresh_token:%s", tokenID)
}

// RateLimitCacheKey builds a cache key for rate limiting.
func RateLimitCacheKey(identifier, action string) string {
	return fmt.Sprintf("ratelimit:%s:%s", identifier, action)
}

// CacheTTL constants for common TTL values.
const (
	// CacheTTLShort is 5 minutes.
	CacheTTLShort = 5 * time.Minute
	// CacheTTLMedium is 15 minutes.
	CacheTTLMedium = 15 * time.Minute
	// CacheTTLLong is 1 hour.
	CacheTTLLong = 1 * time.Hour
	// CacheTTLDay is 24 hours.
	CacheTTLDay = 24 * time.Hour
	// CacheTTLWeek is 7 days.
	CacheTTLWeek = 7 * 24 * time.Hour

	// SessionTTL is 6 hours (access token lifetime).
	SessionTTL = 6 * time.Hour
	// RefreshTokenTTL is 30 days.
	RefreshTokenTTL = 30 * 24 * time.Hour
)
