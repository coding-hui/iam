// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"encoding/json"

	"github.com/coding-hui/iam/internal/apiserver/infrastructure/cache"
	genericoptions "github.com/coding-hui/iam/internal/pkg/options"
	"github.com/coding-hui/iam/internal/pkg/server"
	"github.com/coding-hui/iam/pkg/log"

	cliflag "github.com/coding-hui/common/cli/flag"
	"github.com/coding-hui/common/util/idutil"
)

// Options runs an iam api server.
type Options struct {
	GenericServerRunOptions *genericoptions.ServerRunOptions       `json:"server"         mapstructure:"server"`
	GRPCOptions             *genericoptions.GRPCOptions            `json:"grpc"           mapstructure:"grpc"`
	InsecureServing         *genericoptions.InsecureServingOptions `json:"insecure"       mapstructure:"insecure"`
	SecureServing           *genericoptions.SecureServingOptions   `json:"secure"         mapstructure:"secure"`
	MySQLOptions            *genericoptions.MySQLOptions           `json:"mysql"          mapstructure:"mysql"`
	RedisOptions            *genericoptions.RedisOptions           `json:"redis"          mapstructure:"redis"`
	JwtOptions              *genericoptions.JwtOptions             `json:"jwt"            mapstructure:"jwt"`
	LogOptions              *log.Options                           `json:"log"            mapstructure:"log"`
	FeatureOptions          *genericoptions.FeatureOptions         `json:"feature"        mapstructure:"feature"`
	AuthenticationOptions   *genericoptions.AuthenticationOptions  `json:"authentication" mapstructure:"authentication"`
	CacheOptions            *cache.Options                         `json:"cache"          mapstructure:"cache"`
}

// ApplyTo applies the run options to the method receiver and returns self.
func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.JwtOptions.AddFlags(fss.FlagSet("jwt"))
	o.GRPCOptions.AddFlags(fss.FlagSet("grpc"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.FeatureOptions.AddFlags(fss.FlagSet("features"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure"))
	o.SecureServing.AddFlags(fss.FlagSet("secure"))
	o.LogOptions.AddFlags(fss.FlagSet("log"))
	o.AuthenticationOptions.AddFlags(fss.FlagSet("authentication"))

	return fss
}

func (o *Options) String() string {
	data, err := json.Marshal(o)
	if err != nil {
		log.Errorf("failed to marshal iam-apiserver options. err: %w", err)
		return ""
	}

	return string(data)
}

// Complete set default Options.
func (o *Options) Complete() error {
	if o.JwtOptions.Key == "" {
		o.JwtOptions.Key = idutil.NewSecretKey()
	}
	if o.AuthenticationOptions.JwtSecret == "" {
		o.AuthenticationOptions.JwtSecret = idutil.NewSecretKey()
	}

	return o.SecureServing.Complete()
}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	o := Options{
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		GRPCOptions:             genericoptions.NewGRPCOptions(),
		InsecureServing:         genericoptions.NewInsecureServingOptions(),
		SecureServing:           genericoptions.NewSecureServingOptions(),
		MySQLOptions:            genericoptions.NewMySQLOptions(),
		RedisOptions:            genericoptions.NewRedisOptions(),
		JwtOptions:              genericoptions.NewJwtOptions(),
		LogOptions:              log.NewOptions(),
		FeatureOptions:          genericoptions.NewFeatureOptions(),
		AuthenticationOptions:   genericoptions.NewAuthenticationOptions(),
		CacheOptions:            cache.NewCacheOptions(),
	}

	return &o
}
