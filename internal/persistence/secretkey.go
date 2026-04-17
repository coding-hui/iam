// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package persistence

import (
	"context"
	"time"
)

// SecretKey represents a secret key in the system.
// Domain model with no persistence-specific tags (Ory style).
type SecretKey struct {
	ID         string
	IdentityID string
	KeyID      string
	SecretHash string
	Name       string
	ExpiresAt  *time.Time
	CreatedAt  time.Time
}

// SecretKeyPersister defines the interface for secret key persistence operations.
type SecretKeyPersister interface {
	GetSecretKey(ctx context.Context, id string) (*SecretKey, error)
	GetSecretKeyByKeyID(ctx context.Context, keyID string) (*SecretKey, error)
	ListSecretKeysByIdentityID(ctx context.Context, identityID string) ([]*SecretKey, error)
	CreateSecretKey(ctx context.Context, secretKey *SecretKey) error
	DeleteSecretKey(ctx context.Context, id string) error
}
