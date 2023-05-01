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
	RegisterModel(&Resource{}, &Action{})
}

// Resource represents api resource restful resource. It is also used as gorm model.
type Resource struct {
	// Standard object's metadata.
	metav1alpha1.ObjectMeta `json:"metadata,omitempty"`

	// Type resource type.
	Type string `json:"type" gorm:"column:type" validate:"required"`

	// Type unique identification of resource API.
	Api string `json:"api" gorm:"column:api" validate:"required"`

	// Method resource method.
	Method string `json:"method" gorm:"column:method" validate:"required"`

	IsDefault bool `json:"isDefault" gorm:"column:is_default"`

	// Description resource description.
	Description string `json:"description" gorm:"column:description"`

	// Actions resource access mode.
	Actions []Action `json:"actions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName maps to mysql table name.
func (r *Resource) TableName() string {
	return TableNamePrefix + "resource"
}

func (r *Resource) ShortTableName() string {
	return TableNamePrefix + "res"
}

// AfterCreate run after create database record.
func (r *Resource) AfterCreate(tx *gorm.DB) error {
	r.InstanceID = idutil.GetInstanceID(r.ID, "resource-")

	return tx.Save(r).Error
}

type Action struct {
	ID          uint64 `json:"id,omitempty"         gorm:"primary_key;AUTO_INCREMENT;column:id"`
	ResourceId  uint64 `json:"resourceId,omitempty" gorm:"column:resource_id"`
	Name        string `json:"name"                 gorm:"column:name"`
	Description string `json:"description"          gorm:"column:description"`
}

// TableName maps to mysql table name.
func (a *Action) TableName() string {
	return TableNamePrefix + "resource_actions"
}

func (a *Action) ShortTableName() string {
	return TableNamePrefix + "act"
}
