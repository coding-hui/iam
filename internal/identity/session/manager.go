// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package session

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ManagerImpl implements session.Manager.
type ManagerImpl struct {
	pool     Pool
	privPool PrivilegedPool
}

// NewManagerImpl creates a new session manager.
func NewManagerImpl(pool Pool, privPool PrivilegedPool) *ManagerImpl {
	return &ManagerImpl{
		pool:     pool,
		privPool: privPool,
	}
}

// CreateSession creates a new session.
func (m *ManagerImpl) CreateSession(ctx context.Context, req *CreateSessionRequest) (*Session, error) {
	now := time.Now()
	sess := &Session{
		ID:              uuid.New(),
		IdentityID:      req.IdentityID,
		Active:          true,
		ExpiresAt:       now.Add(req.TTL),
		AuthenticatedAt: now,
		UserAgent:       req.UserAgent,
		ClientIP:        req.ClientIP,
		Extra:           req.Extra,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := m.privPool.CreateSession(ctx, sess); err != nil {
		return nil, err
	}

	return sess, nil
}

// GetSession retrieves a session by ID.
func (m *ManagerImpl) GetSession(ctx context.Context, id uuid.UUID) (*Session, error) {
	return m.pool.GetSession(ctx, id)
}

// ListSessions lists all active sessions for an identity.
func (m *ManagerImpl) ListSessions(ctx context.Context, identityID uuid.UUID) ([]*Session, error) {
	return m.pool.GetActiveSessions(ctx, identityID)
}

// RevokeSession revokes a single session.
func (m *ManagerImpl) RevokeSession(ctx context.Context, id uuid.UUID) error {
	networkID := uuid.Nil
	return m.privPool.RevokeSession(ctx, networkID, id)
}

// RevokeAllSessions revokes all sessions for an identity.
func (m *ManagerImpl) RevokeAllSessions(ctx context.Context, identityID uuid.UUID) error {
	networkID := uuid.Nil
	return m.privPool.RevokeAllSessions(ctx, networkID, identityID)
}

// ExtendSession extends a session's expiration.
func (m *ManagerImpl) ExtendSession(ctx context.Context, id uuid.UUID, extend time.Duration) (*Session, error) {
	sess, err := m.pool.GetSession(ctx, id)
	if err != nil {
		return nil, err
	}

	sess.ExpiresAt = sess.ExpiresAt.Add(extend)
	sess.UpdatedAt = time.Now()

	if err := m.privPool.UpdateSession(ctx, sess); err != nil {
		return nil, err
	}

	return sess, nil
}

// Ensure ManagerImpl implements Manager.
var _ Manager = (*ManagerImpl)(nil)
