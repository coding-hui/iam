// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package identityprovider

import (
	"net/http"

	"github.com/coding-hui/iam/internal/pkg/options"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

type OAuthProvider interface {
	// IdentityExchangeCallback handle oauth callback, exchange identity from remote server
	IdentityExchangeCallback(req *http.Request) (Identity, error)
}

type OAuthProviderFactory interface {
	// Type unique type of the provider
	Type() v1.ProviderType
	// Create Apply the dynamic options
	Create(options options.DynamicOptions) (OAuthProvider, error)
}
