// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"context"
	"time"

	"github.com/coding-hui/iam/internal/persistence"
)

// TokenModel represents a token in the database.
type TokenModel struct {
	ID         string    `gorm:"primaryKey;column:id" json:"id"`
	IdentityID string    `gorm:"column:identity_id;index" json:"identity_id"`
	Type       string    `gorm:"column:type"           json:"type"`
	Value      string    `gorm:"column:value;uniqueIndex" json:"value"`
	ExpiresAt  time.Time `gorm:"column:expires_at"     json:"expires_at"`
	CreatedAt  time.Time `gorm:"column:created_at"     json:"created_at"`
}

// TableName returns the table name for TokenModel.
func (TokenModel) TableName() string {
	return "iam_tokens"
}

// TokenPool implements persistence.TokenPersister using GORM.
type TokenPool struct {
	db *Persister
}

// NewTokenPool creates a new token pool.
func NewTokenPool(db *Persister) *TokenPool {
	return &TokenPool{db: db}
}

// GetToken retrieves a token by ID.
func (p *TokenPool) GetToken(ctx context.Context, id string) (*persistence.Token, error) {
	var m TokenModel
	if err := p.db.Connection(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return p.modelToDomain(&m), nil
}

// GetTokenByValue retrieves a token by its value.
func (p *TokenPool) GetTokenByValue(ctx context.Context, value string) (*persistence.Token, error) {
	var m TokenModel
	if err := p.db.Connection(ctx).Where("value = ?", value).First(&m).Error; err != nil {
		return nil, err
	}
	return p.modelToDomain(&m), nil
}

// CreateToken creates a new token.
func (p *TokenPool) CreateToken(ctx context.Context, token *persistence.Token) error {
	m := p.domainToModel(token)
	return p.db.Connection(ctx).Create(m).Error
}

// DeleteToken deletes a token.
func (p *TokenPool) DeleteToken(ctx context.Context, id string) error {
	return p.db.Connection(ctx).Where("id = ?", id).Delete(&TokenModel{}).Error
}

// DeleteExpiredTokens deletes all expired tokens.
func (p *TokenPool) DeleteExpiredTokens(ctx context.Context) error {
	return p.db.Connection(ctx).Where("expires_at < ?", time.Now()).Delete(&TokenModel{}).Error
}

func (p *TokenPool) modelToDomain(m *TokenModel) *persistence.Token {
	return &persistence.Token{
		ID:         m.ID,
		IdentityID: m.IdentityID,
		Type:       m.Type,
		Value:      m.Value,
		ExpiresAt:  m.ExpiresAt,
		CreatedAt:  m.CreatedAt,
	}
}

func (p *TokenPool) domainToModel(t *persistence.Token) *TokenModel {
	return &TokenModel{
		ID:         t.ID,
		IdentityID: t.IdentityID,
		Type:       t.Type,
		Value:      t.Value,
		ExpiresAt:  t.ExpiresAt,
		CreatedAt:  t.CreatedAt,
	}
}

// Ensure TokenPool implements persistence.TokenPersister.
var _ persistence.TokenPersister = (*TokenPool)(nil)
