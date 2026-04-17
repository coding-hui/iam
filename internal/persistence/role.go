// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package persistence

import (
	"context"
	"time"
)

// Role represents a role in the system.
// Domain model with no persistence-specific tags (Ory style).
type Role struct {
	ID          string
	NetworkID   string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RolePersister defines the interface for role persistence operations.
type RolePersister interface {
	GetRole(ctx context.Context, id string) (*Role, error)
	ListRoles(ctx context.Context, networkID string, limit, offset int) ([]*Role, int, error)
	CreateRole(ctx context.Context, role *Role) error
	UpdateRole(ctx context.Context, role *Role) error
	DeleteRole(ctx context.Context, id string) error
}
