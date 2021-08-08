// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package writerService

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// WriterServiceClient is the client API for WriterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WriterServiceClient interface {
	CreateProduct(ctx context.Context, in *CreateProductReq, opts ...grpc.CallOption) (*CreateProductRes, error)
	UpdateProduct(ctx context.Context, in *UpdateProductReq, opts ...grpc.CallOption) (*UpdateProductRes, error)
	GetProductById(ctx context.Context, in *GetProductByIdReq, opts ...grpc.CallOption) (*GetProductByIdRes, error)
}

type writerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWriterServiceClient(cc grpc.ClientConnInterface) WriterServiceClient {
	return &writerServiceClient{cc}
}

func (c *writerServiceClient) CreateProduct(ctx context.Context, in *CreateProductReq, opts ...grpc.CallOption) (*CreateProductRes, error) {
	out := new(CreateProductRes)
	err := c.cc.Invoke(ctx, "/writerService.writerService/CreateProduct", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) UpdateProduct(ctx context.Context, in *UpdateProductReq, opts ...grpc.CallOption) (*UpdateProductRes, error) {
	out := new(UpdateProductRes)
	err := c.cc.Invoke(ctx, "/writerService.writerService/UpdateProduct", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerServiceClient) GetProductById(ctx context.Context, in *GetProductByIdReq, opts ...grpc.CallOption) (*GetProductByIdRes, error) {
	out := new(GetProductByIdRes)
	err := c.cc.Invoke(ctx, "/writerService.writerService/GetProductById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WriterServiceServer is the server API for WriterService service.
// All implementations should embed UnimplementedWriterServiceServer
// for forward compatibility
type WriterServiceServer interface {
	CreateProduct(context.Context, *CreateProductReq) (*CreateProductRes, error)
	UpdateProduct(context.Context, *UpdateProductReq) (*UpdateProductRes, error)
	GetProductById(context.Context, *GetProductByIdReq) (*GetProductByIdRes, error)
}

// UnimplementedWriterServiceServer should be embedded to have forward compatible implementations.
type UnimplementedWriterServiceServer struct {
}

func (UnimplementedWriterServiceServer) CreateProduct(context.Context, *CreateProductReq) (*CreateProductRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProduct not implemented")
}
func (UnimplementedWriterServiceServer) UpdateProduct(context.Context, *UpdateProductReq) (*UpdateProductRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProduct not implemented")
}
func (UnimplementedWriterServiceServer) GetProductById(context.Context, *GetProductByIdReq) (*GetProductByIdRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProductById not implemented")
}

// UnsafeWriterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WriterServiceServer will
// result in compilation errors.
type UnsafeWriterServiceServer interface {
	mustEmbedUnimplementedWriterServiceServer()
}

func RegisterWriterServiceServer(s grpc.ServiceRegistrar, srv WriterServiceServer) {
	s.RegisterService(&_WriterService_serviceDesc, srv)
}

func _WriterService_CreateProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProductReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).CreateProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/writerService.writerService/CreateProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).CreateProduct(ctx, req.(*CreateProductReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_UpdateProduct_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProductReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).UpdateProduct(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/writerService.writerService/UpdateProduct",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).UpdateProduct(ctx, req.(*UpdateProductReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _WriterService_GetProductById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProductByIdReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServiceServer).GetProductById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/writerService.writerService/GetProductById",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServiceServer).GetProductById(ctx, req.(*GetProductByIdReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _WriterService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "writerService.writerService",
	HandlerType: (*WriterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateProduct",
			Handler:    _WriterService_CreateProduct_Handler,
		},
		{
			MethodName: "UpdateProduct",
			Handler:    _WriterService_UpdateProduct_Handler,
		},
		{
			MethodName: "GetProductById",
			Handler:    _WriterService_GetProductById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "product_writer.proto",
}
