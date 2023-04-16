// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1alpha1

import (
	"github.com/coding-hui/common/validation"
	"github.com/coding-hui/common/validation/field"
)

// Validate validates that a user object is valid.
func (u *CreateUserRequest) Validate() field.ErrorList {
	val := validation.NewValidator(u)
	allErrs := val.Validate()

	if err := validation.IsValidPassword(u.Password); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("password"), err.Error(), ""))
	}

	return allErrs
}

// ValidateUpdate validates that a user object is valid when update.
// Like User.Validate but not validate password.
func (u *UpdateUserRequest) ValidateUpdate() field.ErrorList {
	val := validation.NewValidator(u)
	allErrs := val.Validate()

	return allErrs
}
