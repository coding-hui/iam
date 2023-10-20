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

type applicationRepositoryImpl struct {
	client *Client
}

// newApplicationRepository new Application Repository.
func newApplicationRepository(client *Client) repository.ApplicationRepository {
	return &applicationRepositoryImpl{client}
}

func (p *applicationRepositoryImpl) Create(ctx context.Context, app *model.Application, _ metav1.CreateOptions) error {
	if err := p.client.WithCtx(ctx).Create(&app).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return datastore.ErrRecordExist
		}
		return err
	}

	return nil
}

func (p *applicationRepositoryImpl) CreateBatch(ctx context.Context, apps []*model.Application, _ metav1.CreateOptions) error {
	if err := p.client.WithCtx(ctx).CreateInBatches(&apps, 500).Error; err != nil {
		return err
	}

	return nil
}

func (p *applicationRepositoryImpl) Update(ctx context.Context, app *model.Application, _ metav1.UpdateOptions) error {
	err := p.client.WithCtx(ctx).Model(app).Save(app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return datastore.ErrRecordNotExist
		}

		return err
	}

	return nil
}

func (p *applicationRepositoryImpl) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	db := p.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	err := db.Where("instance_id = ?", name).Delete(&model.Application{}).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return datastore.ErrRecordNotExist
		}

		return err
	}

	return nil
}

func (p *applicationRepositoryImpl) DeleteCollection(ctx context.Context, names []string, opts metav1.DeleteOptions) error {
	db := p.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}

	return db.Where("name in (?)", names).Delete(&model.Application{}).Error
}

func (p *applicationRepositoryImpl) GetByName(ctx context.Context, name string, _ metav1.GetOptions) (app *model.Application, err error) {
	err = p.client.WithCtx(ctx).
		Preload("IdentityProviders").
		Where("name = ?", name).
		First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, datastore.ErrRecordNotExist
		}

		return nil, err
	}

	return app, err
}

func (p *applicationRepositoryImpl) GetByInstanceId(
	ctx context.Context,
	instanceId string,
	_ metav1.GetOptions,
) (app *model.Application, err error) {
	err = p.client.WithCtx(ctx).
		Preload("IdentityProviders").
		Where("instance_id = ?", instanceId).
		First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, datastore.ErrRecordNotExist
		}

		return nil, err
	}

	return app, err
}

func (p *applicationRepositoryImpl) GetByInstanceIdOrName(
	ctx context.Context,
	idOrName string,
	_ metav1.GetOptions,
) (app *model.Application, err error) {
	err = p.client.WithCtx(ctx).
		Preload("IdentityProviders").
		Where("instance_id = ? or name = ?", idOrName, idOrName).
		First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, datastore.ErrRecordNotExist
		}

		return nil, err
	}

	return app, err
}

func (p *applicationRepositoryImpl) List(ctx context.Context, opts metav1.ListOptions) ([]model.Application, error) {
	var apps []model.Application
	err := p.client.WithCtx(ctx).Model(model.Application{}).
		Scopes(
			makeCondition(opts),
			paginate(opts),
		).
		Order("id desc").
		Find(&apps).Error
	if err != nil {
		return nil, datastore.NewDBError(err, "failed to list apps")
	}

	return apps, err
}

func (p *applicationRepositoryImpl) Count(ctx context.Context, opts metav1.ListOptions) (int64, error) {
	var totalCount int64
	err := p.client.WithCtx(ctx).Model(model.Application{}).
		Scopes(
			makeCondition(opts),
			paginate(opts),
		).
		Order("id desc").
		Count(&totalCount).Error
	if err != nil {
		return 0, datastore.NewDBError(err, "failed to count apps")
	}

	return totalCount, err
}
