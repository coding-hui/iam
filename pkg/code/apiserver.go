// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package code

//go:generate codegen -type=int

// iam-apiserver: user errors.
const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 110001

	// ErrUserAlreadyExist - 400: User already exists.
	ErrUserAlreadyExist

	// ErrUserNameIsEmpty - 400: Username cannot be empty.
	ErrUserNameIsEmpty

	// ErrDeleteOneself - 400: You cannot delete your own account.
	ErrDeleteOneself

	// ErrUserAlreadyDisabled - 400: The user is already disabled.
	ErrUserAlreadyDisabled

	// ErrUserAlreadyEnabled - 400: The user is already enabled.
	ErrUserAlreadyEnabled

	// ErrUserHasDisabled - 401: Your account has been disabled. Please contact the administrator.
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

	// ErrPolicyAlreadyExist - 400: Policy already exists.
	ErrPolicyAlreadyExist

	// ErrPolicyNameIsEmpty - 400: Policy name cannot be empty.
	ErrPolicyNameIsEmpty
)

// iam-apiserver: resource errors.
const (
	// ErrResourceNotFound - 404: Resource not found.
	ErrResourceNotFound int = iota + 110301

	// ErrResourceAlreadyExist - 400: Resource already exists.
	ErrResourceAlreadyExist

	// ErrResourceNameIsEmpty - 400: Resource name cannot be empty.
	ErrResourceNameIsEmpty

	// ErrResourceInstanceIdIsEmpty - 400: Resource instance ID cannot be empty.
	ErrResourceInstanceIdIsEmpty

	// ErrResourceHasAssignedPolicy - 400: The resource has been assigned permission policies and cannot be deleted.
	ErrResourceHasAssignedPolicy
)

// iam-apiserver: role errors.
const (
	// ErrRoleNotFound - 404: Role not found.
	ErrRoleNotFound int = iota + 110401

	// ErrRoleAlreadyExist - 400: Role already exists.
	ErrRoleAlreadyExist

	// ErrRoleNameIsEmpty - 400: Role name cannot be empty.
	ErrRoleNameIsEmpty

	// ErrRoleHasAssignedUser - 400: The role has been assigned to users and cannot be deleted.
	ErrRoleHasAssignedUser

	// ErrAssignRoleFailed - 400: Failed to assign role. Please check the role status or contact the administrator.
	ErrAssignRoleFailed

	// ErrUnsupportedAssignTarget - 400: Unsupported assignment target. Only users or departments are supported.
	ErrUnsupportedAssignTarget

	// ErrRevokeRoleFailed - 400: Failed to revoke role. Please check the role status or contact the administrator.
	ErrRevokeRoleFailed

	// ErrUnsupportedRevokeTarget - 400: Unsupported revoke target. Only users or departments are supported.
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

// iam-apiserver: external account binding errors.
const (
	// ErrBindExternalAccount - 400: Failed to bind external account.
	ErrBindExternalAccount int = iota + 110701

	// ErrUnbindExternalAccount - 400: Failed to unbind external account.
	ErrUnbindExternalAccount

	// ErrExternalAccountAlreadyBound - 400: External account already bound to another user.
	ErrExternalAccountAlreadyBound
)

// iam-apiserver: api key errors.
const (
	// ErrApiKeyNotFound - 404: API Key not found.
	ErrApiKeyNotFound int = iota + 110801

	// ErrApiKeyAlreadyExist - 400: API Key already exists.
	ErrApiKeyAlreadyExist

	// ErrApiKeyInvalid - 401: Invalid API Key. Please check the key format and ensure it starts with 'sk-'.
	ErrApiKeyInvalid

	// ErrApiKeyInactive - 401: API Key is not active. Please enable the API Key to use it.
	ErrApiKeyInactive

	// ErrApiKeyExpired - 401: API Key has expired. Please regenerate a new API Key.
	ErrApiKeyExpired

	// ErrApiKeyAlreadyEnabled - 400: API Key is already enabled.
	ErrApiKeyAlreadyEnabled

	// ErrApiKeyAlreadyDisabled - 400: API Key is already disabled.
	ErrApiKeyAlreadyDisabled

	// ErrApiKeyGenerationFailed - 500: Failed to generate unique API Key after multiple attempts. Please try again.
	ErrApiKeyGenerationFailed
)

// iam-apiserver: device authorization errors.
const (
	// ErrDeviceCodeInvalid - 400: Invalid device code.
	ErrDeviceCodeInvalid int = iota + 110901

	// ErrAuthorizationPending - 400: Authorization pending.
	ErrAuthorizationPending

	// ErrAuthorizationDenied - 400: Authorization denied by user.
	ErrAuthorizationDenied

	// ErrDeviceCodeExpired - 400: Device code has expired.
	ErrDeviceCodeExpired

	// ErrClientNotFound - 400: OAuth client not found.
	ErrClientNotFound
)

// iam-apiserver: email template errors.
const (
	// ErrCannotDeleteSystemEmailTemplate - 403: Cannot delete system email template.
	ErrCannotDeleteSystemEmailTemplate int = iota + 111001

	// ErrCannotDeleteDefaultEmailTemplate - 403: Cannot delete default email template.
	ErrCannotDeleteDefaultEmailTemplate

	// ErrCannotDeleteSystemEmailTemplateCategory - 403: Cannot delete system email template category.
	ErrCannotDeleteSystemEmailTemplateCategory
)
