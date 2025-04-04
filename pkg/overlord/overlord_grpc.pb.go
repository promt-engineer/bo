// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package overlord

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

// OverlordClient is the client API for Overlord service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OverlordClient interface {
	GetIntegratorConfig(ctx context.Context, in *GetIntegratorConfigIn, opts ...grpc.CallOption) (*GetIntegratorConfigOut, error)
	SaveParams(ctx context.Context, in *SaveParamsIn, opts ...grpc.CallOption) (*SaveParamsOut, error)
	HealthCheck(ctx context.Context, opts ...grpc.CallOption) (Overlord_HealthCheckClient, error)
}

type overlordClient struct {
	cc grpc.ClientConnInterface
}

func NewOverlordClient(cc grpc.ClientConnInterface) OverlordClient {
	return &overlordClient{cc}
}

func (c *overlordClient) GetIntegratorConfig(ctx context.Context, in *GetIntegratorConfigIn, opts ...grpc.CallOption) (*GetIntegratorConfigOut, error) {
	out := new(GetIntegratorConfigOut)
	err := c.cc.Invoke(ctx, "/overlord.Overlord/GetIntegratorConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *overlordClient) SaveParams(ctx context.Context, in *SaveParamsIn, opts ...grpc.CallOption) (*SaveParamsOut, error) {
	out := new(SaveParamsOut)
	err := c.cc.Invoke(ctx, "/overlord.Overlord/SaveParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *overlordClient) HealthCheck(ctx context.Context, opts ...grpc.CallOption) (Overlord_HealthCheckClient, error) {
	stream, err := c.cc.NewStream(ctx, &Overlord_ServiceDesc.Streams[0], "/overlord.Overlord/HealthCheck", opts...)
	if err != nil {
		return nil, err
	}
	x := &overlordHealthCheckClient{stream}
	return x, nil
}

type Overlord_HealthCheckClient interface {
	Send(*Status) error
	Recv() (*Status, error)
	grpc.ClientStream
}

type overlordHealthCheckClient struct {
	grpc.ClientStream
}

func (x *overlordHealthCheckClient) Send(m *Status) error {
	return x.ClientStream.SendMsg(m)
}

func (x *overlordHealthCheckClient) Recv() (*Status, error) {
	m := new(Status)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// OverlordServer is the server API for Overlord service.
// All implementations should embed UnimplementedOverlordServer
// for forward compatibility
type OverlordServer interface {
	GetIntegratorConfig(context.Context, *GetIntegratorConfigIn) (*GetIntegratorConfigOut, error)
	SaveParams(context.Context, *SaveParamsIn) (*SaveParamsOut, error)
	HealthCheck(Overlord_HealthCheckServer) error
}

// UnimplementedOverlordServer should be embedded to have forward compatible implementations.
type UnimplementedOverlordServer struct {
}

func (UnimplementedOverlordServer) GetIntegratorConfig(context.Context, *GetIntegratorConfigIn) (*GetIntegratorConfigOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIntegratorConfig not implemented")
}
func (UnimplementedOverlordServer) SaveParams(context.Context, *SaveParamsIn) (*SaveParamsOut, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveParams not implemented")
}
func (UnimplementedOverlordServer) HealthCheck(Overlord_HealthCheckServer) error {
	return status.Errorf(codes.Unimplemented, "method HealthCheck not implemented")
}

// UnsafeOverlordServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OverlordServer will
// result in compilation errors.
type UnsafeOverlordServer interface {
	mustEmbedUnimplementedOverlordServer()
}

func RegisterOverlordServer(s grpc.ServiceRegistrar, srv OverlordServer) {
	s.RegisterService(&Overlord_ServiceDesc, srv)
}

func _Overlord_GetIntegratorConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetIntegratorConfigIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OverlordServer).GetIntegratorConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/overlord.Overlord/GetIntegratorConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OverlordServer).GetIntegratorConfig(ctx, req.(*GetIntegratorConfigIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _Overlord_SaveParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveParamsIn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OverlordServer).SaveParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/overlord.Overlord/SaveParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OverlordServer).SaveParams(ctx, req.(*SaveParamsIn))
	}
	return interceptor(ctx, in, info, handler)
}

func _Overlord_HealthCheck_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(OverlordServer).HealthCheck(&overlordHealthCheckServer{stream})
}

type Overlord_HealthCheckServer interface {
	Send(*Status) error
	Recv() (*Status, error)
	grpc.ServerStream
}

type overlordHealthCheckServer struct {
	grpc.ServerStream
}

func (x *overlordHealthCheckServer) Send(m *Status) error {
	return x.ServerStream.SendMsg(m)
}

func (x *overlordHealthCheckServer) Recv() (*Status, error) {
	m := new(Status)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Overlord_ServiceDesc is the grpc.ServiceDesc for Overlord service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Overlord_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "overlord.Overlord",
	HandlerType: (*OverlordServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetIntegratorConfig",
			Handler:    _Overlord_GetIntegratorConfig_Handler,
		},
		{
			MethodName: "SaveParams",
			Handler:    _Overlord_SaveParams_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "HealthCheck",
			Handler:       _Overlord_HealthCheck_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pkg/overlord/overlord.proto",
}
