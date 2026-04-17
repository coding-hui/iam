// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package identity

import "errors"

var (
	// ErrIdentityNotFound is returned when an identity is not found.
	ErrIdentityNotFound = errors.New("identity not found")

	// ErrCredentialsNotFound is returned when credentials are not found.
	ErrCredentialsNotFound = errors.New("credentials not found")

	// ErrInvalidCredentials is returned when credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrDuplicateCredentials is returned when credentials already exist.
	ErrDuplicateCredentials = errors.New("credentials already exist")

	// ErrIdentityAlreadyExists is returned when an identity already exists.
	ErrIdentityAlreadyExists = errors.New("identity already exists")
)
