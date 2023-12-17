// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/coding-hui/common/errors"
	metav1 "github.com/coding-hui/common/meta/v1"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/pkg/code"
	"github.com/coding-hui/iam/pkg/api"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

type policy struct {
	PolicyService service.PolicyService `inject:""`
}

// NewPolicy is the of policy.
func NewPolicy() Interface {
	return &policy{}
}

func (p *policy) RegisterApiGroup(g *gin.Engine) {
	apiv1 := g.Group(versionPrefix+"/policies").Use(autoAuthCheck.AuthFunc(), permissionCheckFunc("policies"))
	{
		apiv1.POST("", p.createPolicy)
		apiv1.PUT("/:instanceId", p.updatePolicy)
		apiv1.DELETE("/:instanceId", p.deletePolicy)
		apiv1.GET("/:instanceId", p.policyCheckFilter, p.detailPolicy)
		apiv1.GET("", p.listPolicies)
	}
}

//	@Tags			Policies
//	@Summary		CreatePolicy
//	@Description	Create a policy
//	@Accept			application/json
//	@Product		application/json
//	@Param			data	body		v1.CreatePolicyRequest	true	"policy rule"
//	@Success		200		{object}	api.Response			"Create a nnw policy"
//	@Router			/api/v1/policies [post]
//	@Security		BearerTokenAuth
//
// createPolicy create a new policy.
func (p *policy) createPolicy(c *gin.Context) {
	createReq := v1.CreatePolicyRequest{}
	err := c.ShouldBindJSON(&createReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := createReq.Validate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = p.PolicyService.CreatePolicy(c.Request.Context(), createReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Policies
//	@Summary		UpdatePolicy
//	@Description	Update a Policy
//	@Accept			application/json
//	@Product		application/json
//	@Param			data		body		v1.UpdatePolicyRequest	true	"Policy rule"
//	@Param			instanceId	path		string					true	"identifier of a policy"
//	@Success		200			{object}	api.Response			"Update policy info"
//	@Router			/api/v1/policies/{instanceId} [put]
//	@Security		BearerTokenAuth
//
// updatePolicy update policy info.
func (p *policy) updatePolicy(c *gin.Context) {
	updateReq := v1.UpdatePolicyRequest{}
	err := c.ShouldBindJSON(&updateReq)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	if errs := updateReq.ValidateUpdate(); errs != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrValidation, errs.ToAggregate().Error()), c)
		return
	}
	err = p.PolicyService.UpdatePolicy(c.Request.Context(), c.Param("instanceId"), updateReq)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Policies
//	@Summary		DeletePolicy
//	@Description	DeleteByInstanceId policy by instanceId
//	@Param			instanceId	path		string			true	"identifier of a policy"
//	@Success		200			{object}	api.Response	"Policy successfully deleted"
//	@Router			/api/v1/policies/{instanceId} [DELETE]
//	@Security		BearerTokenAuth
//
// deletePolicy delete policy by identifier.
func (p *policy) deletePolicy(c *gin.Context) {
	err := p.PolicyService.DeletePolicy(
		c.Request.Context(),
		c.Param("instanceId"),
		metav1.DeleteOptions{},
	)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.Ok(c)
}

//	@Tags			Policies
//	@Summary		GetPolicyInfo
//	@Description	GetByName a policy by name
//	@Param			instanceId	path		string										true	"name of a policy"
//	@Success		200			{object}	api.Response{data=v1.DetailPolicyResponse}	"Policy detail"
//	@Router			/api/v1/policies/{instanceId} [get]
//	@Security		BearerTokenAuth
//
// detailPolicy get policy detail info.
func (p *policy) detailPolicy(c *gin.Context) {
	policy := c.Request.Context().Value(&v1.CtxKeyPolicy).(*model.Policy)
	detail, err := p.PolicyService.DetailPolicy(c.Request.Context(), policy, metav1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithData(detail, c)
}

//	@Tags			Policies
//	@Summary		ListPolicies
//	@Description	List policies
//	@Param			name	query		string								false	"fuzzy search based on name"
//	@Param			offset	query		int									false	"query the page number"
//	@Param			limit	query		int									false	"query the page size number"
//	@Success		200		{object}	api.Response{data=v1.PolicyList}	"policies"
//	@Router			/api/v1/policies [get]
//	@Security		BearerTokenAuth
//
// listPolicies list policy page.
func (p *policy) listPolicies(c *gin.Context) {
	var opts metav1.ListOptions
	err := c.ShouldBindQuery(&opts)
	if err != nil {
		api.FailWithErrCode(errors.WithCode(code.ErrBind, err.Error()), c)
		return
	}
	resp, err := p.PolicyService.ListPolicies(c.Request.Context(), opts)
	if err != nil {
		api.FailWithErrCode(err, c)
		return
	}

	api.OkWithPage(resp.Items, resp.TotalCount, c)
}

func (p *policy) policyCheckFilter(c *gin.Context) {
	policy, err := p.PolicyService.GetPolicy(c.Request.Context(), c.Param("instanceId"), metav1.GetOptions{})
	if err != nil {
		api.FailWithErrCode(err, c)
		c.Abort()
		return
	}
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), &v1.CtxKeyPolicy, policy))
	c.Next()
}
