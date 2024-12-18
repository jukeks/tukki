// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: tukki/rpc/kv/v1/service.proto

package kvv1

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
	KvService_Query_FullMethodName       = "/tukki.rpc.kv.v1.KvService/Query"
	KvService_QueryRange_FullMethodName  = "/tukki.rpc.kv.v1.KvService/QueryRange"
	KvService_Set_FullMethodName         = "/tukki.rpc.kv.v1.KvService/Set"
	KvService_Delete_FullMethodName      = "/tukki.rpc.kv.v1.KvService/Delete"
	KvService_DeleteRange_FullMethodName = "/tukki.rpc.kv.v1.KvService/DeleteRange"
)

// KvServiceClient is the client API for KvService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KvServiceClient interface {
	Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (*QueryResponse, error)
	QueryRange(ctx context.Context, in *QueryRangeRequest, opts ...grpc.CallOption) (KvService_QueryRangeClient, error)
	Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	DeleteRange(ctx context.Context, in *DeleteRangeRequest, opts ...grpc.CallOption) (*DeleteRangeResponse, error)
}

type kvServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKvServiceClient(cc grpc.ClientConnInterface) KvServiceClient {
	return &kvServiceClient{cc}
}

func (c *kvServiceClient) Query(ctx context.Context, in *QueryRequest, opts ...grpc.CallOption) (*QueryResponse, error) {
	out := new(QueryResponse)
	err := c.cc.Invoke(ctx, KvService_Query_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kvServiceClient) QueryRange(ctx context.Context, in *QueryRangeRequest, opts ...grpc.CallOption) (KvService_QueryRangeClient, error) {
	stream, err := c.cc.NewStream(ctx, &KvService_ServiceDesc.Streams[0], KvService_QueryRange_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &kvServiceQueryRangeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type KvService_QueryRangeClient interface {
	Recv() (*QueryRangeResponse, error)
	grpc.ClientStream
}

type kvServiceQueryRangeClient struct {
	grpc.ClientStream
}

func (x *kvServiceQueryRangeClient) Recv() (*QueryRangeResponse, error) {
	m := new(QueryRangeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *kvServiceClient) Set(ctx context.Context, in *SetRequest, opts ...grpc.CallOption) (*SetResponse, error) {
	out := new(SetResponse)
	err := c.cc.Invoke(ctx, KvService_Set_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kvServiceClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, KvService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *kvServiceClient) DeleteRange(ctx context.Context, in *DeleteRangeRequest, opts ...grpc.CallOption) (*DeleteRangeResponse, error) {
	out := new(DeleteRangeResponse)
	err := c.cc.Invoke(ctx, KvService_DeleteRange_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KvServiceServer is the server API for KvService service.
// All implementations must embed UnimplementedKvServiceServer
// for forward compatibility
type KvServiceServer interface {
	Query(context.Context, *QueryRequest) (*QueryResponse, error)
	QueryRange(*QueryRangeRequest, KvService_QueryRangeServer) error
	Set(context.Context, *SetRequest) (*SetResponse, error)
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	DeleteRange(context.Context, *DeleteRangeRequest) (*DeleteRangeResponse, error)
	mustEmbedUnimplementedKvServiceServer()
}

// UnimplementedKvServiceServer must be embedded to have forward compatible implementations.
type UnimplementedKvServiceServer struct {
}

func (UnimplementedKvServiceServer) Query(context.Context, *QueryRequest) (*QueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (UnimplementedKvServiceServer) QueryRange(*QueryRangeRequest, KvService_QueryRangeServer) error {
	return status.Errorf(codes.Unimplemented, "method QueryRange not implemented")
}
func (UnimplementedKvServiceServer) Set(context.Context, *SetRequest) (*SetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedKvServiceServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedKvServiceServer) DeleteRange(context.Context, *DeleteRangeRequest) (*DeleteRangeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRange not implemented")
}
func (UnimplementedKvServiceServer) mustEmbedUnimplementedKvServiceServer() {}

// UnsafeKvServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KvServiceServer will
// result in compilation errors.
type UnsafeKvServiceServer interface {
	mustEmbedUnimplementedKvServiceServer()
}

func RegisterKvServiceServer(s grpc.ServiceRegistrar, srv KvServiceServer) {
	s.RegisterService(&KvService_ServiceDesc, srv)
}

func _KvService_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KvServiceServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KvService_Query_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KvServiceServer).Query(ctx, req.(*QueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KvService_QueryRange_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(QueryRangeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(KvServiceServer).QueryRange(m, &kvServiceQueryRangeServer{stream})
}

type KvService_QueryRangeServer interface {
	Send(*QueryRangeResponse) error
	grpc.ServerStream
}

type kvServiceQueryRangeServer struct {
	grpc.ServerStream
}

func (x *kvServiceQueryRangeServer) Send(m *QueryRangeResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _KvService_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KvServiceServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KvService_Set_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KvServiceServer).Set(ctx, req.(*SetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KvService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KvServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KvService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KvServiceServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KvService_DeleteRange_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KvServiceServer).DeleteRange(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: KvService_DeleteRange_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KvServiceServer).DeleteRange(ctx, req.(*DeleteRangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KvService_ServiceDesc is the grpc.ServiceDesc for KvService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KvService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tukki.rpc.kv.v1.KvService",
	HandlerType: (*KvServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Query",
			Handler:    _KvService_Query_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _KvService_Set_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _KvService_Delete_Handler,
		},
		{
			MethodName: "DeleteRange",
			Handler:    _KvService_DeleteRange_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "QueryRange",
			Handler:       _KvService_QueryRange_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "tukki/rpc/kv/v1/service.proto",
}
