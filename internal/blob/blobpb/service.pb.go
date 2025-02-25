// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/blob/blobpb/*.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.3
// 	protoc        v3.12.4
// source: internal/blob/blobpb/service.proto

package blobpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_internal_blob_blobpb_service_proto protoreflect.FileDescriptor

var file_internal_blob_blobpb_service_proto_rawDesc = []byte{
	0x0a, 0x22, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x62, 0x6c, 0x6f, 0x62, 0x2f,
	0x62, 0x6c, 0x6f, 0x62, 0x70, 0x62, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x62, 0x6c, 0x6f, 0x62, 0x70, 0x62, 0x1a, 0x22, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x62, 0x6c, 0x6f, 0x62, 0x2f, 0x62, 0x6c, 0x6f, 0x62,
	0x70, 0x62, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x32, 0xb9, 0x03, 0x0a, 0x0e, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x3f, 0x0a, 0x0c, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x4f, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x12, 0x15, 0x2e, 0x62, 0x6c, 0x6f, 0x62, 0x70, 0x62, 0x2e, 0x55, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x62, 0x6c, 0x6f,
	0x62, 0x70, 0x62, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x28, 0x01, 0x12, 0x45, 0x0a, 0x0e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64,
	0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x17, 0x2e, 0x62, 0x6c, 0x6f, 0x62, 0x70, 0x62, 0x2e,
	0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x18, 0x2e, 0x62, 0x6c, 0x6f, 0x62, 0x70, 0x62, 0x2e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61,
	0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x3d, 0x0a, 0x0c, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x15, 0x2e, 0x62, 0x6c,
	0x6f, 0x62, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x16, 0x2e, 0x62, 0x6c, 0x6f, 0x62, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x46, 0x0a, 0x0b, 0x53, 0x68,
	0x61, 0x72, 0x65, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x1a, 0x2e, 0x62, 0x6c, 0x6f, 0x62,
	0x70, 0x62, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x65, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x62, 0x6c, 0x6f, 0x62, 0x70, 0x62, 0x2e, 0x53,
	0x68, 0x61, 0x72, 0x65, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x51, 0x0a, 0x14, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x68,
	0x61, 0x72, 0x65, 0x64, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x12, 0x1d, 0x2e, 0x62, 0x6c, 0x6f,
	0x62, 0x70, 0x62, 0x2e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x68, 0x61, 0x72,
	0x65, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x62, 0x6c, 0x6f, 0x62,
	0x70, 0x62, 0x2e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x45, 0x0a, 0x0a, 0x4c, 0x69, 0x73, 0x74, 0x4f, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x12, 0x19, 0x2e, 0x62, 0x6c, 0x6f, 0x62, 0x70, 0x62, 0x2e, 0x4c, 0x69, 0x73,
	0x74, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a,
	0x2e, 0x62, 0x6c, 0x6f, 0x62, 0x70, 0x62, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x42, 0x0a, 0x5a, 0x08,
	0x2e, 0x2f, 0x62, 0x6c, 0x6f, 0x62, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_internal_blob_blobpb_service_proto_goTypes = []any{
	(*UploadRequest)(nil),         // 0: blobpb.UploadRequest
	(*DownloadRequest)(nil),       // 1: blobpb.DownloadRequest
	(*DeleteRequest)(nil),         // 2: blobpb.DeleteRequest
	(*ShareObjectRequest)(nil),    // 3: blobpb.ShareObjectRequest
	(*DownloadSharedRequest)(nil), // 4: blobpb.DownloadSharedRequest
	(*ListObjectRequest)(nil),     // 5: blobpb.ListObjectRequest
	(*UploadResponse)(nil),        // 6: blobpb.UploadResponse
	(*DownloadResponse)(nil),      // 7: blobpb.DownloadResponse
	(*DeleteResponse)(nil),        // 8: blobpb.DeleteResponse
	(*ShareObjectResponse)(nil),   // 9: blobpb.ShareObjectResponse
	(*ListObjectResponse)(nil),    // 10: blobpb.ListObjectResponse
}
var file_internal_blob_blobpb_service_proto_depIdxs = []int32{
	0,  // 0: blobpb.StorageService.UploadObject:input_type -> blobpb.UploadRequest
	1,  // 1: blobpb.StorageService.DownloadObject:input_type -> blobpb.DownloadRequest
	2,  // 2: blobpb.StorageService.DeleteObject:input_type -> blobpb.DeleteRequest
	3,  // 3: blobpb.StorageService.ShareObject:input_type -> blobpb.ShareObjectRequest
	4,  // 4: blobpb.StorageService.DownloadSharedObject:input_type -> blobpb.DownloadSharedRequest
	5,  // 5: blobpb.StorageService.ListObject:input_type -> blobpb.ListObjectRequest
	6,  // 6: blobpb.StorageService.UploadObject:output_type -> blobpb.UploadResponse
	7,  // 7: blobpb.StorageService.DownloadObject:output_type -> blobpb.DownloadResponse
	8,  // 8: blobpb.StorageService.DeleteObject:output_type -> blobpb.DeleteResponse
	9,  // 9: blobpb.StorageService.ShareObject:output_type -> blobpb.ShareObjectResponse
	7,  // 10: blobpb.StorageService.DownloadSharedObject:output_type -> blobpb.DownloadResponse
	10, // 11: blobpb.StorageService.ListObject:output_type -> blobpb.ListObjectResponse
	6,  // [6:12] is the sub-list for method output_type
	0,  // [0:6] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_internal_blob_blobpb_service_proto_init() }
func file_internal_blob_blobpb_service_proto_init() {
	if File_internal_blob_blobpb_service_proto != nil {
		return
	}
	file_internal_blob_blobpb_message_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_blob_blobpb_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_blob_blobpb_service_proto_goTypes,
		DependencyIndexes: file_internal_blob_blobpb_service_proto_depIdxs,
	}.Build()
	File_internal_blob_blobpb_service_proto = out.File
	file_internal_blob_blobpb_service_proto_rawDesc = nil
	file_internal_blob_blobpb_service_proto_goTypes = nil
	file_internal_blob_blobpb_service_proto_depIdxs = nil
}
