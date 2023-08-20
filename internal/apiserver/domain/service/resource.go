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

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1"
	"github.com/coding-hui/iam/internal/pkg/code"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
)

const (
	ApiResourceDir = "api/swagger/swagger.json"
)

// ResourceService Resource manage api.
type ResourceService interface {
	CreateResource(ctx context.Context, req v1.CreateResourceRequest) error
	UpdateResource(ctx context.Context, instanceId string, req v1.UpdateResourceRequest) error
	DeleteResource(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error
	BatchDeleteResources(ctx context.Context, names []string, opts metav1.DeleteOptions) error
	GetResource(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.Resource, error)
	DetailResource(ctx context.Context, resource *model.Resource, opts metav1.GetOptions) (*v1.DetailResourceResponse, error)
	ListResources(ctx context.Context, opts metav1.ListOptions) (*v1.ResourceList, error)
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
	routes := ctx.Value(&v1.CtxKeyRoutes).(gin.RoutesInfo)
	apiPrefix := ctx.Value(&v1.CtxKeyApiPrefix).([]string)
	if len(routes) == 0 {
		log.Warnf("Failed to get the registered route from the init context.")
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
		createReq := v1.CreateResourceRequest{
			Name:        apiName,
			Method:      route.Method,
			Type:        string(v1.API),
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
		_, err := r.Store.ResourceRepository().GetByName(ctx, createReq.Name, metav1.GetOptions{})
		if err != nil && errors.IsCode(err, code.ErrResourceNotFound) {
			if err := r.CreateResource(ctx, createReq); err != nil {
				log.Warnf("Failed to create api resource. [Api: %s Method: %s Handler: %s]", route.Path, route.Method, route.Handler)
			}
		}
	}
	log.Info("initialize system default api resource done")

	return nil
}

// CreateResource create a new resource.
func (r *resourceServiceImpl) CreateResource(ctx context.Context, req v1.CreateResourceRequest) error {
	resource := assembler.ConvertResourceModel(req)
	err := r.Store.ResourceRepository().Create(ctx, resource, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// UpdateResource update resources.
func (r *resourceServiceImpl) UpdateResource(ctx context.Context, instanceId string, req v1.UpdateResourceRequest) error {
	resource, err := r.GetResource(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if req.Api != "" {
		resource.Api = req.Api
	}
	if req.Description != "" {
		resource.Description = req.Description
	}
	resource.Actions = assembler.ConvertToActionModel(resource.Name, req.Actions)
	if err := r.Store.ResourceRepository().Update(ctx, resource, metav1.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}

// DeleteResource delete resources.
func (r *resourceServiceImpl) DeleteResource(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error {
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
func (r *resourceServiceImpl) BatchDeleteResources(ctx context.Context, names []string, opts metav1.DeleteOptions) error {
	if err := r.Store.ResourceRepository().DeleteCollection(ctx, names, opts); err != nil {
		return err
	}

	return nil
}

// GetResource get resource.
func (r *resourceServiceImpl) GetResource(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.Resource, error) {
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
	_ metav1.GetOptions,
) (*v1.DetailResourceResponse, error) {
	base := assembler.ConvertResourceModelToBase(resource)
	detail := &v1.DetailResourceResponse{
		ResourceBase: *base,
	}

	return detail, nil
}

// ListResources list resources.
func (r *resourceServiceImpl) ListResources(ctx context.Context, opts metav1.ListOptions) (*v1.ResourceList, error) {
	resources, err := r.Store.ResourceRepository().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return resources, nil
}
