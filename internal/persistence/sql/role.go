// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"context"
	"time"

	"github.com/coding-hui/iam/internal/persistence"
)

// RoleModel represents a role in the database.
type RoleModel struct {
	ID          string    `gorm:"primaryKey;column:id" json:"id"`
	NetworkID   string    `gorm:"column:nid;index"     json:"network_id"`
	Name        string    `gorm:"column:name"           json:"name"`
	Description string    `gorm:"column:description"   json:"description"`
	CreatedAt   time.Time `gorm:"column:created_at"    json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"    json:"updated_at"`
}

// TableName returns the table name for RoleModel.
func (RoleModel) TableName() string {
	return "iam_roles"
}

// RolePool implements persistence.RolePersister using GORM.
type RolePool struct {
	db *Persister
}

// NewRolePool creates a new role pool.
func NewRolePool(db *Persister) *RolePool {
	return &RolePool{db: db}
}

// GetRole retrieves a role by ID.
func (p *RolePool) GetRole(ctx context.Context, id string) (*persistence.Role, error) {
	var m RoleModel
	if err := p.db.Connection(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return p.modelToDomain(&m), nil
}

// ListRoles lists roles with pagination.
func (p *RolePool) ListRoles(ctx context.Context, networkID string, limit, offset int) ([]*persistence.Role, int, error) {
	var ms []RoleModel
	var total int64

	query := p.db.Connection(ctx)
	if networkID != "" {
		query = query.Where("nid = ?", networkID)
	}

	if err := query.Model(&RoleModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&ms).Error; err != nil {
		return nil, 0, err
	}

	roles := make([]*persistence.Role, len(ms))
	for i := range ms {
		roles[i] = p.modelToDomain(&ms[i])
	}
	return roles, int(total), nil
}

// CreateRole creates a new role.
func (p *RolePool) CreateRole(ctx context.Context, role *persistence.Role) error {
	m := p.domainToModel(role)
	return p.db.Connection(ctx).Create(m).Error
}

// UpdateRole updates a role.
func (p *RolePool) UpdateRole(ctx context.Context, role *persistence.Role) error {
	m := p.domainToModel(role)
	return p.db.Connection(ctx).Model(m).Where("id = ?", role.ID).Updates(m).Error
}

// DeleteRole deletes a role.
func (p *RolePool) DeleteRole(ctx context.Context, id string) error {
	return p.db.Connection(ctx).Where("id = ?", id).Delete(&RoleModel{}).Error
}

func (p *RolePool) modelToDomain(m *RoleModel) *persistence.Role {
	return &persistence.Role{
		ID:          m.ID,
		NetworkID:   m.NetworkID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (p *RolePool) domainToModel(r *persistence.Role) *RoleModel {
	return &RoleModel{
		ID:          r.ID,
		NetworkID:   r.NetworkID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// Ensure RolePool implements persistence.RolePersister.
var _ persistence.RolePersister = (*RolePool)(nil)
