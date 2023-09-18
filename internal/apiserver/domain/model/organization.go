// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"gorm.io/gorm"

	metav1 "github.com/coding-hui/common/meta/v1"
	"github.com/coding-hui/common/util/idutil"
)

const (
	// DefaultOrganization default organization name.
	DefaultOrganization string = "built-in"
)

func init() {
	RegisterModel(&Organization{})
	RegisterModel(&DepartmentMember{})
}

// Organization represents a organization restful resource. It is also used as gorm model.
type Organization struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata,omitempty"`

	DisplayName string `json:"displayName" gorm:"column:display_name;varchar(100)"`
	WebsiteUrl  string `json:"websiteUrl"  gorm:"column:website_url;varchar(100)"`
	Favicon     string `json:"favicon"     gorm:"column:favicon;varchar(100)"`
	Disabled    bool   `json:"disabled"    gorm:"column:disabled;type:bool"`
	Description string `json:"description" gorm:"column:description;type:varchar(512)"`
}

// TableName maps to mysql table name.
func (o *Organization) TableName() string {
	return TableNamePrefix + "organization"
}

// AfterCreate run after create database record.
func (o *Organization) AfterCreate(tx *gorm.DB) error {
	o.InstanceID = idutil.GetInstanceID(o.ID, "org-")

	return tx.Save(o).Error
}

func (o *Organization) IsSystemBuiltIn() bool {
	return o.Name == DefaultOrganization
}

// DepartmentMember represents a organization member restful resource. It is also used as gorm model.
type DepartmentMember struct {
	DepartmentId string    `json:"departmentId"        gorm:"primary_key;column:department_id;type:varchar(64)"`
	MemberId     string    `json:"memberId"            gorm:"primary_key;column:member_id;type:varchar(64)"`
	CreatedAt    time.Time `json:"createdAt,omitempty" gorm:"column:created_at"`
}

// TableName maps to mysql table name.
func (o *DepartmentMember) TableName() string {
	return TableNamePrefix + "department_member"
}
