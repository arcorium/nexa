// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: reaction/v1/reaction.proto

package reactionv1

import (
	common "github.com/arcorium/nexa/proto/gen/go/common"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Type int32

const (
	Type_POST    Type = 0
	Type_COMMENT Type = 1
)

// Enum value maps for Type.
var (
	Type_name = map[int32]string{
		0: "POST",
		1: "COMMENT",
	}
	Type_value = map[string]int32{
		"POST":    0,
		"COMMENT": 1,
	}
)

func (x Type) Enum() *Type {
	p := new(Type)
	*p = x
	return p
}

func (x Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Type) Descriptor() protoreflect.EnumDescriptor {
	return file_reaction_v1_reaction_proto_enumTypes[0].Descriptor()
}

func (Type) Type() protoreflect.EnumType {
	return &file_reaction_v1_reaction_proto_enumTypes[0]
}

func (x Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Type.Descriptor instead.
func (Type) EnumDescriptor() ([]byte, []int) {
	return file_reaction_v1_reaction_proto_rawDescGZIP(), []int{0}
}

type ReactionType int32

const (
	ReactionType_LIKE    ReactionType = 0
	ReactionType_DISLIKE ReactionType = 1
)

// Enum value maps for ReactionType.
var (
	ReactionType_name = map[int32]string{
		0: "LIKE",
		1: "DISLIKE",
	}
	ReactionType_value = map[string]int32{
		"LIKE":    0,
		"DISLIKE": 1,
	}
)

func (x ReactionType) Enum() *ReactionType {
	p := new(ReactionType)
	*p = x
	return p
}

func (x ReactionType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ReactionType) Descriptor() protoreflect.EnumDescriptor {
	return file_reaction_v1_reaction_proto_enumTypes[1].Descriptor()
}

func (ReactionType) Type() protoreflect.EnumType {
	return &file_reaction_v1_reaction_proto_enumTypes[1]
}

func (x ReactionType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ReactionType.Descriptor instead.
func (ReactionType) EnumDescriptor() ([]byte, []int) {
	return file_reaction_v1_reaction_proto_rawDescGZIP(), []int{1}
}

type Reaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId       string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ReactionType ReactionType           `protobuf:"varint,2,opt,name=reaction_type,json=reactionType,proto3,enum=nexa.reaction.v1.ReactionType" json:"reaction_type,omitempty"`
	CreatedAt    *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (x *Reaction) Reset() {
	*x = Reaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reaction_v1_reaction_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Reaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Reaction) ProtoMessage() {}

func (x *Reaction) ProtoReflect() protoreflect.Message {
	mi := &file_reaction_v1_reaction_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Reaction.ProtoReflect.Descriptor instead.
func (*Reaction) Descriptor() ([]byte, []int) {
	return file_reaction_v1_reaction_proto_rawDescGZIP(), []int{0}
}

func (x *Reaction) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Reaction) GetReactionType() ReactionType {
	if x != nil {
		return x.ReactionType
	}
	return ReactionType_LIKE
}

func (x *Reaction) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

type LikeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ItemType Type   `protobuf:"varint,1,opt,name=item_type,json=itemType,proto3,enum=nexa.reaction.v1.Type" json:"item_type,omitempty"`
	ItemId   string `protobuf:"bytes,2,opt,name=item_id,json=itemId,proto3" json:"item_id,omitempty"`
}

func (x *LikeRequest) Reset() {
	*x = LikeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reaction_v1_reaction_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LikeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LikeRequest) ProtoMessage() {}

func (x *LikeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_reaction_v1_reaction_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LikeRequest.ProtoReflect.Descriptor instead.
func (*LikeRequest) Descriptor() ([]byte, []int) {
	return file_reaction_v1_reaction_proto_rawDescGZIP(), []int{1}
}

func (x *LikeRequest) GetItemType() Type {
	if x != nil {
		return x.ItemType
	}
	return Type_POST
}

func (x *LikeRequest) GetItemId() string {
	if x != nil {
		return x.ItemId
	}
	return ""
}

type DislikeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ItemType Type   `protobuf:"varint,1,opt,name=item_type,json=itemType,proto3,enum=nexa.reaction.v1.Type" json:"item_type,omitempty"`
	ItemId   string `protobuf:"bytes,2,opt,name=item_id,json=itemId,proto3" json:"item_id,omitempty"`
}

func (x *DislikeRequest) Reset() {
	*x = DislikeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reaction_v1_reaction_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DislikeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DislikeRequest) ProtoMessage() {}

func (x *DislikeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_reaction_v1_reaction_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DislikeRequest.ProtoReflect.Descriptor instead.
func (*DislikeRequest) Descriptor() ([]byte, []int) {
	return file_reaction_v1_reaction_proto_rawDescGZIP(), []int{2}
}

func (x *DislikeRequest) GetItemType() Type {
	if x != nil {
		return x.ItemType
	}
	return Type_POST
}

func (x *DislikeRequest) GetItemId() string {
	if x != nil {
		return x.ItemId
	}
	return ""
}

type GetItemReactionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ItemType Type                      `protobuf:"varint,1,opt,name=item_type,json=itemType,proto3,enum=nexa.reaction.v1.Type" json:"item_type,omitempty"`
	ItemId   string                    `protobuf:"bytes,2,opt,name=item_id,json=itemId,proto3" json:"item_id,omitempty"`
	Details  *common.PagedElementInput `protobuf:"bytes,3,opt,name=details,proto3" json:"details,omitempty"`
}

func (x *GetItemReactionsRequest) Reset() {
	*x = GetItemReactionsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reaction_v1_reaction_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetItemReactionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetItemReactionsRequest) ProtoMessage() {}

func (x *GetItemReactionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_reaction_v1_reaction_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetItemReactionsRequest.ProtoReflect.Descriptor instead.
func (*GetItemReactionsRequest) Descriptor() ([]byte, []int) {
	return file_reaction_v1_reaction_proto_rawDescGZIP(), []int{3}
}

func (x *GetItemReactionsRequest) GetItemType() Type {
	if x != nil {
		return x.ItemType
	}
	return Type_POST
}

func (x *GetItemReactionsRequest) GetItemId() string {
	if x != nil {
		return x.ItemId
	}
	return ""
}

func (x *GetItemReactionsRequest) GetDetails() *common.PagedElementInput {
	if x != nil {
		return x.Details
	}
	return nil
}

type GetItemReactionsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ItemType  Type                       `protobuf:"varint,1,opt,name=item_type,json=itemType,proto3,enum=nexa.reaction.v1.Type" json:"item_type,omitempty"`
	ItemId    string                     `protobuf:"bytes,2,opt,name=item_id,json=itemId,proto3" json:"item_id,omitempty"`
	Reactions []*Reaction                `protobuf:"bytes,3,rep,name=reactions,proto3" json:"reactions,omitempty"`
	Details   *common.PagedElementOutput `protobuf:"bytes,4,opt,name=details,proto3" json:"details,omitempty"`
}

func (x *GetItemReactionsResponse) Reset() {
	*x = GetItemReactionsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reaction_v1_reaction_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetItemReactionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetItemReactionsResponse) ProtoMessage() {}

func (x *GetItemReactionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_reaction_v1_reaction_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetItemReactionsResponse.ProtoReflect.Descriptor instead.
func (*GetItemReactionsResponse) Descriptor() ([]byte, []int) {
	return file_reaction_v1_reaction_proto_rawDescGZIP(), []int{4}
}

func (x *GetItemReactionsResponse) GetItemType() Type {
	if x != nil {
		return x.ItemType
	}
	return Type_POST
}

func (x *GetItemReactionsResponse) GetItemId() string {
	if x != nil {
		return x.ItemId
	}
	return ""
}

func (x *GetItemReactionsResponse) GetReactions() []*Reaction {
	if x != nil {
		return x.Reactions
	}
	return nil
}

func (x *GetItemReactionsResponse) GetDetails() *common.PagedElementOutput {
	if x != nil {
		return x.Details
	}
	return nil
}

type GetCountRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ItemType Type     `protobuf:"varint,1,opt,name=item_type,json=itemType,proto3,enum=nexa.reaction.v1.Type" json:"item_type,omitempty"`
	ItemIds  []string `protobuf:"bytes,2,rep,name=item_ids,json=itemIds,proto3" json:"item_ids,omitempty"`
}

func (x *GetCountRequest) Reset() {
	*x = GetCountRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reaction_v1_reaction_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCountRequest) ProtoMessage() {}

func (x *GetCountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_reaction_v1_reaction_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCountRequest.ProtoReflect.Descriptor instead.
func (*GetCountRequest) Descriptor() ([]byte, []int) {
	return file_reaction_v1_reaction_proto_rawDescGZIP(), []int{5}
}

func (x *GetCountRequest) GetItemType() Type {
	if x != nil {
		return x.ItemType
	}
	return Type_POST
}

func (x *GetCountRequest) GetItemIds() []string {
	if x != nil {
		return x.ItemIds
	}
	return nil
}

type Count struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TotalLikes    uint64 `protobuf:"varint,1,opt,name=total_likes,json=totalLikes,proto3" json:"total_likes,omitempty"`
	TotalDislikes uint64 `protobuf:"varint,2,opt,name=total_dislikes,json=totalDislikes,proto3" json:"total_dislikes,omitempty"`
}

func (x *Count) Reset() {
	*x = Count{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reaction_v1_reaction_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Count) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Count) ProtoMessage() {}

func (x *Count) ProtoReflect() protoreflect.Message {
	mi := &file_reaction_v1_reaction_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Count.ProtoReflect.Descriptor instead.
func (*Count) Descriptor() ([]byte, []int) {
	return file_reaction_v1_reaction_proto_rawDescGZIP(), []int{6}
}

func (x *Count) GetTotalLikes() uint64 {
	if x != nil {
		return x.TotalLikes
	}
	return 0
}

func (x *Count) GetTotalDislikes() uint64 {
	if x != nil {
		return x.TotalDislikes
	}
	return 0
}

type GetCountResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Counts []*Count `protobuf:"bytes,1,rep,name=counts,proto3" json:"counts,omitempty"`
}

func (x *GetCountResponse) Reset() {
	*x = GetCountResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reaction_v1_reaction_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetCountResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetCountResponse) ProtoMessage() {}

func (x *GetCountResponse) ProtoReflect() protoreflect.Message {
	mi := &file_reaction_v1_reaction_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetCountResponse.ProtoReflect.Descriptor instead.
func (*GetCountResponse) Descriptor() ([]byte, []int) {
	return file_reaction_v1_reaction_proto_rawDescGZIP(), []int{7}
}

func (x *GetCountResponse) GetCounts() []*Count {
	if x != nil {
		return x.Counts
	}
	return nil
}

type ClearUsersReactionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *ClearUsersReactionsRequest) Reset() {
	*x = ClearUsersReactionsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reaction_v1_reaction_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClearUsersReactionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClearUsersReactionsRequest) ProtoMessage() {}

func (x *ClearUsersReactionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_reaction_v1_reaction_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClearUsersReactionsRequest.ProtoReflect.Descriptor instead.
func (*ClearUsersReactionsRequest) Descriptor() ([]byte, []int) {
	return file_reaction_v1_reaction_proto_rawDescGZIP(), []int{8}
}

func (x *ClearUsersReactionsRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type DeleteItemsReactionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ItemType Type     `protobuf:"varint,1,opt,name=item_type,json=itemType,proto3,enum=nexa.reaction.v1.Type" json:"item_type,omitempty"`
	ItemIds  []string `protobuf:"bytes,2,rep,name=item_ids,json=itemIds,proto3" json:"item_ids,omitempty"`
}

func (x *DeleteItemsReactionsRequest) Reset() {
	*x = DeleteItemsReactionsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_reaction_v1_reaction_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteItemsReactionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteItemsReactionsRequest) ProtoMessage() {}

func (x *DeleteItemsReactionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_reaction_v1_reaction_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteItemsReactionsRequest.ProtoReflect.Descriptor instead.
func (*DeleteItemsReactionsRequest) Descriptor() ([]byte, []int) {
	return file_reaction_v1_reaction_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteItemsReactionsRequest) GetItemType() Type {
	if x != nil {
		return x.ItemType
	}
	return Type_POST
}

func (x *DeleteItemsReactionsRequest) GetItemIds() []string {
	if x != nil {
		return x.ItemIds
	}
	return nil
}

var File_reaction_v1_reaction_proto protoreflect.FileDescriptor

var file_reaction_v1_reaction_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x6e, 0x65,
	0x78, 0x61, 0x2e, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x70, 0x61, 0x67, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xa3, 0x01, 0x0a, 0x08, 0x52, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x43, 0x0a, 0x0d, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e,
	0x6e, 0x65, 0x78, 0x61, 0x2e, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31,
	0x2e, 0x52, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0c, 0x72,
	0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x22, 0x5b, 0x0a, 0x0b, 0x4c, 0x69, 0x6b, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x09, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e,
	0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x08, 0x69, 0x74, 0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x74,
	0x65, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x69, 0x74, 0x65,
	0x6d, 0x49, 0x64, 0x22, 0x5e, 0x0a, 0x0e, 0x44, 0x69, 0x73, 0x6c, 0x69, 0x6b, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x09, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e,
	0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x79, 0x70, 0x65,
	0x52, 0x08, 0x69, 0x74, 0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x74,
	0x65, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x69, 0x74, 0x65,
	0x6d, 0x49, 0x64, 0x22, 0xa1, 0x01, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52,
	0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x33, 0x0a, 0x09, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x16, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x69, 0x74, 0x65, 0x6d,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x12, 0x38, 0x0a,
	0x07, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x50, 0x61, 0x67,
	0x65, 0x64, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x52, 0x07,
	0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x22, 0xdd, 0x01, 0x0a, 0x18, 0x47, 0x65, 0x74, 0x49,
	0x74, 0x65, 0x6d, 0x52, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x33, 0x0a, 0x09, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x72,
	0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52,
	0x08, 0x69, 0x74, 0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x74, 0x65,
	0x6d, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x69, 0x74, 0x65, 0x6d,
	0x49, 0x64, 0x12, 0x38, 0x0a, 0x09, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x72, 0x65, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x09, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x39, 0x0a, 0x07,
	0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e,
	0x6e, 0x65, 0x78, 0x61, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x50, 0x61, 0x67, 0x65,
	0x64, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x52, 0x07,
	0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x22, 0x61, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x09, 0x69, 0x74,
	0x65, 0x6d, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e,
	0x6e, 0x65, 0x78, 0x61, 0x2e, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31,
	0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x69, 0x74, 0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x19, 0x0a, 0x08, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x07, 0x69, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x73, 0x22, 0x4f, 0x0a, 0x05, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x6c, 0x69, 0x6b,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x4c,
	0x69, 0x6b, 0x65, 0x73, 0x12, 0x25, 0x0a, 0x0e, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x64, 0x69,
	0x73, 0x6c, 0x69, 0x6b, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x44, 0x69, 0x73, 0x6c, 0x69, 0x6b, 0x65, 0x73, 0x22, 0x43, 0x0a, 0x10, 0x47,
	0x65, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2f, 0x0a, 0x06, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x17, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x06, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73,
	0x22, 0x35, 0x0a, 0x1a, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x55, 0x73, 0x65, 0x72, 0x73, 0x52, 0x65,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x6d, 0x0a, 0x1b, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x09, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x16, 0x2e, 0x6e, 0x65, 0x78, 0x61,
	0x2e, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x08, 0x69, 0x74, 0x65, 0x6d, 0x54, 0x79, 0x70, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x69,
	0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x69,
	0x74, 0x65, 0x6d, 0x49, 0x64, 0x73, 0x2a, 0x1d, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08,
	0x0a, 0x04, 0x50, 0x4f, 0x53, 0x54, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f, 0x4d, 0x4d,
	0x45, 0x4e, 0x54, 0x10, 0x01, 0x2a, 0x25, 0x0a, 0x0c, 0x52, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x4c, 0x49, 0x4b, 0x45, 0x10, 0x00, 0x12,
	0x0b, 0x0a, 0x07, 0x44, 0x49, 0x53, 0x4c, 0x49, 0x4b, 0x45, 0x10, 0x01, 0x32, 0xf5, 0x03, 0x0a,
	0x0f, 0x52, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x3d, 0x0a, 0x04, 0x4c, 0x69, 0x6b, 0x65, 0x12, 0x1d, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e,
	0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x6b, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12,
	0x43, 0x0a, 0x07, 0x44, 0x69, 0x73, 0x6c, 0x69, 0x6b, 0x65, 0x12, 0x20, 0x2e, 0x6e, 0x65, 0x78,
	0x61, 0x2e, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x69,
	0x73, 0x6c, 0x69, 0x6b, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x12, 0x61, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x73,
	0x12, 0x29, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2a, 0x2e, 0x6e, 0x65,
	0x78, 0x61, 0x2e, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x47,
	0x65, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x51, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x21, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x72, 0x65, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x72, 0x65,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x54, 0x0a, 0x0b, 0x44, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x2d, 0x2e, 0x6e, 0x65, 0x78, 0x61,
	0x2e, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x73, 0x52, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x12, 0x52, 0x0a, 0x0a, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x55, 0x73, 0x65, 0x72, 0x73, 0x12, 0x2c,
	0x2e, 0x6e, 0x65, 0x78, 0x61, 0x2e, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x55, 0x73, 0x65, 0x72, 0x73, 0x52, 0x65, 0x61, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x42, 0xc5, 0x01, 0x0a, 0x14, 0x63, 0x6f, 0x6d, 0x2e, 0x6e, 0x65, 0x78,
	0x61, 0x2e, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x42, 0x0d, 0x52,
	0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3c,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x72, 0x63, 0x6f, 0x72,
	0x69, 0x75, 0x6d, 0x2f, 0x6e, 0x65, 0x78, 0x61, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67,
	0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76,
	0x31, 0x3b, 0x72, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x4e,
	0x52, 0x58, 0xaa, 0x02, 0x10, 0x4e, 0x65, 0x78, 0x61, 0x2e, 0x52, 0x65, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x10, 0x4e, 0x65, 0x78, 0x61, 0x5c, 0x52, 0x65, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x1c, 0x4e, 0x65, 0x78, 0x61, 0x5c,
	0x52, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x12, 0x4e, 0x65, 0x78, 0x61, 0x3a, 0x3a,
	0x52, 0x65, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_reaction_v1_reaction_proto_rawDescOnce sync.Once
	file_reaction_v1_reaction_proto_rawDescData = file_reaction_v1_reaction_proto_rawDesc
)

func file_reaction_v1_reaction_proto_rawDescGZIP() []byte {
	file_reaction_v1_reaction_proto_rawDescOnce.Do(func() {
		file_reaction_v1_reaction_proto_rawDescData = protoimpl.X.CompressGZIP(file_reaction_v1_reaction_proto_rawDescData)
	})
	return file_reaction_v1_reaction_proto_rawDescData
}

var file_reaction_v1_reaction_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_reaction_v1_reaction_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_reaction_v1_reaction_proto_goTypes = []interface{}{
	(Type)(0),                           // 0: nexa.reaction.v1.Type
	(ReactionType)(0),                   // 1: nexa.reaction.v1.ReactionType
	(*Reaction)(nil),                    // 2: nexa.reaction.v1.Reaction
	(*LikeRequest)(nil),                 // 3: nexa.reaction.v1.LikeRequest
	(*DislikeRequest)(nil),              // 4: nexa.reaction.v1.DislikeRequest
	(*GetItemReactionsRequest)(nil),     // 5: nexa.reaction.v1.GetItemReactionsRequest
	(*GetItemReactionsResponse)(nil),    // 6: nexa.reaction.v1.GetItemReactionsResponse
	(*GetCountRequest)(nil),             // 7: nexa.reaction.v1.GetCountRequest
	(*Count)(nil),                       // 8: nexa.reaction.v1.Count
	(*GetCountResponse)(nil),            // 9: nexa.reaction.v1.GetCountResponse
	(*ClearUsersReactionsRequest)(nil),  // 10: nexa.reaction.v1.ClearUsersReactionsRequest
	(*DeleteItemsReactionsRequest)(nil), // 11: nexa.reaction.v1.DeleteItemsReactionsRequest
	(*timestamppb.Timestamp)(nil),       // 12: google.protobuf.Timestamp
	(*common.PagedElementInput)(nil),    // 13: nexa.common.PagedElementInput
	(*common.PagedElementOutput)(nil),   // 14: nexa.common.PagedElementOutput
	(*emptypb.Empty)(nil),               // 15: google.protobuf.Empty
}
var file_reaction_v1_reaction_proto_depIdxs = []int32{
	1,  // 0: nexa.reaction.v1.Reaction.reaction_type:type_name -> nexa.reaction.v1.ReactionType
	12, // 1: nexa.reaction.v1.Reaction.created_at:type_name -> google.protobuf.Timestamp
	0,  // 2: nexa.reaction.v1.LikeRequest.item_type:type_name -> nexa.reaction.v1.Type
	0,  // 3: nexa.reaction.v1.DislikeRequest.item_type:type_name -> nexa.reaction.v1.Type
	0,  // 4: nexa.reaction.v1.GetItemReactionsRequest.item_type:type_name -> nexa.reaction.v1.Type
	13, // 5: nexa.reaction.v1.GetItemReactionsRequest.details:type_name -> nexa.common.PagedElementInput
	0,  // 6: nexa.reaction.v1.GetItemReactionsResponse.item_type:type_name -> nexa.reaction.v1.Type
	2,  // 7: nexa.reaction.v1.GetItemReactionsResponse.reactions:type_name -> nexa.reaction.v1.Reaction
	14, // 8: nexa.reaction.v1.GetItemReactionsResponse.details:type_name -> nexa.common.PagedElementOutput
	0,  // 9: nexa.reaction.v1.GetCountRequest.item_type:type_name -> nexa.reaction.v1.Type
	8,  // 10: nexa.reaction.v1.GetCountResponse.counts:type_name -> nexa.reaction.v1.Count
	0,  // 11: nexa.reaction.v1.DeleteItemsReactionsRequest.item_type:type_name -> nexa.reaction.v1.Type
	3,  // 12: nexa.reaction.v1.ReactionService.Like:input_type -> nexa.reaction.v1.LikeRequest
	4,  // 13: nexa.reaction.v1.ReactionService.Dislike:input_type -> nexa.reaction.v1.DislikeRequest
	5,  // 14: nexa.reaction.v1.ReactionService.GetItems:input_type -> nexa.reaction.v1.GetItemReactionsRequest
	7,  // 15: nexa.reaction.v1.ReactionService.GetCount:input_type -> nexa.reaction.v1.GetCountRequest
	11, // 16: nexa.reaction.v1.ReactionService.DeleteItems:input_type -> nexa.reaction.v1.DeleteItemsReactionsRequest
	10, // 17: nexa.reaction.v1.ReactionService.ClearUsers:input_type -> nexa.reaction.v1.ClearUsersReactionsRequest
	15, // 18: nexa.reaction.v1.ReactionService.Like:output_type -> google.protobuf.Empty
	15, // 19: nexa.reaction.v1.ReactionService.Dislike:output_type -> google.protobuf.Empty
	6,  // 20: nexa.reaction.v1.ReactionService.GetItems:output_type -> nexa.reaction.v1.GetItemReactionsResponse
	9,  // 21: nexa.reaction.v1.ReactionService.GetCount:output_type -> nexa.reaction.v1.GetCountResponse
	15, // 22: nexa.reaction.v1.ReactionService.DeleteItems:output_type -> google.protobuf.Empty
	15, // 23: nexa.reaction.v1.ReactionService.ClearUsers:output_type -> google.protobuf.Empty
	18, // [18:24] is the sub-list for method output_type
	12, // [12:18] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_reaction_v1_reaction_proto_init() }
func file_reaction_v1_reaction_proto_init() {
	if File_reaction_v1_reaction_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_reaction_v1_reaction_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Reaction); i {
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
		file_reaction_v1_reaction_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LikeRequest); i {
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
		file_reaction_v1_reaction_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DislikeRequest); i {
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
		file_reaction_v1_reaction_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetItemReactionsRequest); i {
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
		file_reaction_v1_reaction_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetItemReactionsResponse); i {
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
		file_reaction_v1_reaction_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCountRequest); i {
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
		file_reaction_v1_reaction_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Count); i {
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
		file_reaction_v1_reaction_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetCountResponse); i {
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
		file_reaction_v1_reaction_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClearUsersReactionsRequest); i {
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
		file_reaction_v1_reaction_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteItemsReactionsRequest); i {
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
			RawDescriptor: file_reaction_v1_reaction_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_reaction_v1_reaction_proto_goTypes,
		DependencyIndexes: file_reaction_v1_reaction_proto_depIdxs,
		EnumInfos:         file_reaction_v1_reaction_proto_enumTypes,
		MessageInfos:      file_reaction_v1_reaction_proto_msgTypes,
	}.Build()
	File_reaction_v1_reaction_proto = out.File
	file_reaction_v1_reaction_proto_rawDesc = nil
	file_reaction_v1_reaction_proto_goTypes = nil
	file_reaction_v1_reaction_proto_depIdxs = nil
}