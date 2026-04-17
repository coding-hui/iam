// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cache

import (
	"context"
	"sync"
	"time"
)

// memItem represents a cached item with expiration.
type memItem struct {
	value    []byte
	expireAt time.Time
}

// MemoryCache is an in-memory cache implementation.
// It is suitable for single-instance deployments and testing.
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]*memItem
}

// NewMemoryCache creates a new in-memory cache.
func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{items: make(map[string]*memItem)}
	go c.cleanup()
	return c
}

// Get retrieves a value from memory cache.
// Returns nil, nil if the key does not exist or has expired.
func (c *MemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return nil, nil
	}

	// Check if expired
	if !item.expireAt.IsZero() && item.expireAt.Before(time.Now()) {
		return nil, nil
	}

	// Return a copy to prevent external modification
	result := make([]byte, len(item.value))
	copy(result, item.value)
	return result, nil
}

// Set stores a value in memory cache with the given TTL in nanoseconds.
func (c *MemoryCache) Set(ctx context.Context, key string, value []byte, ttl int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	expireAt := time.Time{}
	if ttl > 0 {
		expireAt = time.Now().Add(time.Duration(ttl))
	}

	c.items[key] = &memItem{
		value:    value,
		expireAt: expireAt,
	}
	return nil
}

// Delete removes a key from memory cache.
func (c *MemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	return nil
}

// Keys returns all keys in memory cache.
// Note: This is not efficient for large caches.
func (c *MemoryCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Simple glob matching
	var result []string
	for k := range c.items {
		if matchPattern(k, pattern) {
			result = append(result, k)
		}
	}
	return result, nil
}

// Ping always returns nil (in-memory cache is always available).
func (c *MemoryCache) Ping(ctx context.Context) error {
	return nil
}

// cleanup periodically removes expired items.
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for k, v := range c.items {
			if !v.expireAt.IsZero() && v.expireAt.Before(now) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}

// matchPattern performs simple glob-style matching.
func matchPattern(key, pattern string) bool {
	// Simple implementation - supports * wildcard only
	if pattern == "*" {
		return true
	}
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		return len(key) > 0 && key[:len(key)-1] == pattern[:len(pattern)-1]
	}
	return key == pattern
}

// Ensure MemoryCache implements Cache.
var _ Cache = (*MemoryCache)(nil)
