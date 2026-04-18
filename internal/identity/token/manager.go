// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package token

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/google/uuid"
)

// ManagerImpl implements token.Manager.
type ManagerImpl struct {
	pool     Pool
	privPool PrivilegedPool
}

// NewManagerImpl creates a new token manager.
func NewManagerImpl(pool Pool, privPool PrivilegedPool) *ManagerImpl {
	return &ManagerImpl{
		pool:     pool,
		privPool: privPool,
	}
}

// CreateToken creates a new token.
func (m *ManagerImpl) CreateToken(ctx context.Context, req *CreateTokenRequest) (*Token, error) {
	if req.TTL == 0 {
		req.TTL = 1 * time.Hour
	}

	value, err := generateToken()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	t := &Token{
		ID:         uuid.New(),
		IdentityID: req.IdentityID,
		Type:       req.Type,
		Value:      value,
		ExpiresAt:  now.Add(req.TTL),
		CreatedAt:  now,
	}

	if err := m.privPool.CreateToken(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

// GetToken retrieves a token by ID.
func (m *ManagerImpl) GetToken(ctx context.Context, id uuid.UUID) (*Token, error) {
	return m.pool.GetToken(ctx, id)
}

// IntrospectToken introspects a token by value.
func (m *ManagerImpl) IntrospectToken(ctx context.Context, value string) (*Token, error) {
	t, err := m.pool.GetTokenByValue(ctx, value)
	if err != nil {
		return nil, err
	}
	if time.Now().After(t.ExpiresAt) {
		return nil, ErrTokenExpired
	}
	return t, nil
}

// RevokeToken revokes a token by ID.
func (m *ManagerImpl) RevokeToken(ctx context.Context, id uuid.UUID) error {
	return m.privPool.DeleteToken(ctx, id)
}

// RevokeAllTokens revokes all tokens for an identity.
func (m *ManagerImpl) RevokeAllTokens(ctx context.Context, identityID uuid.UUID) error {
	tokens, err := m.pool.ListTokens(ctx, identityID)
	if err != nil {
		return err
	}
	for _, t := range tokens {
		if err := m.privPool.DeleteToken(ctx, t.ID); err != nil {
			return err
		}
	}
	return nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Ensure ManagerImpl implements Manager.
var _ Manager = (*ManagerImpl)(nil)
