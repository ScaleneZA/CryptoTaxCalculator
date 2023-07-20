// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.22.2
// source: conversionrate.proto

package conversionratepb

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
	Conversionrate_ValueAtTime_FullMethodName = "/conversionratepb.Conversionrate/ValueAtTime"
)

// ConversionrateClient is the client API for Conversionrate service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConversionrateClient interface {
	ValueAtTime(ctx context.Context, in *ValueAtTimeRequest, opts ...grpc.CallOption) (*ValueAtTimeResponse, error)
}

type conversionrateClient struct {
	cc grpc.ClientConnInterface
}

func NewConversionrateClient(cc grpc.ClientConnInterface) ConversionrateClient {
	return &conversionrateClient{cc}
}

func (c *conversionrateClient) ValueAtTime(ctx context.Context, in *ValueAtTimeRequest, opts ...grpc.CallOption) (*ValueAtTimeResponse, error) {
	out := new(ValueAtTimeResponse)
	err := c.cc.Invoke(ctx, Conversionrate_ValueAtTime_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConversionrateServer is the server API for Conversionrate service.
// All implementations must embed UnimplementedConversionrateServer
// for forward compatibility
type ConversionrateServer interface {
	ValueAtTime(context.Context, *ValueAtTimeRequest) (*ValueAtTimeResponse, error)
	mustEmbedUnimplementedConversionrateServer()
}

// UnimplementedConversionrateServer must be embedded to have forward compatible implementations.
type UnimplementedConversionrateServer struct {
}

func (UnimplementedConversionrateServer) ValueAtTime(context.Context, *ValueAtTimeRequest) (*ValueAtTimeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValueAtTime not implemented")
}
func (UnimplementedConversionrateServer) mustEmbedUnimplementedConversionrateServer() {}

// UnsafeConversionrateServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConversionrateServer will
// result in compilation errors.
type UnsafeConversionrateServer interface {
	mustEmbedUnimplementedConversionrateServer()
}

func RegisterConversionrateServer(s grpc.ServiceRegistrar, srv ConversionrateServer) {
	s.RegisterService(&Conversionrate_ServiceDesc, srv)
}

func _Conversionrate_ValueAtTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValueAtTimeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConversionrateServer).ValueAtTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Conversionrate_ValueAtTime_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConversionrateServer).ValueAtTime(ctx, req.(*ValueAtTimeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Conversionrate_ServiceDesc is the grpc.ServiceDesc for Conversionrate service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Conversionrate_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "conversionratepb.Conversionrate",
	HandlerType: (*ConversionrateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ValueAtTime",
			Handler:    _Conversionrate_ValueAtTime_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "conversionrate.proto",
}