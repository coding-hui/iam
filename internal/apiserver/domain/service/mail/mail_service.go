// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mail

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/pkg/log"
	pkgoptions "github.com/coding-hui/iam/pkg/options"

	metav1 "github.com/coding-hui/common/meta/v1"
	"github.com/spf13/viper"
)

//go:embed templates/welcome_email.html
var defaultWelcomeEmailTemplate string

//go:embed templates/password_reset_email.html
var defaultPasswordResetEmailTemplate string

// Service defines the mail service interface
type Service interface {
	SendWelcomeEmail(user *model.User, password string) error
	SendPasswordResetEmail(user *model.User, resetToken string) error
	SendTestEmail(recipient, templateType string, templateData interface{}) error
}

// TemplateManager defines the template manager interface
type TemplateManager interface {
	RenderWelcomeEmail(data *WelcomeEmailData) (string, error)
	RenderPasswordResetEmail(data *PasswordResetEmailData) (string, error)
}

// GetTemplateManager returns the template manager from a service instance
func GetTemplateManager(s Service) TemplateManager {
	if impl, ok := s.(*mailServiceImpl); ok {
		return impl.templateManager
	}
	return nil
}

// WelcomeEmailData represents the data for welcome email template
type WelcomeEmailData struct {
	Username string
	Email    string
	Password string
	System   string
}

// PasswordResetEmailData represents the data for password reset email template
type PasswordResetEmailData struct {
	Username   string
	Email      string
	ResetToken string
	System     string
}

// Config mail configuration
type Config struct {
	Enabled   bool
	Host      string
	Port      int
	Username  string
	Password  string
	From      string
	FromName  string
	Templates *TemplateConfig
}

// TemplateConfig template configuration
type TemplateConfig struct {
	WelcomeEmailTemplate       string `mapstructure:"welcomeEmailTemplate"`
	PasswordResetEmailTemplate string `mapstructure:"passwordResetEmailTemplate"`
}

// mailServiceImpl implements Service
type mailServiceImpl struct {
	smtpHost        string
	smtpPort        int
	smtpUsername    string
	smtpPassword    string
	fromEmail       string
	fromName        string
	templateManager TemplateManager
}

// templateManagerImpl implements TemplateManager
type templateManagerImpl struct {
	welcomeEmailTemplate         *template.Template
	passwordResetEmailTemplate   *template.Template
	defaultWelcomeTemplate       string
	defaultPasswordResetTemplate string
}

// NewService creates a new mail service
func NewService(cfg *Config) Service {
	if cfg == nil {
		cfg = &Config{}
	}

	templateManager := NewTemplateManager(cfg.Templates)

	return &mailServiceImpl{
		smtpHost:        cfg.Host,
		smtpPort:        cfg.Port,
		smtpUsername:    cfg.Username,
		smtpPassword:    cfg.Password,
		fromEmail:       cfg.From,
		fromName:        cfg.FromName,
		templateManager: templateManager,
	}
}

// NewServiceWithConfig creates a new mail service with configuration
func NewServiceWithConfig(c config.Config) Service {
	mailOpts := c.MailOptions()
	if mailOpts == nil || !mailOpts.Enabled {
		return NewService(&Config{Enabled: false}) // Disabled mail service
	}

	return NewService(&Config{
		Enabled:   mailOpts.Enabled,
		Host:      mailOpts.Host,
		Port:      mailOpts.Port,
		Username:  mailOpts.Username,
		Password:  mailOpts.Password,
		From:      mailOpts.From,
		FromName:  mailOpts.FromName,
		Templates: convertTemplateConfig(mailOpts.Templates),
	})
}

// NewTemplateManager creates a new template manager
func NewTemplateManager(templateConfig *TemplateConfig) TemplateManager {
	tm := &templateManagerImpl{
		defaultWelcomeTemplate:       defaultWelcomeEmailTemplate,
		defaultPasswordResetTemplate: defaultPasswordResetEmailTemplate,
	}

	// Load custom templates if provided
	if templateConfig != nil {
		tm.loadCustomTemplates(templateConfig)
	}

	return tm
}

// loadCustomTemplates loads custom templates from configuration
func (tm *templateManagerImpl) loadCustomTemplates(config *TemplateConfig) {
	if config.WelcomeEmailTemplate != "" {
		t, err := template.New("welcome").Parse(config.WelcomeEmailTemplate)
		if err != nil {
			log.Errorf("Failed to parse custom welcome email template: %v, using default", err)
		} else {
			tm.welcomeEmailTemplate = t
		}
	}

	if config.PasswordResetEmailTemplate != "" {
		t, err := template.New("password_reset").Parse(config.PasswordResetEmailTemplate)
		if err != nil {
			log.Errorf("Failed to parse custom password reset email template: %v, using default", err)
		} else {
			tm.passwordResetEmailTemplate = t
		}
	}
}

// RenderWelcomeEmail renders the welcome email template
func (tm *templateManagerImpl) RenderWelcomeEmail(data *WelcomeEmailData) (string, error) {
	return tm.renderTemplate(data, tm.welcomeEmailTemplate, tm.defaultWelcomeTemplate, "welcome")
}

// RenderPasswordResetEmail renders the password reset email template
func (tm *templateManagerImpl) RenderPasswordResetEmail(data *PasswordResetEmailData) (string, error) {
	return tm.renderTemplate(data, tm.passwordResetEmailTemplate, tm.defaultPasswordResetTemplate, "password_reset")
}

// renderTemplate renders email template with fallback to default
func (tm *templateManagerImpl) renderTemplate(data interface{}, customTemplate *template.Template, defaultTemplate, name string) (string, error) {
	var t *template.Template
	var err error

	if customTemplate != nil {
		t = customTemplate
	} else {
		t, err = template.New(name).Parse(defaultTemplate)
		if err != nil {
			return "", fmt.Errorf("failed to parse default %s email template: %w", name, err)
		}
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute %s email template: %w", name, err)
	}

	return buf.String(), nil
}

// SendWelcomeEmail sends a welcome email to the newly created user
func (m *mailServiceImpl) SendWelcomeEmail(user *model.User, password string) error {
	if user.Email == "" {
		log.Warnf("User %s has no email address, skip sending welcome email", user.Name)
		return nil
	}

	subject := "Welcome to WeCoding IAM System"
	data := &WelcomeEmailData{
		Username: user.Name,
		Email:    user.Email,
		Password: password,
		System:   "WeCoding IAM System",
	}

	body, err := m.templateManager.RenderWelcomeEmail(data)
	if err != nil {
		return fmt.Errorf("failed to render welcome email: %w", err)
	}

	return m.sendEmail(user.Email, subject, body)
}

// SendPasswordResetEmail sends a password reset email
func (m *mailServiceImpl) SendPasswordResetEmail(user *model.User, resetToken string) error {
	if user.Email == "" {
		log.Warnf("User %s has no email address, skip sending password reset email", user.Name)
		return nil
	}

	subject := "Reset Your WeCoding IAM Password"
	data := &PasswordResetEmailData{
		Username:   user.Name,
		Email:      user.Email,
		ResetToken: resetToken,
		System:     "WeCoding IAM System",
	}

	body, err := m.templateManager.RenderPasswordResetEmail(data)
	if err != nil {
		return fmt.Errorf("failed to render password reset email: %w", err)
	}

	return m.sendEmail(user.Email, subject, body)
}

// sendEmail sends an email using SMTP
func (m *mailServiceImpl) sendEmail(to, subject, body string) error {
	if m.smtpHost == "" || m.smtpPort == 0 {
		log.Warn("SMTP configuration not set, skip sending email")
		return nil
	}

	// Set default values if not provided
	fromEmail := m.fromEmail
	if fromEmail == "" {
		fromEmail = m.smtpUsername
	}
	fromName := m.fromName
	if fromName == "" {
		fromName = "WeCoding IAM System"
	}

	// Create message
	message := fmt.Sprintf("From: %s <%s>\r\n", fromName, fromEmail)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n"
	message += "\r\n" + body

	// Authentication
	auth := smtp.PlainAuth("", m.smtpUsername, m.smtpPassword, m.smtpHost)

	// Send email
	addr := fmt.Sprintf("%s:%d", m.smtpHost, m.smtpPort)
	err := smtp.SendMail(addr, auth, fromEmail, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	log.Infof("Successfully sent welcome email to %s", to)
	return nil
}

// SendTestEmail sends a test email using the specified template
func (m *mailServiceImpl) SendTestEmail(recipient, templateType string, templateData interface{}) error {
	if recipient == "" {
		return fmt.Errorf("recipient email is required")
	}

	// Extract template data based on template type
	dataMap, ok := templateData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid template data format")
	}

	switch templateType {
	case "welcome":
		// Create mock user for welcome email
		mockUser := &model.User{
			ObjectMeta: metav1.ObjectMeta{
				Name: getStringFromMap(dataMap, "username", "testuser"),
			},
			Email: recipient,
		}
		password := getStringFromMap(dataMap, "password", "initialpassword")
		return m.SendWelcomeEmail(mockUser, password)

	case "password_reset":
		// Create mock user for password reset email
		mockUser := &model.User{
			ObjectMeta: metav1.ObjectMeta{
				Name: getStringFromMap(dataMap, "username", "testuser"),
			},
			Email: recipient,
		}
		resetToken := getStringFromMap(dataMap, "resetToken", "test-reset-token")
		return m.SendPasswordResetEmail(mockUser, resetToken)

	default:
		return fmt.Errorf("unsupported template type: %s", templateType)
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

// convertTemplateConfig converts pkg/options TemplateConfig to domain TemplateConfig
func convertTemplateConfig(pkgConfig *pkgoptions.TemplateConfig) *TemplateConfig {
	if pkgConfig == nil {
		return nil
	}

	return &TemplateConfig{
		WelcomeEmailTemplate:       pkgConfig.WelcomeEmailTemplate,
		PasswordResetEmailTemplate: pkgConfig.PasswordResetEmailTemplate,
	}
}

// LoadConfigFromViper loads mail configuration from viper
func LoadConfigFromViper(v *viper.Viper) *Config {
	templateConfig := &TemplateConfig{
		WelcomeEmailTemplate:       v.GetString("mail.templates.welcomeEmailTemplate"),
		PasswordResetEmailTemplate: v.GetString("mail.templates.passwordResetEmailTemplate"),
	}

	return &Config{
		Host:      v.GetString("mail.host"),
		Port:      v.GetInt("mail.port"),
		Username:  v.GetString("mail.username"),
		Password:  v.GetString("mail.password"),
		From:      v.GetString("mail.from"),
		FromName:  v.GetString("mail.fromName"),
		Templates: templateConfig,
	}
}
