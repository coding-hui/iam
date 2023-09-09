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
)

// OrganizationService Organization manage api.
type OrganizationService interface {
	CreateOrganization(ctx context.Context, req v1.CreateOrganizationRequest, opts metav1.CreateOptions) error
	UpdateOrganization(ctx context.Context, instanceId string, req v1.UpdateOrganizationRequest, opts metav1.UpdateOptions) error
	DeleteOrganization(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error
	GetOrganization(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.Organization, error)
	DetailOrganization(ctx context.Context, org *model.Organization) (*v1.DetailOrganizationResponse, error)
	ListOrganizations(ctx context.Context, opts metav1.ListOptions) (*v1.OrganizationList, error)
	DisableOrganization(ctx context.Context, instanceId string) error
	EnableOrganization(ctx context.Context, instanceId string) error
	Init(ctx context.Context) error
}

type organizationServiceImpl struct {
	Store repository.Factory `inject:"repository"`
}

// NewOrganizationService new Organization service.
func NewOrganizationService(c config.Config) OrganizationService {
	return &organizationServiceImpl{}
}

// Init initialize default org.
func (o *organizationServiceImpl) Init(ctx context.Context) error {
	old, err := o.Store.OrganizationRepository().GetByName(ctx, model.DefaultOrganization, metav1.GetOptions{})
	if err != nil && !errors.IsCode(err, code.ErrOrgNotFound) {
		return err
	}
	if old != nil {
		return nil
	}
	createReq := v1.CreateOrganizationRequest{
		Name:        model.DefaultOrganization,
		DisplayName: "Platform",
		WebsiteUrl:  "http://iam.wecoding.top",
		Favicon:     "",
		Disabled:    false,
		Description: "System Build-in Organization",
	}
	err = o.CreateOrganization(ctx, createReq, metav1.CreateOptions{})
	if err != nil {
		return errors.WithMessagef(err, "Failed to initialize system organization")
	}
	log.Info("initialize system organization done")

	return nil
}

func (o *organizationServiceImpl) CreateOrganization(
	ctx context.Context,
	req v1.CreateOrganizationRequest,
	opts metav1.CreateOptions,
) error {
	org := &model.Organization{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
		DisplayName: req.DisplayName,
		WebsiteUrl:  req.WebsiteUrl,
		Favicon:     req.Favicon,
		Disabled:    req.Disabled,
		Description: req.Description,
	}
	if org.DisplayName == "" {
		org.DisplayName = req.Name
	}
	err := o.Store.OrganizationRepository().Create(ctx, org, opts)
	if err != nil {
		return err
	}

	return nil
}

func (o *organizationServiceImpl) UpdateOrganization(
	ctx context.Context,
	idOrName string,
	req v1.UpdateOrganizationRequest,
	opts metav1.UpdateOptions,
) error {
	org, err := o.Store.OrganizationRepository().GetByInstanceId(ctx, idOrName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if req.DisplayName != "" {
		org.DisplayName = req.DisplayName
	}
	if req.Favicon != "" {
		org.Favicon = req.Favicon
	}
	if req.WebsiteUrl != "" {
		org.WebsiteUrl = req.WebsiteUrl
	}
	if req.Description != "" {
		org.Description = req.Description
	}
	err = o.Store.OrganizationRepository().Update(ctx, org, opts)
	if err != nil {
		return err
	}

	return nil
}

func (o *organizationServiceImpl) DeleteOrganization(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error {
	org, err := o.Store.OrganizationRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if org.IsSystemBuiltIn() {
		return errors.WithCode(code.ErrCannotDeleteBuiltInOrg, "")
	}
	err = o.Store.OrganizationRepository().DeleteByInstanceId(ctx, instanceId, opts)
	if err != nil {
		return err
	}

	return nil
}

func (o *organizationServiceImpl) GetOrganization(
	ctx context.Context,
	instanceId string,
	opts metav1.GetOptions,
) (*model.Organization, error) {
	org, err := o.Store.OrganizationRepository().GetByInstanceId(ctx, instanceId, opts)
	if err != nil {
		return nil, err
	}

	return org, nil
}

func (o *organizationServiceImpl) DetailOrganization(_ context.Context, org *model.Organization) (*v1.DetailOrganizationResponse, error) {
	base := *assembler.ConvertOrganizationModelToBase(org)
	return &v1.DetailOrganizationResponse{
		OrganizationBase: base,
	}, nil
}

func (o *organizationServiceImpl) ListOrganizations(ctx context.Context, opts metav1.ListOptions) (*v1.OrganizationList, error) {
	items, err := o.Store.OrganizationRepository().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (o *organizationServiceImpl) DisableOrganization(ctx context.Context, instanceId string) error {
	org, err := o.Store.OrganizationRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if org.IsSystemBuiltIn() {
		return errors.WithCode(code.ErrCannotDisableBuiltInOrg, "")
	}
	if org.Disabled {
		return errors.WithCode(code.ErrOrgAlreadyDisabled, "The organization [%s] is already disabled.", org.Name)
	}
	org.Disabled = true

	return o.Store.OrganizationRepository().Update(ctx, org, metav1.UpdateOptions{})
}

func (o *organizationServiceImpl) EnableOrganization(ctx context.Context, instanceId string) error {
	org, err := o.Store.OrganizationRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if !org.Disabled {
		return errors.WithCode(code.ErrOrgAlreadyEnabled, "The organization [%s] is already enabled.", org.Name)
	}
	org.Disabled = false

	return o.Store.OrganizationRepository().Update(ctx, org, metav1.UpdateOptions{})
}
