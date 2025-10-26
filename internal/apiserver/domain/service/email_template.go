// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package service

import (
	"context"
	_ "embed"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/pkg/request"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/code"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"
)

const (
	// SystemCategoryName default system category name.
	SystemCategoryName = "System"
	// CustomCategoryName default custom category name.
	CustomCategoryName = "Custom"
)

// EmailTemplateService defines operations for email templates and categories.
// Handles basic CRUD; versioning and variable management can be extended later.
type EmailTemplateService interface {
	// Template operations
	CreateTemplate(ctx context.Context, req *v1.CreateEmailTemplateRequest) (*v1.EmailTemplateBase, error)
	UpdateTemplate(ctx context.Context, instanceId string, req *v1.UpdateEmailTemplateRequest) (*v1.EmailTemplateBase, error)
	DeleteTemplate(ctx context.Context, instanceId string) error
	GetTemplate(ctx context.Context, instanceId string) (*v1.EmailTemplateBase, error)
	ListTemplates(ctx context.Context, opts v1.ListEmailTemplateOptions) (*v1.EmailTemplateList, error)

	// Category operations
	CreateCategory(ctx context.Context, req *v1.CreateEmailTemplateCategoryRequest) (*v1.EmailTemplateCategoryBase, error)
	UpdateCategory(ctx context.Context, instanceId string, req *v1.UpdateEmailTemplateCategoryRequest) (*v1.EmailTemplateCategoryBase, error)
	DeleteCategory(ctx context.Context, instanceId string) error
	GetCategory(ctx context.Context, instanceId string) (*v1.EmailTemplateCategoryBase, error)
	ListCategories(ctx context.Context, opts metav1.ListOptions) (*v1.EmailTemplateCategoryList, error)

	// Initialization
	Init(ctx context.Context) error
}

// emailTemplateServiceImpl implements EmailTemplateService.
type emailTemplateServiceImpl struct {
	Store repository.Factory `inject:"repository"`
}

// NewEmailTemplateService returns a new EmailTemplateService implementation.
func NewEmailTemplateService() EmailTemplateService { return &emailTemplateServiceImpl{} }

// CreateTemplate creates a new email template.
func (s *emailTemplateServiceImpl) CreateTemplate(ctx context.Context, req *v1.CreateEmailTemplateRequest) (*v1.EmailTemplateBase, error) {
	// Get current user from context
	currentUser, ok := request.UserFrom(ctx)
	if !ok {
		return nil, errors.WithCode(code.ErrPermissionDenied, "Failed to obtain the current user")
	}

	tpl := &model.EmailTemplate{
		ObjectMeta:       metav1.ObjectMeta{Name: req.Name},
		Subject:          req.Subject,
		Content:          req.Content,
		PlainTextContent: req.PlainTextContent,
		Status:           model.EmailTemplateStatus(req.Status),
		CategoryID:       req.CategoryID,
		Owner:            currentUser.Name, // Use current user as owner
		Description:      req.Description,
		IsDefault:        req.IsDefault,
	}
	if err := tpl.Validate(); err != nil {
		return nil, err
	}

	if err := s.Store.EmailTemplateRepository().Create(ctx, tpl, metav1.CreateOptions{}); err != nil {
		return nil, err
	}
	return s.convertTemplateModelToBase(ctx, tpl), nil
}

// UpdateTemplate updates an existing email template.
func (s *emailTemplateServiceImpl) UpdateTemplate(ctx context.Context, instanceId string, req *v1.UpdateEmailTemplateRequest) (*v1.EmailTemplateBase, error) {
	existed, err := s.Store.EmailTemplateRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if req.Subject != "" {
		existed.Subject = req.Subject
	}
	if req.Content != "" {
		existed.Content = req.Content
	}
	if req.PlainTextContent != "" {
		existed.PlainTextContent = req.PlainTextContent
	}
	if req.Status != "" {
		existed.Status = model.EmailTemplateStatus(req.Status)
	}
	if req.CategoryID != "" {
		existed.CategoryID = req.CategoryID
	}
	// Owner should not be changed after creation
	if req.Description != "" {
		existed.Description = req.Description
	}
	existed.IsDefault = req.IsDefault || existed.IsDefault

	if err := existed.Validate(); err != nil {
		return nil, err
	}

	if err := s.Store.EmailTemplateRepository().Update(ctx, existed, metav1.UpdateOptions{}); err != nil {
		return nil, err
	}
	return s.convertTemplateModelToBase(ctx, existed), nil
}

// DeleteTemplate deletes an email template by instance ID.
func (s *emailTemplateServiceImpl) DeleteTemplate(ctx context.Context, instanceId string) error {
	// First, get the template to check if it's a system template
	tpl, err := s.Store.EmailTemplateRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Check if the template belongs to a System category
	if tpl.CategoryID != "" {
		category, err := s.Store.EmailTemplateCategoryRepository().GetByInstanceId(ctx, tpl.CategoryID, metav1.GetOptions{})
		if err != nil {
			// If category not found, allow deletion (orphaned template)
			if !errors.IsCode(err, code.ErrRecordNotExist) {
				return err
			}
		} else {
			// Prevent deletion of templates in System category
			if category.Type == model.EmailTemplateCategoryTypeSystem {
				return errors.WithCode(code.ErrCannotDeleteSystemEmailTemplate, "Cannot delete system email template")
			}
		}
	}

	// Also check if it's a default template (additional protection)
	if tpl.IsDefault {
		return errors.WithCode(code.ErrCannotDeleteDefaultEmailTemplate, "Cannot delete default email template")
	}

	return s.Store.EmailTemplateRepository().DeleteByInstanceId(ctx, instanceId, metav1.DeleteOptions{})
}

// GetTemplate retrieves an email template by instance ID.
func (s *emailTemplateServiceImpl) GetTemplate(ctx context.Context, instanceId string) (*v1.EmailTemplateBase, error) {
	tpl, err := s.Store.EmailTemplateRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return s.convertTemplateModelToBase(ctx, tpl), nil
}

// ListTemplates lists email templates with pagination and filters.
func (s *emailTemplateServiceImpl) ListTemplates(ctx context.Context, opts v1.ListEmailTemplateOptions) (*v1.EmailTemplateList, error) {
	// Forward v1 list options to metav1.ListOptions, preserving label/field filters.
	listOpts := opts.ListOptions
	tpls, err := s.Store.EmailTemplateRepository().List(ctx, listOpts)
	if err != nil {
		return nil, err
	}
	total, err := s.Store.EmailTemplateRepository().Count(ctx, listOpts)
	if err != nil {
		return nil, err
	}

	items := make([]*v1.EmailTemplateBase, 0, len(tpls))
	for i := range tpls {
		items = append(items, s.convertTemplateModelToBase(ctx, &tpls[i]))
	}

	return &v1.EmailTemplateList{
		ListMeta: metav1.ListMeta{TotalCount: total},
		Items:    items,
	}, nil
}

// CreateCategory creates a new email template category.
func (s *emailTemplateServiceImpl) CreateCategory(ctx context.Context, req *v1.CreateEmailTemplateCategoryRequest) (*v1.EmailTemplateCategoryBase, error) {
	// Get current user from context
	currentUser, ok := request.UserFrom(ctx)
	if !ok {
		return nil, errors.WithCode(code.ErrPermissionDenied, "Failed to obtain the current user")
	}

	cat := &model.EmailTemplateCategory{
		ObjectMeta:  metav1.ObjectMeta{Name: req.Name},
		ParentID:    req.ParentID,
		Owner:       currentUser.Name, // Use current user as owner
		Description: req.Description,
	}
	if err := cat.Validate(); err != nil {
		return nil, err
	}
	if err := s.Store.EmailTemplateCategoryRepository().Create(ctx, cat, metav1.CreateOptions{}); err != nil {
		return nil, err
	}
	return convertCategoryModelToBase(cat), nil
}

// UpdateCategory updates an existing email template category.
func (s *emailTemplateServiceImpl) UpdateCategory(ctx context.Context, instanceId string, req *v1.UpdateEmailTemplateCategoryRequest) (*v1.EmailTemplateCategoryBase, error) {
	cat, err := s.Store.EmailTemplateCategoryRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		cat.Name = req.Name
	}
	if req.ParentID != "" {
		cat.ParentID = req.ParentID
	}
	// Owner should not be changed after creation
	if req.Description != "" {
		cat.Description = req.Description
	}
	if err := cat.Validate(); err != nil {
		return nil, err
	}
	if err := s.Store.EmailTemplateCategoryRepository().Update(ctx, cat, metav1.UpdateOptions{}); err != nil {
		return nil, err
	}
	return convertCategoryModelToBase(cat), nil
}

// DeleteCategory deletes a category by instance ID.
func (s *emailTemplateServiceImpl) DeleteCategory(ctx context.Context, instanceId string) error {
	// First, get the category to check if it's a system category
	category, err := s.Store.EmailTemplateCategoryRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Prevent deletion of System category
	if category.Type == model.EmailTemplateCategoryTypeSystem {
		return errors.WithCode(code.ErrCannotDeleteSystemEmailTemplateCategory, "Cannot delete system email template category")
	}

	return s.Store.EmailTemplateCategoryRepository().DeleteByInstanceId(ctx, instanceId, metav1.DeleteOptions{})
}

// GetCategory retrieves a category by instance ID.
func (s *emailTemplateServiceImpl) GetCategory(ctx context.Context, instanceId string) (*v1.EmailTemplateCategoryBase, error) {
	cat, err := s.Store.EmailTemplateCategoryRepository().GetByInstanceId(ctx, instanceId, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return convertCategoryModelToBase(cat), nil
}

// ListCategories lists categories with pagination and filters.
func (s *emailTemplateServiceImpl) ListCategories(ctx context.Context, opts metav1.ListOptions) (*v1.EmailTemplateCategoryList, error) {
	cats, err := s.Store.EmailTemplateCategoryRepository().List(ctx, opts)
	if err != nil {
		return nil, err
	}
	total, err := s.Store.EmailTemplateCategoryRepository().Count(ctx, opts)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.EmailTemplateCategoryBase, 0, len(cats))
	for i := range cats {
		items = append(items, convertCategoryModelToBase(&cats[i]))
	}
	return &v1.EmailTemplateCategoryList{ListMeta: metav1.ListMeta{TotalCount: total}, Items: items}, nil
}

// Conversion helpers.
func (s *emailTemplateServiceImpl) convertTemplateModelToBase(ctx context.Context, t *model.EmailTemplate) *v1.EmailTemplateBase {
	if t == nil {
		return nil
	}

	base := &v1.EmailTemplateBase{
		ObjectMeta:       t.ObjectMeta,
		Subject:          t.Subject,
		Content:          t.Content,
		PlainTextContent: t.PlainTextContent,
		Status:           string(t.Status),
		CategoryID:       t.CategoryID,
		Owner:            t.Owner,
		Description:      t.Description,
		IsDefault:        t.IsDefault,
		LastPublishedAt:  t.LastPublishedAt,
	}

	// Load category information if CategoryID is present
	if t.CategoryID != "" {
		if category, err := s.Store.EmailTemplateCategoryRepository().GetByInstanceId(ctx, t.CategoryID, metav1.GetOptions{}); err == nil {
			base.Category = convertCategoryModelToBase(category)
		}
	}

	return base
}

func convertCategoryModelToBase(c *model.EmailTemplateCategory) *v1.EmailTemplateCategoryBase {
	if c == nil {
		return nil
	}
	return &v1.EmailTemplateCategoryBase{
		ObjectMeta:  c.ObjectMeta,
		Type:        string(c.Type),
		ParentID:    c.ParentID,
		Owner:       c.Owner,
		Description: c.Description,
	}
}

//go:embed mail/templates/welcome_email.html
var initWelcomeEmailTemplate string

//go:embed mail/templates/password_reset_email.html
var initPasswordResetEmailTemplate string

// Init initializes default email template categories and templates.
func (s *emailTemplateServiceImpl) Init(ctx context.Context) error {
	// 1) Ensure default categories exist
	type categoryDef struct {
		name         string
		description  string
		categoryType model.EmailTemplateCategoryType
	}

	defaultCategories := []categoryDef{
		{name: SystemCategoryName, description: "System default email templates", categoryType: model.EmailTemplateCategoryTypeSystem},
		{name: CustomCategoryName, description: "User custom email templates", categoryType: model.EmailTemplateCategoryTypeCustom},
	}

	var systemCategoryInstanceID string

	for _, catDef := range defaultCategories {
		cat, err := s.Store.EmailTemplateCategoryRepository().GetByName(ctx, catDef.name, metav1.GetOptions{})
		if err != nil {
			if errors.IsCode(err, code.ErrRecordNotExist) {
				// create category
				newCat := &model.EmailTemplateCategory{
					ObjectMeta:  metav1.ObjectMeta{Name: catDef.name},
					Type:        catDef.categoryType,
					ParentID:    "",
					Owner:       DefaultAdmin,
					Description: catDef.description,
				}
				if vErr := newCat.Validate(); vErr != nil {
					return vErr
				}
				if cErr := s.Store.EmailTemplateCategoryRepository().Create(ctx, newCat, metav1.CreateOptions{}); cErr != nil {
					return cErr
				}
				log.Infof("initialize email template category %s done", catDef.name)

				// Store System category instance ID for template creation
				if catDef.name == SystemCategoryName {
					systemCategoryInstanceID = newCat.InstanceID
				}
			} else {
				return err
			}
		} else {
			// Store System category instance ID for template creation
			if catDef.name == SystemCategoryName {
				systemCategoryInstanceID = cat.InstanceID
			}
		}
	}

	// 2) Seed default templates if not exist
	type seedTpl struct {
		name    string
		subject string
		content string
	}
	seeds := []seedTpl{
		{name: "WelcomeEmail", subject: "Welcome to WeCoding IAM System", content: initWelcomeEmailTemplate},
		{name: "PasswordResetEmail", subject: "Reset Your WeCoding IAM Password", content: initPasswordResetEmailTemplate},
	}

	for _, sdef := range seeds {
		// check existence
		_, err := s.Store.EmailTemplateRepository().GetByName(ctx, sdef.name, metav1.GetOptions{})
		if err != nil {
			if errors.IsCode(err, code.ErrRecordNotExist) {
				tpl := &model.EmailTemplate{
					ObjectMeta:       metav1.ObjectMeta{Name: sdef.name},
					Subject:          sdef.subject,
					Content:          sdef.content,
					PlainTextContent: "",
					Status:           model.EmailTemplateStatusActive,
					CategoryID:       systemCategoryInstanceID,
					Owner:            DefaultAdmin,
					Description:      "System default email template",
					IsDefault:        true,
				}
				if vErr := tpl.Validate(); vErr != nil {
					return vErr
				}
				if cErr := s.Store.EmailTemplateRepository().Create(ctx, tpl, metav1.CreateOptions{}); cErr != nil {
					return cErr
				}
				log.Infof("initialize email template %s done", sdef.name)
			} else {
				return err
			}
		}
	}

	return nil
}
