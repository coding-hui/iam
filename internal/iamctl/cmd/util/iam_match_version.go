// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

import (
	"context"
	"fmt"
	"sync"

	"github.com/spf13/pflag"

	"github.com/coding-hui/common/runtime"
	"github.com/coding-hui/common/scheme"
	"github.com/coding-hui/common/version"
	"github.com/coding-hui/wecoding-sdk-go/rest"
	"github.com/coding-hui/wecoding-sdk-go/tools/clientcmd"

	"github.com/coding-hui/iam/pkg/cli/genericclioptions"
)

const (
	flagMatchBinaryVersion = "match-server-version"
)

// MatchVersionFlags is for setting the "match server version" function.
type MatchVersionFlags struct {
	Delegate genericclioptions.RESTClientGetter

	RequireMatchedServerVersion bool
	checkServerVersion          sync.Once
	matchesServerVersionErr     error
}

var _ genericclioptions.RESTClientGetter = &MatchVersionFlags{}

func (f *MatchVersionFlags) checkMatchingServerVersion() error {
	f.checkServerVersion.Do(func() {
		if !f.RequireMatchedServerVersion {
			return
		}

		clientConfig, err := f.Delegate.ToRESTConfig()
		if err != nil {
			f.matchesServerVersionErr = err
			return
		}

		_ = setIAMDefaults(clientConfig)
		restClient, err := rest.RESTClientFor(clientConfig)
		if err != nil {
			f.matchesServerVersionErr = err
			return
		}

		var sVer *version.Info
		if err := restClient.Get().AbsPath("/version").Do(context.TODO()).Into(&sVer); err != nil {
			f.matchesServerVersionErr = err
			return
		}

		clientVersion := version.Get()

		// GitVersion includes GitCommit and GitTreeState, but best to be safe?
		if clientVersion.GitVersion != sVer.GitVersion || clientVersion.GitCommit != sVer.GitCommit ||
			clientVersion.GitTreeState != sVer.GitTreeState {
			f.matchesServerVersionErr = fmt.Errorf(
				"server version (%#v) differs from client version (%#v)",
				sVer,
				version.Get(),
			)
		}
	})

	return f.matchesServerVersionErr
}

// ToRESTConfig implements RESTClientGetter.
// Returns a REST client configuration based on a provided path
// to a .iamconfig file, loading rules, and config flag overrides.
// Expects the AddFlags method to have been called.
func (f *MatchVersionFlags) ToRESTConfig() (*rest.Config, error) {
	if err := f.checkMatchingServerVersion(); err != nil {
		return nil, err
	}
	clientConfig, err := f.Delegate.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	// TODO we should not have to do this.  It smacks of something going wrong.
	if err = setIAMDefaults(clientConfig); err != nil {
		return nil, err
	}
	return clientConfig, nil
}

func (f *MatchVersionFlags) ToRawIAMConfigLoader() clientcmd.ClientConfig {
	return f.Delegate.ToRawIAMConfigLoader()
}

func (f *MatchVersionFlags) AddFlags(flags *pflag.FlagSet) {
	flags.BoolVar(
		&f.RequireMatchedServerVersion,
		flagMatchBinaryVersion,
		f.RequireMatchedServerVersion,
		"Require server version to match client version",
	)
}

func NewMatchVersionFlags(delegate genericclioptions.RESTClientGetter) *MatchVersionFlags {
	return &MatchVersionFlags{
		Delegate: delegate,
	}
}

// setIAMDefaults sets default values on the provided client config for accessing the
// IAM API or returns an error if any of the defaults are impossible or invalid.
// TODO this isn't what we want.  Each iamclient should be setting defaults as it sees fit.
func setIAMDefaults(config *rest.Config) error {
	// TODO remove this hack.  This is allowing the GetOptions to be serialized.
	config.GroupVersion = &scheme.GroupVersion{Group: "iam.api", Version: "v1"}

	if config.APIPath == "" {
		config.APIPath = "/api"
	}
	if config.Negotiator == nil {
		// This codec factory ensures the resources are not converted. Therefore, resources
		// will not be round-tripped through internal versions. Defaulting does not happen
		// on the client.
		config.Negotiator = runtime.NewSimpleClientNegotiator()
	}
	return rest.SetIAMDefaults(config)
}
