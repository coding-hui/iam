package service

import (
	"context"

	"github.com/wecoding/iam/pkg/apiserver/infrastructure/datastore"
	apisv1 "github.com/wecoding/iam/pkg/apiserver/interfaces/api/dto/v1"
)

// UserService User manage api
type UserService interface {
	ListUsers(ctx context.Context, page, pageSize int, listOptions apisv1.ListUserOptions) (*apisv1.ListUserResponse, error)
	Init(ctx context.Context) error
}

type userServiceImpl struct {
}

// NewUserService new User service
func NewUserService() UserService {
	return &userServiceImpl{}
}

func (u *userServiceImpl) Init(ctx context.Context) error {
	return nil
}

func (u *userServiceImpl) ListUsers(ctx context.Context, page, pageSize int, listOptions apisv1.ListUserOptions) (*apisv1.ListUserResponse, error) {
	var queries []datastore.FuzzyQueryOption
	if listOptions.Name != "" {
		queries = append(queries, datastore.FuzzyQueryOption{Key: "name", Query: listOptions.Name})
	}
	if listOptions.Email != "" {
		queries = append(queries, datastore.FuzzyQueryOption{Key: "email", Query: listOptions.Email})
	}
	if listOptions.Alias != "" {
		queries = append(queries, datastore.FuzzyQueryOption{Key: "alias", Query: listOptions.Alias})
	}

	var userList []*apisv1.DetailUserResponse

	return &apisv1.ListUserResponse{
		Users: userList,
		Total: 0,
	}, nil
}
