// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"time"
)

// SecretKey represents a secret key in the database.
type SecretKey struct {
	ID         string     `gorm:"primaryKey;column:id"      json:"id"`
	IdentityID string     `gorm:"column:identity_id;index"  json:"identity_id"`
	KeyID      string     `gorm:"column:key_id;uniqueIndex" json:"key_id"`
	SecretHash string     `gorm:"column:secret_hash"        json:"secret_hash"`
	Name       string     `gorm:"column:name"               json:"name"`
	ExpiresAt  *time.Time `gorm:"column:expires_at"         json:"expires_at"`
	CreatedAt  time.Time  `gorm:"column:created_at"         json:"created_at"`
}

// TableName returns the table name for SecretKey.
func (SecretKey) TableName() string {
	return "iam_secret_keys"
}
