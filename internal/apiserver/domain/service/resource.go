// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/coding-hui/common/errors"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1alpha1"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"
)

// ResourceService Resource manage api.
type ResourceService interface {
	Create(ctx context.Context, req v1alpha1.CreateResourceRequest) error
	Update(ctx context.Context, name string, req v1alpha1.UpdateResourceRequest) error
	Delete(ctx context.Context, name string, opts metav1alpha1.DeleteOptions) error
	DeleteCollection(ctx context.Context, names []string, opts metav1alpha1.DeleteOptions) error
	Get(ctx context.Context, name string, opts metav1alpha1.GetOptions) (*model.Resource, error)
	List(ctx context.Context, opts metav1alpha1.ListOptions) (*v1alpha1.ResourceList, error)
}

type resourceServiceImpl struct {
	Store repository.Factory `inject:"repository"`
}

// NewResourceService new Resource service.
func NewResourceService() ResourceService {
	return &resourceServiceImpl{}
}

func (r *resourceServiceImpl) Create(ctx context.Context, req v1alpha1.CreateResourceRequest) error {
	resource := assembler.CreateResourceModel(req)
	err := r.Store.ResourceRepository().Create(ctx, resource, metav1alpha1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (r *resourceServiceImpl) Update(ctx context.Context, name string, req v1alpha1.UpdateResourceRequest) error {
	resource, err := r.Get(ctx, name, metav1alpha1.GetOptions{})
	if err != nil {
		return err
	}
	if req.Api != "" {
		resource.Api = req.Api
	}
	if req.Description != "" {
		resource.Description = req.Description
	}
	resource.Actions = assembler.ConvertToActionModel(req.Actions)
	if err := r.Store.ResourceRepository().Update(ctx, resource, metav1alpha1.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}

func (r *resourceServiceImpl) Delete(ctx context.Context, name string, opts metav1alpha1.DeleteOptions) error {
	if err := r.Store.ResourceRepository().Delete(ctx, name, opts); err != nil {
		return err
	}

	return nil
}

func (r *resourceServiceImpl) DeleteCollection(ctx context.Context, names []string, opts metav1alpha1.DeleteOptions) error {
	if err := r.Store.ResourceRepository().DeleteCollection(ctx, names, opts); err != nil {
		return err
	}

	return nil
}

func (r *resourceServiceImpl) Get(ctx context.Context, name string, opts metav1alpha1.GetOptions) (*model.Resource, error) {
	resource, err := r.Store.ResourceRepository().Get(ctx, name, opts)
	if err != nil {
		return nil, errors.WrapC(err, code.ErrResourceNotFound, "Resource `%s` not found", name)
	}

	return resource, nil
}

// List list resources.
func (r *resourceServiceImpl) List(ctx context.Context, opts metav1alpha1.ListOptions) (*v1alpha1.ResourceList, error) {
	users, err := r.Store.ResourceRepository().List(ctx, metav1alpha1.ListOptions{
		Offset: opts.Offset,
		Limit:  opts.Limit,
	})
	if err != nil {
		return nil, err
	}

	return users, nil
}
