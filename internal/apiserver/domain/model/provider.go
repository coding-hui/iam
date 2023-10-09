// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"gorm.io/gorm"

	metav1 "github.com/coding-hui/common/meta/v1"
	"github.com/coding-hui/common/util/idutil"
)

func init() {
	RegisterModel(&Provider{})
}

// Provider represent a provider model.
type Provider struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Type        string `json:"type"        gorm:"column:type;type:varchar(20)"`
	Category    string `json:"category"    gorm:"column:category;type:varchar(20)"`
	Status      string `json:"status"      gorm:"column:status;type:varchar(20);default:0"`
	Owner       string `json:"owner"       gorm:"column:owner;type:varchar(100)"`
	DisplayName string `json:"displayName" gorm:"column:display_name;type:varchar(100)"`
	Description string `json:"description" gorm:"column:description;type:varchar(100)"`
}

// TableName maps to mysql table name.
func (p *Provider) TableName() string {
	return TableNamePrefix + "provider"
}

// AfterCreate run after create database record.
func (p *Provider) AfterCreate(tx *gorm.DB) error {
	p.InstanceID = idutil.GetInstanceID(p.ID, "p-")

	return tx.Save(p).Error
}
