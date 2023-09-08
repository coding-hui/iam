// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	pb "github.com/coding-hui/iam/pkg/api/proto/apiserver/v1alpha1"
)

// ConvertUserModelToBase assemble the User model to DTO.
func ConvertUserModelToBase(user *model.User) *v1.UserBase {
	return &v1.UserBase{
		ObjectMeta:    user.ObjectMeta,
		TenantId:      user.TenantId,
		Status:        user.Status,
		Alias:         user.Alias,
		Email:         user.Email,
		Phone:         user.Phone,
		Avatar:        user.Avatar,
		UserType:      user.UserType,
		LastLoginTime: user.LastLoginTime,
		Disabled:      user.Disabled,
	}
}

// ConvertResourceModelToBase assemble the Resource model to DTO.
func ConvertResourceModelToBase(resource *model.Resource) *v1.ResourceBase {
	return &v1.ResourceBase{
		ObjectMeta:  resource.ObjectMeta,
		Method:      resource.Method,
		Api:         resource.Api,
		Type:        resource.Type,
		Description: resource.Description,
		Actions:     ConvertToActions(resource.Actions),
	}
}

// ConvertToActions assemble the Action model to DTO.
func ConvertToActions(actions []model.Action) []v1.Action {
	list := make([]v1.Action, 0, len(actions))
	for _, act := range actions {
		list = append(list, v1.Action{Name: act.Name, Description: act.Description})
	}

	return list
}

// ConvertRoleModelToBase assemble the Role model to DTO.
func ConvertRoleModelToBase(role *model.Role) *v1.RoleBase {
	return &v1.RoleBase{
		ObjectMeta:  role.ObjectMeta,
		Owner:       role.Owner,
		Description: role.Description,
		Disabled:    role.Disabled,
	}
}

// ConvertPolicyModelToBase assemble the Policy model to DTO.
func ConvertPolicyModelToBase(policy *model.Policy) *v1.PolicyBase {
	statements := make([]v1.Statement, 0, len(policy.Statements))
	for _, statement := range policy.Statements {
		statements = append(statements, v1.Statement{
			Effect:             statement.Effect,
			Resource:           statement.Resource,
			ResourceIdentifier: statement.ResourceIdentifier,
			Actions:            statement.Actions,
		})
	}
	return &v1.PolicyBase{
		ObjectMeta:  policy.ObjectMeta,
		Subjects:    policy.Subjects,
		Statements:  statements,
		Type:        policy.Type,
		Status:      policy.Status,
		Owner:       policy.Owner,
		Description: policy.Description,
		PolicyRules: policy.GetPolicyRules(),
	}
}

// ConvertPolicyModelToProtoInfo assemble the Policy to rpc info.
func ConvertPolicyModelToProtoInfo(policy *model.Policy) *pb.PolicyInfo {
	return &pb.PolicyInfo{
		Name:        policy.Name,
		Subjects:    policy.Subjects,
		Type:        policy.Type,
		Status:      policy.Status,
		Owner:       policy.Owner,
		Description: policy.Description,
		Adapter:     policy.Adapter,
	}
}

// ConvertOrganizationModelToBase assemble the Organization model to DTO.
func ConvertOrganizationModelToBase(org *model.Organization) *v1.OrganizationBase {
	return &v1.OrganizationBase{
		ObjectMeta:  org.ObjectMeta,
		DisplayName: org.DisplayName,
		WebsiteUrl:  org.WebsiteUrl,
		Favicon:     org.Favicon,
		Disabled:    org.Disabled,
		Description: org.Description,
	}
}
