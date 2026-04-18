// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package token

// Token events.
const (
	EventTokenCreated = "token.created"
	EventTokenRevoked = "token.revoked"
	EventTokenExpired = "token.expired"
)

// TokenEvent represents a token-related event.
type TokenEvent struct {
	Type       string
	TokenID    string
	IdentityID string
	TokenType  string
	Outcome    string
	Metadata   map[string]any
}
