// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.27.2
// source: pkg/overlord/overlord.proto

package overlord

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

type Status struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *Status) Reset() {
	*x = Status{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_overlord_overlord_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Status) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Status) ProtoMessage() {}

func (x *Status) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_overlord_overlord_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Status.ProtoReflect.Descriptor instead.
func (*Status) Descriptor() ([]byte, []int) {
	return file_pkg_overlord_overlord_proto_rawDescGZIP(), []int{0}
}

func (x *Status) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type GetIntegratorConfigIn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Integrator string `protobuf:"bytes,1,opt,name=integrator,proto3" json:"integrator,omitempty"`
	Game       string `protobuf:"bytes,2,opt,name=game,proto3" json:"game,omitempty"`
}

func (x *GetIntegratorConfigIn) Reset() {
	*x = GetIntegratorConfigIn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_overlord_overlord_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetIntegratorConfigIn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetIntegratorConfigIn) ProtoMessage() {}

func (x *GetIntegratorConfigIn) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_overlord_overlord_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetIntegratorConfigIn.ProtoReflect.Descriptor instead.
func (*GetIntegratorConfigIn) Descriptor() ([]byte, []int) {
	return file_pkg_overlord_overlord_proto_rawDescGZIP(), []int{1}
}

func (x *GetIntegratorConfigIn) GetIntegrator() string {
	if x != nil {
		return x.Integrator
	}
	return ""
}

func (x *GetIntegratorConfigIn) GetGame() string {
	if x != nil {
		return x.Game
	}
	return ""
}

type GetIntegratorConfigOut struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DefaultWager int64            `protobuf:"varint,1,opt,name=default_wager,json=defaultWager,proto3" json:"default_wager,omitempty"`
	WagerLevels  []int64          `protobuf:"varint,2,rep,packed,name=wager_levels,json=wagerLevels,proto3" json:"wager_levels,omitempty"`
	Multipliers  map[string]int64 `protobuf:"bytes,3,rep,name=multipliers,proto3" json:"multipliers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
}

func (x *GetIntegratorConfigOut) Reset() {
	*x = GetIntegratorConfigOut{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_overlord_overlord_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetIntegratorConfigOut) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetIntegratorConfigOut) ProtoMessage() {}

func (x *GetIntegratorConfigOut) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_overlord_overlord_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetIntegratorConfigOut.ProtoReflect.Descriptor instead.
func (*GetIntegratorConfigOut) Descriptor() ([]byte, []int) {
	return file_pkg_overlord_overlord_proto_rawDescGZIP(), []int{2}
}

func (x *GetIntegratorConfigOut) GetDefaultWager() int64 {
	if x != nil {
		return x.DefaultWager
	}
	return 0
}

func (x *GetIntegratorConfigOut) GetWagerLevels() []int64 {
	if x != nil {
		return x.WagerLevels
	}
	return nil
}

func (x *GetIntegratorConfigOut) GetMultipliers() map[string]int64 {
	if x != nil {
		return x.Multipliers
	}
	return nil
}

type SaveParamsIn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Integrator   string  `protobuf:"bytes,1,opt,name=integrator,proto3" json:"integrator,omitempty"`
	Game         string  `protobuf:"bytes,2,opt,name=game,proto3" json:"game,omitempty"`
	Rtp          *int64  `protobuf:"varint,3,opt,name=rtp,proto3,oneof" json:"rtp,omitempty"`
	Wagers       []int64 `protobuf:"varint,4,rep,packed,name=wagers,proto3" json:"wagers,omitempty"`
	BuyBonus     bool    `protobuf:"varint,5,opt,name=buy_bonus,json=buyBonus,proto3" json:"buy_bonus,omitempty"`
	Gamble       bool    `protobuf:"varint,6,opt,name=gamble,proto3" json:"gamble,omitempty"`
	DoubleChance bool    `protobuf:"varint,7,opt,name=double_chance,json=doubleChance,proto3" json:"double_chance,omitempty"`
	SessionId    string  `protobuf:"bytes,8,opt,name=session_id,json=sessionId,proto3" json:"session_id,omitempty"`
	Volatility   *string `protobuf:"bytes,9,opt,name=volatility,proto3,oneof" json:"volatility,omitempty"`
	IsDemo       bool    `protobuf:"varint,10,opt,name=is_demo,json=isDemo,proto3" json:"is_demo,omitempty"`
	Currency     string  `protobuf:"bytes,11,opt,name=currency,proto3" json:"currency,omitempty"`
	UserId       string  `protobuf:"bytes,12,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	UserLocale   string  `protobuf:"bytes,13,opt,name=user_locale,json=userLocale,proto3" json:"user_locale,omitempty"`
	DefaultWager *int64  `protobuf:"varint,14,opt,name=default_wager,json=defaultWager,proto3,oneof" json:"default_wager,omitempty"`
	Jurisdiction string  `protobuf:"bytes,15,opt,name=jurisdiction,proto3" json:"jurisdiction,omitempty"`
	LobbyUrl     string  `protobuf:"bytes,16,opt,name=lobby_url,json=lobbyUrl,proto3" json:"lobby_url,omitempty"`
	ShowCheats   bool    `protobuf:"varint,17,opt,name=show_cheats,json=showCheats,proto3" json:"show_cheats,omitempty"`
	LowBalance   bool    `protobuf:"varint,18,opt,name=low_balance,json=lowBalance,proto3" json:"low_balance,omitempty"`
	ShortLink    bool    `protobuf:"varint,19,opt,name=short_link,json=shortLink,proto3" json:"short_link,omitempty"`
}

func (x *SaveParamsIn) Reset() {
	*x = SaveParamsIn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_overlord_overlord_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveParamsIn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveParamsIn) ProtoMessage() {}

func (x *SaveParamsIn) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_overlord_overlord_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveParamsIn.ProtoReflect.Descriptor instead.
func (*SaveParamsIn) Descriptor() ([]byte, []int) {
	return file_pkg_overlord_overlord_proto_rawDescGZIP(), []int{3}
}

func (x *SaveParamsIn) GetIntegrator() string {
	if x != nil {
		return x.Integrator
	}
	return ""
}

func (x *SaveParamsIn) GetGame() string {
	if x != nil {
		return x.Game
	}
	return ""
}

func (x *SaveParamsIn) GetRtp() int64 {
	if x != nil && x.Rtp != nil {
		return *x.Rtp
	}
	return 0
}

func (x *SaveParamsIn) GetWagers() []int64 {
	if x != nil {
		return x.Wagers
	}
	return nil
}

func (x *SaveParamsIn) GetBuyBonus() bool {
	if x != nil {
		return x.BuyBonus
	}
	return false
}

func (x *SaveParamsIn) GetGamble() bool {
	if x != nil {
		return x.Gamble
	}
	return false
}

func (x *SaveParamsIn) GetDoubleChance() bool {
	if x != nil {
		return x.DoubleChance
	}
	return false
}

func (x *SaveParamsIn) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *SaveParamsIn) GetVolatility() string {
	if x != nil && x.Volatility != nil {
		return *x.Volatility
	}
	return ""
}

func (x *SaveParamsIn) GetIsDemo() bool {
	if x != nil {
		return x.IsDemo
	}
	return false
}

func (x *SaveParamsIn) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

func (x *SaveParamsIn) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *SaveParamsIn) GetUserLocale() string {
	if x != nil {
		return x.UserLocale
	}
	return ""
}

func (x *SaveParamsIn) GetDefaultWager() int64 {
	if x != nil && x.DefaultWager != nil {
		return *x.DefaultWager
	}
	return 0
}

func (x *SaveParamsIn) GetJurisdiction() string {
	if x != nil {
		return x.Jurisdiction
	}
	return ""
}

func (x *SaveParamsIn) GetLobbyUrl() string {
	if x != nil {
		return x.LobbyUrl
	}
	return ""
}

func (x *SaveParamsIn) GetShowCheats() bool {
	if x != nil {
		return x.ShowCheats
	}
	return false
}

func (x *SaveParamsIn) GetLowBalance() bool {
	if x != nil {
		return x.LowBalance
	}
	return false
}

func (x *SaveParamsIn) GetShortLink() bool {
	if x != nil {
		return x.ShortLink
	}
	return false
}

type SaveParamsOut struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *SaveParamsOut) Reset() {
	*x = SaveParamsOut{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_overlord_overlord_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveParamsOut) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveParamsOut) ProtoMessage() {}

func (x *SaveParamsOut) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_overlord_overlord_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveParamsOut.ProtoReflect.Descriptor instead.
func (*SaveParamsOut) Descriptor() ([]byte, []int) {
	return file_pkg_overlord_overlord_proto_rawDescGZIP(), []int{4}
}

var File_pkg_overlord_overlord_proto protoreflect.FileDescriptor

var file_pkg_overlord_overlord_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x6b, 0x67, 0x2f, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x6f, 0x72, 0x64, 0x2f, 0x6f,
	0x76, 0x65, 0x72, 0x6c, 0x6f, 0x72, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x6f,
	0x76, 0x65, 0x72, 0x6c, 0x6f, 0x72, 0x64, 0x22, 0x20, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x4b, 0x0a, 0x15, 0x47, 0x65, 0x74,
	0x49, 0x6e, 0x74, 0x65, 0x67, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x49, 0x6e, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x6e, 0x74, 0x65, 0x67, 0x72, 0x61, 0x74, 0x6f, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x6e, 0x74, 0x65, 0x67, 0x72, 0x61, 0x74,
	0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x67, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x67, 0x61, 0x6d, 0x65, 0x22, 0xf5, 0x01, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x49, 0x6e,
	0x74, 0x65, 0x67, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4f, 0x75,
	0x74, 0x12, 0x23, 0x0a, 0x0d, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x77, 0x61, 0x67,
	0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c,
	0x74, 0x57, 0x61, 0x67, 0x65, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x77, 0x61, 0x67, 0x65, 0x72, 0x5f,
	0x6c, 0x65, 0x76, 0x65, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x03, 0x52, 0x0b, 0x77, 0x61,
	0x67, 0x65, 0x72, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x73, 0x12, 0x53, 0x0a, 0x0b, 0x6d, 0x75, 0x6c,
	0x74, 0x69, 0x70, 0x6c, 0x69, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x31,
	0x2e, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x6f, 0x72, 0x64, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x74,
	0x65, 0x67, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4f, 0x75, 0x74,
	0x2e, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x69, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x0b, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x69, 0x65, 0x72, 0x73, 0x1a, 0x3e,
	0x0a, 0x10, 0x4d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x69, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xf3,
	0x04, 0x0a, 0x0c, 0x53, 0x61, 0x76, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x49, 0x6e, 0x12,
	0x1e, 0x0a, 0x0a, 0x69, 0x6e, 0x74, 0x65, 0x67, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x69, 0x6e, 0x74, 0x65, 0x67, 0x72, 0x61, 0x74, 0x6f, 0x72, 0x12,
	0x12, 0x0a, 0x04, 0x67, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x67,
	0x61, 0x6d, 0x65, 0x12, 0x15, 0x0a, 0x03, 0x72, 0x74, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x48, 0x00, 0x52, 0x03, 0x72, 0x74, 0x70, 0x88, 0x01, 0x01, 0x12, 0x16, 0x0a, 0x06, 0x77, 0x61,
	0x67, 0x65, 0x72, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x03, 0x52, 0x06, 0x77, 0x61, 0x67, 0x65,
	0x72, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x62, 0x75, 0x79, 0x5f, 0x62, 0x6f, 0x6e, 0x75, 0x73, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x62, 0x75, 0x79, 0x42, 0x6f, 0x6e, 0x75, 0x73, 0x12,
	0x16, 0x0a, 0x06, 0x67, 0x61, 0x6d, 0x62, 0x6c, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x06, 0x67, 0x61, 0x6d, 0x62, 0x6c, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x64, 0x6f, 0x75, 0x62, 0x6c,
	0x65, 0x5f, 0x63, 0x68, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c,
	0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x0a, 0x76,
	0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x01, 0x52, 0x0a, 0x76, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x88, 0x01, 0x01,
	0x12, 0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f, 0x64, 0x65, 0x6d, 0x6f, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x06, 0x69, 0x73, 0x44, 0x65, 0x6d, 0x6f, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x63, 0x79, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x75, 0x72,
	0x72, 0x65, 0x6e, 0x63, 0x79, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1f,
	0x0a, 0x0b, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x65, 0x18, 0x0d, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x65, 0x12,
	0x28, 0x0a, 0x0d, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x77, 0x61, 0x67, 0x65, 0x72,
	0x18, 0x0e, 0x20, 0x01, 0x28, 0x03, 0x48, 0x02, 0x52, 0x0c, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c,
	0x74, 0x57, 0x61, 0x67, 0x65, 0x72, 0x88, 0x01, 0x01, 0x12, 0x22, 0x0a, 0x0c, 0x6a, 0x75, 0x72,
	0x69, 0x73, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x6a, 0x75, 0x72, 0x69, 0x73, 0x64, 0x69, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a,
	0x09, 0x6c, 0x6f, 0x62, 0x62, 0x79, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x10, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x6c, 0x6f, 0x62, 0x62, 0x79, 0x55, 0x72, 0x6c, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x68,
	0x6f, 0x77, 0x5f, 0x63, 0x68, 0x65, 0x61, 0x74, 0x73, 0x18, 0x11, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x0a, 0x73, 0x68, 0x6f, 0x77, 0x43, 0x68, 0x65, 0x61, 0x74, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x6c,
	0x6f, 0x77, 0x5f, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x12, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0a, 0x6c, 0x6f, 0x77, 0x42, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x1d, 0x0a, 0x0a,
	0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x6c, 0x69, 0x6e, 0x6b, 0x18, 0x13, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x09, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x4c, 0x69, 0x6e, 0x6b, 0x42, 0x06, 0x0a, 0x04, 0x5f,
	0x72, 0x74, 0x70, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x76, 0x6f, 0x6c, 0x61, 0x74, 0x69, 0x6c, 0x69,
	0x74, 0x79, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x77,
	0x61, 0x67, 0x65, 0x72, 0x22, 0x0f, 0x0a, 0x0d, 0x53, 0x61, 0x76, 0x65, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x73, 0x4f, 0x75, 0x74, 0x32, 0xe0, 0x01, 0x0a, 0x08, 0x4f, 0x76, 0x65, 0x72, 0x6c, 0x6f,
	0x72, 0x64, 0x12, 0x5a, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x67, 0x72, 0x61,
	0x74, 0x6f, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x1f, 0x2e, 0x6f, 0x76, 0x65, 0x72,
	0x6c, 0x6f, 0x72, 0x64, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x67, 0x72, 0x61, 0x74,
	0x6f, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x49, 0x6e, 0x1a, 0x20, 0x2e, 0x6f, 0x76, 0x65,
	0x72, 0x6c, 0x6f, 0x72, 0x64, 0x2e, 0x47, 0x65, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x67, 0x72, 0x61,
	0x74, 0x6f, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x4f, 0x75, 0x74, 0x22, 0x00, 0x12, 0x3f,
	0x0a, 0x0a, 0x53, 0x61, 0x76, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x16, 0x2e, 0x6f,
	0x76, 0x65, 0x72, 0x6c, 0x6f, 0x72, 0x64, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x73, 0x49, 0x6e, 0x1a, 0x17, 0x2e, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x6f, 0x72, 0x64, 0x2e,
	0x53, 0x61, 0x76, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x4f, 0x75, 0x74, 0x22, 0x00, 0x12,
	0x37, 0x0a, 0x0b, 0x48, 0x65, 0x61, 0x6c, 0x74, 0x68, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x12, 0x10,
	0x2e, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x6f, 0x72, 0x64, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x1a, 0x10, 0x2e, 0x6f, 0x76, 0x65, 0x72, 0x6c, 0x6f, 0x72, 0x64, 0x2e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x0c, 0x5a, 0x0a, 0x2e, 0x2f, 0x6f, 0x76,
	0x65, 0x72, 0x6c, 0x6f, 0x72, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_overlord_overlord_proto_rawDescOnce sync.Once
	file_pkg_overlord_overlord_proto_rawDescData = file_pkg_overlord_overlord_proto_rawDesc
)

func file_pkg_overlord_overlord_proto_rawDescGZIP() []byte {
	file_pkg_overlord_overlord_proto_rawDescOnce.Do(func() {
		file_pkg_overlord_overlord_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_overlord_overlord_proto_rawDescData)
	})
	return file_pkg_overlord_overlord_proto_rawDescData
}

var file_pkg_overlord_overlord_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_pkg_overlord_overlord_proto_goTypes = []interface{}{
	(*Status)(nil),                 // 0: overlord.Status
	(*GetIntegratorConfigIn)(nil),  // 1: overlord.GetIntegratorConfigIn
	(*GetIntegratorConfigOut)(nil), // 2: overlord.GetIntegratorConfigOut
	(*SaveParamsIn)(nil),           // 3: overlord.SaveParamsIn
	(*SaveParamsOut)(nil),          // 4: overlord.SaveParamsOut
	nil,                            // 5: overlord.GetIntegratorConfigOut.MultipliersEntry
}
var file_pkg_overlord_overlord_proto_depIdxs = []int32{
	5, // 0: overlord.GetIntegratorConfigOut.multipliers:type_name -> overlord.GetIntegratorConfigOut.MultipliersEntry
	1, // 1: overlord.Overlord.GetIntegratorConfig:input_type -> overlord.GetIntegratorConfigIn
	3, // 2: overlord.Overlord.SaveParams:input_type -> overlord.SaveParamsIn
	0, // 3: overlord.Overlord.HealthCheck:input_type -> overlord.Status
	2, // 4: overlord.Overlord.GetIntegratorConfig:output_type -> overlord.GetIntegratorConfigOut
	4, // 5: overlord.Overlord.SaveParams:output_type -> overlord.SaveParamsOut
	0, // 6: overlord.Overlord.HealthCheck:output_type -> overlord.Status
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pkg_overlord_overlord_proto_init() }
func file_pkg_overlord_overlord_proto_init() {
	if File_pkg_overlord_overlord_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_overlord_overlord_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Status); i {
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
		file_pkg_overlord_overlord_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetIntegratorConfigIn); i {
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
		file_pkg_overlord_overlord_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetIntegratorConfigOut); i {
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
		file_pkg_overlord_overlord_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveParamsIn); i {
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
		file_pkg_overlord_overlord_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveParamsOut); i {
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
	file_pkg_overlord_overlord_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_overlord_overlord_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_overlord_overlord_proto_goTypes,
		DependencyIndexes: file_pkg_overlord_overlord_proto_depIdxs,
		MessageInfos:      file_pkg_overlord_overlord_proto_msgTypes,
	}.Build()
	File_pkg_overlord_overlord_proto = out.File
	file_pkg_overlord_overlord_proto_rawDesc = nil
	file_pkg_overlord_overlord_proto_goTypes = nil
	file_pkg_overlord_overlord_proto_depIdxs = nil
}
