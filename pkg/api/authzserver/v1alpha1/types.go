// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package v1alpha1

import "encoding/json"

// Context is used as request's context.
type Context map[string]interface{}

// Request is the authz request object.
type Request struct {
	// Resource is the resource that access is requested to.
	Resource string `json:"resource"`

	// Action is the action that is requested on the resource.
	Action string `json:"action"`

	// Subejct is the subject that is requesting access.
	Subject string `json:"subject"`

	// Context is the request's environmental context.
	Context Context `json:"context"`
}

// Response represents the response of subject access review request.
type Response struct {
	Allowed bool   `json:"allowed"`
	Denied  bool   `json:"denied,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ToString marshal Response struct to a json string.
func (rsp *Response) ToString() string {
	data, _ := json.Marshal(rsp)

	return string(data)
}
