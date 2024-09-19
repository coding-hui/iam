// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"sync"

	"github.com/lib/pq"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	assembler "github.com/coding-hui/iam/internal/apiserver/interfaces/api/assembler/v1"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/code"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
)

// PolicyService Policy rule manage api.
type PolicyService interface {
	CreatePolicy(ctx context.Context, req v1.CreatePolicyRequest) error
	UpdatePolicy(ctx context.Context, idOrName string, req v1.UpdatePolicyRequest) error
	DeletePolicy(ctx context.Context, name string, opts metav1.DeleteOptions) error
	GetPolicy(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.Policy, error)
	DetailPolicy(ctx context.Context, policy *model.Policy, opts metav1.GetOptions) (*v1.DetailPolicyResponse, error)
	ListPolicies(ctx context.Context, opts metav1.ListOptions) (*v1.PolicyList, error)
	ListPolicyRules(ctx context.Context, opts metav1.ListOptions) ([]model.PolicyRule, error)
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
	platformRole, err := p.Store.RoleRepository().GetByName(ctx, v1.PlatformAdmin.String(), metav1.GetOptions{})
	if err != nil {
		return errors.WithMessagef(err, "Failed to get %s role info.", v1.PlatformAdmin.String())
	}
	createReq := v1.CreatePolicyRequest{
		Name:        DefaultAdmin,
		Subjects:    []string{platformRole.InstanceID},
		Type:        string(v1.SystemBuildInPolicy),
		Owner:       DefaultAdmin,
		Description: "System default admin policies",
		Statements: []v1.Statement{
			{
				Effect:             v1.AllowAccess,
				Resource:           "*",
				ResourceIdentifier: "*:*",
				Actions:            pq.StringArray{"*"},
			},
		},
	}
	_, err = p.Store.PolicyRepository().GetByName(ctx, createReq.Name, metav1.GetOptions{})
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
func (p *policyServiceImpl) CreatePolicy(ctx context.Context, req v1.CreatePolicyRequest) error {
	if len(req.Statements) == 0 {
		return nil
	}
	policy := assembler.ConvertPolicyModel(req)
	if len(policy.Type) == 0 {
		policy.Type = string(v1.CustomPolicy)
	}
	err := p.Store.PolicyRepository().Create(ctx, policy, metav1.CreateOptions{})
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
func (p *policyServiceImpl) UpdatePolicy(ctx context.Context, idOrName string, req v1.UpdatePolicyRequest) error {
	oldPolicy, err := p.Store.PolicyRepository().GetByInstanceId(ctx, idOrName, metav1.GetOptions{})
	if err != nil {
		oldPolicy, err = p.Store.PolicyRepository().GetByName(ctx, idOrName, metav1.GetOptions{})
		if err != nil {
			return err
		}
	}
	// delete casbin rules before
	e := p.Store.CasbinRepository().SyncedEnforcer()
	_, err = e.RemovePolicies(oldPolicy.GetPolicyRules())
	if err != nil {
		return err
	}

	// update policy info
	oldPolicy.Subjects = req.Subjects
	oldPolicy.Description = req.Description
	oldPolicy.Status = req.Status
	oldPolicy.Owner = req.Owner
	oldPolicy.Statements = assembler.ConvertToStatementModel(req.Statements)
	err = p.Store.PolicyRepository().Update(ctx, oldPolicy, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	// add casbin rules
	_, err = e.AddPolicies(oldPolicy.GetPolicyRules())
	if err != nil {
		return err
	}

	return nil
}

// DeletePolicy delete policy by id.
func (p *policyServiceImpl) DeletePolicy(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	policy, err := p.Store.PolicyRepository().GetByName(ctx, name, metav1.GetOptions{})
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
func (p *policyServiceImpl) GetPolicy(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.Policy, error) {
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
	_ metav1.GetOptions,
) (*v1.DetailPolicyResponse, error) {
	wg := sync.WaitGroup{}
	errChan := make(chan error, 1)
	finished := make(chan bool, 1)

	var m sync.Map

	for _, statement := range policy.Statements {
		if statement.Resource == "*" {
			continue
		}
		wg.Add(1)
		go func(s model.Statement) {
			defer wg.Done()

			r, err := p.Store.ResourceRepository().GetByInstanceId(ctx, s.Resource, metav1.GetOptions{})
			if err != nil {
				errChan <- errors.WithMessagef(err, "load resource [%s] failed.", s.ResourceIdentifier)
				return
			}
			m.Store(s.ID, r)
		}(statement)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChan:
		log.Errorf("failed to load resources: %v", err)
	}

	resources := make([]v1.ResourceBase, 0, len(policy.Statements))
	for _, statement := range policy.Statements {
		r, ok := m.Load(statement.ID)
		if ok {
			resources = append(resources, *assembler.ConvertResourceModelToBase(r.(*model.Resource)))
		}
	}
	base := assembler.ConvertPolicyModelToBase(policy)

	return &v1.DetailPolicyResponse{
		PolicyBase: *base,
		Resources:  resources,
	}, nil
}

// ListPolicies list policies.
func (p *policyServiceImpl) ListPolicies(ctx context.Context, opts metav1.ListOptions) (*v1.PolicyList, error) {
	policies, err := p.Store.PolicyRepository().List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return policies, nil
}

// ListPolicyRules list policy rules.
func (p *policyServiceImpl) ListPolicyRules(ctx context.Context, opts metav1.ListOptions) ([]model.PolicyRule, error) {
	e := p.Store.CasbinRepository().SyncedEnforcer()

	var rules []model.PolicyRule

	pRules, err := e.GetPolicy()
	if err != nil {
		return nil, err
	}
	for _, v := range pRules {
		line := savePolicyLine("p", v)
		rules = append(rules, line)
	}

	gRules, err := e.GetGroupingPolicy()
	if err != nil {
		return nil, err
	}
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
