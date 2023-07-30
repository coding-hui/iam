// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"

	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
)

// TenantRepository defines the tenant repository interface.
type TenantRepository interface {
	Create(ctx context.Context, tenant *model.Tenant, opts metav1alpha1.CreateOptions) error
	Update(ctx context.Context, tenant *model.Tenant, opts metav1alpha1.UpdateOptions) error
}
