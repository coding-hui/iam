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

//	@Tags			System
//	@Summary		Ping
//	@Description	Check service is running
//	@Success		200	{object}	api.Response	"pong"
//	@Router			/ping [get]
//
// ping.
func (p *ping) ping(c *gin.Context) {
	api.OkWithData("pong", c)
}
