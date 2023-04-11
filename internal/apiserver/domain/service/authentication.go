// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/coding-hui/common/errors"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	convert "github.com/coding-hui/iam/internal/apiserver/interfaces/api/convert/v1alpha1"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"
)

const (
	jwtIssuer = "iam-issuer"
	audience  = "iam.api.wecoding.top"

	// GrantTypeAccess is the grant type for access token
	GrantTypeAccess = "access"
	// GrantTypeRefresh is the grant type for refresh token
	GrantTypeRefresh = "refresh"
)

var signedKey string

// AuthenticationService authentication service
type AuthenticationService interface {
	Authenticate(ctx context.Context, loginReq v1alpha1.AuthenticateRequest) (*v1alpha1.AuthenticateResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*v1alpha1.RefreshTokenResponse, error)
}

type authenticationServiceImpl struct {
	cfg         config.Config
	Store       repository.Factory `inject:"repository"`
	UserService UserService        `inject:""`
}

// NewAuthenticationService new authentication service
func NewAuthenticationService(c config.Config) AuthenticationService {
	signedKey = c.JwtOptions.Key
	return &authenticationServiceImpl{cfg: c}
}

type authHandler interface {
	authenticate(ctx context.Context) (*v1alpha1.UserBase, error)
}

type localHandlerImpl struct {
	store       repository.Factory
	userService UserService
	username    string
	password    string
}

func (a *authenticationServiceImpl) newLocalHandler(loginReq v1alpha1.AuthenticateRequest) (*localHandlerImpl, error) {
	if loginReq.Username == "" || loginReq.Password == "" {
		return nil, errors.WithCode(code.ErrMissingLoginValues, "Missing Username or Password")
	}

	return &localHandlerImpl{
		store:       a.Store,
		userService: a.UserService,
		username:    loginReq.Username,
		password:    loginReq.Password,
	}, nil
}

func (a *authenticationServiceImpl) Authenticate(ctx context.Context, loginReq v1alpha1.AuthenticateRequest) (*v1alpha1.AuthenticateResponse, error) {
	var handler authHandler
	var err error
	handler, err = a.newLocalHandler(loginReq)
	if err != nil {
		return nil, err
	}
	userBase, err := handler.authenticate(ctx)
	if err != nil {
		return nil, err
	}
	accessToken, err := a.generateJWTToken(userBase.Name, GrantTypeAccess, a.cfg.JwtOptions.Timeout)
	if err != nil {
		return nil, err
	}
	refreshToken, err := a.generateJWTToken(userBase.Name, GrantTypeRefresh, a.cfg.JwtOptions.MaxRefresh)
	if err != nil {
		return nil, err
	}

	return &v1alpha1.AuthenticateResponse{
		User:        userBase,
		AccessToken: accessToken,
		// The OAuth 2.0 token_type response parameter value MUST be Bearer,
		// as specified in OAuth 2.0 Bearer Token Usage [RFC6750]
		TokenType:    "Bearer",
		RefreshToken: refreshToken,
		ExpiresIn:    int(a.cfg.JwtOptions.Timeout.Seconds()),
	}, nil
}

func (a *authenticationServiceImpl) RefreshToken(_ context.Context, refreshToken string) (*v1alpha1.RefreshTokenResponse, error) {
	claim, err := ParseToken(refreshToken)
	if err != nil {
		if errors.IsCode(err, code.ErrExpired) {
			return nil, errors.WithCode(code.ErrExpired, jwt.ErrTokenExpired.Error())
		}
		return nil, err
	}
	if claim.GrantType == GrantTypeRefresh {
		accessToken, err := a.generateJWTToken(claim.Username, GrantTypeAccess, time.Hour)
		if err != nil {
			return nil, err
		}
		return &v1alpha1.RefreshTokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil
	}
	return nil, err
}

// ParseToken parses and verifies a token
func ParseToken(tokenString string) (*model.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signedKey), nil
	})
	if err != nil {
		var ve *jwt.ValidationError
		if jwtErr := errors.As(err, &ve); jwtErr {
			switch ve.Errors {
			case jwt.ValidationErrorExpired:
				return nil, errors.WithCode(code.ErrExpired, jwt.ErrTokenExpired.Error())
			case jwt.ValidationErrorNotValidYet:
				return nil, errors.WithCode(code.ErrSignatureInvalid, jwt.ErrTokenNotValidYet.Error())
			case jwt.ValidationErrorMalformed:
				return nil, errors.WithCode(code.ErrEncrypt, jwt.ErrTokenMalformed.Error())
			default:
				return nil, errors.WithCode(code.ErrSignatureInvalid, err.Error())
			}
		}
		return nil, err
	}
	if claims, ok := token.Claims.(*model.CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.WithCode(code.ErrSignatureInvalid, err.Error())
}

func (a *authenticationServiceImpl) generateJWTToken(username, grantType string, expiresIn time.Duration) (string, error) {
	issueAt := time.Now()
	claims := model.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			IssuedAt:  jwt.NewNumericDate(issueAt),
			Audience:  jwt.ClaimStrings{audience},
			Subject:   username,
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(issueAt.Add(expiresIn)),
		},
		Username:  username,
		GrantType: grantType,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signedKey))
}

func (l *localHandlerImpl) authenticate(ctx context.Context) (*v1alpha1.UserBase, error) {
	user, err := l.userService.Get(ctx, l.username, metav1alpha1.GetOptions{})
	if err != nil {
		if errors.IsCode(err, code.ErrUserNotFound) {
			return nil, errors.WithCode(code.ErrPasswordIncorrect, "Password was incorrect")
		}
		return nil, err
	}
	if err := passwordVerify(user.Password, l.password); err != nil {
		return nil, err
	}

	return convert.ConvertUserModelToBase(user), nil
}

func passwordVerify(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return errors.WithCode(code.ErrPasswordIncorrect, "Password was incorrect")
	}

	return err
}
