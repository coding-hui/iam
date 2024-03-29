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
	RegisterModel(&Role{})
}

// Role represents a role restful resource. It is also used as gorm model.
type Role struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	DisplayName string `json:"displayName"     gorm:"column:display_name;type:varchar(100)"`
	Owner       string `json:"owner,omitempty" gorm:"column:owner;type:varchar(64)"`
	Disabled    bool   `json:"disabled"        gorm:"column:disabled;type:bool"`
	Description string `json:"description"     gorm:"column:description;type:varchar(512)"`

	Users []User `json:"-" gorm:"many2many:iam_user_role;references:instance_id"`
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
