// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

// BindExternalAccountRequest 绑定第三方账号请求
// BindExternalAccountRequest bind external account request
type BindExternalAccountRequest struct {
	// Username 用户名
	// Username
	Username string `json:"username" validate:"required"`

	// Password 密码
	// Password
	Password string `json:"password" validate:"required"`

	// Provider 第三方认证提供商
	// Provider identity provider
	Provider string `json:"provider" validate:"required"`

	// ExternalUID 第三方用户ID
	// ExternalUID external user ID
	ExternalUID string `json:"externalUID" validate:"required"`

	// ExternalInfo 第三方用户信息
	// ExternalInfo external user info
	ExternalInfo *ExternalUserInfo `json:"externalInfo,omitempty"`
}

// UnbindExternalAccountRequest 解绑第三方账号请求
// UnbindExternalAccountRequest unbind external account request
type UnbindExternalAccountRequest struct {
	// Provider 第三方认证提供商
	// Provider identity provider
	Provider string `json:"provider" validate:"required"`
}

// ExternalUserInfo 第三方用户信息
// ExternalUserInfo external user info
type ExternalUserInfo struct {
	// Username 用户名
	// Username
	Username string `json:"username"`

	// Email 邮箱
	// Email
	Email string `json:"email,omitempty"`

	// Avatar 头像
	// Avatar
	Avatar string `json:"avatar,omitempty"`

	// DisplayName 显示名称
	// DisplayName
	DisplayName string `json:"displayName,omitempty"`
}

// BindExternalAccountResponse 绑定第三方账号响应
// BindExternalAccountResponse bind external account response
type BindExternalAccountResponse struct {
	// Success 是否绑定成功
	// Success whether the binding is successful
	Success bool `json:"success"`

	// Message 提示信息
	// Message message
	Message string `json:"message,omitempty"`
}
