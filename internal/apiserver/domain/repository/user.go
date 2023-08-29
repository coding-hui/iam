// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"

	metav1 "github.com/coding-hui/common/meta/v1"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

// UserRepository defines the user repository interface.
type UserRepository interface {
	Create(ctx context.Context, user *model.User, opts metav1.CreateOptions) (*model.User, error)
	Update(ctx context.Context, user *model.User, opts metav1.UpdateOptions) error
	DeleteByInstanceId(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, usernames []string, opts metav1.DeleteOptions) error
	GetByName(ctx context.Context, username string, opts metav1.GetOptions) (*model.User, error)
	GetByInstanceId(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.User, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error)
}
