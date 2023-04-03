package v1

import (
	"github.com/wecoding/iam/pkg/apiserver/domain/model"
	apisv1 "github.com/wecoding/iam/pkg/apiserver/interfaces/api/dto/v1"
)

func convertBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func ConvertUserModel(user *model.User) *apisv1.DetailUserResponse {
	return &apisv1.DetailUserResponse{
		UserBase: *convertUserBase(user),
	}
}

func convertUserBase(user *model.User) *apisv1.UserBase {
	return &apisv1.UserBase{
		Name:          user.Name,
		Alias:         user.Alias,
		Email:         user.Email,
		CreateTime:    user.CreateTime,
		LastLoginTime: user.LastLoginTime,
		Disabled:      user.Disabled,
	}
}
