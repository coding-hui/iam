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
// @Tags System
// @Summary check service is running
// @Description check service is running
// @Success   200   {object}  api.Response "{"code": "000", "data": [...]}
// @Router /ping [get]
func (p *ping) ping(c *gin.Context) {
	api.OkWithData("pong", c)
}
