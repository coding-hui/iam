// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package wechatmini

import (
	"github.com/mitchellh/mapstructure"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram/encryptor"

	"github.com/coding-hui/iam/internal/apiserver/domain/identityprovider"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/internal/pkg/options"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
	"github.com/coding-hui/common/util/idutil"
)

const (
	wechatMiniIdentityProvider = "WeChatMiniProgramIdentityProvider"
)

func init() {
	identityprovider.RegisterGenericProvider(&wechatMiniProviderFactory{})
}

type wechatMiniProvider struct {
	AppID     string `json:"app_id"     mapstructure:"app-id"`     // appid
	AppSecret string `json:"app_secret" mapstructure:"app-secret"` // appSecret

	miniprogram *miniprogram.MiniProgram
	cache       cache.Cache
}

type wechatMiniProviderFactory struct {
}

func (w *wechatMiniProviderFactory) Type() string {
	return wechatMiniIdentityProvider
}

func (w *wechatMiniProviderFactory) Create(opts options.DynamicOptions) (identityprovider.GenericProvider, error) {
	var provider wechatMiniProvider
	if err := mapstructure.Decode(opts, &provider); err != nil {
		return nil, err
	}
	provider.cache = cache.NewMemory()
	cfg := &miniConfig.Config{
		AppID:     provider.AppID,
		AppSecret: provider.AppSecret,
		Cache:     provider.cache,
	}
	provider.miniprogram = wechat.NewWechat().GetMiniProgram(cfg)
	return &provider, nil
}

// miniprogramIdentity 用户信息/手机号信息
type miniprogramIdentity struct {
	OpenID      string `json:"openId"`
	UnionID     string `json:"unionId"`
	NickName    string `json:"nickName"`
	AvatarURL   string `json:"avatarUrl"`
	PhoneNumber string `json:"phoneNumber"`
}

func (l *miniprogramIdentity) GetUserID() string {
	return l.OpenID
}

func (l *miniprogramIdentity) GetUsername() string {
	return l.NickName
}

func (l *miniprogramIdentity) GetEmail() string {
	return ""
}

func (l *miniprogramIdentity) GetAvatar() string {
	return l.AvatarURL
}

func (w wechatMiniProvider) Authenticate(req v1.AuthenticateRequest) (identityprovider.Identity, error) {
	payload := req.WechatMiniAppCodePayload
	session, err := w.miniprogram.GetAuth().Code2Session(payload.Code)
	if err != nil {
		log.Errorf("get wechat mini program session error: %v", err)
		return nil, err
	}
	sessionKey, openId, unionId := session.SessionKey, session.OpenID, session.UnionID
	if openId == "" && unionId == "" {
		return nil, errors.WithCode(code.ErrInvalidRequest, "the wechat mini program session is invalid")
	}
	var identity *encryptor.PlainData
	encryptedData, iv := payload.EncryptedData, payload.Iv
	if encryptedData != "" && iv != "" {
		identity, err = w.miniprogram.GetEncryptor().Decrypt(sessionKey, encryptedData, iv)
		if err != nil {
			log.Errorf("encryptor wechat mini program data error: %v", err)
			return nil, err
		}
	}

	return &miniprogramIdentity{
		OpenID:      openId,
		UnionID:     unionId,
		NickName:    idutil.GetUUID36("wxid_"),
		AvatarURL:   identity.AvatarURL,
		PhoneNumber: identity.PhoneNumber,
	}, nil
}
