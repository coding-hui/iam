// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authn

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/coding-hui/iam/internal/identity"
	"github.com/coding-hui/iam/internal/session"
)

// Authenticator defines the interface for authentication.
type Authenticator interface {
	Authenticate(ctx context.Context, req *AuthenticateRequest) (*AuthenticateResponse, error)
}

// AuthenticateRequest holds data for authentication.
type AuthenticateRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
	UserAgent  string `json:"user_agent"`
	ClientIP   string `json:"client_ip"`
}

// AuthenticateResponse holds the result of authentication.
type AuthenticateResponse struct {
	SessionID  uuid.UUID `json:"session_id"`
	IdentityID uuid.UUID `json:"identity_id"`
	ExpiresAt  int64     `json:"expires_at"`
}

// PasswordAuthenticator implements password-based authentication.
type PasswordAuthenticator struct {
	identityPool identity.PrivilegedPool
	sessionPool  session.PrivilegedPool
	hasher       identity.Hasher
}

// NewPasswordAuthenticator creates a new password authenticator.
func NewPasswordAuthenticator(
	identityPool identity.PrivilegedPool,
	sessionPool session.PrivilegedPool,
	hasher identity.Hasher,
) *PasswordAuthenticator {
	return &PasswordAuthenticator{
		identityPool: identityPool,
		sessionPool:  sessionPool,
		hasher:       hasher,
	}
}

// Authenticate authenticates a user using identifier and password.
func (a *PasswordAuthenticator) Authenticate(ctx context.Context, req *AuthenticateRequest) (*AuthenticateResponse, error) {
	// 1. Find identity and credentials by identifier
	_, cred, err := a.identityPool.FindCredentialsByIdentifier(ctx, identity.CredentialsTypePassword, req.Identifier)
	if err != nil {
		return nil, identity.ErrInvalidCredentials
	}

	// 2. Verify password
	if err := a.hasher.Verify(req.Password, cred.Config); err != nil {
		return nil, identity.ErrInvalidCredentials
	}

	// 3. Create session
	sess := &session.Session{
		ID:              uuid.New(),
		IdentityID:      cred.IdentityID,
		Active:          true,
		ExpiresAt:       time.Now().Add(24 * time.Hour),
		AuthenticatedAt: time.Now(),
		UserAgent:       req.UserAgent,
		ClientIP:        req.ClientIP,
	}

	if err := a.sessionPool.CreateSession(ctx, sess); err != nil {
		return nil, err
	}

	return &AuthenticateResponse{
		SessionID:  sess.ID,
		IdentityID: sess.IdentityID,
		ExpiresAt:  sess.ExpiresAt.Unix(),
	}, nil
}
