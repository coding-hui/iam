// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"fmt"
	"time"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/infrastructure/cache"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/internal/pkg/token"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
)

// TokenService token service.
type TokenService interface {
	// Verify the given token and returns token.VerifiedResponse
	Verify(token string) (*token.VerifiedResponse, error)
	// IssueTo issue a token for the specified user
	IssueTo(request *token.IssueRequest) (string, error)
	// Revoke revoke the specified token
	Revoke(token string) error
	// RevokeAllUserTokens revoke all user tokens
	RevokeAllUserTokens(username string) error
	// Keys hold encryption and signing keys.
	Keys() *token.Keys
}

type tokenServiceImpl struct {
	cfg    config.Config
	issuer token.Issuer

	Cache cache.Interface `inject:"cache"`
}

// NewTokenService new token service.
func NewTokenService(c config.Config, issuer token.Issuer) TokenService {
	return &tokenServiceImpl{
		cfg:    c,
		issuer: issuer,
	}
}

func (t *tokenServiceImpl) Verify(tokenStr string) (*token.VerifiedResponse, error) {
	response, err := t.issuer.Verify(tokenStr)
	if err != nil {
		return nil, err
	}
	if t.cfg.AuthenticationOptions.OAuthOptions.AccessTokenMaxAge == 0 ||
		response.TokenType == token.StaticToken {
		return response, nil
	}
	if err := t.tokenCacheValidate(response.User.GetName(), tokenStr); err != nil {
		return nil, errors.WithCode(code.ErrTokenInvalid, err.Error())
	}
	return response, nil
}

func (t *tokenServiceImpl) IssueTo(request *token.IssueRequest) (string, error) {
	tokenStr, err := t.issuer.IssueTo(request)
	if err != nil {
		log.Errorf("failed to issue token: %v", err)
		return "", err
	}
	if request.ExpiresIn > 0 {
		if err = t.cacheToken(request.User.GetName(), tokenStr, request.ExpiresIn); err != nil {
			return "", err
		}
	}
	return tokenStr, nil
}

func (t *tokenServiceImpl) Revoke(token string) error {
	pattern := fmt.Sprintf("iam:user:*:token:%s", token)
	if keys, err := t.Cache.Keys(pattern); err != nil {
		log.Errorf("failed to get token keys: %v", err)
		return err
	} else if len(keys) > 0 {
		if err := t.Cache.Del(keys...); err != nil {
			log.Errorf("failed to del token cache: %v", err)
			return err
		}
	}
	return nil
}

func (t *tokenServiceImpl) RevokeAllUserTokens(username string) error {
	pattern := fmt.Sprintf("iam:user:%s:token:*", username)
	if keys, err := t.Cache.Keys(pattern); err != nil {
		log.Errorf("failed to get token keys: %v", err)
		return err
	} else if len(keys) > 0 {
		if err := t.Cache.Del(keys...); err != nil {
			log.Errorf("failed to del token cache: %v", err)
			return err
		}
	}
	return nil
}

func (t *tokenServiceImpl) Keys() *token.Keys {
	return t.issuer.Keys()
}

// tokenCacheValidate verify that the token is in the cache.
func (t *tokenServiceImpl) tokenCacheValidate(username, token string) error {
	key := fmt.Sprintf("iam:user:%s:token:%s", username, token)
	if exist, err := t.Cache.Exists(key); err != nil {
		return err
	} else if !exist {
		return errors.New("token not found in cache")
	}
	return nil
}

// cacheToken cache the token for a period of time.
func (t *tokenServiceImpl) cacheToken(username, token string, duration time.Duration) error {
	key := fmt.Sprintf("iam:user:%s:token:%s", username, token)
	if err := t.Cache.Set(key, token, duration); err != nil {
		log.Errorf("failed to cache token: %v", err)
		return err
	}
	return nil
}
