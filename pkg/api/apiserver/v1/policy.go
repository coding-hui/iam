// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	metav1 "github.com/coding-hui/common/meta/v1"
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
	// SystemBuildInPolicy system default policy type.
	SystemBuildInPolicy PolicyType = "SYSTEM"
	// CustomPolicy user custom policy type.
	CustomPolicy PolicyType = "CUSTOM"
)

type Statement struct {
	Effect             string   `json:"effect"             validate:"required"`
	Resource           string   `json:"resource"           validate:"required"`
	ResourceIdentifier string   `json:"resourceIdentifier" validate:"required"`
	Actions            []string `json:"actions"            validate:"required"`
}

// CreatePolicyRequest create policy request.
type CreatePolicyRequest struct {
	Name        string      `json:"name"        validate:"required,name"`
	Description string      `json:"description" validate:"min=1,max=30"  optional:"true"`
	Type        string      `json:"type"        validate:"required"`
	Statements  []Statement `json:"statements"  validate:"required"`
	Subjects    []string    `json:"subjects"    validate:"required"`
	Status      string      `json:"status"                               optional:"true"`
	Owner       string      `json:"owner"                                optional:"true"`
	Meta        string      `json:"meta"                                 optional:"true"`
}

// UpdatePolicyRequest update policy request.
type UpdatePolicyRequest struct {
	Description string      `json:"description" validate:"min=1,max=30" optional:"true"`
	Type        string      `json:"type"        validate:"required"`
	Subjects    []string    `json:"subjects"    validate:"required"`
	Statements  []Statement `json:"statements"  validate:"required"`
	Status      string      `json:"status"                              optional:"true"`
	Owner       string      `json:"owner"                               optional:"true"`
	Meta        string      `json:"meta"                                optional:"true"`
}

// PolicyBase represents a policy restful resource.
type PolicyBase struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Subjects   []string    `json:"subjects"`
	Statements []Statement `json:"statements"`

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
	Resources []ResourceBase `json:"resources,omitempty"`
}

// PolicyList is the whole list of all policies which have been stored in stroage.
type PolicyList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// +optional
	metav1.ListMeta `json:",inline"`

	Items []*PolicyBase `json:"items"`
}
