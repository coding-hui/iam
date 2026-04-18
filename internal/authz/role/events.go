// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package role

// Role events.
const (
	EventRoleCreated = "role.created"
	EventRoleUpdated = "role.updated"
	EventRoleDeleted = "role.deleted"
)

// RoleEvent represents a role-related event.
type RoleEvent struct {
	Type      string
	RoleID    string
	NetworkID string
	Outcome   string
	Metadata  map[string]any
}
