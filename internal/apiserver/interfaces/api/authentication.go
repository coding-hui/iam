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
	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/pkg/api"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"

	"github.com/coding-hui/common/errors"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
)

type authentication struct {
	UserService           service.UserService           `inject:""`
	AuthenticationService service.AuthenticationService `inject:""`

	cfg config.Config
}

// NewAuthentication is the  of authentication
func NewAuthentication(c config.Config) Interface {
	return &authentication{cfg: c}
}

func (a *authentication) GetApiGroup() InitApiGroup {
	return InitApiGroup{
		BaseUrl: versionPrefix,
		Apis: []InitApi{
			{
				Method:  POST,
				Path:    "/login",
				Handler: a.authenticate,
			},
			{
				Method:  GET,
				Path:    "/auth/refresh-token",
				Handler: a.refreshToken,
			},
			{
				Method:  GET,
				Path:    "/auth/user-info",
				Filters: gin.HandlersChain{authCheckFilter},
				Handler: a.userInfo,
			},
		},
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

	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), &v1alpha1.CtxKeyUserName, token.Username))

	c.Next()
}

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

func (a *authentication) refreshToken(c *gin.Context) {
	base, err := a.AuthenticationService.RefreshToken(c.Request.Context(), c.GetHeader("RefreshToken"))
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(base, c)
}

func (a *authentication) userInfo(c *gin.Context) {
	userName, ok := c.Request.Context().Value(&v1alpha1.CtxKeyUserName).(string)
	if !ok {
		api.FailWithErrCode(errors.WithCode(code.ErrMissingHeader, "The Authorization header was empty"), c)
		return
	}
	user, err := a.UserService.Get(c.Request.Context(), userName, metav1alpha1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(user, c)
}
