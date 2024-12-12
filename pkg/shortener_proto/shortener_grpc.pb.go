// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.0
// source: pkg/shortener_proto/shortener.proto

package __

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	GrpcShortener_UserRegister_FullMethodName    = "/shortener_proto.GrpcShortener/UserRegister"
	GrpcShortener_UserAuth_FullMethodName        = "/shortener_proto.GrpcShortener/UserAuth"
	GrpcShortener_EncodeLink_FullMethodName      = "/shortener_proto.GrpcShortener/EncodeLink"
	GrpcShortener_BatchEncodeLink_FullMethodName = "/shortener_proto.GrpcShortener/BatchEncodeLink"
	GrpcShortener_UserLinks_FullMethodName       = "/shortener_proto.GrpcShortener/UserLinks"
	GrpcShortener_DeleteUserLinks_FullMethodName = "/shortener_proto.GrpcShortener/DeleteUserLinks"
	GrpcShortener_InternalStats_FullMethodName   = "/shortener_proto.GrpcShortener/InternalStats"
	GrpcShortener_DBPing_FullMethodName          = "/shortener_proto.GrpcShortener/DBPing"
)

// GrpcShortenerClient is the client API for GrpcShortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// ---------------- сервер -------------------------------------------
type GrpcShortenerClient interface {
	UserRegister(ctx context.Context, in *UserRegisterRequest, opts ...grpc.CallOption) (*UserRegisterResponse, error)
	UserAuth(ctx context.Context, in *UserAuthRequest, opts ...grpc.CallOption) (*UserAuthResponse, error)
	EncodeLink(ctx context.Context, in *EncodeLinkRequest, opts ...grpc.CallOption) (*EncodeLinkResponse, error)
	BatchEncodeLink(ctx context.Context, in *EncodeBatchLinksRequest, opts ...grpc.CallOption) (*EncodeBatchLinksResponse, error)
	UserLinks(ctx context.Context, in *UserLinksRequest, opts ...grpc.CallOption) (*UserLinksResponse, error)
	DeleteUserLinks(ctx context.Context, in *DeleteUserLinksRequest, opts ...grpc.CallOption) (*DeleteUserLinksResponse, error)
	InternalStats(ctx context.Context, in *InternalStatsRequest, opts ...grpc.CallOption) (*InternalStatsResponse, error)
	DBPing(ctx context.Context, in *DBPingRequest, opts ...grpc.CallOption) (*DBPingResponse, error)
}

type grpcShortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewGrpcShortenerClient(cc grpc.ClientConnInterface) GrpcShortenerClient {
	return &grpcShortenerClient{cc}
}

func (c *grpcShortenerClient) UserRegister(ctx context.Context, in *UserRegisterRequest, opts ...grpc.CallOption) (*UserRegisterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserRegisterResponse)
	err := c.cc.Invoke(ctx, GrpcShortener_UserRegister_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcShortenerClient) UserAuth(ctx context.Context, in *UserAuthRequest, opts ...grpc.CallOption) (*UserAuthResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserAuthResponse)
	err := c.cc.Invoke(ctx, GrpcShortener_UserAuth_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcShortenerClient) EncodeLink(ctx context.Context, in *EncodeLinkRequest, opts ...grpc.CallOption) (*EncodeLinkResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EncodeLinkResponse)
	err := c.cc.Invoke(ctx, GrpcShortener_EncodeLink_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcShortenerClient) BatchEncodeLink(ctx context.Context, in *EncodeBatchLinksRequest, opts ...grpc.CallOption) (*EncodeBatchLinksResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EncodeBatchLinksResponse)
	err := c.cc.Invoke(ctx, GrpcShortener_BatchEncodeLink_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcShortenerClient) UserLinks(ctx context.Context, in *UserLinksRequest, opts ...grpc.CallOption) (*UserLinksResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserLinksResponse)
	err := c.cc.Invoke(ctx, GrpcShortener_UserLinks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcShortenerClient) DeleteUserLinks(ctx context.Context, in *DeleteUserLinksRequest, opts ...grpc.CallOption) (*DeleteUserLinksResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteUserLinksResponse)
	err := c.cc.Invoke(ctx, GrpcShortener_DeleteUserLinks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcShortenerClient) InternalStats(ctx context.Context, in *InternalStatsRequest, opts ...grpc.CallOption) (*InternalStatsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(InternalStatsResponse)
	err := c.cc.Invoke(ctx, GrpcShortener_InternalStats_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *grpcShortenerClient) DBPing(ctx context.Context, in *DBPingRequest, opts ...grpc.CallOption) (*DBPingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DBPingResponse)
	err := c.cc.Invoke(ctx, GrpcShortener_DBPing_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GrpcShortenerServer is the server API for GrpcShortener service.
// All implementations must embed UnimplementedGrpcShortenerServer
// for forward compatibility.
//
// ---------------- сервер -------------------------------------------
type GrpcShortenerServer interface {
	UserRegister(context.Context, *UserRegisterRequest) (*UserRegisterResponse, error)
	UserAuth(context.Context, *UserAuthRequest) (*UserAuthResponse, error)
	EncodeLink(context.Context, *EncodeLinkRequest) (*EncodeLinkResponse, error)
	BatchEncodeLink(context.Context, *EncodeBatchLinksRequest) (*EncodeBatchLinksResponse, error)
	UserLinks(context.Context, *UserLinksRequest) (*UserLinksResponse, error)
	DeleteUserLinks(context.Context, *DeleteUserLinksRequest) (*DeleteUserLinksResponse, error)
	InternalStats(context.Context, *InternalStatsRequest) (*InternalStatsResponse, error)
	DBPing(context.Context, *DBPingRequest) (*DBPingResponse, error)
	mustEmbedUnimplementedGrpcShortenerServer()
}

// UnimplementedGrpcShortenerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGrpcShortenerServer struct{}

func (UnimplementedGrpcShortenerServer) UserRegister(context.Context, *UserRegisterRequest) (*UserRegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserRegister not implemented")
}
func (UnimplementedGrpcShortenerServer) UserAuth(context.Context, *UserAuthRequest) (*UserAuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserAuth not implemented")
}
func (UnimplementedGrpcShortenerServer) EncodeLink(context.Context, *EncodeLinkRequest) (*EncodeLinkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EncodeLink not implemented")
}
func (UnimplementedGrpcShortenerServer) BatchEncodeLink(context.Context, *EncodeBatchLinksRequest) (*EncodeBatchLinksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchEncodeLink not implemented")
}
func (UnimplementedGrpcShortenerServer) UserLinks(context.Context, *UserLinksRequest) (*UserLinksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserLinks not implemented")
}
func (UnimplementedGrpcShortenerServer) DeleteUserLinks(context.Context, *DeleteUserLinksRequest) (*DeleteUserLinksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserLinks not implemented")
}
func (UnimplementedGrpcShortenerServer) InternalStats(context.Context, *InternalStatsRequest) (*InternalStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InternalStats not implemented")
}
func (UnimplementedGrpcShortenerServer) DBPing(context.Context, *DBPingRequest) (*DBPingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DBPing not implemented")
}
func (UnimplementedGrpcShortenerServer) mustEmbedUnimplementedGrpcShortenerServer() {}
func (UnimplementedGrpcShortenerServer) testEmbeddedByValue()                       {}

// UnsafeGrpcShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GrpcShortenerServer will
// result in compilation errors.
type UnsafeGrpcShortenerServer interface {
	mustEmbedUnimplementedGrpcShortenerServer()
}

func RegisterGrpcShortenerServer(s grpc.ServiceRegistrar, srv GrpcShortenerServer) {
	// If the following call pancis, it indicates UnimplementedGrpcShortenerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&GrpcShortener_ServiceDesc, srv)
}

func _GrpcShortener_UserRegister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcShortenerServer).UserRegister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GrpcShortener_UserRegister_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcShortenerServer).UserRegister(ctx, req.(*UserRegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GrpcShortener_UserAuth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserAuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcShortenerServer).UserAuth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GrpcShortener_UserAuth_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcShortenerServer).UserAuth(ctx, req.(*UserAuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GrpcShortener_EncodeLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EncodeLinkRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcShortenerServer).EncodeLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GrpcShortener_EncodeLink_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcShortenerServer).EncodeLink(ctx, req.(*EncodeLinkRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GrpcShortener_BatchEncodeLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EncodeBatchLinksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcShortenerServer).BatchEncodeLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GrpcShortener_BatchEncodeLink_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcShortenerServer).BatchEncodeLink(ctx, req.(*EncodeBatchLinksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GrpcShortener_UserLinks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserLinksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcShortenerServer).UserLinks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GrpcShortener_UserLinks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcShortenerServer).UserLinks(ctx, req.(*UserLinksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GrpcShortener_DeleteUserLinks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserLinksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcShortenerServer).DeleteUserLinks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GrpcShortener_DeleteUserLinks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcShortenerServer).DeleteUserLinks(ctx, req.(*DeleteUserLinksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GrpcShortener_InternalStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InternalStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcShortenerServer).InternalStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GrpcShortener_InternalStats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcShortenerServer).InternalStats(ctx, req.(*InternalStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GrpcShortener_DBPing_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DBPingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GrpcShortenerServer).DBPing(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GrpcShortener_DBPing_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GrpcShortenerServer).DBPing(ctx, req.(*DBPingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GrpcShortener_ServiceDesc is the grpc.ServiceDesc for GrpcShortener service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GrpcShortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener_proto.GrpcShortener",
	HandlerType: (*GrpcShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserRegister",
			Handler:    _GrpcShortener_UserRegister_Handler,
		},
		{
			MethodName: "UserAuth",
			Handler:    _GrpcShortener_UserAuth_Handler,
		},
		{
			MethodName: "EncodeLink",
			Handler:    _GrpcShortener_EncodeLink_Handler,
		},
		{
			MethodName: "BatchEncodeLink",
			Handler:    _GrpcShortener_BatchEncodeLink_Handler,
		},
		{
			MethodName: "UserLinks",
			Handler:    _GrpcShortener_UserLinks_Handler,
		},
		{
			MethodName: "DeleteUserLinks",
			Handler:    _GrpcShortener_DeleteUserLinks_Handler,
		},
		{
			MethodName: "InternalStats",
			Handler:    _GrpcShortener_InternalStats_Handler,
		},
		{
			MethodName: "DBPing",
			Handler:    _GrpcShortener_DBPing_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/shortener_proto/shortener.proto",
}
