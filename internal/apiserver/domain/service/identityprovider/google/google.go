// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package google

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/oauth2"

	"github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/options"
)

const (
	userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	authURL     = "https://accounts.google.com/o/oauth2/auth"
	tokenURL    = "https://oauth2.googleapis.com/token"
)

func init() {
	identityprovider.RegisterOAuthProvider(&googleProviderFactory{})
}

type google struct {
	// ClientID is the application's ID.
	ClientID string `json:"clientID" mapstructure:"clientID"`

	// ClientSecret is the application's secret.
	ClientSecret string `json:"clientSecret" mapstructure:"clientSecret"`

	// Endpoint contains the resource server's token endpoint
	// URLs. These are constants specific to each server and are
	// often available via site-specific packages, such as
	// google.Endpoint or github.endpoint.
	Endpoint endpoint `json:"endpoint" mapstructure:"endpoint"`

	// RedirectURL is the URL to redirect users going through
	// the OAuth flow, after the resource owner's URLs.
	RedirectURL string `json:"redirectURL" mapstructure:"redirectURL"`

	// Used to turn off TLS certificate checks
	InsecureSkipVerify bool `json:"insecureSkipVerify" mapstructure:"insecureSkipVerify"`

	// HTTP proxy to use for requests
	ProxyURL string `json:"proxyURL" mapstructure:"proxyURL"`

	// Scope specifies optional requested permissions.
	Scopes []string `json:"scopes" mapstructure:"scopes"`

	Config *oauth2.Config `json:"-" mapstructure:"-"`
}

// endpoint represents an OAuth 2.0 provider's authorization and token
// endpoint URLs.
type endpoint struct {
	AuthURL     string `json:"authURL" mapstructure:"authURL"`
	TokenURL    string `json:"tokenURL" mapstructure:"tokenURL"`
	UserInfoURL string `json:"userInfoURL" mapstructure:"userInfoURL"`
}

type googleIdentity struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	HD            string `json:"hd"`
}

type googleProviderFactory struct{}

func (g *googleProviderFactory) Type() v1.IdentityProviderType {
	return v1.GoogleIdentityProvider
}

func (g *googleProviderFactory) Create(opts options.DynamicOptions) (identityprovider.OAuthProvider, error) {
	var google google
	if err := mapstructure.Decode(opts, &google); err != nil {
		return nil, err
	}

	if google.Endpoint.AuthURL == "" {
		google.Endpoint.AuthURL = authURL
	}
	if google.Endpoint.TokenURL == "" {
		google.Endpoint.TokenURL = tokenURL
	}
	if google.Endpoint.UserInfoURL == "" {
		google.Endpoint.UserInfoURL = userInfoURL
	}
	// Set default scopes if none provided
	if len(google.Scopes) == 0 {
		google.Scopes = []string{"profile", "email"}
	}
	// fixed options
	opts["endpoint"] = options.DynamicOptions{
		"authURL":     google.Endpoint.AuthURL,
		"tokenURL":    google.Endpoint.TokenURL,
		"userInfoURL": google.Endpoint.UserInfoURL,
	}
	google.Config = &oauth2.Config{
		ClientID:     google.ClientID,
		ClientSecret: google.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  google.Endpoint.AuthURL,
			TokenURL: google.Endpoint.TokenURL,
		},
		RedirectURL: google.RedirectURL,
		Scopes:      google.Scopes,
	}
	return &google, nil
}

func (g googleIdentity) GetUserID() string {
	return g.ID
}

func (g googleIdentity) GetUsername() string {
	return g.Email
}

func (g googleIdentity) GetEmail() string {
	return g.Email
}

func (g googleIdentity) GetAvatar() string {
	return g.Picture
}

func (g *google) IdentityExchangeCallback(req *http.Request) (identityprovider.Identity, error) {
	// OAuth2 callback, see also https://tools.ietf.org/html/rfc6749#section-4.1.2
	code := req.URL.Query().Get("code")
	ctx := req.Context()

	// Create HTTP client with TLS configuration and proxy if needed
	var transport *http.Transport
	if g.InsecureSkipVerify {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	} else {
		transport = &http.Transport{}
	}

	// Configure proxy
	// Priority: 1. ProxyURL from config, 2. Environment variables (via ProxyFromEnvironment)
	transport.Proxy = http.ProxyFromEnvironment
	if g.ProxyURL != "" {
		parsedProxyURL, err := url.Parse(g.ProxyURL)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy URL: %w", err)
		}
		transport.Proxy = http.ProxyURL(parsedProxyURL)
	}

	httpClient := &http.Client{
		Transport: transport,
	}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)

	// Exchange the authorization code for an access token
	token, err := g.Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	// Create request to get user info
	userInfoReq, err := http.NewRequest("GET", g.Endpoint.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}
	userInfoReq.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// Make the request using the same HTTP client
	resp, err := httpClient.Do(userInfoReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: status %d, body: %s", resp.StatusCode, string(body))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var googleIdentity googleIdentity
	err = json.Unmarshal(data, &googleIdentity)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	// Validate that we have at least an ID
	if googleIdentity.ID == "" {
		return nil, errors.New("invalid user info: missing user ID")
	}

	return googleIdentity, nil
}
