// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package authzserver

import (
	"context"

	"k8s.io/klog/v2"

	"github.com/coding-hui/iam/internal/authzserver/config"
	"github.com/coding-hui/iam/internal/pkg/middleware"
	genericapiserver "github.com/coding-hui/iam/internal/pkg/server"
	"github.com/coding-hui/iam/pkg/shutdown"
	"github.com/coding-hui/iam/pkg/shutdown/shutdownmanagers/posixsignal"
)

// AuthzServer interface for call iam-authzserver.
type AuthzServer interface {
	Run(context.Context, chan error) error
}

// authzServer rest iam-authzserver.
type authzServer struct {
	cfg       config.Config
	gs        *shutdown.GracefulShutdown
	webServer *genericapiserver.GenericAPIServer
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
		webServer: genericServer,
		cfg:       *cfg,
		gs:        gs,
	}

	return server, nil
}

func (s *authzServer) Run(ctx context.Context, errChan chan error) error {
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		s.webServer.Close()
		return nil
	}))

	// register apis
	s.registerAPIRoute()

	return s.startAuthzServer()
}

// registerAPIRoute register the API route.
func (s *authzServer) registerAPIRoute() {
	// Init middleware
	middleware.InitMiddleware(s.webServer.Engine)

	klog.Infof("register API route successful")
}

// startAuthzServer start authz server.
func (s *authzServer) startAuthzServer() error {
	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		klog.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	// web server
	return s.webServer.Run()
}
