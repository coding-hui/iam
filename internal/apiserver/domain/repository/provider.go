// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"

	metav1 "github.com/coding-hui/common/meta/v1"
)

// ProviderRepository defines the provider repository interface.
type ProviderRepository interface {
	Create(ctx context.Context, p *model.Provider, opts metav1.CreateOptions) error
	CreateBatch(ctx context.Context, p []*model.Provider, opts metav1.CreateOptions) error
	Update(ctx context.Context, p *model.Provider, opts metav1.UpdateOptions) error
	Delete(ctx context.Context, p string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, names []string, opts metav1.DeleteOptions) error
	GetByName(ctx context.Context, name string, opts metav1.GetOptions) (*model.Provider, error)
	GetByInstanceId(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.Provider, error)
	List(ctx context.Context, opts metav1.ListOptions) ([]model.Provider, error)
	Count(ctx context.Context, opts metav1.ListOptions) (int64, error)
}
