// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gitee

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider"
	"github.com/coding-hui/iam/internal/pkg/options"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

const (
	userInfoURL = "https://gitee.com/api/v5/user"
	authURL     = "https://gitee.com/oauth/authorize"
	tokenURL    = "https://gitee.com/oauth/token"
)

func init() {
	identityprovider.RegisterOAuthProvider(&giteeProviderFactory{})
}

type gitee struct {
	// ClientID is the application's ID.
	ClientID string `json:"clientID" yaml:"clientID"`

	// ClientSecret is the application's secret.
	ClientSecret string `json:"clientSecret" yaml:"clientSecret"`

	// Endpoint contains the resource server's token endpoint
	// URLs. These are constants specific to each server and are
	// often available via site-specific packages, such as
	// google.Endpoint or gitee.endpoint.
	Endpoint endpoint `json:"endpoint" yaml:"endpoint"`

	// RedirectURL is the URL to redirect users going through
	// the OAuth flow, after the resource owner's URLs.
	RedirectURL string `json:"redirectURL" yaml:"redirectURL"`

	// Used to turn off TLS certificate checks
	InsecureSkipVerify bool `json:"insecureSkipVerify" yaml:"insecureSkipVerify"`

	// Scope specifies optional requested permissions.
	Scopes []string `json:"scopes" yaml:"scopes"`

	Config *oauth2.Config `json:"-" yaml:"-"`
}

// endpoint represents an OAuth 2.0 provider's authorization and token
// endpoint URLs.
type endpoint struct {
	AuthURL     string `json:"authURL"     yaml:"authURL"`
	TokenURL    string `json:"tokenURL"    yaml:"tokenURL"`
	UserInfoURL string `json:"userInfoURL" yaml:"userInfoURL"`
}

type giteeIdentity struct {
	AvatarUrl         string `json:"avatar_url"`
	Bio               string `json:"bio"`
	Blog              string `json:"blog"`
	CreatedAt         string `json:"created_at"`
	Email             string `json:"email"`
	EventsUrl         string `json:"events_url"`
	Followers         int    `json:"followers"`
	FollowersUrl      string `json:"followers_url"`
	Following         int    `json:"following"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	HtmlUrl           string `json:"html_url"`
	Id                int    `json:"id"`
	Login             string `json:"login"`
	MemberRole        string `json:"member_role"`
	Name              string `json:"name"`
	OrganizationsUrl  string `json:"organizations_url"`
	PublicGists       int    `json:"public_gists"`
	PublicRepos       int    `json:"public_repos"`
	ReceivedEventsUrl string `json:"received_events_url"`
	ReposUrl          string `json:"repos_url"`
	Stared            int    `json:"stared"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	Type              string `json:"type"`
	UpdatedAt         string `json:"updated_at"`
	Url               string `json:"url"`
	Watched           int    `json:"watched"`
	Weibo             string `json:"weibo"`
}

type giteeProviderFactory struct{}

func (g *giteeProviderFactory) Type() v1.IdentityProviderType {
	return v1.GiteeIdentityProvider
}

func (g *giteeProviderFactory) Create(opts options.DynamicOptions) (identityprovider.OAuthProvider, error) {
	var gitee gitee
	if err := opts.To(&gitee); err != nil {
		return nil, err
	}

	if gitee.Endpoint.AuthURL == "" {
		gitee.Endpoint.AuthURL = authURL
	}
	if gitee.Endpoint.TokenURL == "" {
		gitee.Endpoint.TokenURL = tokenURL
	}
	if gitee.Endpoint.UserInfoURL == "" {
		gitee.Endpoint.UserInfoURL = userInfoURL
	}
	// fixed options
	opts["endpoint"] = options.DynamicOptions{
		"authURL":     gitee.Endpoint.AuthURL,
		"tokenURL":    gitee.Endpoint.TokenURL,
		"userInfoURL": gitee.Endpoint.UserInfoURL,
	}
	gitee.Config = &oauth2.Config{
		ClientID:     gitee.ClientID,
		ClientSecret: gitee.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  gitee.Endpoint.AuthURL,
			TokenURL: gitee.Endpoint.TokenURL,
		},
		RedirectURL: gitee.RedirectURL,
		Scopes:      gitee.Scopes,
	}
	return &gitee, nil
}

func (g giteeIdentity) GetUserID() string {
	return g.Login
}

func (g giteeIdentity) GetUsername() string {
	return g.Name
}

func (g giteeIdentity) GetEmail() string {
	return g.Email
}

func (g giteeIdentity) GetAvatar() string {
	return g.AvatarUrl
}

func (g *gitee) IdentityExchangeCallback(req *http.Request) (identityprovider.Identity, error) {
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

	var giteeIdentity giteeIdentity
	err = json.Unmarshal(data, &giteeIdentity)
	if err != nil {
		return nil, err
	}

	return giteeIdentity, nil
}
