// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/coding-hui/iam/internal/authzserver/authorization"
	"github.com/coding-hui/iam/internal/authzserver/config"
	"github.com/coding-hui/iam/pkg/api"
	authzv1 "github.com/coding-hui/iam/pkg/api/authzserver/v1"
	"github.com/coding-hui/iam/pkg/code"

	"github.com/coding-hui/common/errors"
)

type auth struct {
	Authorizer authorization.Authorization `inject:"authorizer"`
	cfg        config.Config
}

// NewAuth authorization api.
func NewAuth(c config.Config) Interface {
	return &auth{cfg: c}
}

func (a *auth) RegisterApiGroup(g *gin.Engine) {
	v1 := g.Group(versionPrefix)
	{
		v1.POST("/authz", a.authorize)
	}
}

func (a *auth) authorize(c *gin.Context) {
	req := &authzv1.Request{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}
	resp := a.Authorizer.Authorize(req)

	api.OkWithData(resp, c)
}
