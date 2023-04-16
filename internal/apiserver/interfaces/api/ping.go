// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/coding-hui/iam/internal/pkg/api"
	"github.com/gin-gonic/gin"
)

type ping struct{}

// NewPing ping
func NewPing() Interface {
	return &ping{}
}

func (p *ping) GetApiGroup() InitApiGroup {
	return InitApiGroup{
		BaseUrl: "/ping",
		Apis: []InitApi{
			{
				Method:  GET,
				Path:    "",
				Handler: p.ping,
			},
		},
	}
}

// ping
func (p *ping) ping(c *gin.Context) {
	api.OkWithData("pong", c)
}
