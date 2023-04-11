// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1alpha1

import (
	"time"

	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
)

var (
	// CtxKeyUserName request context key of username
	CtxKeyUserName = "username"
)

type UserState string

// These are the valid phases of a user.
const (
	// UserActive means the user is available.
	UserActive UserState = "Active"
	// UserDisabled means the user is disabled.
	UserDisabled UserState = "Disabled"
	// UserAuthLimitExceeded means restrict user login.
	UserAuthLimitExceeded UserState = "AuthLimitExceeded"

	AuthenticatedSuccessfully = "authenticated successfully"
)

// UserBase represents a user restful resource.
type UserBase struct {
	// Standard object's metadata.
	metav1alpha1.ObjectMeta `json:"metadata,omitempty"`

	TenantId uint64 `json:"tenantId,omitempty"`

	Status int `json:"status" validate:"omitempty"`

	Alias string `json:"alias" validate:"min=1,max=30"`

	// Required: true
	Password string `json:"password,omitempty" validate:"required"`

	// Required: true
	Email string `json:"email" validate:"required,email,min=1,max=100"`

	Phone string `json:"phone" validate:"omitempty"`

	Disabled bool `json:"disabled"`

	LastLoginTime *time.Time `json:"last_login_time,omitempty"`
}

// UserList is the whole list of all users which have been stored in stroage.
type UserList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// +optional
	metav1alpha1.ListMeta `json:",inline"`

	Items []*UserBase `json:"items"`
}

type LoginRecord struct {
	metav1alpha1.TypeMeta   `json:",inline"`
	metav1alpha1.ObjectMeta `json:"metadata,omitempty"`
	Spec                    LoginRecordSpec `json:"spec"`
}

type LoginRecordSpec struct {
	// Which authentication method used, BasicAuth/OAuth/Token
	Type LoginType `json:"type"`
	// Provider of authentication, Local/Ldap/Github etc.
	Provider string `json:"provider"`
	// Source IP of client
	SourceIP string `json:"sourceIP"`
	// User agent of login attempt
	UserAgent string `json:"userAgent,omitempty"`
	// Successful login attempt or not
	Success bool `json:"success"`
	// States failed login attempt reason
	Reason string `json:"reason"`
}

type LoginType string

const (
	BasicAuth LoginType = "Basic"
	OAuth     LoginType = "OAuth"
	Token     LoginType = "Token"
)

// AuthenticateRequest is the request body for login
type AuthenticateRequest struct {
	Username string `json:"username,omitempty" optional:"true"`
	Password string `json:"password,omitempty" optional:"true"`
}

// AuthenticateResponse is the response of login request
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

// RefreshTokenResponse is the response of refresh token request
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
