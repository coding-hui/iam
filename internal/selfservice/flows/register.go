// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package flows

import (
	"time"

	"github.com/google/uuid"
)

// RegisterFlow represents a registration flow.
type RegisterFlow struct {
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
	// SchemaID is the identity schema to use for this registration.
	SchemaID string `json:"schema_id"`
	// CreatedIdentityID is the identity ID created after successful registration.
	CreatedIdentityID uuid.UUID `json:"created_identity_id,omitempty"`
}

// NewRegisterFlow creates a new registration flow.
func NewRegisterFlow(requestURL, schemaID string) *RegisterFlow {
	if schemaID == "" {
		schemaID = "default"
	}
	return &RegisterFlow{
		ID:         uuid.New(),
		ExpiresAt:  time.Now().Add(30 * time.Minute),
		IssuedAt:   time.Now(),
		RequestURL: requestURL,
		State:      FlowStateChooseMethod,
		SchemaID:   schemaID,
	}
}

// IsExpired checks if the flow has expired.
func (f *RegisterFlow) IsExpired() bool {
	return time.Now().After(f.ExpiresAt)
}

// Completed checks if the flow has completed.
func (f *RegisterFlow) Completed() bool {
	return f.State == FlowStateCompleted
}
