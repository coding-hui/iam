// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authz

import "errors"

var (
	// ErrNoMatchingPolicy is returned when no matching policy is found.
	ErrNoMatchingPolicy = errors.New("no matching policy")
)
