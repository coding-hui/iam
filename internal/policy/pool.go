// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package policy

import (
	"context"

	"github.com/google/uuid"

	"github.com/coding-hui/iam/internal/persistence"
)

// policyPool implements Pool using persistence.PolicyPersister.
type policyPool struct {
	persister policyPersister
}

// policyPersister is the persistence interface for policy operations.
type policyPersister interface {
	GetPolicy(ctx context.Context, id string) (*persistence.Policy, error)
	ListPolicies(ctx context.Context, networkID string, limit, offset int) ([]*persistence.Policy, int, error)
	CreatePolicy(ctx context.Context, policy *persistence.Policy) error
	UpdatePolicy(ctx context.Context, policy *persistence.Policy) error
	DeletePolicy(ctx context.Context, id string) error
}

// NewPool creates a new policy pool.
func NewPool(p policyPersister) Pool {
	return &policyPool{persister: p}
}

// GetPolicy retrieves a policy by ID.
func (p *policyPool) GetPolicy(ctx context.Context, id uuid.UUID) (*Policy, error) {
	m, err := p.persister.GetPolicy(ctx, id.String())
	if err != nil {
		return nil, err
	}
	return p.modelToDomain(m), nil
}

// ListPolicies lists policies with pagination.
func (p *policyPool) ListPolicies(ctx context.Context, networkID uuid.UUID, limit, offset int) ([]*Policy, int, error) {
	ms, total, err := p.persister.ListPolicies(ctx, networkID.String(), limit, offset)
	if err != nil {
		return nil, 0, err
	}
	policies := make([]*Policy, len(ms))
	for i := range ms {
		policies[i] = p.modelToDomain(ms[i])
	}
	return policies, total, nil
}

func (p *policyPool) modelToDomain(m *persistence.Policy) *Policy {
	if m == nil {
		return nil
	}
	return &Policy{
		ID:         parseUUID(m.ID),
		NetworkID:  parseUUID(m.NetworkID),
		Name:       m.Name,
		Type:       PolicyType(m.Type),
		Subjects:   split(m.Subjects),
		Effect:     Effect(m.Effect),
		Actions:    split(m.Actions),
		Resources:  split(m.Resources),
		Conditions: m.Conditions,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func split(s string) []string {
	if s == "" {
		return nil
	}
	return []string{s}
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

// Ensure policyPool implements Pool.
var _ Pool = (*policyPool)(nil)

// privilegedPool implements PrivilegedPool.
type privilegedPool struct {
	*policyPool
}

// NewPrivilegedPool creates a new policy privileged pool.
func NewPrivilegedPool(p policyPersister) PrivilegedPool {
	return &privilegedPool{
		policyPool: &policyPool{persister: p},
	}
}

// CreatePolicy creates a new policy.
func (p *privilegedPool) CreatePolicy(ctx context.Context, r *Policy) error {
	m := p.domainToModel(r)
	return p.persister.CreatePolicy(ctx, m)
}

// UpdatePolicy updates a policy.
func (p *privilegedPool) UpdatePolicy(ctx context.Context, r *Policy) error {
	m := p.domainToModel(r)
	return p.persister.UpdatePolicy(ctx, m)
}

// DeletePolicy deletes a policy.
func (p *privilegedPool) DeletePolicy(ctx context.Context, networkID, id uuid.UUID) error {
	return p.persister.DeletePolicy(ctx, id.String())
}

func (p *privilegedPool) domainToModel(r *Policy) *persistence.Policy {
	return &persistence.Policy{
		ID:         r.ID.String(),
		NetworkID:  r.NetworkID.String(),
		Name:       r.Name,
		Type:       string(r.Type),
		Subjects:   join(r.Subjects),
		Effect:     string(r.Effect),
		Actions:    join(r.Actions),
		Resources:  join(r.Resources),
		Conditions: r.Conditions,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

func join(ss []string) string {
	if ss == nil {
		return ""
	}
	result := ""
	for i, s := range ss {
		if i > 0 {
			result += ","
		}
		result += s
	}
	return result
}

// Ensure privilegedPool implements PrivilegedPool.
var _ PrivilegedPool = (*privilegedPool)(nil)
