//
// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.14.0
// source: filters/concatenate_bytes/v1alpha1/concatenate_bytes.proto

package quilkin_extensions_filters_concatenate_bytes_v1alpha1

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ConcatenateBytes_Strategy int32

const (
	ConcatenateBytes_DoNothing ConcatenateBytes_Strategy = 0
	ConcatenateBytes_Append    ConcatenateBytes_Strategy = 1
	ConcatenateBytes_Prepend   ConcatenateBytes_Strategy = 2
)

// Enum value maps for ConcatenateBytes_Strategy.
var (
	ConcatenateBytes_Strategy_name = map[int32]string{
		0: "DoNothing",
		1: "Append",
		2: "Prepend",
	}
	ConcatenateBytes_Strategy_value = map[string]int32{
		"DoNothing": 0,
		"Append":    1,
		"Prepend":   2,
	}
)

func (x ConcatenateBytes_Strategy) Enum() *ConcatenateBytes_Strategy {
	p := new(ConcatenateBytes_Strategy)
	*p = x
	return p
}

func (x ConcatenateBytes_Strategy) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ConcatenateBytes_Strategy) Descriptor() protoreflect.EnumDescriptor {
	return file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_enumTypes[0].Descriptor()
}

func (ConcatenateBytes_Strategy) Type() protoreflect.EnumType {
	return &file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_enumTypes[0]
}

func (x ConcatenateBytes_Strategy) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ConcatenateBytes_Strategy.Descriptor instead.
func (ConcatenateBytes_Strategy) EnumDescriptor() ([]byte, []int) {
	return file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_rawDescGZIP(), []int{0, 0}
}

type ConcatenateBytes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OnWrite *ConcatenateBytes_StrategyValue `protobuf:"bytes,1,opt,name=on_write,json=onWrite,proto3" json:"on_write,omitempty"`
	OnRead  *ConcatenateBytes_StrategyValue `protobuf:"bytes,2,opt,name=on_read,json=onRead,proto3" json:"on_read,omitempty"`
	Bytes   []byte                          `protobuf:"bytes,3,opt,name=bytes,proto3" json:"bytes,omitempty"`
}

func (x *ConcatenateBytes) Reset() {
	*x = ConcatenateBytes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConcatenateBytes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConcatenateBytes) ProtoMessage() {}

func (x *ConcatenateBytes) ProtoReflect() protoreflect.Message {
	mi := &file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConcatenateBytes.ProtoReflect.Descriptor instead.
func (*ConcatenateBytes) Descriptor() ([]byte, []int) {
	return file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_rawDescGZIP(), []int{0}
}

func (x *ConcatenateBytes) GetOnWrite() *ConcatenateBytes_StrategyValue {
	if x != nil {
		return x.OnWrite
	}
	return nil
}

func (x *ConcatenateBytes) GetOnRead() *ConcatenateBytes_StrategyValue {
	if x != nil {
		return x.OnRead
	}
	return nil
}

func (x *ConcatenateBytes) GetBytes() []byte {
	if x != nil {
		return x.Bytes
	}
	return nil
}

type ConcatenateBytes_StrategyValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value ConcatenateBytes_Strategy `protobuf:"varint,1,opt,name=value,proto3,enum=quilkin.extensions.filters.concatenate_bytes.v1alpha1.ConcatenateBytes_Strategy" json:"value,omitempty"`
}

func (x *ConcatenateBytes_StrategyValue) Reset() {
	*x = ConcatenateBytes_StrategyValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConcatenateBytes_StrategyValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConcatenateBytes_StrategyValue) ProtoMessage() {}

func (x *ConcatenateBytes_StrategyValue) ProtoReflect() protoreflect.Message {
	mi := &file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConcatenateBytes_StrategyValue.ProtoReflect.Descriptor instead.
func (*ConcatenateBytes_StrategyValue) Descriptor() ([]byte, []int) {
	return file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ConcatenateBytes_StrategyValue) GetValue() ConcatenateBytes_Strategy {
	if x != nil {
		return x.Value
	}
	return ConcatenateBytes_DoNothing
}

var File_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto protoreflect.FileDescriptor

var file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_rawDesc = []byte{
	0x0a, 0x3a, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2f, 0x63, 0x6f, 0x6e, 0x63, 0x61, 0x74,
	0x65, 0x6e, 0x61, 0x74, 0x65, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x2f, 0x76, 0x31, 0x61, 0x6c,
	0x70, 0x68, 0x61, 0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x63, 0x61, 0x74, 0x65, 0x6e, 0x61, 0x74, 0x65,
	0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x35, 0x71, 0x75,
	0x69, 0x6c, 0x6b, 0x69, 0x6e, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x63, 0x61, 0x74, 0x65,
	0x6e, 0x61, 0x74, 0x65, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x22, 0xb7, 0x03, 0x0a, 0x10, 0x43, 0x6f, 0x6e, 0x63, 0x61, 0x74, 0x65, 0x6e,
	0x61, 0x74, 0x65, 0x42, 0x79, 0x74, 0x65, 0x73, 0x12, 0x70, 0x0a, 0x08, 0x6f, 0x6e, 0x5f, 0x77,
	0x72, 0x69, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x55, 0x2e, 0x71, 0x75, 0x69,
	0x6c, 0x6b, 0x69, 0x6e, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x63, 0x61, 0x74, 0x65, 0x6e,
	0x61, 0x74, 0x65, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70, 0x68,
	0x61, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x63, 0x61, 0x74, 0x65, 0x6e, 0x61, 0x74, 0x65, 0x42, 0x79,
	0x74, 0x65, 0x73, 0x2e, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x07, 0x6f, 0x6e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x12, 0x6e, 0x0a, 0x07, 0x6f, 0x6e,
	0x5f, 0x72, 0x65, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x55, 0x2e, 0x71, 0x75,
	0x69, 0x6c, 0x6b, 0x69, 0x6e, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73,
	0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x63, 0x6f, 0x6e, 0x63, 0x61, 0x74, 0x65,
	0x6e, 0x61, 0x74, 0x65, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x61, 0x6c, 0x70,
	0x68, 0x61, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x63, 0x61, 0x74, 0x65, 0x6e, 0x61, 0x74, 0x65, 0x42,
	0x79, 0x74, 0x65, 0x73, 0x2e, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x06, 0x6f, 0x6e, 0x52, 0x65, 0x61, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x79,
	0x74, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x62, 0x79, 0x74, 0x65, 0x73,
	0x1a, 0x77, 0x0a, 0x0d, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65, 0x67, 0x79, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x12, 0x66, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x50, 0x2e, 0x71, 0x75, 0x69, 0x6c, 0x6b, 0x69, 0x6e, 0x2e, 0x65, 0x78, 0x74, 0x65, 0x6e,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x73, 0x2e, 0x63, 0x6f,
	0x6e, 0x63, 0x61, 0x74, 0x65, 0x6e, 0x61, 0x74, 0x65, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x2e,
	0x76, 0x31, 0x61, 0x6c, 0x70, 0x68, 0x61, 0x31, 0x2e, 0x43, 0x6f, 0x6e, 0x63, 0x61, 0x74, 0x65,
	0x6e, 0x61, 0x74, 0x65, 0x42, 0x79, 0x74, 0x65, 0x73, 0x2e, 0x53, 0x74, 0x72, 0x61, 0x74, 0x65,
	0x67, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x32, 0x0a, 0x08, 0x53, 0x74, 0x72,
	0x61, 0x74, 0x65, 0x67, 0x79, 0x12, 0x0d, 0x0a, 0x09, 0x44, 0x6f, 0x4e, 0x6f, 0x74, 0x68, 0x69,
	0x6e, 0x67, 0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x10, 0x01,
	0x12, 0x0b, 0x0a, 0x07, 0x50, 0x72, 0x65, 0x70, 0x65, 0x6e, 0x64, 0x10, 0x02, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_rawDescOnce sync.Once
	file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_rawDescData = file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_rawDesc
)

func file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_rawDescGZIP() []byte {
	file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_rawDescOnce.Do(func() {
		file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_rawDescData = protoimpl.X.CompressGZIP(file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_rawDescData)
	})
	return file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_rawDescData
}

var file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_goTypes = []interface{}{
	(ConcatenateBytes_Strategy)(0),         // 0: quilkin.extensions.filters.concatenate_bytes.v1alpha1.ConcatenateBytes.Strategy
	(*ConcatenateBytes)(nil),               // 1: quilkin.extensions.filters.concatenate_bytes.v1alpha1.ConcatenateBytes
	(*ConcatenateBytes_StrategyValue)(nil), // 2: quilkin.extensions.filters.concatenate_bytes.v1alpha1.ConcatenateBytes.StrategyValue
}
var file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_depIdxs = []int32{
	2, // 0: quilkin.extensions.filters.concatenate_bytes.v1alpha1.ConcatenateBytes.on_write:type_name -> quilkin.extensions.filters.concatenate_bytes.v1alpha1.ConcatenateBytes.StrategyValue
	2, // 1: quilkin.extensions.filters.concatenate_bytes.v1alpha1.ConcatenateBytes.on_read:type_name -> quilkin.extensions.filters.concatenate_bytes.v1alpha1.ConcatenateBytes.StrategyValue
	0, // 2: quilkin.extensions.filters.concatenate_bytes.v1alpha1.ConcatenateBytes.StrategyValue.value:type_name -> quilkin.extensions.filters.concatenate_bytes.v1alpha1.ConcatenateBytes.Strategy
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_init() }
func file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_init() {
	if File_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConcatenateBytes); i {
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
		file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConcatenateBytes_StrategyValue); i {
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
			RawDescriptor: file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_goTypes,
		DependencyIndexes: file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_depIdxs,
		EnumInfos:         file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_enumTypes,
		MessageInfos:      file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_msgTypes,
	}.Build()
	File_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto = out.File
	file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_rawDesc = nil
	file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_goTypes = nil
	file_filters_concatenate_bytes_v1alpha1_concatenate_bytes_proto_depIdxs = nil
}
