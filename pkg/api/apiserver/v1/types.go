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
	// CtxKeyRole request context key of role.
	CtxKeyRole = "role"
	// CtxKeyPolicy request context key of policy.
	CtxKeyPolicy = "policy"
	// CtxKeyResource request context key of resource.
	CtxKeyResource = "resource"
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

// ResourceType resource type.
type ResourceType string

const (
	// API api resource.
	API ResourceType = "API"
)

// ListUserOptions list user options.
type ListUserOptions struct {
	metav1.ListOptions         `       json:",inline"`
	DepartmentID               string `json:"departmentId,omitempty"               form:"departmentId"`
	IncludeChildrenDepartments bool   `json:"includeChildrenDepartments,omitempty" form:"includeChildrenDepartments"`
}

// CreateUserRequest create user request.
type CreateUserRequest struct {
	Name             string   `json:"name"                       validate:"required,name"`
	Password         string   `json:"password"                   validate:"required"`
	Alias            string   `json:"alias,omitempty"`
	Email            string   `json:"email,omitempty"`
	Phone            string   `json:"phone,omitempty"`
	UserType         string   `json:"userType,omitempty"`
	Avatar           string   `json:"avatar,omitempty"`
	IdentifyProvider string   `json:"identifyProvider,omitempty"`
	ExternalUID      string   `json:"externalUID,omitempty"`
	RoleIds          []string `json:"roleIds,omitempty"`
	DepartmentIds    []string `json:"departmentIds,omitempty"`
}

// CreateUserResponse create user response.
type CreateUserResponse struct {
	UserBase `json:",inline"`
}

// UpdateUserRequest update user request.
type UpdateUserRequest struct {
	Alias    string `json:"alias"              validate:"required,min=1,max=30"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Password string `json:"password,omitempty"`
}

// UpdateUserResponse update user response.
type UpdateUserResponse struct {
	UserBase
}

// UserBase represents a user restful resource.
type UserBase struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            int        `json:"status"`
	Alias             string     `json:"alias"`
	Password          string     `json:"password,omitempty"`
	Email             string     `json:"email"`
	Phone             string     `json:"phone"`
	UserType          string     `json:"userType"`
	Disabled          bool       `json:"disabled"`
	Avatar            string     `json:"avatar"`
	DepartmentIds     []string   `json:"departmentIds,omitempty"`
	LastLoginTime     *time.Time `json:"lastLoginTime,omitempty"`
}

// DetailUserResponse is the response of user detail.
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

	Items []*DetailUserResponse `json:"items"`
}

// TenantBase represents a tenant restful resource.
type TenantBase struct {
	// Standard object's metadata.
	metav1.ObjectMeta `       json:"metadata,omitempty"`
	Owner             string `json:"owner,omitempty"    validate:"required"`
	Disabled          bool   `json:"disabled"`
	Description       string `json:"description"`
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
	Name        string `json:"name"                  validate:"required,name"`
	Owner       string `json:"owner,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateRoleRequest update role request.
type UpdateRoleRequest struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description,omitempty"`
	Owner       string `json:"owner,omitempty"`
}

// RoleBase represents a role restful resource.
type RoleBase struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Owner       string `json:"owner,omitempty"`
	Disabled    bool   `json:"disabled"`
	Description string `json:"description"`
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

type WechatMiniAppCodePayload struct {
	// Iv 对称解密算法初始向量，由微信返回
	Iv string `json:"iv"`
	// EncryptedData 获取微信开放数据返回的加密数据（encryptedData）
	EncryptedData string `json:"encryptedData"`
	// Code wx.login 接口返回的用户 code
	Code string `json:"code"`
}

// AuthenticateRequest is the request body for login.
type AuthenticateRequest struct {
	Username                 string                   `json:"username,omitempty"`
	Password                 string                   `json:"password,omitempty"`
	Provider                 string                   `json:"provider,omitempty"`
	WechatMiniAppCodePayload WechatMiniAppCodePayload `json:"wechatMiniAppCodePayload,omitempty"`
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
	Type        string   `json:"type"                validate:"required"`
	Api         string   `json:"api"                 validate:"required"`
	Method      string   `json:"method"`
	Description string   `json:"description"`
	IsDefault   bool     `json:"isDefault,omitempty"`
	Actions     []Action `json:"actions,omitempty"`
}

// UpdateResourceRequest update resource request.
type UpdateResourceRequest struct {
	Name        string   `json:"name"                validate:"required,name"`
	Type        string   `json:"type"                validate:"required"`
	Api         string   `json:"api"                 validate:"required"`
	Method      string   `json:"method"`
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

// CreateOrganizationRequest create organization request.
type CreateOrganizationRequest struct {
	Name        string `json:"name"                  validate:"required,name"`
	DisplayName string `json:"displayName,omitempty"`
	WebsiteUrl  string `json:"websiteUrl,omitempty"`
	Favicon     string `json:"favicon,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateOrganizationRequest update organization request.
type UpdateOrganizationRequest struct {
	DisplayName string `json:"displayName"           validate:"required,min=1,max=30"`
	WebsiteUrl  string `json:"websiteUrl,omitempty"`
	Favicon     string `json:"favicon,omitempty"`
	Description string `json:"description,omitempty"`
}

// OrganizationBase represents a organization restful resource.
type OrganizationBase struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	DisplayName  string `json:"displayName"`
	WebsiteUrl   string `json:"websiteUrl"`
	Favicon      string `json:"favicon"`
	Disabled     bool   `json:"disabled"`
	Description  string `json:"description"`
	IsLeaf       bool   `json:"isLeaf"`
	ParentID     string `json:"parentId"`
	Owner        string `json:"owner"`
	MembersCount int64  `json:"membersCount"`
}

// DetailOrganizationResponse org detail.
type DetailOrganizationResponse struct {
	OrganizationBase
}

// OrganizationList is the whole list of all org which have been stored in stroage.
type OrganizationList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*DetailOrganizationResponse `json:"items"`
}

const (
	// AllowAccess should be used as effect for policies that allow access.
	AllowAccess string = "allow"
	// DenyAccess should be used as effect for policies that deny access.
	DenyAccess string = "deny"
)

// PolicyType define policy type.
type PolicyType string

const (
	// SystemBuildInPolicy system default policy type.
	SystemBuildInPolicy PolicyType = "SYSTEM"
	// CustomPolicy user custom policy type.
	CustomPolicy PolicyType = "CUSTOM"
)

type Statement struct {
	Effect             string   `json:"effect"             validate:"required"`
	Resource           string   `json:"resource"           validate:"required"`
	ResourceIdentifier string   `json:"resourceIdentifier" validate:"required"`
	Actions            []string `json:"actions"            validate:"required"`
}

// CreatePolicyRequest create policy request.
type CreatePolicyRequest struct {
	Name        string      `json:"name"        validate:"required,name"`
	Type        string      `json:"type"        validate:"required"`
	Statements  []Statement `json:"statements"  validate:"required"`
	Subjects    []string    `json:"subjects"    validate:"required"`
	Description string      `json:"description"`
	Status      string      `json:"status"`
	Owner       string      `json:"owner"`
	Meta        string      `json:"meta"`
}

// UpdatePolicyRequest update policy request.
type UpdatePolicyRequest struct {
	Subjects    []string    `json:"subjects"    validate:"required"`
	Statements  []Statement `json:"statements"  validate:"required"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Status      string      `json:"status"`
	Owner       string      `json:"owner"`
	Meta        string      `json:"meta"`
}

// PolicyBase represents a policy restful resource.
type PolicyBase struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Subjects   []string    `json:"subjects"`
	Statements []Statement `json:"statements"`

	Type        string `json:"type"`
	Status      string `json:"status"`
	Owner       string `json:"owner"`
	Description string `json:"description"`

	// casbin required
	Adapter     string     `json:"adapter"`
	Model       string     `json:"model"`
	PolicyRules [][]string `json:"policyRules"`
}

// DetailPolicyResponse policy detail.
type DetailPolicyResponse struct {
	PolicyBase
	Resources []ResourceBase `json:"resources,omitempty"`
}

// PolicyList is the whole list of all policies which have been stored in stroage.
type PolicyList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*PolicyBase `json:"items"`
}

// DetailDepartmentResponse department detail.
type DetailDepartmentResponse struct {
	OrganizationBase `json:",inline"`

	OrganizationID string `json:"organizationId"`
}

// DepartmentList is the whole list of all departments which have been stored in stroage.
type DepartmentList struct {
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*DetailDepartmentResponse `json:"items"`
}

// CreateDepartmentRequest create organization department request.
type CreateDepartmentRequest struct {
	Name           string `json:"name"                  validate:"required,name"`
	OrganizationID string `json:"organizationId"        validate:"required"`
	ParentID       string `json:"parentId"              validate:"required"`
	DisplayName    string `json:"displayName,omitempty"`
	WebsiteUrl     string `json:"websiteUrl,omitempty"`
	Favicon        string `json:"favicon,omitempty"`
	Description    string `json:"description,omitempty"`
}

// UpdateDepartmentRequest update organization department request.
type UpdateDepartmentRequest struct {
	OrganizationID string `json:"organizationId"        validate:"required"`
	ParentID       string `json:"parentId"              validate:"required"`
	DisplayName    string `json:"displayName,omitempty"`
	WebsiteUrl     string `json:"websiteUrl,omitempty"`
	Favicon        string `json:"favicon,omitempty"`
	Description    string `json:"description,omitempty"`
}

// DepartmentMember department member.
type DepartmentMember struct {
	MemberID   string `json:"memberId"             validate:"required"`
	MemberType string `json:"memberType,omitempty"`
}

// AddDepartmentMemberRequest add department members request.
type AddDepartmentMemberRequest struct {
	DepartmentMember `json:",inline"`
}

// BatchAddDepartmentMemberRequest batch add department members request.
type BatchAddDepartmentMemberRequest struct {
	Members []DepartmentMember `json:"members"`
}

// RemoveDepartmentMemberRequest remove department members request.
type RemoveDepartmentMemberRequest struct {
	DepartmentMember `json:",inline"`
}

// BatchRemoveDepartmentMemberRequest batch remove department members request.
type BatchRemoveDepartmentMemberRequest struct {
	Members []DepartmentMember `json:"members"`
}

// DepartmentMemberList is the whole list of all department members which have been stored in stroage.
type DepartmentMemberList struct {
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Members []*DepartmentMember `json:"members"`
}

type (
	ProviderCategory     string
	GrantHandlerType     string
	MappingMethod        string
	IdentityProviderType string
)

const (
	// GrantHandlerAuto auto-approves client authorization grant requests.
	GrantHandlerAuto GrantHandlerType = "auto"
	// GrantHandlerPrompt prompts the user to approve new client authorization grant requests.
	GrantHandlerPrompt GrantHandlerType = "prompt"
	// GrantHandlerDeny auto-denies client authorization grant requests.
	GrantHandlerDeny GrantHandlerType = "deny"
	// MappingMethodAuto  The default value.
	// The user will automatically create and mapping when login successful.
	// Fails if a user with that username is already mapped to another identity.
	MappingMethodAuto MappingMethod = "auto"
	// MappingMethodLookup Looks up an existing identity, user identity mapping, and user, but does not automatically
	// provision users or identities. Using this method requires you to manually provision users.
	MappingMethodLookup MappingMethod = "lookup"
	// MappingMethodMixed  A user entity can be mapped with multiple identifyProvider.
	// not supported yet.
	MappingMethodMixed MappingMethod = "mixed"

	OAuth ProviderCategory = "OAuth"
	Email ProviderCategory = "Email"

	GithubIdentityProvider            IdentityProviderType = "GitHub"
	CodingIdentityProvider            IdentityProviderType = "Coding"
	GiteeIdentityProvider             IdentityProviderType = "Gitee"
	LDAPIdentityProvider              IdentityProviderType = "LDAP"
	WeChatMiniProgramIdentityProvider IdentityProviderType = "WeChatMiniProgram"
	GoogleIdentityProvider            IdentityProviderType = "Google"
)

// CreateIdentityProviderRequest create IdentityProvider request.
type CreateIdentityProviderRequest struct {
	Type          IdentityProviderType `json:"type"                    validate:"required"`
	Category      ProviderCategory     `json:"category"`
	MappingMethod MappingMethod        `json:"mappingMethod,omitempty"`

	Name        string `json:"name"        validate:"required,name"`
	Status      string `json:"status"`
	Owner       string `json:"owner"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	CallbackURL string `json:"callbackURL"`

	Config metav1.Extend `json:"config,omitempty" validate:"omitempty"`
}

// UpdateIdentityProviderRequest update IdentityProvider request.
type UpdateIdentityProviderRequest struct {
	Category      ProviderCategory `json:"category"`
	MappingMethod MappingMethod    `json:"mappingMethod,omitempty"`

	Status      string `json:"status"`
	Owner       string `json:"owner"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	CallbackURL string `json:"callbackURL"`

	Config metav1.Extend `json:"config,omitempty" validate:"omitempty"`
}

type IdentityProviderBase struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Type          IdentityProviderType `json:"type"`
	Category      ProviderCategory     `json:"category"`
	MappingMethod MappingMethod        `json:"mappingMethod"`

	Status      string `json:"status"`
	Owner       string `json:"owner"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	CallbackURL string `json:"callbackURL"`

	Config metav1.Extend `json:"config"`
}

// IdentityProviderList is the whole list of all IdentityProvider which have been stored in stroage.
type IdentityProviderList struct {
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*DetailIdentityProviderResponse `json:"items"`
}

// DetailIdentityProviderResponse IdentityProvider detail.
type DetailIdentityProviderResponse struct {
	IdentityProviderBase `json:",inline"`
}

// CreateApplicationRequest create app request.
type CreateApplicationRequest struct {
	Name                string   `json:"name"                validate:"required,name"`
	Status              string   `json:"status"`
	Owner               string   `json:"owner"`
	DisplayName         string   `json:"displayName"`
	Description         string   `json:"description"`
	Logo                string   `json:"logo"`
	HomepageUrl         string   `json:"homepageUrl"`
	CallbackURL         string   `json:"callbackURL"`
	LoginURL            string   `json:"loginURL"`
	IdentityProviderIds []string `json:"identityProviderIds"`
}

// UpdateApplicationRequest update app request.
type UpdateApplicationRequest struct {
	DisplayName         string   `json:"displayName"`
	Status              string   `json:"status"`
	Owner               string   `json:"owner"`
	Description         string   `json:"description"`
	Logo                string   `json:"logo"`
	HomepageUrl         string   `json:"homepageUrl"`
	CallbackURL         string   `json:"callbackURL"`
	LoginURL            string   `json:"loginURL"`
	IdentityProviderIds []string `json:"identityProviderIds"`
	RefreshAppSecret    bool     `json:"refreshAppSecret"`
}

// ApplicationBase represents a application restful resource.
type ApplicationBase struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status      string `json:"status"`
	Owner       string `json:"owner"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	HomepageUrl string `json:"homepageUrl"`
	AppID       string `json:"appId"`
	AppSecret   string `json:"appSecret"`
	CallbackURL string `json:"callbackURL"`
	LoginURL    string `json:"loginURL"`

	IdentityProviders []IdentityProviderBase `json:"identityProviders"`
}

// DetailApplicationResponse application detail.
type DetailApplicationResponse struct {
	ApplicationBase `json:",inline"`
}

// ApplicationList is the whole list of all Application which have been stored in stroage.
type ApplicationList struct {
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*DetailApplicationResponse `json:"items"`
}

// ListApiKeyOptions list API Key options.
type ListApiKeyOptions struct {
	metav1.ListOptions `       json:",inline"`
	UserID             string `json:"userId,omitempty" form:"userId"`
	Status             int    `json:"status,omitempty" form:"status"`
}

// CreateApiKeyRequest create API Key request.
type CreateApiKeyRequest struct {
	Name      string    `json:"name"                validate:"required,name"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}

// CreateApiKeyResponse create API Key response.
type CreateApiKeyResponse struct {
	ApiKeyBase `json:",inline"`
}

// UpdateApiKeyRequest update API Key request.
type UpdateApiKeyRequest struct {
	Name      string    `json:"name"                validate:"required,name"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
	Status    int       `json:"status,omitempty"`
}

// ApiKeyBase represents an API Key restful resource.
type ApiKeyBase struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Name       string    `json:"name"`
	Key        string    `json:"key"`
	UserID     string    `json:"userId"`
	ExpiresAt  time.Time `json:"expiresAt,omitempty"`
	Status     int       `json:"status"`
	LastUsedAt time.Time `json:"lastUsedAt,omitempty"`
	UsageCount int64     `json:"usageCount"`
}

// ApiKeyList is the whole list of all API Keys which have been stored in storage.
type ApiKeyList struct {
	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*ApiKeyBase `json:"items"`
}

type IdentityProviderConfig struct {
	OAuthConfig `yaml:",inline" json:",inline"`
}

type OAuthConfig struct {
	// ClientID is the application's ID.
	ClientID string `json:"clientID" yaml:"clientID"`

	// ClientSecret is the application's secret.
	ClientSecret string `json:"-" yaml:"clientSecret"`

	// Endpoint contains the resource server's token endpoint
	// URLs. These are constants specific to each server and are
	// often available via site-specific packages, such as
	// google.Endpoint or gitee.endpoint.
	Endpoint struct {
		AuthURL     string `json:"authURL"     yaml:"authURL"`
		TokenURL    string `json:"tokenURL"    yaml:"tokenURL"`
		UserInfoURL string `json:"userInfoURL" yaml:"userInfoURL"`
	} `json:"endpoint" yaml:"endpoint"`

	// RedirectURL is the URL to redirect users going through
	// the OAuth flow, after the resource owner's URLs.
	RedirectURL string `json:"redirectURL" yaml:"redirectURL"`

	// Used to turn off TLS certificate checks
	InsecureSkipVerify bool `json:"insecureSkipVerify" yaml:"insecureSkipVerify"`

	// Scope specifies optional requested permissions.
	Scopes []string `json:"scopes" yaml:"scopes"`
}
