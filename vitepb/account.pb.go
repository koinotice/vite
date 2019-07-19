// Code generated by protoc-gen-go. DO NOT EDIT.
// source: vitepb/account.proto

package vitepb

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

type Account struct {
	AccountId            uint64   `protobuf:"varint,1,opt,name=accountId,proto3" json:"accountId,omitempty"`
	PublicKey            []byte   `protobuf:"bytes,2,opt,name=publicKey,proto3" json:"publicKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Account) Reset()         { *m = Account{} }
func (m *Account) String() string { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()    {}
func (*Account) Descriptor() ([]byte, []int) {
	return fileDescriptor_d4b28117fc611ac5, []int{0}
}

func (m *Account) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Account.Unmarshal(m, b)
}
func (m *Account) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Account.Marshal(b, m, deterministic)
}
func (m *Account) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Account.Merge(m, src)
}
func (m *Account) XXX_Size() int {
	return xxx_messageInfo_Account.Size(m)
}
func (m *Account) XXX_DiscardUnknown() {
	xxx_messageInfo_Account.DiscardUnknown(m)
}

var xxx_messageInfo_Account proto.InternalMessageInfo

func (m *Account) GetAccountId() uint64 {
	if m != nil {
		return m.AccountId
	}
	return 0
}

func (m *Account) GetPublicKey() []byte {
	if m != nil {
		return m.PublicKey
	}
	return nil
}

func init() {
	proto.RegisterType((*Account)(nil), "vitepb.Account")
}

func init() { proto.RegisterFile("vitepb/account.proto", fileDescriptor_d4b28117fc611ac5) }

var fileDescriptor_d4b28117fc611ac5 = []byte{
	// 99 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x29, 0xcb, 0x2c, 0x49,
	0x2d, 0x48, 0xd2, 0x4f, 0x4c, 0x4e, 0xce, 0x2f, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9,
	0x17, 0x62, 0x83, 0x88, 0x2a, 0xb9, 0x72, 0xb1, 0x3b, 0x42, 0x24, 0x84, 0x64, 0xb8, 0x38, 0xa1,
	0x6a, 0x3c, 0x53, 0x24, 0x18, 0x15, 0x18, 0x35, 0x58, 0x82, 0x10, 0x02, 0x20, 0xd9, 0x82, 0xd2,
	0xa4, 0x9c, 0xcc, 0x64, 0xef, 0xd4, 0x4a, 0x09, 0x26, 0x05, 0x46, 0x0d, 0x9e, 0x20, 0x84, 0x40,
	0x12, 0x1b, 0xd8, 0x54, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x01, 0x1f, 0xfd, 0x4a, 0x6d,
	0x00, 0x00, 0x00,
}
