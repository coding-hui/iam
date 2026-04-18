// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import (
	"github.com/coding-hui/common/errors"

	"github.com/coding-hui/iam/internal/config"
)

// Options embeds config.Config.
type Options struct {
	config.Config `mapstructure:",squash"`
}

// NewOptions creates a new Options instance with default values.
func NewOptions() *Options {
	return &Options{
		Config: config.Config{
			Server: config.ServerConfig{
				Port: 8080,
			},
			Database: config.DatabaseConfig{
				Driver: "sqlite",
			},
		},
	}
}

// Complete sets default values.
func (o *Options) Complete() error {
	if o.Server.Port == 0 {
		o.Server.Port = 8080
	}
	return nil
}

// Validate validates the required fields.
func (o *Options) Validate() []error {
	var errs []error

	if o.Database.DSN == "" {
		errs = append(errs, errors.New("database DSN is required"))
	}
	if o.Server.Port < 1 || o.Server.Port > 65535 {
		errs = append(errs, errors.New("server port must be between 1 and 65535"))
	}

	return errs
}
