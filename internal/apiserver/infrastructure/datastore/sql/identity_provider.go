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

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
)

type identityProviderRepositoryImpl struct {
	client *Client
}

// newProviderRepository new IdentityProvider Repository.
func newProviderRepository(client *Client) repository.IdentityProviderRepository {
	return &identityProviderRepositoryImpl{client}
}

func (p *identityProviderRepositoryImpl) Create(ctx context.Context, policy *model.IdentityProvider, _ metav1.CreateOptions) error {
	if err := p.client.WithCtx(ctx).Create(&policy).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return datastore.ErrRecordExist
		}
		return err
	}

	return nil
}

func (p *identityProviderRepositoryImpl) CreateBatch(
	ctx context.Context,
	providers []*model.IdentityProvider,
	_ metav1.CreateOptions,
) error {
	if err := p.client.WithCtx(ctx).CreateInBatches(&providers, 500).Error; err != nil {
		return err
	}

	return nil
}

func (p *identityProviderRepositoryImpl) Update(ctx context.Context, provider *model.IdentityProvider, _ metav1.UpdateOptions) error {
	err := p.client.WithCtx(ctx).Model(provider).Save(provider).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return datastore.ErrRecordNotExist
		}

		return err
	}

	return nil
}

func (p *identityProviderRepositoryImpl) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	db := p.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	err := db.Where("instance_id = ?", name).Delete(&model.IdentityProvider{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return datastore.ErrRecordNotExist
		}

		return err
	}

	return nil
}

func (p *identityProviderRepositoryImpl) DeleteCollection(ctx context.Context, names []string, opts metav1.DeleteOptions) error {
	db := p.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}

	return db.Where("name in (?)", names).Delete(&model.IdentityProvider{}).Error
}

func (p *identityProviderRepositoryImpl) GetByName(
	ctx context.Context,
	name string,
	_ metav1.GetOptions,
) (provider *model.IdentityProvider, err error) {
	err = p.client.WithCtx(ctx).Where("name = ?", name).First(&provider).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, datastore.ErrRecordNotExist
		}

		return nil, err
	}

	return provider, err
}

func (p *identityProviderRepositoryImpl) GetByInstanceId(ctx context.Context, instanceId string, _ metav1.GetOptions) (provider *model.IdentityProvider, err error) {
	err = p.client.WithCtx(ctx).Where("instance_id = ?", instanceId).First(&provider).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, datastore.ErrRecordNotExist
		}

		return nil, err
	}

	return provider, err
}

func (p *identityProviderRepositoryImpl) GetByInstanceIdOrName(ctx context.Context, instanceIdOrName string, opts metav1.GetOptions) (idp *model.IdentityProvider, err error) {
	err = p.client.WithCtx(ctx).
		Where("instance_id = ? or name = ?", instanceIdOrName, instanceIdOrName).
		First(&idp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, datastore.ErrRecordNotExist
		}

		return nil, err
	}

	return idp, err
}

func (p *identityProviderRepositoryImpl) List(ctx context.Context, opts metav1.ListOptions) ([]model.IdentityProvider, error) {
	var providers []model.IdentityProvider
	err := p.client.WithCtx(ctx).Model(model.IdentityProvider{}).
		Scopes(
			makeCondition(opts),
			paginate(opts),
		).
		Order("id desc").
		Find(&providers).Error
	if err != nil {
		return nil, datastore.NewDBError(err, "failed to list providers")
	}

	return providers, err
}

func (p *identityProviderRepositoryImpl) Count(ctx context.Context, opts metav1.ListOptions) (int64, error) {
	var totalCount int64
	err := p.client.WithCtx(ctx).Model(model.IdentityProvider{}).
		Scopes(
			makeCondition(opts),
			paginate(opts),
		).
		Order("id desc").
		Count(&totalCount).Error
	if err != nil {
		return 0, datastore.NewDBError(err, "failed to count providers")
	}

	return totalCount, err
}
