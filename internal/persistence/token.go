// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package persistence

import (
	"context"
	"time"
)

// Token represents a token in the system.
// Domain model with no persistence-specific tags (Ory style).
type Token struct {
	ID         string
	IdentityID string
	Type       string
	Value      string
	ExpiresAt  time.Time
	CreatedAt  time.Time
}

// TokenPersister defines the interface for token persistence operations.
type TokenPersister interface {
	GetToken(ctx context.Context, id string) (*Token, error)
	GetTokenByValue(ctx context.Context, value string) (*Token, error)
	CreateToken(ctx context.Context, token *Token) error
	DeleteToken(ctx context.Context, id string) error
	DeleteExpiredTokens(ctx context.Context) error
}
