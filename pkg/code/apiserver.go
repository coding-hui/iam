// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package code

//go:generate codegen -type=int

// iam-apiserver: user errors.
const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 110001

	// ErrUserAlreadyExist - 400: User already exist.
	ErrUserAlreadyExist

	// ErrUserNameIsEmpty - 400: Username is empty.
	ErrUserNameIsEmpty

	// ErrDeleteOneself - 400: Unable to delete yourself.
	ErrDeleteOneself

	// ErrUserAlreadyDisabled - 400: The user is already disabled.
	ErrUserAlreadyDisabled

	// ErrUserAlreadyEnabled - 400: The user is already enabled.
	ErrUserAlreadyEnabled

	// ErrUserHasDisabled - 401: The account has been disabled.
	ErrUserHasDisabled
)

// iam-apiserver: secret errors.
const (
	// ErrReachMaxCount - 400: Secret reach the max count.
	ErrReachMaxCount int = iota + 110101

	// ErrSecretNotFound - 404: Secret not found.
	ErrSecretNotFound
)

// iam-apiserver: policy errors.
const (
	// ErrPolicyNotFound - 404: Policy not found.
	ErrPolicyNotFound int = iota + 110201

	// ErrPolicyAlreadyExist - 400: Policy already exist.
	ErrPolicyAlreadyExist

	// ErrPolicyNameIsEmpty - 400: Policy name is empty.
	ErrPolicyNameIsEmpty
)

// iam-apiserver: resource errors.
const (
	// ErrResourceNotFound - 404: Resource not found.
	ErrResourceNotFound int = iota + 110301

	// ErrResourceAlreadyExist - 400: Resource already exist.
	ErrResourceAlreadyExist

	// ErrResourceNameIsEmpty - 400: Resource name is empty.
	ErrResourceNameIsEmpty

	// ErrResourceInstanceIdIsEmpty - 400: Resource instanceId is empty.
	ErrResourceInstanceIdIsEmpty

	// ErrResourceHasAssignedPolicy - 400: The resource has been assigned permission policies.
	ErrResourceHasAssignedPolicy
)

// iam-apiserver: role errors.
const (
	// ErrRoleNotFound - 404: Role not found.
	ErrRoleNotFound int = iota + 110401

	// ErrRoleAlreadyExist - 400: Role already exist.
	ErrRoleAlreadyExist

	// ErrRoleNameIsEmpty - 400: Role name is empty.
	ErrRoleNameIsEmpty

	// ErrRoleHasAssignedUser - 400: The role has been assigned to a user.
	ErrRoleHasAssignedUser

	// ErrAssignRoleFailed - 400: User role assignment fails. Please check the role status or contact the administrator.
	ErrAssignRoleFailed

	// ErrUnsupportedAssignTarget - 400: The assignment target is not supported. Only user or department are supported.
	ErrUnsupportedAssignTarget

	// ErrRevokeRoleFailed - 400: User role revoke fails. Please check the role status or contact the administrator.
	ErrRevokeRoleFailed

	// ErrUnsupportedRevokeTarget - 400: The revoke target is not supported. Only user or department are supported.
	ErrUnsupportedRevokeTarget
)

// iam-apiserver: organization errors.
const (
	// ErrOrgNotFound - 404: Organization not found.
	ErrOrgNotFound int = iota + 110501

	// ErrOrgAlreadyExist - 400: Organization already exist.
	ErrOrgAlreadyExist

	// ErrOrgAlreadyDisabled - 400: The organization is already disabled.
	ErrOrgAlreadyDisabled

	// ErrOrgAlreadyEnabled - 400: The organization is already enabled.
	ErrOrgAlreadyEnabled

	// ErrOrgHasDisabled - 401: The organization has been disabled.
	ErrOrgHasDisabled

	// ErrCannotDeleteBuiltInOrg - 400: Built-in organizations cannot be deleted.
	ErrCannotDeleteBuiltInOrg

	// ErrCannotDisableBuiltInOrg - 400: Built-in organizations cannot be disabled.
	ErrCannotDisableBuiltInOrg

	// ErrMaxDepartmentsReached  - 400: The number of departments has reached its limit.
	ErrMaxDepartmentsReached
)

// iam-apiserver: organization department errors.
const (
	// ErrMemberAlreadyInDepartment - 400: Member is already in department.
	ErrMemberAlreadyInDepartment int = iota + 110601

	// ErrSubDepartmentsExist - 400: Sub departments exist and cannot be deleted.
	ErrSubDepartmentsExist
)
