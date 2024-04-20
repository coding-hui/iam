// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import (
	"context"
	"fmt"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/coding-hui/iam/api/swagger"

	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/apiserver/event"
	"github.com/coding-hui/iam/internal/apiserver/infrastructure/cache"
	"github.com/coding-hui/iam/internal/apiserver/infrastructure/datastore/mysqldb"
	apisv1 "github.com/coding-hui/iam/internal/apiserver/interfaces/api"
	"github.com/coding-hui/iam/internal/pkg/token"
	"github.com/coding-hui/iam/pkg/container"
	"github.com/coding-hui/iam/pkg/log"
	genericapiserver "github.com/coding-hui/iam/pkg/server"
	"github.com/coding-hui/iam/pkg/shutdown"
	"github.com/coding-hui/iam/pkg/shutdown/shutdownmanagers/posixsignal"
)

//	@title			IAM API
//	@version		v1
//	@description	IAM ApiService API Doc.
//	@termsOfService	http://iam.wecoding.top/

//	@contact.name	coding-hui
//	@contact.url	https://github.com/coding-hui/iam
//	@contact.email	wecoding@yeah.net

//	@license.name	MIT
//	@license.url	https://en.wikipedia.org/wiki/MIT_License

//	@host	localhost:8000

//	@securityDefinitions.apikey	BearerTokenAuth
//	@in							header
//	@name						Authorization
//	@description				Set Bearer token in the request header

// APIServer interface for call iam-apiserver.
type APIServer interface {
	Run(context.Context, chan error) error
}

// apiServer rest iam-apiserver.
type apiServer struct {
	cfg              config.Config
	gs               *shutdown.GracefulShutdown
	webServer        *genericapiserver.GenericAPIServer
	gRPCAPIServer    *grpcAPIServer
	beanContainer    *container.Container
	repositoryFactor repository.Factory
	eventBus         *event.Bus

	// entity that issues tokens
	issuer token.Issuer
}

// New create iam-apiserver with config data.
func New(cfg *config.Config) (a APIServer, err error) {
	gs := shutdown.New()
	gs.AddShutdownManager(posixsignal.NewPosixSignalManager())

	genericConfig, err := cfg.BuildGenericConfig(cfg)
	if err != nil {
		return nil, err
	}

	gRPCConfig, err := buildGRPCConfig(cfg)
	if err != nil {
		return nil, err
	}

	genericServer, err := genericConfig.Complete().New()
	if err != nil {
		return nil, err
	}

	gRPCAPIServer, err := gRPCConfig.complete().New()
	if err != nil {
		return nil, err
	}

	server := &apiServer{
		webServer:     genericServer,
		gRPCAPIServer: gRPCAPIServer,
		beanContainer: container.NewContainer(),
		cfg:           *cfg,
		gs:            gs,
	}

	return server, nil
}

func (s *apiServer) Run(ctx context.Context, errChan chan error) (lastErr error) {
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		s.webServer.Close()
		if s.cfg.GRPCOptions.BindPort > 0 {
			s.gRPCAPIServer.Close()
		}
		if s.repositoryFactor != nil {
			_ = s.repositoryFactor.Close()
		}
		if s.eventBus != nil {
			_ = (*s.eventBus).CloseWait()
		}

		return nil
	}))

	// create token issuer
	s.issuer, lastErr = token.NewIssuer(s.cfg.AuthenticationOptions)
	if lastErr != nil {
		return fmt.Errorf("unable to create issuer: %w", lastErr)
	}

	// build the Ioc Container
	if lastErr = s.buildIoCContainer(ctx); lastErr != nil {
		return fmt.Errorf("failed to build IoCContainer %w", lastErr)
	}

	// register apis
	s.registerAPIRoute()

	// init database
	if lastErr = service.InitData(ctx); lastErr != nil {
		return fmt.Errorf("failed to init database %w", lastErr)
	}

	return s.startAPIServer()
}

// buildIoCContainer build ioc container.
func (s *apiServer) buildIoCContainer(ctx context.Context) (err error) {
	// infrastructure
	if err = s.beanContainer.ProvideWithName("apiServer", s); err != nil {
		return fmt.Errorf("fail to provides the apiServer bean to the container: %w", err)
	}

	// datastore repository
	var factory repository.Factory
	if s.cfg.MySQLOptions != nil {
		factory, err = mysqldb.New(context.Background(), s.cfg)
		if err != nil {
			return fmt.Errorf("create mysqldb datastore instance failure %w", err)
		}
	}

	s.repositoryFactor = factory
	if err = s.beanContainer.ProvideWithName("repository", s.repositoryFactor); err != nil {
		return fmt.Errorf("fail to provides the datastore bean to the container: %w", err)
	}
	repository.SetClient(factory)

	// cache
	var cacheClient cache.Interface
	if cacheClient, err = cache.New(s.cfg.CacheOptions, ctx.Done()); err != nil {
		return fmt.Errorf("failed to create cache, error: %w", err)
	}
	if err = s.beanContainer.ProvideWithName("cache", cacheClient); err != nil {
		return fmt.Errorf("fail to provides the cache bean to the container: %w", err)
	}

	// domain
	if err = s.beanContainer.Provides(service.InitServiceBean(s.cfg, s.issuer)...); err != nil {
		return fmt.Errorf("fail to provides the service bean to the container: %w", err)
	}

	// interfaces
	if err = s.beanContainer.Provides(apisv1.InitAPIBean(s.cfg)...); err != nil {
		return fmt.Errorf("fail to provides the api bean to the container: %w", err)
	}

	// event
	eventBus, listeners := event.InitEvent(s.cfg)
	if err = s.beanContainer.ProvideWithName("eventBus", eventBus); err != nil {
		return fmt.Errorf("fail to provides the event bus bean to the container: %w", err)
	}
	if err = s.beanContainer.Provides(listeners...); err != nil {
		return fmt.Errorf("fail to provides the event listener bean to the container: %w", err)
	}
	s.eventBus = &eventBus

	if err = s.beanContainer.Populate(); err != nil {
		return fmt.Errorf("fail to populate the bean container: %w", err)
	}
	log.Infof("build IoC Container successful")

	return nil
}

// registerAPIRoute register the API route.
func (s *apiServer) registerAPIRoute() {
	// swagger router
	s.configSwagger()

	// Register all custom api
	for _, api := range apisv1.GetRegisteredAPI() {
		api.RegisterApiGroup(s.webServer.Engine)
	}

	log.Infof("register API route successful")
}

func (s *apiServer) configSwagger() {
	s.webServer.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.NewHandler()))
}

// startAPIServer start api server.
func (s *apiServer) startAPIServer() error {
	// start gRPC server
	if s.cfg.GRPCOptions.BindPort > 0 {
		go s.gRPCAPIServer.Run()
	}

	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		log.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	// web server
	return s.webServer.Run()
}
