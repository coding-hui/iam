// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authzserver

import (
	"context"
	"fmt"

	"github.com/coding-hui/iam/internal/authzserver/adapter"
	"github.com/coding-hui/iam/internal/authzserver/authorization"
	"github.com/coding-hui/iam/internal/authzserver/authorization/authorizer"
	"github.com/coding-hui/iam/internal/authzserver/config"
	apisv1 "github.com/coding-hui/iam/internal/authzserver/interfaces/api"
	"github.com/coding-hui/iam/internal/authzserver/store/apiserver"
	"github.com/coding-hui/iam/internal/pkg/middleware"
	genericapiserver "github.com/coding-hui/iam/internal/pkg/server"
	"github.com/coding-hui/iam/internal/pkg/utils/container"
	"github.com/coding-hui/iam/pkg/log"
	"github.com/coding-hui/iam/pkg/shutdown"
	"github.com/coding-hui/iam/pkg/shutdown/shutdownmanagers/posixsignal"
)

// AuthzServer interface for call iam-authzserver.
type AuthzServer interface {
	Run(context.Context, chan error) error
}

// authzServer rest iam-authzserver.
type authzServer struct {
	cfg           config.Config
	gs            *shutdown.GracefulShutdown
	webServer     *genericapiserver.GenericAPIServer
	beanContainer *container.Container
}

// New create iam-authzserver with config data.
func New(cfg *config.Config) (a AuthzServer, err error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	genericConfig, err := cfg.BuildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	server := &authzServer{
		webServer:     genericServer,
		cfg:           *cfg,
		gs:            gs,
		beanContainer: container.NewContainer(),
	}

	return server, nil
}

func (s *authzServer) Run(ctx context.Context, errChan chan error) error {
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		s.webServer.Close()
		return nil
	}))

	// build the Ioc Container
	if err := s.buildIoCContainer(); err != nil {
		return fmt.Errorf("failed to build IoCContainer %w", err)
	}

	// register apis
	s.registerAPIRoute()

	return s.startAuthzServer()
}

// buildIoCContainer build ioc container.
func (s *authzServer) buildIoCContainer() (err error) {
	// authz server
	if err = s.beanContainer.ProvideWithName("authzServer", s); err != nil {
		return fmt.Errorf("fail to provides the authzServer bean to the container: %w", err)
	}

	// infrastructure store
	storeFactory := apiserver.GetAPIServerFactoryOrDie(s.cfg.RPCServer, s.cfg.ClientCA)
	if err = s.beanContainer.ProvideWithName("store", storeFactory); err != nil {
		return fmt.Errorf("fail to provides the datastore bean to the container: %w", err)
	}

	// interfaces
	if err = s.beanContainer.Provides(apisv1.InitAPIBean(s.cfg)...); err != nil {
		return fmt.Errorf("fail to provides the api bean to the container: %w", err)
	}

	// adapter
	var a *adapter.Adapter
	a, err = adapter.NewAdapter(storeFactory)
	if err != nil {
		return fmt.Errorf("fail to create adapter bean: %w", err)
	}

	// authorizer service
	var authz authorization.Authorization
	authz, err = authorizer.NewAuthorizer(s.cfg, a)
	if err != nil {
		return fmt.Errorf("fail to create authorizer bean: %w", err)
	}
	if err = s.beanContainer.ProvideWithName("authorizer", authz); err != nil {
		return fmt.Errorf("fail to provides the authorizer bean to the container: %w", err)
	}

	if err = s.beanContainer.Populate(); err != nil {
		return fmt.Errorf("fail to populate the bean container: %w", err)
	}
	log.Infof("build IoC Container successful")

	return nil
}

// registerAPIRoute register the API route.
func (s *authzServer) registerAPIRoute() {
	// Init middleware
	middleware.InitMiddleware(s.webServer.Engine)

	// Register all custom api
	for _, api := range apisv1.GetRegisteredAPI() {
		api.RegisterApiGroup(s.webServer.Engine)
	}

	log.Infof("register API route successful")
}

// startAuthzServer start authz server.
func (s *authzServer) startAuthzServer() error {
	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	// web server
	return s.webServer.Run()
}
