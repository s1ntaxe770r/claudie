// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

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

// KuberServiceClient is the client API for KuberService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KuberServiceClient interface {
	SetUpStorage(ctx context.Context, in *SetUpStorageRequest, opts ...grpc.CallOption) (*SetUpStorageResponse, error)
}

type kuberServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKuberServiceClient(cc grpc.ClientConnInterface) KuberServiceClient {
	return &kuberServiceClient{cc}
}

func (c *kuberServiceClient) SetUpStorage(ctx context.Context, in *SetUpStorageRequest, opts ...grpc.CallOption) (*SetUpStorageResponse, error) {
	out := new(SetUpStorageResponse)
	err := c.cc.Invoke(ctx, "/platform.KuberService/SetUpStorage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KuberServiceServer is the server API for KuberService service.
// All implementations must embed UnimplementedKuberServiceServer
// for forward compatibility
type KuberServiceServer interface {
	SetUpStorage(context.Context, *SetUpStorageRequest) (*SetUpStorageResponse, error)
	mustEmbedUnimplementedKuberServiceServer()
}

// UnimplementedKuberServiceServer must be embedded to have forward compatible implementations.
type UnimplementedKuberServiceServer struct {
}

func (UnimplementedKuberServiceServer) SetUpStorage(context.Context, *SetUpStorageRequest) (*SetUpStorageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetUpStorage not implemented")
}
func (UnimplementedKuberServiceServer) mustEmbedUnimplementedKuberServiceServer() {}

// UnsafeKuberServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KuberServiceServer will
// result in compilation errors.
type UnsafeKuberServiceServer interface {
	mustEmbedUnimplementedKuberServiceServer()
}

func RegisterKuberServiceServer(s grpc.ServiceRegistrar, srv KuberServiceServer) {
	s.RegisterService(&KuberService_ServiceDesc, srv)
}

func _KuberService_SetUpStorage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetUpStorageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KuberServiceServer).SetUpStorage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/platform.KuberService/SetUpStorage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KuberServiceServer).SetUpStorage(ctx, req.(*SetUpStorageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KuberService_ServiceDesc is the grpc.ServiceDesc for KuberService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KuberService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "platform.KuberService",
	HandlerType: (*KuberServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetUpStorage",
			Handler:    _KuberService_SetUpStorage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/kuber.proto",
}