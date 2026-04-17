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

	"github.com/coding-hui/iam/internal/api"
	"github.com/coding-hui/iam/internal/apiserver/config"
	"github.com/coding-hui/iam/internal/driver"
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
	cfg           config.Config
	gs            *shutdown.GracefulShutdown
	webServer     *genericapiserver.GenericAPIServer
	gRPCAPIServer *grpcAPIServer
	registry      *driver.RegistryDefault
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
		cfg:           *cfg,
		gs:            gs,
		webServer:     genericServer,
		gRPCAPIServer: gRPCAPIServer,
	}

	return server, nil
}

// Run runs the apiserver.
func (s *apiServer) Run(ctx context.Context, errChan chan error) (lastErr error) {
	s.gs.AddShutdownCallback(shutdown.ShutdownFunc(func(string) error {
		s.webServer.Close()
		if s.cfg.GRPCOptions.BindPort > 0 {
			s.gRPCAPIServer.Close()
		}
		return nil
	}))

	// Create driver config from options
	driverConfig := s.createDriverConfig()

	// Create registry and initialize
	s.registry = driver.NewRegistry(driverConfig)
	if lastErr = s.registry.Init(ctx); lastErr != nil {
		return fmt.Errorf("failed to initialize registry: %w", lastErr)
	}

	// Create router with registry
	router := api.NewRouter(s.registry)

	// Replace the embedded gin.Engine with our new router
	s.webServer.Engine = router

	// Run database migrations
	if lastErr = s.registry.MigrateUp(ctx); lastErr != nil {
		return fmt.Errorf("failed to migrate database: %w", lastErr)
	}

	// Start servers
	return s.startAPIServer()
}

// createDriverConfig creates a driver.Config from the apiserver config.
func (s *apiServer) createDriverConfig() *driver.Config {
	opts := s.cfg.Options

	serverMode := "debug"
	switch opts.GenericServerRunOptions.Mode {
	case "release":
		serverMode = "release"
	case "test":
		serverMode = "test"
	}

	var dbDriver string
	var dbDSN string
	if opts.SQLiteOptions != nil && opts.SQLiteOptions.Database != "" {
		dbDriver = "sqlite"
		dbDSN = opts.SQLiteOptions.Database
	} else if opts.MySQLOptions != nil && opts.MySQLOptions.Database != "" {
		dbDriver = "mysql"
		// Construct DSN from individual fields
		dbDSN = opts.MySQLOptions.Username + ":" + opts.MySQLOptions.Password +
			"@tcp(" + opts.MySQLOptions.Host + ")/" + opts.MySQLOptions.Database +
			"?charset=utf8&parseTime=True&loc=Local"
	}

	redisAddr := ""
	redisPassword := ""
	redisDB := 0
	if opts.RedisOptions != nil && len(opts.RedisOptions.Addrs) > 0 {
		redisAddr = opts.RedisOptions.Addrs[0]
		redisPassword = opts.RedisOptions.Password
		redisDB = opts.RedisOptions.Database
	}

	return &driver.Config{
		Server: driver.ServerConfig{
			Host:        opts.InsecureServing.BindAddress,
			Port:        opts.InsecureServing.BindPort,
			Mode:        serverMode,
			Healthz:     true,
			Middlewares: opts.GenericServerRunOptions.Middlewares,
		},
		Database: driver.DatabaseConfig{
			Driver: dbDriver,
			DSN:    dbDSN,
		},
		Redis: driver.RedisConfig{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       redisDB,
		},
	}
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

	// Configure swagger before running
	s.webServer.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.NewHandler()))

	// web server
	return s.webServer.Run()
}
