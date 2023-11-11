// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1"
	"github.com/coding-hui/iam/internal/pkg/code"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
	"github.com/coding-hui/common/util/idutil"
)

// ApplicationService Application manage api.
type ApplicationService interface {
	CreateApplication(ctx context.Context, req v1.CreateApplicationRequest) error
	UpdateApplication(ctx context.Context, app string, req v1.UpdateApplicationRequest) error
	DeleteApplication(ctx context.Context, app string, opts metav1.DeleteOptions) error
	GetApplication(ctx context.Context, app string, opts metav1.GetOptions) (*v1.DetailApplicationResponse, error)
	ListApplications(ctx context.Context, opts metav1.ListOptions) (*v1.ApplicationList, error)
	Init(ctx context.Context) error
}

type applicationServiceImpl struct {
	cfg   config.Config
	Store repository.Factory `inject:"repository"`
}

// NewApplicationService new Application service.
func NewApplicationService(c config.Config) ApplicationService {
	return &applicationServiceImpl{cfg: c}
}

// Init initialize default app data.
func (a *applicationServiceImpl) Init(ctx context.Context) error {
	old, err := a.Store.ApplicationRepository().GetByName(ctx, model.DefaultApplication, metav1.GetOptions{})
	if err != nil && !errors.IsCode(err, code.ErrRecordNotExist) {
		return err
	}
	if old != nil {
		return nil
	}
	createReq := v1.CreateApplicationRequest{
		Name:        model.DefaultApplication,
		DisplayName: "Built-in Application",
		HomepageUrl: "http://iam.wecoding.top",
		Logo:        "",
		Description: "System Build-in Application",
	}
	err = a.CreateApplication(ctx, createReq)
	if err != nil {
		return errors.WithMessagef(err, "Failed to initialize system built-in application")
	}
	log.Info("initialize system built-in application done")

	return nil
}

func (a *applicationServiceImpl) CreateApplication(ctx context.Context, req v1.CreateApplicationRequest) error {
	app := &model.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
		Status:      req.Status,
		Owner:       req.Owner,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Logo:        req.Logo,
		HomepageUrl: req.HomepageUrl,
		AppID:       idutil.NewSecretID(),
		AppSecret:   idutil.NewSecretKey(),
	}
	if len(req.IdentityProviderIds) != 0 {
		for _, idpId := range req.IdentityProviderIds {
			idp, err := a.Store.IdentityProviderRepository().GetByInstanceId(ctx, idpId, metav1.GetOptions{})
			if err != nil {
				log.Warnf("Failed to get the IdentityProvider [%s]: %v", idpId, err)
				continue
			}
			app.IdentityProviders = append(app.IdentityProviders, *idp)
		}
	}
	err := a.Store.ApplicationRepository().Create(ctx, app, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return err
}

func (a *applicationServiceImpl) UpdateApplication(ctx context.Context, app string, req v1.UpdateApplicationRequest) error {
	oldApp, err := a.GetApplication(ctx, app, metav1.GetOptions{})
	if err != nil {
		return nil
	}
	newApp := assembler.ConvertUpdateAppReqToModel(req, oldApp)
	if newApp.AppID == "" {
		newApp.AppID = idutil.NewSecretID()
	}
	if newApp.AppSecret == "" || req.RefreshAppSecret {
		newApp.AppSecret = idutil.NewSecretID()
	}
	if len(req.IdentityProviderIds) > 0 {
		for _, idpId := range req.IdentityProviderIds {
			idp, err := a.Store.IdentityProviderRepository().GetByInstanceId(ctx, idpId, metav1.GetOptions{})
			if err != nil {
				log.Warnf("Failed to get the IdentityProvider [%s]: %v", idpId, err)
				continue
			}
			newApp.IdentityProviders = append(newApp.IdentityProviders, *idp)
		}
	}
	return a.Store.ExecTx(ctx, func(ctx context.Context) error {
		err = a.Store.ApplicationRepository().Update(ctx, newApp, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		return nil
	})
}

func (a *applicationServiceImpl) DeleteApplication(ctx context.Context, app string, opts metav1.DeleteOptions) error {
	return a.Store.ApplicationRepository().Delete(ctx, app, opts)
}

func (a *applicationServiceImpl) GetApplication(ctx context.Context, idOrName string, opts metav1.GetOptions) (*v1.DetailApplicationResponse, error) {
	app, err := a.Store.ApplicationRepository().GetByInstanceIdOrName(ctx, idOrName, opts)
	if err != nil {
		log.Errorf("failed to get the app [%s]: %v", idOrName, err)
		return nil, err
	}
	base := assembler.ConvertModelToApplicationBase(app)

	return &v1.DetailApplicationResponse{
		ApplicationBase: *base,
	}, nil
}

func (a *applicationServiceImpl) ListApplications(ctx context.Context, opts metav1.ListOptions) (*v1.ApplicationList, error) {
	var appList []*v1.DetailApplicationResponse
	appRepo := a.Store.ApplicationRepository()
	apps, err := appRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	for _, v := range apps {
		base := assembler.ConvertModelToApplicationBase(&v)
		appList = append(appList, &v1.DetailApplicationResponse{
			ApplicationBase: *base,
		})
	}
	count, err := appRepo.Count(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &v1.ApplicationList{
		Items: appList,
		ListMeta: metav1.ListMeta{
			TotalCount: count,
		},
	}, nil
}
