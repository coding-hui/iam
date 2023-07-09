// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1alpha1

import (
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"
)

// ConvertResourceModel assemble the DTO to Resource Model.
func ConvertResourceModel(req v1alpha1.CreateResourceRequest) *model.Resource {
	return &model.Resource{
		ObjectMeta: metav1alpha1.ObjectMeta{
			Name: req.Name,
		},
		Method:      req.Method,
		Type:        req.Type,
		Api:         req.Api,
		IsDefault:   req.IsDefault,
		Description: req.Description,
		Actions:     ConvertToActionModel(req.Actions),
	}
}

// ConvertToActionModel assemble the DTO to Action Model.
func ConvertToActionModel(actions []v1alpha1.Action) []model.Action {
	list := make([]model.Action, 0, len(actions))
	for _, act := range actions {
		list = append(list, model.Action{Name: act.Name, Description: act.Description})
	}

	return list
}

// ConvertPolicyModel assemble the DTO to Policy Model.
func ConvertPolicyModel(req v1alpha1.CreatePolicyRequest) *model.Policy {
	return &model.Policy{
		ObjectMeta: metav1alpha1.ObjectMeta{
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
func ConvertToStatementModel(statements []v1alpha1.Statement) []model.Statement {
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
