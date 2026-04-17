// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"time"
)

// Token represents a token in the database.
type Token struct {
	ID         string    `gorm:"primaryKey;column:id"     json:"id"`
	IdentityID string    `gorm:"column:identity_id;index" json:"identity_id"`
	Type       string    `gorm:"column:type"              json:"type"` // access, refresh
	Value      string    `gorm:"column:value;uniqueIndex" json:"value"`
	ExpiresAt  time.Time `gorm:"column:expires_at"        json:"expires_at"`
	CreatedAt  time.Time `gorm:"column:created_at"        json:"created_at"`
}

// TableName returns the table name for Token.
func (Token) TableName() string {
	return "iam_tokens"
}
