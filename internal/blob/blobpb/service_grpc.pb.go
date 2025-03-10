// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/blob/blobpb/*.proto

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: internal/blob/blobpb/service.proto

package blobpb

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
	StorageService_UploadObject_FullMethodName         = "/blobpb.StorageService/UploadObject"
	StorageService_DownloadObject_FullMethodName       = "/blobpb.StorageService/DownloadObject"
	StorageService_DeleteObject_FullMethodName         = "/blobpb.StorageService/DeleteObject"
	StorageService_ShareObject_FullMethodName          = "/blobpb.StorageService/ShareObject"
	StorageService_DownloadSharedObject_FullMethodName = "/blobpb.StorageService/DownloadSharedObject"
	StorageService_ListObject_FullMethodName           = "/blobpb.StorageService/ListObject"
)

// StorageServiceClient is the client API for StorageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StorageServiceClient interface {
	UploadObject(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[UploadRequest, UploadResponse], error)
	DownloadObject(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DownloadResponse], error)
	DeleteObject(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	ShareObject(ctx context.Context, in *ShareObjectRequest, opts ...grpc.CallOption) (*ShareObjectResponse, error)
	DownloadSharedObject(ctx context.Context, in *DownloadSharedRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DownloadResponse], error)
	ListObject(ctx context.Context, in *ListObjectRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ListObjectResponse], error)
}

type storageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStorageServiceClient(cc grpc.ClientConnInterface) StorageServiceClient {
	return &storageServiceClient{cc}
}

func (c *storageServiceClient) UploadObject(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[UploadRequest, UploadResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &StorageService_ServiceDesc.Streams[0], StorageService_UploadObject_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[UploadRequest, UploadResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StorageService_UploadObjectClient = grpc.ClientStreamingClient[UploadRequest, UploadResponse]

func (c *storageServiceClient) DownloadObject(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DownloadResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &StorageService_ServiceDesc.Streams[1], StorageService_DownloadObject_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[DownloadRequest, DownloadResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StorageService_DownloadObjectClient = grpc.ServerStreamingClient[DownloadResponse]

func (c *storageServiceClient) DeleteObject(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, StorageService_DeleteObject_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) ShareObject(ctx context.Context, in *ShareObjectRequest, opts ...grpc.CallOption) (*ShareObjectResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ShareObjectResponse)
	err := c.cc.Invoke(ctx, StorageService_ShareObject_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageServiceClient) DownloadSharedObject(ctx context.Context, in *DownloadSharedRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DownloadResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &StorageService_ServiceDesc.Streams[2], StorageService_DownloadSharedObject_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[DownloadSharedRequest, DownloadResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StorageService_DownloadSharedObjectClient = grpc.ServerStreamingClient[DownloadResponse]

func (c *storageServiceClient) ListObject(ctx context.Context, in *ListObjectRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ListObjectResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &StorageService_ServiceDesc.Streams[3], StorageService_ListObject_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ListObjectRequest, ListObjectResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StorageService_ListObjectClient = grpc.ServerStreamingClient[ListObjectResponse]

// StorageServiceServer is the server API for StorageService service.
// All implementations must embed UnimplementedStorageServiceServer
// for forward compatibility.
type StorageServiceServer interface {
	UploadObject(grpc.ClientStreamingServer[UploadRequest, UploadResponse]) error
	DownloadObject(*DownloadRequest, grpc.ServerStreamingServer[DownloadResponse]) error
	DeleteObject(context.Context, *DeleteRequest) (*DeleteResponse, error)
	ShareObject(context.Context, *ShareObjectRequest) (*ShareObjectResponse, error)
	DownloadSharedObject(*DownloadSharedRequest, grpc.ServerStreamingServer[DownloadResponse]) error
	ListObject(*ListObjectRequest, grpc.ServerStreamingServer[ListObjectResponse]) error
	mustEmbedUnimplementedStorageServiceServer()
}

// UnimplementedStorageServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedStorageServiceServer struct{}

func (UnimplementedStorageServiceServer) UploadObject(grpc.ClientStreamingServer[UploadRequest, UploadResponse]) error {
	return status.Errorf(codes.Unimplemented, "method UploadObject not implemented")
}
func (UnimplementedStorageServiceServer) DownloadObject(*DownloadRequest, grpc.ServerStreamingServer[DownloadResponse]) error {
	return status.Errorf(codes.Unimplemented, "method DownloadObject not implemented")
}
func (UnimplementedStorageServiceServer) DeleteObject(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteObject not implemented")
}
func (UnimplementedStorageServiceServer) ShareObject(context.Context, *ShareObjectRequest) (*ShareObjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShareObject not implemented")
}
func (UnimplementedStorageServiceServer) DownloadSharedObject(*DownloadSharedRequest, grpc.ServerStreamingServer[DownloadResponse]) error {
	return status.Errorf(codes.Unimplemented, "method DownloadSharedObject not implemented")
}
func (UnimplementedStorageServiceServer) ListObject(*ListObjectRequest, grpc.ServerStreamingServer[ListObjectResponse]) error {
	return status.Errorf(codes.Unimplemented, "method ListObject not implemented")
}
func (UnimplementedStorageServiceServer) mustEmbedUnimplementedStorageServiceServer() {}
func (UnimplementedStorageServiceServer) testEmbeddedByValue()                        {}

// UnsafeStorageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StorageServiceServer will
// result in compilation errors.
type UnsafeStorageServiceServer interface {
	mustEmbedUnimplementedStorageServiceServer()
}

func RegisterStorageServiceServer(s grpc.ServiceRegistrar, srv StorageServiceServer) {
	// If the following call pancis, it indicates UnimplementedStorageServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&StorageService_ServiceDesc, srv)
}

func _StorageService_UploadObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(StorageServiceServer).UploadObject(&grpc.GenericServerStream[UploadRequest, UploadResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StorageService_UploadObjectServer = grpc.ClientStreamingServer[UploadRequest, UploadResponse]

func _StorageService_DownloadObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServiceServer).DownloadObject(m, &grpc.GenericServerStream[DownloadRequest, DownloadResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StorageService_DownloadObjectServer = grpc.ServerStreamingServer[DownloadResponse]

func _StorageService_DeleteObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).DeleteObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StorageService_DeleteObject_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).DeleteObject(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_ShareObject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShareObjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServiceServer).ShareObject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StorageService_ShareObject_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServiceServer).ShareObject(ctx, req.(*ShareObjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StorageService_DownloadSharedObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DownloadSharedRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServiceServer).DownloadSharedObject(m, &grpc.GenericServerStream[DownloadSharedRequest, DownloadResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StorageService_DownloadSharedObjectServer = grpc.ServerStreamingServer[DownloadResponse]

func _StorageService_ListObject_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ListObjectRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StorageServiceServer).ListObject(m, &grpc.GenericServerStream[ListObjectRequest, ListObjectResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type StorageService_ListObjectServer = grpc.ServerStreamingServer[ListObjectResponse]

// StorageService_ServiceDesc is the grpc.ServiceDesc for StorageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StorageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "blobpb.StorageService",
	HandlerType: (*StorageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteObject",
			Handler:    _StorageService_DeleteObject_Handler,
		},
		{
			MethodName: "ShareObject",
			Handler:    _StorageService_ShareObject_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadObject",
			Handler:       _StorageService_UploadObject_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "DownloadObject",
			Handler:       _StorageService_DownloadObject_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "DownloadSharedObject",
			Handler:       _StorageService_DownloadSharedObject_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ListObject",
			Handler:       _StorageService_ListObject_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "internal/blob/blobpb/service.proto",
}
