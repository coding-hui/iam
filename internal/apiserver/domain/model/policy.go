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

	Subjects  pq.StringArray `json:"subjects"  gorm:"column:subjects;type:mediumtext"`
	Resources pq.StringArray `json:"resources" gorm:"column:resources;type:mediumtext"`
	Actions   pq.StringArray `json:"actions"   gorm:"column:actions;type:mediumtext"`
	Effect    string         `json:"effect"    gorm:"column:effect;type:varchar(10)"`

	Type        string `json:"type"        gorm:"column:type;type:varchar(20)"`
	Status      string `json:"status"      gorm:"column:status;type:varchar(20)"`
	Owner       string `json:"owner"       gorm:"column:owner;type:varchar(100)"`
	Description string `json:"description" gorm:"column:description;type:varchar(100)"`

	// casbin required
	Adapter string `json:"adapter" gorm:"column:adapter;type:varchar(100)"`
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

// AllowAccess returns true if the policy effect is allow, otherwise false.
func (p *Policy) AllowAccess() bool {
	return p.Effect == v1alpha1.AllowAccess
}

func (p *Policy) GetPolicyRules() [][]string {
	var rules [][]string
	for _, sub := range p.Subjects {
		for _, obj := range p.Resources {
			for _, act := range p.Actions {
				rules = append(rules, []string{sub, obj, strings.ToLower(act)})
			}
		}
	}

	return rules
}
