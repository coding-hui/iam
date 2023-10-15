// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"gorm.io/gorm"

	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"

	metav1 "github.com/coding-hui/common/meta/v1"
	"github.com/coding-hui/common/util/idutil"
)

func init() {
	RegisterModel(&IdentityProvider{})
}

// IdentityProvider represent a provider model.
type IdentityProvider struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Type          v1.IdentityProviderType `json:"type"          gorm:"column:type;type:varchar(64)"`
	Category      v1.ProviderCategory     `json:"category"      gorm:"column:category;type:varchar(20)"`
	MappingMethod v1.MappingMethod        `json:"mappingMethod" gorm:"column:mapping_method;type:varchar(20)"`

	CallbackURL string `json:"callbackURL"   gorm:"column:callback_url;type:varchar(200)"`
	Status      string `json:"status"        gorm:"column:status;type:varchar(20);default:0"`
	Owner       string `json:"owner"         gorm:"column:owner;type:varchar(100)"`
	DisplayName string `json:"displayName"   gorm:"column:display_name;type:varchar(100)"`
	Description string `json:"description"   gorm:"column:description;type:varchar(100)"`
}

// TableName maps to mysql table name.
func (p *IdentityProvider) TableName() string {
	return TableNamePrefix + "identity_provider"
}

// AfterCreate run after create database record.
func (p *IdentityProvider) AfterCreate(tx *gorm.DB) error {
	p.InstanceID = idutil.GetInstanceID(p.ID, "idp-")

	return tx.Save(p).Error
}
