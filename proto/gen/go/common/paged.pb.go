// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: common/paged.proto

package common

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

type PagedElementInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Element uint64 `protobuf:"varint,1,opt,name=element,proto3" json:"element,omitempty"`
	Page    uint64 `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
}

func (x *PagedElementInput) Reset() {
	*x = PagedElementInput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_paged_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PagedElementInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PagedElementInput) ProtoMessage() {}

func (x *PagedElementInput) ProtoReflect() protoreflect.Message {
	mi := &file_common_paged_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PagedElementInput.ProtoReflect.Descriptor instead.
func (*PagedElementInput) Descriptor() ([]byte, []int) {
	return file_common_paged_proto_rawDescGZIP(), []int{0}
}

func (x *PagedElementInput) GetElement() uint64 {
	if x != nil {
		return x.Element
	}
	return 0
}

func (x *PagedElementInput) GetPage() uint64 {
	if x != nil {
		return x.Page
	}
	return 0
}

type PagedElementOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Element       uint64 `protobuf:"varint,1,opt,name=element,proto3" json:"element,omitempty"`
	Page          uint64 `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	TotalElements uint64 `protobuf:"varint,3,opt,name=totalElements,proto3" json:"totalElements,omitempty"`
	TotalPages    uint64 `protobuf:"varint,4,opt,name=totalPages,proto3" json:"totalPages,omitempty"`
}

func (x *PagedElementOutput) Reset() {
	*x = PagedElementOutput{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_paged_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PagedElementOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PagedElementOutput) ProtoMessage() {}

func (x *PagedElementOutput) ProtoReflect() protoreflect.Message {
	mi := &file_common_paged_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PagedElementOutput.ProtoReflect.Descriptor instead.
func (*PagedElementOutput) Descriptor() ([]byte, []int) {
	return file_common_paged_proto_rawDescGZIP(), []int{1}
}

func (x *PagedElementOutput) GetElement() uint64 {
	if x != nil {
		return x.Element
	}
	return 0
}

func (x *PagedElementOutput) GetPage() uint64 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *PagedElementOutput) GetTotalElements() uint64 {
	if x != nil {
		return x.TotalElements
	}
	return 0
}

func (x *PagedElementOutput) GetTotalPages() uint64 {
	if x != nil {
		return x.TotalPages
	}
	return 0
}

var File_common_paged_proto protoreflect.FileDescriptor

var file_common_paged_proto_rawDesc = []byte{
	0x0a, 0x12, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x61, 0x67, 0x65, 0x64, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x22, 0x41, 0x0a, 0x11, 0x50, 0x61, 0x67, 0x65, 0x64, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04,
	0x70, 0x61, 0x67, 0x65, 0x22, 0x88, 0x01, 0x0a, 0x12, 0x50, 0x61, 0x67, 0x65, 0x64, 0x45, 0x6c,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x65,
	0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x07, 0x65, 0x6c,
	0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0d, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12,
	0x1e, 0x0a, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x61, 0x67, 0x65, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x50, 0x61, 0x67, 0x65, 0x73, 0x42,
	0x98, 0x01, 0x0a, 0x0f, 0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x63, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0x42, 0x0a, 0x50, 0x61, 0x67, 0x65, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50,
	0x01, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x72,
	0x63, 0x6f, 0x72, 0x69, 0x75, 0x6d, 0x2f, 0x6e, 0x65, 0x78, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0xa2,
	0x02, 0x03, 0x4e, 0x43, 0x58, 0xaa, 0x02, 0x0b, 0x4e, 0x65, 0x78, 0x61, 0x2e, 0x43, 0x6f, 0x6d,
	0x6d, 0x6f, 0x6e, 0xca, 0x02, 0x0b, 0x4e, 0x65, 0x78, 0x61, 0x5c, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0xe2, 0x02, 0x17, 0x4e, 0x65, 0x78, 0x61, 0x5c, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c,
	0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0c, 0x4e, 0x65,
	0x78, 0x61, 0x3a, 0x3a, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_common_paged_proto_rawDescOnce sync.Once
	file_common_paged_proto_rawDescData = file_common_paged_proto_rawDesc
)

func file_common_paged_proto_rawDescGZIP() []byte {
	file_common_paged_proto_rawDescOnce.Do(func() {
		file_common_paged_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_paged_proto_rawDescData)
	})
	return file_common_paged_proto_rawDescData
}

var file_common_paged_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_common_paged_proto_goTypes = []interface{}{
	(*PagedElementInput)(nil),  // 0: nexa.common.PagedElementInput
	(*PagedElementOutput)(nil), // 1: nexa.common.PagedElementOutput
}
var file_common_paged_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_common_paged_proto_init() }
func file_common_paged_proto_init() {
	if File_common_paged_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_paged_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PagedElementInput); i {
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
		file_common_paged_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PagedElementOutput); i {
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
			RawDescriptor: file_common_paged_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_paged_proto_goTypes,
		DependencyIndexes: file_common_paged_proto_depIdxs,
		MessageInfos:      file_common_paged_proto_msgTypes,
	}.Build()
	File_common_paged_proto = out.File
	file_common_paged_proto_rawDesc = nil
	file_common_paged_proto_goTypes = nil
	file_common_paged_proto_depIdxs = nil
}
