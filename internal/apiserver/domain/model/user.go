// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"
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
	metav1alpha1.ObjectMeta `json:"metadata,omitempty"`

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

// CustomClaims is the custom claims.
type CustomClaims struct {
	jwt.RegisteredClaims
	// Private Claim Names
	// TokenType defined the type of the token
	TokenType      string `json:"token_type,omitempty"`
	UserInstanceId string `json:"user_instance_id"`
	UserType       string `json:"user_type"`
	GrantType      string `json:"grant_type"`
}
