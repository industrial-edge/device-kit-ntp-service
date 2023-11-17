// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: Ntp.proto

package siemens_iedge_dmapi_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	NtpService_SetNtpServer_FullMethodName = "/siemens.iedge.dmapi.ntp.v1.NtpService/SetNtpServer"
	NtpService_GetNtpServer_FullMethodName = "/siemens.iedge.dmapi.ntp.v1.NtpService/GetNtpServer"
	NtpService_GetStatus_FullMethodName    = "/siemens.iedge.dmapi.ntp.v1.NtpService/GetStatus"
)

// NtpServiceClient is the client API for NtpService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NtpServiceClient interface {
	// Set ntp server
	SetNtpServer(ctx context.Context, in *Ntp, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Returns ntp servers
	GetNtpServer(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Ntp, error)
	// Returns NTP Status message.
	GetStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Status, error)
}

type ntpServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNtpServiceClient(cc grpc.ClientConnInterface) NtpServiceClient {
	return &ntpServiceClient{cc}
}

func (c *ntpServiceClient) SetNtpServer(ctx context.Context, in *Ntp, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, NtpService_SetNtpServer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ntpServiceClient) GetNtpServer(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Ntp, error) {
	out := new(Ntp)
	err := c.cc.Invoke(ctx, NtpService_GetNtpServer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ntpServiceClient) GetStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, NtpService_GetStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NtpServiceServer is the server API for NtpService service.
// All implementations must embed UnimplementedNtpServiceServer
// for forward compatibility
type NtpServiceServer interface {
	// Set ntp server
	SetNtpServer(context.Context, *Ntp) (*emptypb.Empty, error)
	// Returns ntp servers
	GetNtpServer(context.Context, *emptypb.Empty) (*Ntp, error)
	// Returns NTP Status message.
	GetStatus(context.Context, *emptypb.Empty) (*Status, error)
	mustEmbedUnimplementedNtpServiceServer()
}

// UnimplementedNtpServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNtpServiceServer struct {
}

func (UnimplementedNtpServiceServer) SetNtpServer(context.Context, *Ntp) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetNtpServer not implemented")
}
func (UnimplementedNtpServiceServer) GetNtpServer(context.Context, *emptypb.Empty) (*Ntp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNtpServer not implemented")
}
func (UnimplementedNtpServiceServer) GetStatus(context.Context, *emptypb.Empty) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedNtpServiceServer) mustEmbedUnimplementedNtpServiceServer() {}

// UnsafeNtpServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NtpServiceServer will
// result in compilation errors.
type UnsafeNtpServiceServer interface {
	mustEmbedUnimplementedNtpServiceServer()
}

func RegisterNtpServiceServer(s grpc.ServiceRegistrar, srv NtpServiceServer) {
	s.RegisterService(&NtpService_ServiceDesc, srv)
}

func _NtpService_SetNtpServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ntp)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NtpServiceServer).SetNtpServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NtpService_SetNtpServer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NtpServiceServer).SetNtpServer(ctx, req.(*Ntp))
	}
	return interceptor(ctx, in, info, handler)
}

func _NtpService_GetNtpServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NtpServiceServer).GetNtpServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NtpService_GetNtpServer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NtpServiceServer).GetNtpServer(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _NtpService_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NtpServiceServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: NtpService_GetStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NtpServiceServer).GetStatus(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// NtpService_ServiceDesc is the grpc.ServiceDesc for NtpService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NtpService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "siemens.iedge.dmapi.ntp.v1.NtpService",
	HandlerType: (*NtpServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetNtpServer",
			Handler:    _NtpService_SetNtpServer_Handler,
		},
		{
			MethodName: "GetNtpServer",
			Handler:    _NtpService_GetNtpServer_Handler,
		},
		{
			MethodName: "GetStatus",
			Handler:    _NtpService_GetStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Ntp.proto",
}
