// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"

	"github.com/lib/pq"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1alpha1"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
)

// PolicyService Policy rule manage api.
type PolicyService interface {
	CreatePolicy(ctx context.Context, req v1alpha1.CreatePolicyRequest) error
	UpdatePolicy(ctx context.Context, idOrName string, req v1alpha1.UpdatePolicyRequest) error
	DeletePolicy(ctx context.Context, name string, opts metav1alpha1.DeleteOptions) error
	GetPolicy(ctx context.Context, instanceId string, opts metav1alpha1.GetOptions) (*model.Policy, error)
	DetailPolicy(ctx context.Context, policy *model.Policy, opts metav1alpha1.GetOptions) (*v1alpha1.DetailPolicyResponse, error)
	ListPolicies(ctx context.Context, opts metav1alpha1.ListOptions) (*v1alpha1.PolicyList, error)
	ListPolicyRules(ctx context.Context, opts metav1alpha1.ListOptions) ([]model.PolicyRule, error)
	Init(ctx context.Context) error
}

type policyServiceImpl struct {
	Store repository.Factory `inject:"repository"`
}

// NewPolicyService new Policy service.
func NewPolicyService() PolicyService {
	return &policyServiceImpl{}
}

// Init initialize resource data.
func (p *policyServiceImpl) Init(ctx context.Context) error {
	// find platform role
	platformRole, err := p.Store.RoleRepository().GetByName(ctx, v1alpha1.PlatformAdmin.String(), metav1alpha1.GetOptions{})
	if err != nil {
		return errors.WithMessagef(err, "Failed to get %s role info.", v1alpha1.PlatformAdmin.String())
	}
	createReq := v1alpha1.CreatePolicyRequest{
		Name:        DefaultAdmin,
		Subjects:    []string{platformRole.InstanceID},
		Type:        string(v1alpha1.SystemBuildInPolicy),
		Owner:       DefaultAdmin,
		Description: "System default admin policies",
		Statements: []v1alpha1.Statement{
			{
				Effect:             v1alpha1.AllowAccess,
				Resource:           "*",
				ResourceIdentifier: "*:*",
				Actions:            pq.StringArray{"*"},
			},
		},
	}
	_, err = p.Store.PolicyRepository().GetByName(ctx, createReq.Name, metav1alpha1.GetOptions{})
	if err != nil && errors.IsCode(err, code.ErrPolicyNotFound) {
		if err := p.CreatePolicy(ctx, createReq); err != nil {
			log.Warnf("Failed to create admin policy.")
			return err
		}
	}
	log.Info("initialize system default policies done")

	return nil
}

// CreatePolicy create a new policy.
func (p *policyServiceImpl) CreatePolicy(ctx context.Context, req v1alpha1.CreatePolicyRequest) error {
	if len(req.Statements) < 0 {
		return nil
	}
	policy := assembler.ConvertPolicyModel(req)
	if len(policy.Type) == 0 {
		policy.Type = string(v1alpha1.CustomPolicy)
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
		log.Warnf("The authorization rule %s already exists and cannot be added", policy.Name)
	}

	return err
}

// UpdatePolicy update policy.
func (p *policyServiceImpl) UpdatePolicy(ctx context.Context, idOrName string, req v1alpha1.UpdatePolicyRequest) error {
	oldPolicy, err := p.Store.PolicyRepository().GetByInstanceId(ctx, idOrName, metav1alpha1.GetOptions{})
	if err != nil {
		oldPolicy, err = p.Store.PolicyRepository().GetByName(ctx, idOrName, metav1alpha1.GetOptions{})
		if err != nil {
			return err
		}
	}
	oldPolicy.Subjects = req.Subjects
	oldPolicy.Description = req.Description
	oldPolicy.Status = req.Status
	oldPolicy.Owner = req.Owner
	oldPolicy.Statements = assembler.ConvertToStatementModel(req.Statements)
	err = p.Store.PolicyRepository().Update(ctx, oldPolicy, metav1alpha1.UpdateOptions{})
	if err != nil {
		return err
	}
	// delete rules before
	e := p.Store.CasbinRepository().SyncedEnforcer()
	_, err = e.RemovePolicies(oldPolicy.GetPolicyRules())
	if err != nil {
		return err
	}
	_, err = e.AddPolicies(oldPolicy.GetPolicyRules())
	if err != nil {
		return err
	}

	return nil
}

// DeletePolicy delete policy by id.
func (p *policyServiceImpl) DeletePolicy(ctx context.Context, name string, opts metav1alpha1.DeleteOptions) error {
	policy, err := p.Store.PolicyRepository().GetByName(ctx, name, metav1alpha1.GetOptions{})
	if err != nil {
		return err
	}

	e := p.Store.CasbinRepository().SyncedEnforcer()
	areRulesRemoved, err := e.RemovePolicies(policy.GetPolicyRules())
	if err != nil {
		return err
	}
	if !areRulesRemoved {
		log.Warnf("The rules is not removed. Check whether it exists. policyName: %s", policy.Name)
	}

	err = p.Store.PolicyRepository().Delete(ctx, name, opts)
	if err != nil {
		return err
	}

	return err
}

// GetPolicy get policy by id.
func (p *policyServiceImpl) GetPolicy(ctx context.Context, instanceId string, opts metav1alpha1.GetOptions) (*model.Policy, error) {
	policy, err := p.Store.PolicyRepository().GetByInstanceId(ctx, instanceId, opts)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

// DetailPolicy get policy details.
func (p *policyServiceImpl) DetailPolicy(
	ctx context.Context,
	policy *model.Policy,
	_ metav1alpha1.GetOptions,
) (*v1alpha1.DetailPolicyResponse, error) {
	var resources []v1alpha1.ResourceBase
	for _, statement := range policy.Statements {
		if statement.Resource == "*" {
			continue
		}
		r, _ := p.Store.ResourceRepository().GetByInstanceId(ctx, statement.Resource, metav1alpha1.GetOptions{})
		resources = append(resources, *assembler.ConvertResourceModelToBase(r))
	}
	base := assembler.ConvertPolicyModelToBase(policy)
	detail := &v1alpha1.DetailPolicyResponse{
		PolicyBase: *base,
		Resources:  resources,
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

// ListPolicyRules list policy rules.
func (p *policyServiceImpl) ListPolicyRules(ctx context.Context, opts metav1alpha1.ListOptions) ([]model.PolicyRule, error) {
	e := p.Store.CasbinRepository().SyncedEnforcer()

	var rules []model.PolicyRule

	pRules := e.GetPolicy()
	for _, v := range pRules {
		line := savePolicyLine("p", v)
		rules = append(rules, line)
	}

	gRules := e.GetGroupingPolicy()
	for _, v := range gRules {
		line := savePolicyLine("g", v)
		rules = append(rules, line)
	}

	return rules, nil
}

func savePolicyLine(ptype string, rule []string) model.PolicyRule {
	line := model.PolicyRule{}

	line.PType = ptype
	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}
