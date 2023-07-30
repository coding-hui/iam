// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/coding-hui/iam/internal/apiserver/config"
	apisv1 "github.com/coding-hui/iam/internal/apiserver/interfaces/api"
	genericoptions "github.com/coding-hui/iam/internal/pkg/options"
	pb "github.com/coding-hui/iam/pkg/api/proto/apiserver/v1alpha1"
	"github.com/coding-hui/iam/pkg/log"
)

// gRPCConfig defines extra configuration for the iam-apiserver.
type gRPCConfig struct {
	Addr       string
	MaxMsgSize int
	ServerCert genericoptions.GeneratableKeyCert
}

type grpcAPIServer struct {
	*grpc.Server
	address string
}

func (s *grpcAPIServer) Run() {
	listen, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatalf("failed to listen: %s", err.Error())
	}

	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatalf("failed to start grpc server: %s", err.Error())
		}
	}()

	log.Infof("start grpc server at %s", s.address)
}

func (s *grpcAPIServer) Close() {
	s.GracefulStop()
	log.Infof("GRPC server on %s stopped", s.address)
}

func buildGRPCConfig(cfg *config.Config) (*gRPCConfig, error) {
	return &gRPCConfig{
		Addr:       fmt.Sprintf("%s:%d", cfg.GRPCOptions.BindAddress, cfg.GRPCOptions.BindPort),
		MaxMsgSize: cfg.GRPCOptions.MaxMsgSize,
		ServerCert: cfg.SecureServing.ServerCert,
	}, nil
}

type completedGRPCConfig struct {
	*gRPCConfig
}

// Complete fills in any fields not set that are required to have valid data and can be derived from other fields.
func (c *gRPCConfig) complete() *completedGRPCConfig {
	if c.Addr == "" {
		c.Addr = "127.0.0.1:8081"
	}

	return &completedGRPCConfig{c}
}

// New create a grpcAPIServer instance.
func (c *completedGRPCConfig) New() (*grpcAPIServer, error) {
	creds, err := credentials.NewServerTLSFromFile(c.ServerCert.CertKey.CertFile, c.ServerCert.CertKey.KeyFile)
	if err != nil {
		log.Fatalf("Failed to generate credentials %s", err.Error())
	}
	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(c.MaxMsgSize), grpc.Creds(creds)}

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterCacheServer(grpcServer, apisv1.NewCacheServer())

	reflection.Register(grpcServer)

	return &grpcAPIServer{Server: grpcServer, address: c.Addr}, nil
}
