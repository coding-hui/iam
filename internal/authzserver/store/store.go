// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package store

import (
	"context"

	pb "github.com/coding-hui/iam/pkg/api/proto/apiserver/v1alpha1"
)

// PolicyStore defines the policy storage interface.
type PolicyStore interface {
	DetailPolicy(ctx context.Context, name string) (*pb.PolicyInfo, error)
	ListPolicies(ctx context.Context) (*pb.ListPoliciesResponse, error)
	ListPolicyRules(ctx context.Context) (*pb.ListPolicyRulesResponse, error)
}
