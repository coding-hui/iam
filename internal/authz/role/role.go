// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package role

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Role represents a role in the system.
type Role struct {
	ID          uuid.UUID       `json:"id"`
	NetworkID   uuid.UUID       `json:"network_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InheritFrom []uuid.UUID     `json:"inherit_from,omitempty"`
	Extra       json.RawMessage `json:"extra,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// Pool defines the interface for reading role data.
type Pool interface {
	GetRole(ctx context.Context, id uuid.UUID) (*Role, error)
	GetRoleByNetworkID(ctx context.Context, networkID, id uuid.UUID) (*Role, error)
	ListRoles(ctx context.Context, networkID uuid.UUID, limit, offset int) ([]*Role, int, error)
}

// PrivilegedPool defines the interface for writing role data.
type PrivilegedPool interface {
	Pool

	CreateRole(ctx context.Context, r *Role) error
	UpdateRole(ctx context.Context, r *Role) error
	DeleteRole(ctx context.Context, networkID, id uuid.UUID) error
}

// Manager defines the interface for role business logic.
type Manager interface {
	CreateRole(ctx context.Context, req *CreateRoleRequest) (*Role, error)
	GetRole(ctx context.Context, id uuid.UUID) (*Role, error)
	ListRoles(ctx context.Context, networkID uuid.UUID) ([]*Role, error)
	UpdateRole(ctx context.Context, id uuid.UUID, req *UpdateRoleRequest) (*Role, error)
	DeleteRole(ctx context.Context, id uuid.UUID) error
}

// CreateRoleRequest holds data for creating a new role.
type CreateRoleRequest struct {
	NetworkID   uuid.UUID       `json:"network_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InheritFrom []uuid.UUID     `json:"inherit_from,omitempty"`
	Extra       json.RawMessage `json:"extra,omitempty"`
}

// UpdateRoleRequest holds data for updating a role.
type UpdateRoleRequest struct {
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	InheritFrom []uuid.UUID     `json:"inherit_from,omitempty"`
	Extra       json.RawMessage `json:"extra,omitempty"`
}
