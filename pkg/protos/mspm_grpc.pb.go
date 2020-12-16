// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protos

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// MspmClient is the client API for Mspm service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MspmClient interface {
	SetLabels(ctx context.Context, in *SetLabelRequest, opts ...grpc.CallOption) (*PackageInformation, error)
	GetPackageInformation(ctx context.Context, in *PackageInformationRequest, opts ...grpc.CallOption) (*PackageInformationResponse, error)
	UploadPackage(ctx context.Context, in *NewPackage, opts ...grpc.CallOption) (*PackageInformation, error)
	GetPackage(ctx context.Context, in *GetPackageRequest, opts ...grpc.CallOption) (*GetPackageResponse, error)
}

type mspmClient struct {
	cc grpc.ClientConnInterface
}

func NewMspmClient(cc grpc.ClientConnInterface) MspmClient {
	return &mspmClient{cc}
}

func (c *mspmClient) SetLabels(ctx context.Context, in *SetLabelRequest, opts ...grpc.CallOption) (*PackageInformation, error) {
	out := new(PackageInformation)
	err := c.cc.Invoke(ctx, "/mspm.Mspm/SetLabels", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mspmClient) GetPackageInformation(ctx context.Context, in *PackageInformationRequest, opts ...grpc.CallOption) (*PackageInformationResponse, error) {
	out := new(PackageInformationResponse)
	err := c.cc.Invoke(ctx, "/mspm.Mspm/GetPackageInformation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mspmClient) UploadPackage(ctx context.Context, in *NewPackage, opts ...grpc.CallOption) (*PackageInformation, error) {
	out := new(PackageInformation)
	err := c.cc.Invoke(ctx, "/mspm.Mspm/UploadPackage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mspmClient) GetPackage(ctx context.Context, in *GetPackageRequest, opts ...grpc.CallOption) (*GetPackageResponse, error) {
	out := new(GetPackageResponse)
	err := c.cc.Invoke(ctx, "/mspm.Mspm/GetPackage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MspmServer is the server API for Mspm service.
// All implementations must embed UnimplementedMspmServer
// for forward compatibility
type MspmServer interface {
	SetLabels(context.Context, *SetLabelRequest) (*PackageInformation, error)
	GetPackageInformation(context.Context, *PackageInformationRequest) (*PackageInformationResponse, error)
	UploadPackage(context.Context, *NewPackage) (*PackageInformation, error)
	GetPackage(context.Context, *GetPackageRequest) (*GetPackageResponse, error)
	mustEmbedUnimplementedMspmServer()
}

// UnimplementedMspmServer must be embedded to have forward compatible implementations.
type UnimplementedMspmServer struct {
}

func (UnimplementedMspmServer) SetLabels(context.Context, *SetLabelRequest) (*PackageInformation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetLabels not implemented")
}
func (UnimplementedMspmServer) GetPackageInformation(context.Context, *PackageInformationRequest) (*PackageInformationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPackageInformation not implemented")
}
func (UnimplementedMspmServer) UploadPackage(context.Context, *NewPackage) (*PackageInformation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadPackage not implemented")
}
func (UnimplementedMspmServer) GetPackage(context.Context, *GetPackageRequest) (*GetPackageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPackage not implemented")
}
func (UnimplementedMspmServer) mustEmbedUnimplementedMspmServer() {}

// UnsafeMspmServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MspmServer will
// result in compilation errors.
type UnsafeMspmServer interface {
	mustEmbedUnimplementedMspmServer()
}

func RegisterMspmServer(s grpc.ServiceRegistrar, srv MspmServer) {
	s.RegisterService(&_Mspm_serviceDesc, srv)
}

func _Mspm_SetLabels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetLabelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MspmServer).SetLabels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mspm.Mspm/SetLabels",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MspmServer).SetLabels(ctx, req.(*SetLabelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mspm_GetPackageInformation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PackageInformationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MspmServer).GetPackageInformation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mspm.Mspm/GetPackageInformation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MspmServer).GetPackageInformation(ctx, req.(*PackageInformationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mspm_UploadPackage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewPackage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MspmServer).UploadPackage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mspm.Mspm/UploadPackage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MspmServer).UploadPackage(ctx, req.(*NewPackage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mspm_GetPackage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPackageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MspmServer).GetPackage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mspm.Mspm/GetPackage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MspmServer).GetPackage(ctx, req.(*GetPackageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Mspm_serviceDesc = grpc.ServiceDesc{
	ServiceName: "mspm.Mspm",
	HandlerType: (*MspmServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetLabels",
			Handler:    _Mspm_SetLabels_Handler,
		},
		{
			MethodName: "GetPackageInformation",
			Handler:    _Mspm_GetPackageInformation_Handler,
		},
		{
			MethodName: "UploadPackage",
			Handler:    _Mspm_UploadPackage_Handler,
		},
		{
			MethodName: "GetPackage",
			Handler:    _Mspm_GetPackage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mspm.proto",
}
