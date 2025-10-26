// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/coding-hui/common/errors"

	"github.com/coding-hui/iam/internal/apiserver/domain/service/mail"
	"github.com/coding-hui/iam/pkg/api"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/code"
)

type mailTemplate struct {
	MailTemplateService mail.TemplateService `inject:""`
}

// NewMailTemplate is the of mail template.
func NewMailTemplate() Interface {
	return &mailTemplate{}
}

func (m *mailTemplate) RegisterApiGroup(g *gin.Engine) {
	apiv1 := g.Group(versionPrefix+"/mail-templates").Use(autoAuthCheck.AuthFunc(), permissionCheckFunc("mail-templates"))
	{
		apiv1.GET("", m.getMailTemplates)
		apiv1.PUT("", m.updateMailTemplates)
		apiv1.POST("/preview", m.previewMailTemplate)
		apiv1.POST("/test", m.sendTestEmail)
	}
}

//	@Tags			MailTemplates
//	@Summary		GetMailTemplates
//	@Description	Get current mail templates
//	@Accept			application/json
//	@Product		application/json
//	@Success		200	{object}	api.Response{data=v1.MailTemplateResponse}	"mail templates"
//	@Router			/api/v1/mail-templates [get]
//	@Security		BearerTokenAuth
//
// getMailTemplates get current mail templates.
func (m *mailTemplate) getMailTemplates(c *gin.Context) {
	templates, err := m.MailTemplateService.GetMailTemplates(c.Request.Context())
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(v1.MailTemplateResponse{
		Data: templates,
	}, c)
}

//	@Tags			MailTemplates
//	@Summary		UpdateMailTemplates
//	@Description	Update mail templates
//	@Accept			application/json
//	@Product		application/json
//	@Param			data	body		v1.MailTemplate	true	"mail templates"
//	@Success		200		{object}	api.Response	"update mail templates"
//	@Router			/api/v1/mail-templates [put]
//	@Security		BearerTokenAuth
//
// updateMailTemplates update mail templates.
func (m *mailTemplate) updateMailTemplates(c *gin.Context) {
	updateReq := v1.MailTemplate{}
	err := c.ShouldBindJSON(&updateReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}

	err = m.MailTemplateService.UpdateMailTemplates(c.Request.Context(), &updateReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			MailTemplates
//	@Summary		PreviewMailTemplate
//	@Description	Preview mail template with sample data
//	@Accept			application/json
//	@Product		application/json
//	@Param			data	body		v1.PreviewMailTemplateRequest						true	"preview request"
//	@Success		200		{object}	api.Response{data=v1.PreviewMailTemplateResponse}	"rendered template"
//	@Router			/api/v1/mail-templates/preview [post]
//	@Security		BearerTokenAuth
//
// previewMailTemplate preview mail template with sample data.
func (m *mailTemplate) previewMailTemplate(c *gin.Context) {
	previewReq := v1.PreviewMailTemplateRequest{}
	err := c.ShouldBindJSON(&previewReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}

	if err = previewReq.Validate(); err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, "%s", err.Error()), c)
		return
	}

	result, err := m.MailTemplateService.PreviewMailTemplate(c.Request.Context(), &previewReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(v1.PreviewMailTemplateResponse{
		Data: result,
	}, c)
}

//	@Tags			MailTemplates
//	@Summary		SendTestEmail
//	@Description	Send test email with specified template
//	@Accept			json
//	@Produce		json
//	@Param			data	body		v1.TestEmailRequest						true	"test email request"
//	@Success		200		{object}	api.Response{data=v1.TestEmailResponse}	"test email result"
//	@Router			/api/v1/mail-templates/test [post]
//
// sendTestEmail send test email with specified template
func (m *mailTemplate) sendTestEmail(c *gin.Context) {
	testReq := v1.TestEmailRequest{}
	if err := c.ShouldBindJSON(&testReq); err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}

	if err := testReq.Validate(); err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, "%s", err.Error()), c)
		return
	}

	if err := m.MailTemplateService.SendTestEmail(c.Request.Context(), &testReq); err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrUnknown, "%s", err.Error()), c)
		return
	}

	api.OkWithData(v1.TestEmailResponse{
		Message:      "Test email sent successfully",
		Recipient:    testReq.Recipient,
		TemplateType: testReq.TemplateType,
	}, c)
}
