package config

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/wecoding/iam/pkg/apiserver/infrastructure/datastore"
	"github.com/wecoding/iam/pkg/env"
)

// Config config for apiserver
type Config struct {
	// api apiserver bind address
	BindAddr string
	// monitor metric path
	MetricPath string
	// Datastore config
	Datastore datastore.Config
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
			Type:     "mongodb",
			Database: "wecoding",
			URL:      "",
		},
		Mode:         env.ModeDev.String(),
		PrintVersion: false,
	}
}

// Validate validate generic apiserver run options
func (s *Config) Validate() []error {
	var errs []error

	if s.Datastore.Type != "mongodb" && s.Datastore.Type != "mysql" {
		errs = append(errs, fmt.Errorf("not support datastore type %s", s.Datastore.Type))
	}

	return errs
}

// AddFlags adds flags to the specified FlagSet
func (s *Config) AddFlags(fs *pflag.FlagSet, c *Config) {
	fs.StringVar(&s.BindAddr, "bind-addr", c.BindAddr, "The bind address used to serve the http APIs.")
	fs.StringVar(&s.MetricPath, "metrics-path", c.MetricPath, "The path to expose the metrics.")
	fs.StringVar(&s.Datastore.Type, "datastore-type", c.Datastore.Type, "Metadata storage driver type, support kubeapi and mongodb")
	fs.StringVar(&s.Datastore.Database, "datastore-database", c.Datastore.Database, "Metadata storage database name, takes effect when the storage driver is mongodb.")
	fs.StringVar(&s.Datastore.URL, "datastore-url", c.Datastore.URL, "Metadata storage database url,takes effect when the storage driver is mongodb.")
}
