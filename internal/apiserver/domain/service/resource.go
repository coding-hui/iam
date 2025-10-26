// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"fmt"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/code"
	"github.com/coding-hui/iam/pkg/log"
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
	if err := r.initSystemResources(ctx); err != nil {
		log.Warnf("Failed to initialize system resources: %v", err)
		return err
	}
	return nil
}

// initSystemResources initializes all system built-in API resources
func (r *resourceServiceImpl) initSystemResources(ctx context.Context) error {
	// Load statically configured API resources from model domain
	resources := model.SystemBuiltInResources()

	// Create resources if they don't exist
	for _, resource := range resources {
		if err := r.createSystemResource(ctx, resource.Name, resource.Description, resource.Actions, resource.Api); err != nil {
			log.Warnf("Failed to create system resource %s: %v", resource.Name, err)
			// Continue with other resources even if one fails
		}
	}

	log.Info("initialize system built-in API resources done")
	return nil
}

// createSystemResource creates a system resource with specified actions
func (r *resourceServiceImpl) createSystemResource(ctx context.Context, resourceName, description string, actions []string, apiPath string) error {
	// Check if resource already exists
	_, err := r.Store.ResourceRepository().GetByName(ctx, resourceName, metav1.GetOptions{})
	if err == nil {
		// Resource already exists, skip creation
		return nil
	}
	if !errors.IsCode(err, code.ErrResourceNotFound) {
		return err
	}

	// Convert string actions to Action models
	actionModels := make([]model.Action, 0, len(actions))
	for _, action := range actions {
		actionModels = append(actionModels, model.Action{
			Name:        fmt.Sprintf("%s:%s", resourceName, action),
			Description: fmt.Sprintf("%s operation for %s resource", action, resourceName),
		})
	}

	// Create the resource
	resource := &model.Resource{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceName,
		},
		Type:        string(v1.API),
		Api:         apiPath,
		Method:      "ALL",
		IsDefault:   true,
		Description: description,
		Actions:     actionModels,
	}

	if err := r.Store.ResourceRepository().Create(ctx, resource, metav1.CreateOptions{}); err != nil {
		return errors.WithMessagef(err, "Failed to create system resource %s", resourceName)
	}

	log.Infof("Created system resource: %s with actions: %v", resourceName, actions)
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
