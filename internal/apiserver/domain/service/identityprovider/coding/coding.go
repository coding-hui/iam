// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package coding

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/oauth2"

	"github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider"
	"github.com/coding-hui/iam/internal/pkg/options"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

const (
	userInfoURL = "https://%s.coding.net/api/me"
	authURL     = "https://%s.coding.net/oauth_authorize.html"
	tokenURL    = "https://%s.coding.net/api/oauth/access_token"
)

func init() {
	identityprovider.RegisterOAuthProvider(&codingProviderFactory{})
}

type coding struct {
	// ClientID is the application's ID.
	ClientID string `json:"clientID" yaml:"clientID"`

	// ClientSecret is the application's secret.
	ClientSecret string `json:"clientSecret" yaml:"clientSecret"`

	// Endpoint contains the resource server's token endpoint
	// URLs. These are constants specific to each server and are
	// often available via site-specific packages, such as
	// google.Endpoint or github.endpoint.
	Endpoint endpoint `json:"endpoint" yaml:"endpoint"`

	// RedirectURL is the URL to redirect users going through
	// the OAuth flow, after the resource owner's URLs.
	RedirectURL string `json:"redirectURL" yaml:"redirectURL"`

	// Used to turn off TLS certificate checks
	InsecureSkipVerify bool `json:"insecureSkipVerify" yaml:"insecureSkipVerify"`

	// Scope specifies optional requested permissions.
	Scopes []string `json:"scopes" yaml:"scopes"`

	// Team team domain prefix.
	Team string `json:"team" yaml:"team"`

	Config *oauth2.Config `json:"-" yaml:"-"`
}

// endpoint represents an OAuth 2.0 provider's authorization and token
// endpoint URLs.
type endpoint struct {
	AuthURL     string `json:"authURL"     yaml:"authURL"`
	TokenURL    string `json:"tokenURL"    yaml:"tokenURL"`
	UserInfoURL string `json:"userInfoURL" yaml:"userInfoURL"`
}

type codingIdentity struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	NamePinyin      string `json:"name_pinyin"`
	Path            string `json:"path"`
	Team            string `json:"team"`
	EmailValidation int    `json:"email_validation"`
	Avatar          string `json:"avatar"`
	CreatedAt       uint64 `json:"created_at"`
	UpdatedAt       uint64 `json:"updated_at"`
	GlobalKey       string `json:"global_key"`
}

type codingProviderFactory struct {
}

func (g *codingProviderFactory) Type() v1.IdentityProviderType {
	return v1.CodingIdentityProvider
}

func (g *codingProviderFactory) Create(opts options.DynamicOptions) (identityprovider.OAuthProvider, error) {
	var coding coding
	if err := mapstructure.Decode(opts, &coding); err != nil {
		return nil, err
	}

	if coding.Endpoint.AuthURL == "" {
		coding.Endpoint.AuthURL = fmt.Sprintf(authURL, coding.Team)
	}
	if coding.Endpoint.TokenURL == "" {
		coding.Endpoint.TokenURL = fmt.Sprintf(tokenURL, coding.Team)
	}
	if coding.Endpoint.UserInfoURL == "" {
		coding.Endpoint.UserInfoURL = fmt.Sprintf(userInfoURL, coding.Team)
	}
	// fixed options
	opts["endpoint"] = options.DynamicOptions{
		"authURL":     coding.Endpoint.AuthURL,
		"tokenURL":    coding.Endpoint.TokenURL,
		"userInfoURL": coding.Endpoint.UserInfoURL,
	}
	coding.Config = &oauth2.Config{
		ClientID:     coding.ClientID,
		ClientSecret: coding.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  coding.Endpoint.AuthURL,
			TokenURL: coding.Endpoint.TokenURL,
		},
		RedirectURL: coding.RedirectURL,
		Scopes:      coding.Scopes,
	}
	return &coding, nil
}

func (c codingIdentity) GetUserID() string {
	return strconv.Itoa(c.ID)
}

func (c codingIdentity) GetUsername() string {
	return c.Name
}

func (c codingIdentity) GetEmail() string {
	return ""
}

func (c codingIdentity) GetAvatar() string {
	return c.Avatar
}

func (c *coding) IdentityExchangeCallback(req *http.Request) (identityprovider.Identity, error) {
	// OAuth2 callback, see also https://tools.ietf.org/html/rfc6749#section-4.1.2
	code := req.URL.Query().Get("code")
	ctx := req.Context()
	if c.InsecureSkipVerify {
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		ctx = context.WithValue(ctx, oauth2.HTTPClient, client)
	}
	token, err := c.Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	resp, err := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token)).Get(c.Endpoint.UserInfoURL)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var codingIdentity codingIdentity
	err = json.Unmarshal(data, &codingIdentity)
	if err != nil {
		return nil, err
	}

	return codingIdentity, nil
}
