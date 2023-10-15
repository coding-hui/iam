// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	"fmt"
	"strings"

	metav1 "github.com/coding-hui/common/meta/v1"
	"github.com/coding-hui/common/util/auth"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

// ConvertResourceModel assemble the DTO to Resource Model.
func ConvertResourceModel(req v1.CreateResourceRequest) *model.Resource {
	return &model.Resource{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
		Method:      req.Method,
		Type:        req.Type,
		Api:         req.Api,
		IsDefault:   req.IsDefault,
		Description: req.Description,
		Actions:     ConvertToActionModel(req.Name, req.Actions),
	}
}

// ConvertToActionModel assemble the DTO to Action Model.
func ConvertToActionModel(resource string, actions []v1.Action) []model.Action {
	list := make([]model.Action, 0, len(actions))
	for _, act := range actions {
		actName := act.Name
		if !strings.HasPrefix(actName, resource) {
			actName = fmt.Sprintf("%s:%s", resource, actName)
		}
		list = append(list, model.Action{Name: actName, Description: act.Description})
	}

	return list
}

// ConvertPolicyModel assemble the DTO to Policy Model.
func ConvertPolicyModel(req v1.CreatePolicyRequest) *model.Policy {
	return &model.Policy{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
		Subjects:    req.Subjects,
		Type:        req.Type,
		Status:      req.Status,
		Owner:       req.Owner,
		Description: req.Description,
		Statements:  ConvertToStatementModel(req.Statements),
	}
}

// ConvertToStatementModel assemble the DTO to Statements Model.
func ConvertToStatementModel(statements []v1.Statement) []model.Statement {
	list := make([]model.Statement, 0, len(statements))
	for _, sta := range statements {
		list = append(list, model.Statement{
			Effect:             sta.Effect,
			Resource:           sta.Resource,
			ResourceIdentifier: sta.ResourceIdentifier,
			Actions:            sta.Actions,
		})
	}

	return list
}

// ConvertCreateUserReqToUserModel assemble the create user request to User Model.
func ConvertCreateUserReqToUserModel(req v1.CreateUserRequest, external *model.UserExternal) *model.User {
	encryptPassword, _ := auth.Encrypt(req.Password)
	return &model.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
		Password: encryptPassword,
		Alias:    req.Alias,
		Email:    req.Email,
		Avatar:   req.Avatar,
		UserType: req.UserType,
		Disabled: false,
		External: external,
	}
}

// ConvertCreateDeptReqToModel assemble the create org dept request to Organization Model.
func ConvertCreateDeptReqToModel(req v1.CreateDepartmentRequest, parent *model.Organization) *model.Organization {
	return &model.Organization{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
		Ancestors:   parent.Ancestors + "," + parent.InstanceID,
		ParentID:    req.ParentID,
		DisplayName: req.DisplayName,
		WebsiteUrl:  req.WebsiteUrl,
		Favicon:     req.Favicon,
		Disabled:    false,
		IsLeaf:      false,
		Owner:       req.OrganizationID,
		Description: req.Description,
	}
}

// ConvertUpdateAppReqToModel assemble the update app request to Application Model.
func ConvertUpdateAppReqToModel(req v1.UpdateApplicationRequest, detail *v1.DetailApplicationResponse) *model.Application {
	app := &model.Application{
		ObjectMeta:  detail.ObjectMeta,
		Status:      detail.Status,
		Owner:       detail.Owner,
		Logo:        detail.Logo,
		HomepageUrl: detail.HomepageUrl,
		DisplayName: detail.DisplayName,
		Description: detail.Description,
	}
	if req.DisplayName != "" {
		app.DisplayName = req.DisplayName
	}
	if req.Status != "" {
		app.Status = req.Status
	}
	if req.Owner != "" {
		app.Owner = req.Owner
	}
	if req.Logo != "" {
		app.Logo = req.Logo
	}
	if req.HomepageUrl != "" {
		app.HomepageUrl = req.HomepageUrl
	}
	if req.Description != "" {
		app.Description = req.Description
	}
	return app
}
