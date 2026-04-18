// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package strategies

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"

	"github.com/coding-hui/iam/internal/identity"
	"github.com/coding-hui/iam/internal/identity/session"
)

// OAuthProvider represents an OAuth2 provider.
type OAuthProvider string

const (
	OAuthProviderGitHub OAuthProvider = "github"
	OAuthProviderGoogle OAuthProvider = "google"
)

// OAuthState stores OAuth flow state.
type OAuthState struct {
	State       string    `json:"state"`
	Provider    string    `json:"provider"`
	IdentityID  uuid.UUID `json:"identity_id,omitempty"`
	RedirectURL string    `json:"redirect_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// OAuthAuthenticator implements OAuth2-based authentication.
type OAuthAuthenticator struct {
	identityPool identity.PrivilegedPool
	sessionPool  sessionPool
	hasher       identity.Hasher
	providers    map[OAuthProvider]*oauth2.Config
	stateStore   map[string]*OAuthState
}

// NewOAuthAuthenticator creates a new OAuth authenticator.
func NewOAuthAuthenticator(
	identityPool identity.PrivilegedPool,
	sessionPool sessionPool,
	hasher identity.Hasher,
) *OAuthAuthenticator {
	return &OAuthAuthenticator{
		identityPool: identityPool,
		sessionPool:  sessionPool,
		hasher:       hasher,
		providers:    make(map[OAuthProvider]*oauth2.Config),
		stateStore:   make(map[string]*OAuthState),
	}
}

// RegisterProvider registers an OAuth2 provider configuration.
func (a *OAuthAuthenticator) RegisterProvider(provider OAuthProvider, config *oauth2.Config) {
	a.providers[provider] = config
}

// InitiateOAuthFlow initiates an OAuth2 flow and returns the authorization URL.
func (a *OAuthAuthenticator) InitiateOAuthFlow(ctx context.Context, provider OAuthProvider, redirectURL string) (string, *OAuthState, error) {
	cfg, ok := a.providers[provider]
	if !ok {
		return "", nil, fmt.Errorf("unsupported OAuth provider: %s", provider)
	}

	stateStr, err := generateState()
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate state: %w", err)
	}

	state := &OAuthState{
		State:       stateStr,
		Provider:    string(provider),
		RedirectURL: redirectURL,
		CreatedAt:   time.Now(),
	}

	// Store state in memory
	a.stateStore[stateStr] = state

	authURL := cfg.AuthCodeURL(stateStr,
		oauth2.AccessTypeOnline,
		oauth2.SetAuthURLParam("redirect_uri", cfg.RedirectURL),
	)

	return authURL, state, nil
}

// HandleOAuthCallback handles the OAuth2 callback.
func (a *OAuthAuthenticator) HandleOAuthCallback(ctx context.Context, provider OAuthProvider, code string) (*AuthenticateResponse, error) {
	cfg, ok := a.providers[provider]
	if !ok {
		return nil, fmt.Errorf("unsupported OAuth provider: %s", provider)
	}

	token, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Get user info from provider
	userInfo, err := a.fetchUserInfo(ctx, provider, token)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %w", err)
	}

	// Find or create identity
	ident, err := a.findOrCreateIdentity(ctx, provider, userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create identity: %w", err)
	}

	// Create session
	sess, err := a.createSession(ctx, ident, userInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &AuthenticateResponse{
		SessionID:  sess.ID,
		IdentityID: ident.ID,
		ExpiresAt:  sess.ExpiresAt.Unix(),
	}, nil
}

type oauthUserInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Provider string `json:"provider"`
}

func (a *OAuthAuthenticator) fetchUserInfo(ctx context.Context, provider OAuthProvider, token *oauth2.Token) (*oauthUserInfo, error) {
	switch provider {
	case OAuthProviderGitHub:
		return a.fetchGitHubUserInfo(ctx, token)
	case OAuthProviderGoogle:
		return a.fetchGoogleUserInfo(ctx, token)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

func (a *OAuthAuthenticator) fetchGitHubUserInfo(ctx context.Context, token *oauth2.Token) (*oauthUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo oauthUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	userInfo.Provider = string(OAuthProviderGitHub)
	return &userInfo, nil
}

func (a *OAuthAuthenticator) fetchGoogleUserInfo(ctx context.Context, token *oauth2.Token) (*oauthUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo oauthUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	userInfo.Provider = string(OAuthProviderGoogle)
	return &userInfo, nil
}

func (a *OAuthAuthenticator) findOrCreateIdentity(ctx context.Context, provider OAuthProvider, userInfo *oauthUserInfo) (*identity.Identity, error) {
	identifier := fmt.Sprintf("%s:%s", provider, userInfo.ID)

	// Find existing identity by OAuth provider ID
	existing, _, err := a.sessionPool.FindCredentialsByIdentifier(ctx, identity.CredentialsTypeAPIKey, identifier)
	if err == nil && existing != nil {
		return existing, nil
	}

	// Create new identity
	traits := map[string]any{
		"email":   userInfo.Email,
		"name":    userInfo.Name,
		"picture": userInfo.Picture,
	}
	traitsJSON, _ := json.Marshal(traits)

	newIdentity := &identity.Identity{
		ID:        uuid.New(),
		SchemaID:  "oauth",
		Traits:    traitsJSON,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.identityPool.CreateIdentity(ctx, newIdentity); err != nil {
		return nil, err
	}

	// Create OAuth credentials
	credConfig, _ := json.Marshal(map[string]any{
		"provider":    provider,
		"provider_id": userInfo.ID,
	})
	newCred := &identity.Credentials{
		ID:          uuid.New(),
		IdentityID:  newIdentity.ID,
		Type:        identity.CredentialsTypeAPIKey,
		Identifiers: []string{identifier},
		Config:      credConfig,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := a.identityPool.CreateCredentials(ctx, newCred); err != nil {
		return nil, err
	}

	return newIdentity, nil
}

func (a *OAuthAuthenticator) createSession(ctx context.Context, ident *identity.Identity, userInfo *oauthUserInfo) (*session.Session, error) {
	sess := &session.Session{
		ID:              uuid.New(),
		IdentityID:      ident.ID,
		Active:          true,
		ExpiresAt:       time.Now().Add(24 * time.Hour),
		AuthenticatedAt: time.Now(),
		UserAgent:       "OAuth/" + userInfo.Provider,
		ClientIP:        "",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := a.sessionPool.CreateSession(ctx, sess); err != nil {
		return nil, err
	}

	return sess, nil
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// sessionPool interface for session operations used by OAuth authenticator.
type sessionPool interface {
	FindCredentialsByIdentifier(ctx context.Context, credType identity.CredentialsType, identifier string) (*identity.Identity, *identity.Credentials, error)
	CreateSession(ctx context.Context, s *session.Session) error
}
