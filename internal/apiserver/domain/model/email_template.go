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

func init() {
	RegisterModel(&EmailTemplate{})
	RegisterModel(&EmailTemplateCategory{})
	RegisterModel(&EmailTemplateVersion{})
	RegisterModel(&EmailTemplateVariable{})
}

// EmailTemplateStatus represents the status of an email template.
type EmailTemplateStatus string

const (
	// EmailTemplateStatusDraft indicates the template is in draft mode.
	EmailTemplateStatusDraft EmailTemplateStatus = "draft"
	// EmailTemplateStatusActive indicates the template is active and can be used.
	EmailTemplateStatusActive EmailTemplateStatus = "active"
	// EmailTemplateStatusDisabled indicates the template is disabled and cannot be used.
	EmailTemplateStatusDisabled EmailTemplateStatus = "disabled"
	// EmailTemplateStatusArchived indicates the template is archived and cannot be used.
	EmailTemplateStatusArchived EmailTemplateStatus = "archived"
)

// EmailTemplate represents an email template resource. It is also used as gorm model.
type EmailTemplate struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata"`

	// Subject is the email subject line.
	Subject string `json:"subject" gorm:"column:subject;type:varchar(200)"`

	// Content is the HTML content of the email template.
	Content string `json:"content" gorm:"column:content;type:text"`

	// PlainTextContent is the plain text version of the email.
	PlainTextContent string `json:"plainTextContent" gorm:"column:plain_text_content;type:text"`

	// Status indicates the current status of the template.
	Status EmailTemplateStatus `json:"status" gorm:"column:status;type:varchar(20);default:'draft'"`

	// CategoryID is the ID of the category this template belongs to.
	CategoryID string `json:"categoryId" gorm:"column:category_id;type:varchar(64);index"`

	// Owner is the user ID of the template owner.
	Owner string `json:"owner" gorm:"column:owner;type:varchar(64)"`

	// IsDefault indicates if this is a system default template.
	IsDefault bool `json:"isDefault" gorm:"column:is_default;type:bool;default:false"`

	// Description provides additional information about the template.
	Description string `json:"description" gorm:"column:description;type:varchar(512)"`

	// LastPublishedAt is the timestamp when this template was last published.
	LastPublishedAt *time.Time `json:"lastPublishedAt,omitempty" gorm:"column:last_published_at"`

	// Variables is a list of variables used in this template (not stored in DB).
	Variables []EmailTemplateVariable `json:"variables,omitempty" gorm:"-"`

	// Category is the template category (not stored in DB).
	Category *EmailTemplateCategory `json:"category,omitempty" gorm:"-"`

	// Versions is a list of historical versions of this template (not stored in DB).
	Versions []EmailTemplateVersion `json:"versions,omitempty" gorm:"-"`
}

// TableName maps to mysql table name.
func (t *EmailTemplate) TableName() string {
	return TableNamePrefix + "email_template"
}

// AfterCreate run after create database record.
func (t *EmailTemplate) AfterCreate(tx *gorm.DB) error {
	t.InstanceID = idutil.GetInstanceID(t.ID, "etpl-")
	return tx.Save(t).Error
}

// EmailTemplateCategoryType represents the type of email template category.
type EmailTemplateCategoryType string

const (
	// EmailTemplateCategoryTypeSystem indicates system-managed categories.
	EmailTemplateCategoryTypeSystem EmailTemplateCategoryType = "system"
	// EmailTemplateCategoryTypeCustom indicates user-created categories.
	EmailTemplateCategoryTypeCustom EmailTemplateCategoryType = "custom"
)

// EmailTemplateCategory represents a category for organizing email templates.
type EmailTemplateCategory struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata"`

	// Type indicates whether this is a system or custom category.
	Type EmailTemplateCategoryType `json:"type" gorm:"column:type;type:varchar(20);default:'custom'"`

	// ParentID is the ID of the parent category (for hierarchical organization).
	ParentID string `json:"parentId" gorm:"column:parent_id;type:varchar(64);index"`

	// Owner is the user ID of the category owner.
	Owner string `json:"owner" gorm:"column:owner;type:varchar(64)"`

	// Description provides additional information about the category.
	Description string `json:"description" gorm:"column:description;type:varchar(512)"`

	// Templates is a list of templates in this category (not stored in DB).
	Templates []EmailTemplate `json:"templates,omitempty" gorm:"-"`

	// Children is a list of child categories (not stored in DB).
	Children []EmailTemplateCategory `json:"children,omitempty" gorm:"-"`
}

// TableName maps to mysql table name.
func (c *EmailTemplateCategory) TableName() string {
	return TableNamePrefix + "email_template_category"
}

// AfterCreate run after create database record.
func (c *EmailTemplateCategory) AfterCreate(tx *gorm.DB) error {
	c.InstanceID = idutil.GetInstanceID(c.ID, "ecat-")
	return tx.Save(c).Error
}

// EmailTemplateVersion represents a historical version of an email template.
type EmailTemplateVersion struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata"`

	// TemplateID is the ID of the template this version belongs to.
	TemplateID string `json:"templateId" gorm:"column:template_id;type:varchar(64);index"`

	// VersionNumber is the sequential version number.
	VersionNumber int `json:"versionNumber" gorm:"column:version_number;type:int"`

	// Subject is the email subject line for this version.
	Subject string `json:"subject" gorm:"column:subject;type:varchar(200)"`

	// Content is the HTML content of this version.
	Content string `json:"content" gorm:"column:content;type:text"`

	// PlainTextContent is the plain text version of the email for this version.
	PlainTextContent string `json:"plainTextContent" gorm:"column:plain_text_content;type:text"`

	// PublishedBy is the user ID who published this version.
	PublishedBy string `json:"publishedBy" gorm:"column:published_by;type:varchar(64)"`

	// Comment is an optional comment about this version.
	Comment string `json:"comment" gorm:"column:comment;type:varchar(512)"`
}

// TableName maps to mysql table name.
func (v *EmailTemplateVersion) TableName() string {
	return TableNamePrefix + "email_template_version"
}

// AfterCreate run after create database record.
func (v *EmailTemplateVersion) AfterCreate(tx *gorm.DB) error {
	v.InstanceID = idutil.GetInstanceID(v.ID, "ever-")
	return tx.Save(v).Error
}

// EmailTemplateVariable represents a variable that can be used in email templates.
type EmailTemplateVariable struct {
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata"`

	// TemplateID is the ID of the template this variable belongs to.
	TemplateID string `json:"templateId" gorm:"column:template_id;type:varchar(64);index"`

	// DefaultValue is the default value for this variable.
	DefaultValue string `json:"defaultValue" gorm:"column:default_value;type:varchar(255)"`

	// Description explains the purpose of this variable.
	Description string `json:"description" gorm:"column:description;type:varchar(512)"`

	// Required indicates if this variable must be provided when using the template.
	Required bool `json:"required" gorm:"column:required;type:bool;default:false"`

	// Type indicates the data type of this variable (string, number, date, etc.).
	Type string `json:"type" gorm:"column:type;type:varchar(20);default:'string'"`
}

// TableName maps to mysql table name.
func (v *EmailTemplateVariable) TableName() string {
	return TableNamePrefix + "email_template_variable"
}

// AfterCreate run after create database record.
func (v *EmailTemplateVariable) AfterCreate(tx *gorm.DB) error {
	v.InstanceID = idutil.GetInstanceID(v.ID, "evar-")
	return tx.Save(v).Error
}
