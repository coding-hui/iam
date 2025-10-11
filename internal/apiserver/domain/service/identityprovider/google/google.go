// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package google

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"

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
	if g.InsecureSkipVerify {
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		ctx = context.WithValue(ctx, oauth2.HTTPClient, client)
	}
	token, err := g.Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	resp, err := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token)).Get(g.Endpoint.UserInfoURL)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var googleIdentity googleIdentity
	err = json.Unmarshal(data, &googleIdentity)
	if err != nil {
		return nil, err
	}

	return googleIdentity, nil
}
