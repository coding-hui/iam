// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package identity

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Identity represents a user identity in the system.
type Identity struct {
	ID        uuid.UUID       `json:"id"`
	NetworkID uuid.UUID       `json:"network_id"`
	SchemaID  string          `json:"schema_id"`
	Traits    json.RawMessage `json:"traits"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// Credentials represents authentication credentials for an identity.
type Credentials struct {
	ID          uuid.UUID       `json:"id"`
	IdentityID  uuid.UUID       `json:"identity_id"`
	Type        CredentialsType `json:"type"`
	Identifiers []string        `json:"identifiers"`
	Config      json.RawMessage `json:"config"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// CredentialsType represents the type of credentials.
type CredentialsType string

const (
	CredentialsTypePassword CredentialsType = "password"
	CredentialsTypeAPIKey   CredentialsType = "api_key"
	CredentialsTypeTOTP     CredentialsType = "totp"
)

// Pool defines the interface for reading identity data.
type Pool interface {
	GetIdentity(ctx context.Context, id uuid.UUID) (*Identity, error)
	GetIdentityByNetworkID(ctx context.Context, networkID, id uuid.UUID) (*Identity, error)
	GetIdentityByIdentifier(ctx context.Context, identifier string) (*Identity, error)
	ListIdentities(ctx context.Context, networkID uuid.UUID, limit, offset int) ([]*Identity, int, error)
	FindCredentialsByIdentifier(ctx context.Context, credType CredentialsType, identifier string) (*Identity, *Credentials, error)
}

// PrivilegedPool defines the interface for writing identity data.
type PrivilegedPool interface {
	Pool

	CreateIdentity(ctx context.Context, i *Identity) error
	UpdateIdentity(ctx context.Context, i *Identity) error
	DeleteIdentity(ctx context.Context, networkID, id uuid.UUID) error

	CreateCredentials(ctx context.Context, c *Credentials) error
	UpdateCredentials(ctx context.Context, c *Credentials) error
	DeleteCredentials(ctx context.Context, networkID uuid.UUID, id uuid.UUID, credType CredentialsType) error
}

// ListIdentitiesParams holds parameters for listing identities.
type ListIdentitiesParams struct {
	SchemaID string
	Page     int
	PageSize int
	Filters  map[string]string
}

// Manager defines the interface for identity business logic.
type Manager interface {
	CreateIdentity(ctx context.Context, req *CreateIdentityRequest) (*Identity, error)
	GetIdentity(ctx context.Context, id uuid.UUID) (*Identity, error)
	ListIdentities(ctx context.Context, networkID uuid.UUID, params ListIdentitiesParams) ([]*Identity, int, error)
	UpdateIdentity(ctx context.Context, id uuid.UUID, req *UpdateIdentityRequest) (*Identity, error)
	DeleteIdentity(ctx context.Context, id uuid.UUID) error

	AddCredentials(ctx context.Context, id uuid.UUID, req *AddCredentialsRequest) error
	DeleteCredentials(ctx context.Context, id uuid.UUID, credType CredentialsType) error
}

// CreateIdentityRequest holds data for creating a new identity.
type CreateIdentityRequest struct {
	SchemaID string          `json:"schema_id"`
	Traits   json.RawMessage `json:"traits"`
	Password string          `json:"password,omitempty"`
}

// UpdateIdentityRequest holds data for updating an identity.
type UpdateIdentityRequest struct {
	Traits json.RawMessage `json:"traits"`
}

// AddCredentialsRequest holds data for adding credentials.
type AddCredentialsRequest struct {
	Type        CredentialsType `json:"type"`
	Identifiers []string        `json:"identifiers"`
	Config      json.RawMessage `json:"config"`
	Password    string          `json:"password,omitempty"` // used for password type
}
