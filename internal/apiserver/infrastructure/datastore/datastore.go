// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package datastore

import (
	"fmt"
)

var (
	// ErrGetClient Error that get datastore client.
	ErrGetClient = NewDBError(fmt.Errorf("get datastore client faield"))

	// ErrPrimaryEmpty Error that primary key is empty.
	ErrPrimaryEmpty = NewDBError(fmt.Errorf("entity primary is empty"))

	// ErrNilEntity Error that entity is nil
	ErrNilEntity = NewDBError(fmt.Errorf("entity is nil"))

	// ErrRecordExist Error that entity primary key is exist
	ErrRecordExist = NewDBError(fmt.Errorf("data record is exist"))

	// ErrRecordNotExist Error that entity primary key is not exist
	ErrRecordNotExist = NewDBError(fmt.Errorf("data record is not exist"))

	// ErrIndexInvalid Error that entity index is invalid
	ErrIndexInvalid = NewDBError(fmt.Errorf("entity index is invalid"))

	// ErrEntityInvalid Error that entity is invalid
	ErrEntityInvalid = NewDBError(fmt.Errorf("entity is invalid"))
)

// DBError datastore error
type DBError struct {
	err error
}

func (d *DBError) Error() string {
	return d.err.Error()
}

// NewDBError new datastore error
func NewDBError(err error) error {
	return &DBError{err: err}
}

// Config datastore config
type Config struct {
	Type     string
	URL      string
	Database string
}
