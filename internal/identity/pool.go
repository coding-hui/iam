// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package identity

import (
	"context"

	"github.com/google/uuid"

	"github.com/coding-hui/iam/internal/persistence"
)

// identityPool implements Pool using persistence.IdentityPersister.
type identityPool struct {
	persister identityPersister
}

// identityPersister is the persistence interface for identity operations.
type identityPersister interface {
	GetIdentity(ctx context.Context, id string) (*persistence.Identity, error)
	ListIdentities(ctx context.Context, networkID string, limit, offset int) ([]*persistence.Identity, int, error)
	CreateIdentity(ctx context.Context, identity *persistence.Identity) error
	UpdateIdentity(ctx context.Context, identity *persistence.Identity) error
	DeleteIdentity(ctx context.Context, id string) error
}

// NewPool creates a new identity pool.
func NewPool(p identityPersister) Pool {
	return &identityPool{persister: p}
}

// GetIdentity retrieves an identity by ID.
func (p *identityPool) GetIdentity(ctx context.Context, id uuid.UUID) (*Identity, error) {
	m, err := p.persister.GetIdentity(ctx, id.String())
	if err != nil {
		return nil, err
	}
	return p.modelToDomain(m), nil
}

// GetIdentityByNetworkID retrieves an identity by network ID and identity ID.
func (p *identityPool) GetIdentityByNetworkID(ctx context.Context, networkID, id uuid.UUID) (*Identity, error) {
	m, err := p.persister.GetIdentity(ctx, id.String())
	if err != nil {
		return nil, err
	}
	if m.NetworkID != networkID.String() {
		return nil, ErrIdentityNotFound
	}
	return p.modelToDomain(m), nil
}

// GetIdentityByIdentifier retrieves an identity by identifier (e.g., email).
func (p *identityPool) GetIdentityByIdentifier(ctx context.Context, identifier string) (*Identity, error) {
	return nil, ErrIdentityNotFound
}

// ListIdentities lists identities with pagination.
func (p *identityPool) ListIdentities(ctx context.Context, networkID uuid.UUID, limit, offset int) ([]*Identity, int, error) {
	ms, total, err := p.persister.ListIdentities(ctx, networkID.String(), limit, offset)
	if err != nil {
		return nil, 0, err
	}
	identities := make([]*Identity, len(ms))
	for i := range ms {
		identities[i] = p.modelToDomain(ms[i])
	}
	return identities, total, nil
}

// FindCredentialsByIdentifier finds an identity and credentials by identifier.
func (p *identityPool) FindCredentialsByIdentifier(ctx context.Context, credType CredentialsType, identifier string) (*Identity, *Credentials, error) {
	return nil, nil, ErrIdentityNotFound
}

func (p *identityPool) modelToDomain(m *persistence.Identity) *Identity {
	if m == nil {
		return nil
	}
	return &Identity{
		ID:        parseUUID(m.ID),
		NetworkID: parseUUID(m.NetworkID),
		SchemaID:  m.SchemaID,
		Traits:    m.Traits,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
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

// Ensure identityPool implements Pool.
var _ Pool = (*identityPool)(nil)
