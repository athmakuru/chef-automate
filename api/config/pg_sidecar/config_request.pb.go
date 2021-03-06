// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/config/pg_sidecar/config_request.proto

package pg_sidecar // import "github.com/chef/automate/api/config/pg_sidecar"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import shared "github.com/chef/automate/api/config/shared"
import _ "github.com/chef/automate/components/automate-grpc/protoc-gen-a2-config/api/a2conf"
import wrappers "github.com/golang/protobuf/ptypes/wrappers"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ConfigRequest struct {
	V1                   *ConfigRequest_V1 `protobuf:"bytes,1,opt,name=v1,proto3" json:"v1,omitempty" toml:"v1,omitempty" mapstructure:"v1,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-" toml:"-" mapstructure:"-,omitempty"`
	XXX_unrecognized     []byte            `json:"-" toml:"-" mapstructure:"-,omitempty"`
	XXX_sizecache        int32             `json:"-" toml:"-" mapstructure:"-,omitempty"`
}

func (m *ConfigRequest) Reset()         { *m = ConfigRequest{} }
func (m *ConfigRequest) String() string { return proto.CompactTextString(m) }
func (*ConfigRequest) ProtoMessage()    {}
func (*ConfigRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_request_19df0529236e4a9b, []int{0}
}
func (m *ConfigRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigRequest.Unmarshal(m, b)
}
func (m *ConfigRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigRequest.Marshal(b, m, deterministic)
}
func (dst *ConfigRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigRequest.Merge(dst, src)
}
func (m *ConfigRequest) XXX_Size() int {
	return xxx_messageInfo_ConfigRequest.Size(m)
}
func (m *ConfigRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigRequest proto.InternalMessageInfo

func (m *ConfigRequest) GetV1() *ConfigRequest_V1 {
	if m != nil {
		return m.V1
	}
	return nil
}

type ConfigRequest_V1 struct {
	Sys                  *ConfigRequest_V1_System  `protobuf:"bytes,1,opt,name=sys,proto3" json:"sys,omitempty" toml:"sys,omitempty" mapstructure:"sys,omitempty"`
	Svc                  *ConfigRequest_V1_Service `protobuf:"bytes,2,opt,name=svc,proto3" json:"svc,omitempty" toml:"svc,omitempty" mapstructure:"svc,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-" toml:"-" mapstructure:"-,omitempty"`
	XXX_unrecognized     []byte                    `json:"-" toml:"-" mapstructure:"-,omitempty"`
	XXX_sizecache        int32                     `json:"-" toml:"-" mapstructure:"-,omitempty"`
}

func (m *ConfigRequest_V1) Reset()         { *m = ConfigRequest_V1{} }
func (m *ConfigRequest_V1) String() string { return proto.CompactTextString(m) }
func (*ConfigRequest_V1) ProtoMessage()    {}
func (*ConfigRequest_V1) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_request_19df0529236e4a9b, []int{0, 0}
}
func (m *ConfigRequest_V1) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigRequest_V1.Unmarshal(m, b)
}
func (m *ConfigRequest_V1) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigRequest_V1.Marshal(b, m, deterministic)
}
func (dst *ConfigRequest_V1) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigRequest_V1.Merge(dst, src)
}
func (m *ConfigRequest_V1) XXX_Size() int {
	return xxx_messageInfo_ConfigRequest_V1.Size(m)
}
func (m *ConfigRequest_V1) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigRequest_V1.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigRequest_V1 proto.InternalMessageInfo

func (m *ConfigRequest_V1) GetSys() *ConfigRequest_V1_System {
	if m != nil {
		return m.Sys
	}
	return nil
}

func (m *ConfigRequest_V1) GetSvc() *ConfigRequest_V1_Service {
	if m != nil {
		return m.Svc
	}
	return nil
}

type ConfigRequest_V1_System struct {
	Mlsa                 *shared.Mlsa                     `protobuf:"bytes,1,opt,name=mlsa,proto3" json:"mlsa,omitempty" toml:"mlsa,omitempty" mapstructure:"mlsa,omitempty"`
	Tls                  *shared.TLSCredentials           `protobuf:"bytes,2,opt,name=tls,proto3" json:"tls,omitempty" toml:"tls,omitempty" mapstructure:"tls,omitempty"`
	Service              *ConfigRequest_V1_System_Service `protobuf:"bytes,3,opt,name=service,proto3" json:"service,omitempty" toml:"service,omitempty" mapstructure:"service,omitempty"`
	Log                  *ConfigRequest_V1_System_Log     `protobuf:"bytes,4,opt,name=log,proto3" json:"log,omitempty" toml:"log,omitempty" mapstructure:"log,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                         `json:"-" toml:"-" mapstructure:"-,omitempty"`
	XXX_unrecognized     []byte                           `json:"-" toml:"-" mapstructure:"-,omitempty"`
	XXX_sizecache        int32                            `json:"-" toml:"-" mapstructure:"-,omitempty"`
}

func (m *ConfigRequest_V1_System) Reset()         { *m = ConfigRequest_V1_System{} }
func (m *ConfigRequest_V1_System) String() string { return proto.CompactTextString(m) }
func (*ConfigRequest_V1_System) ProtoMessage()    {}
func (*ConfigRequest_V1_System) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_request_19df0529236e4a9b, []int{0, 0, 0}
}
func (m *ConfigRequest_V1_System) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigRequest_V1_System.Unmarshal(m, b)
}
func (m *ConfigRequest_V1_System) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigRequest_V1_System.Marshal(b, m, deterministic)
}
func (dst *ConfigRequest_V1_System) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigRequest_V1_System.Merge(dst, src)
}
func (m *ConfigRequest_V1_System) XXX_Size() int {
	return xxx_messageInfo_ConfigRequest_V1_System.Size(m)
}
func (m *ConfigRequest_V1_System) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigRequest_V1_System.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigRequest_V1_System proto.InternalMessageInfo

func (m *ConfigRequest_V1_System) GetMlsa() *shared.Mlsa {
	if m != nil {
		return m.Mlsa
	}
	return nil
}

func (m *ConfigRequest_V1_System) GetTls() *shared.TLSCredentials {
	if m != nil {
		return m.Tls
	}
	return nil
}

func (m *ConfigRequest_V1_System) GetService() *ConfigRequest_V1_System_Service {
	if m != nil {
		return m.Service
	}
	return nil
}

func (m *ConfigRequest_V1_System) GetLog() *ConfigRequest_V1_System_Log {
	if m != nil {
		return m.Log
	}
	return nil
}

type ConfigRequest_V1_System_Service struct {
	Host                 *wrappers.StringValue `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty" toml:"host,omitempty" mapstructure:"host,omitempty"`
	Port                 *wrappers.Int32Value  `protobuf:"bytes,2,opt,name=port,proto3" json:"port,omitempty" toml:"port,omitempty" mapstructure:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-" toml:"-" mapstructure:"-,omitempty"`
	XXX_unrecognized     []byte                `json:"-" toml:"-" mapstructure:"-,omitempty"`
	XXX_sizecache        int32                 `json:"-" toml:"-" mapstructure:"-,omitempty"`
}

func (m *ConfigRequest_V1_System_Service) Reset()         { *m = ConfigRequest_V1_System_Service{} }
func (m *ConfigRequest_V1_System_Service) String() string { return proto.CompactTextString(m) }
func (*ConfigRequest_V1_System_Service) ProtoMessage()    {}
func (*ConfigRequest_V1_System_Service) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_request_19df0529236e4a9b, []int{0, 0, 0, 0}
}
func (m *ConfigRequest_V1_System_Service) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigRequest_V1_System_Service.Unmarshal(m, b)
}
func (m *ConfigRequest_V1_System_Service) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigRequest_V1_System_Service.Marshal(b, m, deterministic)
}
func (dst *ConfigRequest_V1_System_Service) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigRequest_V1_System_Service.Merge(dst, src)
}
func (m *ConfigRequest_V1_System_Service) XXX_Size() int {
	return xxx_messageInfo_ConfigRequest_V1_System_Service.Size(m)
}
func (m *ConfigRequest_V1_System_Service) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigRequest_V1_System_Service.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigRequest_V1_System_Service proto.InternalMessageInfo

func (m *ConfigRequest_V1_System_Service) GetHost() *wrappers.StringValue {
	if m != nil {
		return m.Host
	}
	return nil
}

func (m *ConfigRequest_V1_System_Service) GetPort() *wrappers.Int32Value {
	if m != nil {
		return m.Port
	}
	return nil
}

type ConfigRequest_V1_System_Log struct {
	Level                *wrappers.StringValue `protobuf:"bytes,1,opt,name=level,proto3" json:"level,omitempty" toml:"level,omitempty" mapstructure:"level,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-" toml:"-" mapstructure:"-,omitempty"`
	XXX_unrecognized     []byte                `json:"-" toml:"-" mapstructure:"-,omitempty"`
	XXX_sizecache        int32                 `json:"-" toml:"-" mapstructure:"-,omitempty"`
}

func (m *ConfigRequest_V1_System_Log) Reset()         { *m = ConfigRequest_V1_System_Log{} }
func (m *ConfigRequest_V1_System_Log) String() string { return proto.CompactTextString(m) }
func (*ConfigRequest_V1_System_Log) ProtoMessage()    {}
func (*ConfigRequest_V1_System_Log) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_request_19df0529236e4a9b, []int{0, 0, 0, 1}
}
func (m *ConfigRequest_V1_System_Log) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigRequest_V1_System_Log.Unmarshal(m, b)
}
func (m *ConfigRequest_V1_System_Log) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigRequest_V1_System_Log.Marshal(b, m, deterministic)
}
func (dst *ConfigRequest_V1_System_Log) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigRequest_V1_System_Log.Merge(dst, src)
}
func (m *ConfigRequest_V1_System_Log) XXX_Size() int {
	return xxx_messageInfo_ConfigRequest_V1_System_Log.Size(m)
}
func (m *ConfigRequest_V1_System_Log) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigRequest_V1_System_Log.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigRequest_V1_System_Log proto.InternalMessageInfo

func (m *ConfigRequest_V1_System_Log) GetLevel() *wrappers.StringValue {
	if m != nil {
		return m.Level
	}
	return nil
}

type ConfigRequest_V1_Service struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-" toml:"-" mapstructure:"-,omitempty"`
	XXX_unrecognized     []byte   `json:"-" toml:"-" mapstructure:"-,omitempty"`
	XXX_sizecache        int32    `json:"-" toml:"-" mapstructure:"-,omitempty"`
}

func (m *ConfigRequest_V1_Service) Reset()         { *m = ConfigRequest_V1_Service{} }
func (m *ConfigRequest_V1_Service) String() string { return proto.CompactTextString(m) }
func (*ConfigRequest_V1_Service) ProtoMessage()    {}
func (*ConfigRequest_V1_Service) Descriptor() ([]byte, []int) {
	return fileDescriptor_config_request_19df0529236e4a9b, []int{0, 0, 1}
}
func (m *ConfigRequest_V1_Service) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConfigRequest_V1_Service.Unmarshal(m, b)
}
func (m *ConfigRequest_V1_Service) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConfigRequest_V1_Service.Marshal(b, m, deterministic)
}
func (dst *ConfigRequest_V1_Service) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConfigRequest_V1_Service.Merge(dst, src)
}
func (m *ConfigRequest_V1_Service) XXX_Size() int {
	return xxx_messageInfo_ConfigRequest_V1_Service.Size(m)
}
func (m *ConfigRequest_V1_Service) XXX_DiscardUnknown() {
	xxx_messageInfo_ConfigRequest_V1_Service.DiscardUnknown(m)
}

var xxx_messageInfo_ConfigRequest_V1_Service proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ConfigRequest)(nil), "chef.automate.infra.pg_sidecar.ConfigRequest")
	proto.RegisterType((*ConfigRequest_V1)(nil), "chef.automate.infra.pg_sidecar.ConfigRequest.V1")
	proto.RegisterType((*ConfigRequest_V1_System)(nil), "chef.automate.infra.pg_sidecar.ConfigRequest.V1.System")
	proto.RegisterType((*ConfigRequest_V1_System_Service)(nil), "chef.automate.infra.pg_sidecar.ConfigRequest.V1.System.Service")
	proto.RegisterType((*ConfigRequest_V1_System_Log)(nil), "chef.automate.infra.pg_sidecar.ConfigRequest.V1.System.Log")
	proto.RegisterType((*ConfigRequest_V1_Service)(nil), "chef.automate.infra.pg_sidecar.ConfigRequest.V1.Service")
}

func init() {
	proto.RegisterFile("api/config/pg_sidecar/config_request.proto", fileDescriptor_config_request_19df0529236e4a9b)
}

var fileDescriptor_config_request_19df0529236e4a9b = []byte{
	// 493 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0xdf, 0x6a, 0xd4, 0x40,
	0x14, 0xc6, 0xd9, 0x4d, 0xb6, 0xc5, 0x11, 0x41, 0x06, 0x91, 0x90, 0xca, 0x52, 0xbc, 0x92, 0x42,
	0x26, 0xbb, 0xa9, 0xe0, 0x5f, 0x50, 0xda, 0x0b, 0x69, 0xd9, 0x7a, 0x91, 0x95, 0x05, 0xbd, 0x29,
	0xb3, 0xd9, 0xb3, 0xb3, 0x81, 0xc9, 0xcc, 0x38, 0x33, 0x89, 0xf4, 0xde, 0x17, 0xf0, 0x95, 0xfa,
	0x34, 0xde, 0xd7, 0x7b, 0x25, 0x93, 0xd9, 0x55, 0x69, 0x29, 0x6d, 0x2f, 0x93, 0x73, 0xbe, 0xdf,
	0x39, 0xdf, 0xc7, 0x49, 0xd0, 0x1e, 0x55, 0x65, 0x5a, 0x48, 0xb1, 0x2c, 0x59, 0xaa, 0xd8, 0xa9,
	0x29, 0x17, 0x50, 0x50, 0xed, 0xdf, 0x9c, 0x6a, 0xf8, 0x5a, 0x83, 0xb1, 0x44, 0x69, 0x69, 0x25,
	0x1e, 0x16, 0x2b, 0x58, 0x12, 0x5a, 0x5b, 0x59, 0x51, 0x0b, 0xa4, 0x14, 0x4b, 0x4d, 0xc9, 0x5f,
	0x51, 0x3c, 0xfc, 0x87, 0x65, 0x56, 0x54, 0xc3, 0x22, 0x65, 0x5c, 0xce, 0x29, 0xef, 0xf4, 0xf1,
	0xce, 0xe5, 0xba, 0xe5, 0xc6, 0x17, 0x8f, 0x0b, 0x59, 0x29, 0x29, 0x40, 0x58, 0x93, 0xae, 0x47,
	0x24, 0x4c, 0xab, 0x22, 0x75, 0xf5, 0x22, 0x61, 0x20, 0x12, 0x9a, 0x25, 0x5e, 0xdf, 0xa2, 0x68,
	0xd6, 0x3e, 0xa4, 0x54, 0x08, 0x69, 0xa9, 0x2d, 0xa5, 0x58, 0xb3, 0x86, 0x4c, 0x4a, 0xc6, 0xa1,
	0x53, 0xce, 0xeb, 0x65, 0xfa, 0x4d, 0x53, 0xa5, 0x40, 0xfb, 0xfa, 0xd3, 0xdf, 0x03, 0xf4, 0xe0,
	0xd0, 0x71, 0xf2, 0xce, 0x20, 0x7e, 0x8f, 0xfa, 0xcd, 0x38, 0xea, 0xed, 0xf6, 0x9e, 0xdd, 0xcf,
	0x46, 0xe4, 0x7a, 0x9f, 0xe4, 0x3f, 0x29, 0x99, 0x8d, 0xf3, 0x7e, 0x33, 0x8e, 0x7f, 0x0c, 0x50,
	0x7f, 0x36, 0xc6, 0x47, 0x28, 0x30, 0x67, 0xc6, 0x93, 0x5e, 0xdc, 0x96, 0x44, 0xa6, 0x67, 0xc6,
	0x42, 0x95, 0xb7, 0x0c, 0x7c, 0x8c, 0x02, 0xd3, 0x14, 0x51, 0xdf, 0xa1, 0x5e, 0xde, 0x1e, 0x05,
	0xba, 0x29, 0x0b, 0xc8, 0x5b, 0x48, 0xfc, 0x33, 0x40, 0x5b, 0x1d, 0x1b, 0x3f, 0x47, 0x61, 0xc5,
	0x0d, 0xf5, 0x2b, 0xee, 0x5e, 0xc9, 0xed, 0x42, 0x26, 0x27, 0xdc, 0xd0, 0xdc, 0x75, 0xe3, 0xb7,
	0x28, 0xb0, 0xdc, 0xf8, 0x65, 0xf6, 0xae, 0x13, 0x7d, 0x9a, 0x4c, 0x0f, 0x35, 0x2c, 0x40, 0xd8,
	0x92, 0x72, 0x93, 0xb7, 0x32, 0xfc, 0x19, 0x6d, 0x9b, 0x6e, 0x9d, 0x28, 0x70, 0x84, 0x77, 0x77,
	0x4c, 0x66, 0xe3, 0x6a, 0xcd, 0xc3, 0x27, 0x28, 0xe0, 0x92, 0x45, 0xa1, 0xc3, 0xbe, 0xb9, 0x2b,
	0x76, 0x22, 0x59, 0xde, 0x72, 0xe2, 0xef, 0x3d, 0xb4, 0xed, 0x67, 0xe0, 0x11, 0x0a, 0x57, 0xd2,
	0x58, 0x9f, 0xd4, 0x13, 0xd2, 0x5d, 0x15, 0x59, 0x5f, 0x15, 0x99, 0x5a, 0x5d, 0x0a, 0x36, 0xa3,
	0xbc, 0x86, 0xdc, 0x75, 0xe2, 0x0f, 0x28, 0x54, 0x52, 0x5b, 0x1f, 0xd3, 0xce, 0x25, 0xc5, 0x91,
	0xb0, 0xfb, 0x99, 0x13, 0x1c, 0x3c, 0x3e, 0xbf, 0x88, 0xf0, 0x26, 0x96, 0x87, 0xbf, 0x3e, 0xc6,
	0x61, 0x7b, 0xed, 0xb9, 0x03, 0xc4, 0xaf, 0x50, 0x30, 0x91, 0x0c, 0x67, 0x68, 0xc0, 0xa1, 0x01,
	0x7e, 0xa3, 0x15, 0xba, 0xd6, 0xf8, 0xde, 0xc6, 0xc0, 0xeb, 0xe8, 0xfc, 0x22, 0x7a, 0x84, 0xb0,
	0x62, 0x89, 0xf7, 0x9f, 0xf8, 0x71, 0x07, 0xa3, 0x2f, 0x84, 0x95, 0x76, 0x55, 0xcf, 0x49, 0x21,
	0xab, 0xb4, 0x0d, 0x6d, 0xf3, 0xd1, 0xa5, 0x57, 0xfe, 0x11, 0xe6, 0x5b, 0x6e, 0xe6, 0xfe, 0x9f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x26, 0xf2, 0xa2, 0xae, 0x31, 0x04, 0x00, 0x00,
}
