// Code generated by protoc-gen-go. DO NOT EDIT.
// source: vitepb/snapshot_block.proto

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

type SnapshotBlock struct {
	Hash                 []byte   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	PrevHash             []byte   `protobuf:"bytes,2,opt,name=prevHash,proto3" json:"prevHash,omitempty"`
	Height               uint64   `protobuf:"varint,3,opt,name=height,proto3" json:"height,omitempty"`
	PublicKey            []byte   `protobuf:"bytes,4,opt,name=publicKey,proto3" json:"publicKey,omitempty"`
	Signature            []byte   `protobuf:"bytes,5,opt,name=signature,proto3" json:"signature,omitempty"`
	Timestamp            int64    `protobuf:"varint,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Seed                 uint64   `protobuf:"varint,7,opt,name=Seed,proto3" json:"Seed,omitempty"`
	SeedHash             []byte   `protobuf:"bytes,8,opt,name=SeedHash,proto3" json:"SeedHash,omitempty"`
	SnapshotContent      []byte   `protobuf:"bytes,9,opt,name=snapshotContent,proto3" json:"snapshotContent,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SnapshotBlock) Reset()         { *m = SnapshotBlock{} }
func (m *SnapshotBlock) String() string { return proto.CompactTextString(m) }
func (*SnapshotBlock) ProtoMessage()    {}
func (*SnapshotBlock) Descriptor() ([]byte, []int) {
	return fileDescriptor_14ed8e66c18c4fa1, []int{0}
}

func (m *SnapshotBlock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SnapshotBlock.Unmarshal(m, b)
}
func (m *SnapshotBlock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SnapshotBlock.Marshal(b, m, deterministic)
}
func (m *SnapshotBlock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SnapshotBlock.Merge(m, src)
}
func (m *SnapshotBlock) XXX_Size() int {
	return xxx_messageInfo_SnapshotBlock.Size(m)
}
func (m *SnapshotBlock) XXX_DiscardUnknown() {
	xxx_messageInfo_SnapshotBlock.DiscardUnknown(m)
}

var xxx_messageInfo_SnapshotBlock proto.InternalMessageInfo

func (m *SnapshotBlock) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *SnapshotBlock) GetPrevHash() []byte {
	if m != nil {
		return m.PrevHash
	}
	return nil
}

func (m *SnapshotBlock) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *SnapshotBlock) GetPublicKey() []byte {
	if m != nil {
		return m.PublicKey
	}
	return nil
}

func (m *SnapshotBlock) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *SnapshotBlock) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *SnapshotBlock) GetSeed() uint64 {
	if m != nil {
		return m.Seed
	}
	return 0
}

func (m *SnapshotBlock) GetSeedHash() []byte {
	if m != nil {
		return m.SeedHash
	}
	return nil
}

func (m *SnapshotBlock) GetSnapshotContent() []byte {
	if m != nil {
		return m.SnapshotContent
	}
	return nil
}

func init() {
	proto.RegisterType((*SnapshotBlock)(nil), "vitepb.SnapshotBlock")
}

func init() { proto.RegisterFile("vitepb/snapshot_block.proto", fileDescriptor_14ed8e66c18c4fa1) }

var fileDescriptor_14ed8e66c18c4fa1 = []byte{
	// 221 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xcf, 0x4a, 0xc4, 0x30,
	0x10, 0xc6, 0xc9, 0x6e, 0x8d, 0xbb, 0x41, 0x11, 0x72, 0x90, 0xa0, 0x1e, 0x8a, 0xa7, 0x9c, 0xf4,
	0xe0, 0x1b, 0xe8, 0x45, 0xf0, 0xd6, 0x7d, 0x00, 0x49, 0xd6, 0x61, 0x13, 0xdc, 0x26, 0xa1, 0x99,
	0x16, 0x7c, 0x06, 0x5f, 0x5a, 0x32, 0xfd, 0x07, 0x9e, 0x3a, 0xdf, 0xef, 0x37, 0x25, 0x1f, 0x23,
	0xee, 0x07, 0x8f, 0x90, 0xec, 0x73, 0x0e, 0x26, 0x65, 0x17, 0xf1, 0xd3, 0x9e, 0xe3, 0xf1, 0xfb,
	0x29, 0x75, 0x11, 0xa3, 0xe4, 0xa3, 0x7c, 0xfc, 0xdd, 0x88, 0xeb, 0xc3, 0xb4, 0xf0, 0x5a, 0xbc,
	0x94, 0xa2, 0x72, 0x26, 0x3b, 0xc5, 0x6a, 0xa6, 0xaf, 0x1a, 0x9a, 0xe5, 0x9d, 0xd8, 0xa5, 0x0e,
	0x86, 0xf7, 0xc2, 0x37, 0xc4, 0x97, 0x2c, 0x6f, 0x05, 0x77, 0xe0, 0x4f, 0x0e, 0xd5, 0xb6, 0x66,
	0xba, 0x6a, 0xa6, 0x24, 0x1f, 0xc4, 0x3e, 0xf5, 0xf6, 0xec, 0x8f, 0x1f, 0xf0, 0xa3, 0x2a, 0xfa,
	0x69, 0x05, 0xc5, 0x66, 0x7f, 0x0a, 0x06, 0xfb, 0x0e, 0xd4, 0xc5, 0x68, 0x17, 0x50, 0x2c, 0xfa,
	0x16, 0x32, 0x9a, 0x36, 0x29, 0x5e, 0x33, 0xbd, 0x6d, 0x56, 0x50, 0x1a, 0x1e, 0x00, 0xbe, 0xd4,
	0x25, 0xbd, 0x47, 0x73, 0x69, 0x58, 0xbe, 0xd4, 0x70, 0x37, 0x36, 0x9c, 0xb3, 0xd4, 0xe2, 0x66,
	0xbe, 0xc1, 0x5b, 0x0c, 0x08, 0x01, 0xd5, 0x9e, 0x56, 0xfe, 0x63, 0xcb, 0xe9, 0x38, 0x2f, 0x7f,
	0x01, 0x00, 0x00, 0xff, 0xff, 0x4b, 0xeb, 0x42, 0xf0, 0x3b, 0x01, 0x00, 0x00,
}
