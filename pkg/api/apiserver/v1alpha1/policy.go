// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1alpha1

import (
	"encoding/json"

	"github.com/coding-hui/common/errors"
)

// AllowAccess should be used as effect for policies that allow access.
const AllowAccess = "allow"

// DenyAccess should be used as effect for policies that deny access.
const DenyAccess = "deny"

// CreatePolicyRequest create policy request.
type CreatePolicyRequest struct {
	Type        string     `json:"type"`
	Description string     `json:"description"`
	Subjects    []string   `json:"subjects"`
	Effect      string     `json:"effect"`
	Resources   []string   `json:"resources"`
	Actions     []string   `json:"actions"`
	Conditions  Conditions `json:"conditions"`
	Meta        []byte     `json:"meta"`
}

// Policy represent a policy model.
type Policy struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Description string     `json:"description"`
	Subjects    []string   `json:"subjects"`
	Effect      string     `json:"effect"`
	Resources   []string   `json:"resources"`
	Actions     []string   `json:"actions"`
	Conditions  Conditions `json:"conditions"`
	Meta        []byte     `json:"meta"`
}

// UnmarshalJSON overwrite own policy with values of the given in policy in JSON format
func (p *Policy) UnmarshalJSON(data []byte) error {
	var pol = struct {
		ID          string     `json:"id" gorethink:"id"`
		Type        string     `json:"type" gorethink:"type"`
		Description string     `json:"description" gorethink:"description"`
		Subjects    []string   `json:"subjects" gorethink:"subjects"`
		Effect      string     `json:"effect" gorethink:"effect"`
		Resources   []string   `json:"resources" gorethink:"resources"`
		Actions     []string   `json:"actions" gorethink:"actions"`
		Conditions  Conditions `json:"conditions" gorethink:"conditions"`
		Meta        []byte     `json:"meta" gorethink:"meta"`
	}{
		Conditions: Conditions{},
	}

	if err := json.Unmarshal(data, &pol); err != nil {
		return errors.WithStack(err)
	}

	*p = *&Policy{
		ID:          pol.ID,
		Type:        pol.Type,
		Description: pol.Description,
		Subjects:    pol.Subjects,
		Effect:      pol.Effect,
		Resources:   pol.Resources,
		Actions:     pol.Actions,
		Conditions:  pol.Conditions,
		Meta:        pol.Meta,
	}
	return nil
}

// UnmarshalMeta parses the policies []byte encoded metadata and stores the result in the value pointed to by v.
func (p *Policy) UnmarshalMeta(v interface{}) error {
	if err := json.Unmarshal(p.Meta, &v); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// GetID returns the policies id.
func (p *Policy) GetID() string {
	return p.ID
}

// GetType returns the policies type.
func (p *Policy) GetType() string {
	return p.Type
}

// GetDescription returns the policies description.
func (p *Policy) GetDescription() string {
	return p.Description
}

// GetSubjects returns the policies subjects.
func (p *Policy) GetSubjects() []string {
	return p.Subjects
}

// AllowAccess returns true if the policy effect is allow, otherwise false.
func (p *Policy) AllowAccess() bool {
	return p.Effect == AllowAccess
}

// GetEffect returns the policies effect which might be 'allow' or 'deny'.
func (p *Policy) GetEffect() string {
	return p.Effect
}

// GetResources returns the policies resources.
func (p *Policy) GetResources() []string {
	return p.Resources
}

// GetActions returns the policies actions.
func (p *Policy) GetActions() []string {
	return p.Actions
}

// GetConditions returns the policies conditions.
func (p *Policy) GetConditions() Conditions {
	return p.Conditions
}

// GetMeta returns the policies arbitrary metadata set by the user.
func (p *Policy) GetMeta() []byte {
	return p.Meta
}

// GetEndDelimiter returns the delimiter which identifies the end of a regular expression.
func (p *Policy) GetEndDelimiter() byte {
	return '>'
}

// GetStartDelimiter returns the delimiter which identifies the beginning of a regular expression.
func (p *Policy) GetStartDelimiter() byte {
	return '<'
}
