// Code generated by protoc-gen-go. DO NOT EDIT.
// source: annotations.proto

/*
Package google_api is a generated protocol buffer package.

It is generated from these files:
	annotations.proto
	http.proto

It has these top-level messages:
	Http
	HttpRule
	CustomHttpPattern
*/
package google_api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

var E_Http = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MethodOptions)(nil),
	ExtensionType: (*HttpRule)(nil),
	Field:         72295729,
	Name:          "google.api.http",
	Tag:           "bytes,72295729,opt,name=http",
	Filename:      "annotations.proto",
}

func init() {
	proto.RegisterExtension(E_Http)
}

func init() { proto.RegisterFile("annotations.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 161 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4c, 0xcc, 0xcb, 0xcb,
	0x2f, 0x49, 0x2c, 0xc9, 0xcc, 0xcf, 0x2b, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x4a,
	0xcf, 0xcf, 0x4f, 0xcf, 0x49, 0xd5, 0x4b, 0x2c, 0xc8, 0x94, 0xe2, 0xca, 0x28, 0x29, 0x29, 0x80,
	0x88, 0x4b, 0x29, 0x40, 0xc4, 0xf5, 0xc1, 0xbc, 0xa4, 0xd2, 0x34, 0xfd, 0x94, 0xd4, 0xe2, 0xe4,
	0xa2, 0xcc, 0x82, 0x92, 0xfc, 0x22, 0x88, 0x0a, 0x2b, 0x6f, 0x2e, 0x16, 0x90, 0x7a, 0x21, 0x39,
	0x3d, 0xa8, 0x11, 0x30, 0xa5, 0x7a, 0xbe, 0xa9, 0x25, 0x19, 0xf9, 0x29, 0xfe, 0x05, 0x60, 0x7b,
	0x24, 0x36, 0x9c, 0xda, 0xa3, 0xa4, 0xc0, 0xa8, 0xc1, 0x6d, 0x24, 0xa2, 0x87, 0xb0, 0x4b, 0xcf,
	0xa3, 0xa4, 0xa4, 0x20, 0xa8, 0x34, 0x27, 0x35, 0x08, 0x6c, 0x88, 0x93, 0x0a, 0x17, 0x5f, 0x72,
	0x7e, 0x2e, 0x92, 0x02, 0x27, 0x01, 0x47, 0x84, 0x5b, 0x03, 0x40, 0x26, 0x07, 0x30, 0x26, 0xb1,
	0x81, 0xad, 0x30, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x11, 0x08, 0x01, 0xa2, 0xc8, 0x00, 0x00,
	0x00,
}