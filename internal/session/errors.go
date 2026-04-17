// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package session

import "errors"

var (
	// ErrSessionNotFound is returned when a session is not found.
	ErrSessionNotFound = errors.New("session not found")

	// ErrSessionExpired is returned when a session has expired.
	ErrSessionExpired = errors.New("session expired")

	// ErrSessionInactive is returned when a session is inactive.
	ErrSessionInactive = errors.New("session inactive")
)
