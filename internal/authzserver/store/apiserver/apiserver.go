// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package apiserver

import (
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"k8s.io/klog/v2"

	"github.com/coding-hui/iam/internal/authzserver/store"
	pb "github.com/coding-hui/iam/pkg/api/proto/apiserver/v1alpha1"
)

type datastore struct {
	cli pb.CacheClient
}

func (ds *datastore) Policies() store.PolicyStore {
	return newPolicies(ds)
}

var (
	apiServerFactory store.Factory
	once             sync.Once
)

// GetAPIServerFactoryOrDie return cache instance and panics on any error.
func GetAPIServerFactoryOrDie(address string, clientCA string) store.Factory {
	once.Do(func() {
		var (
			err   error
			conn  *grpc.ClientConn
			creds credentials.TransportCredentials
		)

		creds, err = credentials.NewClientTLSFromFile(clientCA, "")
		if err != nil {
			klog.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
		}

		conn, err = grpc.Dial(address, grpc.WithBlock(), grpc.WithTransportCredentials(creds))
		if err != nil {
			klog.Fatalf("Connect to grpc server failed, error: %s", err.Error())
		}

		apiServerFactory = &datastore{pb.NewCacheClient(conn)}
		klog.Infof("Connected to grpc server, address: %s", address)
	})

	if apiServerFactory == nil {
		klog.Fatalf("failed to get apiserver store fatory")
	}

	return apiServerFactory
}
