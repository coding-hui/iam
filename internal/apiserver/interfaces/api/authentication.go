// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/pkg/api"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/internal/pkg/middleware"
	"github.com/coding-hui/iam/internal/pkg/middleware/auth"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
)

// autoAuthCheck authentication strategy which can automatically choose between Basic and Bearer
var autoAuthCheck middleware.AuthStrategy

type authentication struct {
	UserService           service.UserService           `inject:""`
	AuthenticationService service.AuthenticationService `inject:""`

	cfg config.Config
}

// NewAuthentication is the  of authentication.
func NewAuthentication(c config.Config) Interface {
	return &authentication{cfg: c}
}

func (a *authentication) RegisterApiGroup(g *gin.Engine) {
	autoAuthCheck = auth.NewAutoStrategy(
		newBasicAuth(a.AuthenticationService).(auth.BasicStrategy),
		newJWTAuth(a.cfg.JwtOptions.Key).(auth.JWTStrategy),
	)
	apiv1 := g.Group(versionPrefix)
	{
		apiv1.POST("/login", a.authenticate)
		apiv1.GET("/auth/refresh-token", a.refreshToken)
		apiv1.GET("/auth/user-info", autoAuthCheck.AuthFunc(), a.userInfo)
	}
}

func newBasicAuth(authentication service.AuthenticationService) middleware.AuthStrategy {
	return auth.NewBasicStrategy(func(username string, password string) (*v1.AuthenticateResponse, error) {
		login := v1.AuthenticateRequest{
			Username: username,
			Password: password,
		}
		response, err := authentication.Authenticate(context.TODO(), login)
		if err != nil {
			return nil, err
		}

		return response, nil
	})
}

func newJWTAuth(signedKey string) middleware.AuthStrategy {
	return auth.NewJWTStrategy(signedKey)
}

func permissionCheckFunc(r string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, ok := c.Request.Context().Value(&v1.CtxKeyUserType).(string)
		if ok && userType == v1.PlatformAdmin.String() {
			c.Next()
			return
		}

		sub, ok := c.Request.Context().Value(&v1.CtxKeyUserInstanceId).(string)
		if !ok {
			api.FailWithErrCode(errors.WithCode(code.ErrPermissionDenied, "Failed to obtain the current user role"), c)
			c.Abort()
			return
		}
		url := c.Request.URL.Path
		obj := fmt.Sprintf("%s:%s", r, url)
		act := strings.ToLower(c.Request.Method)

		e := repository.Client().CasbinRepository().SyncedEnforcer()
		pass, err := e.Enforce(sub, obj, act)
		if err != nil {
			api.FailWithErrCode(err, c)
			c.Abort()
			return
		}
		if !pass {
			api.FailWithErrCode(errors.WithCode(
				code.ErrPermissionDenied,
				"Permission denied. obj: [%s] sub: [%s] act: [%s]", obj, sub, act),
				c)
			c.Abort()
			return
		}
		log.Infof("Permission verification. path: %s", c.Request.URL.Path)

		c.Next()
	}
}

//	@Tags			Authentication
//	@Summary		LoginSystem
//	@Description	Login by user account and password
//	@Accept			application/json
//	@Product		application/json
//	@Param			data	body		v1.AuthenticateRequest						true	"login request"
//	@Success		200		{object}	api.Response{data=v1.AuthenticateResponse}	"token info"
//	@Router			/api/v1/login [post]
//
// authenticate login by user.
func (a *authentication) authenticate(c *gin.Context) {
	var login v1.AuthenticateRequest
	var err error

	// support header and body both
	if c.Request.Header.Get("Authorization") != "" {
		login, err = parseWithHeader(c)
	}
	if c.Request.Header.Get("Authorization") == "" || err != nil {
		login, err = parseWithBody(c)
	}

	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	var resp *v1.AuthenticateResponse

	if login.Username != "" && login.Password != "" {
		resp, err = a.AuthenticationService.Authenticate(c.Request.Context(), login)
		if err != nil {
			api.FailWithErrCode(err, c)
			return
		}

		api.OkWithData(resp, c)
		return
	}

	resp, err = a.AuthenticationService.AuthenticateByProvider(c.Request.Context(), login)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(resp, c)
}

func parseWithHeader(c *gin.Context) (v1.AuthenticateRequest, error) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		return v1.AuthenticateRequest{}, errors.WithCode(code.ErrPasswordIncorrect, "")
	}

	return v1.AuthenticateRequest{
		Username: username,
		Password: password,
	}, nil
}

func parseWithBody(c *gin.Context) (v1.AuthenticateRequest, error) {
	var login v1.AuthenticateRequest
	if err := c.ShouldBindJSON(&login); err != nil {
		return v1.AuthenticateRequest{}, errors.WithCode(code.ErrPasswordIncorrect, "")
	}

	return login, nil
}

//	@Tags			Authentication
//	@Summary		RefreshToken
//	@Description	RefreshToken
//	@Accept			application/json
//	@Product		application/json
//	@Param			RefreshToken	header		string										true	"refresh token"
//	@Success		200				{object}	api.Response{data=v1.RefreshTokenResponse}	"token info"
//	@Router			/api/v1/auth/refresh-token [get]
//
// refreshToken refresh token.
func (a *authentication) refreshToken(c *gin.Context) {
	base, err := a.AuthenticationService.RefreshToken(
		c.Request.Context(),
		c.GetHeader("RefreshToken"),
	)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(base, c)
}

func (a *authentication) userInfo(c *gin.Context) {
	instanceId, ok := c.Request.Context().Value(&v1.CtxKeyUserInstanceId).(string)
	if !ok {
		api.FailWithErrCode(
			errors.WithCode(code.ErrMissingHeader, "The Authorization header was empty"),
			c,
		)
		return
	}
	user, err := a.UserService.GetUserByInstanceId(c.Request.Context(), instanceId, metav1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}
	resp, err := a.UserService.DetailUser(c.Request.Context(), user)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(resp, c)
}
