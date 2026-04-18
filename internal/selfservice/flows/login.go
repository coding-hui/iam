// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package flows

import (
	"time"

	"github.com/google/uuid"
)

// FlowState represents the state of a self-service flow.
type FlowState string

const (
	// FlowStateChooseMethod - user chooses authentication method
	FlowStateChooseMethod FlowState = "choose_method"
	// FlowStateSentCode - verification code sent
	FlowStateSentCode FlowState = "sent_code"
	// FlowStatePassedChallenge - authentication challenge passed
	FlowStatePassedChallenge FlowState = "passed_challenge"
	// FlowStateCompleted - flow completed successfully
	FlowStateCompleted FlowState = "completed"
)

// LoginFlow represents a login flow.
type LoginFlow struct {
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
	// Active authenticator identity ID (set after successful authentication).
	ActiveIdentityID uuid.UUID `json:"active_identity_id,omitempty"`
	// RequestedAcr is the Authentication Context Class Reference requested.
	RequestedAcr string `json:"requested_acr,omitempty"`
}

// NewLoginFlow creates a new login flow.
func NewLoginFlow(requestURL string) *LoginFlow {
	return &LoginFlow{
		ID:         uuid.New(),
		ExpiresAt:  time.Now().Add(30 * time.Minute),
		IssuedAt:   time.Now(),
		RequestURL: requestURL,
		State:      FlowStateChooseMethod,
	}
}

// IsExpired checks if the flow has expired.
func (f *LoginFlow) IsExpired() bool {
	return time.Now().After(f.ExpiresAt)
}

// Completed checks if the flow has completed.
func (f *LoginFlow) Completed() bool {
	return f.State == FlowStateCompleted
}
