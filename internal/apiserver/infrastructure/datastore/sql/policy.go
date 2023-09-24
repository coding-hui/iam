// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"context"

	"gorm.io/gorm"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/apiserver/infrastructure/datastore"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1"
	"github.com/coding-hui/iam/internal/pkg/code"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
)

type policyRepositoryImpl struct {
	client *Client
}

// newPolicyRepository new Policy Repository.
func newPolicyRepository(client *Client) repository.PolicyRepository {
	return &policyRepositoryImpl{client}
}

// Create creates a new policy.
func (p *policyRepositoryImpl) Create(ctx context.Context, policy *model.Policy, _ metav1.CreateOptions) error {
	if oldPolicy, _ := p.GetByName(ctx, policy.Name, metav1.GetOptions{}); oldPolicy != nil {
		return errors.WithCode(code.ErrPolicyAlreadyExist, "Policy %s already exist", policy.Name)
	}
	if err := p.client.WithCtx(ctx).Create(&policy).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.WithCode(code.ErrPolicyAlreadyExist, err.Error())
		}
		return err
	}

	return nil
}

// CreateBatch creates a new policy.
func (p *policyRepositoryImpl) CreateBatch(ctx context.Context, policy []*model.Policy, _ metav1.CreateOptions) error {
	if err := p.client.WithCtx(ctx).CreateInBatches(&policy, 500).Error; err != nil {
		return err
	}

	return nil
}

// Update updates an policy information.
func (p *policyRepositoryImpl) Update(ctx context.Context, policy *model.Policy, _ metav1.UpdateOptions) error {
	err := p.client.WithCtx(ctx).Model(policy).Association("Statements").Replace(policy.Statements)
	if err != nil {
		return err
	}
	err = p.client.WithCtx(ctx).Model(policy).Save(policy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrPolicyNotFound, err.Error())
		}

		return err
	}

	return nil
}

// Delete deletes the policy by the user identifier.
func (p *policyRepositoryImpl) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	db := p.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	err := db.Where("name = ?", name).Delete(&model.Policy{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(code.ErrPolicyNotFound, err.Error())
		}

		return err
	}

	return nil
}

// DeleteCollection batch deletes the policies.
func (p *policyRepositoryImpl) DeleteCollection(ctx context.Context, names []string, opts metav1.DeleteOptions) error {
	db := p.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}

	return db.Where("name in (?)", names).Delete(&model.User{}).Error
}

// GetByName get policy.
func (p *policyRepositoryImpl) GetByName(ctx context.Context, name string, _ metav1.GetOptions) (policy *model.Policy, err error) {
	err = p.client.WithCtx(ctx).Preload("Statements").Where("name = ?", name).First(&policy).Error
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
	_ metav1.GetOptions,
) (policy *model.Policy, err error) {
	err = p.client.WithCtx(ctx).Debug().Preload("Statements").Where("instance_id = ?", instanceId).First(&policy).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrPolicyNotFound, err.Error())
		}

		return nil, err
	}

	return policy, err
}

// List list policies.
func (p *policyRepositoryImpl) List(ctx context.Context, opts metav1.ListOptions) (*v1.PolicyList, error) {
	var policies []*model.Policy
	var totalCount int64
	err := p.client.WithCtx(ctx).Model(model.Policy{}).
		Scopes(
			makeCondition(opts),
			paginate(opts),
		).
		Order("id desc").
		Find(&policies).Offset(-1).Limit(-1).
		Count(&totalCount).Error
	if err != nil {
		return nil, datastore.NewDBError(err, "failed to list policies")
	}
	items := make([]*v1.PolicyBase, 0, len(policies))
	for _, policy := range policies {
		items = append(items, assembler.ConvertPolicyModelToBase(policy))
	}

	return &v1.PolicyList{
		ListMeta: metav1.ListMeta{TotalCount: totalCount},
		Items:    items,
	}, err
}

// CountStatementByResource get statement count by resource.
func (p *policyRepositoryImpl) CountStatementByResource(ctx context.Context, resource ...string) (int64, error) {
	var count int64
	db := p.client.WithCtx(ctx).Model(&model.Statement{})
	db.Where("resource in ?", resource)
	db.Where("policy_id is not null")
	err := db.Count(&count).Error
	if err != nil {
		return -1, err
	}

	return count, nil
}
