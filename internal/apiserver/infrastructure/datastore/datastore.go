// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package datastore

import (
	"github.com/coding-hui/iam/pkg/code"

	"github.com/coding-hui/common/errors"
)

var (
	// ErrPrimaryEmpty Error that primary key is empty.
	ErrPrimaryEmpty = errors.WithCode(code.ErrPrimaryEmpty, "entity primary is empty")

	// ErrTableNameEmpty Error that table name is empty.
	ErrTableNameEmpty = errors.WithCode(code.ErrTableNameEmpty, "entity table name is empty")

	// ErrNilEntity Error that entity is nil.
	ErrNilEntity = errors.WithCode(code.ErrNilEntity, "entity is nil")

	// ErrRecordExist Error that entity primary key is exist.
	ErrRecordExist = errors.WithCode(code.ErrRecordExist, "data record is exist")

	// ErrRecordNotExist Error that entity primary key is not exist.
	ErrRecordNotExist = errors.WithCode(code.ErrRecordNotExist, "data record is not exist")

	// ErrIndexInvalid Error that entity index is invalid.
	ErrIndexInvalid = errors.WithCode(code.ErrIndexInvalid, "entity index is invalid")

	// ErrEntityInvalid Error that entity is invalid.
	ErrEntityInvalid = errors.WithCode(code.ErrEntityInvalid, "entity is invalid")
)

// NewDBError new datastore error.
func NewDBError(err error, format string, args ...interface{}) error {
	return errors.WrapC(err, code.ErrDatabase, format, args...)
}
