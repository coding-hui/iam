// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"context"

	"gorm.io/gorm"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/apiserver/infrastructure/datastore"
	"github.com/coding-hui/iam/pkg/code"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
)

// ---------------------- EmailTemplate ----------------------

type emailTemplateRepositoryImpl struct {
	client *Client
}

func newEmailTemplateRepository(client *Client) repository.EmailTemplateRepository {
	return &emailTemplateRepositoryImpl{client: client}
}

// Create creates a new email template.
func (r *emailTemplateRepositoryImpl) Create(ctx context.Context, tpl *model.EmailTemplate, _ metav1.CreateOptions) error {
	if err := r.client.WithCtx(ctx).Create(&tpl).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return datastore.ErrRecordExist
		}
		return err
	}
	return nil
}

// Update updates an existing email template.
func (r *emailTemplateRepositoryImpl) Update(ctx context.Context, tpl *model.EmailTemplate, _ metav1.UpdateOptions) error {
	if err := r.client.WithCtx(ctx).Save(tpl).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return datastore.ErrRecordNotExist
		}
		return err
	}
	return nil
}

// DeleteByInstanceId deletes a template by instance ID.
func (r *emailTemplateRepositoryImpl) DeleteByInstanceId(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error {
	db := r.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	if err := db.Where("instance_id = ?", instanceId).Delete(&model.EmailTemplate{}).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return datastore.ErrRecordNotExist
		}
		return err
	}
	return nil
}

// GetByInstanceId retrieves a template by instance ID.
func (r *emailTemplateRepositoryImpl) GetByInstanceId(ctx context.Context, instanceId string, _ metav1.GetOptions) (*model.EmailTemplate, error) {
	var tpl *model.EmailTemplate
	if err := r.client.WithCtx(ctx).Where("instance_id = ?", instanceId).First(&tpl).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, datastore.ErrRecordNotExist
		}
		return nil, err
	}
	return tpl, nil
}

// GetByName retrieves a template by name.
func (r *emailTemplateRepositoryImpl) GetByName(ctx context.Context, name string, _ metav1.GetOptions) (*model.EmailTemplate, error) {
	var tpl *model.EmailTemplate
	if err := r.client.WithCtx(ctx).Where("name = ?", name).First(&tpl).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, datastore.ErrRecordNotExist
		}
		return nil, err
	}
	return tpl, nil
}

// List retrieves templates with pagination and filtering.
func (r *emailTemplateRepositoryImpl) List(ctx context.Context, opts metav1.ListOptions) ([]model.EmailTemplate, error) {
	var tpls []model.EmailTemplate
	if err := r.client.WithCtx(ctx).Model(model.EmailTemplate{}).
		Scopes(
			makeCondition(opts),
			paginate(opts),
		).
		Order("id desc").
		Find(&tpls).Error; err != nil {
		return nil, datastore.NewDBError(err, "failed to list email templates")
	}
	return tpls, nil
}

// Count returns the total count of templates.
func (r *emailTemplateRepositoryImpl) Count(ctx context.Context, opts metav1.ListOptions) (int64, error) {
	var count int64
	if err := r.client.WithCtx(ctx).Model(model.EmailTemplate{}).Scopes(makeCondition(opts)).Count(&count).Error; err != nil {
		return 0, datastore.NewDBError(err, "failed to count email templates")
	}
	return count, nil
}

// CreateVersion creates a new version for a template.
func (r *emailTemplateRepositoryImpl) CreateVersion(ctx context.Context, ver *model.EmailTemplateVersion, _ metav1.CreateOptions) error {
	// ensure template exists
	if ver.TemplateID == "" {
		return errors.WithCode(code.ErrEntityInvalid, "templateId is empty")
	}
	if err := r.client.WithCtx(ctx).Create(ver).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return datastore.ErrRecordExist
		}
		return err
	}
	return nil
}

// ListVersions lists versions of a template by template instance ID.
func (r *emailTemplateRepositoryImpl) ListVersions(ctx context.Context, templateInstanceId string, opts metav1.ListOptions) ([]model.EmailTemplateVersion, error) {
	var versions []model.EmailTemplateVersion
	if err := r.client.WithCtx(ctx).
		Model(model.EmailTemplateVersion{}).
		Where("template_id = ?", templateInstanceId).
		Scopes(
			makeCondition(opts),
			paginate(opts),
		).
		Order("version_number desc").
		Find(&versions).Error; err != nil {
		return nil, datastore.NewDBError(err, "failed to list template versions")
	}
	return versions, nil
}

// GetVersionByNumber gets a specific version by template instance ID and version number.
func (r *emailTemplateRepositoryImpl) GetVersionByNumber(
	ctx context.Context,
	templateInstanceId string,
	versionNumber int,
	_ metav1.GetOptions,
) (*model.EmailTemplateVersion, error) {
	var ver *model.EmailTemplateVersion
	if err := r.client.WithCtx(ctx).
		Where("template_id = ? AND version_number = ?", templateInstanceId, versionNumber).
		First(&ver).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, datastore.ErrRecordNotExist
		}
		return nil, err
	}
	return ver, nil
}

// DeleteVersion deletes a specific version.
func (r *emailTemplateRepositoryImpl) DeleteVersion(ctx context.Context, templateInstanceId string, versionNumber int, opts metav1.DeleteOptions) error {
	db := r.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	if err := db.Where("template_id = ? AND version_number = ?", templateInstanceId, versionNumber).
		Delete(&model.EmailTemplateVersion{}).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return datastore.ErrRecordNotExist
		}
		return err
	}
	return nil
}

// AddVariables batch creates variables for a template.
func (r *emailTemplateRepositoryImpl) AddVariables(ctx context.Context, vars []*model.EmailTemplateVariable, _ metav1.CreateOptions) error {
	if len(vars) == 0 {
		return nil
	}
	if err := r.client.WithCtx(ctx).CreateInBatches(&vars, 500).Error; err != nil {
		return datastore.NewDBError(err, "failed to add template variables")
	}
	return nil
}

// ListVariables lists variables for a template.
func (r *emailTemplateRepositoryImpl) ListVariables(ctx context.Context, templateInstanceId string, opts metav1.ListOptions) ([]model.EmailTemplateVariable, error) {
	var vars []model.EmailTemplateVariable
	if err := r.client.WithCtx(ctx).
		Model(model.EmailTemplateVariable{}).
		Where("template_id = ?", templateInstanceId).
		Scopes(
			makeCondition(opts),
			paginate(opts),
		).
		Order("id desc").
		Find(&vars).Error; err != nil {
		return nil, datastore.NewDBError(err, "failed to list template variables")
	}
	return vars, nil
}

// DeleteVariable deletes a variable by name under a template.
func (r *emailTemplateRepositoryImpl) DeleteVariable(ctx context.Context, templateInstanceId string, varName string, opts metav1.DeleteOptions) error {
	db := r.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	if err := db.Where("template_id = ? AND name = ?", templateInstanceId, varName).Delete(&model.EmailTemplateVariable{}).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return datastore.ErrRecordNotExist
		}
		return err
	}
	return nil
}

// DeleteVariablesByTemplate deletes all variables for a template.
func (r *emailTemplateRepositoryImpl) DeleteVariablesByTemplate(ctx context.Context, templateInstanceId string, opts metav1.DeleteOptions) error {
	db := r.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	if err := db.Where("template_id = ?", templateInstanceId).Delete(&model.EmailTemplateVariable{}).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return datastore.ErrRecordNotExist
		}
		return err
	}
	return nil
}

// ---------------------- EmailTemplateCategory ----------------------

type emailTemplateCategoryRepositoryImpl struct {
	client *Client
}

func newEmailTemplateCategoryRepository(client *Client) repository.EmailTemplateCategoryRepository {
	return &emailTemplateCategoryRepositoryImpl{client: client}
}

// Create creates a new template category.
func (r *emailTemplateCategoryRepositoryImpl) Create(ctx context.Context, cat *model.EmailTemplateCategory, _ metav1.CreateOptions) error {
	if err := r.client.WithCtx(ctx).Create(cat).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return datastore.ErrRecordExist
		}
		return err
	}
	return nil
}

// Update updates an existing template category.
func (r *emailTemplateCategoryRepositoryImpl) Update(ctx context.Context, cat *model.EmailTemplateCategory, _ metav1.UpdateOptions) error {
	if err := r.client.WithCtx(ctx).Save(cat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return datastore.ErrRecordNotExist
		}
		return err
	}
	return nil
}

// DeleteByInstanceId deletes a category by instance ID.
func (r *emailTemplateCategoryRepositoryImpl) DeleteByInstanceId(ctx context.Context, instanceId string, opts metav1.DeleteOptions) error {
	db := r.client.WithCtx(ctx)
	if opts.Unscoped {
		db = db.Unscoped()
	}
	if err := db.Where("instance_id = ?", instanceId).Delete(&model.EmailTemplateCategory{}).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return datastore.ErrRecordNotExist
		}
		return err
	}
	return nil
}

// GetByInstanceId retrieves a category by instance ID.
func (r *emailTemplateCategoryRepositoryImpl) GetByInstanceId(ctx context.Context, instanceId string, _ metav1.GetOptions) (*model.EmailTemplateCategory, error) {
	var cat *model.EmailTemplateCategory
	if err := r.client.WithCtx(ctx).Where("instance_id = ?", instanceId).First(&cat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, datastore.ErrRecordNotExist
		}
		return nil, err
	}
	return cat, nil
}

// GetByName retrieves a category by name.
func (r *emailTemplateCategoryRepositoryImpl) GetByName(ctx context.Context, name string, _ metav1.GetOptions) (*model.EmailTemplateCategory, error) {
	var cat *model.EmailTemplateCategory
	if err := r.client.WithCtx(ctx).Where("name = ?", name).First(&cat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, datastore.ErrRecordNotExist
		}
		return nil, err
	}
	return cat, nil
}

// List retrieves categories with pagination and filtering.
func (r *emailTemplateCategoryRepositoryImpl) List(ctx context.Context, opts metav1.ListOptions) ([]model.EmailTemplateCategory, error) {
	var cats []model.EmailTemplateCategory
	if err := r.client.WithCtx(ctx).Model(model.EmailTemplateCategory{}).
		Scopes(
			makeCondition(opts),
			paginate(opts),
		).
		Order("id desc").
		Find(&cats).Error; err != nil {
		return nil, datastore.NewDBError(err, "failed to list template categories")
	}
	return cats, nil
}

// Count returns the total count of categories.
func (r *emailTemplateCategoryRepositoryImpl) Count(ctx context.Context, opts metav1.ListOptions) (int64, error) {
	var count int64
	if err := r.client.WithCtx(ctx).Model(model.EmailTemplateCategory{}).Scopes(makeCondition(opts)).Count(&count).Error; err != nil {
		return 0, datastore.NewDBError(err, "failed to count template categories")
	}
	return count, nil
}

// ListChildren lists categories under a given parent category instance ID.
func (r *emailTemplateCategoryRepositoryImpl) ListChildren(ctx context.Context, parentInstanceId string, opts metav1.ListOptions) ([]model.EmailTemplateCategory, error) {
	var cats []model.EmailTemplateCategory
	if err := r.client.WithCtx(ctx).
		Model(model.EmailTemplateCategory{}).
		Where("parent_id = ?", parentInstanceId).
		Scopes(
			makeCondition(opts),
			paginate(opts),
		).
		Order("id desc").
		Find(&cats).Error; err != nil {
		return nil, datastore.NewDBError(err, "failed to list child template categories")
	}
	return cats, nil
}

// CountChildren counts categories under a given parent category instance ID.
func (r *emailTemplateCategoryRepositoryImpl) CountChildren(ctx context.Context, parentInstanceId string, opts metav1.ListOptions) (int64, error) {
	var count int64
	if err := r.client.WithCtx(ctx).
		Model(model.EmailTemplateCategory{}).
		Where("parent_id = ?", parentInstanceId).
		Scopes(makeCondition(opts)).
		Count(&count).Error; err != nil {
		return 0, datastore.NewDBError(err, "failed to count child template categories")
	}
	return count, nil
}
