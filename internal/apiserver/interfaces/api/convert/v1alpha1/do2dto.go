// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1alpha1

import (
	iamv1alpha1 "github.com/coding-hui/api/iam/v1alpha1"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"
)

// ConvertUserModelToBase assemble the User model to DTO
func ConvertUserModelToBase(user *iamv1alpha1.User) *v1alpha1.UserBase {
	return &v1alpha1.UserBase{
		ObjectMeta:    user.ObjectMeta,
		TenantId:      user.TenantId,
		Status:        user.Status,
		Alias:         user.Alias,
		Email:         user.Email,
		LastLoginTime: user.LastLoginTime,
		Disabled:      user.Disabled,
	}
}
