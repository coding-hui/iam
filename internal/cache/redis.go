// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrCacheMiss is returned when a key is not found in cache.
var ErrCacheMiss = errors.New("cache miss")

// RedisCache is a Redis-backed cache implementation.
type RedisCache struct {
	client *redis.Client
	prefix string
}

// NewRedisCache creates a new Redis cache with the given prefix.
func NewRedisCache(client *redis.Client, prefix string) *RedisCache {
	return &RedisCache{
		client: client,
		prefix: prefix,
	}
}

// Get retrieves a value from Redis.
func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.client.Get(ctx, c.prefix+key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Set stores a value in Redis with the given TTL.
func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl int64) error {
	return c.client.Set(ctx, c.prefix+key, value, time.Duration(ttl)).Err()
}

// Delete removes a key from Redis.
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, c.prefix+key).Err()
}

// Keys returns all keys matching the pattern.
func (c *RedisCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	return c.client.Keys(ctx, c.prefix+pattern).Result()
}

// Ping checks if Redis is available.
func (c *RedisCache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// NewRedisClient creates a new Redis client.
func NewRedisClient(addr, password string, db, poolSize int) *redis.Client {
	if poolSize <= 0 {
		poolSize = 10
	}

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	})
}

// Ensure RedisCache implements Cache.
var _ Cache = (*RedisCache)(nil)
