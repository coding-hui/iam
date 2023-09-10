// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

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

// Validate validates that a organization object is valid.
func (o *CreateOrganizationRequest) Validate() field.ErrorList {
	val := validation.NewValidator(o)
	allErrs := val.Validate()

	return allErrs
}

// ValidateUpdate validates that a organization object is valid when update.
func (o *UpdateOrganizationRequest) ValidateUpdate() field.ErrorList {
	val := validation.NewValidator(o)
	allErrs := val.Validate()

	return allErrs
}

// Validate validates that a role object is valid.
func (r *CreateRoleRequest) Validate() field.ErrorList {
	val := validation.NewValidator(r)
	allErrs := val.Validate()

	return allErrs
}

// ValidateUpdate validates that a role object is valid when update.
func (r *UpdateRoleRequest) ValidateUpdate() field.ErrorList {
	val := validation.NewValidator(r)
	allErrs := val.Validate()

	return allErrs
}

// Validate validates that a resource object is valid.
func (r *CreateResourceRequest) Validate() field.ErrorList {
	val := validation.NewValidator(r)
	allErrs := val.Validate()

	return allErrs
}

// ValidateUpdate validates that a user resource is valid when update.
func (r *UpdateResourceRequest) ValidateUpdate() field.ErrorList {
	val := validation.NewValidator(r)
	allErrs := val.Validate()

	return allErrs
}

// Validate validates that a policy object is valid.
func (p *CreatePolicyRequest) Validate() field.ErrorList {
	val := validation.NewValidator(p)
	allErrs := val.Validate()

	return allErrs
}

// ValidateUpdate validates that a policy object is valid when update.
func (p *UpdatePolicyRequest) ValidateUpdate() field.ErrorList {
	val := validation.NewValidator(p)
	allErrs := val.Validate()

	return allErrs
}
