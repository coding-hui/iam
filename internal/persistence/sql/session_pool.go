// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"context"
	"time"

	"github.com/coding-hui/iam/internal/persistence"
)

// SessionModel represents a session in the database.
type SessionModel struct {
	ID              string    `gorm:"primaryKey;column:id"     json:"id"`
	IdentityID      string    `gorm:"column:identity_id;index" json:"identity_id"`
	Active          bool      `gorm:"column:active"            json:"active"`
	ExpiresAt       time.Time `gorm:"column:expires_at"        json:"expires_at"`
	AuthenticatedAt time.Time `gorm:"column:authenticated_at"  json:"authenticated_at"`
	UserAgent       string    `gorm:"column:user_agent"        json:"user_agent"`
	ClientIP        string    `gorm:"column:client_ip"         json:"client_ip"`
	Extra           []byte    `gorm:"column:extra"             json:"extra"`
	CreatedAt       time.Time `gorm:"column:created_at"        json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"        json:"updated_at"`
}

// TableName returns the table name for SessionModel.
func (SessionModel) TableName() string {
	return "iam_sessions"
}

// SessionPool implements persistence.SessionPersister using GORM.
type SessionPool struct {
	db *Persister
}

// NewSessionPool creates a new session pool.
func NewSessionPool(db *Persister) *SessionPool {
	return &SessionPool{db: db}
}

// GetSession retrieves a session by ID.
func (p *SessionPool) GetSession(ctx context.Context, id string) (*persistence.Session, error) {
	var m SessionModel
	if err := p.db.Connection(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return p.modelToDomain(&m), nil
}

// ListSessionsByIdentityID lists sessions for an identity.
func (p *SessionPool) ListSessionsByIdentityID(ctx context.Context, identityID string) ([]*persistence.Session, error) {
	var ms []SessionModel
	if err := p.db.Connection(ctx).
		Where("identity_id = ? AND active = ?", identityID, true).
		Order("created_at DESC").
		Find(&ms).Error; err != nil {
		return nil, err
	}

	sessions := make([]*persistence.Session, len(ms))
	for i := range ms {
		sessions[i] = p.modelToDomain(&ms[i])
	}
	return sessions, nil
}

// CreateSession creates a new session.
func (p *SessionPool) CreateSession(ctx context.Context, session *persistence.Session) error {
	m := p.domainToModel(session)
	return p.db.Connection(ctx).Create(m).Error
}

// UpdateSession updates a session.
func (p *SessionPool) UpdateSession(ctx context.Context, session *persistence.Session) error {
	m := p.domainToModel(session)
	return p.db.Connection(ctx).Model(m).Where("id = ?", session.ID).Updates(m).Error
}

// DeleteSession deletes a session.
func (p *SessionPool) DeleteSession(ctx context.Context, id string) error {
	return p.db.Connection(ctx).Where("id = ?", id).Delete(&SessionModel{}).Error
}

func (p *SessionPool) modelToDomain(m *SessionModel) *persistence.Session {
	return &persistence.Session{
		ID:              m.ID,
		IdentityID:      m.IdentityID,
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

func (p *SessionPool) domainToModel(s *persistence.Session) *SessionModel {
	return &SessionModel{
		ID:              s.ID,
		IdentityID:      s.IdentityID,
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

// Ensure SessionPool implements persistence.SessionPersister.
var _ persistence.SessionPersister = (*SessionPool)(nil)
