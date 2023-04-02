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
			Type:     "mysqldb",
			Database: "iam",
			//URL:      "mongodb://root:root@localhost:27017/",
			URL: "root:123456@tcp(127.0.0.1:3306)/iam_1?charset=utf8&parseTime=True&loc=Local&timeout=1000ms",
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
