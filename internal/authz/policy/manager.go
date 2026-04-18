// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package policy

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ManagerImpl implements policy.Manager.
type ManagerImpl struct {
	pool     Pool
	privPool PrivilegedPool
}

// NewManagerImpl creates a new policy manager.
func NewManagerImpl(pool Pool, privPool PrivilegedPool) *ManagerImpl {
	return &ManagerImpl{
		pool:     pool,
		privPool: privPool,
	}
}

// CreatePolicy creates a new policy.
func (m *ManagerImpl) CreatePolicy(ctx context.Context, req *CreatePolicyRequest) (*Policy, error) {
	now := time.Now()
	r := &Policy{
		ID:         uuid.New(),
		NetworkID:  req.NetworkID,
		Name:       req.Name,
		Type:       req.Type,
		Subjects:   req.Subjects,
		Effect:     req.Effect,
		Actions:    req.Actions,
		Resources:  req.Resources,
		Conditions: req.Conditions,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := m.privPool.CreatePolicy(ctx, r); err != nil {
		return nil, err
	}

	return r, nil
}

// GetPolicy retrieves a policy by ID.
func (m *ManagerImpl) GetPolicy(ctx context.Context, id uuid.UUID) (*Policy, error) {
	return m.pool.GetPolicy(ctx, id)
}

// ListPolicies lists all policies for a network.
func (m *ManagerImpl) ListPolicies(ctx context.Context, networkID uuid.UUID) ([]*Policy, error) {
	policies, _, err := m.pool.ListPolicies(ctx, networkID, 100, 0)
	return policies, err
}

// UpdatePolicy updates a policy.
func (m *ManagerImpl) UpdatePolicy(ctx context.Context, id uuid.UUID, req *UpdatePolicyRequest) (*Policy, error) {
	r, err := m.pool.GetPolicy(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		r.Name = req.Name
	}
	if req.Subjects != nil {
		r.Subjects = req.Subjects
	}
	if req.Effect != "" {
		r.Effect = req.Effect
	}
	if req.Actions != nil {
		r.Actions = req.Actions
	}
	if req.Resources != nil {
		r.Resources = req.Resources
	}
	if req.Conditions != nil {
		r.Conditions = req.Conditions
	}
	r.UpdatedAt = time.Now()

	if err := m.privPool.UpdatePolicy(ctx, r); err != nil {
		return nil, err
	}

	return r, nil
}

// DeletePolicy deletes a policy.
func (m *ManagerImpl) DeletePolicy(ctx context.Context, id uuid.UUID) error {
	networkID := uuid.Nil
	return m.privPool.DeletePolicy(ctx, networkID, id)
}

// Ensure ManagerImpl implements Manager.
var _ Manager = (*ManagerImpl)(nil)
