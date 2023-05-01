// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1alpha1

import (
	metav1alpha1 "github.com/coding-hui/common/meta/v1alpha1"

	"github.com/coding-hui/iam/internal/apiserver/domain/model"
	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"
)

func CreateResourceModel(req v1alpha1.CreateResourceRequest) *model.Resource {
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

func ConvertToActionModel(actions []v1alpha1.Action) []model.Action {
	var list []model.Action
	for _, act := range actions {
		list = append(list, model.Action{Name: act.Name, Description: act.Description})
	}

	return list
}
