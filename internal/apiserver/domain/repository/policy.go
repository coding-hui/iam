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

// PolicyRepository defines the policy repository interface.
type PolicyRepository interface {
	Create(ctx context.Context, p *model.Policy, opts metav1.CreateOptions) error
	CreateBatch(ctx context.Context, p []*model.Policy, opts metav1.CreateOptions) error
	Update(ctx context.Context, p *model.Policy, opts metav1.UpdateOptions) error
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, names []string, opts metav1.DeleteOptions) error
	GetByName(ctx context.Context, name string, opts metav1.GetOptions) (*model.Policy, error)
	GetByInstanceId(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.Policy, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.PolicyList, error)
	CountStatementByResource(ctx context.Context, resource ...string) (int64, error)
}
