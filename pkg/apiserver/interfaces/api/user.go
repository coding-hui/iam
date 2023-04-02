package api

import (
	"github.com/gin-gonic/gin"

	"github.com/wecoding/iam/pkg/api"
	"github.com/wecoding/iam/pkg/apiserver/domain/service"
	apisv1 "github.com/wecoding/iam/pkg/apiserver/interfaces/api/dto/v1"
)

type user struct {
	UserService service.UserService `inject:""`
}

// NewUser is the of user
func NewUser() Interface {
	return &user{}
}

func (u *user) GetApiGroup() InitApiGroup {
	v1 := InitApiGroup{
		BaseUrl: versionPrefix + "/users",
		Apis: []InitApi{
			{
				Method:  "GET",
				Handler: u.listUser,
			},
		},
	}

	return v1
}

// listUser
// @Tags 用户管理
// @Summary 分页获取用户列表
// @Description 分页获取用户列表
// @Param name query string false "名称"
// @Param alias query string false "别名"
// @Param email query string false "邮箱"
// @Param page_size query int false "页条数"
// @Param page query int false "页码"
// @Success   200   {object}  api.Response{data=api.PageResponse{list=[]apisv1.DetailUserResponse}} "{"code": "000", "data": [...]} "分页获取用户列表,返回包括列表,总数,页码,每页数量"
// @Router /api/v1/users [get]
// @Security Bearer
func (u *user) listUser(c *gin.Context) {
	api.OkWithPage(apisv1.ListUserResponse{}, 0, 10, 1, c)
}
