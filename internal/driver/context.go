// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package driver

import (
	"context"

	"github.com/google/uuid"
)

// networkIDKey is the context key for network/tenant ID.
type networkIDKey struct{}

// NetworkID returns the current network/tenant ID from context.
// Returns uuid.Nil if not set.
func NetworkID(ctx context.Context) uuid.UUID {
	if v := ctx.Value(networkIDKey{}); v != nil {
		return v.(uuid.UUID)
	}
	return uuid.Nil
}

// WithNetworkID returns a new context with the given network ID.
func WithNetworkID(ctx context.Context, nid uuid.UUID) context.Context {
	return context.WithValue(ctx, networkIDKey{}, nid)
}
