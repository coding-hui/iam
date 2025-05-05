// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import (
	"context"

	"github.com/AlekSi/pointer"
	"github.com/avast/retry-go"

	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"

	pb "github.com/coding-hui/iam/pkg/api/proto/apiserver/v1"
)

type policies struct {
	cli pb.CacheClient
}

func newPolicies(ds *datastore) *policies {
	return &policies{ds.cli}
}

// DetailPolicy get policy detail.
func (p *policies) DetailPolicy(ctx context.Context, name string) (*pb.PolicyInfo, error) {
	var policy *pb.PolicyInfo
	err := retry.Do(
		func() error {
			var listErr error
			policy, listErr = p.cli.DetailPolicy(ctx, &pb.GetPolicyRequest{Name: name})
			if listErr != nil {
				return listErr
			}

			return nil
		}, retry.Attempts(3),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get policy %s detail", name)
	}

	log.Infof("Policy %s found", name)

	return policy, nil
}

// ListPolicies returns all the authorization policies.
func (p *policies) ListPolicies(ctx context.Context) (*pb.ListPoliciesResponse, error) {
	req := &pb.ListPoliciesRequest{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(-1),
	}

	var resp *pb.ListPoliciesResponse
	err := retry.Do(
		func() error {
			var listErr error
			resp, listErr = p.cli.ListPolicies(ctx, req)
			if listErr != nil {
				return listErr
			}

			return nil
		}, retry.Attempts(3),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list policies")
	}

	log.Infof("Policies found (%d total)", len(resp.GetItems()))

	return resp, nil
}

// ListPolicyRules returns all the authorization policy rules.
func (p *policies) ListPolicyRules(ctx context.Context) (*pb.ListPolicyRulesResponse, error) {
	req := &pb.ListPolicyRulesRequest{
		Offset: pointer.ToInt64(0),
		Limit:  pointer.ToInt64(-1),
	}

	var resp *pb.ListPolicyRulesResponse
	err := retry.Do(
		func() error {
			var listErr error
			resp, listErr = p.cli.ListPolicyRules(ctx, req)
			if listErr != nil {
				return listErr
			}

			return nil
		}, retry.Attempts(3),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list policy rules")
	}

	log.Infof("PolicyRules found (%d total)", len(resp.GetItems()))

	return resp, nil
}
