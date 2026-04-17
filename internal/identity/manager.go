// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package identity

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ManagerImpl implements identity.Manager using identity.Pool.
type ManagerImpl struct {
	pool     Pool
	privPool PrivilegedPool
	hasher   Hasher
}

// NewManagerImpl creates a new identity manager.
func NewManagerImpl(pool Pool, privPool PrivilegedPool, hasher Hasher) *ManagerImpl {
	return &ManagerImpl{
		pool:     pool,
		privPool: privPool,
		hasher:   hasher,
	}
}

// CreateIdentity creates a new identity.
func (m *ManagerImpl) CreateIdentity(ctx context.Context, req *CreateIdentityRequest) (*Identity, error) {
	if req.SchemaID == "" {
		req.SchemaID = "default"
	}

	identity := &Identity{
		ID:        uuid.New(),
		SchemaID:  req.SchemaID,
		Traits:    req.Traits,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := m.privPool.CreateIdentity(ctx, identity); err != nil {
		return nil, err
	}

	return identity, nil
}

// GetIdentity retrieves an identity by ID.
func (m *ManagerImpl) GetIdentity(ctx context.Context, id uuid.UUID) (*Identity, error) {
	return m.pool.GetIdentity(ctx, id)
}

// ListIdentities lists identities with pagination.
func (m *ManagerImpl) ListIdentities(ctx context.Context, networkID uuid.UUID, params ListIdentitiesParams) ([]*Identity, int, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 20
	}
	offset := (params.Page - 1) * params.PageSize

	return m.pool.ListIdentities(ctx, networkID, params.PageSize, offset)
}

// UpdateIdentity updates an identity's traits.
func (m *ManagerImpl) UpdateIdentity(ctx context.Context, id uuid.UUID, req *UpdateIdentityRequest) (*Identity, error) {
	identity, err := m.pool.GetIdentity(ctx, id)
	if err != nil {
		return nil, err
	}

	identity.Traits = req.Traits
	identity.UpdatedAt = time.Now()

	if err := m.privPool.UpdateIdentity(ctx, identity); err != nil {
		return nil, err
	}

	return identity, nil
}

// DeleteIdentity deletes an identity.
func (m *ManagerImpl) DeleteIdentity(ctx context.Context, id uuid.UUID) error {
	networkID := uuid.Nil
	return m.privPool.DeleteIdentity(ctx, networkID, id)
}

// AddCredentials adds credentials to an identity.
func (m *ManagerImpl) AddCredentials(ctx context.Context, id uuid.UUID, req *AddCredentialsRequest) error {
	identity, err := m.pool.GetIdentity(ctx, id)
	if err != nil {
		return err
	}

	var config json.RawMessage
	if req.Type == CredentialsTypePassword && req.Config == nil {
		hash, err := m.hasher.Hash(req.Password)
		if err != nil {
			return err
		}
		config = hash
	} else {
		config = req.Config
	}

	cred := &Credentials{
		ID:          uuid.New(),
		IdentityID:  identity.ID,
		Type:        req.Type,
		Identifiers: req.Identifiers,
		Config:      config,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return m.privPool.CreateCredentials(ctx, cred)
}

// DeleteCredentials deletes credentials from an identity.
func (m *ManagerImpl) DeleteCredentials(ctx context.Context, id uuid.UUID, credType CredentialsType) error {
	networkID := uuid.Nil
	return m.privPool.DeleteCredentials(ctx, networkID, id, credType)
}

// Ensure ManagerImpl implements Manager.
var _ Manager = (*ManagerImpl)(nil)
