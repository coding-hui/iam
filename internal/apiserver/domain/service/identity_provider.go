// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"gopkg.in/square/go-jose.v2/json"

	_ "github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider/coding"
	_ "github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider/gitee"
	_ "github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider/github"
	_ "github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider/ldap"
	_ "github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider/wechatmini"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/apiserver/domain/service/identityprovider"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1"
	"github.com/coding-hui/iam/internal/pkg/code"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
)

// IdentityProviderService IdentityProvider manage api.
type IdentityProviderService interface {
	CreateIdentityProvider(ctx context.Context, req v1.CreateIdentityProviderRequest) error
	UpdateIdentityProvider(ctx context.Context, identifier string, req v1.UpdateIdentityProviderRequest) error
	DeleteIdentityProvider(ctx context.Context, identifier string, opts metav1.DeleteOptions) error
	GetIdentityProvider(ctx context.Context, identifier string, opts metav1.GetOptions) (*model.IdentityProvider, error)
	DetailIdentityProvider(ctx context.Context, idp *model.IdentityProvider, opts metav1.GetOptions) (*v1.DetailIdentityProviderResponse, error)
	ListIdentityProviders(ctx context.Context, opts metav1.ListOptions) (*v1.IdentityProviderList, error)
	Init(ctx context.Context) error
}

type identityProviderServiceImpl struct {
	cfg   config.Config
	Store repository.Factory `inject:"repository"`
}

// NewIdentityProviderService new IdentityProvider service.
func NewIdentityProviderService(c config.Config) IdentityProviderService {
	return &identityProviderServiceImpl{cfg: c}
}

// Init initialize resource data.
func (i *identityProviderServiceImpl) Init(ctx context.Context) error {
	oauthOpts := i.cfg.AuthenticationOptions.OAuthOptions
	for _, idp := range oauthOpts.IdentityProviders {
		_, err := i.Store.IdentityProviderRepository().GetByName(ctx, idp.Name, metav1.GetOptions{})
		if err != nil && errors.IsCode(err, code.ErrRecordNotExist) {
			createReq := v1.CreateIdentityProviderRequest{
				Name:          idp.Name,
				Type:          idp.Type,
				Category:      idp.Category,
				DisplayName:   idp.Name,
				MappingMethod: idp.MappingMethod,
				CallbackURL:   idp.CallbackURL,
				Description:   "Built-in IdentityProvider",
				Config:        metav1.Extend(idp.Provider),
			}
			err = i.CreateIdentityProvider(ctx, createReq)
			if err != nil {
				return errors.WithMessagef(err, "Failed to initialize default IdentityProvider %s", idp.Name)
			}
			log.Infof("initialize %s IdentityProvider done", idp.Name)
		}
	}

	return nil
}

func (i *identityProviderServiceImpl) CreateIdentityProvider(ctx context.Context, req v1.CreateIdentityProviderRequest) error {
	idp := &model.IdentityProvider{
		ObjectMeta: metav1.ObjectMeta{
			Name:   req.Name,
			Extend: req.Config,
		},
		Type:          req.Type,
		Category:      req.Category,
		Status:        req.Status,
		Owner:         req.Owner,
		DisplayName:   req.DisplayName,
		Description:   req.Description,
		MappingMethod: req.MappingMethod,
		CallbackURL:   req.CallbackURL,
	}
	if idp.MappingMethod == "" {
		idp.MappingMethod = v1.MappingMethodAuto
	}
	if err := i.applyProviderConfig(idp); err != nil {
		return err
	}
	err := i.Store.IdentityProviderRepository().Create(ctx, idp, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return err
}

func (i *identityProviderServiceImpl) UpdateIdentityProvider(
	ctx context.Context,
	identifier string,
	req v1.UpdateIdentityProviderRequest,
) error {
	idp, err := i.Store.IdentityProviderRepository().GetByInstanceIdOrName(ctx, identifier, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if req.MappingMethod != "" {
		idp.MappingMethod = req.MappingMethod
	}
	if req.Category != "" {
		idp.Category = req.Category
	}
	if req.Status != "" {
		idp.Status = req.Status
	}
	if req.Owner != "" {
		idp.Owner = req.Owner
	}
	if req.DisplayName != "" {
		idp.DisplayName = req.DisplayName
	}
	if req.Description != "" {
		idp.Description = req.Description
	}
	if req.CallbackURL != "" {
		idp.CallbackURL = req.CallbackURL
	}
	if req.Config != nil {
		idp.Extend.Merge(req.Config.String())
	}
	if err := i.applyProviderConfig(idp); err != nil {
		return err
	}
	err = i.Store.IdentityProviderRepository().Update(ctx, idp, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (i *identityProviderServiceImpl) DeleteIdentityProvider(ctx context.Context, identifier string, opts metav1.DeleteOptions) error {
	return i.Store.IdentityProviderRepository().Delete(ctx, identifier, opts)
}

func (i *identityProviderServiceImpl) GetIdentityProvider(
	ctx context.Context,
	identifier string,
	opts metav1.GetOptions,
) (*model.IdentityProvider, error) {
	idp, err := i.Store.IdentityProviderRepository().GetByInstanceIdOrName(ctx, identifier, opts)
	if err != nil {
		return nil, err
	}
	return idp, err
}

func (i *identityProviderServiceImpl) DetailIdentityProvider(
	_ context.Context,
	idp *model.IdentityProvider,
	_ metav1.GetOptions,
) (*v1.DetailIdentityProviderResponse, error) {
	base := assembler.ConvertModelToIdentityProviderBase(idp)
	return &v1.DetailIdentityProviderResponse{
		IdentityProviderBase: *base,
	}, nil
}

func (i *identityProviderServiceImpl) ListIdentityProviders(
	ctx context.Context,
	opts metav1.ListOptions,
) (*v1.IdentityProviderList, error) {
	var idpList []*v1.DetailIdentityProviderResponse
	idpRepo := i.Store.IdentityProviderRepository()
	identityProviders, err := idpRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}
	for _, v := range identityProviders {
		idpList = append(idpList, &v1.DetailIdentityProviderResponse{
			IdentityProviderBase: *assembler.ConvertModelToIdentityProviderBase(&v),
		})
	}
	count, err := idpRepo.Count(ctx, opts)
	if err != nil {
		return nil, err
	}

	return &v1.IdentityProviderList{
		Items: idpList,
		ListMeta: metav1.ListMeta{
			TotalCount: count,
		},
	}, nil
}

func (i *identityProviderServiceImpl) applyProviderConfig(idp *model.IdentityProvider) error {
	var providerConf string
	if idp.Category == v1.OAuth {
		oauthProviderConfig, err := identityprovider.SetOAuthProvider(idp)
		if err != nil {
			log.Warnf("Failed to get the oauth IdentityProvider config for [%s]: %v", idp.Name, err)
		}
		if oauthProviderConfig != nil {
			conf, err := json.Marshal(oauthProviderConfig)
			if err != nil {
				log.Errorf("Failed to marshal the oauth IdentityProvider config for [%s]: %v", idp.Name, err)
				return err
			}
			providerConf = string(conf)
		}
	} else {
		genericProvider, err := identityprovider.SetGenericProvider(idp)
		if err != nil {
			log.Warnf("Failed to get the generic IdentityProvider config for [%s]: %v", idp.Name, err)
		}
		if genericProvider != nil {
			conf, err := json.Marshal(genericProvider)
			if err != nil {
				log.Errorf("Failed to marshal the generic IdentityProvider config for [%s]: %v", idp.Name, err)
				return err
			}
			providerConf = string(conf)
		}
	}
	if len(providerConf) == 0 {
		log.Warnf("The IdentityProvider [%s] config is empty", idp.Name)
		return nil
	}
	idp.Extend = make(metav1.Extend)
	idp.Extend.Merge(providerConf)

	return nil
}
