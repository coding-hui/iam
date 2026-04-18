// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"context"
	"time"

	"github.com/coding-hui/iam/internal/persistence"
)

// PolicyModel represents a policy in the database.
type PolicyModel struct {
	ID         string    `gorm:"primaryKey;column:id" json:"id"`
	NetworkID  string    `gorm:"column:nid;index"     json:"network_id"`
	Name       string    `gorm:"column:name"          json:"name"`
	Type       string    `gorm:"column:type"          json:"type"`
	Subjects   string    `gorm:"column:subjects"      json:"subjects"`
	Effect     string    `gorm:"column:effect"        json:"effect"`
	Actions    string    `gorm:"column:actions"       json:"actions"`
	Resources  string    `gorm:"column:resources"     json:"resources"`
	Conditions []byte    `gorm:"column:conditions"    json:"conditions"`
	CreatedAt  time.Time `gorm:"column:created_at"    json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"    json:"updated_at"`
}

// TableName returns the table name for PolicyModel.
func (PolicyModel) TableName() string {
	return "iam_policies"
}

// PolicyPool implements persistence.PolicyPersister using GORM.
type PolicyPool struct {
	db *Persister
}

// NewPolicyPool creates a new policy pool.
func NewPolicyPool(db *Persister) *PolicyPool {
	return &PolicyPool{db: db}
}

// GetPolicy retrieves a policy by ID.
func (p *PolicyPool) GetPolicy(ctx context.Context, id string) (*persistence.Policy, error) {
	var m PolicyModel
	if err := p.db.Connection(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return p.modelToDomain(&m), nil
}

// ListPolicies lists policies with pagination.
func (p *PolicyPool) ListPolicies(ctx context.Context, networkID string, limit, offset int) ([]*persistence.Policy, int, error) {
	var ms []PolicyModel
	var total int64

	query := p.db.Connection(ctx)
	if networkID != "" {
		query = query.Where("nid = ?", networkID)
	}

	if err := query.Model(&PolicyModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&ms).Error; err != nil {
		return nil, 0, err
	}

	policies := make([]*persistence.Policy, len(ms))
	for i := range ms {
		policies[i] = p.modelToDomain(&ms[i])
	}
	return policies, int(total), nil
}

// CreatePolicy creates a new policy.
func (p *PolicyPool) CreatePolicy(ctx context.Context, policy *persistence.Policy) error {
	m := p.domainToModel(policy)
	return p.db.Connection(ctx).Create(m).Error
}

// UpdatePolicy updates a policy.
func (p *PolicyPool) UpdatePolicy(ctx context.Context, policy *persistence.Policy) error {
	m := p.domainToModel(policy)
	return p.db.Connection(ctx).Model(m).Where("id = ?", policy.ID).Updates(m).Error
}

// DeletePolicy deletes a policy.
func (p *PolicyPool) DeletePolicy(ctx context.Context, id string) error {
	return p.db.Connection(ctx).Where("id = ?", id).Delete(&PolicyModel{}).Error
}

func (p *PolicyPool) modelToDomain(m *PolicyModel) *persistence.Policy {
	return &persistence.Policy{
		ID:         m.ID,
		NetworkID:  m.NetworkID,
		Name:       m.Name,
		Type:       m.Type,
		Subjects:   m.Subjects,
		Effect:     m.Effect,
		Actions:    m.Actions,
		Resources:  m.Resources,
		Conditions: m.Conditions,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func (p *PolicyPool) domainToModel(r *persistence.Policy) *PolicyModel {
	return &PolicyModel{
		ID:         r.ID,
		NetworkID:  r.NetworkID,
		Name:       r.Name,
		Type:       r.Type,
		Subjects:   r.Subjects,
		Effect:     r.Effect,
		Actions:    r.Actions,
		Resources:  r.Resources,
		Conditions: r.Conditions,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

// Ensure PolicyPool implements persistence.PolicyPersister.
var _ persistence.PolicyPersister = (*PolicyPool)(nil)
