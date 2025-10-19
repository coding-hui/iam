// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"fmt"

	"github.com/spf13/pflag"
)

// MailOptions defines configuration for mail service
type MailOptions struct {
	// Enabled specifies whether mail service is enabled
	Enabled bool `json:"enabled"   mapstructure:"enabled"`
	// Host is the SMTP server host
	Host string `json:"host"      mapstructure:"host"`
	// Port is the SMTP server port
	Port int `json:"port"      mapstructure:"port"`
	// Username is the SMTP username
	Username string `json:"username"  mapstructure:"username"`
	// Password is the SMTP password
	Password string `json:"password"  mapstructure:"password"`
	// From is the sender email address
	From string `json:"from"      mapstructure:"from"`
	// FromName is the sender name
	FromName string `json:"fromName"  mapstructure:"fromName"`
	// Templates contains email template configurations
	Templates *TemplateConfig `json:"templates" mapstructure:"templates"`
}

// TemplateConfig defines email template configurations
type TemplateConfig struct {
	// WelcomeEmailTemplate custom welcome email template
	WelcomeEmailTemplate string `json:"welcomeEmailTemplate"       mapstructure:"welcomeEmailTemplate"`
	// PasswordResetEmailTemplate custom password reset email template
	PasswordResetEmailTemplate string `json:"passwordResetEmailTemplate" mapstructure:"passwordResetEmailTemplate"`
}

// NewMailOptions create a `zero` value instance.
func NewMailOptions() *MailOptions {
	return &MailOptions{
		Enabled:  false,
		Host:     "",
		Port:     587,
		Username: "",
		Password: "",
		From:     "",
		FromName: "WeCoding IAM System",
		Templates: &TemplateConfig{
			WelcomeEmailTemplate:       "",
			PasswordResetEmailTemplate: "",
		},
	}
}

// Validate verifies flags passed to MailOptions.
func (o *MailOptions) Validate() []error {
	var errs []error

	if o.Enabled {
		if o.Host == "" {
			errs = append(errs, fmt.Errorf("--mail.host must be specified when mail is enabled"))
		}
		if o.Port <= 0 || o.Port > 65535 {
			errs = append(errs, fmt.Errorf("--mail.port must be between 1 and 65535"))
		}
		if o.Username == "" {
			errs = append(errs, fmt.Errorf("--mail.username must be specified when mail is enabled"))
		}
		if o.From == "" {
			errs = append(errs, fmt.Errorf("--mail.from must be specified when mail is enabled"))
		}
	}

	return errs
}

// AddFlags adds flags related to mail service for a specific APIServer to the specified FlagSet.
func (o *MailOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.BoolVar(&o.Enabled, "mail.enabled", o.Enabled, "Enable mail service for sending notifications")
	fs.StringVar(&o.Host, "mail.host", o.Host, "SMTP server host")
	fs.IntVar(&o.Port, "mail.port", o.Port, "SMTP server port")
	fs.StringVar(&o.Username, "mail.username", o.Username, "SMTP username")
	fs.StringVar(&o.Password, "mail.password", o.Password, "SMTP password")
	fs.StringVar(&o.From, "mail.from", o.From, "Sender email address")
	fs.StringVar(&o.FromName, "mail.fromName", o.FromName, "Sender name")
	fs.StringVar(&o.Templates.WelcomeEmailTemplate, "mail.templates.welcome-email-template", o.Templates.WelcomeEmailTemplate, "Custom welcome email template")
	fs.StringVar(
		&o.Templates.PasswordResetEmailTemplate,
		"mail.templates.password-reset-email-template",
		o.Templates.PasswordResetEmailTemplate,
		"Custom password reset email template",
	)
}
