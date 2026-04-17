// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package sql

import (
	"time"
)

// Identity represents an identity in the database.
type Identity struct {
	ID        string    `gorm:"primaryKey;column:id" json:"id"`
	NetworkID string    `gorm:"column:nid;index"     json:"network_id"`
	SchemaID  string    `gorm:"column:schema_id"     json:"schema_id"`
	Traits    []byte    `gorm:"column:traits"        json:"traits"`
	CreatedAt time.Time `gorm:"column:created_at"    json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"    json:"updated_at"`
}

// TableName returns the table name for Identity.
func (Identity) TableName() string {
	return "iam_identities"
}
