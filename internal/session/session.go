// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package session

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Session represents a user session.
type Session struct {
	ID              uuid.UUID       `json:"id"`
	IdentityID      uuid.UUID       `json:"identity_id"`
	Active          bool            `json:"active"`
	ExpiresAt       time.Time       `json:"expires_at"`
	AuthenticatedAt time.Time       `json:"authenticated_at"`
	UserAgent       string          `json:"user_agent"`
	ClientIP        string          `json:"client_ip"`
	Extra           json.RawMessage `json:"extra,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// Pool defines the interface for reading session data.
type Pool interface {
	GetSession(ctx context.Context, id uuid.UUID) (*Session, error)
	GetSessionByNetworkID(ctx context.Context, networkID, id uuid.UUID) (*Session, error)
	GetActiveSessions(ctx context.Context, identityID uuid.UUID) ([]*Session, error)
	ListSessions(ctx context.Context, networkID uuid.UUID, identityID uuid.UUID, limit, offset int) ([]*Session, int, error)
}

// PrivilegedPool defines the interface for writing session data.
type PrivilegedPool interface {
	Pool

	CreateSession(ctx context.Context, s *Session) error
	UpdateSession(ctx context.Context, s *Session) error
	RevokeSession(ctx context.Context, networkID uuid.UUID, id uuid.UUID) error
	RevokeAllSessions(ctx context.Context, networkID, identityID uuid.UUID) error
}

// Manager defines the interface for session business logic.
type Manager interface {
	CreateSession(ctx context.Context, req *CreateSessionRequest) (*Session, error)
	GetSession(ctx context.Context, id uuid.UUID) (*Session, error)
	ListSessions(ctx context.Context, identityID uuid.UUID) ([]*Session, error)
	RevokeSession(ctx context.Context, id uuid.UUID) error
	RevokeAllSessions(ctx context.Context, identityID uuid.UUID) error
	ExtendSession(ctx context.Context, id uuid.UUID, extend time.Duration) (*Session, error)
}

// CreateSessionRequest holds data for creating a new session.
type CreateSessionRequest struct {
	IdentityID uuid.UUID       `json:"identity_id"`
	UserAgent  string          `json:"user_agent"`
	ClientIP   string          `json:"client_ip"`
	TTL        time.Duration   `json:"ttl"` // default 24h
	Extra      json.RawMessage `json:"extra,omitempty"`
}
