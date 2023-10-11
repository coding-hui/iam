// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/pkg/code"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
)

// ProviderService Provider manage api.
type ProviderService interface {
	CreateProvider(ctx context.Context, req v1.CreateProviderRequest) error
	UpdateProvider(ctx context.Context, idOrName string, req v1.UpdateProviderRequest) error
	DeleteProvider(ctx context.Context, name string, opts metav1.DeleteOptions) error
	GetProvider(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.Provider, error)
	Init(ctx context.Context) error
}

type providerServiceImpl struct {
	cfg   config.Config
	Store repository.Factory `inject:"repository"`
}

// NewProviderService new Provider service.
func NewProviderService(c config.Config) ProviderService {
	return &providerServiceImpl{cfg: c}
}

// Init initialize resource data.
func (p *providerServiceImpl) Init(ctx context.Context) error {
	for _, idp := range p.cfg.AuthenticationOptions.OAuthOptions.IdentityProviders {
		_, err := p.Store.ProviderRepository().GetByName(ctx, idp.Name, metav1.GetOptions{})
		if err != nil && errors.IsCode(err, code.ErrRecordNotExist) {
			createReq := v1.CreateProviderRequest{
				Name:          idp.Name,
				Type:          idp.Type,
				Category:      idp.Category,
				DisplayName:   idp.Name,
				MappingMethod: idp.MappingMethod,
				Description:   "Built-in IdentityProvider",
				Extend:        metav1.Extend(idp.Provider),
			}
			err = p.CreateProvider(ctx, createReq)
			if err != nil {
				return errors.WithMessagef(err, "Failed to initialize default IdentityProvider %s", idp.Name)
			}
			log.Infof("initialize %s IdentityProvider done", idp.Name)
		}
	}

	return nil
}

func (p *providerServiceImpl) CreateProvider(ctx context.Context, req v1.CreateProviderRequest) error {
	provider := &model.Provider{
		ObjectMeta: metav1.ObjectMeta{
			Name:   req.Name,
			Extend: req.Extend,
		},
		Type:          req.Type,
		Category:      req.Category,
		Status:        req.Status,
		Owner:         req.Owner,
		DisplayName:   req.DisplayName,
		Description:   req.Description,
		MappingMethod: req.MappingMethod,
	}
	err := p.Store.ProviderRepository().Create(ctx, provider, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return err
}

func (p *providerServiceImpl) UpdateProvider(ctx context.Context, idOrName string, req v1.UpdateProviderRequest) error {
	return nil
}

func (p *providerServiceImpl) DeleteProvider(ctx context.Context, idOrName string, opts metav1.DeleteOptions) error {
	return p.Store.ProviderRepository().Delete(ctx, idOrName, opts)
}

func (p *providerServiceImpl) GetProvider(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.Provider, error) {
	provider, err := p.Store.ProviderRepository().GetByInstanceId(ctx, instanceId, opts)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
