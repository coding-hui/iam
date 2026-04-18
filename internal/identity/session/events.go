// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package session

// Session events.
const (
	EventSessionCreated  = "session.created"
	EventSessionRevoked  = "session.revoked"
	EventSessionExtended = "session.extended"
)

// SessionEvent represents a session-related event.
type SessionEvent struct {
	Type       string
	SessionID  string
	IdentityID string
	NetworkID  string
	Outcome    string
	Metadata   map[string]any
}
