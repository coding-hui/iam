// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package identityprovider

import (
	"context"
	"fmt"
	"sync"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/internal/pkg/options"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
)

var (
	oauthProviderFactories   = make(map[v1.ProviderType]OAuthProviderFactory)
	genericProviderFactories = make(map[v1.ProviderType]GenericProviderFactory)
	oauthProviders           = make(map[string]OAuthProvider)
	genericProviders         = make(map[string]GenericProvider)

	lock = sync.Mutex{}
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
func GetGenericProvider(idp *model.Provider) (GenericProvider, error) {
	if provider, ok := genericProviders[idp.Name]; ok {
		return provider, nil
	}
	providerRepo := repository.Client().ProviderRepository()
	providerInfo, err := providerRepo.GetByName(context.Background(), idp.Name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.WithCode(code.ErrIdentityProviderNotFound, "identity provider [%s] not found", idp.Name)
	}
	if genericProviderFactories[providerInfo.Type] == nil {
		err := fmt.Errorf("identity provider %s with type %s is not supported", providerInfo.Name, providerInfo.Type)
		return nil, err
	}
	if factory, ok := genericProviderFactories[providerInfo.Type]; ok {
		if provider, err := factory.Create(options.DynamicOptions(idp.Extend)); err != nil {
			log.Errorf("failed to create identity provider %s: %s", providerInfo.Name, err)
		} else {
			lock.TryLock()
			defer lock.Unlock()
			genericProviders[providerInfo.Name] = provider
			log.Infof("create identity provider %s successfully", providerInfo.Name)
			return provider, nil
		}
	}
	return nil, errors.WithCode(code.ErrIdentityProviderNotFound, "identity provider [%s] not found", idp.Name)
}

// GetOAuthProvider returns OAuthProvider with given name
func GetOAuthProvider(idp *model.Provider) (OAuthProvider, error) {
	if provider, ok := oauthProviders[idp.Name]; ok {
		return provider, nil
	}
	if oauthProviderFactories[idp.Type] == nil {
		err := fmt.Errorf("identity provider %s with type %s is not supported", idp.Name, idp.Type)
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
	return nil, errors.WithCode(code.ErrIdentityProviderNotFound, "identity provider [%s] not found", idp.Name)
}

// RegisterOAuthProvider register OAuthProviderFactory with the specified type
func RegisterOAuthProvider(factory OAuthProviderFactory) {
	oauthProviderFactories[factory.Type()] = factory
}

// RegisterGenericProvider registers GenericProviderFactory with the specified type
func RegisterGenericProvider(factory GenericProviderFactory) {
	genericProviderFactories[factory.Type()] = factory
}
