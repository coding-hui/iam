// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"

	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1alpha1"
	pb "github.com/coding-hui/iam/pkg/api/proto/apiserver/v1alpha1"

	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
)

var (
	cacheServer *Cache
	once        sync.Once
)

// Cache defines a cache service used to list all secrets and policies.
type Cache struct {
	*pb.UnimplementedCacheServer
	PolicyService service.PolicyService `inject:""`
}

func (c *Cache) RegisterApiGroup(_ *gin.Engine) {}

// NewCacheServer is the of cache server.
func NewCacheServer() *Cache {
	once.Do(func() {
		cacheServer = &Cache{}
	})
	return cacheServer
}

// DetailPolicy returns policy details.
func (c *Cache) DetailPolicy(ctx context.Context, r *pb.GetPolicyRequest) (*pb.PolicyInfo, error) {
	klog.Info("get policy detail function called.")
	policy, err := c.PolicyService.GetPolicy(ctx, r.GetName(), metav1alpha1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return assembler.ConvertPolicyModelToProtoInfo(policy), nil
}

// ListPolicies returns all policies.
func (c *Cache) ListPolicies(ctx context.Context, r *pb.ListPoliciesRequest) (*pb.ListPoliciesResponse, error) {
	klog.Info("list policies function called.")
	opts := metav1alpha1.ListOptions{
		Offset: r.Offset,
		Limit:  r.Limit,
	}

	policies, err := c.PolicyService.ListPolicies(ctx, opts)
	if err != nil {
		return nil, err
	}

	items := make([]*pb.PolicyInfo, 0)
	for _, pol := range policies.Items {
		items = append(items, &pb.PolicyInfo{
			Name:        pol.Name,
			Subjects:    pol.Subjects,
			Resources:   pol.Resources,
			Actions:     pol.Actions,
			Effect:      pol.Effect,
			Type:        pol.Type,
			Status:      pol.Status,
			Owner:       pol.Owner,
			Description: pol.Description,
			Adapter:     pol.Adapter,
			Model:       pol.Model,
		})
	}

	return &pb.ListPoliciesResponse{
		TotalCount: policies.TotalCount,
		Items:      items,
	}, nil
}
