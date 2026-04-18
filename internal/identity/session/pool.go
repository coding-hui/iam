// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package session

import (
	"context"

	"github.com/google/uuid"

	"github.com/coding-hui/iam/internal/persistence"
)

// sessionPool implements Pool using persistence.SessionPersister.
type sessionPool struct {
	persister sessionPersister
}

// sessionPersister is the persistence interface for session operations.
type sessionPersister interface {
	GetSession(ctx context.Context, id string) (*persistence.Session, error)
	ListSessionsByIdentityID(ctx context.Context, identityID string) ([]*persistence.Session, error)
	CreateSession(ctx context.Context, session *persistence.Session) error
	UpdateSession(ctx context.Context, session *persistence.Session) error
	DeleteSession(ctx context.Context, id string) error
}

// NewPool creates a new session pool.
func NewPool(p sessionPersister) Pool {
	return &sessionPool{persister: p}
}

// GetSession retrieves a session by ID.
func (p *sessionPool) GetSession(ctx context.Context, id uuid.UUID) (*Session, error) {
	m, err := p.persister.GetSession(ctx, id.String())
	if err != nil {
		return nil, err
	}
	return p.modelToDomain(m), nil
}

// GetSessionByNetworkID retrieves a session by network ID and session ID.
func (p *sessionPool) GetSessionByNetworkID(ctx context.Context, networkID, id uuid.UUID) (*Session, error) {
	return p.GetSession(ctx, id)
}

// GetActiveSessions retrieves all active sessions for an identity.
func (p *sessionPool) GetActiveSessions(ctx context.Context, identityID uuid.UUID) ([]*Session, error) {
	ms, err := p.persister.ListSessionsByIdentityID(ctx, identityID.String())
	if err != nil {
		return nil, err
	}
	sessions := make([]*Session, len(ms))
	for i := range ms {
		sessions[i] = p.modelToDomain(ms[i])
	}
	return sessions, nil
}

// ListSessions lists sessions with pagination.
func (p *sessionPool) ListSessions(ctx context.Context, networkID uuid.UUID, identityID uuid.UUID, limit, offset int) ([]*Session, int, error) {
	ms, err := p.persister.ListSessionsByIdentityID(ctx, identityID.String())
	if err != nil {
		return nil, 0, err
	}
	sessions := make([]*Session, len(ms))
	for i := range ms {
		sessions[i] = p.modelToDomain(ms[i])
	}
	return sessions, len(sessions), nil
}

func (p *sessionPool) modelToDomain(m *persistence.Session) *Session {
	if m == nil {
		return nil
	}
	return &Session{
		ID:              parseUUID(m.ID),
		IdentityID:      parseUUID(m.IdentityID),
		Active:          m.Active,
		ExpiresAt:       m.ExpiresAt,
		AuthenticatedAt: m.AuthenticatedAt,
		UserAgent:       m.UserAgent,
		ClientIP:        m.ClientIP,
		Extra:           m.Extra,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
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

// Ensure sessionPool implements Pool.
var _ Pool = (*sessionPool)(nil)
