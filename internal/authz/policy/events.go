// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package policy

// Policy events.
const (
	EventPolicyCreated = "policy.created"
	EventPolicyUpdated = "policy.updated"
	EventPolicyDeleted = "policy.deleted"
)

// PolicyEvent represents a policy-related event.
type PolicyEvent struct {
	Type      string
	PolicyID  string
	NetworkID string
	Outcome   string
	Metadata  map[string]any
}
