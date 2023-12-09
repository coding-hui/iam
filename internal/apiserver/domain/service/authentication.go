// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/internal/pkg/token"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/log"
)

// AuthenticationService authentication service.
type AuthenticationService interface {
	Login(ctx context.Context, loginReq v1.AuthenticateRequest) (*v1.AuthenticateResponse, error)
	LoginByProvider(ctx context.Context, loginReq v1.AuthenticateRequest) (*v1.AuthenticateResponse, error)
	LoginByOAuthProvider(ctx context.Context, idp *model.IdentityProvider, req *http.Request) (*v1.AuthenticateResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*v1.RefreshTokenResponse, error)
}

type authenticationServiceImpl struct {
	cfg             config.Config
	Store           repository.Factory      `inject:"repository"`
	UserService     UserService             `inject:""`
	TokenService    TokenService            `inject:""`
	ProviderService IdentityProviderService `inject:""`
}

// NewAuthenticationService new authentication service.
func NewAuthenticationService(c config.Config) AuthenticationService {
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

func (a *authenticationServiceImpl) Login(ctx context.Context, loginReq v1.AuthenticateRequest) (*v1.AuthenticateResponse, error) {
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
		return nil, errors.WithCode(code.ErrUserHasDisabled, "The account [%s] has been disabled.", userBase.Name)
	}
	accessToken, err := a.TokenService.IssueTo(&token.IssueRequest{
		User:      *userBase,
		Claims:    token.Claims{TokenType: token.AccessToken},
		ExpiresIn: a.cfg.AuthenticationOptions.OAuthOptions.AccessTokenMaxAge,
	})
	if err != nil {
		return nil, err
	}
	refreshToken, err := a.TokenService.IssueTo(&token.IssueRequest{
		User:      *userBase,
		Claims:    token.Claims{TokenType: token.RefreshToken},
		ExpiresIn: a.cfg.AuthenticationOptions.OAuthOptions.AccessTokenMaxAge + a.cfg.AuthenticationOptions.OAuthOptions.AccessTokenInactivityTimeout,
	})
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
		ExpiresIn:    int(a.cfg.AuthenticationOptions.OAuthOptions.AccessTokenMaxAge.Seconds()),
	}, nil
}

func (a *authenticationServiceImpl) LoginByProvider(
	ctx context.Context,
	loginReq v1.AuthenticateRequest,
) (*v1.AuthenticateResponse, error) {
	provider, err := a.Store.IdentityProviderRepository().GetByName(ctx, loginReq.Provider, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	genericProvider, err := identityprovider.GetGenericProvider(provider)
	if err != nil {
		return nil, err
	}
	authenticated, err := genericProvider.Authenticate(loginReq)
	if err != nil {
		return nil, err
	}
	linkedAccount, err := a.Store.UserRepository().
		GetByExternalId(ctx, authenticated.GetUserID(), provider.Name, metav1.GetOptions{})
	if err != nil && !errors.IsCode(err, code.ErrUserNotFound) {
		return nil, err
	}

	var userBase *v1.UserBase

	if linkedAccount != nil {
		userBase = assembler.ConvertUserModelToBase(linkedAccount)
	}
	// the user will automatically create and mapping when login successful.
	if userBase == nil && provider.MappingMethod == v1.MappingMethodAuto {
		createResp, err := a.UserService.CreateUser(ctx, mappedUser(provider.Name, authenticated))
		if err != nil {
			return nil, err
		}
		userBase = &createResp.UserBase
	}

	go func() {
		err := a.Store.UserRepository().FlushLastLoginTime(ctx, userBase.Name)
		if err != nil {
			log.Errorf("Failed to flush user [%s] last login time: %v", userBase.Name, err)
		}
	}()

	accessToken, err := a.TokenService.IssueTo(&token.IssueRequest{
		User:      *userBase,
		Claims:    token.Claims{TokenType: token.AccessToken},
		ExpiresIn: a.cfg.AuthenticationOptions.OAuthOptions.AccessTokenMaxAge,
	})
	if err != nil {
		return nil, err
	}
	refreshToken, err := a.TokenService.IssueTo(&token.IssueRequest{
		User:      *userBase,
		Claims:    token.Claims{TokenType: token.RefreshToken},
		ExpiresIn: a.cfg.AuthenticationOptions.OAuthOptions.AccessTokenMaxAge + a.cfg.AuthenticationOptions.OAuthOptions.AccessTokenInactivityTimeout,
	})
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
		ExpiresIn:    int(a.cfg.AuthenticationOptions.OAuthOptions.AccessTokenMaxAge.Seconds()),
	}, nil
}

func (a *authenticationServiceImpl) LoginByOAuthProvider(
	ctx context.Context,
	idp *model.IdentityProvider,
	req *http.Request,
) (*v1.AuthenticateResponse, error) {
	oauthProvider, err := identityprovider.GetOAuthProvider(idp)
	if err != nil {
		return nil, err
	}
	authenticated, err := oauthProvider.IdentityExchangeCallback(req)
	if err != nil {
		return nil, err
	}
	linkedAccount, err := a.Store.UserRepository().
		GetByExternalId(ctx, authenticated.GetUserID(), idp.Name, metav1.GetOptions{})
	if err != nil && !errors.IsCode(err, code.ErrUserNotFound) {
		return nil, err
	}

	var userBase *v1.UserBase

	if linkedAccount != nil {
		userBase = assembler.ConvertUserModelToBase(linkedAccount)
	}
	// the user will automatically create and mapping when login successful.
	if userBase == nil && idp.MappingMethod == v1.MappingMethodAuto {
		createResp, err := a.UserService.CreateUser(ctx, mappedUser(idp.Name, authenticated))
		if err != nil {
			return nil, err
		}
		userBase = &createResp.UserBase
	}
	if userBase == nil {
		userBase = &v1.UserBase{
			ObjectMeta: metav1.ObjectMeta{
				Name: authenticated.GetUsername(),
			},
			Avatar: authenticated.GetAvatar(),
			Email:  authenticated.GetEmail(),
		}
	}
	go func() {
		err := a.Store.UserRepository().FlushLastLoginTime(ctx, userBase.Name)
		if err != nil {
			log.Errorf("Failed to flush user [%s] last login time: %v", userBase.Name, err)
		}
	}()

	accessToken, err := a.TokenService.IssueTo(&token.IssueRequest{
		User:      *userBase,
		Claims:    token.Claims{TokenType: token.AccessToken},
		ExpiresIn: a.cfg.AuthenticationOptions.OAuthOptions.AccessTokenMaxAge,
	})
	if err != nil {
		return nil, err
	}
	refreshToken, err := a.TokenService.IssueTo(&token.IssueRequest{
		User:      *userBase,
		Claims:    token.Claims{TokenType: token.RefreshToken},
		ExpiresIn: a.cfg.AuthenticationOptions.OAuthOptions.AccessTokenMaxAge + a.cfg.AuthenticationOptions.OAuthOptions.AccessTokenInactivityTimeout,
	})
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
		ExpiresIn:    int(a.cfg.AuthenticationOptions.OAuthOptions.AccessTokenMaxAge.Seconds()),
	}, nil
}

func (a *authenticationServiceImpl) RefreshToken(_ context.Context, refreshToken string) (*v1.RefreshTokenResponse, error) {
	claim, err := a.TokenService.Verify(refreshToken)
	if err != nil {
		if errors.IsCode(err, code.ErrExpired) {
			return nil, errors.WithCode(code.ErrExpired, jwt.ErrTokenExpired.Error())
		}
		return nil, err
	}
	if claim.TokenType == token.RefreshToken {
		accessToken, err := a.TokenService.IssueTo(&token.IssueRequest{
			User:      claim.User,
			Claims:    token.Claims{TokenType: token.AccessToken},
			ExpiresIn: a.cfg.AuthenticationOptions.OAuthOptions.AccessTokenMaxAge,
		})
		if err != nil {
			return nil, err
		}
		return &v1.RefreshTokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil
	}
	return nil, errors.WithCode(code.ErrInvalidRefreshToken,
		"Incorrect refresh token format %s, expected in refresh format", claim.TokenType)
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
	go func() {
		err := l.store.UserRepository().FlushLastLoginTime(ctx, user.Name)
		if err != nil {
			log.Errorf("Failed to flush user [%s] last login time: %v", user.Name, err)
		}
	}()

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
