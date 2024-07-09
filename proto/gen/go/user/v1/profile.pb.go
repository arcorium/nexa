// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: user/v1/profile.proto

package userv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type FindProfileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserIds []string `protobuf:"bytes,1,rep,name=user_ids,json=userIds,proto3" json:"user_ids,omitempty"`
}

func (x *FindProfileRequest) Reset() {
	*x = FindProfileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_v1_profile_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindProfileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindProfileRequest) ProtoMessage() {}

func (x *FindProfileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_v1_profile_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindProfileRequest.ProtoReflect.Descriptor instead.
func (*FindProfileRequest) Descriptor() ([]byte, []int) {
	return file_user_v1_profile_proto_rawDescGZIP(), []int{0}
}

func (x *FindProfileRequest) GetUserIds() []string {
	if x != nil {
		return x.UserIds
	}
	return nil
}

type FindProfileResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Profiles []*Profile `protobuf:"bytes,1,rep,name=profiles,proto3" json:"profiles,omitempty"`
}

func (x *FindProfileResponse) Reset() {
	*x = FindProfileResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_v1_profile_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindProfileResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindProfileResponse) ProtoMessage() {}

func (x *FindProfileResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_v1_profile_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindProfileResponse.ProtoReflect.Descriptor instead.
func (*FindProfileResponse) Descriptor() ([]byte, []int) {
	return file_user_v1_profile_proto_rawDescGZIP(), []int{1}
}

func (x *FindProfileResponse) GetProfiles() []*Profile {
	if x != nil {
		return x.Profiles
	}
	return nil
}

type Profile struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FirstName string `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName  string `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Bio       string `protobuf:"bytes,4,opt,name=bio,proto3" json:"bio,omitempty"`
	ImagePath string `protobuf:"bytes,5,opt,name=image_path,json=imagePath,proto3" json:"image_path,omitempty"`
}

func (x *Profile) Reset() {
	*x = Profile{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_v1_profile_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Profile) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Profile) ProtoMessage() {}

func (x *Profile) ProtoReflect() protoreflect.Message {
	mi := &file_user_v1_profile_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Profile.ProtoReflect.Descriptor instead.
func (*Profile) Descriptor() ([]byte, []int) {
	return file_user_v1_profile_proto_rawDescGZIP(), []int{2}
}

func (x *Profile) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Profile) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *Profile) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *Profile) GetBio() string {
	if x != nil {
		return x.Bio
	}
	return ""
}

func (x *Profile) GetImagePath() string {
	if x != nil {
		return x.ImagePath
	}
	return ""
}

type UpdateProfileRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FirstName *string `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3,oneof" json:"first_name,omitempty"`
	LastName  *string `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3,oneof" json:"last_name,omitempty"`
	Bio       *string `protobuf:"bytes,4,opt,name=bio,proto3,oneof" json:"bio,omitempty"`
}

func (x *UpdateProfileRequest) Reset() {
	*x = UpdateProfileRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_v1_profile_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateProfileRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateProfileRequest) ProtoMessage() {}

func (x *UpdateProfileRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_v1_profile_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateProfileRequest.ProtoReflect.Descriptor instead.
func (*UpdateProfileRequest) Descriptor() ([]byte, []int) {
	return file_user_v1_profile_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateProfileRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateProfileRequest) GetFirstName() string {
	if x != nil && x.FirstName != nil {
		return *x.FirstName
	}
	return ""
}

func (x *UpdateProfileRequest) GetLastName() string {
	if x != nil && x.LastName != nil {
		return *x.LastName
	}
	return ""
}

func (x *UpdateProfileRequest) GetBio() string {
	if x != nil && x.Bio != nil {
		return *x.Bio
	}
	return ""
}

type UpdateProfileAvatarRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Filename string `protobuf:"bytes,2,opt,name=filename,proto3" json:"filename,omitempty"`
	Chunk    []byte `protobuf:"bytes,3,opt,name=chunk,proto3" json:"chunk,omitempty"`
}

func (x *UpdateProfileAvatarRequest) Reset() {
	*x = UpdateProfileAvatarRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_v1_profile_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateProfileAvatarRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateProfileAvatarRequest) ProtoMessage() {}

func (x *UpdateProfileAvatarRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_v1_profile_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateProfileAvatarRequest.ProtoReflect.Descriptor instead.
func (*UpdateProfileAvatarRequest) Descriptor() ([]byte, []int) {
	return file_user_v1_profile_proto_rawDescGZIP(), []int{4}
}

func (x *UpdateProfileAvatarRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateProfileAvatarRequest) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *UpdateProfileAvatarRequest) GetChunk() []byte {
	if x != nil {
		return x.Chunk
	}
	return nil
}

var File_user_v1_profile_proto protoreflect.FileDescriptor

var file_user_v1_profile_proto_rawDesc = []byte{
	0x0a, 0x15, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x75, 0x73,
	0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x2f, 0x0a, 0x12, 0x46, 0x69, 0x6e, 0x64, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x73, 0x22, 0x48, 0x0a, 0x13, 0x46, 0x69, 0x6e, 0x64, 0x50, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x08, 0x70, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x6e,
	0x65, 0x78, 0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x72, 0x6f, 0x66,
	0x69, 0x6c, 0x65, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x22, 0x86, 0x01,
	0x0a, 0x07, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x72,
	0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66,
	0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x61, 0x73, 0x74,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x61, 0x73,
	0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x62, 0x69, 0x6f, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x62, 0x69, 0x6f, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x6d, 0x61, 0x67, 0x65,
	0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x50, 0x61, 0x74, 0x68, 0x22, 0xa8, 0x01, 0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x22, 0x0a, 0x0a, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x20, 0x0a, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e, 0x61,
	0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x15, 0x0a, 0x03, 0x62, 0x69, 0x6f, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x02, 0x52, 0x03, 0x62, 0x69, 0x6f, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b,
	0x5f, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x0c, 0x0a, 0x0a, 0x5f,
	0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x5f, 0x62, 0x69,
	0x6f, 0x22, 0x5e, 0x0a, 0x1a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x66, 0x69,
	0x6c, 0x65, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63,
	0x68, 0x75, 0x6e, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x63, 0x68, 0x75, 0x6e,
	0x6b, 0x32, 0xf7, 0x01, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0x4b, 0x0a, 0x04, 0x46, 0x69, 0x6e, 0x64, 0x12, 0x20, 0x2e, 0x6e,
	0x65, 0x78, 0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x6e, 0x64,
	0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21,
	0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69,
	0x6e, 0x64, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x44, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x22, 0x2e, 0x6e, 0x65,
	0x78, 0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x52, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x12, 0x28, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x28, 0x01, 0x42, 0xa8, 0x01, 0x0a, 0x10,
	0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x42, 0x0c, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x72, 0x63,
	0x6f, 0x72, 0x69, 0x75, 0x6d, 0x2f, 0x6e, 0x65, 0x78, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b,
	0x75, 0x73, 0x65, 0x72, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x4e, 0x55, 0x58, 0xaa, 0x02, 0x0c, 0x4e,
	0x65, 0x78, 0x61, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0c, 0x4e, 0x65,
	0x78, 0x61, 0x5c, 0x55, 0x73, 0x65, 0x72, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x18, 0x4e, 0x65, 0x78,
	0x61, 0x5c, 0x55, 0x73, 0x65, 0x72, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0e, 0x4e, 0x65, 0x78, 0x61, 0x3a, 0x3a, 0x55, 0x73,
	0x65, 0x72, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_user_v1_profile_proto_rawDescOnce sync.Once
	file_user_v1_profile_proto_rawDescData = file_user_v1_profile_proto_rawDesc
)

func file_user_v1_profile_proto_rawDescGZIP() []byte {
	file_user_v1_profile_proto_rawDescOnce.Do(func() {
		file_user_v1_profile_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_v1_profile_proto_rawDescData)
	})
	return file_user_v1_profile_proto_rawDescData
}

var file_user_v1_profile_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_user_v1_profile_proto_goTypes = []interface{}{
	(*FindProfileRequest)(nil),         // 0: nexa.user.v1.FindProfileRequest
	(*FindProfileResponse)(nil),        // 1: nexa.user.v1.FindProfileResponse
	(*Profile)(nil),                    // 2: nexa.user.v1.Profile
	(*UpdateProfileRequest)(nil),       // 3: nexa.user.v1.UpdateProfileRequest
	(*UpdateProfileAvatarRequest)(nil), // 4: nexa.user.v1.UpdateProfileAvatarRequest
	(*emptypb.Empty)(nil),              // 5: google.protobuf.Empty
}
var file_user_v1_profile_proto_depIdxs = []int32{
	2, // 0: nexa.user.v1.FindProfileResponse.profiles:type_name -> nexa.user.v1.Profile
	0, // 1: nexa.user.v1.ProfileService.Find:input_type -> nexa.user.v1.FindProfileRequest
	3, // 2: nexa.user.v1.ProfileService.Update:input_type -> nexa.user.v1.UpdateProfileRequest
	4, // 3: nexa.user.v1.ProfileService.UpdateAvatar:input_type -> nexa.user.v1.UpdateProfileAvatarRequest
	1, // 4: nexa.user.v1.ProfileService.Find:output_type -> nexa.user.v1.FindProfileResponse
	5, // 5: nexa.user.v1.ProfileService.Update:output_type -> google.protobuf.Empty
	5, // 6: nexa.user.v1.ProfileService.UpdateAvatar:output_type -> google.protobuf.Empty
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_user_v1_profile_proto_init() }
func file_user_v1_profile_proto_init() {
	if File_user_v1_profile_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_user_v1_profile_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindProfileRequest); i {
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
		file_user_v1_profile_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindProfileResponse); i {
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
		file_user_v1_profile_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Profile); i {
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
		file_user_v1_profile_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateProfileRequest); i {
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
		file_user_v1_profile_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateProfileAvatarRequest); i {
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
	file_user_v1_profile_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_user_v1_profile_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_v1_profile_proto_goTypes,
		DependencyIndexes: file_user_v1_profile_proto_depIdxs,
		MessageInfos:      file_user_v1_profile_proto_msgTypes,
	}.Build()
	File_user_v1_profile_proto = out.File
	file_user_v1_profile_proto_rawDesc = nil
	file_user_v1_profile_proto_goTypes = nil
	file_user_v1_profile_proto_depIdxs = nil
}
