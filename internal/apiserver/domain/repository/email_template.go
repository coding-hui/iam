// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"context"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"

	metav1 "github.com/coding-hui/common/meta/v1"
)

// EmailTemplateRepository defines repository operations for email templates,
// including managing versions and variables under a template aggregate.
type EmailTemplateRepository interface {
	// Template CRUD
	Create(ctx context.Context, tpl *model.EmailTemplate, opts metav1.CreateOptions) error
	Update(ctx context.Context, tpl *model.EmailTemplate, opts metav1.UpdateOptions) error
	DeleteByInstanceId(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error
	GetByInstanceId(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.EmailTemplate, error)
	GetByName(ctx context.Context, name string, opts metav1.GetOptions) (*model.EmailTemplate, error)
	List(ctx context.Context, opts metav1.ListOptions) ([]model.EmailTemplate, error)
	Count(ctx context.Context, opts metav1.ListOptions) (int64, error)

	// Versions
	CreateVersion(ctx context.Context, ver *model.EmailTemplateVersion, opts metav1.CreateOptions) error
	ListVersions(ctx context.Context, templateInstanceId string, opts metav1.ListOptions) ([]model.EmailTemplateVersion, error)
	GetVersionByNumber(ctx context.Context, templateInstanceId string, versionNumber int, opts metav1.GetOptions) (*model.EmailTemplateVersion, error)
	DeleteVersion(ctx context.Context, templateInstanceId string, versionNumber int, opts metav1.DeleteOptions) error

	// Variables
	AddVariables(ctx context.Context, vars []*model.EmailTemplateVariable, opts metav1.CreateOptions) error
	ListVariables(ctx context.Context, templateInstanceId string, opts metav1.ListOptions) ([]model.EmailTemplateVariable, error)
	DeleteVariable(ctx context.Context, templateInstanceId string, varName string, opts metav1.DeleteOptions) error
	DeleteVariablesByTemplate(ctx context.Context, templateInstanceId string, opts metav1.DeleteOptions) error
}

// EmailTemplateCategoryRepository defines repository operations for email template categories.
type EmailTemplateCategoryRepository interface {
	Create(ctx context.Context, cat *model.EmailTemplateCategory, opts metav1.CreateOptions) error
	Update(ctx context.Context, cat *model.EmailTemplateCategory, opts metav1.UpdateOptions) error
	DeleteByInstanceId(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error
	GetByInstanceId(ctx context.Context, instanceId string, opts metav1.GetOptions) (*model.EmailTemplateCategory, error)
	GetByName(ctx context.Context, name string, opts metav1.GetOptions) (*model.EmailTemplateCategory, error)
	List(ctx context.Context, opts metav1.ListOptions) ([]model.EmailTemplateCategory, error)
	Count(ctx context.Context, opts metav1.ListOptions) (int64, error)
	ListChildren(ctx context.Context, parentInstanceId string, opts metav1.ListOptions) ([]model.EmailTemplateCategory, error)
	CountChildren(ctx context.Context, parentInstanceId string, opts metav1.ListOptions) (int64, error)
}