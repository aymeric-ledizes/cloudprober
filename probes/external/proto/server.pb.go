// This package defines protocol for use with cloudprober's external
// probe (in server mode).

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0
// 	protoc        v3.11.2
// source: github.com/google/cloudprober/probes/external/proto/server.proto

package proto

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

// ProbeRequest is the message that cloudprober sends to the external probe
// server.
type ProbeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The unique identifier for this request.  This is unique across
	// an execution of the probe server.  It starts at 1.
	RequestId *int32 `protobuf:"varint,1,req,name=request_id,json=requestId" json:"request_id,omitempty"`
	// How long to allow for the execution of this request, in
	// milliseconds.  If the time limit is exceeded, the server
	// should abort the request, but *not* send back a reply.  The
	// client will have to do timeouts anyway.
	TimeLimit *int32                 `protobuf:"varint,2,req,name=time_limit,json=timeLimit" json:"time_limit,omitempty"`
	Options   []*ProbeRequest_Option `protobuf:"bytes,3,rep,name=options" json:"options,omitempty"`
}

func (x *ProbeRequest) Reset() {
	*x = ProbeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_google_cloudprober_probes_external_proto_server_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProbeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProbeRequest) ProtoMessage() {}

func (x *ProbeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_google_cloudprober_probes_external_proto_server_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProbeRequest.ProtoReflect.Descriptor instead.
func (*ProbeRequest) Descriptor() ([]byte, []int) {
	return file_github_com_google_cloudprober_probes_external_proto_server_proto_rawDescGZIP(), []int{0}
}

func (x *ProbeRequest) GetRequestId() int32 {
	if x != nil && x.RequestId != nil {
		return *x.RequestId
	}
	return 0
}

func (x *ProbeRequest) GetTimeLimit() int32 {
	if x != nil && x.TimeLimit != nil {
		return *x.TimeLimit
	}
	return 0
}

func (x *ProbeRequest) GetOptions() []*ProbeRequest_Option {
	if x != nil {
		return x.Options
	}
	return nil
}

// ProbeReply is the message that external probe server sends back to the
// cloudprober.
type ProbeReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The sequence number for this request.
	RequestId *int32 `protobuf:"varint,1,req,name=request_id,json=requestId" json:"request_id,omitempty"`
	// For a normal result, this is not present.
	// If it is present, it indicates that the probe failed.
	ErrorMessage *string `protobuf:"bytes,2,opt,name=error_message,json=errorMessage" json:"error_message,omitempty"`
	// The result of the probe. Cloudprober parses the payload to retrieve
	// variables from it. It expects variables in the following format:
	// var1 value1 (for example: total_errors 589)
	// TODO(manugarg): Add an option to export mapped variables, for example:
	// client-errors map:lang java:200 python:20 golang:3
	Payload *string `protobuf:"bytes,3,opt,name=payload" json:"payload,omitempty"`
}

func (x *ProbeReply) Reset() {
	*x = ProbeReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_google_cloudprober_probes_external_proto_server_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProbeReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProbeReply) ProtoMessage() {}

func (x *ProbeReply) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_google_cloudprober_probes_external_proto_server_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProbeReply.ProtoReflect.Descriptor instead.
func (*ProbeReply) Descriptor() ([]byte, []int) {
	return file_github_com_google_cloudprober_probes_external_proto_server_proto_rawDescGZIP(), []int{1}
}

func (x *ProbeReply) GetRequestId() int32 {
	if x != nil && x.RequestId != nil {
		return *x.RequestId
	}
	return 0
}

func (x *ProbeReply) GetErrorMessage() string {
	if x != nil && x.ErrorMessage != nil {
		return *x.ErrorMessage
	}
	return ""
}

func (x *ProbeReply) GetPayload() string {
	if x != nil && x.Payload != nil {
		return *x.Payload
	}
	return ""
}

type ProbeRequest_Option struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name  *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	Value *string `protobuf:"bytes,2,req,name=value" json:"value,omitempty"`
}

func (x *ProbeRequest_Option) Reset() {
	*x = ProbeRequest_Option{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_google_cloudprober_probes_external_proto_server_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProbeRequest_Option) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProbeRequest_Option) ProtoMessage() {}

func (x *ProbeRequest_Option) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_google_cloudprober_probes_external_proto_server_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProbeRequest_Option.ProtoReflect.Descriptor instead.
func (*ProbeRequest_Option) Descriptor() ([]byte, []int) {
	return file_github_com_google_cloudprober_probes_external_proto_server_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ProbeRequest_Option) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *ProbeRequest_Option) GetValue() string {
	if x != nil && x.Value != nil {
		return *x.Value
	}
	return ""
}

var File_github_com_google_cloudprober_probes_external_proto_server_proto protoreflect.FileDescriptor

var file_github_com_google_cloudprober_probes_external_proto_server_proto_rawDesc = []byte{
	0x0a, 0x40, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x70, 0x72, 0x6f, 0x62, 0x65, 0x72, 0x2f,
	0x70, 0x72, 0x6f, 0x62, 0x65, 0x73, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0b, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x70, 0x72, 0x6f, 0x62, 0x65, 0x72, 0x22,
	0xbc, 0x01, 0x0a, 0x0c, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x02, 0x28, 0x05, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x12,
	0x1d, 0x0a, 0x0a, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20,
	0x02, 0x28, 0x05, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x3a,
	0x0a, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x20, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x70, 0x72, 0x6f, 0x62, 0x65, 0x72, 0x2e, 0x50, 0x72,
	0x6f, 0x62, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x07, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x1a, 0x32, 0x0a, 0x06, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x02,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x02, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x6a,
	0x0a, 0x0a, 0x50, 0x72, 0x6f, 0x62, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1d, 0x0a, 0x0a,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x02, 0x28, 0x05,
	0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x35, 0x5a, 0x33, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x63, 0x6c, 0x6f, 0x75, 0x64, 0x70, 0x72, 0x6f, 0x62, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x62,
	0x65, 0x73, 0x2f, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f,
}

var (
	file_github_com_google_cloudprober_probes_external_proto_server_proto_rawDescOnce sync.Once
	file_github_com_google_cloudprober_probes_external_proto_server_proto_rawDescData = file_github_com_google_cloudprober_probes_external_proto_server_proto_rawDesc
)

func file_github_com_google_cloudprober_probes_external_proto_server_proto_rawDescGZIP() []byte {
	file_github_com_google_cloudprober_probes_external_proto_server_proto_rawDescOnce.Do(func() {
		file_github_com_google_cloudprober_probes_external_proto_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_google_cloudprober_probes_external_proto_server_proto_rawDescData)
	})
	return file_github_com_google_cloudprober_probes_external_proto_server_proto_rawDescData
}

var file_github_com_google_cloudprober_probes_external_proto_server_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_github_com_google_cloudprober_probes_external_proto_server_proto_goTypes = []interface{}{
	(*ProbeRequest)(nil),        // 0: cloudprober.ProbeRequest
	(*ProbeReply)(nil),          // 1: cloudprober.ProbeReply
	(*ProbeRequest_Option)(nil), // 2: cloudprober.ProbeRequest.Option
}
var file_github_com_google_cloudprober_probes_external_proto_server_proto_depIdxs = []int32{
	2, // 0: cloudprober.ProbeRequest.options:type_name -> cloudprober.ProbeRequest.Option
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_github_com_google_cloudprober_probes_external_proto_server_proto_init() }
func file_github_com_google_cloudprober_probes_external_proto_server_proto_init() {
	if File_github_com_google_cloudprober_probes_external_proto_server_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_google_cloudprober_probes_external_proto_server_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProbeRequest); i {
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
		file_github_com_google_cloudprober_probes_external_proto_server_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProbeReply); i {
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
		file_github_com_google_cloudprober_probes_external_proto_server_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProbeRequest_Option); i {
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
			RawDescriptor: file_github_com_google_cloudprober_probes_external_proto_server_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_google_cloudprober_probes_external_proto_server_proto_goTypes,
		DependencyIndexes: file_github_com_google_cloudprober_probes_external_proto_server_proto_depIdxs,
		MessageInfos:      file_github_com_google_cloudprober_probes_external_proto_server_proto_msgTypes,
	}.Build()
	File_github_com_google_cloudprober_probes_external_proto_server_proto = out.File
	file_github_com_google_cloudprober_probes_external_proto_server_proto_rawDesc = nil
	file_github_com_google_cloudprober_probes_external_proto_server_proto_goTypes = nil
	file_github_com_google_cloudprober_probes_external_proto_server_proto_depIdxs = nil
}
