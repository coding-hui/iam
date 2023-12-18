// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"encoding/json"

	"github.com/coding-hui/iam/pkg/log"
	"github.com/coding-hui/iam/pkg/options"
	"github.com/coding-hui/iam/pkg/server"

	cliflag "github.com/coding-hui/common/cli/flag"
)

// Options runs an iam api server.
type Options struct {
	RPCServer               string                          `json:"rpcserver"      mapstructure:"rpcserver"`
	ClientCA                string                          `json:"client-ca-file" mapstructure:"client-ca-file"`
	GenericServerRunOptions *options.ServerRunOptions       `json:"server"         mapstructure:"server"`
	InsecureServing         *options.InsecureServingOptions `json:"insecure"       mapstructure:"insecure"`
	SecureServing           *options.SecureServingOptions   `json:"secure"         mapstructure:"secure"`
	RedisOptions            *options.RedisOptions           `json:"redis"          mapstructure:"redis"`
	LogOptions              *log.Options                    `json:"log"            mapstructure:"log"`
	FeatureOptions          *options.FeatureOptions         `json:"feature"        mapstructure:"feature"`
}

// ApplyTo applies the run options to the method receiver and returns self.
func (o *Options) ApplyTo(c *server.Config) error {
	return nil
}

// Flags returns flags for a specific AuthzServer by section name.
func (o *Options) Flags() (fss cliflag.NamedFlagSets) {
	o.GenericServerRunOptions.AddFlags(fss.FlagSet("generic"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.FeatureOptions.AddFlags(fss.FlagSet("features"))
	o.InsecureServing.AddFlags(fss.FlagSet("insecure serving"))
	o.SecureServing.AddFlags(fss.FlagSet("secure serving"))
	o.LogOptions.AddFlags(fss.FlagSet("log"))

	// Note: the weird ""+ in below lines seems to be the only way to get gofmt to
	// arrange these text blocks sensibly. Grrr.
	fs := fss.FlagSet("misc")
	fs.StringVar(&o.RPCServer, "rpcserver", o.RPCServer, "The address of iam rpc server. "+
		"The rpc server can provide all the secrets and policies to use.")
	fs.StringVar(&o.ClientCA, "client-ca-file", o.ClientCA, ""+
		"If set, any request presenting a client certificate signed by one of "+
		"the authorities in the client-ca-file is authenticated with an identity "+
		"corresponding to the CommonName of the client certificate.")
	return fss
}

func (o *Options) String() string {
	data, err := json.Marshal(o)
	if err != nil {
		log.Errorf("failed to marshal iam-authzserver options. err: %w", err)
		return ""
	}

	return string(data)
}

// Complete set default Options.
func (o *Options) Complete() error {
	return o.SecureServing.Complete()
}

// NewOptions creates a new Options object with default parameters.
func NewOptions() *Options {
	o := Options{
		RPCServer:               "127.0.0.1:8081",
		ClientCA:                "",
		GenericServerRunOptions: options.NewServerRunOptions(),
		InsecureServing:         options.NewInsecureServingOptions(),
		SecureServing:           options.NewSecureServingOptions(),
		RedisOptions:            options.NewRedisOptions(),
		LogOptions:              log.NewOptions(),
		FeatureOptions:          options.NewFeatureOptions(),
	}

	return &o
}
