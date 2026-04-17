// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package token

import (
	"context"

	"github.com/google/uuid"

	"github.com/coding-hui/iam/internal/persistence"
)

// tokenPool implements Pool using persistence.TokenPersister.
type tokenPool struct {
	persister tokenPersister
}

// tokenPersister is the persistence interface for token operations.
type tokenPersister interface {
	GetToken(ctx context.Context, id string) (*persistence.Token, error)
	GetTokenByValue(ctx context.Context, value string) (*persistence.Token, error)
	CreateToken(ctx context.Context, token *persistence.Token) error
	DeleteToken(ctx context.Context, id string) error
	DeleteExpiredTokens(ctx context.Context) error
}

// NewPool creates a new token pool.
func NewPool(p tokenPersister) Pool {
	return &tokenPool{persister: p}
}

// GetToken retrieves a token by ID.
func (p *tokenPool) GetToken(ctx context.Context, id uuid.UUID) (*Token, error) {
	m, err := p.persister.GetToken(ctx, id.String())
	if err != nil {
		return nil, err
	}
	return p.modelToDomain(m), nil
}

// GetTokenByValue retrieves a token by its value.
func (p *tokenPool) GetTokenByValue(ctx context.Context, value string) (*Token, error) {
	m, err := p.persister.GetTokenByValue(ctx, value)
	if err != nil {
		return nil, err
	}
	return p.modelToDomain(m), nil
}

// ListTokens lists all tokens for an identity.
func (p *tokenPool) ListTokens(ctx context.Context, identityID uuid.UUID) ([]*Token, error) {
	// Not in persister interface - return empty for now
	return []*Token{}, nil
}

func (p *tokenPool) modelToDomain(m *persistence.Token) *Token {
	if m == nil {
		return nil
	}
	return &Token{
		ID:         parseUUID(m.ID),
		IdentityID: parseUUID(m.IdentityID),
		Type:       TokenType(m.Type),
		Value:      m.Value,
		ExpiresAt:  m.ExpiresAt,
		CreatedAt:  m.CreatedAt,
	}
}

func parseUUID(s string) uuid.UUID {
	if s == "" {
		return uuid.Nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}

// Ensure tokenPool implements Pool.
var _ Pool = (*tokenPool)(nil)

// privilegedPool implements PrivilegedPool.
type privilegedPool struct {
	*tokenPool
}

// NewPrivilegedPool creates a new token privileged pool.
func NewPrivilegedPool(p tokenPersister) PrivilegedPool {
	return &privilegedPool{
		tokenPool: &tokenPool{persister: p},
	}
}

// CreateToken creates a new token.
func (p *privilegedPool) CreateToken(ctx context.Context, t *Token) error {
	return p.persister.CreateToken(ctx, p.domainToModel(t))
}

// DeleteToken deletes a token.
func (p *privilegedPool) DeleteToken(ctx context.Context, id uuid.UUID) error {
	return p.persister.DeleteToken(ctx, id.String())
}

// DeleteExpiredTokens deletes all expired tokens.
func (p *privilegedPool) DeleteExpiredTokens(ctx context.Context) error {
	return p.persister.DeleteExpiredTokens(ctx)
}

func (p *privilegedPool) domainToModel(t *Token) *persistence.Token {
	return &persistence.Token{
		ID:         t.ID.String(),
		IdentityID: t.IdentityID.String(),
		Type:       string(t.Type),
		Value:      t.Value,
		ExpiresAt:  t.ExpiresAt,
		CreatedAt:  t.CreatedAt,
	}
}

// Ensure privilegedPool implements PrivilegedPool.
var _ PrivilegedPool = (*privilegedPool)(nil)
