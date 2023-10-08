// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package identityprovider

import (
	"fmt"

	"github.com/coding-hui/common/errors"

	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/internal/pkg/options"
	"github.com/coding-hui/iam/pkg/log"
)

var (
	oauthProviderFactories   = make(map[string]OAuthProviderFactory)
	genericProviderFactories = make(map[string]GenericProviderFactory)
	oauthProviders           = make(map[string]OAuthProvider)
	genericProviders         = make(map[string]GenericProvider)
)

type UserInfo interface {
}

// Identity represents the account mapped to iam
type Identity interface {
	// GetUserID required
	// Identifier for the End-User at the Issuer.
	GetUserID() string
	// GetUsername optional
	// The username which the End-User wishes to be referred to iam.
	GetUsername() string
	// GetEmail optional
	GetEmail() string
	// GetAvatar optional
	GetAvatar() string
}

// SetupWithOptions will verify the configuration and initialize the identityProviders
func SetupWithOptions(options []options.IdentityProviderOptions) error {
	// Clear all providers when reloading configuration
	oauthProviders = make(map[string]OAuthProvider)
	genericProviders = make(map[string]GenericProvider)

	for _, o := range options {
		if oauthProviders[o.Name] != nil || genericProviders[o.Name] != nil {
			err := fmt.Errorf("duplicate identity provider found: %s, name must be unique", o.Name)
			log.Error(err.Error())
			return err
		}
		if genericProviderFactories[o.Type] == nil && oauthProviderFactories[o.Type] == nil {
			err := fmt.Errorf("identity provider %s with type %s is not supported", o.Name, o.Type)
			log.Error(err.Error())
			return err
		}
		if factory, ok := oauthProviderFactories[o.Type]; ok {
			if provider, err := factory.Create(o.Provider); err != nil {
				// don’t return errors, decoupling external dependencies
				log.Errorf("failed to create identity provider %s: %s", o.Name, err)
			} else {
				oauthProviders[o.Name] = provider
				log.Infof("create identity provider %s successfully", o.Name)
			}
		}
		if factory, ok := genericProviderFactories[o.Type]; ok {
			if provider, err := factory.Create(o.Provider); err != nil {
				log.Errorf("failed to create identity provider %s: %s", o.Name, err)
			} else {
				genericProviders[o.Name] = provider
				log.Infof("create identity provider %s successfully", o.Name)
			}
		}
	}
	return nil
}

// GetGenericProvider returns GenericProvider with given name
func GetGenericProvider(providerName string) (GenericProvider, error) {
	if provider, ok := genericProviders[providerName]; ok {
		return provider, nil
	}
	return nil, errors.WithCode(code.ErrIdentityProviderNotFound, "identity provider [%s] not found", providerName)
}

// GetOAuthProvider returns OAuthProvider with given name
func GetOAuthProvider(providerName string) (OAuthProvider, error) {
	if provider, ok := oauthProviders[providerName]; ok {
		return provider, nil
	}
	return nil, errors.WithCode(code.ErrIdentityProviderNotFound, "identity provider [%s] not found", providerName)
}

// RegisterOAuthProvider register OAuthProviderFactory with the specified type
func RegisterOAuthProvider(factory OAuthProviderFactory) {
	oauthProviderFactories[factory.Type()] = factory
}

// RegisterGenericProvider registers GenericProviderFactory with the specified type
func RegisterGenericProvider(factory GenericProviderFactory) {
	genericProviderFactories[factory.Type()] = factory
}
