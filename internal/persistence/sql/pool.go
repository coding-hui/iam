// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"context"
	"time"

	"github.com/coding-hui/iam/internal/persistence"
)

// IdentityModel represents an identity in the database.
type IdentityModel struct {
	ID        string    `gorm:"primaryKey;column:id" json:"id"`
	NetworkID string    `gorm:"column:nid;index"     json:"network_id"`
	SchemaID  string    `gorm:"column:schema_id"     json:"schema_id"`
	Traits    []byte    `gorm:"column:traits"        json:"traits"`
	CreatedAt time.Time `gorm:"column:created_at"    json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"    json:"updated_at"`
}

// TableName returns the table name for IdentityModel.
func (IdentityModel) TableName() string {
	return "iam_identities"
}

// IdentityPool implements persistence.IdentityPersister using GORM.
type IdentityPool struct {
	db *Persister
}

// NewIdentityPool creates a new identity pool.
func NewIdentityPool(db *Persister) *IdentityPool {
	return &IdentityPool{db: db}
}

// GetIdentity retrieves an identity by ID.
func (p *IdentityPool) GetIdentity(ctx context.Context, id string) (*persistence.Identity, error) {
	var m IdentityModel
	if err := p.db.Connection(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return p.modelToDomain(&m), nil
}

// ListIdentities lists identities with pagination.
func (p *IdentityPool) ListIdentities(ctx context.Context, networkID string, limit, offset int) ([]*persistence.Identity, int, error) {
	var ms []IdentityModel
	var total int64

	query := p.db.Connection(ctx)
	if networkID != "" {
		query = query.Where("nid = ?", networkID)
	}

	if err := query.Model(&IdentityModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&ms).Error; err != nil {
		return nil, 0, err
	}

	identities := make([]*persistence.Identity, len(ms))
	for i := range ms {
		identities[i] = p.modelToDomain(&ms[i])
	}
	return identities, int(total), nil
}

// CreateIdentity creates a new identity.
func (p *IdentityPool) CreateIdentity(ctx context.Context, identity *persistence.Identity) error {
	m := p.domainToModel(identity)
	return p.db.Connection(ctx).Create(m).Error
}

// UpdateIdentity updates an identity.
func (p *IdentityPool) UpdateIdentity(ctx context.Context, identity *persistence.Identity) error {
	m := p.domainToModel(identity)
	return p.db.Connection(ctx).Model(m).Where("id = ?", identity.ID).Updates(m).Error
}

// DeleteIdentity deletes an identity.
func (p *IdentityPool) DeleteIdentity(ctx context.Context, id string) error {
	return p.db.Connection(ctx).Where("id = ?", id).Delete(&IdentityModel{}).Error
}

func (p *IdentityPool) modelToDomain(m *IdentityModel) *persistence.Identity {
	return &persistence.Identity{
		ID:        m.ID,
		NetworkID: m.NetworkID,
		SchemaID:  m.SchemaID,
		Traits:    m.Traits,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (p *IdentityPool) domainToModel(i *persistence.Identity) *IdentityModel {
	return &IdentityModel{
		ID:        i.ID,
		NetworkID: i.NetworkID,
		SchemaID:  i.SchemaID,
		Traits:    i.Traits,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}

// Ensure IdentityPool implements persistence.IdentityPersister.
var _ persistence.IdentityPersister = (*IdentityPool)(nil)
