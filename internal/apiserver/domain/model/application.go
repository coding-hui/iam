// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"gorm.io/gorm"

	metav1 "github.com/coding-hui/common/meta/v1"
	"github.com/coding-hui/common/util/idutil"
)

const (
	// DefaultApplication default app name.
	DefaultApplication string = "built-in-app"
)

func init() {
	RegisterModel(&Application{})
}

// Application represent a application model.
type Application struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status      string `json:"status"      gorm:"column:status;type:varchar(20);default:0"`
	Owner       string `json:"owner"       gorm:"column:owner;type:varchar(100)"`
	Logo        string `json:"icon"        gorm:"column:logo;type:varchar(200)"`
	HomepageUrl string `json:"homepageUrl" gorm:"column:homepage_url;type:varchar(100)"`
	DisplayName string `json:"displayName" gorm:"column:display_name;type:varchar(100)"`
	Description string `json:"description" gorm:"column:description;type:varchar(100)"`
	AppID       string `json:"appId"       gorm:"column:app_id;type:varchar(64)"`
	AppSecret   string `json:"appSecret"   gorm:"column:app_secret;type:varchar(64)"`
	CallbackURL string `json:"callbackURL" gorm:"column:callback_url;type:varchar(200)"`
	LoginURL    string `json:"loginURL"    gorm:"column:login_url;type:varchar(200)"`

	IdentityProviders []IdentityProvider `json:"identityProviders" gorm:"many2many:iam_application_idp;references:instance_id"`
}

// TableName maps to mysql table name.
func (a *Application) TableName() string {
	return TableNamePrefix + "application"
}

// AfterCreate run after create database record.
func (a *Application) AfterCreate(tx *gorm.DB) error {
	a.InstanceID = idutil.GetInstanceID(a.ID, "app-")

	return tx.Save(a).Error
}
