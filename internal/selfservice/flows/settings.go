// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package flows

import (
	"time"

	"github.com/google/uuid"
)

// SettingsFlow represents an account settings flow.
type SettingsFlow struct {
	// ID is the unique identifier of the flow.
	ID uuid.UUID `json:"id"`
	// ExpiresAt is the time when the flow expires.
	ExpiresAt time.Time `json:"expires_at"`
	// IssuedAt is the time when the flow was issued.
	IssuedAt time.Time `json:"issued_at"`
	// RequestURL is the URL where the flow was initiated.
	RequestURL string `json:"request_url"`
	// State is the current state of the flow.
	State FlowState `json:"state"`
	// IdentityID is the identity being updated.
	IdentityID uuid.UUID `json:"identity_id"`
	// SchemaID is the schema to use for validation.
	SchemaID string `json:"schema_id"`
}

// NewSettingsFlow creates a new account settings flow.
func NewSettingsFlow(requestURL string, identityID uuid.UUID, schemaID string) *SettingsFlow {
	if schemaID == "" {
		schemaID = "default"
	}
	return &SettingsFlow{
		ID:         uuid.New(),
		ExpiresAt:  time.Now().Add(30 * time.Minute),
		IssuedAt:   time.Now(),
		RequestURL: requestURL,
		State:      FlowStateChooseMethod,
		IdentityID: identityID,
		SchemaID:   schemaID,
	}
}

// IsExpired checks if the flow has expired.
func (f *SettingsFlow) IsExpired() bool {
	return time.Now().After(f.ExpiresAt)
}

// Completed checks if the flow has completed.
func (f *SettingsFlow) Completed() bool {
	return f.State == FlowStateCompleted
}
