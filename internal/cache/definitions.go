// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cache

import (
	"context"
)

// Cache defines the interface for caching operations.
// It supports both in-memory and distributed (Redis) implementations.
type Cache interface {
	// Get retrieves a value by key.
	// Returns nil, nil if the key does not exist.
	Get(ctx context.Context, key string) ([]byte, error)

	// Set stores a value with the given TTL in nanoseconds.
	// A TTL of 0 means no expiration.
	Set(ctx context.Context, key string, value []byte, ttl int64) error

	// Delete removes a key.
	Delete(ctx context.Context, key string) error

	// Keys returns all keys matching the pattern.
	// Pattern uses glob-style matching (e.g., "user:*").
	Keys(ctx context.Context, pattern string) ([]string, error)

	// Ping checks if the cache is available.
	Ping(ctx context.Context) error
}
