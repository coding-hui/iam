// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"strings"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"

	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
	"github.com/coding-hui/common/util/idutil"
)

func init() {
	RegisterModel(&Policy{})
	RegisterModel(&Statement{})
}

// PolicyRule is used to determine which policy line to load.
type PolicyRule struct {
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

// Policy represent a policy model.
type Policy struct {
	// May add TypeMeta in the future.
	// metav1alpha1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1alpha1.ObjectMeta `json:"metadata,omitempty"`

	Subjects    pq.StringArray `json:"subjects"    gorm:"column:subjects;type:mediumtext"`
	Type        string         `json:"type"        gorm:"column:type;type:varchar(20)"`
	Status      string         `json:"status"      gorm:"column:status;type:varchar(20)"`
	Owner       string         `json:"owner"       gorm:"column:owner;type:varchar(100)"`
	Description string         `json:"description" gorm:"column:description;type:varchar(100)"`

	// casbin required
	Adapter    string      `json:"adapter"    gorm:"column:adapter;type:varchar(100)"`
	Statements []Statement `json:"statements" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName maps to mysql table name.
func (p *Policy) TableName() string {
	return TableNamePrefix + "policy"
}

// AfterCreate run after create database record.
func (p *Policy) AfterCreate(tx *gorm.DB) error {
	p.InstanceID = idutil.GetInstanceID(p.ID, "policy-")

	return tx.Save(p).Error
}

// GetPolicyRules get policy all casbin rules
func (p *Policy) GetPolicyRules() [][]string {
	var rules [][]string
	for _, sub := range p.Subjects {
		for _, statement := range p.Statements {
			for _, act := range statement.Actions {
				rules = append(rules, []string{sub, statement.ResourceIdentifier, strings.ToLower(act)})
			}
		}
	}

	return rules
}

// Statement resource policy statement.
type Statement struct {
	ID                 uint64         `json:"id"                 gorm:"primary_key;AUTO_INCREMENT;column:id"`
	PolicyId           uint64         `json:"policyId"           gorm:"column:policy_id;type:varchar(64)"`
	Effect             string         `json:"effect"             gorm:"column:effect;type:varchar(10)"`
	Resource           string         `json:"resource"           gorm:"column:resource;type:varchar(64)"`
	ResourceIdentifier string         `json:"resourceIdentifier" gorm:"column:resource_identifier;type:varchar(64)"`
	Actions            pq.StringArray `json:"actions"            gorm:"column:actions;type:mediumtext"`
}

// TableName maps to mysql table name.
func (s *Statement) TableName() string {
	return TableNamePrefix + "policy_statement"
}

// AllowAccess returns true if the policy effect is allow, otherwise false.
func (s *Statement) AllowAccess() bool {
	return s.Effect == v1alpha1.AllowAccess
}
