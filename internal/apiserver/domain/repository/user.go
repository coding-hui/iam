// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"

	metav1 "github.com/coding-hui/common/meta/v1"
)

// UserRepository defines the user repository interface.
type UserRepository interface {
	Create(ctx context.Context, user *model.User, opts metav1.CreateOptions) (*model.User, error)
	Update(ctx context.Context, user *model.User, opts metav1.UpdateOptions) error
	DeleteByInstanceId(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error
	BatchDelete(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error
	GetByName(ctx context.Context, username string, opts metav1.GetOptions) (*model.User, error)
	GetByInstanceId(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.User, error)
	GetByNameOrInstanceId(ctx context.Context, nameOrId string, opts metav1.GetOptions) (*model.User, error)
	FlushLastLoginTime(ctx context.Context, nameOrId string) error
	GetByExternalId(ctx context.Context, openId, externalId string, opts metav1.GetOptions) (*model.User, error)
	CreateExternalUser(ctx context.Context, externalUser *model.UserExternal) error
	DeleteExternalUser(ctx context.Context, userId, externalUid, idp string) error
	DeleteExternalUserByProvider(ctx context.Context, userId, provider string) error
	List(ctx context.Context, opts v1.ListUserOptions) ([]model.User, error)
	Count(ctx context.Context, opts v1.ListUserOptions) (int64, error)
}
