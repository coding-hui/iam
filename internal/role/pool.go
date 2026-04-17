// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package role

import (
	"context"

	"github.com/google/uuid"

	"github.com/coding-hui/iam/internal/persistence"
)

// rolePool implements Pool using persistence.RolePersister.
type rolePool struct {
	persister rolePersister
}

// rolePersister is the persistence interface for role operations.
type rolePersister interface {
	GetRole(ctx context.Context, id string) (*persistence.Role, error)
	ListRoles(ctx context.Context, networkID string, limit, offset int) ([]*persistence.Role, int, error)
	CreateRole(ctx context.Context, role *persistence.Role) error
	UpdateRole(ctx context.Context, role *persistence.Role) error
	DeleteRole(ctx context.Context, id string) error
}

// NewPool creates a new role pool.
func NewPool(p rolePersister) Pool {
	return &rolePool{persister: p}
}

// GetRole retrieves a role by ID.
func (p *rolePool) GetRole(ctx context.Context, id uuid.UUID) (*Role, error) {
	m, err := p.persister.GetRole(ctx, id.String())
	if err != nil {
		return nil, err
	}
	return p.modelToDomain(m), nil
}

// GetRoleByNetworkID retrieves a role by network ID and role ID.
func (p *rolePool) GetRoleByNetworkID(ctx context.Context, networkID, id uuid.UUID) (*Role, error) {
	m, err := p.persister.GetRole(ctx, id.String())
	if err != nil {
		return nil, err
	}
	if m.NetworkID != networkID.String() {
		return nil, ErrRoleNotFound
	}
	return p.modelToDomain(m), nil
}

// ListRoles lists roles with pagination.
func (p *rolePool) ListRoles(ctx context.Context, networkID uuid.UUID, limit, offset int) ([]*Role, int, error) {
	ms, total, err := p.persister.ListRoles(ctx, networkID.String(), limit, offset)
	if err != nil {
		return nil, 0, err
	}
	roles := make([]*Role, len(ms))
	for i := range ms {
		roles[i] = p.modelToDomain(ms[i])
	}
	return roles, total, nil
}

func (p *rolePool) modelToDomain(m *persistence.Role) *Role {
	if m == nil {
		return nil
	}
	return &Role{
		ID:          parseUUID(m.ID),
		NetworkID:   parseUUID(m.NetworkID),
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
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

// Ensure rolePool implements Pool.
var _ Pool = (*rolePool)(nil)
