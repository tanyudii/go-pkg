package protobuf

import (
	"encoding/json"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"pkg.tanyudii.me/go-pkg/go-mon/pointer"
	"time"
)

func ToPbTimestamp(val *time.Time) *timestamppb.Timestamp {
	if val == nil {
		return nil
	}
	return &timestamppb.Timestamp{Seconds: val.Unix()}
}

func ToTimePointer(val *timestamppb.Timestamp) *time.Time {
	if val == nil {
		return nil
	}
	return pointer.Val(val.AsTime())
}

func ToPbDoubleValue(val *float64) *wrappers.DoubleValue {
	if val == nil {
		return nil
	}
	return &wrappers.DoubleValue{
		Value: *val,
	}
}

func ToFloat64Pointer(val *wrappers.DoubleValue) *float64 {
	if val == nil {
		return nil
	}
	return pointer.Val(val.Value)
}

func ToPbFloatValue(val *float32) *wrappers.FloatValue {
	if val == nil {
		return nil
	}
	return &wrappers.FloatValue{
		Value: *val,
	}
}

func ToFloat32Pointer(val *wrappers.FloatValue) *float32 {
	if val == nil {
		return nil
	}
	return pointer.Val(val.Value)
}

func ToPbInt64Value(val *int64) *wrappers.Int64Value {
	if val == nil {
		return nil
	}
	return &wrappers.Int64Value{
		Value: *val,
	}
}

func ToInt64Pointer(val *wrappers.Int64Value) *int64 {
	if val == nil {
		return nil
	}
	return pointer.Val(val.Value)
}

func ToPbUInt64Value(val *uint64) *wrappers.UInt64Value {
	if val == nil {
		return nil
	}
	return &wrappers.UInt64Value{
		Value: *val,
	}
}

func ToUInt64Pointer(val *wrappers.UInt64Value) *uint64 {
	if val == nil {
		return nil
	}
	return pointer.Val(val.Value)
}

func ToPbInt32Value(val *int32) *wrappers.Int32Value {
	if val == nil {
		return nil
	}
	return &wrappers.Int32Value{
		Value: *val,
	}
}

func ToInt32Pointer(val *wrappers.Int32Value) *int32 {
	if val == nil {
		return nil
	}
	return pointer.Val(val.Value)
}

func ToPbUInt32Value(val *uint32) *wrappers.UInt32Value {
	if val == nil {
		return nil
	}
	return &wrappers.UInt32Value{
		Value: *val,
	}
}

func ToUInt32Pointer(val *wrappers.UInt32Value) *uint32 {
	if val == nil {
		return nil
	}
	return pointer.Val(val.Value)
}

func ToPbBoolValue(val *bool) *wrapperspb.BoolValue {
	if val == nil {
		return nil
	}
	return &wrapperspb.BoolValue{
		Value: *val,
	}
}

func ToBoolPointer(val *wrapperspb.BoolValue) *bool {
	if val == nil {
		return nil
	}
	return pointer.Val(val.Value)
}

func ToPbStringValue(val *string) *wrappers.StringValue {
	if val == nil {
		return nil
	}
	return &wrappers.StringValue{
		Value: *val,
	}
}

func ToPbStringValueNullable(val string) *wrappers.StringValue {
	if val == "" {
		return nil
	}
	return &wrappers.StringValue{
		Value: val,
	}
}

func ToStringPointer(val *wrappers.StringValue) *string {
	if val == nil {
		return nil
	}
	return pointer.Val(val.Value)
}

func TOPbBytesValue(val *[]byte) *wrappers.BytesValue {
	if val == nil {
		return nil
	}
	return &wrappers.BytesValue{
		Value: *val,
	}
}

func ToBytesPointer(val *wrappers.BytesValue) *[]byte {
	if val == nil {
		return nil
	}
	return pointer.Val(val.Value)
}

func ToPbStruct(val map[string]interface{}) *structpb.Struct {
	jsonVal, _ := json.Marshal(val)
	valPb := structpb.Struct{}
	_ = json.Unmarshal(jsonVal, &valPb)
	return &valPb
}

func ToMapInterface(val *structpb.Struct) map[string]interface{} {
	if val == nil {
		return nil
	}
	return val.AsMap()
}
