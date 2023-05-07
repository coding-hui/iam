// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"fmt"

	"github.com/coding-hui/iam/internal/apiserver/config"
)

// needInitData register the service that need to init data.
var needInitData []DataInit

// InitServiceBean init all service instance.
func InitServiceBean(c config.Config) []interface{} {
	authenticationService := NewAuthenticationService(c)
	userService := NewUserService()
	resourceService := NewResourceService()
	roleService := NewRoleService()
	policyService := NewPolicyService()

	needInitData = []DataInit{userService, roleService, resourceService}

	return []interface{}{userService, authenticationService, resourceService, roleService, policyService}
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
