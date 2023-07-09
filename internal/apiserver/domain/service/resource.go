// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"bytes"
	"context"
	"os"
	"regexp"
	"strings"

	"github.com/bitly/go-simplejson"
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

const (
	ApiResourceDir = "api/swagger/swagger.json"
)

// ResourceService Resource manage api.
type ResourceService interface {
	CreateResource(ctx context.Context, req v1alpha1.CreateResourceRequest) error
	UpdateResource(ctx context.Context, instanceId string, req v1alpha1.UpdateResourceRequest) error
	DeleteResource(ctx context.Context, instanceId string, opts metav1alpha1.DeleteOptions) error
	BatchDeleteResources(ctx context.Context, names []string, opts metav1alpha1.DeleteOptions) error
	GetResource(ctx context.Context, instanceId string, opts metav1alpha1.GetOptions) (*model.Resource, error)
	DetailResource(ctx context.Context, resource *model.Resource, opts metav1alpha1.GetOptions) (*v1alpha1.DetailResourceResponse, error)
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
	jsonFile, _ := os.ReadFile(ApiResourceDir)
	apiDocs, _ := simplejson.NewFromReader(bytes.NewReader(jsonFile))
	for _, route := range routes {
		urlPath := route.Path
		idPatten := "(.*)/:(\\w+)"
		reg, _ := regexp.Compile(idPatten)
		if reg.MatchString(urlPath) {
			urlPath = reg.ReplaceAllString(route.Path, "${1}/{${2}}")
		}
		apiName, _ := apiDocs.Get("paths").Get(urlPath).Get(strings.ToLower(route.Method)).Get("summary").String()
		apiDesc, _ := apiDocs.Get("paths").Get(urlPath).Get(strings.ToLower(route.Method)).Get("description").String()
		createReq := v1alpha1.CreateResourceRequest{
			Name:        apiName,
			Method:      route.Method,
			Type:        string(v1alpha1.API),
			Api:         route.Path,
			Description: apiDesc,
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
		_, err := r.Store.ResourceRepository().GetByName(ctx, createReq.Name, metav1alpha1.GetOptions{})
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
	resource := assembler.ConvertResourceModel(req)
	err := r.Store.ResourceRepository().Create(ctx, resource, metav1alpha1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// UpdateResource update resources.
func (r *resourceServiceImpl) UpdateResource(ctx context.Context, instanceId string, req v1alpha1.UpdateResourceRequest) error {
	resource, err := r.GetResource(ctx, instanceId, metav1alpha1.GetOptions{})
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
func (r *resourceServiceImpl) DeleteResource(ctx context.Context, instanceId string, opts metav1alpha1.DeleteOptions) error {
	count, _ := r.Store.PolicyRepository().CountStatementByResource(ctx, instanceId)
	if count > 0 {
		return errors.WithCode(code.ErrResourceHasAssignedPolicy, "Resource [%s] has been assigned permission policies", instanceId)
	}
	if err := r.Store.ResourceRepository().DeleteByInstanceId(ctx, instanceId, opts); err != nil {
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
func (r *resourceServiceImpl) GetResource(ctx context.Context, instanceId string, opts metav1alpha1.GetOptions) (*model.Resource, error) {
	resource, err := r.Store.ResourceRepository().GetByInstanceId(ctx, instanceId, opts)
	if err != nil {
		return nil, errors.WrapC(err, code.ErrResourceNotFound, "Resource `%s` not found", instanceId)
	}

	return resource, nil
}

// DetailResource get resource detail.
func (r *resourceServiceImpl) DetailResource(
	_ context.Context,
	resource *model.Resource,
	_ metav1alpha1.GetOptions,
) (*v1alpha1.DetailResourceResponse, error) {
	base := assembler.ConvertResourceModelToBase(resource)
	detail := &v1alpha1.DetailResourceResponse{
		ResourceBase: *base,
	}

	return detail, nil
}

// ListResources list resources.
func (r *resourceServiceImpl) ListResources(ctx context.Context, opts metav1alpha1.ListOptions) (*v1alpha1.ResourceList, error) {
	resources, err := r.Store.ResourceRepository().List(ctx, metav1alpha1.ListOptions{
		Offset: opts.Offset,
		Limit:  opts.Limit,
	})
	if err != nil {
		return nil, err
	}

	return resources, nil
}
