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

// ResourceRepository defines the resource repository interface.
type ResourceRepository interface {
	BatchCreate(ctx context.Context, resources []*model.Resource, opts metav1alpha1.CreateOptions) error
	Create(ctx context.Context, resource *model.Resource, opts metav1alpha1.CreateOptions) error
	Update(ctx context.Context, resource *model.Resource, opts metav1alpha1.UpdateOptions) error
	Delete(ctx context.Context, name string, opts metav1alpha1.DeleteOptions) error
	DeleteCollection(ctx context.Context, names []string, opts metav1alpha1.DeleteOptions) error
	Get(ctx context.Context, name string, opts metav1alpha1.GetOptions) (*model.Resource, error)
	List(ctx context.Context, opts metav1alpha1.ListOptions) (*v1alpha1.ResourceList, error)
}
