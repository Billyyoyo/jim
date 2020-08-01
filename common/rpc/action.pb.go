// Code generated by protoc-gen-go. DO NOT EDIT.
// source: action.proto

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

type ActType int32

const (
	ActType_ACT_UNKNOWN  ActType = 0
	ActType_ACT_JOIN     ActType = 1
	ActType_ACT_QUIT     ActType = 2
	ActType_ACT_WITHDRAW ActType = 3
	ActType_ACT_SYNC     ActType = 4
	ActType_ACT_RENAME   ActType = 5
	ActType_ACT_CREATE   ActType = 6
)

var ActType_name = map[int32]string{
	0: "ACT_UNKNOWN",
	1: "ACT_JOIN",
	2: "ACT_QUIT",
	3: "ACT_WITHDRAW",
	4: "ACT_SYNC",
	5: "ACT_RENAME",
	6: "ACT_CREATE",
}

var ActType_value = map[string]int32{
	"ACT_UNKNOWN":  0,
	"ACT_JOIN":     1,
	"ACT_QUIT":     2,
	"ACT_WITHDRAW": 3,
	"ACT_SYNC":     4,
	"ACT_RENAME":   5,
	"ACT_CREATE":   6,
}

func (x ActType) String() string {
	return proto.EnumName(ActType_name, int32(x))
}

func (ActType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_59885c909ad4dfd3, []int{0}
}

type Action struct {
	UserId               int64    `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	RemoteAddr           string   `protobuf:"bytes,2,opt,name=remoteAddr,proto3" json:"remoteAddr,omitempty"`
	RequestId            int64    `protobuf:"varint,3,opt,name=requestId,proto3" json:"requestId,omitempty"`
	Time                 int64    `protobuf:"varint,4,opt,name=time,proto3" json:"time,omitempty"`
	Type                 ActType  `protobuf:"varint,5,opt,name=type,proto3,enum=rpc.ActType" json:"type,omitempty"`
	Data                 []byte   `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Action) Reset()         { *m = Action{} }
func (m *Action) String() string { return proto.CompactTextString(m) }
func (*Action) ProtoMessage()    {}
func (*Action) Descriptor() ([]byte, []int) {
	return fileDescriptor_59885c909ad4dfd3, []int{0}
}

func (m *Action) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Action.Unmarshal(m, b)
}
func (m *Action) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Action.Marshal(b, m, deterministic)
}
func (m *Action) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Action.Merge(m, src)
}
func (m *Action) XXX_Size() int {
	return xxx_messageInfo_Action.Size(m)
}
func (m *Action) XXX_DiscardUnknown() {
	xxx_messageInfo_Action.DiscardUnknown(m)
}

var xxx_messageInfo_Action proto.InternalMessageInfo

func (m *Action) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *Action) GetRemoteAddr() string {
	if m != nil {
		return m.RemoteAddr
	}
	return ""
}

func (m *Action) GetRequestId() int64 {
	if m != nil {
		return m.RequestId
	}
	return 0
}

func (m *Action) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *Action) GetType() ActType {
	if m != nil {
		return m.Type
	}
	return ActType_ACT_UNKNOWN
}

func (m *Action) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type CreateSessionAction struct {
	SessionId            int64       `protobuf:"varint,1,opt,name=sessionId,proto3" json:"sessionId,omitempty"`
	OwnerId              int64       `protobuf:"varint,2,opt,name=ownerId,proto3" json:"ownerId,omitempty"`
	Name                 string      `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	UserIds              []int64     `protobuf:"varint,4,rep,packed,name=userIds,proto3" json:"userIds,omitempty"`
	Type                 SessionType `protobuf:"varint,5,opt,name=type,proto3,enum=rpc.SessionType" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *CreateSessionAction) Reset()         { *m = CreateSessionAction{} }
func (m *CreateSessionAction) String() string { return proto.CompactTextString(m) }
func (*CreateSessionAction) ProtoMessage()    {}
func (*CreateSessionAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_59885c909ad4dfd3, []int{1}
}

func (m *CreateSessionAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateSessionAction.Unmarshal(m, b)
}
func (m *CreateSessionAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateSessionAction.Marshal(b, m, deterministic)
}
func (m *CreateSessionAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateSessionAction.Merge(m, src)
}
func (m *CreateSessionAction) XXX_Size() int {
	return xxx_messageInfo_CreateSessionAction.Size(m)
}
func (m *CreateSessionAction) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateSessionAction.DiscardUnknown(m)
}

var xxx_messageInfo_CreateSessionAction proto.InternalMessageInfo

func (m *CreateSessionAction) GetSessionId() int64 {
	if m != nil {
		return m.SessionId
	}
	return 0
}

func (m *CreateSessionAction) GetOwnerId() int64 {
	if m != nil {
		return m.OwnerId
	}
	return 0
}

func (m *CreateSessionAction) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateSessionAction) GetUserIds() []int64 {
	if m != nil {
		return m.UserIds
	}
	return nil
}

func (m *CreateSessionAction) GetType() SessionType {
	if m != nil {
		return m.Type
	}
	return SessionType_SESSION_UNKNOWN
}

type JoinSessionAction struct {
	SessionId            int64    `protobuf:"varint,1,opt,name=sessionId,proto3" json:"sessionId,omitempty"`
	User                 *User    `protobuf:"bytes,2,opt,name=user,proto3" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JoinSessionAction) Reset()         { *m = JoinSessionAction{} }
func (m *JoinSessionAction) String() string { return proto.CompactTextString(m) }
func (*JoinSessionAction) ProtoMessage()    {}
func (*JoinSessionAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_59885c909ad4dfd3, []int{2}
}

func (m *JoinSessionAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinSessionAction.Unmarshal(m, b)
}
func (m *JoinSessionAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinSessionAction.Marshal(b, m, deterministic)
}
func (m *JoinSessionAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinSessionAction.Merge(m, src)
}
func (m *JoinSessionAction) XXX_Size() int {
	return xxx_messageInfo_JoinSessionAction.Size(m)
}
func (m *JoinSessionAction) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinSessionAction.DiscardUnknown(m)
}

var xxx_messageInfo_JoinSessionAction proto.InternalMessageInfo

func (m *JoinSessionAction) GetSessionId() int64 {
	if m != nil {
		return m.SessionId
	}
	return 0
}

func (m *JoinSessionAction) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type QuitSessionAction struct {
	SessionId            int64    `protobuf:"varint,1,opt,name=sessionId,proto3" json:"sessionId,omitempty"`
	UserId               int64    `protobuf:"varint,2,opt,name=userId,proto3" json:"userId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QuitSessionAction) Reset()         { *m = QuitSessionAction{} }
func (m *QuitSessionAction) String() string { return proto.CompactTextString(m) }
func (*QuitSessionAction) ProtoMessage()    {}
func (*QuitSessionAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_59885c909ad4dfd3, []int{3}
}

func (m *QuitSessionAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QuitSessionAction.Unmarshal(m, b)
}
func (m *QuitSessionAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QuitSessionAction.Marshal(b, m, deterministic)
}
func (m *QuitSessionAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QuitSessionAction.Merge(m, src)
}
func (m *QuitSessionAction) XXX_Size() int {
	return xxx_messageInfo_QuitSessionAction.Size(m)
}
func (m *QuitSessionAction) XXX_DiscardUnknown() {
	xxx_messageInfo_QuitSessionAction.DiscardUnknown(m)
}

var xxx_messageInfo_QuitSessionAction proto.InternalMessageInfo

func (m *QuitSessionAction) GetSessionId() int64 {
	if m != nil {
		return m.SessionId
	}
	return 0
}

func (m *QuitSessionAction) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

type WithdrawMessageAction struct {
	UserId               int64    `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	SessionId            int64    `protobuf:"varint,2,opt,name=sessionId,proto3" json:"sessionId,omitempty"`
	SendNo               int64    `protobuf:"varint,3,opt,name=sendNo,proto3" json:"sendNo,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WithdrawMessageAction) Reset()         { *m = WithdrawMessageAction{} }
func (m *WithdrawMessageAction) String() string { return proto.CompactTextString(m) }
func (*WithdrawMessageAction) ProtoMessage()    {}
func (*WithdrawMessageAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_59885c909ad4dfd3, []int{4}
}

func (m *WithdrawMessageAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WithdrawMessageAction.Unmarshal(m, b)
}
func (m *WithdrawMessageAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WithdrawMessageAction.Marshal(b, m, deterministic)
}
func (m *WithdrawMessageAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WithdrawMessageAction.Merge(m, src)
}
func (m *WithdrawMessageAction) XXX_Size() int {
	return xxx_messageInfo_WithdrawMessageAction.Size(m)
}
func (m *WithdrawMessageAction) XXX_DiscardUnknown() {
	xxx_messageInfo_WithdrawMessageAction.DiscardUnknown(m)
}

var xxx_messageInfo_WithdrawMessageAction proto.InternalMessageInfo

func (m *WithdrawMessageAction) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *WithdrawMessageAction) GetSessionId() int64 {
	if m != nil {
		return m.SessionId
	}
	return 0
}

func (m *WithdrawMessageAction) GetSendNo() int64 {
	if m != nil {
		return m.SendNo
	}
	return 0
}

type SyncMessageAction struct {
	UserId               int64      `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Cond                 string     `protobuf:"bytes,2,opt,name=cond,proto3" json:"cond,omitempty"`
	Messages             []*Message `protobuf:"bytes,3,rep,name=messages,proto3" json:"messages,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *SyncMessageAction) Reset()         { *m = SyncMessageAction{} }
func (m *SyncMessageAction) String() string { return proto.CompactTextString(m) }
func (*SyncMessageAction) ProtoMessage()    {}
func (*SyncMessageAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_59885c909ad4dfd3, []int{5}
}

func (m *SyncMessageAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncMessageAction.Unmarshal(m, b)
}
func (m *SyncMessageAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncMessageAction.Marshal(b, m, deterministic)
}
func (m *SyncMessageAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncMessageAction.Merge(m, src)
}
func (m *SyncMessageAction) XXX_Size() int {
	return xxx_messageInfo_SyncMessageAction.Size(m)
}
func (m *SyncMessageAction) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncMessageAction.DiscardUnknown(m)
}

var xxx_messageInfo_SyncMessageAction proto.InternalMessageInfo

func (m *SyncMessageAction) GetUserId() int64 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *SyncMessageAction) GetCond() string {
	if m != nil {
		return m.Cond
	}
	return ""
}

func (m *SyncMessageAction) GetMessages() []*Message {
	if m != nil {
		return m.Messages
	}
	return nil
}

type RenameSessionAction struct {
	SessionId            int64    `protobuf:"varint,1,opt,name=sessionId,proto3" json:"sessionId,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RenameSessionAction) Reset()         { *m = RenameSessionAction{} }
func (m *RenameSessionAction) String() string { return proto.CompactTextString(m) }
func (*RenameSessionAction) ProtoMessage()    {}
func (*RenameSessionAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_59885c909ad4dfd3, []int{6}
}

func (m *RenameSessionAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RenameSessionAction.Unmarshal(m, b)
}
func (m *RenameSessionAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RenameSessionAction.Marshal(b, m, deterministic)
}
func (m *RenameSessionAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RenameSessionAction.Merge(m, src)
}
func (m *RenameSessionAction) XXX_Size() int {
	return xxx_messageInfo_RenameSessionAction.Size(m)
}
func (m *RenameSessionAction) XXX_DiscardUnknown() {
	xxx_messageInfo_RenameSessionAction.DiscardUnknown(m)
}

var xxx_messageInfo_RenameSessionAction proto.InternalMessageInfo

func (m *RenameSessionAction) GetSessionId() int64 {
	if m != nil {
		return m.SessionId
	}
	return 0
}

func (m *RenameSessionAction) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterEnum("rpc.ActType", ActType_name, ActType_value)
	proto.RegisterType((*Action)(nil), "rpc.Action")
	proto.RegisterType((*CreateSessionAction)(nil), "rpc.CreateSessionAction")
	proto.RegisterType((*JoinSessionAction)(nil), "rpc.JoinSessionAction")
	proto.RegisterType((*QuitSessionAction)(nil), "rpc.QuitSessionAction")
	proto.RegisterType((*WithdrawMessageAction)(nil), "rpc.WithdrawMessageAction")
	proto.RegisterType((*SyncMessageAction)(nil), "rpc.SyncMessageAction")
	proto.RegisterType((*RenameSessionAction)(nil), "rpc.RenameSessionAction")
}

func init() { proto.RegisterFile("action.proto", fileDescriptor_59885c909ad4dfd3) }

var fileDescriptor_59885c909ad4dfd3 = []byte{
	// 472 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x53, 0x5d, 0x6f, 0x94, 0x40,
	0x14, 0x95, 0x8f, 0xa5, 0xe5, 0x2e, 0xa9, 0xec, 0x34, 0x1a, 0xd2, 0x54, 0x43, 0x88, 0x0f, 0xc4,
	0x87, 0x7d, 0xa8, 0xbf, 0x00, 0xd7, 0x8d, 0x52, 0x53, 0x6a, 0x67, 0xd9, 0x6c, 0xf4, 0xc5, 0x20,
	0x4c, 0x2a, 0x89, 0x30, 0x38, 0x33, 0x6b, 0xb3, 0x3f, 0xc7, 0x37, 0x7f, 0xa6, 0x99, 0x61, 0x28,
	0xdd, 0xc4, 0xc4, 0xee, 0xdb, 0x3d, 0xe7, 0x7e, 0x1d, 0x38, 0x77, 0xc0, 0x2b, 0x4a, 0x51, 0xd3,
	0x76, 0xde, 0x31, 0x2a, 0x28, 0xb2, 0x58, 0x57, 0x9e, 0x4d, 0x1b, 0x5a, 0x91, 0x1f, 0x3d, 0x73,
	0xe6, 0x36, 0xfc, 0xb6, 0x0f, 0xa3, 0x3f, 0x06, 0x38, 0x89, 0xaa, 0x46, 0xcf, 0xc1, 0xd9, 0x72,
	0xc2, 0xd2, 0x2a, 0x30, 0x42, 0x23, 0xb6, 0xb0, 0x46, 0xe8, 0x25, 0x00, 0x23, 0x0d, 0x15, 0x24,
	0xa9, 0x2a, 0x16, 0x98, 0xa1, 0x11, 0xbb, 0xf8, 0x01, 0x83, 0xce, 0xc1, 0x65, 0xe4, 0xe7, 0x96,
	0x70, 0x91, 0x56, 0x81, 0xa5, 0x5a, 0x47, 0x02, 0x21, 0xb0, 0x45, 0xdd, 0x90, 0xc0, 0x56, 0x09,
	0x15, 0xa3, 0x10, 0x6c, 0xb1, 0xeb, 0x48, 0x30, 0x09, 0x8d, 0xf8, 0xe4, 0xc2, 0x9b, 0xb3, 0xae,
	0x9c, 0x27, 0xa5, 0xc8, 0x77, 0x1d, 0xc1, 0x2a, 0x23, 0xbb, 0xaa, 0x42, 0x14, 0x81, 0x13, 0x1a,
	0xb1, 0x87, 0x55, 0x1c, 0xfd, 0x36, 0xe0, 0x74, 0xc1, 0x48, 0x21, 0xc8, 0x8a, 0x70, 0x5e, 0xd3,
	0x56, 0xeb, 0x3e, 0x07, 0x97, 0xf7, 0xc4, 0xbd, 0xf4, 0x91, 0x40, 0x01, 0x1c, 0xd1, 0xbb, 0x56,
	0x7d, 0x96, 0xa9, 0x72, 0x03, 0x94, 0x3b, 0xda, 0xa2, 0x21, 0x4a, 0xb2, 0x8b, 0x55, 0x2c, 0xab,
	0xfb, 0xaf, 0xe6, 0x81, 0x1d, 0x5a, 0xb2, 0x5a, 0x43, 0xf4, 0x6a, 0x4f, 0xb3, 0xaf, 0x34, 0x6b,
	0x1d, 0xa3, 0xee, 0xe8, 0x13, 0xcc, 0x2e, 0x69, 0xdd, 0x1e, 0x22, 0xf0, 0x05, 0xd8, 0x72, 0x87,
	0x52, 0x37, 0xbd, 0x70, 0xd5, 0xe0, 0x35, 0x27, 0x0c, 0x2b, 0x3a, 0x4a, 0x61, 0x76, 0xb3, 0xad,
	0xc5, 0x21, 0x13, 0x47, 0x23, 0xcd, 0x87, 0x46, 0x46, 0x04, 0x9e, 0x6d, 0x6a, 0xf1, 0xbd, 0x62,
	0xc5, 0xdd, 0x15, 0xe1, 0xbc, 0xb8, 0x25, 0xff, 0x71, 0x7e, 0x6f, 0x8d, 0xf9, 0x8f, 0x35, 0x9c,
	0xb4, 0x55, 0x46, 0xb5, 0xe9, 0x1a, 0x45, 0x35, 0xcc, 0x56, 0xbb, 0xb6, 0x7c, 0xdc, 0x0a, 0x04,
	0x76, 0x49, 0xdb, 0x4a, 0x9f, 0x95, 0x8a, 0x51, 0x0c, 0xc7, 0x4d, 0xdf, 0xcc, 0x03, 0x2b, 0xb4,
	0xe2, 0xa9, 0x3e, 0x11, 0x3d, 0x11, 0xdf, 0x67, 0xa3, 0xf7, 0x70, 0x8a, 0x89, 0x34, 0xee, 0x90,
	0xdf, 0x33, 0xf8, 0x6e, 0x8e, 0xbe, 0xbf, 0xfe, 0x05, 0x47, 0xfa, 0x00, 0xd1, 0x53, 0x98, 0x26,
	0x8b, 0xfc, 0xeb, 0x3a, 0xfb, 0x98, 0x5d, 0x6f, 0x32, 0xff, 0x09, 0xf2, 0xe0, 0x58, 0x12, 0x97,
	0xd7, 0x69, 0xe6, 0x1b, 0x03, 0xba, 0x59, 0xa7, 0xb9, 0x6f, 0x22, 0x1f, 0x3c, 0x89, 0x36, 0x69,
	0xfe, 0xe1, 0x1d, 0x4e, 0x36, 0xbe, 0x35, 0xe4, 0x57, 0x9f, 0xb3, 0x85, 0x6f, 0xa3, 0x13, 0x00,
	0x89, 0xf0, 0x32, 0x4b, 0xae, 0x96, 0xfe, 0x64, 0xc0, 0x0b, 0xbc, 0x4c, 0xf2, 0xa5, 0xef, 0xbc,
	0x9d, 0x7c, 0x91, 0xaf, 0xf3, 0x9b, 0xa3, 0x1e, 0xe3, 0x9b, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff,
	0x0d, 0xa1, 0xd7, 0x0b, 0xb9, 0x03, 0x00, 0x00,
}
