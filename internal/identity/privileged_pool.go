// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package identity

import (
	"context"

	"github.com/google/uuid"

	"github.com/coding-hui/iam/internal/persistence"
)

// privilegedPool implements PrivilegedPool.
type privilegedPool struct {
	*identityPool
}

// NewPrivilegedPool creates a new identity privileged pool.
func NewPrivilegedPool(p identityPersister) PrivilegedPool {
	return &privilegedPool{
		identityPool: &identityPool{persister: p},
	}
}

// CreateIdentity creates a new identity.
func (p *privilegedPool) CreateIdentity(ctx context.Context, i *Identity) error {
	m := p.domainToModel(i)
	return p.persister.CreateIdentity(ctx, m)
}

// UpdateIdentity updates an identity.
func (p *privilegedPool) UpdateIdentity(ctx context.Context, i *Identity) error {
	m := p.domainToModel(i)
	return p.persister.UpdateIdentity(ctx, m)
}

// DeleteIdentity deletes an identity.
func (p *privilegedPool) DeleteIdentity(ctx context.Context, networkID, id uuid.UUID) error {
	return p.persister.DeleteIdentity(ctx, id.String())
}

// CreateCredentials creates credentials.
func (p *privilegedPool) CreateCredentials(ctx context.Context, c *Credentials) error {
	return nil
}

// UpdateCredentials updates credentials.
func (p *privilegedPool) UpdateCredentials(ctx context.Context, c *Credentials) error {
	return nil
}

// DeleteCredentials deletes credentials.
func (p *privilegedPool) DeleteCredentials(ctx context.Context, networkID uuid.UUID, id uuid.UUID, credType CredentialsType) error {
	return nil
}

func (p *privilegedPool) domainToModel(i *Identity) *persistence.Identity {
	return &persistence.Identity{
		ID:        i.ID.String(),
		NetworkID: i.NetworkID.String(),
		SchemaID:  i.SchemaID,
		Traits:    i.Traits,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
	}
}

// Ensure privilegedPool implements PrivilegedPool.
var _ PrivilegedPool = (*privilegedPool)(nil)
