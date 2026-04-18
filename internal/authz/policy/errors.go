// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package policy

import "errors"

var (
	// ErrPolicyNotFound is returned when a policy is not found.
	ErrPolicyNotFound = errors.New("policy not found")
)
