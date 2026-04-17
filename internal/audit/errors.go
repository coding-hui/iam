// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package audit

import "errors"

var (
	// ErrAuditNotFound is returned when audit events are not found.
	ErrAuditNotFound = errors.New("audit event not found")
)
