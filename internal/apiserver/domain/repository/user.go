// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"

	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"
)

// UserRepository defines the user repository interface.
type UserRepository interface {
	Create(ctx context.Context, user *model.User, opts metav1alpha1.CreateOptions) error
	Update(ctx context.Context, user *model.User, opts metav1alpha1.UpdateOptions) error
	Delete(ctx context.Context, username string, opts metav1alpha1.DeleteOptions) error
	DeleteCollection(ctx context.Context, usernames []string, opts metav1alpha1.DeleteOptions) error
	Get(ctx context.Context, username string, opts metav1alpha1.GetOptions) (*model.User, error)
	List(ctx context.Context, opts metav1alpha1.ListOptions) (*v1alpha1.UserList, error)
}
