// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"encoding/json"

	"github.com/coding-hui/iam/internal/apiserver/infrastructure/cache"
	"github.com/coding-hui/iam/pkg/log"
	"github.com/coding-hui/iam/pkg/options"
	"github.com/coding-hui/iam/pkg/server"

	cliflag "github.com/coding-hui/common/cli/flag"
	"github.com/coding-hui/common/util/idutil"
)

// Options runs an iam api server.
type Options struct {
	GenericServerRunOptions *options.ServerRunOptions       `json:"server"         mapstructure:"server"`
	GRPCOptions             *options.GRPCOptions            `json:"grpc"           mapstructure:"grpc"`
	InsecureServing         *options.InsecureServingOptions `json:"insecure"       mapstructure:"insecure"`
	SecureServing           *options.SecureServingOptions   `json:"secure"         mapstructure:"secure"`
	MySQLOptions            *options.MySQLOptions           `json:"mysql"          mapstructure:"mysql"`
	RedisOptions            *options.RedisOptions           `json:"redis"          mapstructure:"redis"`
	LogOptions              *log.Options                    `json:"log"            mapstructure:"log"`
	FeatureOptions          *options.FeatureOptions         `json:"feature"        mapstructure:"feature"`
	AuthenticationOptions   *options.AuthenticationOptions  `json:"authentication" mapstructure:"authentication"`
	CacheOptions            *cache.Options                  `json:"cache"          mapstructure:"cache"`
	MailOptions             *options.MailOptions            `json:"mail"           mapstructure:"mail"`
}

// ApplyTo applies the run options to the method receiver and returns self.
func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}

// Flags returns flags for a specific APIServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.GRPCOptions.AddFlags(fss.FlagSet("grpc"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.FeatureOptions.AddFlags(fss.FlagSet("features"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure"))
	o.SecureServing.AddFlags(fss.FlagSet("secure"))
	o.LogOptions.AddFlags(fss.FlagSet("log"))
	o.AuthenticationOptions.AddFlags(fss.FlagSet("authentication"))
	o.MailOptions.AddFlags(fss.FlagSet("mail"))

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
	if o.AuthenticationOptions.JwtSecret == "" {
		o.AuthenticationOptions.JwtSecret = idutil.NewSecretKey()
	}

	return o.SecureServing.Complete()
}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	o := Options{
		GenericServerRunOptions: options.NewServerRunOptions(),
		GRPCOptions:             options.NewGRPCOptions(),
		InsecureServing:         options.NewInsecureServingOptions(),
		SecureServing:           options.NewSecureServingOptions(),
		MySQLOptions:            options.NewMySQLOptions(),
		RedisOptions:            options.NewRedisOptions(),
		LogOptions:              log.NewOptions(),
		FeatureOptions:          options.NewFeatureOptions(),
		AuthenticationOptions:   options.NewAuthenticationOptions(),
		CacheOptions:            cache.NewCacheOptions(),
		MailOptions:             options.NewMailOptions(),
	}

	return &o
}
