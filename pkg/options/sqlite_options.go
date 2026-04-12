// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"github.com/spf13/pflag"
)

// SQLiteOptions defines options for sqlite database.
type SQLiteOptions struct {
	Database string `json:"database" mapstructure:"database"`
}

// NewSQLiteOptions create a `zero` value instance.
func NewSQLiteOptions() *SQLiteOptions {
	return &SQLiteOptions{
		Database: "iam.db",
	}
}

// Validate verifies flags passed to SQLiteOptions.
func (o *SQLiteOptions) Validate() []error {
	errs := []error{}

	return errs
}

// AddFlags adds flags related to sqlite storage for a specific APIServer to the specified FlagSet.
func (o *SQLiteOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Database, "sqlite.database", o.Database, ""+
		"Sqlite database file path.")
}
