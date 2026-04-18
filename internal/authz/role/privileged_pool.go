// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package role

import (
	"context"

	"github.com/google/uuid"

	"github.com/coding-hui/iam/internal/persistence"
)

// privilegedPool implements PrivilegedPool.
type privilegedPool struct {
	*rolePool
}

// NewPrivilegedPool creates a new role privileged pool.
func NewPrivilegedPool(p rolePersister) PrivilegedPool {
	return &privilegedPool{
		rolePool: &rolePool{persister: p},
	}
}

// CreateRole creates a new role.
func (p *privilegedPool) CreateRole(ctx context.Context, r *Role) error {
	m := p.domainToModel(r)
	return p.persister.CreateRole(ctx, m)
}

// UpdateRole updates a role.
func (p *privilegedPool) UpdateRole(ctx context.Context, r *Role) error {
	m := p.domainToModel(r)
	return p.persister.UpdateRole(ctx, m)
}

// DeleteRole deletes a role.
func (p *privilegedPool) DeleteRole(ctx context.Context, networkID, id uuid.UUID) error {
	return p.persister.DeleteRole(ctx, id.String())
}

func (p *privilegedPool) domainToModel(r *Role) *persistence.Role {
	return &persistence.Role{
		ID:          r.ID.String(),
		NetworkID:   r.NetworkID.String(),
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// Ensure privilegedPool implements PrivilegedPool.
var _ PrivilegedPool = (*privilegedPool)(nil)
