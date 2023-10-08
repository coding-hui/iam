// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/pkg/token"
	"github.com/coding-hui/iam/pkg/log"
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
}

// NewTokenService new token service.
func NewTokenService(c config.Config, issuer token.Issuer) TokenService {
	return &tokenServiceImpl{
		cfg:    c,
		issuer: issuer,
	}
}

func (t tokenServiceImpl) Verify(tokenStr string) (*token.VerifiedResponse, error) {
	response, err := t.issuer.Verify(tokenStr)
	if err != nil {
		return nil, err
	}
	if t.cfg.AuthenticationOptions.OAuthOptions.AccessTokenMaxAge == 0 ||
		response.TokenType == token.StaticToken {
		return response, nil
	}
	return response, nil
}

func (t tokenServiceImpl) IssueTo(request *token.IssueRequest) (string, error) {
	tokenStr, err := t.issuer.IssueTo(request)
	if err != nil {
		log.Errorf("failed to issue token: %v", err)
		return "", err
	}
	return tokenStr, nil
}

func (t tokenServiceImpl) Revoke(token string) error {
	//TODO implement me
	panic("implement me")
}

func (t tokenServiceImpl) RevokeAllUserTokens(username string) error {
	//TODO implement me
	panic("implement me")
}

func (t tokenServiceImpl) Keys() *token.Keys {
	return t.issuer.Keys()
}
