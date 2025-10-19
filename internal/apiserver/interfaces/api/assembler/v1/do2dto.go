// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1

import (
	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
	pb "github.com/coding-hui/iam/pkg/api/proto/apiserver/v1"

	metav1 "github.com/coding-hui/common/meta/v1"
)

// ConvertUserModelToBase assemble the User model to DTO.
func ConvertUserModelToBase(user *model.User) *v1.UserBase {
	return &v1.UserBase{
		ObjectMeta:    user.ObjectMeta,
		Status:        user.Status,
		Alias:         user.Alias,
		Email:         user.Email,
		Phone:         user.Phone,
		Avatar:        user.Avatar,
		UserType:      user.UserType,
		LastLoginTime: user.LastLoginTime,
		Disabled:      user.Disabled,
		DepartmentIds: user.DepartmentIds,
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

// ConvertModelToOrganizationBase assemble the Organization model to DTO.
func ConvertModelToOrganizationBase(org *model.Organization) *v1.OrganizationBase {
	return &v1.OrganizationBase{
		ObjectMeta:  org.ObjectMeta,
		DisplayName: org.DisplayName,
		WebsiteUrl:  org.WebsiteUrl,
		Favicon:     org.Favicon,
		Disabled:    org.Disabled,
		Description: org.Description,
		IsLeaf:      org.IsLeaf,
		ParentID:    org.ParentID,
		Owner:       org.Owner,
	}
}

// ConvertModelToApplicationBase assemble the Application model to DTO.
func ConvertModelToApplicationBase(app *model.Application) *v1.ApplicationBase {
	var identityProviders []v1.IdentityProviderBase
	if len(app.IdentityProviders) > 0 {
		for _, idp := range app.IdentityProviders {
			identityProviders = append(identityProviders, v1.IdentityProviderBase{
				ObjectMeta: metav1.ObjectMeta{
					Name:       idp.Name,
					InstanceID: idp.InstanceID,
					CreatedAt:  idp.CreatedAt,
					UpdatedAt:  idp.UpdatedAt,
				},
				Type:          idp.Type,
				Category:      idp.Category,
				MappingMethod: idp.MappingMethod,
				Status:        idp.Status,
				Owner:         idp.Owner,
				DisplayName:   idp.DisplayName,
				Description:   idp.Description,
				Config:        idp.Extend,
			})
		}
	}
	return &v1.ApplicationBase{
		ObjectMeta:        app.ObjectMeta,
		DisplayName:       app.DisplayName,
		Description:       app.Description,
		Owner:             app.Owner,
		Status:            app.Status,
		Logo:              app.Logo,
		HomepageUrl:       app.HomepageUrl,
		AppID:             app.AppID,
		AppSecret:         app.AppSecret,
		CallbackURL:       app.CallbackURL,
		LoginURL:          app.LoginURL,
		IdentityProviders: identityProviders,
	}
}

// ConvertModelToIdentityProviderBase assemble the IdentityProvider model to DTO.
func ConvertModelToIdentityProviderBase(idp *model.IdentityProvider) *v1.IdentityProviderBase {
	config := metav1.Extend{}
	config.Merge(idp.Extend.String())
	return &v1.IdentityProviderBase{
		ObjectMeta: metav1.ObjectMeta{
			InstanceID: idp.InstanceID,
			Name:       idp.Name,
			CreatedAt:  idp.CreatedAt,
			UpdatedAt:  idp.UpdatedAt,
		},
		Status:        idp.Status,
		DisplayName:   idp.DisplayName,
		Description:   idp.Description,
		Owner:         idp.Owner,
		Type:          idp.Type,
		Category:      idp.Category,
		MappingMethod: idp.MappingMethod,
		CallbackURL:   idp.CallbackURL,
		Config:        config,
	}
}

// ConvertApiKeyModelToBase assemble the ApiKey model to DTO.
func ConvertApiKeyModelToBase(apiKey *model.ApiKey) *v1.ApiKeyBase {
	base := &v1.ApiKeyBase{
		ObjectMeta:  apiKey.ObjectMeta,
		Name:        apiKey.Name,
		Key:         apiKey.Key,
		UserID:      apiKey.UserID,
		Status:      int(apiKey.Status),
		UsageCount:  apiKey.UsageCount,
		Description: apiKey.Description,
	}

	// Convert *time.Time to time.Time
	if apiKey.ExpiresAt != nil {
		base.ExpiresAt = *apiKey.ExpiresAt
	}

	if apiKey.LastUsedAt != nil {
		base.LastUsedAt = *apiKey.LastUsedAt
	}

	return base
}
