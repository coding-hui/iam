// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"

	"github.com/coding-hui/iam/internal/apiserver/infrastructure/datastore"
	"github.com/coding-hui/iam/internal/pkg/options"
	"github.com/coding-hui/iam/internal/pkg/utils/env"
)

// Config config for apiserver
type Config struct {
	// api apiserver bind address
	BindAddr string
	// monitor metric path
	MetricPath string
	// Datastore config
	Datastore datastore.Config
	// Jwt config
	JwtOptions options.JwtOptions
	// RunMode
	Mode         string
	PrintVersion bool
}

// NewConfig  returns a Config struct with default values
func NewConfig() *Config {
	return &Config{
		BindAddr:   "0.0.0.0:8000",
		MetricPath: "/metrics",
		Datastore: datastore.Config{
			Type:     "mysqldb",
			Database: "iam_1",
			//URL:      "mongodb://root:root@localhost:27017/",
			URL: "root:123456@tcp(127.0.0.1:3306)/iam_1?charset=utf8mb4&parseTime=True&loc=Local&timeout=1000ms",
		},
		JwtOptions: options.JwtOptions{
			Realm:      "iam jwt",
			Key:        "iam-key",
			Timeout:    1 * time.Hour,
			MaxRefresh: 1 * time.Hour,
		},
		Mode:         env.ModeDev.String(),
		PrintVersion: false,
	}
}

// Validate validate generic apiserver run options
func (s *Config) Validate() []error {
	var errs []error

	if s.Datastore.Type != "mongodb" && s.Datastore.Type != "mysqldb" {
		errs = append(errs, fmt.Errorf("not support datastore type %s", s.Datastore.Type))
	}

	return errs
}

// AddFlags adds flags to the specified FlagSet
func (s *Config) AddFlags(fs *pflag.FlagSet, c *Config) {
	fs.StringVar(&s.BindAddr, "bind-addr", c.BindAddr, "The bind address used to serve the http APIs.")
	fs.StringVar(&s.MetricPath, "metrics-path", c.MetricPath, "The path to expose the metrics.")
	fs.StringVar(&s.Datastore.Type, "datastore-type", c.Datastore.Type, "Metadata storage driver type, support mysqldb and mongodb")
	fs.StringVar(&s.Datastore.Database, "datastore-database", c.Datastore.Database, "Metadata storage database name, takes effect when the storage driver is mongodb.")
	fs.StringVar(&s.Datastore.URL, "datastore-url", c.Datastore.URL, "Metadata storage database url,takes effect when the storage driver is mongodb.")
}
