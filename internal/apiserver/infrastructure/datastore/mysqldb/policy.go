// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysqldb

import (
	"context"

	"gorm.io/gorm"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1alpha1"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/internal/pkg/utils/gormutil"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"

	"github.com/coding-hui/common/errors"
	"github.com/coding-hui/common/fields"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
)

type policyRepositoryImpl struct {
	db *gorm.DB
}

// newPolicyRepository new Policy Repository.
func newPolicyRepository(db *gorm.DB) repository.PolicyRepository {
	return &policyRepositoryImpl{db}
}

// Create creates a new policy.
func (p *policyRepositoryImpl) Create(ctx context.Context, policy *model.Policy, _ metav1alpha1.CreateOptions) error {
	if oldPolicy, _ := p.Get(ctx, policy.Name, metav1alpha1.GetOptions{}); oldPolicy != nil {
		return errors.WithCode(code.ErrPolicyAlreadyExist, "Policy %s already exist", policy.Name)
	}
	if err := p.db.WithContext(ctx).Create(&policy).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.WithCode(code.ErrPolicyAlreadyExist, err.Error())
		}
		return err
	}

	return nil
}

// Update updates an policy information.
func (p *policyRepositoryImpl) Update(ctx context.Context, policy *model.Policy, _ metav1alpha1.UpdateOptions) error {
	if err := p.db.WithContext(ctx).Save(policy).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrPolicyNotFound, err.Error())
		}

		return err
	}

	return nil
}

// Delete deletes the policy by the user identifier.
func (p *policyRepositoryImpl) Delete(ctx context.Context, name string, opts metav1alpha1.DeleteOptions) error {
	if opts.Unscoped {
		p.db = p.db.Unscoped()
	}
	err := p.db.WithContext(ctx).Where("name = ?", name).Delete(&model.Policy{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrPolicyNotFound, err.Error())
		}

		return err
	}

	return nil
}

// DeleteCollection batch deletes the policies.
func (p *policyRepositoryImpl) DeleteCollection(ctx context.Context, names []string, opts metav1alpha1.DeleteOptions) error {
	if opts.Unscoped {
		p.db = p.db.Unscoped()
	}

	return p.db.WithContext(ctx).Where("name in (?)", names).Delete(&model.User{}).Error
}

// Get get policy.
func (p *policyRepositoryImpl) Get(ctx context.Context, name string, _ metav1alpha1.GetOptions) (policy *model.Policy, err error) {
	err = p.db.WithContext(ctx).Where("name = ?", name).First(&policy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrPolicyNotFound, err.Error())
		}

		return nil, err
	}

	return policy, err
}

// List list policies.
func (p *policyRepositoryImpl) List(ctx context.Context, opts metav1alpha1.ListOptions) (*v1alpha1.PolicyList, error) {
	var policies []*model.Policy
	var totalCount int64

	ol := gormutil.Unpointer(opts.Offset, opts.Limit)

	db := p.db.WithContext(ctx).Model(model.Policy{})
	selector, _ := fields.ParseSelector(opts.FieldSelector)
	username, _ := selector.RequiresExactMatch("name")
	if username != "" {
		db.Where("name like ?", "%"+username+"%")
	}
	db.Offset(ol.Offset).
		Limit(ol.Limit).
		Order("id desc").
		Find(&policies).
		Offset(-1).
		Limit(-1).
		Count(&totalCount)
	items := make([]*v1alpha1.PolicyBase, 0, len(policies))
	for _, policy := range policies {
		items = append(items, assembler.ConvertPolicyModelToBase(policy))
	}

	return &v1alpha1.PolicyList{
		ListMeta: metav1alpha1.ListMeta{TotalCount: totalCount},
		Items:    items,
	}, db.Error
}
