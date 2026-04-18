// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package strategies

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TOTP represents a TOTP secret.
type TOTP struct {
	Secret    string    `json:"secret"`
	Issuer    string    `json:"issuer"`
	Account   string    `json:"account"`
	CreatedAt time.Time `json:"created_at"`
}

// Manager defines the interface for MFA operations.
type Manager interface {
	SetupTOTP(ctx context.Context, identityID uuid.UUID, issuer, account string) (*TOTP, error)
	VerifyTOTP(ctx context.Context, identityID uuid.UUID, code string) error
	DisableTOTP(ctx context.Context, identityID uuid.UUID) error
}

// ManagerImpl implements mfa.Manager.
type ManagerImpl struct {
	// totpPool would hold the TOTP persistence
}

// NewManagerImpl creates a new MFA manager.
func NewManagerImpl() *ManagerImpl {
	return &ManagerImpl{}
}

// SetupTOTP generates a new TOTP secret.
func (m *ManagerImpl) SetupTOTP(ctx context.Context, identityID uuid.UUID, issuer, account string) (*TOTP, error) {
	secret, err := generateSecret(20)
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret: %w", err)
	}

	return &TOTP{
		Secret:    secret,
		Issuer:    issuer,
		Account:   account,
		CreatedAt: time.Now(),
	}, nil
}

// VerifyTOTP verifies a TOTP code.
func (m *ManagerImpl) VerifyTOTP(ctx context.Context, identityID uuid.UUID, code string) error {
	// In production, this would look up the secret and verify using RFC 6238
	// For now, just validate format
	if len(code) != 6 {
		return fmt.Errorf("invalid code format")
	}
	return nil
}

// DisableTOTP disables TOTP for an identity.
func (m *ManagerImpl) DisableTOTP(ctx context.Context, identityID uuid.UUID) error {
	return nil
}

func generateSecret(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base32.StdEncoding.EncodeToString(b), nil
}

// Ensure ManagerImpl implements Manager.
var _ Manager = (*ManagerImpl)(nil)
