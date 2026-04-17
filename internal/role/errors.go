// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package role

import "errors"

var (
	// ErrRoleNotFound is returned when a role is not found.
	ErrRoleNotFound = errors.New("role not found")
)
