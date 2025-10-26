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
	RegisterModel(&Resource{}, &Action{})
}

// Resource represents api resource restful resource. It is also used as gorm model.
type Resource struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Type resource type.
	Type string `json:"type" gorm:"column:type;type:varchar(20)"`

	// Type unique identification of resource API.
	Api string `json:"api" gorm:"column:api;type:varchar(512)"`

	// Method resource method.
	Method string `json:"method" gorm:"column:method;type:varchar(10)"`

	IsDefault bool `json:"isDefault" gorm:"column:is_default;type:bool"`

	// Description resource description.
	Description string `json:"description" gorm:"column:description;type:varchar(512)"`

	// Actions resource access mode.
	Actions []Action `json:"actions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName maps to mysql table name.
func (r *Resource) TableName() string {
	return TableNamePrefix + "resource"
}

// AfterCreate run after create database record.
func (r *Resource) AfterCreate(tx *gorm.DB) error {
	r.InstanceID = idutil.GetInstanceID(r.ID, "resource-")

	return tx.Save(r).Error
}

// Action resource all operator actions.
type Action struct {
	ID          uint64 `json:"id"          gorm:"primary_key;AUTO_INCREMENT;column:id"`
	ResourceId  uint64 `json:"resourceId"  gorm:"column:resource_id;type:varchar(64)"`
	Name        string `json:"name"        gorm:"column:name;type:varchar(64)"`
	Description string `json:"description" gorm:"column:description;type:varchar(512)"`
}

// TableName maps to mysql table name.
func (a *Action) TableName() string {
	return TableNamePrefix + "resource_actions"
}

// BuiltInResourceSpec defines a system built-in API resource and its actions.
type BuiltInResourceSpec struct {
	Name        string
	Description string
	Actions     []string
	Api         string
}

// SystemBuiltInResources returns statically configured system built-in API resources.
func SystemBuiltInResources() []BuiltInResourceSpec {
	return []BuiltInResourceSpec{
		{
			Name:        "users",
			Description: "System built-in users management API resource",
			Actions:     []string{"GET", "POST", "PUT", "DELETE"},
			Api:         "/api/v1/users",
		},
		{
			Name:        "roles",
			Description: "System built-in roles management API resource",
			Actions:     []string{"GET", "POST", "PUT", "DELETE"},
			Api:         "/api/v1/roles",
		},
		{
			Name:        "resources",
			Description: "System built-in resources management API resource",
			Actions:     []string{"GET", "POST", "PUT", "DELETE"},
			Api:         "/api/v1/resources",
		},
		{
			Name:        "policies",
			Description: "System built-in policies management API resource",
			Actions:     []string{"GET", "POST", "PUT", "DELETE"},
			Api:         "/api/v1/policies",
		},
		{
			Name:        "organizations",
			Description: "System built-in organizations management API resource",
			Actions:     []string{"GET", "POST", "PUT", "DELETE"},
			Api:         "/api/v1/organizations",
		},
		{
			Name:        "departments",
			Description: "System built-in departments management API resource",
			Actions:     []string{"GET", "POST", "PUT", "DELETE"},
			Api:         "/api/v1/departments",
		},
		{
			Name:        "applications",
			Description: "System built-in applications management API resource",
			Actions:     []string{"GET", "POST", "PUT", "DELETE"},
			Api:         "/api/v1/applications",
		},
		{
			Name:        "identity-providers",
			Description: "System built-in identity-providers management API resource",
			Actions:     []string{"GET", "POST", "PUT", "DELETE"},
			Api:         "/api/v1/identity-providers",
		},
		{
			Name:        "apikeys",
			Description: "System built-in apikeys management API resource",
			Actions:     []string{"GET", "POST", "PUT", "DELETE"},
			Api:         "/api/v1/apikeys",
		},
		{
			Name:        "email-templates",
			Description: "System built-in email-templates management API resource",
			Actions:     []string{"GET", "POST", "PUT", "DELETE"},
			Api:         "/api/v1/email-templates",
		},
		{
			Name:        "email-template-categories",
			Description: "System built-in email-template-categories management API resource",
			Actions:     []string{"GET", "POST", "PUT", "DELETE"},
			Api:         "/api/v1/email-template-categories",
		},
	}
}
