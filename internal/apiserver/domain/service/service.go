// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"fmt"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/service/mail"
	"github.com/coding-hui/iam/internal/pkg/token"
)

// needInitData register the service that need to init data.
var needInitData []DataInit

// InitServiceBean init all service instance.
func InitServiceBean(c config.Config, issuer token.Issuer) []interface{} {
	tokenService := NewTokenService(c, issuer)
	authenticationService := NewAuthenticationService(c)
	userService := NewUserService()
	resourceService := NewResourceService()
	roleService := NewRoleService()
	policyService := NewPolicyService()
	organizationService := NewOrganizationService(c)
	providerService := NewIdentityProviderService(c)
	appService := NewApplicationService(c)
	mailService := mail.NewServiceWithConfig(c)
	mailTemplateService := mail.NewTemplateService()

	needInitData = []DataInit{
		organizationService, userService, roleService, resourceService,
		policyService, providerService, appService, authenticationService,
	}

	return []interface{}{
		organizationService, userService, authenticationService,
		resourceService, roleService, policyService,
		tokenService, providerService, appService,
		mailService, mailTemplateService,
	}
}

// DataInit the service set that needs init data.
type DataInit interface {
	Init(ctx context.Context) error
}

// InitData init data.
func InitData(ctx context.Context) error {
	for _, init := range needInitData {
		if err := init.Init(ctx); err != nil {
			return fmt.Errorf("database init failure %w", err)
		}
	}

	return nil
}
