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

type providerRepositoryImpl struct {
	client *Client
}

// newProviderRepository new Provider Repository.
func newProviderRepository(client *Client) repository.ProviderRepository {
	return &providerRepositoryImpl{client}
}

func (p *providerRepositoryImpl) Create(ctx context.Context, policy *model.Provider, _ metav1.CreateOptions) error {
	if err := p.client.WithCtx(ctx).Create(&policy).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return datastore.ErrRecordExist
		}
		return err
	}

	return nil
}

func (p *providerRepositoryImpl) CreateBatch(ctx context.Context, providers []*model.Provider, _ metav1.CreateOptions) error {
	if err := p.client.WithCtx(ctx).CreateInBatches(&providers, 500).Error; err != nil {
		return err
	}

	return nil
}

func (p *providerRepositoryImpl) Update(ctx context.Context, provider *model.Provider, _ metav1.UpdateOptions) error {
	err := p.client.WithCtx(ctx).Model(provider).Save(provider).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return datastore.ErrRecordNotExist
		}

		return err
	}

	return nil
}

func (p *providerRepositoryImpl) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	db := p.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	err := db.Where("instance_id = ?", name).Delete(&model.Policy{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return datastore.ErrRecordNotExist
		}

		return err
	}

	return nil
}

func (p *providerRepositoryImpl) DeleteCollection(ctx context.Context, names []string, opts metav1.DeleteOptions) error {
	db := p.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}

	return db.Where("name in (?)", names).Delete(&model.Provider{}).Error
}

func (p *providerRepositoryImpl) GetByName(ctx context.Context, name string, _ metav1.GetOptions) (provider *model.Provider, err error) {
	err = p.client.WithCtx(ctx).Where("name = ?", name).First(&provider).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, datastore.ErrRecordNotExist
		}

		return nil, err
	}

	return provider, err
}

func (p *providerRepositoryImpl) GetByInstanceId(
	ctx context.Context,
	instanceId string,
	_ metav1.GetOptions,
) (provider *model.Provider, err error) {
	err = p.client.WithCtx(ctx).Where("instance_id = ?", instanceId).First(&provider).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, datastore.ErrRecordNotExist
		}

		return nil, err
	}

	return provider, err
}

func (p *providerRepositoryImpl) List(ctx context.Context, opts metav1.ListOptions) ([]model.Provider, error) {
	var providers []model.Provider
	err := p.client.WithCtx(ctx).Model(model.Provider{}).
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

func (p *providerRepositoryImpl) Count(ctx context.Context, opts metav1.ListOptions) (int64, error) {
	var totalCount int64
	err := p.client.WithCtx(ctx).Model(model.Provider{}).
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
