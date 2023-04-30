// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package code

//go:generate codegen -type=int

// iam-apiserver: user errors.
const (
	// ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 110001

	// ErrUserAlreadyExist - 400: User already exist.
	ErrUserAlreadyExist

	// ErrUserNameIsEmpty - 400: Username is empty.
	ErrUserNameIsEmpty

	// ErrDeleteOneself - 400: Unable to delete yourself.
	ErrDeleteOneself
)

// iam-apiserver: secret errors.
const (
	// ErrReachMaxCount - 400: Secret reach the max count.
	ErrReachMaxCount int = iota + 110101

	// ErrSecretNotFound - 404: Secret not found.
	ErrSecretNotFound
)

// iam-apiserver: policy errors.
const (
	// ErrPolicyNotFound - 404: Policy not found.
	ErrPolicyNotFound int = iota + 110201
)

// iam-apiserver: resource errors.
const (
	// ErrResourceNotFound - 404: Resource not found.
	ErrResourceNotFound int = iota + 110001

	// ErrResourceAlreadyExist - 400: Resource already exist.
	ErrResourceAlreadyExist

	// ErrResourceNameIsEmpty - 400: Resource name is empty.
	ErrResourceNameIsEmpty
)
