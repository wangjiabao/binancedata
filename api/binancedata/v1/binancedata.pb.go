// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.7
// source: api/binancedata/v1/binancedata.proto

package v1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type DownloadBinanceDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DownloadBinanceDataRequest) Reset() {
	*x = DownloadBinanceDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_binancedata_v1_binancedata_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadBinanceDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadBinanceDataRequest) ProtoMessage() {}

func (x *DownloadBinanceDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_binancedata_v1_binancedata_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadBinanceDataRequest.ProtoReflect.Descriptor instead.
func (*DownloadBinanceDataRequest) Descriptor() ([]byte, []int) {
	return file_api_binancedata_v1_binancedata_proto_rawDescGZIP(), []int{0}
}

type DownloadBinanceDataReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DownloadBinanceDataReply) Reset() {
	*x = DownloadBinanceDataReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_binancedata_v1_binancedata_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DownloadBinanceDataReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DownloadBinanceDataReply) ProtoMessage() {}

func (x *DownloadBinanceDataReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_binancedata_v1_binancedata_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DownloadBinanceDataReply.ProtoReflect.Descriptor instead.
func (*DownloadBinanceDataReply) Descriptor() ([]byte, []int) {
	return file_api_binancedata_v1_binancedata_proto_rawDescGZIP(), []int{1}
}

type IntervalMAvgEndPriceDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Start string `protobuf:"bytes,1,opt,name=start,proto3" json:"start,omitempty"`
	End   string `protobuf:"bytes,2,opt,name=end,proto3" json:"end,omitempty"`
	M     int64  `protobuf:"varint,3,opt,name=m,proto3" json:"m,omitempty"`
	N     int64  `protobuf:"varint,4,opt,name=n,proto3" json:"n,omitempty"`
}

func (x *IntervalMAvgEndPriceDataRequest) Reset() {
	*x = IntervalMAvgEndPriceDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_binancedata_v1_binancedata_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IntervalMAvgEndPriceDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IntervalMAvgEndPriceDataRequest) ProtoMessage() {}

func (x *IntervalMAvgEndPriceDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_binancedata_v1_binancedata_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IntervalMAvgEndPriceDataRequest.ProtoReflect.Descriptor instead.
func (*IntervalMAvgEndPriceDataRequest) Descriptor() ([]byte, []int) {
	return file_api_binancedata_v1_binancedata_proto_rawDescGZIP(), []int{2}
}

func (x *IntervalMAvgEndPriceDataRequest) GetStart() string {
	if x != nil {
		return x.Start
	}
	return ""
}

func (x *IntervalMAvgEndPriceDataRequest) GetEnd() string {
	if x != nil {
		return x.End
	}
	return ""
}

func (x *IntervalMAvgEndPriceDataRequest) GetM() int64 {
	if x != nil {
		return x.M
	}
	return 0
}

func (x *IntervalMAvgEndPriceDataRequest) GetN() int64 {
	if x != nil {
		return x.N
	}
	return 0
}

type IntervalMAvgEndPriceDataReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data                []*IntervalMAvgEndPriceDataReply_List  `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	OperationData       []*IntervalMAvgEndPriceDataReply_List2 `protobuf:"bytes,2,rep,name=operationData,proto3" json:"operationData,omitempty"`
	OperationOrderTotal int64                                  `protobuf:"varint,3,opt,name=operationOrderTotal,proto3" json:"operationOrderTotal,omitempty"`
	OperationWinRate    string                                 `protobuf:"bytes,4,opt,name=operationWinRate,proto3" json:"operationWinRate,omitempty"`
	OperationWinAmount  string                                 `protobuf:"bytes,5,opt,name=operationWinAmount,proto3" json:"operationWinAmount,omitempty"`
}

func (x *IntervalMAvgEndPriceDataReply) Reset() {
	*x = IntervalMAvgEndPriceDataReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_binancedata_v1_binancedata_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IntervalMAvgEndPriceDataReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IntervalMAvgEndPriceDataReply) ProtoMessage() {}

func (x *IntervalMAvgEndPriceDataReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_binancedata_v1_binancedata_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IntervalMAvgEndPriceDataReply.ProtoReflect.Descriptor instead.
func (*IntervalMAvgEndPriceDataReply) Descriptor() ([]byte, []int) {
	return file_api_binancedata_v1_binancedata_proto_rawDescGZIP(), []int{3}
}

func (x *IntervalMAvgEndPriceDataReply) GetData() []*IntervalMAvgEndPriceDataReply_List {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *IntervalMAvgEndPriceDataReply) GetOperationData() []*IntervalMAvgEndPriceDataReply_List2 {
	if x != nil {
		return x.OperationData
	}
	return nil
}

func (x *IntervalMAvgEndPriceDataReply) GetOperationOrderTotal() int64 {
	if x != nil {
		return x.OperationOrderTotal
	}
	return 0
}

func (x *IntervalMAvgEndPriceDataReply) GetOperationWinRate() string {
	if x != nil {
		return x.OperationWinRate
	}
	return ""
}

func (x *IntervalMAvgEndPriceDataReply) GetOperationWinAmount() string {
	if x != nil {
		return x.OperationWinAmount
	}
	return ""
}

type IntervalMAvgEndPriceDataReply_List struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartPrice            string `protobuf:"bytes,1,opt,name=start_price,json=startPrice,proto3" json:"start_price,omitempty"`
	EndPrice              string `protobuf:"bytes,2,opt,name=end_price,json=endPrice,proto3" json:"end_price,omitempty"`
	TopPrice              string `protobuf:"bytes,3,opt,name=top_price,json=topPrice,proto3" json:"top_price,omitempty"`
	LowPrice              string `protobuf:"bytes,4,opt,name=low_price,json=lowPrice,proto3" json:"low_price,omitempty"`
	WithBeforeAvgEndPrice string `protobuf:"bytes,5,opt,name=with_before_avg_end_price,json=withBeforeAvgEndPrice,proto3" json:"with_before_avg_end_price,omitempty"`
	Time                  string `protobuf:"bytes,6,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *IntervalMAvgEndPriceDataReply_List) Reset() {
	*x = IntervalMAvgEndPriceDataReply_List{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_binancedata_v1_binancedata_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IntervalMAvgEndPriceDataReply_List) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IntervalMAvgEndPriceDataReply_List) ProtoMessage() {}

func (x *IntervalMAvgEndPriceDataReply_List) ProtoReflect() protoreflect.Message {
	mi := &file_api_binancedata_v1_binancedata_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IntervalMAvgEndPriceDataReply_List.ProtoReflect.Descriptor instead.
func (*IntervalMAvgEndPriceDataReply_List) Descriptor() ([]byte, []int) {
	return file_api_binancedata_v1_binancedata_proto_rawDescGZIP(), []int{3, 0}
}

func (x *IntervalMAvgEndPriceDataReply_List) GetStartPrice() string {
	if x != nil {
		return x.StartPrice
	}
	return ""
}

func (x *IntervalMAvgEndPriceDataReply_List) GetEndPrice() string {
	if x != nil {
		return x.EndPrice
	}
	return ""
}

func (x *IntervalMAvgEndPriceDataReply_List) GetTopPrice() string {
	if x != nil {
		return x.TopPrice
	}
	return ""
}

func (x *IntervalMAvgEndPriceDataReply_List) GetLowPrice() string {
	if x != nil {
		return x.LowPrice
	}
	return ""
}

func (x *IntervalMAvgEndPriceDataReply_List) GetWithBeforeAvgEndPrice() string {
	if x != nil {
		return x.WithBeforeAvgEndPrice
	}
	return ""
}

func (x *IntervalMAvgEndPriceDataReply_List) GetTime() string {
	if x != nil {
		return x.Time
	}
	return ""
}

type IntervalMAvgEndPriceDataReply_List2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartPrice    string `protobuf:"bytes,1,opt,name=start_price,json=startPrice,proto3" json:"start_price,omitempty"`
	EndPrice      string `protobuf:"bytes,2,opt,name=end_price,json=endPrice,proto3" json:"end_price,omitempty"`
	StartTime     string `protobuf:"bytes,3,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	Time          string `protobuf:"bytes,4,opt,name=time,proto3" json:"time,omitempty"`
	EndTime       string `protobuf:"bytes,6,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	Type          string `protobuf:"bytes,7,opt,name=type,proto3" json:"type,omitempty"`
	Status        string `protobuf:"bytes,8,opt,name=status,proto3" json:"status,omitempty"`
	Rate          string `protobuf:"bytes,9,opt,name=rate,proto3" json:"rate,omitempty"`
	CloseEndPrice string `protobuf:"bytes,10,opt,name=CloseEndPrice,proto3" json:"CloseEndPrice,omitempty"`
}

func (x *IntervalMAvgEndPriceDataReply_List2) Reset() {
	*x = IntervalMAvgEndPriceDataReply_List2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_binancedata_v1_binancedata_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IntervalMAvgEndPriceDataReply_List2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IntervalMAvgEndPriceDataReply_List2) ProtoMessage() {}

func (x *IntervalMAvgEndPriceDataReply_List2) ProtoReflect() protoreflect.Message {
	mi := &file_api_binancedata_v1_binancedata_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IntervalMAvgEndPriceDataReply_List2.ProtoReflect.Descriptor instead.
func (*IntervalMAvgEndPriceDataReply_List2) Descriptor() ([]byte, []int) {
	return file_api_binancedata_v1_binancedata_proto_rawDescGZIP(), []int{3, 1}
}

func (x *IntervalMAvgEndPriceDataReply_List2) GetStartPrice() string {
	if x != nil {
		return x.StartPrice
	}
	return ""
}

func (x *IntervalMAvgEndPriceDataReply_List2) GetEndPrice() string {
	if x != nil {
		return x.EndPrice
	}
	return ""
}

func (x *IntervalMAvgEndPriceDataReply_List2) GetStartTime() string {
	if x != nil {
		return x.StartTime
	}
	return ""
}

func (x *IntervalMAvgEndPriceDataReply_List2) GetTime() string {
	if x != nil {
		return x.Time
	}
	return ""
}

func (x *IntervalMAvgEndPriceDataReply_List2) GetEndTime() string {
	if x != nil {
		return x.EndTime
	}
	return ""
}

func (x *IntervalMAvgEndPriceDataReply_List2) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *IntervalMAvgEndPriceDataReply_List2) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *IntervalMAvgEndPriceDataReply_List2) GetRate() string {
	if x != nil {
		return x.Rate
	}
	return ""
}

func (x *IntervalMAvgEndPriceDataReply_List2) GetCloseEndPrice() string {
	if x != nil {
		return x.CloseEndPrice
	}
	return ""
}

var File_api_binancedata_v1_binancedata_proto protoreflect.FileDescriptor

var file_api_binancedata_v1_binancedata_proto_rawDesc = []byte{
	0x0a, 0x24, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x64, 0x61, 0x74,
	0x61, 0x2f, 0x76, 0x31, 0x2f, 0x62, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x61, 0x70, 0x69, 0x2e, 0x62, 0x69, 0x6e, 0x61,
	0x6e, 0x63, 0x65, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61,
	0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x1c, 0x0a, 0x1a, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x69, 0x6e,
	0x61, 0x6e, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0x1a, 0x0a, 0x18, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x69, 0x6e, 0x61, 0x6e,
	0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x65, 0x0a, 0x1f, 0x49,
	0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x4d, 0x41, 0x76, 0x67, 0x45, 0x6e, 0x64, 0x50, 0x72,
	0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x12, 0x0c, 0x0a, 0x01, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x01, 0x6d, 0x12, 0x0c, 0x0a, 0x01, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x01, 0x6e, 0x22, 0xa3, 0x06, 0x0a, 0x1d, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x4d,
	0x41, 0x76, 0x67, 0x45, 0x6e, 0x64, 0x50, 0x72, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x12, 0x4a, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x36, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x62, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65,
	0x64, 0x61, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c,
	0x4d, 0x41, 0x76, 0x67, 0x45, 0x6e, 0x64, 0x50, 0x72, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x12, 0x5d, 0x0a, 0x0d, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74,
	0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x37, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x62, 0x69,
	0x6e, 0x61, 0x6e, 0x63, 0x65, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x74,
	0x65, 0x72, 0x76, 0x61, 0x6c, 0x4d, 0x41, 0x76, 0x67, 0x45, 0x6e, 0x64, 0x50, 0x72, 0x69, 0x63,
	0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x32,
	0x52, 0x0d, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x30, 0x0a, 0x13, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4f, 0x72, 0x64, 0x65,
	0x72, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x13, 0x6f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x54, 0x6f, 0x74, 0x61,
	0x6c, 0x12, 0x2a, 0x0a, 0x10, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x57, 0x69,
	0x6e, 0x52, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x6f, 0x70, 0x65,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x57, 0x69, 0x6e, 0x52, 0x61, 0x74, 0x65, 0x12, 0x2e, 0x0a,
	0x12, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x57, 0x69, 0x6e, 0x41, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x6f, 0x70, 0x65, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x57, 0x69, 0x6e, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x1a, 0xcc, 0x01,
	0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f,
	0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x6e, 0x64, 0x5f, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x50,
	0x72, 0x69, 0x63, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x6f, 0x70, 0x5f, 0x70, 0x72, 0x69, 0x63,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x6f, 0x70, 0x50, 0x72, 0x69, 0x63,
	0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6c, 0x6f, 0x77, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6c, 0x6f, 0x77, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x38,
	0x0a, 0x19, 0x77, 0x69, 0x74, 0x68, 0x5f, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x5f, 0x61, 0x76,
	0x67, 0x5f, 0x65, 0x6e, 0x64, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x15, 0x77, 0x69, 0x74, 0x68, 0x42, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x41, 0x76, 0x67,
	0x45, 0x6e, 0x64, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x1a, 0xf9, 0x01, 0x0a,
	0x05, 0x4c, 0x69, 0x73, 0x74, 0x32, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f,
	0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x6e, 0x64, 0x5f, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x50,
	0x72, 0x69, 0x63, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12,
	0x0a, 0x04, 0x72, 0x61, 0x74, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x72, 0x61,
	0x74, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x45, 0x6e, 0x64, 0x50, 0x72,
	0x69, 0x63, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x43, 0x6c, 0x6f, 0x73, 0x65,
	0x45, 0x6e, 0x64, 0x50, 0x72, 0x69, 0x63, 0x65, 0x32, 0xe3, 0x02, 0x0a, 0x0b, 0x42, 0x69, 0x6e,
	0x61, 0x6e, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x96, 0x01, 0x0a, 0x13, 0x44, 0x6f, 0x77,
	0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x2e, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x62, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x64, 0x61,
	0x74, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x69,
	0x6e, 0x61, 0x6e, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x2c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x62, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x64, 0x61,
	0x74, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x42, 0x69,
	0x6e, 0x61, 0x6e, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x21,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1b, 0x12, 0x19, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x69, 0x6e,
	0x61, 0x6e, 0x63, 0x65, 0x64, 0x61, 0x74, 0x61, 0x2f, 0x64, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61,
	0x64, 0x12, 0xba, 0x01, 0x0a, 0x18, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x4d, 0x41,
	0x76, 0x67, 0x45, 0x6e, 0x64, 0x50, 0x72, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x33,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x62, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x4d, 0x41, 0x76, 0x67,
	0x45, 0x6e, 0x64, 0x50, 0x72, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x31, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x62, 0x69, 0x6e, 0x61, 0x6e, 0x63,
	0x65, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61,
	0x6c, 0x4d, 0x41, 0x76, 0x67, 0x45, 0x6e, 0x64, 0x50, 0x72, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x36, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x30, 0x12, 0x2e,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x64, 0x61, 0x74, 0x61,
	0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x5f, 0x6d, 0x5f, 0x61, 0x76, 0x67, 0x5f,
	0x65, 0x6e, 0x64, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x42, 0x39,
	0x0a, 0x12, 0x61, 0x70, 0x69, 0x2e, 0x62, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x64, 0x61, 0x74,
	0x61, 0x2e, 0x76, 0x31, 0x50, 0x01, 0x5a, 0x21, 0x62, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x64,
	0x61, 0x74, 0x61, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x65, 0x64,
	0x61, 0x74, 0x61, 0x2f, 0x76, 0x31, 0x3b, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_api_binancedata_v1_binancedata_proto_rawDescOnce sync.Once
	file_api_binancedata_v1_binancedata_proto_rawDescData = file_api_binancedata_v1_binancedata_proto_rawDesc
)

func file_api_binancedata_v1_binancedata_proto_rawDescGZIP() []byte {
	file_api_binancedata_v1_binancedata_proto_rawDescOnce.Do(func() {
		file_api_binancedata_v1_binancedata_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_binancedata_v1_binancedata_proto_rawDescData)
	})
	return file_api_binancedata_v1_binancedata_proto_rawDescData
}

var file_api_binancedata_v1_binancedata_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_api_binancedata_v1_binancedata_proto_goTypes = []interface{}{
	(*DownloadBinanceDataRequest)(nil),          // 0: api.binancedata.v1.DownloadBinanceDataRequest
	(*DownloadBinanceDataReply)(nil),            // 1: api.binancedata.v1.DownloadBinanceDataReply
	(*IntervalMAvgEndPriceDataRequest)(nil),     // 2: api.binancedata.v1.IntervalMAvgEndPriceDataRequest
	(*IntervalMAvgEndPriceDataReply)(nil),       // 3: api.binancedata.v1.IntervalMAvgEndPriceDataReply
	(*IntervalMAvgEndPriceDataReply_List)(nil),  // 4: api.binancedata.v1.IntervalMAvgEndPriceDataReply.List
	(*IntervalMAvgEndPriceDataReply_List2)(nil), // 5: api.binancedata.v1.IntervalMAvgEndPriceDataReply.List2
}
var file_api_binancedata_v1_binancedata_proto_depIdxs = []int32{
	4, // 0: api.binancedata.v1.IntervalMAvgEndPriceDataReply.data:type_name -> api.binancedata.v1.IntervalMAvgEndPriceDataReply.List
	5, // 1: api.binancedata.v1.IntervalMAvgEndPriceDataReply.operationData:type_name -> api.binancedata.v1.IntervalMAvgEndPriceDataReply.List2
	0, // 2: api.binancedata.v1.BinanceData.DownloadBinanceData:input_type -> api.binancedata.v1.DownloadBinanceDataRequest
	2, // 3: api.binancedata.v1.BinanceData.IntervalMAvgEndPriceData:input_type -> api.binancedata.v1.IntervalMAvgEndPriceDataRequest
	1, // 4: api.binancedata.v1.BinanceData.DownloadBinanceData:output_type -> api.binancedata.v1.DownloadBinanceDataReply
	3, // 5: api.binancedata.v1.BinanceData.IntervalMAvgEndPriceData:output_type -> api.binancedata.v1.IntervalMAvgEndPriceDataReply
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_binancedata_v1_binancedata_proto_init() }
func file_api_binancedata_v1_binancedata_proto_init() {
	if File_api_binancedata_v1_binancedata_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_binancedata_v1_binancedata_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownloadBinanceDataRequest); i {
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
		file_api_binancedata_v1_binancedata_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DownloadBinanceDataReply); i {
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
		file_api_binancedata_v1_binancedata_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IntervalMAvgEndPriceDataRequest); i {
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
		file_api_binancedata_v1_binancedata_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IntervalMAvgEndPriceDataReply); i {
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
		file_api_binancedata_v1_binancedata_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IntervalMAvgEndPriceDataReply_List); i {
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
		file_api_binancedata_v1_binancedata_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IntervalMAvgEndPriceDataReply_List2); i {
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
			RawDescriptor: file_api_binancedata_v1_binancedata_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_binancedata_v1_binancedata_proto_goTypes,
		DependencyIndexes: file_api_binancedata_v1_binancedata_proto_depIdxs,
		MessageInfos:      file_api_binancedata_v1_binancedata_proto_msgTypes,
	}.Build()
	File_api_binancedata_v1_binancedata_proto = out.File
	file_api_binancedata_v1_binancedata_proto_rawDesc = nil
	file_api_binancedata_v1_binancedata_proto_goTypes = nil
	file_api_binancedata_v1_binancedata_proto_depIdxs = nil
}
