// Copyright (c) 2023 coding-hui. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.5
// source: pkg/api/proto/apiserver/v1alpha1/cache.proto

package v1alpha1

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Cache_DetailPolicy_FullMethodName    = "/proto.apiserver.v1alpha1.Cache/DetailPolicy"
	Cache_ListPolicies_FullMethodName    = "/proto.apiserver.v1alpha1.Cache/ListPolicies"
	Cache_ListPolicyRules_FullMethodName = "/proto.apiserver.v1alpha1.Cache/ListPolicyRules"
)

// CacheClient is the client API for Cache service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CacheClient interface {
	DetailPolicy(ctx context.Context, in *GetPolicyRequest, opts ...grpc.CallOption) (*PolicyInfo, error)
	ListPolicies(ctx context.Context, in *ListPoliciesRequest, opts ...grpc.CallOption) (*ListPoliciesResponse, error)
	ListPolicyRules(ctx context.Context, in *ListPolicyRulesRequest, opts ...grpc.CallOption) (*ListPolicyRulesResponse, error)
}

type cacheClient struct {
	cc grpc.ClientConnInterface
}

func NewCacheClient(cc grpc.ClientConnInterface) CacheClient {
	return &cacheClient{cc}
}

func (c *cacheClient) DetailPolicy(ctx context.Context, in *GetPolicyRequest, opts ...grpc.CallOption) (*PolicyInfo, error) {
	out := new(PolicyInfo)
	err := c.cc.Invoke(ctx, Cache_DetailPolicy_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cacheClient) ListPolicies(ctx context.Context, in *ListPoliciesRequest, opts ...grpc.CallOption) (*ListPoliciesResponse, error) {
	out := new(ListPoliciesResponse)
	err := c.cc.Invoke(ctx, Cache_ListPolicies_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cacheClient) ListPolicyRules(ctx context.Context, in *ListPolicyRulesRequest, opts ...grpc.CallOption) (*ListPolicyRulesResponse, error) {
	out := new(ListPolicyRulesResponse)
	err := c.cc.Invoke(ctx, Cache_ListPolicyRules_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CacheServer is the server API for Cache service.
// All implementations must embed UnimplementedCacheServer
// for forward compatibility
type CacheServer interface {
	DetailPolicy(context.Context, *GetPolicyRequest) (*PolicyInfo, error)
	ListPolicies(context.Context, *ListPoliciesRequest) (*ListPoliciesResponse, error)
	ListPolicyRules(context.Context, *ListPolicyRulesRequest) (*ListPolicyRulesResponse, error)
	mustEmbedUnimplementedCacheServer()
}

// UnimplementedCacheServer must be embedded to have forward compatible implementations.
type UnimplementedCacheServer struct {
}

func (UnimplementedCacheServer) DetailPolicy(context.Context, *GetPolicyRequest) (*PolicyInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetailPolicy not implemented")
}
func (UnimplementedCacheServer) ListPolicies(context.Context, *ListPoliciesRequest) (*ListPoliciesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPolicies not implemented")
}
func (UnimplementedCacheServer) ListPolicyRules(context.Context, *ListPolicyRulesRequest) (*ListPolicyRulesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPolicyRules not implemented")
}
func (UnimplementedCacheServer) mustEmbedUnimplementedCacheServer() {}

// UnsafeCacheServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CacheServer will
// result in compilation errors.
type UnsafeCacheServer interface {
	mustEmbedUnimplementedCacheServer()
}

func RegisterCacheServer(s grpc.ServiceRegistrar, srv CacheServer) {
	s.RegisterService(&Cache_ServiceDesc, srv)
}

func _Cache_DetailPolicy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPolicyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CacheServer).DetailPolicy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cache_DetailPolicy_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CacheServer).DetailPolicy(ctx, req.(*GetPolicyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cache_ListPolicies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPoliciesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CacheServer).ListPolicies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cache_ListPolicies_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CacheServer).ListPolicies(ctx, req.(*ListPoliciesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Cache_ListPolicyRules_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPolicyRulesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CacheServer).ListPolicyRules(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Cache_ListPolicyRules_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CacheServer).ListPolicyRules(ctx, req.(*ListPolicyRulesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Cache_ServiceDesc is the grpc.ServiceDesc for Cache service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Cache_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.apiserver.v1alpha1.Cache",
	HandlerType: (*CacheServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DetailPolicy",
			Handler:    _Cache_DetailPolicy_Handler,
		},
		{
			MethodName: "ListPolicies",
			Handler:    _Cache_ListPolicies_Handler,
		},
		{
			MethodName: "ListPolicyRules",
			Handler:    _Cache_ListPolicyRules_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/api/proto/apiserver/v1alpha1/cache.proto",
}
