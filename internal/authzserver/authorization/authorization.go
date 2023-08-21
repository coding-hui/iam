// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authorization

import authzv1 "github.com/coding-hui/iam/pkg/api/authzserver/v1"

// Authorization define the authorize interface that use local repository to
// authorize the subject access review.
type Authorization interface {
	// Authorize returns nil if subject s can perform action a on resource r with context c or an error otherwise.
	//  if err := guard.Authorize(&Request{Resource: "article/1234", Action: "update", Subject: "peter"}); err != nil {
	//    return errors.New("Not allowed")
	//  }
	Authorize(r *authzv1.Request) *authzv1.Response
}
