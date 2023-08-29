// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	metav1 "github.com/coding-hui/common/meta/v1"
	"github.com/coding-hui/common/util/auth"
	"github.com/coding-hui/common/util/idutil"
)

func init() {
	RegisterModel(&User{})
}

// User represents a user restful resource. It is also used as gorm model.
type User struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	TenantId      uint64     `json:"tenantId,omitempty"      gorm:"column:tenant_id;type:varchar(64)"`
	Status        int        `json:"status"                  gorm:"column:status;"`
	Alias         string     `json:"alias"                   gorm:"column:alias;type:varchar(64)"`
	Password      string     `json:"password,omitempty"      gorm:"column:password;type:varchar(256)"`
	Email         string     `json:"email"                   gorm:"column:email;type:varchar(64)"`
	Phone         string     `json:"phone"                   gorm:"column:phone;type:varchar(20)"`
	UserType      string     `json:"userType"                gorm:"column:user_type;type:varchar(20)"`
	Disabled      bool       `json:"disabled"                gorm:"column:disabled;type:bool"`
	LastLoginTime *time.Time `json:"lastLoginTime,omitempty" gorm:"column:last_login_time"`
}

// TableName maps to mysql table name.
func (u *User) TableName() string {
	return TableNamePrefix + "user"
}

// AfterCreate run after create database record.
func (u *User) AfterCreate(tx *gorm.DB) error {
	u.InstanceID = idutil.GetInstanceID(u.ID, "user-")

	return tx.Save(u).Error
}

// Compare with the plain text password. Returns true if it's the same as the encrypted one (in the `User` struct).
func (u *User) Compare(pwd string) error {
	if err := auth.Compare(u.Password, pwd); err != nil {
		return fmt.Errorf("failed to compile password: %w", err)
	}

	return nil
}
