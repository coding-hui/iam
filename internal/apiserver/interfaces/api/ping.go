// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/coding-hui/iam/internal/pkg/api"
)

type ping struct{}

// NewPing ping.
func NewPing() Interface {
	return &ping{}
}

func (p *ping) RegisterApiGroup(g *gin.Engine) {
	g.GET("/ping", p.ping)
}

// ping.
//
//	@Tags			System
//	@Summary		check service is running
//	@Description	check service is running
//	@Success		200	{object}	api.Response	"{"code": "000", "data": [...]}
//	@Router			/ping [get]
func (p *ping) ping(c *gin.Context) {
	api.OkWithData("pong", c)
}
