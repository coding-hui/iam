// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"context"
	"strings"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/pkg/api"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
)

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
	v1 := g.Group(versionPrefix)
	{
		v1.POST("/login", a.authenticate)
		v1.GET("/auth/refresh-token", a.refreshToken)
		v1.GET("/auth/user-info", authCheckFilter, a.userInfo)
	}
}

func authCheckFilter(c *gin.Context) {
	var tokenValue string
	tokenHeader := c.Request.Header.Get("Authorization")
	if tokenHeader != "" {
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			api.FailWithErrCode(errors.WithCode(code.ErrMissingHeader, "The Authorization header was empty"), c)
			c.Abort()
			return
		}
		tokenValue = splitted[1]
	}
	if tokenValue == "" {
		api.FailWithErrCode(errors.WithCode(code.ErrMissingHeader, "The Authorization header was empty"), c)
		c.Abort()
		return
	}
	token, err := service.ParseToken(tokenValue)
	if err != nil {
		api.FailWithErrCode(err, c)
		c.Abort()
		return
	}
	if token.GrantType != service.GrantTypeAccess {
		api.FailWithErrCode(errors.WithCode(code.ErrPermissionDenied, "Invalid authorization header"), c)
		c.Abort()
		return
	}

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), &v1alpha1.CtxKeyUserInstanceId, token.UserInstanceId))
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), &v1alpha1.CtxKeyUserType, token.UserType))

	c.Next()
}

func permissionCheckFilter(c *gin.Context) {
	userType, ok := c.Request.Context().Value(&v1alpha1.CtxKeyUserType).(string)
	if ok && userType == v1alpha1.PlatformAdmin.String() {
		c.Next()
		return
	}

	sub, ok := c.Request.Context().Value(&v1alpha1.CtxKeyUserInstanceId).(string)
	if !ok {
		api.FailWithErrCode(errors.WithCode(code.ErrPermissionDenied, "Failed to obtain the current user role"), c)
		c.Abort()
		return
	}
	obj := c.Request.URL.Path
	act := c.Request.Method

	e := repository.Client().CasbinRepository().SyncedEnforcer()
	pass, err := e.Enforce(sub, obj, act)
	if err != nil {
		api.FailWithErrCode(err, c)
		c.Abort()
		return
	}
	if !pass {
		api.FailWithErrCode(errors.WithCode(code.ErrPermissionDenied, "Permission denied. role: %s", sub), c)
		c.Abort()
		return
	}
	log.Infof("Permission verification. path: %s", c.Request.URL.Path)

	c.Next()
}

//	@Tags			Authentication
//	@Summary		LoginSystem
//	@Description	Login by user account and password
//	@Accept			application/json
//	@Product		application/json
//	@Param			data	body		v1alpha1.AuthenticateRequest						true	"login request"
//	@Success		200		{object}	api.Response{data=v1alpha1.AuthenticateResponse}	"token info"
//	@Router			/api/v1/login [post]
//
// authenticate login by user.
func (a *authentication) authenticate(c *gin.Context) {
	var login v1alpha1.AuthenticateRequest
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

	response, err := a.AuthenticationService.Authenticate(c.Request.Context(), login)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(response, c)
}

func parseWithHeader(c *gin.Context) (v1alpha1.AuthenticateRequest, error) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		return v1alpha1.AuthenticateRequest{}, jwt.ErrFailedAuthentication
	}

	return v1alpha1.AuthenticateRequest{
		Username: username,
		Password: password,
	}, nil
}

func parseWithBody(c *gin.Context) (v1alpha1.AuthenticateRequest, error) {
	var login v1alpha1.AuthenticateRequest
	if err := c.ShouldBindJSON(&login); err != nil {
		return v1alpha1.AuthenticateRequest{}, jwt.ErrFailedAuthentication
	}

	return login, nil
}

//	@Tags			Authentication
//	@Summary		RefreshToken
//	@Description	RefreshToken
//	@Accept			application/json
//	@Product		application/json
//	@Param			RefreshToken	header		string												true	"refresh token"
//	@Success		200				{object}	api.Response{data=v1alpha1.RefreshTokenResponse}	"token info"
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
	instanceId, ok := c.Request.Context().Value(&v1alpha1.CtxKeyUserInstanceId).(string)
	if !ok {
		api.FailWithErrCode(
			errors.WithCode(code.ErrMissingHeader, "The Authorization header was empty"),
			c,
		)
		return
	}
	user, err := a.UserService.GetUserByInstanceId(c.Request.Context(), instanceId, metav1alpha1.GetOptions{})
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
