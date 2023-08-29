// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package token

import (
	"github.com/golang-jwt/jwt/v4"

	"github.com/coding-hui/common/errors"

	"github.com/coding-hui/iam/internal/pkg/code"
)

// Token is the custom token claims.
type Token struct {
	jwt.RegisteredClaims
	// Private Claim Names
	// TokenType defined the type of the token
	TokenType      string `json:"token_type,omitempty"`
	UserInstanceId string `json:"user_instance_id"`
	UserType       string `json:"user_type"`
	GrantType      string `json:"grant_type"`
}

// ParseToken parses and verifies a token.
func ParseToken(tokenString, signedKey string) (*Token, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Token{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(signedKey), nil
		},
	)
	if err != nil {
		var ve *jwt.ValidationError
		if jwtErr := errors.As(err, &ve); jwtErr {
			switch ve.Errors {
			case jwt.ValidationErrorExpired:
				return nil, errors.WithCode(code.ErrExpired, err.Error())
			case jwt.ValidationErrorNotValidYet:
				return nil, errors.WithCode(code.ErrTokenNotValidYet, err.Error())
			case jwt.ValidationErrorMalformed:
				return nil, errors.WithCode(code.ErrTokenMalformed, err.Error())
			default:
				return nil, errors.WithCode(code.ErrTokenInvalid, err.Error())
			}
		}
		return nil, err
	}
	if claims, ok := token.Claims.(*Token); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.WithCode(code.ErrTokenInvalid, err.Error())
}
