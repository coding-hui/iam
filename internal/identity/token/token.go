// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package token

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// TokenType represents the type of token.
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
	TokenTypeAPIKey  TokenType = "api_key"
)

// Token represents a token in the system.
type Token struct {
	ID         uuid.UUID `json:"id"`
	IdentityID uuid.UUID `json:"identity_id"`
	Type       TokenType `json:"type"`
	Value      string    `json:"value"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}

// Pool defines the interface for reading token data.
type Pool interface {
	GetToken(ctx context.Context, id uuid.UUID) (*Token, error)
	GetTokenByValue(ctx context.Context, value string) (*Token, error)
	ListTokens(ctx context.Context, identityID uuid.UUID) ([]*Token, error)
}

// PrivilegedPool defines the interface for writing token data.
type PrivilegedPool interface {
	Pool

	CreateToken(ctx context.Context, t *Token) error
	DeleteToken(ctx context.Context, id uuid.UUID) error
	DeleteExpiredTokens(ctx context.Context) error
}

// Manager defines the interface for token business logic.
type Manager interface {
	CreateToken(ctx context.Context, req *CreateTokenRequest) (*Token, error)
	GetToken(ctx context.Context, id uuid.UUID) (*Token, error)
	IntrospectToken(ctx context.Context, value string) (*Token, error)
	RevokeToken(ctx context.Context, id uuid.UUID) error
	RevokeAllTokens(ctx context.Context, identityID uuid.UUID) error
}

// CreateTokenRequest holds data for creating a new token.
type CreateTokenRequest struct {
	IdentityID uuid.UUID     `json:"identity_id"`
	Type       TokenType     `json:"type"`
	TTL        time.Duration `json:"ttl"`
}
