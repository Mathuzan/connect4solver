// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.6.1
// source: cache_old.proto

package proto_old

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Ending int32

const (
	Ending_PLAYER_A Ending = 0
	Ending_PLAYER_B Ending = 1
	Ending_TIE      Ending = 2
)

// Enum value maps for Ending.
var (
	Ending_name = map[int32]string{
		0: "PLAYER_A",
		1: "PLAYER_B",
		2: "TIE",
	}
	Ending_value = map[string]int32{
		"PLAYER_A": 0,
		"PLAYER_B": 1,
		"TIE":      2,
	}
)

func (x Ending) Enum() *Ending {
	p := new(Ending)
	*p = x
	return p
}

func (x Ending) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Ending) Descriptor() protoreflect.EnumDescriptor {
	return file_cache_old_proto_enumTypes[0].Descriptor()
}

func (Ending) Type() protoreflect.EnumType {
	return &file_cache_old_proto_enumTypes[0]
}

func (x Ending) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Ending.Descriptor instead.
func (Ending) EnumDescriptor() ([]byte, []int) {
	return file_cache_old_proto_rawDescGZIP(), []int{0}
}

type DepthCaches struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DepthCaches []*DepthCache `protobuf:"bytes,1,rep,name=depthCaches,proto3" json:"depthCaches,omitempty"`
}

func (x *DepthCaches) Reset() {
	*x = DepthCaches{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_old_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DepthCaches) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DepthCaches) ProtoMessage() {}

func (x *DepthCaches) ProtoReflect() protoreflect.Message {
	mi := &file_cache_old_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DepthCaches.ProtoReflect.Descriptor instead.
func (*DepthCaches) Descriptor() ([]byte, []int) {
	return file_cache_old_proto_rawDescGZIP(), []int{0}
}

func (x *DepthCaches) GetDepthCaches() []*DepthCache {
	if x != nil {
		return x.DepthCaches
	}
	return nil
}

type DepthCache struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entries map[uint64]Ending `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3,enum=proto_old.Ending"`
}

func (x *DepthCache) Reset() {
	*x = DepthCache{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cache_old_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DepthCache) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DepthCache) ProtoMessage() {}

func (x *DepthCache) ProtoReflect() protoreflect.Message {
	mi := &file_cache_old_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DepthCache.ProtoReflect.Descriptor instead.
func (*DepthCache) Descriptor() ([]byte, []int) {
	return file_cache_old_proto_rawDescGZIP(), []int{1}
}

func (x *DepthCache) GetEntries() map[uint64]Ending {
	if x != nil {
		return x.Entries
	}
	return nil
}

var File_cache_old_proto protoreflect.FileDescriptor

var file_cache_old_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x63, 0x61, 0x63, 0x68, 0x65, 0x5f, 0x6f, 0x6c, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x5f, 0x6f, 0x6c, 0x64, 0x22, 0x46, 0x0a, 0x0b,
	0x44, 0x65, 0x70, 0x74, 0x68, 0x43, 0x61, 0x63, 0x68, 0x65, 0x73, 0x12, 0x37, 0x0a, 0x0b, 0x64,
	0x65, 0x70, 0x74, 0x68, 0x43, 0x61, 0x63, 0x68, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x5f, 0x6f, 0x6c, 0x64, 0x2e, 0x44, 0x65, 0x70,
	0x74, 0x68, 0x43, 0x61, 0x63, 0x68, 0x65, 0x52, 0x0b, 0x64, 0x65, 0x70, 0x74, 0x68, 0x43, 0x61,
	0x63, 0x68, 0x65, 0x73, 0x22, 0x99, 0x01, 0x0a, 0x0a, 0x44, 0x65, 0x70, 0x74, 0x68, 0x43, 0x61,
	0x63, 0x68, 0x65, 0x12, 0x3c, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x5f, 0x6f, 0x6c, 0x64,
	0x2e, 0x44, 0x65, 0x70, 0x74, 0x68, 0x43, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x45, 0x6e, 0x74, 0x72,
	0x69, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65,
	0x73, 0x1a, 0x4d, 0x0a, 0x0c, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x27, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x5f, 0x6f, 0x6c, 0x64, 0x2e, 0x45,
	0x6e, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01,
	0x2a, 0x2d, 0x0a, 0x06, 0x45, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x4c,
	0x41, 0x59, 0x45, 0x52, 0x5f, 0x41, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x50, 0x4c, 0x41, 0x59,
	0x45, 0x52, 0x5f, 0x42, 0x10, 0x01, 0x12, 0x07, 0x0a, 0x03, 0x54, 0x49, 0x45, 0x10, 0x02, 0x42,
	0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x67,
	0x72, 0x65, 0x6b, 0x35, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x34, 0x73, 0x6f,
	0x6c, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x5f, 0x6f, 0x6c, 0x64, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_cache_old_proto_rawDescOnce sync.Once
	file_cache_old_proto_rawDescData = file_cache_old_proto_rawDesc
)

func file_cache_old_proto_rawDescGZIP() []byte {
	file_cache_old_proto_rawDescOnce.Do(func() {
		file_cache_old_proto_rawDescData = protoimpl.X.CompressGZIP(file_cache_old_proto_rawDescData)
	})
	return file_cache_old_proto_rawDescData
}

var file_cache_old_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_cache_old_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_cache_old_proto_goTypes = []interface{}{
	(Ending)(0),         // 0: proto_old.Ending
	(*DepthCaches)(nil), // 1: proto_old.DepthCaches
	(*DepthCache)(nil),  // 2: proto_old.DepthCache
	nil,                 // 3: proto_old.DepthCache.EntriesEntry
}
var file_cache_old_proto_depIdxs = []int32{
	2, // 0: proto_old.DepthCaches.depthCaches:type_name -> proto_old.DepthCache
	3, // 1: proto_old.DepthCache.entries:type_name -> proto_old.DepthCache.EntriesEntry
	0, // 2: proto_old.DepthCache.EntriesEntry.value:type_name -> proto_old.Ending
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_cache_old_proto_init() }
func file_cache_old_proto_init() {
	if File_cache_old_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cache_old_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DepthCaches); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_cache_old_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DepthCache); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_cache_old_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_cache_old_proto_goTypes,
		DependencyIndexes: file_cache_old_proto_depIdxs,
		EnumInfos:         file_cache_old_proto_enumTypes,
		MessageInfos:      file_cache_old_proto_msgTypes,
	}.Build()
	File_cache_old_proto = out.File
	file_cache_old_proto_rawDesc = nil
	file_cache_old_proto_goTypes = nil
	file_cache_old_proto_depIdxs = nil
}
