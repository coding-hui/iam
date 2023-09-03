// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/internal/pkg/options"
	"github.com/coding-hui/iam/internal/pkg/token"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

const (
	jwtIssuer = "iam-issuer"
	audience  = "iam.api.wecoding.top"

	// GrantTypeAccess is the grant type for access token.
	GrantTypeAccess = "access"
	// GrantTypeRefresh is the grant type for refresh token.
	GrantTypeRefresh = "refresh"
)

var signedKey string

// AuthenticationService authentication service.
type AuthenticationService interface {
	Authenticate(ctx context.Context, loginReq v1.AuthenticateRequest) (*v1.AuthenticateResponse, error)
	AuthenticateByProvider(ctx context.Context, loginReq v1.AuthenticateRequest) (*v1.AuthenticateResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*v1.RefreshTokenResponse, error)
}

type authenticationServiceImpl struct {
	cfg         config.Config
	Store       repository.Factory `inject:"repository"`
	UserService UserService        `inject:""`
}

// NewAuthenticationService new authentication service.
func NewAuthenticationService(c config.Config) AuthenticationService {
	signedKey = c.JwtOptions.Key
	return &authenticationServiceImpl{cfg: c}
}

type authHandler interface {
	authenticate(ctx context.Context) (*v1.UserBase, error)
}

type localHandlerImpl struct {
	store       repository.Factory
	userService UserService
	username    string
	password    string
}

func (a *authenticationServiceImpl) newLocalHandler(loginReq v1.AuthenticateRequest) (*localHandlerImpl, error) {
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

func (a *authenticationServiceImpl) Authenticate(ctx context.Context, loginReq v1.AuthenticateRequest) (*v1.AuthenticateResponse, error) {
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
	if userBase.Disabled {
		return nil, errors.WithCode(code.ErrUserHasDisabled, "The account [%d] has been disabled.", userBase.Name)
	}
	accessToken, err := a.generateJWTToken(userBase.InstanceID, userBase.UserType, GrantTypeAccess, a.cfg.JwtOptions.Timeout)
	if err != nil {
		return nil, err
	}
	refreshToken, err := a.generateJWTToken(userBase.InstanceID, userBase.UserType, GrantTypeRefresh, a.cfg.JwtOptions.MaxRefresh)
	if err != nil {
		return nil, err
	}

	return &v1.AuthenticateResponse{
		User:        userBase,
		AccessToken: accessToken,
		// The OAuth 2.0 token_type response parameter value MUST be Bearer,
		// as specified in OAuth 2.0 Bearer Token Usage [RFC6750]
		TokenType:    "Bearer",
		RefreshToken: refreshToken,
		ExpiresIn:    int(a.cfg.JwtOptions.Timeout.Seconds()),
	}, nil
}

func (a *authenticationServiceImpl) AuthenticateByProvider(
	ctx context.Context,
	loginReq v1.AuthenticateRequest,
) (*v1.AuthenticateResponse, error) {
	providerOptions, err := a.cfg.OAuthOptions.IdentityProviderOptions(loginReq.Provider)
	if err != nil {
		return nil, err
	}
	genericProvider, err := identityprovider.GetGenericProvider(providerOptions.Name)
	if err != nil {
		return nil, err
	}
	authenticated, err := genericProvider.Authenticate(loginReq)
	if err != nil {
		return nil, err
	}
	linkedAccount, err := a.Store.UserRepository().
		GetByExternalId(ctx, authenticated.GetUserID(), providerOptions.Name, metav1.GetOptions{})
	if err != nil && !errors.IsCode(err, code.ErrUserNotFound) {
		return nil, err
	}

	var userBase *v1.UserBase

	if linkedAccount != nil {
		userBase = assembler.ConvertUserModelToBase(linkedAccount)
	}
	// the user will automatically create and mapping when login successful.
	if userBase == nil && providerOptions.MappingMethod == options.MappingMethodAuto {
		createResp, err := a.UserService.CreateUser(ctx, mappedUser(providerOptions.Name, authenticated))
		if err != nil {
			return nil, err
		}
		userBase = &createResp.User
	}

	accessToken, err := a.generateJWTToken(userBase.InstanceID, userBase.UserType, GrantTypeAccess, a.cfg.JwtOptions.Timeout)
	if err != nil {
		return nil, err
	}
	refreshToken, err := a.generateJWTToken(userBase.InstanceID, userBase.UserType, GrantTypeRefresh, a.cfg.JwtOptions.MaxRefresh)
	if err != nil {
		return nil, err
	}

	return &v1.AuthenticateResponse{
		User:        userBase,
		AccessToken: accessToken,
		// The OAuth 2.0 token_type response parameter value MUST be Bearer,
		// as specified in OAuth 2.0 Bearer Token Usage [RFC6750]
		TokenType:    "Bearer",
		RefreshToken: refreshToken,
		ExpiresIn:    int(a.cfg.JwtOptions.Timeout.Seconds()),
	}, nil
}

func (a *authenticationServiceImpl) RefreshToken(_ context.Context, refreshToken string) (*v1.RefreshTokenResponse, error) {
	claim, err := token.ParseToken(refreshToken, signedKey)
	if err != nil {
		if errors.IsCode(err, code.ErrExpired) {
			return nil, errors.WithCode(code.ErrExpired, jwt.ErrTokenExpired.Error())
		}
		return nil, err
	}
	if claim.GrantType == GrantTypeRefresh {
		accessToken, err := a.generateJWTToken(claim.UserInstanceId, claim.UserType, GrantTypeAccess, time.Hour)
		if err != nil {
			return nil, err
		}
		return &v1.RefreshTokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil
	}
	return nil,
		errors.WithCode(code.ErrInvalidRefreshToken,
			"Incorrect refresh token format %s, expected in refresh format", claim.GrantType)
}

func (a *authenticationServiceImpl) generateJWTToken(userInstanceId, userType, grantType string, expiresIn time.Duration) (string, error) {
	issueAt := time.Now()
	claims := token.Token{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			IssuedAt:  jwt.NewNumericDate(issueAt),
			Audience:  jwt.ClaimStrings{audience},
			Subject:   userInstanceId,
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(issueAt.Add(expiresIn)),
		},
		UserInstanceId: userInstanceId,
		GrantType:      grantType,
		UserType:       userType,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return t.SignedString([]byte(signedKey))
}

func (l *localHandlerImpl) authenticate(ctx context.Context) (*v1.UserBase, error) {
	user, err := l.userService.GetUser(ctx, l.username, metav1.GetOptions{})
	if err != nil {
		if errors.IsCode(err, code.ErrUserNotFound) {
			return nil, errors.WithCode(code.ErrPasswordIncorrect, err.Error())
		}
		return nil, err
	}
	if err := user.Compare(l.password); err != nil {
		return nil, err
	}
	if err := l.userService.FlushLastLoginTime(ctx, user); err != nil {
		return nil, err
	}

	return assembler.ConvertUserModelToBase(user), nil
}

func mappedUser(idp string, identity identityprovider.Identity) v1.CreateUserRequest {
	// username convert
	username := strings.ToLower(identity.GetUsername())
	alias := username
	if len(username) > 12 {
		alias = username[:12]
	}
	return v1.CreateUserRequest{
		Name:             username,
		Alias:            alias,
		Email:            identity.GetEmail(),
		Avatar:           identity.GetAvatar(),
		UserType:         v1.Default.String(),
		ExternalUID:      identity.GetUserID(),
		IdentifyProvider: idp,
	}
}
