// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: authentication/v1/authentication.proto

package authenticationv1

import (
	_ "github.com/arcorium/nexa/proto/gen/go/common"
	_ "github.com/arcorium/nexa/proto/gen/go/token/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
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

type Credential struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Device    string `protobuf:"bytes,2,opt,name=device,proto3" json:"device,omitempty"`
	IsCurrent *bool  `protobuf:"varint,3,opt,name=is_current,json=isCurrent,proto3,oneof" json:"is_current,omitempty"`
}

func (x *Credential) Reset() {
	*x = Credential{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authentication_v1_authentication_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Credential) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Credential) ProtoMessage() {}

func (x *Credential) ProtoReflect() protoreflect.Message {
	mi := &file_authentication_v1_authentication_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Credential.ProtoReflect.Descriptor instead.
func (*Credential) Descriptor() ([]byte, []int) {
	return file_authentication_v1_authentication_proto_rawDescGZIP(), []int{0}
}

func (x *Credential) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Credential) GetDevice() string {
	if x != nil {
		return x.Device
	}
	return ""
}

func (x *Credential) GetIsCurrent() bool {
	if x != nil && x.IsCurrent != nil {
		return *x.IsCurrent
	}
	return false
}

type LoginRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email      string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password   string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	DeviceName string `protobuf:"bytes,3,opt,name=device_name,json=deviceName,proto3" json:"device_name,omitempty"`
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authentication_v1_authentication_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_authentication_v1_authentication_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_authentication_v1_authentication_proto_rawDescGZIP(), []int{1}
}

func (x *LoginRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *LoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *LoginRequest) GetDeviceName() string {
	if x != nil {
		return x.DeviceName
	}
	return ""
}

type LoginResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TokenType   string               `protobuf:"bytes,1,opt,name=token_type,json=tokenType,proto3" json:"token_type,omitempty"`
	AccessToken string               `protobuf:"bytes,2,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	ExpiryTime  *durationpb.Duration `protobuf:"bytes,3,opt,name=expiry_time,json=expiryTime,proto3" json:"expiry_time,omitempty"`
}

func (x *LoginResponse) Reset() {
	*x = LoginResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authentication_v1_authentication_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse) ProtoMessage() {}

func (x *LoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_authentication_v1_authentication_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return file_authentication_v1_authentication_proto_rawDescGZIP(), []int{2}
}

func (x *LoginResponse) GetTokenType() string {
	if x != nil {
		return x.TokenType
	}
	return ""
}

func (x *LoginResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *LoginResponse) GetExpiryTime() *durationpb.Duration {
	if x != nil {
		return x.ExpiryTime
	}
	return nil
}

type RegisterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username  string  `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Email     string  `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	Password  string  `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	FirstName string  `protobuf:"bytes,4,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName  *string `protobuf:"bytes,5,opt,name=last_name,json=lastName,proto3,oneof" json:"last_name,omitempty"`
	Bio       *string `protobuf:"bytes,6,opt,name=bio,proto3,oneof" json:"bio,omitempty"`
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authentication_v1_authentication_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_authentication_v1_authentication_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_authentication_v1_authentication_proto_rawDescGZIP(), []int{3}
}

func (x *RegisterRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *RegisterRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *RegisterRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *RegisterRequest) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *RegisterRequest) GetLastName() string {
	if x != nil && x.LastName != nil {
		return *x.LastName
	}
	return ""
}

func (x *RegisterRequest) GetBio() string {
	if x != nil && x.Bio != nil {
		return *x.Bio
	}
	return ""
}

type RefreshTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TokenType   string `protobuf:"bytes,1,opt,name=token_type,json=tokenType,proto3" json:"token_type,omitempty"`
	AccessToken string `protobuf:"bytes,2,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
}

func (x *RefreshTokenRequest) Reset() {
	*x = RefreshTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authentication_v1_authentication_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefreshTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshTokenRequest) ProtoMessage() {}

func (x *RefreshTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_authentication_v1_authentication_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshTokenRequest.ProtoReflect.Descriptor instead.
func (*RefreshTokenRequest) Descriptor() ([]byte, []int) {
	return file_authentication_v1_authentication_proto_rawDescGZIP(), []int{4}
}

func (x *RefreshTokenRequest) GetTokenType() string {
	if x != nil {
		return x.TokenType
	}
	return ""
}

func (x *RefreshTokenRequest) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

type RefreshTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TokenType   string               `protobuf:"bytes,1,opt,name=token_type,json=tokenType,proto3" json:"token_type,omitempty"`
	AccessToken string               `protobuf:"bytes,2,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	ExpiryTime  *durationpb.Duration `protobuf:"bytes,3,opt,name=expiry_time,json=expiryTime,proto3" json:"expiry_time,omitempty"`
}

func (x *RefreshTokenResponse) Reset() {
	*x = RefreshTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authentication_v1_authentication_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RefreshTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshTokenResponse) ProtoMessage() {}

func (x *RefreshTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_authentication_v1_authentication_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshTokenResponse.ProtoReflect.Descriptor instead.
func (*RefreshTokenResponse) Descriptor() ([]byte, []int) {
	return file_authentication_v1_authentication_proto_rawDescGZIP(), []int{5}
}

func (x *RefreshTokenResponse) GetTokenType() string {
	if x != nil {
		return x.TokenType
	}
	return ""
}

func (x *RefreshTokenResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *RefreshTokenResponse) GetExpiryTime() *durationpb.Duration {
	if x != nil {
		return x.ExpiryTime
	}
	return nil
}

type GetCredentialsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId *string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
}

func (x *GetCredentialsRequest) Reset() {
	*x = GetCredentialsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authentication_v1_authentication_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCredentialsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCredentialsRequest) ProtoMessage() {}

func (x *GetCredentialsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_authentication_v1_authentication_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCredentialsRequest.ProtoReflect.Descriptor instead.
func (*GetCredentialsRequest) Descriptor() ([]byte, []int) {
	return file_authentication_v1_authentication_proto_rawDescGZIP(), []int{6}
}

func (x *GetCredentialsRequest) GetUserId() string {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return ""
}

type GetCredentialsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Creds []*Credential `protobuf:"bytes,1,rep,name=creds,proto3" json:"creds,omitempty"`
}

func (x *GetCredentialsResponse) Reset() {
	*x = GetCredentialsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authentication_v1_authentication_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCredentialsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCredentialsResponse) ProtoMessage() {}

func (x *GetCredentialsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_authentication_v1_authentication_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCredentialsResponse.ProtoReflect.Descriptor instead.
func (*GetCredentialsResponse) Descriptor() ([]byte, []int) {
	return file_authentication_v1_authentication_proto_rawDescGZIP(), []int{7}
}

func (x *GetCredentialsResponse) GetCreds() []*Credential {
	if x != nil {
		return x.Creds
	}
	return nil
}

type LogoutRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  *string  `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"` // Allow super role to be able delete arbitrary credentials
	CredIds []string `protobuf:"bytes,2,rep,name=cred_ids,json=credIds,proto3" json:"cred_ids,omitempty"`
}

func (x *LogoutRequest) Reset() {
	*x = LogoutRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authentication_v1_authentication_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogoutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogoutRequest) ProtoMessage() {}

func (x *LogoutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_authentication_v1_authentication_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogoutRequest.ProtoReflect.Descriptor instead.
func (*LogoutRequest) Descriptor() ([]byte, []int) {
	return file_authentication_v1_authentication_proto_rawDescGZIP(), []int{8}
}

func (x *LogoutRequest) GetUserId() string {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return ""
}

func (x *LogoutRequest) GetCredIds() []string {
	if x != nil {
		return x.CredIds
	}
	return nil
}

type LogoutAllRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId *string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3,oneof" json:"user_id,omitempty"`
}

func (x *LogoutAllRequest) Reset() {
	*x = LogoutAllRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_authentication_v1_authentication_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogoutAllRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogoutAllRequest) ProtoMessage() {}

func (x *LogoutAllRequest) ProtoReflect() protoreflect.Message {
	mi := &file_authentication_v1_authentication_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogoutAllRequest.ProtoReflect.Descriptor instead.
func (*LogoutAllRequest) Descriptor() ([]byte, []int) {
	return file_authentication_v1_authentication_proto_rawDescGZIP(), []int{9}
}

func (x *LogoutAllRequest) GetUserId() string {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return ""
}

var File_authentication_v1_authentication_proto protoreflect.FileDescriptor

var file_authentication_v1_authentication_proto_rawDesc = []byte{
	0x0a, 0x26, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x61,
	0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31,
	0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x14, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x61, 0x67, 0x65,
	0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x67, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x61, 0x6c, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x22, 0x0a,
	0x0a, 0x69, 0x73, 0x5f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x08, 0x48, 0x00, 0x52, 0x09, 0x69, 0x73, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x88, 0x01,
	0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x69, 0x73, 0x5f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74,
	0x22, 0x61, 0x0a, 0x0c, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x22, 0x8d, 0x01, 0x0a, 0x0d, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x3a, 0x0a, 0x0b, 0x65, 0x78, 0x70, 0x69, 0x72,
	0x79, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x54,
	0x69, 0x6d, 0x65, 0x22, 0xcd, 0x01, 0x0a, 0x0f, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x72, 0x73, 0x74, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69, 0x72, 0x73, 0x74,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x4e,
	0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x15, 0x0a, 0x03, 0x62, 0x69, 0x6f, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x03, 0x62, 0x69, 0x6f, 0x88, 0x01, 0x01, 0x42, 0x0c, 0x0a,
	0x0a, 0x5f, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x5f,
	0x62, 0x69, 0x6f, 0x22, 0x57, 0x0a, 0x13, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x94, 0x01, 0x0a,
	0x14, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x3a, 0x0a, 0x0b, 0x65, 0x78, 0x70, 0x69, 0x72,
	0x79, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x65, 0x78, 0x70, 0x69, 0x72, 0x79, 0x54,
	0x69, 0x6d, 0x65, 0x22, 0x41, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x61, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x22, 0x52, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x43, 0x72, 0x65,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x38, 0x0a, 0x05, 0x63, 0x72, 0x65, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x22, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x61, 0x6c, 0x52, 0x05, 0x63, 0x72, 0x65, 0x64, 0x73, 0x22, 0x54, 0x0a, 0x0d, 0x4c, 0x6f,
	0x67, 0x6f, 0x75, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x08, 0x63, 0x72, 0x65,
	0x64, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x63, 0x72, 0x65,
	0x64, 0x49, 0x64, 0x73, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x22, 0x3c, 0x0a, 0x10, 0x4c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x41, 0x6c, 0x6c, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x88,
	0x01, 0x01, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x32, 0xae,
	0x04, 0x0a, 0x15, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x54, 0x0a, 0x05, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x12, 0x24, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x25, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x61,
	0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31,
	0x2e, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4b,
	0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x12, 0x27, 0x2e, 0x6e, 0x65, 0x78,
	0x61, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x69, 0x0a, 0x0c, 0x52,
	0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x2b, 0x2e, 0x6e, 0x65,
	0x78, 0x61, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e,
	0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x6f, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x43, 0x72, 0x65,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x12, 0x2d, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e,
	0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x61,
	0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x65, 0x74, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x47, 0x0a, 0x06, 0x4c, 0x6f, 0x67, 0x6f, 0x75,
	0x74, 0x12, 0x25, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x6f, 0x75,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x12, 0x4d, 0x0a, 0x09, 0x4c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x41, 0x6c, 0x6c, 0x12, 0x28, 0x2e,
	0x6e, 0x65, 0x78, 0x61, 0x2e, 0x61, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x6f, 0x67, 0x6f, 0x75, 0x74, 0x41, 0x6c, 0x6c,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42,
	0xf5, 0x01, 0x0a, 0x1a, 0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x61, 0x75, 0x74,
	0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x42, 0x13,
	0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x48, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x61, 0x72, 0x63, 0x6f, 0x72, 0x69, 0x75, 0x6d, 0x2f, 0x6e, 0x65, 0x78, 0x61, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x61, 0x75, 0x74,
	0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x3b, 0x61,
	0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x76, 0x31, 0xa2,
	0x02, 0x03, 0x4e, 0x41, 0x58, 0xaa, 0x02, 0x16, 0x4e, 0x65, 0x78, 0x61, 0x2e, 0x41, 0x75, 0x74,
	0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x56, 0x31, 0xca, 0x02,
	0x16, 0x4e, 0x65, 0x78, 0x61, 0x5c, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x22, 0x4e, 0x65, 0x78, 0x61, 0x5c, 0x41,
	0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5c, 0x56, 0x31,
	0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x18, 0x4e,
	0x65, 0x78, 0x61, 0x3a, 0x3a, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e, 0x74, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_authentication_v1_authentication_proto_rawDescOnce sync.Once
	file_authentication_v1_authentication_proto_rawDescData = file_authentication_v1_authentication_proto_rawDesc
)

func file_authentication_v1_authentication_proto_rawDescGZIP() []byte {
	file_authentication_v1_authentication_proto_rawDescOnce.Do(func() {
		file_authentication_v1_authentication_proto_rawDescData = protoimpl.X.CompressGZIP(file_authentication_v1_authentication_proto_rawDescData)
	})
	return file_authentication_v1_authentication_proto_rawDescData
}

var file_authentication_v1_authentication_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_authentication_v1_authentication_proto_goTypes = []interface{}{
	(*Credential)(nil),             // 0: nexa.authentication.v1.Credential
	(*LoginRequest)(nil),           // 1: nexa.authentication.v1.LoginRequest
	(*LoginResponse)(nil),          // 2: nexa.authentication.v1.LoginResponse
	(*RegisterRequest)(nil),        // 3: nexa.authentication.v1.RegisterRequest
	(*RefreshTokenRequest)(nil),    // 4: nexa.authentication.v1.RefreshTokenRequest
	(*RefreshTokenResponse)(nil),   // 5: nexa.authentication.v1.RefreshTokenResponse
	(*GetCredentialsRequest)(nil),  // 6: nexa.authentication.v1.GetCredentialsRequest
	(*GetCredentialsResponse)(nil), // 7: nexa.authentication.v1.GetCredentialsResponse
	(*LogoutRequest)(nil),          // 8: nexa.authentication.v1.LogoutRequest
	(*LogoutAllRequest)(nil),       // 9: nexa.authentication.v1.LogoutAllRequest
	(*durationpb.Duration)(nil),    // 10: google.protobuf.Duration
	(*emptypb.Empty)(nil),          // 11: google.protobuf.Empty
}
var file_authentication_v1_authentication_proto_depIdxs = []int32{
	10, // 0: nexa.authentication.v1.LoginResponse.expiry_time:type_name -> google.protobuf.Duration
	10, // 1: nexa.authentication.v1.RefreshTokenResponse.expiry_time:type_name -> google.protobuf.Duration
	0,  // 2: nexa.authentication.v1.GetCredentialsResponse.creds:type_name -> nexa.authentication.v1.Credential
	1,  // 3: nexa.authentication.v1.AuthenticationService.Login:input_type -> nexa.authentication.v1.LoginRequest
	3,  // 4: nexa.authentication.v1.AuthenticationService.Register:input_type -> nexa.authentication.v1.RegisterRequest
	4,  // 5: nexa.authentication.v1.AuthenticationService.RefreshToken:input_type -> nexa.authentication.v1.RefreshTokenRequest
	6,  // 6: nexa.authentication.v1.AuthenticationService.GetCredentials:input_type -> nexa.authentication.v1.GetCredentialsRequest
	8,  // 7: nexa.authentication.v1.AuthenticationService.Logout:input_type -> nexa.authentication.v1.LogoutRequest
	9,  // 8: nexa.authentication.v1.AuthenticationService.LogoutAll:input_type -> nexa.authentication.v1.LogoutAllRequest
	2,  // 9: nexa.authentication.v1.AuthenticationService.Login:output_type -> nexa.authentication.v1.LoginResponse
	11, // 10: nexa.authentication.v1.AuthenticationService.Register:output_type -> google.protobuf.Empty
	5,  // 11: nexa.authentication.v1.AuthenticationService.RefreshToken:output_type -> nexa.authentication.v1.RefreshTokenResponse
	7,  // 12: nexa.authentication.v1.AuthenticationService.GetCredentials:output_type -> nexa.authentication.v1.GetCredentialsResponse
	11, // 13: nexa.authentication.v1.AuthenticationService.Logout:output_type -> google.protobuf.Empty
	11, // 14: nexa.authentication.v1.AuthenticationService.LogoutAll:output_type -> google.protobuf.Empty
	9,  // [9:15] is the sub-list for method output_type
	3,  // [3:9] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_authentication_v1_authentication_proto_init() }
func file_authentication_v1_authentication_proto_init() {
	if File_authentication_v1_authentication_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_authentication_v1_authentication_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Credential); i {
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
		file_authentication_v1_authentication_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginRequest); i {
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
		file_authentication_v1_authentication_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginResponse); i {
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
		file_authentication_v1_authentication_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterRequest); i {
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
		file_authentication_v1_authentication_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RefreshTokenRequest); i {
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
		file_authentication_v1_authentication_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RefreshTokenResponse); i {
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
		file_authentication_v1_authentication_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCredentialsRequest); i {
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
		file_authentication_v1_authentication_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCredentialsResponse); i {
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
		file_authentication_v1_authentication_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogoutRequest); i {
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
		file_authentication_v1_authentication_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogoutAllRequest); i {
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
	file_authentication_v1_authentication_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_authentication_v1_authentication_proto_msgTypes[3].OneofWrappers = []interface{}{}
	file_authentication_v1_authentication_proto_msgTypes[6].OneofWrappers = []interface{}{}
	file_authentication_v1_authentication_proto_msgTypes[8].OneofWrappers = []interface{}{}
	file_authentication_v1_authentication_proto_msgTypes[9].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_authentication_v1_authentication_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_authentication_v1_authentication_proto_goTypes,
		DependencyIndexes: file_authentication_v1_authentication_proto_depIdxs,
		MessageInfos:      file_authentication_v1_authentication_proto_msgTypes,
	}.Build()
	File_authentication_v1_authentication_proto = out.File
	file_authentication_v1_authentication_proto_rawDesc = nil
	file_authentication_v1_authentication_proto_goTypes = nil
	file_authentication_v1_authentication_proto_depIdxs = nil
}
