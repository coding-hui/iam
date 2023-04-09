package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/wecoding/iam/docs/apidoc"
	"github.com/wecoding/iam/pkg/apiserver/config"
	"github.com/wecoding/iam/pkg/apiserver/domain/repository"
	"github.com/wecoding/iam/pkg/apiserver/domain/service"
	"github.com/wecoding/iam/pkg/apiserver/event"
	"github.com/wecoding/iam/pkg/apiserver/infrastructure/datastore/mysqldb"
	apisv1 "github.com/wecoding/iam/pkg/apiserver/interfaces/api"
	"github.com/wecoding/iam/pkg/middleware"
	"github.com/wecoding/iam/pkg/utils/container"
	"github.com/wecoding/iam/pkg/utils/env"
)

// APIServer interface for call api apiserver
type APIServer interface {
	Run(context.Context, chan error) error
}

// restServer rest apiserver
type restServer struct {
	webEngine        *gin.Engine
	beanContainer    *container.Container
	cfg              config.Config
	repositoryFactor repository.Factory
	dexProxy         *httputil.ReverseProxy
}

// New create api apiserver with config data
func New(cfg config.Config) (a APIServer) {
	if cfg.Mode == env.ModeProd.String() {
		gin.SetMode(gin.ReleaseMode)
	}
	s := &restServer{
		webEngine:     gin.New(),
		beanContainer: container.NewContainer(),
		cfg:           cfg,
	}
	return s
}

func (s *restServer) buildIoCContainer() (err error) {
	// infrastructure
	if err := s.beanContainer.ProvideWithName("RestServer", s); err != nil {
		return fmt.Errorf("fail to provides the RestServer bean to the container: %w", err)
	}

	// datastore repository
	var factory repository.Factory
	switch s.cfg.Datastore.Type {
	case "mysqldb":
		factory, err = mysqldb.GetMySQLFactory(context.Background(), s.cfg.Datastore)
		if err != nil {
			return fmt.Errorf("create mysqldb datastore instance failure %w", err)
		}
	default:
		return fmt.Errorf("not support datastore type %s", s.cfg.Datastore.Type)
	}
	s.repositoryFactor = factory
	if err := s.beanContainer.ProvideWithName("repository", s.repositoryFactor); err != nil {
		return fmt.Errorf("fail to provides the datastore bean to the container: %w", err)
	}
	repository.SetClient(factory)

	// domain
	if err := s.beanContainer.Provides(service.InitServiceBean(s.cfg)...); err != nil {
		return fmt.Errorf("fail to provides the service bean to the container: %w", err)
	}

	// interfaces
	if err := s.beanContainer.Provides(apisv1.InitAPIBean(s.cfg)...); err != nil {
		return fmt.Errorf("fail to provides the api bean to the container: %w", err)
	}

	// event
	if err := s.beanContainer.Provides(event.InitEvent(s.cfg)...); err != nil {
		return fmt.Errorf("fail to provides the event bean to the container: %w", err)
	}

	if err := s.beanContainer.Populate(); err != nil {
		return fmt.Errorf("fail to populate the bean container: %w", err)
	}
	return nil
}

func (s *restServer) Run(ctx context.Context, errChan chan error) error {
	// build the Ioc Container
	if err := s.buildIoCContainer(); err != nil {
		return err
	}

	// init database
	if err := service.InitData(ctx); err != nil {
		return fmt.Errorf("fail to init database %w", err)
	}

	s.RegisterAPIRoute()

	return s.startHTTP(ctx)
}

// RegisterAPIRoute register the API route
func (s *restServer) RegisterAPIRoute() {
	// Init middleware
	middleware.InitMiddleware(s.webEngine)

	// swagger router
	s.configSwagger()

	// Register all custom api
	for _, api := range apisv1.GetRegisteredAPI() {
		r := s.webEngine.Group(api.GetApiGroup().BaseUrl)
		if len(api.GetApiGroup().Filters) > 0 {
			for _, filter := range api.GetApiGroup().Filters {
				r.Use(filter)
			}
		}
		for _, v := range api.GetApiGroup().Apis {
			r.Handle(v.Method, v.Path, v.Handler)
		}
	}

	klog.Infof("init router successful")
}

func (s *restServer) configSwagger() {
	apidoc.SwaggerInfo.Title = "IAM API Doc"
	apidoc.SwaggerInfo.Description = "IAM ApiService API Doc."
	apidoc.SwaggerInfo.Version = "v1alpha"
	s.webEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.NewHandler()))
}

func (s *restServer) startHTTP(ctx context.Context) error {
	// Start HTTP apiserver
	klog.Infof("HTTP APIs are being served on: %s, ctx: %s", s.cfg.BindAddr, ctx)
	server := &http.Server{Addr: s.cfg.BindAddr, Handler: s.webEngine, ReadHeaderTimeout: 2 * time.Second}
	return server.ListenAndServe()
}
