// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

type WechatMiniAppCodeOptions struct {
	// Iv 对称解密算法初始向量，由微信返回
	Iv string `json:"iv"`
	// EncryptedData 获取微信开放数据返回的加密数据（encryptedData）
	EncryptedData string `json:"encryptedData"`
	// Code wx.login 接口返回的用户 code
	Code string `json:"code"`
}

// LoginByMobileRequest is the request body for mobile login.
type LoginByMobileRequest struct {
	Provider                 string                   `json:"provider"`
	WechatMiniAppCodeOptions WechatMiniAppCodeOptions `json:"wechatMiniAppCodeOptions"`
}
