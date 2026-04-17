// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package identity

// Identity events.
const (
	EventIdentityCreated    = "identity.created"
	EventIdentityUpdated    = "identity.updated"
	EventIdentityDeleted    = "identity.deleted"
	EventCredentialsAdded   = "credentials.added"
	EventCredentialsDeleted = "credentials.deleted"
)

// IdentityEvent represents an identity-related event.
type IdentityEvent struct {
	Type       string
	IdentityID string
	NetworkID  string
	ActorID    string
	Outcome    string
	Metadata   map[string]any
}
