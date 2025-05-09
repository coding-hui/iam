// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package adapter

import (
	"context"
	"errors"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"

	"github.com/coding-hui/iam/internal/authzserver/store"
	pb "github.com/coding-hui/iam/pkg/api/proto/apiserver/v1"
	"github.com/coding-hui/iam/pkg/log"
)

// Adapter represents the Redis adapter for policy storage.
type Adapter struct {
	cli        store.Factory
	isFiltered bool
}

func newAdapter(cli store.Factory) (*Adapter, error) {
	a := &Adapter{}
	a.cli = cli

	return a, nil
}

// NewAdapter is the constructor for Adapter.
func NewAdapter(cli store.Factory) (*Adapter, error) {
	return newAdapter(cli)
}

func toStringPolicyRule(c *pb.PolicyRuleInfo) []string {
	policy := make([]string, 0)
	if c.GetPType() != "" {
		policy = append(policy, c.GetPType())
	}
	if c.GetV0() != "" {
		policy = append(policy, c.GetV0())
	}
	if c.GetV1() != "" {
		policy = append(policy, c.GetV1())
	}
	if c.GetV2() != "" {
		policy = append(policy, c.GetV2())
	}
	if c.GetV3() != "" {
		policy = append(policy, c.GetV3())
	}
	if c.GetV4() != "" {
		policy = append(policy, c.GetV4())
	}
	if c.GetV5() != "" {
		policy = append(policy, c.GetV5())
	}

	return policy
}

// LoadPolicy loads policy from database.
func (a *Adapter) LoadPolicy(model model.Model) error {
	rules, err := a.cli.Policies().ListPolicyRules(context.Background())
	if err != nil {
		return err
	}

	for _, r := range rules.GetItems() {
		line := toStringPolicyRule(r)
		err := persist.LoadPolicyArray(line, model)
		if err != nil {
			log.Warnf("failed to load policy array. policy=%s err=%w", err, r)
		}
	}

	a.isFiltered = false
	return nil
}

// SavePolicy saves policy to database.
func (a *Adapter) SavePolicy(model model.Model) error {
	return nil
}

// AddPolicy adds a policy rule to the storage.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemovePolicy removes a policy rule from the storage.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}

// AddPolicies adds policy rules to the storage.
func (a *Adapter) AddPolicies(sec string, ptype string, rules [][]string) error {
	return errors.New("not implemented")
}

// RemovePolicies removes policy rules from the storage.
func (a *Adapter) RemovePolicies(sec string, ptype string, rules [][]string) error {
	return errors.New("not implemented")
}
