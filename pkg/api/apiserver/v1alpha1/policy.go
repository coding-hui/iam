// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1alpha1

import (
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
)

const (
	// AllowAccess should be used as effect for policies that allow access.
	AllowAccess string = "allow"
	// DenyAccess should be used as effect for policies that deny access.
	DenyAccess string = "deny"
)

// PolicyType define policy type.
type PolicyType string

const (
	// SystemBuildIn system default policy type.
	SystemBuildIn PolicyType = "SYSTEM"
)

// CreatePolicyRequest create policy request.
type CreatePolicyRequest struct {
	Name        string   `json:"name"        validate:"required,name"`
	Description string   `json:"description" validate:"min=1,max=30"  optional:"true"`
	Type        string   `json:"type"        validate:"required"`
	Subjects    []string `json:"subjects"    validate:"required"`
	Effect      string   `json:"effect"                               optional:"true"`
	Resources   []string `json:"resources"   validate:"required"`
	Actions     []string `json:"actions"     validate:"required"`
	Status      string   `json:"status"                               optional:"true"`
	Owner       string   `json:"owner"                                optional:"true"`
	Meta        string   `json:"meta"                                 optional:"true"`
}

// UpdatePolicyRequest update policy request.
type UpdatePolicyRequest struct {
	Description string   `json:"description" validate:"min=1,max=30" optional:"true"`
	Type        string   `json:"type"        validate:"required"`
	Subjects    []string `json:"subjects"    validate:"required"`
	Effect      string   `json:"effect"                              optional:"true"`
	Resources   []string `json:"resources"   validate:"required"`
	Actions     []string `json:"actions"     validate:"required"`
	Meta        string   `json:"meta"                                optional:"true"`
}

// PolicyBase represents a policy restful resource.
type PolicyBase struct {
	// May add TypeMeta in the future.
	// metav1alpha1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1alpha1.ObjectMeta `json:"metadata,omitempty"`

	Subjects  []string `json:"subjects"`
	Resources []string `json:"resources"`
	Actions   []string `json:"actions"`
	Effect    string   `json:"effect"`

	Type        string `json:"type"`
	Status      string `json:"status"`
	Owner       string `json:"owner"`
	Description string `json:"description"`

	// casbin required
	Adapter     string     `json:"adapter"`
	Model       string     `json:"model"`
	PolicyRules [][]string `json:"policyRules"`
}

// DetailPolicyResponse policy detail.
type DetailPolicyResponse struct {
	PolicyBase
}

// PolicyList is the whole list of all policies which have been stored in stroage.
type PolicyList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// +optional
	metav1alpha1.ListMeta `json:",inline"`

	Items []*PolicyBase `json:"items"`
}
