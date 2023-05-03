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
	RegisterModel(&Role{})
}

// Role represents a role restful resource. It is also used as gorm model.
type Role struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1alpha1.ObjectMeta `json:"metadata,omitempty"`

	DisplayName string `json:"displayName" gorm:"column:display_name;varchar(100)"`

	// Required: true
	Owner string `json:"owner,omitempty" gorm:"column:owner;varchar(64)" validate:"required"`

	Disabled bool `json:"disabled" gorm:"column:disabled"`

	Description string `json:"description" gorm:"column:description;varchar(255)"`

	Users []User `json:"-" gorm:"many2many:sys_user_role;references:instance_id"`
}

// TableName maps to mysql table name.
func (u *Role) TableName() string {
	return TableNamePrefix + "role"
}

// AfterCreate run after create database record.
func (u *Role) AfterCreate(tx *gorm.DB) error {
	u.InstanceID = idutil.GetInstanceID(u.ID, "role-")

	return tx.Save(u).Error
}
