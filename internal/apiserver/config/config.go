// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"github.com/coding-hui/iam/cmd/iam-apiserver/app/options"
	pkgoptions "github.com/coding-hui/iam/pkg/options"
	genericapiserver "github.com/coding-hui/iam/pkg/server"
)

// Config config for iam-apiserver.
type Config struct {
	*options.Options
}

// MailOptions returns the mail configuration
func (c *Config) MailOptions() *pkgoptions.MailOptions {
	if c.Options != nil {
		return c.Options.MailOptions
	}
	return nil
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given IAM pump command line or configuration file option.
func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}

// BuildGenericConfig creates a generic server run config.
func (c *Config) BuildGenericConfig(cfg *Config) (genericConfig *genericapiserver.Config, lastErr error) {
	genericConfig = genericapiserver.NewConfig()
	if lastErr = cfg.GenericServerRunOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.FeatureOptions.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.SecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	if lastErr = cfg.InsecureServing.ApplyTo(genericConfig); lastErr != nil {
		return
	}

	return
}
