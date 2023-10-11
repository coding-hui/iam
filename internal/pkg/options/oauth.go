// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"net/url"
	"time"

	"github.com/spf13/pflag"

	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"

	"github.com/coding-hui/common/errors"
	"github.com/coding-hui/common/util/sliceutil"
)

const (
	DefaultIssuer string = "iam-apiserver"
)

var (
	ErrorClientNotFound        = errors.New("the OAuth client was not found")
	ErrorProviderNotFound      = errors.New("the identity provider was not found")
	ErrorRedirectURLNotAllowed = errors.New("redirect URL is not allowed")
)

type OAuthOptions struct {
	// An Issuer Identifier is a case-sensitive URL using the https scheme that contains scheme,
	// host, and optionally, port number and path components and no query or fragment components.
	Issuer string `json:"issuer,omitempty" mapstructure:"issuer,omitempty"`

	// RSA private key file used to sign the id token
	SignKey string `json:"signKey,omitempty" mapstructure:"signKey,omitempty"`

	// Raw RSA private key. Base64 encoded PEM file
	SignKeyData string `json:"-,omitempty" mapstructure:"signKeyData,omitempty"`

	// Register identity providers.
	IdentityProviders []IdentityProviderOptions `json:"identityProviders,omitempty" mapstructure:"identityProviders,omitempty"`

	// Register additional OAuth clients.
	Clients []Client `json:"clients,omitempty" mapstructure:"clients,omitempty"`

	// AccessTokenMaxAgeSeconds  control the lifetime of access tokens. The default lifetime is 24 hours.
	// 0 means no expiration.
	AccessTokenMaxAge time.Duration `json:"accessTokenMaxAge" mapstructure:"accessTokenMaxAge"`

	// Inactivity timeout for tokens
	// The value represents the maximum amount of time that can occur between
	// consecutive uses of the token. Tokens become invalid if they are not
	// used within this temporal window. The user will need to acquire a new
	// token to regain access once a token times out.
	// This value needs to be set only if the default set in configuration is
	// not appropriate for this client. Valid values are:
	// - 0: Tokens for this client never time out
	// - X: Tokens time out if there is no activity
	// The current minimum allowed value for X is 5 minutes
	AccessTokenInactivityTimeout time.Duration `json:"accessTokenInactivityTimeout" mapstructure:"accessTokenInactivityTimeout"`
}

type IdentityProviderOptions struct {
	// The provider name.
	Name string `json:"name" mapstructure:"name"`

	// Defines how new identities are mapped to users when they login. Allowed values are:
	//  - auto:   The default value.The user will automatically create and mapping when login successful.
	//            Fails if a user with that user name is already mapped to another identity.
	//  - lookup: Looks up an existing identity, user identity mapping, and user, but does not automatically
	//            provision users or identities. Using this method requires you to manually provision users.
	//  - mixed:  A user entity can be mapped with multiple identifyProvider.
	MappingMethod v1.MappingMethod `json:"mappingMethod" mapstructure:"mappingMethod"`

	// DisableLoginConfirmation means that when the user login successfully,
	// reconfirm the account information is not required.
	// Username from IDP must math [a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*
	DisableLoginConfirmation bool `json:"disableLoginConfirmation" mapstructure:"disableLoginConfirmation"`

	// The type of identify provider
	// OpenIDIdentityProvider LDAPIdentityProvider GitHubIdentityProvider
	Type string `json:"type" mapstructure:"type"`

	// The category of identify provider
	// OAuth Email Storage
	Category string `json:"category" mapstructure:"category"`

	// The options of identify provider
	Provider DynamicOptions `json:"provider" mapstructure:"provider"`
}

type Client struct {
	// The name of the OAuth client is used as the client_id parameter when making requests to <master>/oauth/authorize
	// and <master>/oauth/token.
	Name string `json:"name,omitempty" mapstructure:"name,omitempty"`

	// Secret is the unique secret associated with a client
	Secret string `json:"-" yaml:"secret,omitempty"`

	// RespondWithChallenges indicates whether the client wants authentication needed responses made
	// in the form of challenges instead of redirects
	RespondWithChallenges bool `json:"respondWithChallenges,omitempty" mapstructure:"respondWithChallenges,omitempty"`

	// RedirectURIs is the valid redirection URIs associated with a client
	RedirectURIs []string `json:"redirectURIs,omitempty" mapstructure:"redirectURIs,omitempty"`

	// GrantMethod determines how to handle grants for this client. If no method is provided, the
	// cluster default grant handling method will be used. Valid grant handling methods are:
	//  - auto:   always approves grant requests, useful for trusted clients
	//  - prompt: prompts the end user for approval of grant requests, useful for third-party clients
	//  - deny:   always denies grant requests, useful for black-listed clients
	GrantMethod v1.GrantHandlerType `json:"grantMethod,omitempty" mapstructure:"grantMethod,omitempty"`

	// ScopeRestrictions describes which scopes this client can request.  Each requested scope
	// is checked against each restriction.  If any restriction matches, then the scope is allowed.
	// If no restriction matches, then the scope is denied.
	ScopeRestrictions []string `json:"scopeRestrictions,omitempty" mapstructure:"scopeRestrictions,omitempty"`

	// AccessTokenMaxAge overrides the default access token max age for tokens granted to this client.
	AccessTokenMaxAge *time.Duration `json:"accessTokenMaxAge,omitempty" mapstructure:"accessTokenMaxAge,omitempty"`

	// AccessTokenInactivityTimeout overrides the default token
	// inactivity timeout for tokens granted to this client.
	AccessTokenInactivityTimeout *time.Duration `json:"accessTokenInactivityTimeout,omitempty" mapstructure:"accessTokenInactivityTimeout,omitempty"`
}

var (
	// AllowAllRedirectURI Allow any redirect URI if the redirectURI is defined in request
	AllowAllRedirectURI = "*"
)

func (o *OAuthOptions) OAuthClient(name string) (Client, error) {
	for _, found := range o.Clients {
		if found.Name == name {
			return found, nil
		}
	}
	return Client{}, ErrorClientNotFound
}

func (o *OAuthOptions) IdentityProviderOptions(name string) (*IdentityProviderOptions, error) {
	for _, found := range o.IdentityProviders {
		if found.Name == name {
			return &found, nil
		}
	}
	return nil, ErrorProviderNotFound
}

func (c Client) anyRedirectAbleURI() []string {
	uris := make([]string, 0)
	for _, uri := range c.RedirectURIs {
		_, err := url.Parse(uri)
		if err == nil {
			uris = append(uris, uri)
		}
	}
	return uris
}

func (c Client) ResolveRedirectURL(expectURL string) (*url.URL, error) {
	// RedirectURIs is empty
	if len(c.RedirectURIs) == 0 {
		return nil, ErrorRedirectURLNotAllowed
	}
	allowAllRedirectURI := sliceutil.HasString(c.RedirectURIs, AllowAllRedirectURI)
	redirectAbleURIs := c.anyRedirectAbleURI()

	if expectURL == "" {
		// Need to specify at least one RedirectURI
		if len(redirectAbleURIs) > 0 {
			return url.Parse(redirectAbleURIs[0])
		} else {
			return nil, ErrorRedirectURLNotAllowed
		}
	}
	if allowAllRedirectURI || sliceutil.HasString(redirectAbleURIs, expectURL) {
		return url.Parse(expectURL)
	}

	return nil, ErrorRedirectURLNotAllowed
}

// NewOauthOptions creates a OAuthOptions object with default parameters.
func NewOauthOptions() *OAuthOptions {
	return &OAuthOptions{
		Issuer:                       DefaultIssuer,
		IdentityProviders:            make([]IdentityProviderOptions, 0),
		Clients:                      make([]Client, 0),
		AccessTokenMaxAge:            time.Hour * 2,
		AccessTokenInactivityTimeout: time.Hour * 2,
	}
}

// Validate is used to parse and validate the parameters entered by the user at
// the command line when the program starts.
func (o *OAuthOptions) Validate() []error {
	var errs []error
	return errs
}

// AddFlags adds flags related to features for a specific api server to the
// specified FlagSet.
func (o *OAuthOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}
}
