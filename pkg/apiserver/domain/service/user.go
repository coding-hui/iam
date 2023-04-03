package service

import (
	"context"

	"github.com/wecoding/iam/pkg/apiserver/domain/model"
	"github.com/wecoding/iam/pkg/apiserver/infrastructure/datastore"
	assembler "github.com/wecoding/iam/pkg/apiserver/interfaces/api/assembler/v1"
	apisv1 "github.com/wecoding/iam/pkg/apiserver/interfaces/api/dto/v1"
)

// UserService User manage api
type UserService interface {
	ListUsers(ctx context.Context, page, pageSize int, listOptions apisv1.ListUserOptions) (*apisv1.ListUserResponse, error)
	Init(ctx context.Context) error
}

type userServiceImpl struct {
	Store datastore.DataStore `inject:"datastore"`
}

// NewUserService new User service
func NewUserService() UserService {
	return &userServiceImpl{}
}

func (u *userServiceImpl) Init(ctx context.Context) error {
	return nil
}

func (u *userServiceImpl) ListUsers(ctx context.Context, page, pageSize int, listOptions apisv1.ListUserOptions) (*apisv1.ListUserResponse, error) {
	user := &model.User{}
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
	fo := datastore.FilterOptions{Queries: queries}

	var userList []*apisv1.DetailUserResponse
	users, err := u.Store.List(ctx, user, &datastore.ListOptions{
		Page:          page,
		PageSize:      pageSize,
		SortBy:        []datastore.SortOption{{Key: "createTime", Order: datastore.SortOrderDescending}},
		FilterOptions: fo,
	})
	if err != nil {
		return nil, err
	}
	for _, v := range users {
		user, ok := v.(*model.User)
		if ok {
			userList = append(userList, assembler.ConvertUserModel(user))
		}
	}
	count, err := u.Store.Count(ctx, user, &fo)
	if err != nil {
		return nil, err
	}

	return &apisv1.ListUserResponse{
		Users: userList,
		Total: count,
	}, nil
}
