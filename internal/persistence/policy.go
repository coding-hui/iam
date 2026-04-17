// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package persistence

import (
	"context"
	"time"
)

// Policy represents a policy in the system.
// Domain model with no persistence-specific tags (Ory style).
type Policy struct {
	ID         string
	NetworkID  string
	Name       string
	Type       string
	Subjects   string
	Effect     string
	Actions    string
	Resources  string
	Conditions []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// PolicyPersister defines the interface for policy persistence operations.
type PolicyPersister interface {
	GetPolicy(ctx context.Context, id string) (*Policy, error)
	ListPolicies(ctx context.Context, networkID string, limit, offset int) ([]*Policy, int, error)
	CreatePolicy(ctx context.Context, policy *Policy) error
	UpdatePolicy(ctx context.Context, policy *Policy) error
	DeletePolicy(ctx context.Context, id string) error
}
