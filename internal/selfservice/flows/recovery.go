// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package flows

import (
	"time"

	"github.com/google/uuid"
)

// RecoveryFlow represents a password recovery flow.
type RecoveryFlow struct {
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
	// IdentityID is the identity being recovered (set after email verification).
	IdentityID uuid.UUID `json:"identity_id,omitempty"`
	// Email is the email address used for recovery.
	Email string `json:"email,omitempty"`
	// RecoveryToken is the token sent to the email.
	RecoveryToken string `json:"recovery_token,omitempty"`
}

// NewRecoveryFlow creates a new password recovery flow.
func NewRecoveryFlow(requestURL string) *RecoveryFlow {
	return &RecoveryFlow{
		ID:         uuid.New(),
		ExpiresAt:  time.Now().Add(60 * time.Minute),
		IssuedAt:   time.Now(),
		RequestURL: requestURL,
		State:      FlowStateChooseMethod,
	}
}

// IsExpired checks if the flow has expired.
func (f *RecoveryFlow) IsExpired() bool {
	return time.Now().After(f.ExpiresAt)
}

// Completed checks if the flow has completed.
func (f *RecoveryFlow) Completed() bool {
	return f.State == FlowStateCompleted
}
