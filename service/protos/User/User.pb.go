// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: User.proto

package User

import (
	Common "service/protos/Common"
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

type AddUserChildRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Model *UserChildModel `protobuf:"bytes,1,opt,name=model,proto3" json:"model,omitempty"`
}

func (x *AddUserChildRequest) Reset() {
	*x = AddUserChildRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_User_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddUserChildRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddUserChildRequest) ProtoMessage() {}

func (x *AddUserChildRequest) ProtoReflect() protoreflect.Message {
	mi := &file_User_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddUserChildRequest.ProtoReflect.Descriptor instead.
func (*AddUserChildRequest) Descriptor() ([]byte, []int) {
	return file_User_proto_rawDescGZIP(), []int{0}
}

func (x *AddUserChildRequest) GetModel() *UserChildModel {
	if x != nil {
		return x.Model
	}
	return nil
}

type GetConfirmUserChildInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UrlKey string `protobuf:"bytes,1,opt,name=url_key,json=urlKey,proto3" json:"url_key,omitempty"`
}

func (x *GetConfirmUserChildInfoRequest) Reset() {
	*x = GetConfirmUserChildInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_User_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetConfirmUserChildInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetConfirmUserChildInfoRequest) ProtoMessage() {}

func (x *GetConfirmUserChildInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_User_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetConfirmUserChildInfoRequest.ProtoReflect.Descriptor instead.
func (*GetConfirmUserChildInfoRequest) Descriptor() ([]byte, []int) {
	return file_User_proto_rawDescGZIP(), []int{1}
}

func (x *GetConfirmUserChildInfoRequest) GetUrlKey() string {
	if x != nil {
		return x.UrlKey
	}
	return ""
}

type GetConfirmUserChildInfoReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Account string `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
}

func (x *GetConfirmUserChildInfoReply) Reset() {
	*x = GetConfirmUserChildInfoReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_User_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetConfirmUserChildInfoReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetConfirmUserChildInfoReply) ProtoMessage() {}

func (x *GetConfirmUserChildInfoReply) ProtoReflect() protoreflect.Message {
	mi := &file_User_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetConfirmUserChildInfoReply.ProtoReflect.Descriptor instead.
func (*GetConfirmUserChildInfoReply) Descriptor() ([]byte, []int) {
	return file_User_proto_rawDescGZIP(), []int{2}
}

func (x *GetConfirmUserChildInfoReply) GetAccount() string {
	if x != nil {
		return x.Account
	}
	return ""
}

type ConfirmUserChildRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UrlKey   string `protobuf:"bytes,1,opt,name=url_key,json=urlKey,proto3" json:"url_key,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *ConfirmUserChildRequest) Reset() {
	*x = ConfirmUserChildRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_User_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConfirmUserChildRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConfirmUserChildRequest) ProtoMessage() {}

func (x *ConfirmUserChildRequest) ProtoReflect() protoreflect.Message {
	mi := &file_User_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConfirmUserChildRequest.ProtoReflect.Descriptor instead.
func (*ConfirmUserChildRequest) Descriptor() ([]byte, []int) {
	return file_User_proto_rawDescGZIP(), []int{3}
}

func (x *ConfirmUserChildRequest) GetUrlKey() string {
	if x != nil {
		return x.UrlKey
	}
	return ""
}

func (x *ConfirmUserChildRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type ResendChildInviteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *ResendChildInviteRequest) Reset() {
	*x = ResendChildInviteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_User_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResendChildInviteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResendChildInviteRequest) ProtoMessage() {}

func (x *ResendChildInviteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_User_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResendChildInviteRequest.ProtoReflect.Descriptor instead.
func (*ResendChildInviteRequest) Descriptor() ([]byte, []int) {
	return file_User_proto_rawDescGZIP(), []int{4}
}

func (x *ResendChildInviteRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type ResendChildInviteReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UrlKey string `protobuf:"bytes,1,opt,name=url_key,json=urlKey,proto3" json:"url_key,omitempty"` // for test
}

func (x *ResendChildInviteReply) Reset() {
	*x = ResendChildInviteReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_User_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResendChildInviteReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResendChildInviteReply) ProtoMessage() {}

func (x *ResendChildInviteReply) ProtoReflect() protoreflect.Message {
	mi := &file_User_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResendChildInviteReply.ProtoReflect.Descriptor instead.
func (*ResendChildInviteReply) Descriptor() ([]byte, []int) {
	return file_User_proto_rawDescGZIP(), []int{5}
}

func (x *ResendChildInviteReply) GetUrlKey() string {
	if x != nil {
		return x.UrlKey
	}
	return ""
}

type CancelChildInviteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *CancelChildInviteRequest) Reset() {
	*x = CancelChildInviteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_User_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CancelChildInviteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelChildInviteRequest) ProtoMessage() {}

func (x *CancelChildInviteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_User_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelChildInviteRequest.ProtoReflect.Descriptor instead.
func (*CancelChildInviteRequest) Descriptor() ([]byte, []int) {
	return file_User_proto_rawDescGZIP(), []int{6}
}

func (x *CancelChildInviteRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type UpdateUserChildRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Model *UserChildModel `protobuf:"bytes,1,opt,name=model,proto3" json:"model,omitempty"`
}

func (x *UpdateUserChildRequest) Reset() {
	*x = UpdateUserChildRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_User_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateUserChildRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserChildRequest) ProtoMessage() {}

func (x *UpdateUserChildRequest) ProtoReflect() protoreflect.Message {
	mi := &file_User_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserChildRequest.ProtoReflect.Descriptor instead.
func (*UpdateUserChildRequest) Descriptor() ([]byte, []int) {
	return file_User_proto_rawDescGZIP(), []int{7}
}

func (x *UpdateUserChildRequest) GetModel() *UserChildModel {
	if x != nil {
		return x.Model
	}
	return nil
}

type GetUserChildListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageInfo *Common.PageInfoRequest `protobuf:"bytes,1,opt,name=pageInfo,proto3" json:"pageInfo,omitempty"`
}

func (x *GetUserChildListRequest) Reset() {
	*x = GetUserChildListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_User_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserChildListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserChildListRequest) ProtoMessage() {}

func (x *GetUserChildListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_User_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserChildListRequest.ProtoReflect.Descriptor instead.
func (*GetUserChildListRequest) Descriptor() ([]byte, []int) {
	return file_User_proto_rawDescGZIP(), []int{8}
}

func (x *GetUserChildListRequest) GetPageInfo() *Common.PageInfoRequest {
	if x != nil {
		return x.PageInfo
	}
	return nil
}

type GetUserChildListReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Models   []*UserChildModel     `protobuf:"bytes,1,rep,name=models,proto3" json:"models,omitempty"`
	PageInfo *Common.PageInfoReply `protobuf:"bytes,2,opt,name=pageInfo,proto3" json:"pageInfo,omitempty"`
}

func (x *GetUserChildListReply) Reset() {
	*x = GetUserChildListReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_User_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserChildListReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserChildListReply) ProtoMessage() {}

func (x *GetUserChildListReply) ProtoReflect() protoreflect.Message {
	mi := &file_User_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserChildListReply.ProtoReflect.Descriptor instead.
func (*GetUserChildListReply) Descriptor() ([]byte, []int) {
	return file_User_proto_rawDescGZIP(), []int{9}
}

func (x *GetUserChildListReply) GetModels() []*UserChildModel {
	if x != nil {
		return x.Models
	}
	return nil
}

func (x *GetUserChildListReply) GetPageInfo() *Common.PageInfoReply {
	if x != nil {
		return x.PageInfo
	}
	return nil
}

type DeleteUserChildRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Model *UserChildModel `protobuf:"bytes,1,opt,name=model,proto3" json:"model,omitempty"`
}

func (x *DeleteUserChildRequest) Reset() {
	*x = DeleteUserChildRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_User_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteUserChildRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserChildRequest) ProtoMessage() {}

func (x *DeleteUserChildRequest) ProtoReflect() protoreflect.Message {
	mi := &file_User_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserChildRequest.ProtoReflect.Descriptor instead.
func (*DeleteUserChildRequest) Descriptor() ([]byte, []int) {
	return file_User_proto_rawDescGZIP(), []int{10}
}

func (x *DeleteUserChildRequest) GetModel() *UserChildModel {
	if x != nil {
		return x.Model
	}
	return nil
}

type UserChildModel struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId            string           `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`                                    // update
	Email             string           `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`                                                    // add
	Password          string           `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`                                              // confirm
	Name              string           `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`                                                      // add, update
	PermissionGroupId string           `protobuf:"bytes,5,opt,name=permission_group_id,json=permissionGroupId,proto3" json:"permission_group_id,omitempty"` // add, update
	MailState         Common.MailState `protobuf:"varint,6,opt,name=mail_state,json=mailState,proto3,enum=Common.MailState" json:"mail_state,omitempty"`    // list
	MailExpireDate    string           `protobuf:"bytes,7,opt,name=mail_expire_date,json=mailExpireDate,proto3" json:"mail_expire_date,omitempty"`          // list
}

func (x *UserChildModel) Reset() {
	*x = UserChildModel{}
	if protoimpl.UnsafeEnabled {
		mi := &file_User_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserChildModel) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserChildModel) ProtoMessage() {}

func (x *UserChildModel) ProtoReflect() protoreflect.Message {
	mi := &file_User_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserChildModel.ProtoReflect.Descriptor instead.
func (*UserChildModel) Descriptor() ([]byte, []int) {
	return file_User_proto_rawDescGZIP(), []int{11}
}

func (x *UserChildModel) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UserChildModel) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *UserChildModel) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *UserChildModel) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UserChildModel) GetPermissionGroupId() string {
	if x != nil {
		return x.PermissionGroupId
	}
	return ""
}

func (x *UserChildModel) GetMailState() Common.MailState {
	if x != nil {
		return x.MailState
	}
	return Common.MailState(0)
}

func (x *UserChildModel) GetMailExpireDate() string {
	if x != nil {
		return x.MailExpireDate
	}
	return ""
}

type ChangePasswordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OldPassword string `protobuf:"bytes,1,opt,name=old_password,json=oldPassword,proto3" json:"old_password,omitempty"`
	NewPassword string `protobuf:"bytes,2,opt,name=new_password,json=newPassword,proto3" json:"new_password,omitempty"`
}

func (x *ChangePasswordRequest) Reset() {
	*x = ChangePasswordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_User_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ChangePasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ChangePasswordRequest) ProtoMessage() {}

func (x *ChangePasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_User_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ChangePasswordRequest.ProtoReflect.Descriptor instead.
func (*ChangePasswordRequest) Descriptor() ([]byte, []int) {
	return file_User_proto_rawDescGZIP(), []int{12}
}

func (x *ChangePasswordRequest) GetOldPassword() string {
	if x != nil {
		return x.OldPassword
	}
	return ""
}

func (x *ChangePasswordRequest) GetNewPassword() string {
	if x != nil {
		return x.NewPassword
	}
	return ""
}

type ForgotPasswordRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Email string `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *ForgotPasswordRequest) Reset() {
	*x = ForgotPasswordRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_User_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ForgotPasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForgotPasswordRequest) ProtoMessage() {}

func (x *ForgotPasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_User_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForgotPasswordRequest.ProtoReflect.Descriptor instead.
func (*ForgotPasswordRequest) Descriptor() ([]byte, []int) {
	return file_User_proto_rawDescGZIP(), []int{13}
}

func (x *ForgotPasswordRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type UpdateUserProgramRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProgramId string `protobuf:"bytes,1,opt,name=program_id,json=programId,proto3" json:"program_id,omitempty"`
}

func (x *UpdateUserProgramRequest) Reset() {
	*x = UpdateUserProgramRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_User_proto_msgTypes[14]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateUserProgramRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserProgramRequest) ProtoMessage() {}

func (x *UpdateUserProgramRequest) ProtoReflect() protoreflect.Message {
	mi := &file_User_proto_msgTypes[14]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserProgramRequest.ProtoReflect.Descriptor instead.
func (*UpdateUserProgramRequest) Descriptor() ([]byte, []int) {
	return file_User_proto_rawDescGZIP(), []int{14}
}

func (x *UpdateUserProgramRequest) GetProgramId() string {
	if x != nil {
		return x.ProgramId
	}
	return ""
}

var File_User_proto protoreflect.FileDescriptor

var file_User_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x55, 0x73,
	0x65, 0x72, 0x1a, 0x0c, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x41, 0x0a, 0x13, 0x41, 0x64, 0x64, 0x55, 0x73, 0x65, 0x72, 0x43, 0x68, 0x69, 0x6c, 0x64,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x52, 0x05, 0x6d, 0x6f,
	0x64, 0x65, 0x6c, 0x22, 0x39, 0x0a, 0x1e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72,
	0x6d, 0x55, 0x73, 0x65, 0x72, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x72, 0x6c, 0x5f, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x72, 0x6c, 0x4b, 0x65, 0x79, 0x22, 0x38,
	0x0a, 0x1c, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x72, 0x6d, 0x55, 0x73, 0x65, 0x72,
	0x43, 0x68, 0x69, 0x6c, 0x64, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x18,
	0x0a, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x4e, 0x0a, 0x17, 0x43, 0x6f, 0x6e, 0x66,
	0x69, 0x72, 0x6d, 0x55, 0x73, 0x65, 0x72, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x72, 0x6c, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x72, 0x6c, 0x4b, 0x65, 0x79, 0x12, 0x1a, 0x0a, 0x08,
	0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x33, 0x0a, 0x18, 0x52, 0x65, 0x73, 0x65,
	0x6e, 0x64, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x31, 0x0a,
	0x16, 0x52, 0x65, 0x73, 0x65, 0x6e, 0x64, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x49, 0x6e, 0x76, 0x69,
	0x74, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x72, 0x6c, 0x5f, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x72, 0x6c, 0x4b, 0x65, 0x79,
	0x22, 0x33, 0x0a, 0x18, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x49,
	0x6e, 0x76, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x44, 0x0a, 0x16, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x55,
	0x73, 0x65, 0x72, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x2a, 0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x4d,
	0x6f, 0x64, 0x65, 0x6c, 0x52, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x22, 0x4e, 0x0a, 0x17, 0x47,
	0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x33, 0x0a, 0x08, 0x70, 0x61, 0x67, 0x65, 0x49, 0x6e,
	0x66, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x50, 0x61, 0x67, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x78, 0x0a, 0x15, 0x47,
	0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x2c, 0x0a, 0x06, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x2e, 0x55, 0x73, 0x65, 0x72,
	0x43, 0x68, 0x69, 0x6c, 0x64, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x52, 0x06, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x73, 0x12, 0x31, 0x0a, 0x08, 0x70, 0x61, 0x67, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x50, 0x61,
	0x67, 0x65, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x52, 0x08, 0x70, 0x61, 0x67,
	0x65, 0x49, 0x6e, 0x66, 0x6f, 0x22, 0x44, 0x0a, 0x16, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x55,
	0x73, 0x65, 0x72, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x2a, 0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x55, 0x73, 0x65, 0x72, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x4d,
	0x6f, 0x64, 0x65, 0x6c, 0x52, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x22, 0xfb, 0x01, 0x0a, 0x0e,
	0x55, 0x73, 0x65, 0x72, 0x43, 0x68, 0x69, 0x6c, 0x64, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x17,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a,
	0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x2e, 0x0a,
	0x13, 0x70, 0x65, 0x72, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x70, 0x65, 0x72, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x12, 0x30, 0x0a,
	0x0a, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x11, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x4d, 0x61, 0x69, 0x6c, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x52, 0x09, 0x6d, 0x61, 0x69, 0x6c, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x28, 0x0a, 0x10, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x5f, 0x64,
	0x61, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6d, 0x61, 0x69, 0x6c, 0x45,
	0x78, 0x70, 0x69, 0x72, 0x65, 0x44, 0x61, 0x74, 0x65, 0x22, 0x5d, 0x0a, 0x15, 0x43, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x6f, 0x6c, 0x64, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f,
	0x72, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6f, 0x6c, 0x64, 0x50, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x6e, 0x65, 0x77, 0x5f, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x6e, 0x65, 0x77,
	0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x2d, 0x0a, 0x15, 0x46, 0x6f, 0x72, 0x67,
	0x6f, 0x74, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x39, 0x0a, 0x18, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x55, 0x73, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x61, 0x6d,
	0x49, 0x64, 0x42, 0x1b, 0x5a, 0x19, 0x4d, 0x45, 0x54, 0x41, 0x52, 0x5f, 0x62, 0x61, 0x63, 0x6b,
	0x65, 0x6e, 0x64, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x55, 0x73, 0x65, 0x72, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_User_proto_rawDescOnce sync.Once
	file_User_proto_rawDescData = file_User_proto_rawDesc
)

func file_User_proto_rawDescGZIP() []byte {
	file_User_proto_rawDescOnce.Do(func() {
		file_User_proto_rawDescData = protoimpl.X.CompressGZIP(file_User_proto_rawDescData)
	})
	return file_User_proto_rawDescData
}

var file_User_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_User_proto_goTypes = []interface{}{
	(*AddUserChildRequest)(nil),            // 0: User.AddUserChildRequest
	(*GetConfirmUserChildInfoRequest)(nil), // 1: User.GetConfirmUserChildInfoRequest
	(*GetConfirmUserChildInfoReply)(nil),   // 2: User.GetConfirmUserChildInfoReply
	(*ConfirmUserChildRequest)(nil),        // 3: User.ConfirmUserChildRequest
	(*ResendChildInviteRequest)(nil),       // 4: User.ResendChildInviteRequest
	(*ResendChildInviteReply)(nil),         // 5: User.ResendChildInviteReply
	(*CancelChildInviteRequest)(nil),       // 6: User.CancelChildInviteRequest
	(*UpdateUserChildRequest)(nil),         // 7: User.UpdateUserChildRequest
	(*GetUserChildListRequest)(nil),        // 8: User.GetUserChildListRequest
	(*GetUserChildListReply)(nil),          // 9: User.GetUserChildListReply
	(*DeleteUserChildRequest)(nil),         // 10: User.DeleteUserChildRequest
	(*UserChildModel)(nil),                 // 11: User.UserChildModel
	(*ChangePasswordRequest)(nil),          // 12: User.ChangePasswordRequest
	(*ForgotPasswordRequest)(nil),          // 13: User.ForgotPasswordRequest
	(*UpdateUserProgramRequest)(nil),       // 14: User.UpdateUserProgramRequest
	(*Common.PageInfoRequest)(nil),         // 15: Common.PageInfoRequest
	(*Common.PageInfoReply)(nil),           // 16: Common.PageInfoReply
	(Common.MailState)(0),                  // 17: Common.MailState
}
var file_User_proto_depIdxs = []int32{
	11, // 0: User.AddUserChildRequest.model:type_name -> User.UserChildModel
	11, // 1: User.UpdateUserChildRequest.model:type_name -> User.UserChildModel
	15, // 2: User.GetUserChildListRequest.pageInfo:type_name -> Common.PageInfoRequest
	11, // 3: User.GetUserChildListReply.models:type_name -> User.UserChildModel
	16, // 4: User.GetUserChildListReply.pageInfo:type_name -> Common.PageInfoReply
	11, // 5: User.DeleteUserChildRequest.model:type_name -> User.UserChildModel
	17, // 6: User.UserChildModel.mail_state:type_name -> Common.MailState
	7,  // [7:7] is the sub-list for method output_type
	7,  // [7:7] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_User_proto_init() }
func file_User_proto_init() {
	if File_User_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_User_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddUserChildRequest); i {
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
		file_User_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetConfirmUserChildInfoRequest); i {
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
		file_User_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetConfirmUserChildInfoReply); i {
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
		file_User_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConfirmUserChildRequest); i {
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
		file_User_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResendChildInviteRequest); i {
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
		file_User_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResendChildInviteReply); i {
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
		file_User_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CancelChildInviteRequest); i {
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
		file_User_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateUserChildRequest); i {
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
		file_User_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserChildListRequest); i {
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
		file_User_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserChildListReply); i {
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
		file_User_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteUserChildRequest); i {
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
		file_User_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserChildModel); i {
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
		file_User_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ChangePasswordRequest); i {
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
		file_User_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ForgotPasswordRequest); i {
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
		file_User_proto_msgTypes[14].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateUserProgramRequest); i {
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
			RawDescriptor: file_User_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_User_proto_goTypes,
		DependencyIndexes: file_User_proto_depIdxs,
		MessageInfos:      file_User_proto_msgTypes,
	}.Build()
	File_User_proto = out.File
	file_User_proto_rawDesc = nil
	file_User_proto_goTypes = nil
	file_User_proto_depIdxs = nil
}
