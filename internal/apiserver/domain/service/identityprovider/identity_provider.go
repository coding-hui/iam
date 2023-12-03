// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package identityprovider

import (
	"fmt"
	"sync"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/pkg/options"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/log"
)

var (
	oauthProviderFactories   = make(map[v1.IdentityProviderType]OAuthProviderFactory)
	genericProviderFactories = make(map[v1.IdentityProviderType]GenericProviderFactory)
	oauthProviders           = make(map[string]OAuthProvider)
	genericProviders         = make(map[string]GenericProvider)

	lock = sync.Mutex{}
)

// Identity represents the account mapped to iam.
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

// GetGenericProvider returns GenericProvider with given name.
func GetGenericProvider(idp *model.IdentityProvider) (GenericProvider, error) {
	if provider, ok := genericProviders[idp.Name]; ok {
		return provider, nil
	}
	if genericProviderFactories[idp.Type] == nil {
		err := fmt.Errorf("identity provider %s with type %s is not supported", idp.Name, idp.Type)
		return nil, err
	}
	return SetGenericProvider(idp)
}

// SetGenericProvider update GenericProvider based on the name.
func SetGenericProvider(idp *model.IdentityProvider) (GenericProvider, error) {
	if idp.Extend == nil {
		err := fmt.Errorf("identity provider %s config is empty", idp.Name)
		return nil, err
	}
	if factory, ok := genericProviderFactories[idp.Type]; ok {
		if provider, err := factory.Create(options.DynamicOptions(idp.Extend)); err != nil {
			log.Errorf("failed to create identity provider %s: %s", idp.Name, err)
		} else {
			lock.TryLock()
			defer lock.Unlock()
			genericProviders[idp.Name] = provider
			log.Infof("create identity provider %s successfully", idp.Name)
			return provider, nil
		}
	}
	return nil, fmt.Errorf("identity provider %s with type %s is not supported", idp.Name, idp.Type)
}

// GetOAuthProvider returns OAuthProvider with given name.
func GetOAuthProvider(idp *model.IdentityProvider) (OAuthProvider, error) {
	if provider, ok := oauthProviders[idp.Name]; ok {
		return provider, nil
	}
	if oauthProviderFactories[idp.Type] == nil {
		err := fmt.Errorf("identity provider %s with type %s is not supported", idp.Name, idp.Type)
		return nil, err
	}
	return SetOAuthProvider(idp)
}

// SetOAuthProvider update OAuthProvider based on the name.
func SetOAuthProvider(idp *model.IdentityProvider) (OAuthProvider, error) {
	if idp.Extend == nil {
		err := fmt.Errorf("identity provider %s config is empty", idp.Name)
		return nil, err
	}
	if factory, ok := oauthProviderFactories[idp.Type]; ok {
		if provider, err := factory.Create(options.DynamicOptions(idp.Extend)); err != nil {
			log.Errorf("failed to create identity provider %s: %s", idp.Name, err)
		} else {
			lock.TryLock()
			defer lock.Unlock()
			oauthProviders[idp.Name] = provider
			log.Infof("create identity provider %s successfully", idp.Name)
			return provider, nil
		}
	}
	return nil, fmt.Errorf("identity provider %s with type %s is not supported", idp.Name, idp.Type)
}

// RegisterOAuthProvider register OAuthProviderFactory with the specified type.
func RegisterOAuthProvider(factory OAuthProviderFactory) {
	oauthProviderFactories[factory.Type()] = factory
}

// RegisterGenericProvider registers GenericProviderFactory with the specified type.
func RegisterGenericProvider(factory GenericProviderFactory) {
	genericProviderFactories[factory.Type()] = factory
}
