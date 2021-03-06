// Code generated by protoc-gen-gogo.
// source: pb/collector.proto
// DO NOT EDIT!

/*
	Package pb is a generated protocol buffer package.

	It is generated from these files:
		pb/collector.proto
		pb/exporter.proto
		pb/metric.proto

	It has these top-level messages:
		MeteringReqResp
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for CollectorService service

type CollectorServiceClient interface {
	Transit(ctx context.Context, in *MeteringReqResp, opts ...grpc.CallOption) (*MeteringReqResp, error)
	BatchTransit(ctx context.Context, opts ...grpc.CallOption) (CollectorService_BatchTransitClient, error)
}

type collectorServiceClient struct {
	cc *grpc.ClientConn
}

func NewCollectorServiceClient(cc *grpc.ClientConn) CollectorServiceClient {
	return &collectorServiceClient{cc}
}

func (c *collectorServiceClient) Transit(ctx context.Context, in *MeteringReqResp, opts ...grpc.CallOption) (*MeteringReqResp, error) {
	out := new(MeteringReqResp)
	err := grpc.Invoke(ctx, "/pb.CollectorService/Transit", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectorServiceClient) BatchTransit(ctx context.Context, opts ...grpc.CallOption) (CollectorService_BatchTransitClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_CollectorService_serviceDesc.Streams[0], c.cc, "/pb.CollectorService/BatchTransit", opts...)
	if err != nil {
		return nil, err
	}
	x := &collectorServiceBatchTransitClient{stream}
	return x, nil
}

type CollectorService_BatchTransitClient interface {
	Send(*MeteringReqResp) error
	Recv() (*MeteringReqResp, error)
	grpc.ClientStream
}

type collectorServiceBatchTransitClient struct {
	grpc.ClientStream
}

func (x *collectorServiceBatchTransitClient) Send(m *MeteringReqResp) error {
	return x.ClientStream.SendMsg(m)
}

func (x *collectorServiceBatchTransitClient) Recv() (*MeteringReqResp, error) {
	m := new(MeteringReqResp)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for CollectorService service

type CollectorServiceServer interface {
	Transit(context.Context, *MeteringReqResp) (*MeteringReqResp, error)
	BatchTransit(CollectorService_BatchTransitServer) error
}

func RegisterCollectorServiceServer(s *grpc.Server, srv CollectorServiceServer) {
	s.RegisterService(&_CollectorService_serviceDesc, srv)
}

func _CollectorService_Transit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MeteringReqResp)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CollectorServiceServer).Transit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.CollectorService/Transit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CollectorServiceServer).Transit(ctx, req.(*MeteringReqResp))
	}
	return interceptor(ctx, in, info, handler)
}

func _CollectorService_BatchTransit_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CollectorServiceServer).BatchTransit(&collectorServiceBatchTransitServer{stream})
}

type CollectorService_BatchTransitServer interface {
	Send(*MeteringReqResp) error
	Recv() (*MeteringReqResp, error)
	grpc.ServerStream
}

type collectorServiceBatchTransitServer struct {
	grpc.ServerStream
}

func (x *collectorServiceBatchTransitServer) Send(m *MeteringReqResp) error {
	return x.ServerStream.SendMsg(m)
}

func (x *collectorServiceBatchTransitServer) Recv() (*MeteringReqResp, error) {
	m := new(MeteringReqResp)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _CollectorService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.CollectorService",
	HandlerType: (*CollectorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Transit",
			Handler:    _CollectorService_Transit_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "BatchTransit",
			Handler:       _CollectorService_BatchTransit_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pb/collector.proto",
}

func init() { proto.RegisterFile("pb/collector.proto", fileDescriptorCollector) }

var fileDescriptorCollector = []byte{
	// 215 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2a, 0x48, 0xd2, 0x4f,
	0xce, 0xcf, 0xc9, 0x49, 0x4d, 0x2e, 0xc9, 0x2f, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0x2a, 0x48, 0x92, 0x92, 0x49, 0xcf, 0xcf, 0x4f, 0xcf, 0x49, 0xd5, 0x4f, 0x2c, 0xc8, 0xd4, 0x4f,
	0xcc, 0xcb, 0xcb, 0x2f, 0x49, 0x2c, 0xc9, 0xcc, 0xcf, 0x2b, 0x86, 0xa8, 0x90, 0xe2, 0x2f, 0x48,
	0xd2, 0xcf, 0x4d, 0x2d, 0x29, 0xca, 0x4c, 0x86, 0x08, 0x18, 0xed, 0x65, 0xe4, 0x12, 0x70, 0x86,
	0x19, 0x13, 0x9c, 0x5a, 0x54, 0x96, 0x99, 0x9c, 0x2a, 0xe4, 0xc3, 0xc5, 0x1e, 0x52, 0x94, 0x98,
	0x57, 0x9c, 0x59, 0x22, 0x24, 0xac, 0x57, 0x90, 0xa4, 0xe7, 0x9b, 0x5a, 0x92, 0x5a, 0x94, 0x99,
	0x97, 0x1e, 0x94, 0x5a, 0x18, 0x94, 0x5a, 0x5c, 0x20, 0x85, 0x4d, 0x50, 0x49, 0xbc, 0xe9, 0xf2,
	0x93, 0xc9, 0x4c, 0x82, 0x4a, 0x3c, 0xfa, 0x65, 0x86, 0x20, 0x3b, 0xc0, 0x92, 0x56, 0x8c, 0x5a,
	0x42, 0xd1, 0x5c, 0x3c, 0x4e, 0x89, 0x25, 0xc9, 0x19, 0xa4, 0x1b, 0x29, 0x0b, 0x36, 0x52, 0x5c,
	0x49, 0x08, 0xd9, 0xc8, 0xe2, 0x92, 0xa2, 0xd4, 0xc4, 0x5c, 0x2b, 0x46, 0x2d, 0x0d, 0x46, 0x03,
	0x46, 0x27, 0x81, 0x13, 0x8f, 0xe4, 0x18, 0x2f, 0x3c, 0x92, 0x63, 0x7c, 0xf0, 0x48, 0x8e, 0x71,
	0xc6, 0x63, 0x39, 0x86, 0x24, 0x36, 0xb0, 0xc7, 0x8c, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x60,
	0x39, 0x89, 0xa3, 0x21, 0x01, 0x00, 0x00,
}
