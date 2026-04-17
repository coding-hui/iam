// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package role

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ManagerImpl implements role.Manager.
type ManagerImpl struct {
	pool     Pool
	privPool PrivilegedPool
}

// NewManagerImpl creates a new role manager.
func NewManagerImpl(pool Pool, privPool PrivilegedPool) *ManagerImpl {
	return &ManagerImpl{
		pool:     pool,
		privPool: privPool,
	}
}

// CreateRole creates a new role.
func (m *ManagerImpl) CreateRole(ctx context.Context, req *CreateRoleRequest) (*Role, error) {
	now := time.Now()
	r := &Role{
		ID:          uuid.New(),
		NetworkID:   req.NetworkID,
		Name:        req.Name,
		Description: req.Description,
		InheritFrom: req.InheritFrom,
		Extra:       req.Extra,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := m.privPool.CreateRole(ctx, r); err != nil {
		return nil, err
	}

	return r, nil
}

// GetRole retrieves a role by ID.
func (m *ManagerImpl) GetRole(ctx context.Context, id uuid.UUID) (*Role, error) {
	return m.pool.GetRole(ctx, id)
}

// ListRoles lists all roles for a network.
func (m *ManagerImpl) ListRoles(ctx context.Context, networkID uuid.UUID) ([]*Role, error) {
	roles, _, err := m.pool.ListRoles(ctx, networkID, 100, 0)
	return roles, err
}

// UpdateRole updates a role.
func (m *ManagerImpl) UpdateRole(ctx context.Context, id uuid.UUID, req *UpdateRoleRequest) (*Role, error) {
	r, err := m.pool.GetRole(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		r.Name = req.Name
	}
	if req.Description != "" {
		r.Description = req.Description
	}
	if req.InheritFrom != nil {
		r.InheritFrom = req.InheritFrom
	}
	if req.Extra != nil {
		r.Extra = req.Extra
	}
	r.UpdatedAt = time.Now()

	if err := m.privPool.UpdateRole(ctx, r); err != nil {
		return nil, err
	}

	return r, nil
}

// DeleteRole deletes a role.
func (m *ManagerImpl) DeleteRole(ctx context.Context, id uuid.UUID) error {
	networkID := uuid.Nil
	return m.privPool.DeleteRole(ctx, networkID, id)
}

// Ensure ManagerImpl implements Manager.
var _ Manager = (*ManagerImpl)(nil)
