// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// this file contains factories with no other dependencies

package util

import (
	restclient "github.com/coding-hui/wecoding-sdk-go/rest"
	"github.com/coding-hui/wecoding-sdk-go/services/iam"
	authzv1 "github.com/coding-hui/wecoding-sdk-go/services/iam/authz/v1"
	"github.com/coding-hui/wecoding-sdk-go/tools/clientcmd"
	"github.com/spf13/viper"

	"github.com/coding-hui/iam/pkg/cli/genericclioptions"
)

type factoryImpl struct {
	clientGetter genericclioptions.RESTClientGetter
}

func NewFactory(clientGetter genericclioptions.RESTClientGetter) Factory {
	if clientGetter == nil {
		panic("attempt to instantiate client_access_factory with nil clientGetter")
	}

	f := &factoryImpl{
		clientGetter: clientGetter,
	}

	return f
}

func (f *factoryImpl) ToRESTConfig() (*restclient.Config, error) {
	return f.clientGetter.ToRESTConfig()
}

func (f *factoryImpl) ToRawIAMConfigLoader() clientcmd.ClientConfig {
	return f.clientGetter.ToRawIAMConfigLoader()
}

func (f *factoryImpl) IAMClient() (*iam.IamClient, error) {
	clientConfig, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	return iam.NewForConfig(clientConfig)
}

func (f *factoryImpl) RESTClient() (*restclient.RESTClient, error) {
	clientConfig, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	_ = setIAMDefaults(clientConfig)
	return restclient.RESTClientFor(clientConfig)
}

func (f *factoryImpl) AuthzV1Client() (authzv1.AuthzV1Interface, error) {
	clientConfig, err := f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	// Override host to authzserver address
	authzServer := viper.GetString("authzserver.address")
	if authzServer == "" {
		authzServer = "http://127.0.0.1:9090"
	}
	// Update the host in the config
	clientConfig.Host = authzServer

	iamClient, err := iam.NewForConfig(clientConfig)
	if err != nil {
		return nil, err
	}

	return iamClient.AuthzV1(), nil
}
