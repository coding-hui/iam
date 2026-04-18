// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package policy

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// PolicyType represents the type of policy.
type PolicyType string

const (
	PolicyTypeRole PolicyType = "role"
	PolicyTypeUser PolicyType = "user"
)

// Effect represents the effect of a policy.
type Effect string

const (
	EffectAllow Effect = "allow"
	EffectDeny  Effect = "deny"
)

// Policy represents a policy in the system.
type Policy struct {
	ID         uuid.UUID       `json:"id"`
	NetworkID  uuid.UUID       `json:"network_id"`
	Name       string          `json:"name"`
	Type       PolicyType      `json:"type"`
	Subjects   []string        `json:"subjects"`
	Effect     Effect          `json:"effect"`
	Actions    []string        `json:"actions"`
	Resources  []string        `json:"resources"`
	Conditions json.RawMessage `json:"conditions,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// Pool defines the interface for reading policy data.
type Pool interface {
	GetPolicy(ctx context.Context, id uuid.UUID) (*Policy, error)
	ListPolicies(ctx context.Context, networkID uuid.UUID, limit, offset int) ([]*Policy, int, error)
}

// PrivilegedPool defines the interface for writing policy data.
type PrivilegedPool interface {
	Pool

	CreatePolicy(ctx context.Context, p *Policy) error
	UpdatePolicy(ctx context.Context, p *Policy) error
	DeletePolicy(ctx context.Context, networkID, id uuid.UUID) error
}

// Manager defines the interface for policy business logic.
type Manager interface {
	CreatePolicy(ctx context.Context, req *CreatePolicyRequest) (*Policy, error)
	GetPolicy(ctx context.Context, id uuid.UUID) (*Policy, error)
	ListPolicies(ctx context.Context, networkID uuid.UUID) ([]*Policy, error)
	UpdatePolicy(ctx context.Context, id uuid.UUID, req *UpdatePolicyRequest) (*Policy, error)
	DeletePolicy(ctx context.Context, id uuid.UUID) error
}

// CreatePolicyRequest holds data for creating a new policy.
type CreatePolicyRequest struct {
	NetworkID  uuid.UUID       `json:"network_id"`
	Name       string          `json:"name"`
	Type       PolicyType      `json:"type"`
	Subjects   []string        `json:"subjects"`
	Effect     Effect          `json:"effect"`
	Actions    []string        `json:"actions"`
	Resources  []string        `json:"resources"`
	Conditions json.RawMessage `json:"conditions,omitempty"`
}

// UpdatePolicyRequest holds data for updating a policy.
type UpdatePolicyRequest struct {
	Name       string          `json:"name,omitempty"`
	Subjects   []string        `json:"subjects,omitempty"`
	Effect     Effect          `json:"effect,omitempty"`
	Actions    []string        `json:"actions,omitempty"`
	Resources  []string        `json:"resources,omitempty"`
	Conditions json.RawMessage `json:"conditions,omitempty"`
}
