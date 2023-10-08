// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

type AuthenticationOptions struct {
	// Token verification maximum time difference, default to 10s.
	// You should consider allowing a clock skew when checking the time-based values.
	// This should be values of a few seconds, and we don’t recommend using more than 30 seconds for this purpose,
	// as this would rather indicate problems with the server, rather than a common clock skew.
	MaximumClockSkew time.Duration `json:"maximumClockSkew" mapstructure:"maximumClockSkew"`
	// secret to sign jwt token
	JwtSecret string `json:"-"                mapstructure:"jwtSecret"`
	// OAuthOptions defines options needed for integrated oauth plugins
	OAuthOptions *OAuthOptions `json:"oauth"            mapstructure:"oauth"`
	// ClientCA is the certificate bundle for all the signers that you'll recognize for incoming client certificates
	ClientCA string `json:"client-ca-file"   mapstructure:"client-ca-file"`
}

func NewAuthenticationOptions() *AuthenticationOptions {
	return &AuthenticationOptions{
		MaximumClockSkew: 10 * time.Second,
		OAuthOptions:     NewOauthOptions(),
		JwtSecret:        "",
		ClientCA:         "",
	}
}

func (o *AuthenticationOptions) Validate() []error {
	var errs []error
	if len(o.JwtSecret) == 0 {
		errs = append(errs, fmt.Errorf("JWT secret MUST not be empty"))
	}
	return errs
}

func (o *AuthenticationOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ClientCA, "client-ca-file", o.ClientCA, ""+
		"If set, any request presenting a client certificate signed by one of "+
		"the authorities in the client-ca-file is authenticated with an identity "+
		"corresponding to the CommonName of the client certificate.")
	fs.StringVar(&o.JwtSecret, "jwt-secret", o.JwtSecret, "Secret to sign jwt token, must not be empty.")
	fs.DurationVar(
		&o.OAuthOptions.AccessTokenMaxAge,
		"access-token-max-age",
		o.OAuthOptions.AccessTokenMaxAge,
		"access-token-max-age control the lifetime of access tokens, 0 means no expiration.",
	)
	fs.DurationVar(
		&o.MaximumClockSkew,
		"maximum-clock-skew",
		o.MaximumClockSkew,
		"The maximum time difference between the system clocks of the iam-apiserver that issued a JWT and the iam-apiserver that verified the JWT.",
	)
}
