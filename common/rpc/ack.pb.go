// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ack.proto

package rpc

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type AckType int32

const (
	AckType_AT_UNKNOWN      AckType = 0
	AckType_AT_NOTIFICATION AckType = 1
	AckType_AT_MESSAGE      AckType = 2
	AckType_AT_ACT          AckType = 3
)

var AckType_name = map[int32]string{
	0: "AT_UNKNOWN",
	1: "AT_NOTIFICATION",
	2: "AT_MESSAGE",
	3: "AT_ACT",
}

var AckType_value = map[string]int32{
	"AT_UNKNOWN":      0,
	"AT_NOTIFICATION": 1,
	"AT_MESSAGE":      2,
	"AT_ACT":          3,
}

func (x AckType) String() string {
	return proto.EnumName(AckType_name, int32(x))
}

func (AckType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_29efde0d93e5101c, []int{0}
}

type Ack struct {
	ObjId                int64    `protobuf:"varint,1,opt,name=objId,proto3" json:"objId,omitempty"`
	Type                 AckType  `protobuf:"varint,2,opt,name=type,proto3,enum=rpc.AckType" json:"type,omitempty"`
	RequestId            int64    `protobuf:"varint,3,opt,name=requestId,proto3" json:"requestId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ack) Reset()         { *m = Ack{} }
func (m *Ack) String() string { return proto.CompactTextString(m) }
func (*Ack) ProtoMessage()    {}
func (*Ack) Descriptor() ([]byte, []int) {
	return fileDescriptor_29efde0d93e5101c, []int{0}
}

func (m *Ack) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ack.Unmarshal(m, b)
}
func (m *Ack) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ack.Marshal(b, m, deterministic)
}
func (m *Ack) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ack.Merge(m, src)
}
func (m *Ack) XXX_Size() int {
	return xxx_messageInfo_Ack.Size(m)
}
func (m *Ack) XXX_DiscardUnknown() {
	xxx_messageInfo_Ack.DiscardUnknown(m)
}

var xxx_messageInfo_Ack proto.InternalMessageInfo

func (m *Ack) GetObjId() int64 {
	if m != nil {
		return m.ObjId
	}
	return 0
}

func (m *Ack) GetType() AckType {
	if m != nil {
		return m.Type
	}
	return AckType_AT_UNKNOWN
}

func (m *Ack) GetRequestId() int64 {
	if m != nil {
		return m.RequestId
	}
	return 0
}

func init() {
	proto.RegisterEnum("rpc.AckType", AckType_name, AckType_value)
	proto.RegisterType((*Ack)(nil), "rpc.Ack")
}

func init() { proto.RegisterFile("ack.proto", fileDescriptor_29efde0d93e5101c) }

var fileDescriptor_29efde0d93e5101c = []byte{
	// 191 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4c, 0x4c, 0xce, 0xd6,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2e, 0x2a, 0x48, 0x56, 0x8a, 0xe6, 0x62, 0x76, 0x4c,
	0xce, 0x16, 0x12, 0xe1, 0x62, 0xcd, 0x4f, 0xca, 0xf2, 0x4c, 0x91, 0x60, 0x54, 0x60, 0xd4, 0x60,
	0x0e, 0x82, 0x70, 0x84, 0x14, 0xb8, 0x58, 0x4a, 0x2a, 0x0b, 0x52, 0x25, 0x98, 0x14, 0x18, 0x35,
	0xf8, 0x8c, 0x78, 0xf4, 0x8a, 0x0a, 0x92, 0xf5, 0x1c, 0x93, 0xb3, 0x43, 0x2a, 0x0b, 0x52, 0x83,
	0xc0, 0x32, 0x42, 0x32, 0x5c, 0x9c, 0x45, 0xa9, 0x85, 0xa5, 0xa9, 0xc5, 0x25, 0x9e, 0x29, 0x12,
	0xcc, 0x60, 0xbd, 0x08, 0x01, 0x2d, 0x2f, 0x2e, 0x76, 0xa8, 0x72, 0x21, 0x3e, 0x2e, 0x2e, 0xc7,
	0x90, 0xf8, 0x50, 0x3f, 0x6f, 0x3f, 0xff, 0x70, 0x3f, 0x01, 0x06, 0x21, 0x61, 0x2e, 0x7e, 0xc7,
	0x90, 0x78, 0x3f, 0xff, 0x10, 0x4f, 0x37, 0x4f, 0x67, 0xc7, 0x10, 0x4f, 0x7f, 0x3f, 0x01, 0x46,
	0xa8, 0x22, 0x5f, 0xd7, 0xe0, 0x60, 0x47, 0x77, 0x57, 0x01, 0x26, 0x21, 0x2e, 0x2e, 0x36, 0xc7,
	0x90, 0x78, 0x47, 0xe7, 0x10, 0x01, 0x66, 0x27, 0xd6, 0x28, 0x90, 0x7b, 0x93, 0xd8, 0xc0, 0x6e,
	0x37, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xdf, 0x38, 0x03, 0x3f, 0xc8, 0x00, 0x00, 0x00,
}
