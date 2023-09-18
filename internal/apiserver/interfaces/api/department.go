// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/gin-gonic/gin"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"

	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/apiserver/utils"
	"github.com/coding-hui/iam/internal/pkg/api"
	"github.com/coding-hui/iam/internal/pkg/code"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

type department struct {
	OrganizationService service.OrganizationService `inject:""`
}

// NewDepartment is the of department department.
func NewDepartment() Interface {
	return &department{}
}

func (d *department) RegisterApiGroup(g *gin.Engine) {
	apiv1 := g.Group(versionPrefix+"/departments").Use(autoAuthCheck.AuthFunc(), permissionCheckFunc("departments"))
	{
		apiv1.POST("", d.createDepartment)
		apiv1.PUT("/:instanceId", d.updateDepartment)
		apiv1.DELETE("/:instanceId", d.deleteDepartment)
		apiv1.GET("/:instanceId", d.getDepartment)
		apiv1.GET("", d.listDepartment)
		apiv1.GET("/:instanceId/disable", d.disableDepartment)
		apiv1.GET("/:instanceId/enable", d.enableDepartment)
		apiv1.POST("/:instanceId/member/add", d.addDepartmentMember)
		apiv1.POST("/:instanceId/member/batch_add", d.batchAddDepartmentMember)
		apiv1.POST("/:instanceId/member/remove", d.removeDepartmentMember)
		apiv1.POST("/:instanceId/member/batch_remove", d.batchRemoveDepartmentMember)
		apiv1.GET("/:instanceId/member", d.listDepartmentMembers)
	}
}

//	@Tags			Departments
//	@Summary		CreateDepartment
//	@Description	Create department
//	@Accept			application/json
//	@Product		application/json
//	@Param			data	body		v1.CreateDepartmentRequest	true	"department info"
//	@Success		200		{object}	api.Response				"create a new department"
//	@Router			/api/v1/departments [post]
//	@Security		BearerTokenAuth
//
// createDepartment create a new department.
func (d *department) createDepartment(c *gin.Context) {
	createReq := v1.CreateOrganizationRequest{}
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := createReq.Validate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = d.OrganizationService.CreateOrganization(c.Request.Context(), createReq, metav1.CreateOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Departments
//	@Summary		UpdateDepartment
//	@Description	Update department info
//	@Accept			application/json
//	@Product		application/json
//	@Param			instanceId	path		string						true	"identifier of a department"
//	@Param			data		body		v1.UpdateDepartmentRequest	true	"department info"
//	@Success		200			{object}	api.Response				"update department info"
//	@Router			/api/v1/departments/{instanceId} [put]
//	@Security		BearerTokenAuth
//
// updateDepartment update department info.
func (d *department) updateDepartment(c *gin.Context) {
	updateReq := v1.UpdateOrganizationRequest{}
	err := c.ShouldBindJSON(&updateReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := updateReq.ValidateUpdate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = d.OrganizationService.UpdateOrganization(c.Request.Context(), c.Param("instanceId"), updateReq, metav1.UpdateOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Departments
//	@Summary		DeleteDepartment
//	@Description	Delete department
//	@Param			instanceId	path		string			true	"identifier of a department"
//	@Success		200			{object}	api.Response	"delete department"
//	@Router			/api/v1/departments/{instanceId} [delete]
//	@Security		BearerTokenAuth
//
// deleteDepartment delete department by identifier.
func (d *department) deleteDepartment(c *gin.Context) {
	err := d.OrganizationService.DeleteOrganization(c.Request.Context(), c.Param("instanceId"), metav1.DeleteOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Departments
//	@Summary		GetDepartmentInfo
//	@Description	Get department info
//	@Param			instanceId	path		string									true	"identifier of a department"
//	@Success		200			{object}	api.Response{data=v1.OrganizationBase}	"department detail"
//	@Router			/api/v1/departments/{instanceId} [get]
//	@Security		BearerTokenAuth
//
// getDepartment get department detail.
func (d *department) getDepartment(c *gin.Context) {
	department, err := d.OrganizationService.GetOrganization(c.Request.Context(), c.Param("instanceId"), metav1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}
	resp, err := d.OrganizationService.DetailOrganization(c.Request.Context(), department)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(resp, c)
}

//	@Tags			Departments
//	@Summary		ListDepartments
//	@Description	List departments
//	@Param			offset	query		int										false	"query the page number"
//	@Param			limit	query		int										false	"query the page size number"
//	@Success		200		{object}	api.Response{data=v1.OrganizationList}	"departments"
//	@Router			/api/v1/departments [get]
//	@Security		BearerTokenAuth
//
// listDepartment list departments page.
func (d *department) listDepartment(c *gin.Context) {
	page, pageSize, err := utils.ExtractPagingParams(c, minPageSize, maxPageSize)
	if err != nil {
		api.Fail(c)
		return
	}
	resp, err := d.OrganizationService.ListOrganizations(c.Request.Context(), metav1.ListOptions{
		Limit:         &pageSize,
		Offset:        &page,
		FieldSelector: c.Query("fieldSelector"),
	})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(resp.Items, resp.TotalCount, c)
}

//	@Tags			Departments
//	@Summary		EnableDepartment
//	@Description	Enable a department
//	@Param			instanceId	path		string	true	"identifier of a department"
//	@Success		200			{object}	api.Response
//	@Router			/api/v1/departments/{instanceId}/disable [get]
//	@Security		BearerTokenAuth
//
// getDepartmentRoles disable department.
func (d *department) disableDepartment(c *gin.Context) {
	err := d.OrganizationService.DisableOrganization(c.Request.Context(), c.Param("instanceId"))
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Departments
//	@Summary		EnableDepartment
//	@Description	Enable a department
//	@Param			instanceId	path		string	true	"identifier of a department"
//	@Success		200			{object}	api.Response
//	@Router			/api/v1/departments/{instanceId}/enable [get]
//	@Security		BearerTokenAuth
//
// getDepartmentRoles enable department.
func (d *department) enableDepartment(c *gin.Context) {
	err := d.OrganizationService.EnableOrganization(c.Request.Context(), c.Param("instanceId"))
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Departments
//	@Summary		AddDepartmentMember
//	@Description	Add user to the department
//	@Param			instanceId	path		string	true	"identifier of a department"
//	@Success		200			{object}	api.Response
//	@Router			/api/v1/departments/{instanceId}/member/add [post]
//	@Security		BearerTokenAuth
//
// addDepartmentMember add user to the department.
func (d *department) addDepartmentMember(c *gin.Context) {
	addReq := v1.AddDepartmentMemberRequest{}
	err := c.ShouldBindJSON(&addReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := addReq.Validate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = d.OrganizationService.BatchAddDepartmentMembers(c.Request.Context(), c.Param("instanceId"),
		v1.BatchAddDepartmentMemberRequest{
			Members: []v1.DepartmentMember{addReq.DepartmentMember}},
	)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Departments
//	@Summary		BatchAddDepartmentMember
//	@Description	Batch add user to the department
//	@Param			instanceId	path		string	true	"identifier of a department"
//	@Success		200			{object}	api.Response
//	@Router			/api/v1/departments/{instanceId}/member/batch_add [post]
//	@Security		BearerTokenAuth
//
// batchAddDepartmentMember batch add user to the department.
func (d *department) batchAddDepartmentMember(c *gin.Context) {
	addReq := v1.BatchAddDepartmentMemberRequest{}
	err := c.ShouldBindJSON(&addReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := addReq.Validate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = d.OrganizationService.BatchAddDepartmentMembers(c.Request.Context(), c.Param("instanceId"), addReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Departments
//	@Summary		RemoveDepartmentMember
//	@Description	Remove user from the department
//	@Param			instanceId	path		string	true	"identifier of a department"
//	@Success		200			{object}	api.Response
//	@Router			/api/v1/departments/{instanceId}/member/remove [post]
//	@Security		BearerTokenAuth
//
// removeDepartmentMember remove user from the department.
func (d *department) removeDepartmentMember(c *gin.Context) {
	removeReq := v1.RemoveDepartmentMemberRequest{}
	err := c.ShouldBindJSON(&removeReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := removeReq.Validate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = d.OrganizationService.BatchRemoveDepartmentMembers(c.Request.Context(), c.Param("instanceId"),
		v1.BatchRemoveDepartmentMemberRequest{
			Members: []v1.DepartmentMember{removeReq.DepartmentMember}},
	)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Departments
//	@Summary		BatchRemoveDepartmentMember
//	@Description	Batch remove user from the department
//	@Param			instanceId	path		string	true	"identifier of a department"
//	@Success		200			{object}	api.Response
//	@Router			/api/v1/departments/{instanceId}/member/batch_remove [post]
//	@Security		BearerTokenAuth
//
// batchRemoveDepartmentMember batch remove user from the department.
func (d *department) batchRemoveDepartmentMember(c *gin.Context) {
	removeReq := v1.BatchRemoveDepartmentMemberRequest{}
	err := c.ShouldBindJSON(&removeReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := removeReq.Validate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = d.OrganizationService.BatchRemoveDepartmentMembers(c.Request.Context(), c.Param("instanceId"), removeReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Departments
//	@Summary		ListDepartmentMembers
//	@Description	List department members
//	@Param			offset	query		int											false	"query the page number"
//	@Param			limit	query		int											false	"query the page size number"
//	@Success		200		{object}	api.Response{data=v1.DepartmentMemberList}	"departments"
//	@Router			/api/v1/departments/{instanceId}/member [get]
//	@Security		BearerTokenAuth
//
// listDepartmentMembers list department members page.
func (d *department) listDepartmentMembers(c *gin.Context) {
	page, pageSize, err := utils.ExtractPagingParams(c, minPageSize, maxPageSize)
	if err != nil {
		api.Fail(c)
		return
	}
	resp, err := d.OrganizationService.ListDepartmentMembers(c.Request.Context(), c.Param("instanceId"),
		metav1.ListOptions{
			Limit:         &pageSize,
			Offset:        &page,
			FieldSelector: c.Query("fieldSelector"),
		},
	)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(resp.Members, resp.TotalCount, c)
}
