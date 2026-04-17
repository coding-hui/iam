// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package persistence

import (
	"context"
	"time"
)

// Identity represents an identity in the system.
// Domain model with no persistence-specific tags (Ory style).
type Identity struct {
	ID        string
	NetworkID string
	SchemaID  string
	Traits    []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IdentityPersister defines the interface for identity persistence operations.
type IdentityPersister interface {
	GetIdentity(ctx context.Context, id string) (*Identity, error)
	ListIdentities(ctx context.Context, networkID string, limit, offset int) ([]*Identity, int, error)
	CreateIdentity(ctx context.Context, identity *Identity) error
	UpdateIdentity(ctx context.Context, identity *Identity) error
	DeleteIdentity(ctx context.Context, id string) error
}
