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

// OrganizationRepository defines the organization repository interface.
type OrganizationRepository interface {
	Create(ctx context.Context, o *model.Organization, opts metav1.CreateOptions) error
	Update(ctx context.Context, o *model.Organization, opts metav1.UpdateOptions) error
	BatchUpdate(ctx context.Context, list []*model.Organization, opts metav1.UpdateOptions) error
	DeleteByName(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteByInstanceId(ctx context.Context, uid string, opts metav1.DeleteOptions) error
	GetByName(ctx context.Context, name string, opts metav1.GetOptions) (*model.Organization, error)
	GetByInstanceId(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.Organization, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.OrganizationList, error)
	CountDepartmentByOrg(ctx context.Context, org string, opts metav1.ListOptions) (int64, error)
	ListChildDepartments(ctx context.Context, org string, opts metav1.ListOptions) ([]model.Organization, error)
	ListDepartmentMembers(ctx context.Context, department string, opts metav1.ListOptions) ([]model.DepartmentMember, error)
	CountDepartmentMembers(ctx context.Context, department string, opts metav1.ListOptions) (int64, error)
	AddDepartmentMembers(ctx context.Context, members []*model.DepartmentMember) error
	RemoveDepartmentMembers(ctx context.Context, members []*model.DepartmentMember) error
}
