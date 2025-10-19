// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mail

import (
	"context"
	"fmt"

	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/log"

	"github.com/coding-hui/common/errors"
)

// TemplateService defines the mail template service interface
type TemplateService interface {
	// GetMailTemplates get current mail templates
	GetMailTemplates(ctx context.Context) (*v1.MailTemplate, error)
	// UpdateMailTemplates update mail templates
	UpdateMailTemplates(ctx context.Context, templates *v1.MailTemplate) error
	// PreviewMailTemplate preview mail template with sample data
	PreviewMailTemplate(ctx context.Context, req *v1.PreviewMailTemplateRequest) (string, error)
}

type templateServiceImpl struct {
}

// NewTemplateService creates a new mail template service
func NewTemplateService() TemplateService {
	return &templateServiceImpl{}
}

// GetMailTemplates get current mail templates
func (m *templateServiceImpl) GetMailTemplates(ctx context.Context) (*v1.MailTemplate, error) {
	// For now return empty templates since they are loaded from configuration
	// This could be extended to read from database or other storage in the future
	return &v1.MailTemplate{
		WelcomeEmailTemplate:       "",
		PasswordResetEmailTemplate: "",
	}, nil
}

// UpdateMailTemplates update mail templates
func (m *templateServiceImpl) UpdateMailTemplates(ctx context.Context, templates *v1.MailTemplate) error {
	// For now, templates are configured via configuration files
	// This could be extended to save to database or other storage in the future
	log.Warnf("Mail template update requested but currently templates are configured via configuration files only")
	return errors.New("Mail templates can only be configured via configuration files")
}

// PreviewMailTemplate preview mail template with sample data
func (m *templateServiceImpl) PreviewMailTemplate(ctx context.Context, req *v1.PreviewMailTemplateRequest) (string, error) {
	// Create a mock mail service for preview
	previewService := NewService(nil)

	// Get the template manager from the service
	templateManager := GetTemplateManager(previewService)
	if templateManager == nil {
		return "", errors.New("Failed to access template rendering methods")
	}

	switch req.TemplateType {
	case "welcome":
		// Extract template data for welcome email
		dataMap, ok := req.TemplateData.(map[string]interface{})
		if !ok {
			return "", errors.New("Invalid template data format for welcome email")
		}

		// Create welcome email data with default values
		welcomeData := &WelcomeEmailData{
			Username: getStringFromMap(dataMap, "username", "testuser"),
			Email:    getStringFromMap(dataMap, "email", "test@example.com"),
			Password: getStringFromMap(dataMap, "password", "initialpassword"),
			System:   getStringFromMap(dataMap, "system", "WeCoding IAM System"),
		}

		return templateManager.RenderWelcomeEmail(welcomeData)

	case "password_reset":
		// Extract template data for password reset email
		dataMap, ok := req.TemplateData.(map[string]interface{})
		if !ok {
			return "", errors.New("Invalid template data format for password reset email")
		}

		// Create password reset email data with default values
		resetData := &PasswordResetEmailData{
			Username:   getStringFromMap(dataMap, "username", "testuser"),
			Email:      getStringFromMap(dataMap, "email", "test@example.com"),
			ResetToken: getStringFromMap(dataMap, "resetToken", "reset-token-12345"),
			System:     getStringFromMap(dataMap, "system", "WeCoding IAM System"),
		}

		return templateManager.RenderPasswordResetEmail(resetData)

	default:
		return "", fmt.Errorf("unsupported template type: %s", req.TemplateType)
	}
}

// getStringFromMap safely gets string value from map with default fallback
func getStringFromMap(data map[string]interface{}, key string, defaultValue string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}
