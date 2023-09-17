// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.24.3
// source: proto/webhandler.proto

package proto

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

// WebHandlerClient is the client API for WebHandler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WebHandlerClient interface {
	GetEndpoints(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*WebEndpointList, error)
	HandleRequest(ctx context.Context, in *WebRequest, opts ...grpc.CallOption) (*WebResponse, error)
}

type webHandlerClient struct {
	cc grpc.ClientConnInterface
}

func NewWebHandlerClient(cc grpc.ClientConnInterface) WebHandlerClient {
	return &webHandlerClient{cc}
}

func (c *webHandlerClient) GetEndpoints(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*WebEndpointList, error) {
	out := new(WebEndpointList)
	err := c.cc.Invoke(ctx, "/interfaces.WebHandler/GetEndpoints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webHandlerClient) HandleRequest(ctx context.Context, in *WebRequest, opts ...grpc.CallOption) (*WebResponse, error) {
	out := new(WebResponse)
	err := c.cc.Invoke(ctx, "/interfaces.WebHandler/HandleRequest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WebHandlerServer is the server API for WebHandler service.
// All implementations must embed UnimplementedWebHandlerServer
// for forward compatibility
type WebHandlerServer interface {
	GetEndpoints(context.Context, *Empty) (*WebEndpointList, error)
	HandleRequest(context.Context, *WebRequest) (*WebResponse, error)
	mustEmbedUnimplementedWebHandlerServer()
}

// UnimplementedWebHandlerServer must be embedded to have forward compatible implementations.
type UnimplementedWebHandlerServer struct {
}

func (UnimplementedWebHandlerServer) GetEndpoints(context.Context, *Empty) (*WebEndpointList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEndpoints not implemented")
}
func (UnimplementedWebHandlerServer) HandleRequest(context.Context, *WebRequest) (*WebResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleRequest not implemented")
}
func (UnimplementedWebHandlerServer) mustEmbedUnimplementedWebHandlerServer() {}

// UnsafeWebHandlerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WebHandlerServer will
// result in compilation errors.
type UnsafeWebHandlerServer interface {
	mustEmbedUnimplementedWebHandlerServer()
}

func RegisterWebHandlerServer(s grpc.ServiceRegistrar, srv WebHandlerServer) {
	s.RegisterService(&WebHandler_ServiceDesc, srv)
}

func _WebHandler_GetEndpoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebHandlerServer).GetEndpoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/interfaces.WebHandler/GetEndpoints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebHandlerServer).GetEndpoints(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _WebHandler_HandleRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WebRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebHandlerServer).HandleRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/interfaces.WebHandler/HandleRequest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebHandlerServer).HandleRequest(ctx, req.(*WebRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// WebHandler_ServiceDesc is the grpc.ServiceDesc for WebHandler service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WebHandler_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "interfaces.WebHandler",
	HandlerType: (*WebHandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetEndpoints",
			Handler:    _WebHandler_GetEndpoints_Handler,
		},
		{
			MethodName: "HandleRequest",
			Handler:    _WebHandler_HandleRequest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/webhandler.proto",
}
