// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import (
	"context"
	"fmt"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"k8s.io/klog/v2"

	"github.com/coding-hui/iam/pkg/api/apiserver/v1alpha1"
	"github.com/coding-hui/iam/pkg/shutdown"
	"github.com/coding-hui/iam/pkg/shutdown/shutdownmanagers/posixsignal"

	_ "github.com/coding-hui/iam/api/swagger"
	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/apiserver/domain/repository"
	"github.com/coding-hui/iam/internal/apiserver/domain/service"
	"github.com/coding-hui/iam/internal/apiserver/event"
	"github.com/coding-hui/iam/internal/apiserver/infrastructure/datastore/mysqldb"
	apisv1 "github.com/coding-hui/iam/internal/apiserver/interfaces/api"
	"github.com/coding-hui/iam/internal/pkg/middleware"
	genericapiserver "github.com/coding-hui/iam/internal/pkg/server"
	"github.com/coding-hui/iam/internal/pkg/utils/container"
)

//	@title			IAM API
//	@version		v1alpha
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
	beanContainer    *container.Container
	repositoryFactor repository.Factory
}

// New create iam-apiserver with config data.
func New(cfg *config.Config) (a APIServer, err error) {
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

	server := &apiServer{
		webServer:     genericServer,
		beanContainer: container.NewContainer(),
		cfg:           *cfg,
		gs:            gs,
	}

	return server, nil
}

func (s *apiServer) Run(ctx context.Context, errChan chan error) error {
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		s.webServer.Close()
		if s.repositoryFactor != nil {
			_ = s.repositoryFactor.Close()
		}

		return nil
	}))

	// build the Ioc Container
	if err := s.buildIoCContainer(); err != nil {
		return fmt.Errorf("failed to build IoCContainer %w", err)
	}

	// register apis
	s.registerAPIRoute()

	// init database
	if err := service.InitData(s.withRoutesContext(ctx)); err != nil {
		return fmt.Errorf("failed to init database %w", err)
	}

	go event.StartEventWorker(ctx, errChan)

	return s.startAPIServer()
}

// buildIoCContainer build ioc container.
func (s *apiServer) buildIoCContainer() (err error) {
	// infrastructure
	if err = s.beanContainer.ProvideWithName("apiServer", s); err != nil {
		return fmt.Errorf("fail to provides the apiServer bean to the container: %w", err)
	}

	// datastore repository
	var factory repository.Factory
	if s.cfg.MySQLOptions != nil {
		factory, err = mysqldb.GetMySQLFactory(context.Background(), s.cfg)
		if err != nil {
			return fmt.Errorf("create mysqldb datastore instance failure %w", err)
		}
	}

	s.repositoryFactor = factory
	if err = s.beanContainer.ProvideWithName("repository", s.repositoryFactor); err != nil {
		return fmt.Errorf("fail to provides the datastore bean to the container: %w", err)
	}
	repository.SetClient(factory)

	// domain
	if err = s.beanContainer.Provides(service.InitServiceBean(s.cfg)...); err != nil {
		return fmt.Errorf("fail to provides the service bean to the container: %w", err)
	}

	// interfaces
	if err = s.beanContainer.Provides(apisv1.InitAPIBean(s.cfg)...); err != nil {
		return fmt.Errorf("fail to provides the api bean to the container: %w", err)
	}

	// event
	if err = s.beanContainer.Provides(event.InitEvent(s.cfg)...); err != nil {
		return fmt.Errorf("fail to provides the event bean to the container: %w", err)
	}

	if err = s.beanContainer.Populate(); err != nil {
		return fmt.Errorf("fail to populate the bean container: %w", err)
	}
	klog.Infof("build IoC Container successful")

	return nil
}

// registerAPIRoute register the API route.
func (s *apiServer) registerAPIRoute() {
	// Init middleware
	middleware.InitMiddleware(s.webServer.Engine)

	// swagger router
	s.configSwagger()

	// Register all custom api
	for _, api := range apisv1.GetRegisteredAPI() {
		api.RegisterApiGroup(s.webServer.Engine)
	}

	klog.Infof("register API route successful")
}

func (s *apiServer) configSwagger() {
	s.webServer.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.NewHandler()))
}

func (s *apiServer) withRoutesContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, &v1alpha1.CtxKeyRoutes, s.webServer.Routes())
	ctx = context.WithValue(ctx, &v1alpha1.CtxKeyApiPrefix, apisv1.GetAPIPrefix())
	return ctx
}

// startAPIServer start api server.
func (s *apiServer) startAPIServer() error {
	// start shutdown managers
	if err := s.gs.Start(); err != nil {
		klog.Fatalf("start shutdown manager failed: %s", err.Error())
	}

	// web server
	return s.webServer.Run()
}
