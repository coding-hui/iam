// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/coding-hui/common/errors"
)

// Validate validates the EmailTemplate fields.
func (t *EmailTemplate) Validate() error {
	if t.Name == "" {
		return errors.New("email template name cannot be empty")
	}

	if len(t.Name) > 100 {
		return errors.New("email template name cannot exceed 100 characters")
	}

	if t.Subject == "" {
		return errors.New("email template subject cannot be empty")
	}

	if len(t.Subject) > 200 {
		return errors.New("email template subject cannot exceed 200 characters")
	}

	if t.Content == "" {
		return errors.New("email template content cannot be empty")
	}

	// Validate status
	validStatus := map[EmailTemplateStatus]bool{
		EmailTemplateStatusDraft:    true,
		EmailTemplateStatusActive:   true,
		EmailTemplateStatusDisabled: true,
		EmailTemplateStatusArchived: true,
	}

	if !validStatus[t.Status] {
		return fmt.Errorf("invalid email template status: %s", t.Status)
	}

	return nil
}

// Validate validates the EmailTemplateCategory fields.
func (c *EmailTemplateCategory) Validate() error {
	if c.Name == "" {
		return errors.New("email template category name cannot be empty")
	}

	if len(c.Name) > 100 {
		return errors.New("email template category name cannot exceed 100 characters")
	}

	if len(c.Description) > 512 {
		return errors.New("email template category description cannot exceed 512 characters")
	}

	return nil
}

// Validate validates the EmailTemplateVersion fields.
func (v *EmailTemplateVersion) Validate() error {
	if v.TemplateID == "" {
		return errors.New("template ID cannot be empty")
	}

	if v.VersionNumber <= 0 {
		return errors.New("version number must be greater than 0")
	}

	if v.Subject == "" {
		return errors.New("email template version subject cannot be empty")
	}

	if len(v.Subject) > 200 {
		return errors.New("email template version subject cannot exceed 200 characters")
	}

	if v.Content == "" {
		return errors.New("email template version content cannot be empty")
	}

	if len(v.Comment) > 512 {
		return errors.New("email template version comment cannot exceed 512 characters")
	}

	return nil
}

// Validate validates the EmailTemplateVariable fields.
func (v *EmailTemplateVariable) Validate() error {
	if v.TemplateID == "" {
		return errors.New("template ID cannot be empty")
	}

	if v.Name == "" {
		return errors.New("variable name cannot be empty")
	}

	if len(v.Name) > 100 {
		return errors.New("variable name cannot exceed 100 characters")
	}

	// Variable names should follow a specific pattern (e.g., alphanumeric with dots)
	validNamePattern := regexp.MustCompile(`^[a-zA-Z0-9_]+(\.?[a-zA-Z0-9_]+)*$`)
	if !validNamePattern.MatchString(v.Name) {
		return errors.New("variable name must contain only alphanumeric characters, underscores, and dots (e.g., user.name)")
	}

	// Validate variable type
	validTypes := map[string]bool{
		"string":  true,
		"number":  true,
		"boolean": true,
		"date":    true,
		"object":  true,
		"array":   true,
	}

	if !validTypes[strings.ToLower(v.Type)] {
		return fmt.Errorf("invalid variable type: %s", v.Type)
	}

	if len(v.DefaultValue) > 255 {
		return errors.New("variable default value cannot exceed 255 characters")
	}

	if len(v.Description) > 512 {
		return errors.New("variable description cannot exceed 512 characters")
	}

	return nil
}
