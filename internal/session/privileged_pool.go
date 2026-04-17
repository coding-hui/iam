// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package session

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/coding-hui/iam/internal/persistence"
)

// privilegedPool implements PrivilegedPool.
type privilegedPool struct {
	*sessionPool
}

// NewPrivilegedPool creates a new session privileged pool.
func NewPrivilegedPool(p sessionPersister) PrivilegedPool {
	return &privilegedPool{
		sessionPool: &sessionPool{persister: p},
	}
}

// CreateSession creates a new session.
func (p *privilegedPool) CreateSession(ctx context.Context, s *Session) error {
	return p.persister.CreateSession(ctx, p.domainToModel(s))
}

// UpdateSession updates a session.
func (p *privilegedPool) UpdateSession(ctx context.Context, s *Session) error {
	return p.persister.UpdateSession(ctx, p.domainToModel(s))
}

// RevokeSession revokes a session by ID.
func (p *privilegedPool) RevokeSession(ctx context.Context, networkID uuid.UUID, id uuid.UUID) error {
	m, err := p.persister.GetSession(ctx, id.String())
	if err != nil {
		return err
	}
	m.Active = false
	m.UpdatedAt = time.Now()
	return p.persister.UpdateSession(ctx, m)
}

// RevokeAllSessions revokes all sessions for an identity.
func (p *privilegedPool) RevokeAllSessions(ctx context.Context, networkID, identityID uuid.UUID) error {
	sessions, err := p.persister.ListSessionsByIdentityID(ctx, identityID.String())
	if err != nil {
		return err
	}
	now := time.Now()
	for _, s := range sessions {
		s.Active = false
		s.UpdatedAt = now
		if err := p.persister.UpdateSession(ctx, s); err != nil {
			return err
		}
	}
	return nil
}

func (p *privilegedPool) domainToModel(s *Session) *persistence.Session {
	return &persistence.Session{
		ID:              s.ID.String(),
		IdentityID:      s.IdentityID.String(),
		Active:          s.Active,
		ExpiresAt:       s.ExpiresAt,
		AuthenticatedAt: s.AuthenticatedAt,
		UserAgent:       s.UserAgent,
		ClientIP:        s.ClientIP,
		Extra:           s.Extra,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}

// Ensure privilegedPool implements PrivilegedPool.
var _ PrivilegedPool = (*privilegedPool)(nil)
