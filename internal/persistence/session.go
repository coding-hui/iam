// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source is governed by a MIT style
// license that can be found in the LICENSE file.

package persistence

import (
	"context"
	"time"
)

// Session represents a session in the system.
// Domain model with no persistence-specific tags (Ory style).
type Session struct {
	ID              string
	IdentityID      string
	Active          bool
	ExpiresAt       time.Time
	AuthenticatedAt time.Time
	UserAgent       string
	ClientIP        string
	Extra           []byte
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// SessionPersister defines the interface for session persistence operations.
type SessionPersister interface {
	GetSession(ctx context.Context, id string) (*Session, error)
	ListSessionsByIdentityID(ctx context.Context, identityID string) ([]*Session, error)
	CreateSession(ctx context.Context, session *Session) error
	UpdateSession(ctx context.Context, session *Session) error
	DeleteSession(ctx context.Context, id string) error
}
