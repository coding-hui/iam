// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mysqldb

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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
	if oldPolicy, _ := p.GetByName(ctx, policy.Name, metav1alpha1.GetOptions{}); oldPolicy != nil {
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

// CreateBatch creates a new policy.
func (p *policyRepositoryImpl) CreateBatch(ctx context.Context, policy []*model.Policy, _ metav1alpha1.CreateOptions) error {
	if err := p.db.WithContext(ctx).CreateInBatches(&policy, 500).Error; err != nil {
		return err
	}

	return nil
}

// Update updates an policy information.
func (p *policyRepositoryImpl) Update(ctx context.Context, policy *model.Policy, _ metav1alpha1.UpdateOptions) error {
	err := p.db.WithContext(ctx).Model(policy).Association("Statements").Replace(policy.Statements)
	if err != nil {
		return err
	}
	err = p.db.WithContext(ctx).Model(policy).Save(policy).Error
	if err != nil {
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

// GetByName get policy.
func (p *policyRepositoryImpl) GetByName(ctx context.Context, name string, _ metav1alpha1.GetOptions) (policy *model.Policy, err error) {
	err = p.db.WithContext(ctx).Preload("Statements").Where("name = ?", name).First(&policy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrPolicyNotFound, err.Error())
		}

		return nil, err
	}

	return policy, err
}

// GetByInstanceId get policy.
func (p *policyRepositoryImpl) GetByInstanceId(
	ctx context.Context,
	instanceId string,
	_ metav1alpha1.GetOptions,
) (policy *model.Policy, err error) {
	err = p.db.WithContext(ctx).Debug().Preload("Statements").Where("instance_id = ?", instanceId).First(&policy).Error
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
	var clauses []clause.Expression
	selector, _ := fields.ParseSelector(opts.FieldSelector)
	clauses = _applyFieldSelector(clauses, selector)
	db.Offset(ol.Offset).
		Limit(ol.Limit).
		Clauses(clauses...).
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

// CountStatementByResource get statement count by resource.
func (p *policyRepositoryImpl) CountStatementByResource(ctx context.Context, resource ...string) (int64, error) {
	var count int64
	db := p.db.WithContext(ctx).Model(&model.Statement{})
	db.Where("resource in ?", resource)
	db.Where("policy_id is not null")
	err := db.Count(&count).Error
	if err != nil {
		return -1, err
	}

	return count, nil
}
