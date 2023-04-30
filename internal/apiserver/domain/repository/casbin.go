// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package repository

import (
	"github.com/casbin/casbin/v2"
)

// CasbinRepository defines the casbin repository interface.
type CasbinRepository interface {
	SyncedEnforcer() *casbin.SyncedEnforcer
}
