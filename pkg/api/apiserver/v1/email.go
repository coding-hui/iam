// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
    "time"

    metav1 "github.com/coding-hui/common/meta/v1"
)

// ListEmailTemplateOptions options for listing email templates
type ListEmailTemplateOptions struct {
    metav1.ListOptions `json:",inline"`
    CategoryID         string `json:"categoryId,omitempty" form:"categoryId" example:"category-001"`
    Status             string `json:"status,omitempty"     form:"status"     example:"active"`
}

// CreateEmailTemplateRequest request payload to create an email template
type CreateEmailTemplateRequest struct {
    Name             string `json:"name"             validate:"required,name" example:"welcome-email"`
    Subject          string `json:"subject"          validate:"required"      example:"Welcome to IAM"`
    Content          string `json:"content"          validate:"required"      example:"<html><body>Welcome, {{username}}</body></html>"`
    PlainTextContent string `json:"plainTextContent,omitempty"                 example:"Welcome, {{username}}"`
    Status           string `json:"status,omitempty"                           example:"active"`
    CategoryID       string `json:"categoryId,omitempty"                       example:"category-001"`
    Owner            string `json:"owner,omitempty"                            example:"ops@example.com"`
    Description      string `json:"description,omitempty"                      example:"Template for welcoming new users"`
    IsDefault        bool   `json:"isDefault,omitempty"                        example:"true"`
}

// UpdateEmailTemplateRequest request payload to update an email template
type UpdateEmailTemplateRequest struct {
    Subject          string `json:"subject,omitempty"          example:"Password Reset"`
    Content          string `json:"content,omitempty"          example:"<html>Reset instructions</html>"`
    PlainTextContent string `json:"plainTextContent,omitempty" example:"Reset instructions"`
    Status           string `json:"status,omitempty"           example:"draft"`
    CategoryID       string `json:"categoryId,omitempty"       example:"category-001"`
    Owner            string `json:"owner,omitempty"            example:"ops@example.com"`
    Description      string `json:"description,omitempty"      example:"Update description"`
    IsDefault        bool   `json:"isDefault,omitempty"        example:"false"`
}

// EmailTemplateBase base information of an email template
type EmailTemplateBase struct {
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Subject          string                     `json:"subject"                     example:"Welcome to IAM"`
    Content          string                     `json:"content"                     example:"<html><body>Welcome, {{username}}</body></html>"`
    PlainTextContent string                     `json:"plainTextContent,omitempty"  example:"Welcome, {{username}}"`
    Status           string                     `json:"status"                      example:"active"`
    CategoryID       string                     `json:"categoryId,omitempty"        example:"category-001"`
    Category         *EmailTemplateCategoryBase `json:"category,omitempty"`
    Owner            string                     `json:"owner,omitempty"             example:"ops@example.com"`
    Description      string                     `json:"description,omitempty"       example:"Template for welcoming new users"`
    IsDefault        bool                       `json:"isDefault"                   example:"true"`
    LastPublishedAt  *time.Time                 `json:"lastPublishedAt,omitempty"   example:"2025-01-15T12:34:56Z"`
}

// DetailEmailTemplateResponse email template detail response
type DetailEmailTemplateResponse struct {
    EmailTemplateBase `json:",inline"`
}

// EmailTemplateList email template list response
type EmailTemplateList struct {
    metav1.ListMeta `json:",inline"`
    Items []*EmailTemplateBase `json:"items"`
}

// CreateEmailTemplateCategoryRequest request payload to create a template category
type CreateEmailTemplateCategoryRequest struct {
    Name        string `json:"name"        validate:"required,name" example:"system-notifications"`
    ParentID    string `json:"parentId,omitempty"                    example:"root"`
    Owner       string `json:"owner,omitempty"                        example:"ops@example.com"`
    Description string `json:"description,omitempty"                  example:"Notifications for system events"`
}

// UpdateEmailTemplateCategoryRequest request payload to update a template category
type UpdateEmailTemplateCategoryRequest struct {
    Name        string `json:"name,omitempty"        example:"marketing"`
    ParentID    string `json:"parentId,omitempty"    example:"root"`
    Owner       string `json:"owner,omitempty"       example:"ops@example.com"`
    Description string `json:"description,omitempty" example:"Marketing templates"`
}

// EmailTemplateCategoryBase base information of a template category
type EmailTemplateCategoryBase struct {
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Type        string `json:"type,omitempty"        example:"system"`
    ParentID    string `json:"parentId,omitempty"    example:"root"`
    Owner       string `json:"owner,omitempty"       example:"ops@example.com"`
    Description string `json:"description,omitempty" example:"Notifications for system events"`
}

// EmailTemplateCategoryList template category list response
type EmailTemplateCategoryList struct {
    metav1.ListMeta `json:",inline"`
    Items []*EmailTemplateCategoryBase `json:"items"`
}