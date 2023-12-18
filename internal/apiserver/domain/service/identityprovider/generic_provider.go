// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package identityprovider

import (
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/options"
)

type GenericProvider interface {
	// Authenticate from remote server
	Authenticate(req v1.AuthenticateRequest) (Identity, error)
}

type GenericProviderFactory interface {
	// Type unique type of the provider
	Type() v1.IdentityProviderType
	// Create Apply the dynamic options from iam-apiserver-config
	Create(options options.DynamicOptions) (GenericProvider, error)
}
