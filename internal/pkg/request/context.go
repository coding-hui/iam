// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package request

import (
	"context"

	v1 "github.com/coding-hui/iam/pkg/api/apiserver/v1"
)

// ctxKeyUser request context key of user.
var ctxKeyUser = "user"

// NewContext instantiates a base context object for request flows.
func NewContext() context.Context {
	return context.TODO()
}

// WithValue returns a copy of parent in which the value associated with key is val.
func WithValue(parent context.Context, key interface{}, val interface{}) context.Context {
	return context.WithValue(parent, key, val)
}

func WithUser(parent context.Context, user v1.UserBase) context.Context {
	return WithValue(parent, ctxKeyUser, user)
}

func UserFrom(ctx context.Context) (v1.UserBase, bool) {
	user, ok := ctx.Value(ctxKeyUser).(v1.UserBase)
	return user, ok
}
