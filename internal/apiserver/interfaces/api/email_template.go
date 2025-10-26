// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/coding-hui/common/errors"

	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/pkg/api"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	"github.com/coding-hui/iam/pkg/code"
)

// emailTemplate API handler
// Manages persistent email templates and categories
type emailTemplate struct {
	EmailTemplateService service.EmailTemplateService `inject:""`
}

// NewEmailTemplate constructs the API handler
func NewEmailTemplate() Interface { return &emailTemplate{} }

func (h *emailTemplate) RegisterApiGroup(g *gin.Engine) {
	tpl := g.Group(versionPrefix+"/email-templates").Use(autoAuthCheck.AuthFunc(), permissionCheckFunc("email-templates"))
	{
		tpl.GET("", h.listTemplates)
		tpl.POST("", h.createTemplate)
		tpl.GET(":instanceId", h.getTemplate)
		tpl.PUT(":instanceId", h.updateTemplate)
		tpl.DELETE(":instanceId", h.deleteTemplate)
	}

	cat := g.Group(versionPrefix+"/email-template-categories").Use(autoAuthCheck.AuthFunc(), permissionCheckFunc("email-template-categories"))
	{
		cat.GET("", h.listCategories)
		cat.POST("", h.createCategory)
		cat.GET(":instanceId", h.getCategory)
		cat.PUT(":instanceId", h.updateCategory)
		cat.DELETE(":instanceId", h.deleteCategory)
	}
}

// Template handlers
// @Tags         EmailTemplates
// @Summary      ListEmailTemplates
// @Description  List email templates
// @Param        query  query  v1.ListEmailTemplateOptions  false  "query parameters"
// @Success      200    {object}  api.Response{data=v1.EmailTemplateList}  "email templates"
// @Router       /api/v1/email-templates [get]
// @Security     BearerTokenAuth
func (h *emailTemplate) listTemplates(c *gin.Context) {
	var opts v1.ListEmailTemplateOptions
	if err := c.ShouldBindQuery(&opts); err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}
	list, err := h.EmailTemplateService.ListTemplates(c.Request.Context(), opts)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}
	api.OkWithData(list, c)
}

// @Tags         EmailTemplates
// @Summary      CreateEmailTemplate
// @Description  Create an email template
// @Accept       application/json
// @Product      application/json
// @Param        data  body  v1.CreateEmailTemplateRequest  true  "email template"
// @Success      200   {object}  api.Response{data=v1.EmailTemplateBase}  "created"
// @Router       /api/v1/email-templates [post]
// @Security     BearerTokenAuth
func (h *emailTemplate) createTemplate(c *gin.Context) {
	var req v1.CreateEmailTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}
	tpl, err := h.EmailTemplateService.CreateTemplate(c.Request.Context(), &req)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}
	api.OkWithData(tpl, c)
}

// @Tags         EmailTemplates
// @Summary      GetEmailTemplate
// @Description  Get an email template by instance ID
// @Param        instanceId  path  string  true  "identifier of an email template" default(template-001)
// @Success      200  {object}  api.Response{data=v1.EmailTemplateBase}  "email template"
// @Router       /api/v1/email-templates/{instanceId} [get]
// @Security     BearerTokenAuth
func (h *emailTemplate) getTemplate(c *gin.Context) {
	instanceId := c.Param("instanceId")
	tpl, err := h.EmailTemplateService.GetTemplate(c.Request.Context(), instanceId)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}
	api.OkWithData(tpl, c)
}

// @Tags         EmailTemplates
// @Summary      UpdateEmailTemplate
// @Description  Update an email template
// @Accept       application/json
// @Product      application/json
// @Param        instanceId  path  string  true  "identifier of an email template" default(template-001)
// @Param        data        body  v1.UpdateEmailTemplateRequest  true  "email template"
// @Success      200  {object}  api.Response{data=v1.EmailTemplateBase}  "updated"
// @Router       /api/v1/email-templates/{instanceId} [put]
// @Security     BearerTokenAuth
func (h *emailTemplate) updateTemplate(c *gin.Context) {
	instanceId := c.Param("instanceId")
	var req v1.UpdateEmailTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}
	tpl, err := h.EmailTemplateService.UpdateTemplate(c.Request.Context(), instanceId, &req)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}
	api.OkWithData(tpl, c)
}

// @Tags         EmailTemplates
// @Summary      DeleteEmailTemplate
// @Description  Delete an email template
// @Param        instanceId  path  string  true  "identifier of an email template" default(template-001)
// @Success      200  {object}  api.Response  "deleted"
// @Router       /api/v1/email-templates/{instanceId} [delete]
// @Security     BearerTokenAuth
func (h *emailTemplate) deleteTemplate(c *gin.Context) {
	instanceId := c.Param("instanceId")
	if err := h.EmailTemplateService.DeleteTemplate(c.Request.Context(), instanceId); err != nil {
		api.FailWithErrCode(err, c)
		return
	}
	api.Ok(c)
}

// Category handlers
// @Tags         EmailTemplateCategories
// @Summary      ListEmailTemplateCategories
// @Description  List email template categories
// @Param        query  query  v1.ListEmailTemplateOptions  false  "query parameters"
// @Success      200  {object}  api.Response{data=v1.EmailTemplateCategoryList}  "categories"
// @Router       /api/v1/email-template-categories [get]
// @Security     BearerTokenAuth
func (h *emailTemplate) listCategories(c *gin.Context) {
	var opts v1.ListEmailTemplateOptions // reuse ListOptions only
	if err := c.ShouldBindQuery(&opts); err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}
	list, err := h.EmailTemplateService.ListCategories(c.Request.Context(), opts.ListOptions)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}
	api.OkWithData(list, c)
}

// @Tags         EmailTemplateCategories
// @Summary      CreateEmailTemplateCategory
// @Description  Create an email template category
// @Accept       application/json
// @Product      application/json
// @Param        data  body  v1.CreateEmailTemplateCategoryRequest  true  "category"
// @Success      200  {object}  api.Response{data=v1.EmailTemplateCategoryBase}  "created"
// @Router       /api/v1/email-template-categories [post]
// @Security     BearerTokenAuth
func (h *emailTemplate) createCategory(c *gin.Context) {
	var req v1.CreateEmailTemplateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}
	cat, err := h.EmailTemplateService.CreateCategory(c.Request.Context(), &req)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}
	api.OkWithData(cat, c)
}

// @Tags         EmailTemplateCategories
// @Summary      GetEmailTemplateCategory
// @Description  Get an email template category by instance ID
// @Param        instanceId  path  string  true  "identifier of a category" default(category-001)
// @Success      200  {object}  api.Response{data=v1.EmailTemplateCategoryBase}  "category"
// @Router       /api/v1/email-template-categories/{instanceId} [get]
// @Security     BearerTokenAuth
func (h *emailTemplate) getCategory(c *gin.Context) {
	instanceId := c.Param("instanceId")
	cat, err := h.EmailTemplateService.GetCategory(c.Request.Context(), instanceId)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}
	api.OkWithData(cat, c)
}

// @Tags         EmailTemplateCategories
// @Summary      UpdateEmailTemplateCategory
// @Description  Update an email template category
// @Accept       application/json
// @Product      application/json
// @Param        instanceId  path  string  true  "identifier of a category" default(category-001)
// @Param        data        body  v1.UpdateEmailTemplateCategoryRequest  true  "category"
// @Success      200  {object}  api.Response{data=v1.EmailTemplateCategoryBase}  "updated"
// @Router       /api/v1/email-template-categories/{instanceId} [put]
// @Security     BearerTokenAuth
func (h *emailTemplate) updateCategory(c *gin.Context) {
	instanceId := c.Param("instanceId")
	var req v1.UpdateEmailTemplateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, "%s", err.Error()), c)
		return
	}
	cat, err := h.EmailTemplateService.UpdateCategory(c.Request.Context(), instanceId, &req)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}
	api.OkWithData(cat, c)
}

// @Tags         EmailTemplateCategories
// @Summary      DeleteEmailTemplateCategory
// @Description  Delete an email template category
// @Param        instanceId  path  string  true  "identifier of a category" default(category-001)
// @Success      200  {object}  api.Response  "deleted"
// @Router       /api/v1/email-template-categories/{instanceId} [delete]
// @Security     BearerTokenAuth
func (h *emailTemplate) deleteCategory(c *gin.Context) {
	instanceId := c.Param("instanceId")
	if err := h.EmailTemplateService.DeleteCategory(c.Request.Context(), instanceId); err != nil {
		api.FailWithErrCode(err, c)
		return
	}
	api.Ok(c)
}
