// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"gorm.io/gorm"

	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
	"github.com/coding-hui/common/util/idutil"
)

func init() {
	//RegisterModel(&Tenant{})
}

// Tenant represents a tenant restful resource. It is also used as gorm model.
type Tenant struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1alpha1.ObjectMeta `json:"metadata,omitempty"`

	// Required: true
	Owner string `json:"owner,omitempty" gorm:"column:owner" validate:"required"`

	Disabled bool `json:"disabled" gorm:"column:disabled"`

	Description string `json:"description" gorm:"column:description"`
}

// TableName maps to mysql table name.
func (u *Tenant) TableName() string {
	return TableNamePrefix + "tenant"
}

func (u *Tenant) ShortTableName() string {
	return TableNamePrefix + "t"
}

// AfterCreate run after create database record.
func (u *Tenant) AfterCreate(tx *gorm.DB) error {
	u.InstanceID = idutil.GetInstanceID(u.ID, "tenant-")

	return tx.Save(u).Error
}