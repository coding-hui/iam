// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	"encoding/json"
	"time"

	metav1 "github.com/coding-hui/common/meta/v1"
)

var (
	// CtxKeyUserRole request context key of user role.
	CtxKeyUserRole = "user_role"
	// CtxKeyUserInstanceId request context key of username.
	CtxKeyUserInstanceId = "user_instance_id"
	// CtxKeyUserType request context key of user type.
	CtxKeyUserType = "user_type"
	// CtxKeyRole request context key of role.
	CtxKeyRole = "role"
	// CtxKeyPolicy request context key of policy.
	CtxKeyPolicy = "policy"
	// CtxKeyResource request context key of resource.
	CtxKeyResource = "resource"

	// CtxKeyRoutes initialize context key of routes.
	CtxKeyRoutes = "routes"
	// CtxKeyApiPrefix initialize context key of api-prefix.
	CtxKeyApiPrefix = "api_prefix"
)

const (
	// UserTarget assign to user.
	UserTarget string = "user"
)

// UserType user type.
type UserType string

// These are the valid phases of a user role.
const (
	// PlatformAdmin platform admin.
	PlatformAdmin UserType = "platform"
	// TenantAdmin tenant admin.
	TenantAdmin UserType = "tenant"
	// Default default user.
	Default UserType = "default"
)

func (r UserType) String() string {
	return string(r)
}

// UserState user account state.
type UserState string

const (
	// UserActive means the user is available.
	UserActive UserState = "Active"
	// UserDisabled means the user is disabled.
	UserDisabled UserState = "Disabled"
	// UserAuthLimitExceeded means restrict user login.
	UserAuthLimitExceeded UserState = "AuthLimitExceeded"
)

// LoginType authenticate type.
type LoginType string

const (
	// BasicAuth basic auth type.
	BasicAuth LoginType = "Basic"
	// OAuth oauth2 auth type.
	OAuth LoginType = "OAuth"
	// Token jwt token auth type.
	Token LoginType = "Token"
)

// ResourceType resource type.
type ResourceType string

const (
	// API api resource
	API ResourceType = "API"
)

// CreateUserRequest create user request.
type CreateUserRequest struct {
	Name     string `json:"name"            validate:"required,name"`
	Password string `json:"password"        validate:"required"`
	Alias    string `json:"alias,omitempty"                          optional:"true"`
	Email    string `json:"email"                                    optional:"true"`
	Phone    string `json:"phone"`
	UserType string `json:"userType"`
}

// CreateUserResponse create user response.
type CreateUserResponse struct {
	UserBase
}

// UpdateUserRequest update user request.
type UpdateUserRequest struct {
	Alias    string `json:"alias,omitempty"    validate:"min=1,max=30"                 optional:"true"`
	Email    string `json:"email"              validate:"required,email,min=1,max=100"`
	Phone    string `json:"phone"`
	Password string `json:"password,omitempty" validate:"required"                     optional:"true"`
}

// UpdateUserResponse update user response.
type UpdateUserResponse struct {
	UserBase
}

// UserBase represents a user restful resource.
type UserBase struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`
	TenantId          uint64     `json:"tenantId,omitempty"`
	Status            int        `json:"status"`
	Alias             string     `json:"alias"`
	Password          string     `json:"password,omitempty"`
	Email             string     `json:"email"`
	Phone             string     `json:"phone"`
	UserType          string     `json:"userType"`
	Disabled          bool       `json:"disabled"`
	LastLoginTime     *time.Time `json:"lastLoginTime,omitempty"`
}

// DetailUserResponse is the response of user detail
type DetailUserResponse struct {
	UserBase
	Roles []RoleBase `json:"roles"`
}

// UserList is the whole list of all users which have been stored in stroage.
type UserList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*UserBase `json:"items"`
}

// TenantBase represents a tenant restful resource.
type TenantBase struct {
	// Standard object's metadata.
	metav1.ObjectMeta `       json:"metadata,omitempty"`
	// Owner tenant owner name.
	Owner string `json:"owner,omitempty"    gorm:"column:owner"       validate:"required"`
	// Disabled tenant state.
	Disabled bool `json:"disabled"           gorm:"column:disabled"`
	// Description tenant description.
	Description string `json:"description"        gorm:"column:description"`
}

// TenantList is the whole list of all tenants which have been stored in stroage.
type TenantList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*TenantBase `json:"items"`
}

// AssignRoleRequest assign role request.
type AssignRoleRequest struct {
	// InstanceID role instanceID.
	InstanceID string `json:"instanceId" validate:"required"`
	// Targets target resource instanceIds.
	Targets []string `json:"targets"`
}

// BatchAssignRoleRequest batch assign role request.
type BatchAssignRoleRequest struct {
	// InstanceIds role instanceIds.
	InstanceIds []string `json:"instanceIds" validate:"required"`
	// Targets target resource instanceIds.
	Targets []string `json:"targets"`
}

// RevokeRoleRequest revoke role request.
type RevokeRoleRequest struct {
	// InstanceID role instanceID.
	InstanceID string `json:"instanceId" validate:"required"`
	// Targets target resource instanceIds.
	Targets []string `json:"targets"`
}

// BatchRevokeRoleRequest batch revoke role request.
type BatchRevokeRoleRequest struct {
	// InstanceIds role instanceIds.
	InstanceIds []string `json:"instanceIds" validate:"required"`
	// Targets target resource instanceIds.
	Targets []string `json:"targets"`
}

// AuthorizeResources authorize role resource request.
type AuthorizeResources struct {
	// InstanceID role instanceID.
	InstanceID string `json:"instanceId" validate:"required"`
	// Resources authorize resource.
	Resources []string `json:"targets"`
}

// CreateRoleRequest create role request.
type CreateRoleRequest struct {
	Name        string `json:"name"            validate:"required,name"`
	Owner       string `json:"owner,omitempty" validate:"min=1,max=30"  optional:"true"`
	DisplayName string `json:"displayName"                              optional:"true"`
	Description string `json:"description"                              optional:"true"`
}

// UpdateRoleRequest update role request.
type UpdateRoleRequest struct {
	DisplayName string `json:"displayName" optional:"true"`
	Description string `json:"description" optional:"true"`
	Owner       string `json:"owner"       optional:"true"`
}

// RoleBase represents a role restful resource.
type RoleBase struct {
	// Standard object's metadata.
	metav1.ObjectMeta `       json:"metadata,omitempty"`
	// Owner tenant owner name.
	Owner string `json:"owner,omitempty"    gorm:"column:owner"       validate:"required"`
	// Disabled tenant state.
	Disabled bool `json:"disabled"           gorm:"column:disabled"`
	// Description tenant description.
	Description string `json:"description"        gorm:"column:description"`
}

// RoleList is the whole list of all roles which have been stored in stroage.
type RoleList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*RoleBase `json:"items"`
}

// DetailRoleResponse role detail.
type DetailRoleResponse struct {
	RoleBase
	Users []UserBase `json:"users"`
}

// AuthenticateRequest is the request body for login.
type AuthenticateRequest struct {
	Username string `json:"username,omitempty" optional:"true"`
	Password string `json:"password,omitempty" optional:"true"`
}

// AuthenticateResponse is the response of login request.
type AuthenticateResponse struct {
	// User user info
	User *UserBase `json:"user"`

	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken string `json:"access_token"`

	// TokenType is the type of token.
	// The Type method returns either this or "Bearer", the default.
	TokenType string `json:"token_type,omitempty"`

	// RefreshToken is a token that's used by the application
	// (as opposed to the user) to refresh the access token
	// if it expires.
	RefreshToken string `json:"refresh_token,omitempty"`

	// ExpiresIn is the optional expiration second of the access token.
	ExpiresIn int `json:"expires_in,omitempty"`
}

// RefreshTokenResponse is the response of refresh token request.
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// CreateResourceRequest create resource request.
type CreateResourceRequest struct {
	Name        string   `json:"name"                validate:"required,name"`
	Method      string   `json:"method"              validate:"required"`
	Type        string   `json:"type"                validate:"required"`
	Api         string   `json:"api"                 validate:"required"`
	Description string   `json:"description"`
	IsDefault   bool     `json:"isDefault,omitempty"`
	Actions     []Action `json:"actions,omitempty"   validate:"required"`
}

// UpdateResourceRequest update resource request.
type UpdateResourceRequest struct {
	Name        string   `json:"name"                validate:"required,name"`
	Method      string   `json:"method"              validate:"required"`
	Type        string   `json:"type"                validate:"required"`
	Api         string   `json:"api"                 validate:"required"`
	Description string   `json:"description"`
	IsDefault   bool     `json:"isDefault,omitempty"`
	Actions     []Action `json:"actions,omitempty"   validate:"required"`
}

// ResourceBase represents a resource restful resource.
type ResourceBase struct {
	// Standard object's metadata.
	metav1.ObjectMeta `         json:"metadata,omitempty"`
	// Method resource method.
	Method string `json:"method"             validate:"required"`
	// Type unique identification of resource API.
	Api string `json:"api"                validate:"required"`
	// Type resource type.
	Type string `json:"type"               validate:"required"`
	// Description resource description.
	Description string `json:"description"`
	// Actions resource access mode.
	Actions []Action `json:"actions,omitempty"  validate:"required"`
}

// DetailResourceResponse resource detail.
type DetailResourceResponse struct {
	ResourceBase
}

// Action resource access mode.
type Action struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (a *Action) String() string {
	data, _ := json.Marshal(a)
	return string(data)
}

// ResourceList is the whole list of all resource which have been stored in stroage.
type ResourceList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*ResourceBase `json:"items"`
}