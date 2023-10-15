// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"

	metav1 "github.com/coding-hui/common/meta/v1"
)

// ApplicationRepository defines the app repository interface.
type ApplicationRepository interface {
	Create(ctx context.Context, app *model.Application, opts metav1.CreateOptions) error
	CreateBatch(ctx context.Context, app []*model.Application, opts metav1.CreateOptions) error
	Update(ctx context.Context, app *model.Application, opts metav1.UpdateOptions) error
	Delete(ctx context.Context, app string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, apps []string, opts metav1.DeleteOptions) error
	GetByName(ctx context.Context, name string, opts metav1.GetOptions) (*model.Application, error)
	GetByInstanceId(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.Application, error)
	List(ctx context.Context, opts metav1.ListOptions) ([]model.Application, error)
	Count(ctx context.Context, opts metav1.ListOptions) (int64, error)
}
