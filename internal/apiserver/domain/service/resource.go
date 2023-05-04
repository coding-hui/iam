// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"

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
	CreateResource(ctx context.Context, req v1alpha1.CreateResourceRequest) error
	UpdateResource(ctx context.Context, name string, req v1alpha1.UpdateResourceRequest) error
	DeleteResource(ctx context.Context, name string, opts metav1alpha1.DeleteOptions) error
	BatchDeleteResources(ctx context.Context, names []string, opts metav1alpha1.DeleteOptions) error
	GetResource(ctx context.Context, name string, opts metav1alpha1.GetOptions) (*model.Resource, error)
	ListResources(ctx context.Context, opts metav1alpha1.ListOptions) (*v1alpha1.ResourceList, error)
	Init(ctx context.Context) error
}

type resourceServiceImpl struct {
	Store repository.Factory `inject:"repository"`
}

// NewResourceService new Resource service.
func NewResourceService() ResourceService {
	return &resourceServiceImpl{}
}

// Init initialize resource data.
func (r *resourceServiceImpl) Init(ctx context.Context) error {
	routes := ctx.Value(&v1alpha1.CtxKeyRoutes).(gin.RoutesInfo)
	apiPrefix := ctx.Value(&v1alpha1.CtxKeyApiPrefix).([]string)
	if len(routes) == 0 {
		klog.Warning("Failed to get the registered route from the init context.")
		return nil
	}
	for _, route := range routes {
		createReq := v1alpha1.CreateResourceRequest{
			Name:        route.Handler[strings.LastIndex(route.Handler, ".")+1:],
			Method:      route.Method,
			Type:        "API",
			Api:         route.Path,
			Description: route.Handler,
			IsDefault:   true,
			Actions:     nil,
		}
		found := false
		for _, prefix := range apiPrefix {
			found = strings.Contains(route.Path, prefix)
		}
		if !found {
			continue
		}
		_, err := r.GetResource(ctx, createReq.Name, metav1alpha1.GetOptions{})
		if err != nil && errors.IsCode(err, code.ErrResourceNotFound) {
			if err := r.CreateResource(ctx, createReq); err != nil {
				klog.Warningf("Failed to create api resource. [Api: %s Method: %s Handler: %s]", route.Path, route.Method, route.Handler)
			}
		}
	}
	klog.Info("initialize system default api resource done")

	return nil
}

// CreateResource create a new resource.
func (r *resourceServiceImpl) CreateResource(ctx context.Context, req v1alpha1.CreateResourceRequest) error {
	resource := assembler.CreateResourceModel(req)
	err := r.Store.ResourceRepository().Create(ctx, resource, metav1alpha1.CreateOptions{})
	if err != nil {
		return err
	}
	e := r.Store.CasbinRepository().SyncedEnforcer()
	_, err = e.AddNamedPolicy("p", v1alpha1.PlatformAdmin.String(), resource.Api, resource.Method)
	if err != nil {
		return err
	}

	return nil
}

// UpdateResource update resources.
func (r *resourceServiceImpl) UpdateResource(ctx context.Context, name string, req v1alpha1.UpdateResourceRequest) error {
	resource, err := r.GetResource(ctx, name, metav1alpha1.GetOptions{})
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

// DeleteResource delete resources.
func (r *resourceServiceImpl) DeleteResource(ctx context.Context, name string, opts metav1alpha1.DeleteOptions) error {
	if err := r.Store.ResourceRepository().Delete(ctx, name, opts); err != nil {
		return err
	}

	return nil
}

// BatchDeleteResources batch delete resources.
func (r *resourceServiceImpl) BatchDeleteResources(ctx context.Context, names []string, opts metav1alpha1.DeleteOptions) error {
	if err := r.Store.ResourceRepository().DeleteCollection(ctx, names, opts); err != nil {
		return err
	}

	return nil
}

// GetResource get resource.
func (r *resourceServiceImpl) GetResource(ctx context.Context, name string, opts metav1alpha1.GetOptions) (*model.Resource, error) {
	resource, err := r.Store.ResourceRepository().Get(ctx, name, opts)
	if err != nil {
		return nil, errors.WrapC(err, code.ErrResourceNotFound, "Resource `%s` not found", name)
	}

	return resource, nil
}

// ListResources list resources.
func (r *resourceServiceImpl) ListResources(ctx context.Context, opts metav1alpha1.ListOptions) (*v1alpha1.ResourceList, error) {
	users, err := r.Store.ResourceRepository().List(ctx, metav1alpha1.ListOptions{
		Offset: opts.Offset,
		Limit:  opts.Limit,
	})
	if err != nil {
		return nil, err
	}

	return users, nil
}
