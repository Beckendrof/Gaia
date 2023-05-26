// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package apostolis

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

// ApostolisClient is the client API for Apostolis service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ApostolisClient interface {
	System(ctx context.Context, in *ApostolisRequest, opts ...grpc.CallOption) (Apostolis_SystemClient, error)
}

type apostolisClient struct {
	cc grpc.ClientConnInterface
}

func NewApostolisClient(cc grpc.ClientConnInterface) ApostolisClient {
	return &apostolisClient{cc}
}

func (c *apostolisClient) System(ctx context.Context, in *ApostolisRequest, opts ...grpc.CallOption) (Apostolis_SystemClient, error) {
	stream, err := c.cc.NewStream(ctx, &Apostolis_ServiceDesc.Streams[0], "/gaia.apostolis.api.Apostolis/System", opts...)
	if err != nil {
		return nil, err
	}
	x := &apostolisSystemClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Apostolis_SystemClient interface {
	Recv() (*ApostolisReply, error)
	grpc.ClientStream
}

type apostolisSystemClient struct {
	grpc.ClientStream
}

func (x *apostolisSystemClient) Recv() (*ApostolisReply, error) {
	m := new(ApostolisReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ApostolisServer is the server API for Apostolis service.
// All implementations must embed UnimplementedApostolisServer
// for forward compatibility
type ApostolisServer interface {
	System(*ApostolisRequest, Apostolis_SystemServer) error
	mustEmbedUnimplementedApostolisServer()
}

// UnimplementedApostolisServer must be embedded to have forward compatible implementations.
type UnimplementedApostolisServer struct {
}

func (UnimplementedApostolisServer) System(*ApostolisRequest, Apostolis_SystemServer) error {
	return status.Errorf(codes.Unimplemented, "method System not implemented")
}
func (UnimplementedApostolisServer) mustEmbedUnimplementedApostolisServer() {}

// UnsafeApostolisServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ApostolisServer will
// result in compilation errors.
type UnsafeApostolisServer interface {
	mustEmbedUnimplementedApostolisServer()
}

func RegisterApostolisServer(s grpc.ServiceRegistrar, srv ApostolisServer) {
	s.RegisterService(&Apostolis_ServiceDesc, srv)
}

func _Apostolis_System_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ApostolisRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ApostolisServer).System(m, &apostolisSystemServer{stream})
}

type Apostolis_SystemServer interface {
	Send(*ApostolisReply) error
	grpc.ServerStream
}

type apostolisSystemServer struct {
	grpc.ServerStream
}

func (x *apostolisSystemServer) Send(m *ApostolisReply) error {
	return x.ServerStream.SendMsg(m)
}

// Apostolis_ServiceDesc is the grpc.ServiceDesc for Apostolis service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Apostolis_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gaia.apostolis.api.Apostolis",
	HandlerType: (*ApostolisServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "System",
			Handler:       _Apostolis_System_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "src/models/protos/apostolis/apostolis_service.proto",
}