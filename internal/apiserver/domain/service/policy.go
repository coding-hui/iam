// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"k8s.io/klog/v2"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1alpha1"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"

	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
)

// PolicyService Policy rule manage api.
type PolicyService interface {
	CreatePolicy(ctx context.Context, req v1alpha1.CreatePolicyRequest) error
	UpdatePolicy(ctx context.Context, idOrName string, req v1alpha1.UpdatePolicyRequest) error
	DeletePolicy(ctx context.Context, name string, opts metav1alpha1.DeleteOptions) error
	GetPolicy(ctx context.Context, name string, opts metav1alpha1.GetOptions) (*model.Policy, error)
	DetailPolicy(ctx context.Context, policy *model.Policy, opts metav1alpha1.GetOptions) (*v1alpha1.DetailPolicyResponse, error)
	ListPolicies(ctx context.Context, opts metav1alpha1.ListOptions) (*v1alpha1.PolicyList, error)
}

type policyServiceImpl struct {
	Store repository.Factory `inject:"repository"`
}

// NewPolicyService new Policy service.
func NewPolicyService() PolicyService {
	return &policyServiceImpl{}
}

// CreatePolicy create a new policy.
func (p *policyServiceImpl) CreatePolicy(ctx context.Context, req v1alpha1.CreatePolicyRequest) error {
	policy := &model.Policy{
		ObjectMeta: metav1alpha1.ObjectMeta{
			Name: req.Name,
		},
		Subjects:    req.Subjects,
		Resources:   req.Resources,
		Actions:     req.Actions,
		Effect:      req.Effect,
		Type:        req.Type,
		Status:      req.Status,
		Owner:       req.Owner,
		Description: req.Description,
	}
	err := p.Store.PolicyRepository().Create(ctx, policy, metav1alpha1.CreateOptions{})
	if err != nil {
		return err
	}

	e := p.Store.CasbinRepository().SyncedEnforcer()
	res, err := e.AddPolicies(policy.GetPolicyRules())
	if err != nil {
		return err
	}
	if !res {
		klog.Warning("The authorization rule %s already exists and cannot be added", policy.Name)
	}

	return err
}

// UpdatePolicy update policy.
func (p *policyServiceImpl) UpdatePolicy(ctx context.Context, name string, req v1alpha1.UpdatePolicyRequest) error {
	oldPolicy, err := p.Store.PolicyRepository().Get(ctx, name, metav1alpha1.GetOptions{})
	if err != nil {
		return err
	}
	policy := &model.Policy{
		ObjectMeta:  oldPolicy.ObjectMeta,
		Subjects:    req.Subjects,
		Resources:   req.Resources,
		Actions:     req.Actions,
		Effect:      req.Effect,
		Type:        req.Type,
		Description: req.Description,
	}
	err = p.Store.PolicyRepository().Update(ctx, policy, metav1alpha1.UpdateOptions{})
	if err != nil {
		return err
	}
	// delete rules before
	e := p.Store.CasbinRepository().SyncedEnforcer()
	_, err = e.RemovePolicies(oldPolicy.GetPolicyRules())
	if err != nil {
		return err
	}
	_, err = e.AddPolicies(policy.GetPolicyRules())
	if err != nil {
		return err
	}

	return nil
}

// DeletePolicy delete policy by id.
func (p *policyServiceImpl) DeletePolicy(ctx context.Context, name string, opts metav1alpha1.DeleteOptions) error {
	policy, err := p.Store.PolicyRepository().Get(ctx, name, metav1alpha1.GetOptions{})
	if err != nil {
		return err
	}

	e := p.Store.CasbinRepository().SyncedEnforcer()
	areRulesRemoved, err := e.RemovePolicies(policy.GetPolicyRules())
	if err != nil {
		return err
	}
	if !areRulesRemoved {
		klog.Warning("The rules is not removed. Check whether it exists. policyName: %s", policy.Name)
	}

	err = p.Store.PolicyRepository().Delete(ctx, name, opts)
	if err != nil {
		return err
	}

	return err
}

// GetPolicy get policy by id.
func (p *policyServiceImpl) GetPolicy(ctx context.Context, name string, opts metav1alpha1.GetOptions) (*model.Policy, error) {
	policy, err := p.Store.PolicyRepository().Get(ctx, name, opts)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

// DetailPolicy get policy details.
func (p *policyServiceImpl) DetailPolicy(
	_ context.Context,
	policy *model.Policy,
	opts metav1alpha1.GetOptions,
) (*v1alpha1.DetailPolicyResponse, error) {
	base := assembler.ConvertPolicyModelToBase(policy)
	detail := &v1alpha1.DetailPolicyResponse{
		PolicyBase: *base,
	}

	return detail, nil
}

// ListPolicies list policies.
func (p *policyServiceImpl) ListPolicies(ctx context.Context, opts metav1alpha1.ListOptions) (*v1alpha1.PolicyList, error) {
	policies, err := p.Store.PolicyRepository().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return policies, nil
}
