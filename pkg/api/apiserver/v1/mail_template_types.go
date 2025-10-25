// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import "fmt"

// MailTemplate represents mail template configuration
//
//	@ValueObject
//	@Description	Mail template configuration
//
//	@Example		{
//	  "welcomeEmailTemplate": "<!DOCTYPE html>...",
//	  "passwordResetEmailTemplate": "<!DOCTYPE html>..."
//	}
type MailTemplate struct {
	// WelcomeEmailTemplate custom welcome email template
	WelcomeEmailTemplate string `json:"welcomeEmailTemplate,omitempty"       example:"<!DOCTYPE html>..."`
	// PasswordResetEmailTemplate custom password reset email template
	PasswordResetEmailTemplate string `json:"passwordResetEmailTemplate,omitempty" example:"<!DOCTYPE html>..."`
}

// MailTemplateResponse represents mail template response
//
//	@ValueObject
//	@Description	Mail template response
//
//	@Example		{
//	  "data": {
//	    "welcomeEmailTemplate": "<!DOCTYPE html>...",
//	    "passwordResetEmailTemplate": "<!DOCTYPE html>..."
//	  }
//	}
type MailTemplateResponse struct {
	Data *MailTemplate `json:"data"`
}

// PreviewMailTemplateRequest represents mail template preview request
//
//	@ValueObject
//	@Description	Mail template preview request
//
//	@Example		{
//	  "templateType": "welcome",
//	  "templateData": {
//	    "username": "testuser",
//	    "email": "test@example.com",
//	    "password": "initialpassword",
//	    "system": "WeCoding IAM 系统"
//	  }
//	}
type PreviewMailTemplateRequest struct {
	// TemplateType template type (welcome or password_reset)
	TemplateType string `json:"templateType" binding:"required,oneof=welcome password_reset" example:"welcome"`
	// TemplateData template data for preview
	TemplateData interface{} `json:"templateData" binding:"required"`
}

// Validate validates the preview request
func (r *PreviewMailTemplateRequest) Validate() error {
	if r.TemplateType == "" {
		return fmt.Errorf("templateType is required")
	}
	if r.TemplateType != "welcome" && r.TemplateType != "password_reset" {
		return fmt.Errorf("templateType must be 'welcome' or 'password_reset'")
	}
	if r.TemplateData == nil {
		return fmt.Errorf("templateData is required")
	}
	return nil
}

// PreviewMailTemplateResponse represents mail template preview response
//
//	@ValueObject
//	@Description	Mail template preview response
//
//	@Example		{
//	  "data": "<!DOCTYPE html>..."
//	}
type PreviewMailTemplateResponse struct {
	// Data rendered template content
	Data string `json:"data" example:"<!DOCTYPE html>..."`
}

// TestEmailRequest represents email test request
//
//	@ValueObject
//	@Description	Email test request
//
//	@Example		{
//	  "recipient": "test@example.com",
//	  "templateType": "welcome",
//	  "templateData": {
//	    "username": "testuser",
//	    "email": "test@example.com",
//	    "password": "initialpassword",
//	    "system": "WeCoding IAM 系统"
//	  }
//	}
type TestEmailRequest struct {
	// Recipient email address to send test email
	Recipient string `json:"recipient"    binding:"required,email"                        example:"test@example.com"`
	// TemplateType template type (welcome or password_reset)
	TemplateType string `json:"templateType" binding:"required,oneof=welcome password_reset" example:"welcome"`
	// TemplateData template data for the email
	TemplateData interface{} `json:"templateData" binding:"required"`
}

// Validate validates the test email request
func (r *TestEmailRequest) Validate() error {
	if r.Recipient == "" {
		return fmt.Errorf("recipient is required")
	}
	if r.TemplateType == "" {
		return fmt.Errorf("templateType is required")
	}
	if r.TemplateType != "welcome" && r.TemplateType != "password_reset" {
		return fmt.Errorf("templateType must be 'welcome' or 'password_reset'")
	}
	if r.TemplateData == nil {
		return fmt.Errorf("templateData is required")
	}
	return nil
}

// TestEmailResponse represents email test response
//
//	@ValueObject
//	@Description	Email test response
//
//	@Example		{
//	  "message": "Test email sent successfully",
//	  "recipient": "test@example.com",
//	  "templateType": "welcome"
//	}
type TestEmailResponse struct {
	// Message result message
	Message string `json:"message"      example:"Test email sent successfully"`
	// Recipient email address that received the test email
	Recipient string `json:"recipient"    example:"test@example.com"`
	// TemplateType template type used for the test
	TemplateType string `json:"templateType" example:"welcome"`
}
