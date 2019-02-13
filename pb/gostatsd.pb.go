// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pb/gostatsd.proto

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EventV1_EventPriority int32

const (
	EventV1_Normal EventV1_EventPriority = 0
	EventV1_Low    EventV1_EventPriority = 1
)

var EventV1_EventPriority_name = map[int32]string{
	0: "Normal",
	1: "Low",
}
var EventV1_EventPriority_value = map[string]int32{
	"Normal": 0,
	"Low":    1,
}

func (x EventV1_EventPriority) String() string {
	return proto.EnumName(EventV1_EventPriority_name, int32(x))
}
func (EventV1_EventPriority) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_gostatsd_62c098c42c58da7d, []int{1, 0}
}

type EventV1_AlertType int32

const (
	EventV1_Info    EventV1_AlertType = 0
	EventV1_Warning EventV1_AlertType = 1
	EventV1_Error   EventV1_AlertType = 2
	EventV1_Success EventV1_AlertType = 3
)

var EventV1_AlertType_name = map[int32]string{
	0: "Info",
	1: "Warning",
	2: "Error",
	3: "Success",
}
var EventV1_AlertType_value = map[string]int32{
	"Info":    0,
	"Warning": 1,
	"Error":   2,
	"Success": 3,
}

func (x EventV1_AlertType) String() string {
	return proto.EnumName(EventV1_AlertType_name, int32(x))
}
func (EventV1_AlertType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_gostatsd_62c098c42c58da7d, []int{1, 1}
}

type RawMessageV1 struct {
	RawMetrics           []*RawMetricV1 `protobuf:"bytes,1,rep,name=RawMetrics,proto3" json:"RawMetrics,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *RawMessageV1) Reset()         { *m = RawMessageV1{} }
func (m *RawMessageV1) String() string { return proto.CompactTextString(m) }
func (*RawMessageV1) ProtoMessage()    {}
func (*RawMessageV1) Descriptor() ([]byte, []int) {
	return fileDescriptor_gostatsd_62c098c42c58da7d, []int{0}
}
func (m *RawMessageV1) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RawMessageV1.Unmarshal(m, b)
}
func (m *RawMessageV1) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RawMessageV1.Marshal(b, m, deterministic)
}
func (dst *RawMessageV1) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RawMessageV1.Merge(dst, src)
}
func (m *RawMessageV1) XXX_Size() int {
	return xxx_messageInfo_RawMessageV1.Size(m)
}
func (m *RawMessageV1) XXX_DiscardUnknown() {
	xxx_messageInfo_RawMessageV1.DiscardUnknown(m)
}

var xxx_messageInfo_RawMessageV1 proto.InternalMessageInfo

func (m *RawMessageV1) GetRawMetrics() []*RawMetricV1 {
	if m != nil {
		return m.RawMetrics
	}
	return nil
}

type EventV1 struct {
	Title                string                `protobuf:"bytes,1,opt,name=Title,proto3" json:"Title,omitempty"`
	Text                 string                `protobuf:"bytes,2,opt,name=Text,proto3" json:"Text,omitempty"`
	DateHappened         int64                 `protobuf:"varint,3,opt,name=DateHappened,proto3" json:"DateHappened,omitempty"`
	Hostname             string                `protobuf:"bytes,4,opt,name=Hostname,proto3" json:"Hostname,omitempty"`
	AggregationKey       string                `protobuf:"bytes,5,opt,name=AggregationKey,proto3" json:"AggregationKey,omitempty"`
	SourceTypeName       string                `protobuf:"bytes,6,opt,name=SourceTypeName,proto3" json:"SourceTypeName,omitempty"`
	Tags                 []string              `protobuf:"bytes,7,rep,name=Tags,proto3" json:"Tags,omitempty"`
	SourceIP             string                `protobuf:"bytes,8,opt,name=SourceIP,proto3" json:"SourceIP,omitempty"`
	Priority             EventV1_EventPriority `protobuf:"varint,9,opt,name=Priority,proto3,enum=pb.EventV1_EventPriority" json:"Priority,omitempty"`
	Type                 EventV1_AlertType     `protobuf:"varint,10,opt,name=Type,proto3,enum=pb.EventV1_AlertType" json:"Type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *EventV1) Reset()         { *m = EventV1{} }
func (m *EventV1) String() string { return proto.CompactTextString(m) }
func (*EventV1) ProtoMessage()    {}
func (*EventV1) Descriptor() ([]byte, []int) {
	return fileDescriptor_gostatsd_62c098c42c58da7d, []int{1}
}
func (m *EventV1) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventV1.Unmarshal(m, b)
}
func (m *EventV1) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventV1.Marshal(b, m, deterministic)
}
func (dst *EventV1) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventV1.Merge(dst, src)
}
func (m *EventV1) XXX_Size() int {
	return xxx_messageInfo_EventV1.Size(m)
}
func (m *EventV1) XXX_DiscardUnknown() {
	xxx_messageInfo_EventV1.DiscardUnknown(m)
}

var xxx_messageInfo_EventV1 proto.InternalMessageInfo

func (m *EventV1) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *EventV1) GetText() string {
	if m != nil {
		return m.Text
	}
	return ""
}

func (m *EventV1) GetDateHappened() int64 {
	if m != nil {
		return m.DateHappened
	}
	return 0
}

func (m *EventV1) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

func (m *EventV1) GetAggregationKey() string {
	if m != nil {
		return m.AggregationKey
	}
	return ""
}

func (m *EventV1) GetSourceTypeName() string {
	if m != nil {
		return m.SourceTypeName
	}
	return ""
}

func (m *EventV1) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *EventV1) GetSourceIP() string {
	if m != nil {
		return m.SourceIP
	}
	return ""
}

func (m *EventV1) GetPriority() EventV1_EventPriority {
	if m != nil {
		return m.Priority
	}
	return EventV1_Normal
}

func (m *EventV1) GetType() EventV1_AlertType {
	if m != nil {
		return m.Type
	}
	return EventV1_Info
}

type RawMetricV1 struct {
	Name     string   `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
	Tags     []string `protobuf:"bytes,2,rep,name=Tags,proto3" json:"Tags,omitempty"`
	Hostname string   `protobuf:"bytes,3,opt,name=Hostname,proto3" json:"Hostname,omitempty"`
	// Types that are valid to be assigned to M:
	//	*RawMetricV1_Counter
	//	*RawMetricV1_Gauge
	//	*RawMetricV1_Set
	//	*RawMetricV1_Timer
	M                    isRawMetricV1_M `protobuf_oneof:"M"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *RawMetricV1) Reset()         { *m = RawMetricV1{} }
func (m *RawMetricV1) String() string { return proto.CompactTextString(m) }
func (*RawMetricV1) ProtoMessage()    {}
func (*RawMetricV1) Descriptor() ([]byte, []int) {
	return fileDescriptor_gostatsd_62c098c42c58da7d, []int{2}
}
func (m *RawMetricV1) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RawMetricV1.Unmarshal(m, b)
}
func (m *RawMetricV1) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RawMetricV1.Marshal(b, m, deterministic)
}
func (dst *RawMetricV1) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RawMetricV1.Merge(dst, src)
}
func (m *RawMetricV1) XXX_Size() int {
	return xxx_messageInfo_RawMetricV1.Size(m)
}
func (m *RawMetricV1) XXX_DiscardUnknown() {
	xxx_messageInfo_RawMetricV1.DiscardUnknown(m)
}

var xxx_messageInfo_RawMetricV1 proto.InternalMessageInfo

func (m *RawMetricV1) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *RawMetricV1) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func (m *RawMetricV1) GetHostname() string {
	if m != nil {
		return m.Hostname
	}
	return ""
}

type isRawMetricV1_M interface {
	isRawMetricV1_M()
}

type RawMetricV1_Counter struct {
	Counter *RawCounterV1 `protobuf:"bytes,12,opt,name=counter,proto3,oneof"`
}

type RawMetricV1_Gauge struct {
	Gauge *RawGaugeV1 `protobuf:"bytes,13,opt,name=gauge,proto3,oneof"`
}

type RawMetricV1_Set struct {
	Set *RawSetV1 `protobuf:"bytes,14,opt,name=set,proto3,oneof"`
}

type RawMetricV1_Timer struct {
	Timer *RawTimerV1 `protobuf:"bytes,15,opt,name=timer,proto3,oneof"`
}

func (*RawMetricV1_Counter) isRawMetricV1_M() {}

func (*RawMetricV1_Gauge) isRawMetricV1_M() {}

func (*RawMetricV1_Set) isRawMetricV1_M() {}

func (*RawMetricV1_Timer) isRawMetricV1_M() {}

func (m *RawMetricV1) GetM() isRawMetricV1_M {
	if m != nil {
		return m.M
	}
	return nil
}

func (m *RawMetricV1) GetCounter() *RawCounterV1 {
	if x, ok := m.GetM().(*RawMetricV1_Counter); ok {
		return x.Counter
	}
	return nil
}

func (m *RawMetricV1) GetGauge() *RawGaugeV1 {
	if x, ok := m.GetM().(*RawMetricV1_Gauge); ok {
		return x.Gauge
	}
	return nil
}

func (m *RawMetricV1) GetSet() *RawSetV1 {
	if x, ok := m.GetM().(*RawMetricV1_Set); ok {
		return x.Set
	}
	return nil
}

func (m *RawMetricV1) GetTimer() *RawTimerV1 {
	if x, ok := m.GetM().(*RawMetricV1_Timer); ok {
		return x.Timer
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*RawMetricV1) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _RawMetricV1_OneofMarshaler, _RawMetricV1_OneofUnmarshaler, _RawMetricV1_OneofSizer, []interface{}{
		(*RawMetricV1_Counter)(nil),
		(*RawMetricV1_Gauge)(nil),
		(*RawMetricV1_Set)(nil),
		(*RawMetricV1_Timer)(nil),
	}
}

func _RawMetricV1_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*RawMetricV1)
	// M
	switch x := m.M.(type) {
	case *RawMetricV1_Counter:
		b.EncodeVarint(12<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Counter); err != nil {
			return err
		}
	case *RawMetricV1_Gauge:
		b.EncodeVarint(13<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Gauge); err != nil {
			return err
		}
	case *RawMetricV1_Set:
		b.EncodeVarint(14<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Set); err != nil {
			return err
		}
	case *RawMetricV1_Timer:
		b.EncodeVarint(15<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Timer); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("RawMetricV1.M has unexpected type %T", x)
	}
	return nil
}

func _RawMetricV1_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*RawMetricV1)
	switch tag {
	case 12: // M.counter
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RawCounterV1)
		err := b.DecodeMessage(msg)
		m.M = &RawMetricV1_Counter{msg}
		return true, err
	case 13: // M.gauge
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RawGaugeV1)
		err := b.DecodeMessage(msg)
		m.M = &RawMetricV1_Gauge{msg}
		return true, err
	case 14: // M.set
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RawSetV1)
		err := b.DecodeMessage(msg)
		m.M = &RawMetricV1_Set{msg}
		return true, err
	case 15: // M.timer
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(RawTimerV1)
		err := b.DecodeMessage(msg)
		m.M = &RawMetricV1_Timer{msg}
		return true, err
	default:
		return false, nil
	}
}

func _RawMetricV1_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*RawMetricV1)
	// M
	switch x := m.M.(type) {
	case *RawMetricV1_Counter:
		s := proto.Size(x.Counter)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *RawMetricV1_Gauge:
		s := proto.Size(x.Gauge)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *RawMetricV1_Set:
		s := proto.Size(x.Set)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *RawMetricV1_Timer:
		s := proto.Size(x.Timer)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type RawCounterV1 struct {
	Value                float64  `protobuf:"fixed64,1,opt,name=Value,proto3" json:"Value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RawCounterV1) Reset()         { *m = RawCounterV1{} }
func (m *RawCounterV1) String() string { return proto.CompactTextString(m) }
func (*RawCounterV1) ProtoMessage()    {}
func (*RawCounterV1) Descriptor() ([]byte, []int) {
	return fileDescriptor_gostatsd_62c098c42c58da7d, []int{3}
}
func (m *RawCounterV1) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RawCounterV1.Unmarshal(m, b)
}
func (m *RawCounterV1) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RawCounterV1.Marshal(b, m, deterministic)
}
func (dst *RawCounterV1) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RawCounterV1.Merge(dst, src)
}
func (m *RawCounterV1) XXX_Size() int {
	return xxx_messageInfo_RawCounterV1.Size(m)
}
func (m *RawCounterV1) XXX_DiscardUnknown() {
	xxx_messageInfo_RawCounterV1.DiscardUnknown(m)
}

var xxx_messageInfo_RawCounterV1 proto.InternalMessageInfo

func (m *RawCounterV1) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type RawGaugeV1 struct {
	Value                float64  `protobuf:"fixed64,1,opt,name=Value,proto3" json:"Value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RawGaugeV1) Reset()         { *m = RawGaugeV1{} }
func (m *RawGaugeV1) String() string { return proto.CompactTextString(m) }
func (*RawGaugeV1) ProtoMessage()    {}
func (*RawGaugeV1) Descriptor() ([]byte, []int) {
	return fileDescriptor_gostatsd_62c098c42c58da7d, []int{4}
}
func (m *RawGaugeV1) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RawGaugeV1.Unmarshal(m, b)
}
func (m *RawGaugeV1) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RawGaugeV1.Marshal(b, m, deterministic)
}
func (dst *RawGaugeV1) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RawGaugeV1.Merge(dst, src)
}
func (m *RawGaugeV1) XXX_Size() int {
	return xxx_messageInfo_RawGaugeV1.Size(m)
}
func (m *RawGaugeV1) XXX_DiscardUnknown() {
	xxx_messageInfo_RawGaugeV1.DiscardUnknown(m)
}

var xxx_messageInfo_RawGaugeV1 proto.InternalMessageInfo

func (m *RawGaugeV1) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type RawSetV1 struct {
	Value                string   `protobuf:"bytes,1,opt,name=Value,proto3" json:"Value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RawSetV1) Reset()         { *m = RawSetV1{} }
func (m *RawSetV1) String() string { return proto.CompactTextString(m) }
func (*RawSetV1) ProtoMessage()    {}
func (*RawSetV1) Descriptor() ([]byte, []int) {
	return fileDescriptor_gostatsd_62c098c42c58da7d, []int{5}
}
func (m *RawSetV1) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RawSetV1.Unmarshal(m, b)
}
func (m *RawSetV1) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RawSetV1.Marshal(b, m, deterministic)
}
func (dst *RawSetV1) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RawSetV1.Merge(dst, src)
}
func (m *RawSetV1) XXX_Size() int {
	return xxx_messageInfo_RawSetV1.Size(m)
}
func (m *RawSetV1) XXX_DiscardUnknown() {
	xxx_messageInfo_RawSetV1.DiscardUnknown(m)
}

var xxx_messageInfo_RawSetV1 proto.InternalMessageInfo

func (m *RawSetV1) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type RawTimerV1 struct {
	Value                float64  `protobuf:"fixed64,1,opt,name=Value,proto3" json:"Value,omitempty"`
	Rate                 float64  `protobuf:"fixed64,2,opt,name=Rate,proto3" json:"Rate,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RawTimerV1) Reset()         { *m = RawTimerV1{} }
func (m *RawTimerV1) String() string { return proto.CompactTextString(m) }
func (*RawTimerV1) ProtoMessage()    {}
func (*RawTimerV1) Descriptor() ([]byte, []int) {
	return fileDescriptor_gostatsd_62c098c42c58da7d, []int{6}
}
func (m *RawTimerV1) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RawTimerV1.Unmarshal(m, b)
}
func (m *RawTimerV1) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RawTimerV1.Marshal(b, m, deterministic)
}
func (dst *RawTimerV1) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RawTimerV1.Merge(dst, src)
}
func (m *RawTimerV1) XXX_Size() int {
	return xxx_messageInfo_RawTimerV1.Size(m)
}
func (m *RawTimerV1) XXX_DiscardUnknown() {
	xxx_messageInfo_RawTimerV1.DiscardUnknown(m)
}

var xxx_messageInfo_RawTimerV1 proto.InternalMessageInfo

func (m *RawTimerV1) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

func (m *RawTimerV1) GetRate() float64 {
	if m != nil {
		return m.Rate
	}
	return 0
}

func init() {
	proto.RegisterType((*RawMessageV1)(nil), "pb.RawMessageV1")
	proto.RegisterType((*EventV1)(nil), "pb.EventV1")
	proto.RegisterType((*RawMetricV1)(nil), "pb.RawMetricV1")
	proto.RegisterType((*RawCounterV1)(nil), "pb.RawCounterV1")
	proto.RegisterType((*RawGaugeV1)(nil), "pb.RawGaugeV1")
	proto.RegisterType((*RawSetV1)(nil), "pb.RawSetV1")
	proto.RegisterType((*RawTimerV1)(nil), "pb.RawTimerV1")
	proto.RegisterEnum("pb.EventV1_EventPriority", EventV1_EventPriority_name, EventV1_EventPriority_value)
	proto.RegisterEnum("pb.EventV1_AlertType", EventV1_AlertType_name, EventV1_AlertType_value)
}

func init() { proto.RegisterFile("pb/gostatsd.proto", fileDescriptor_gostatsd_62c098c42c58da7d) }

var fileDescriptor_gostatsd_62c098c42c58da7d = []byte{
	// 516 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x93, 0x4f, 0x6f, 0xd3, 0x4c,
	0x10, 0xc6, 0xbb, 0x71, 0xfe, 0x4e, 0xd2, 0xd4, 0xef, 0xea, 0x45, 0x5a, 0x38, 0x59, 0x56, 0x55,
	0x05, 0x09, 0xa5, 0x4a, 0x10, 0x1c, 0xb8, 0xa0, 0x02, 0x15, 0xa9, 0xa0, 0x55, 0xe5, 0x44, 0xe1,
	0xbc, 0x49, 0x07, 0xcb, 0x52, 0xe2, 0xb5, 0x76, 0x37, 0x94, 0x7c, 0x32, 0xbe, 0x19, 0x67, 0xb4,
	0xb3, 0x8e, 0x15, 0x57, 0x70, 0x8a, 0xe7, 0x99, 0xdf, 0xcc, 0x4e, 0x9e, 0xd9, 0x85, 0xff, 0x8a,
	0xd5, 0x65, 0xaa, 0x8c, 0x95, 0xd6, 0x3c, 0x8c, 0x0b, 0xad, 0xac, 0xe2, 0x8d, 0x62, 0x15, 0xbf,
	0x87, 0x41, 0x22, 0x1f, 0x6f, 0xd1, 0x18, 0x99, 0xe2, 0x72, 0xc2, 0x2f, 0x01, 0x28, 0xb6, 0x3a,
	0x5b, 0x1b, 0xc1, 0xa2, 0x60, 0xd4, 0x9f, 0x9e, 0x8d, 0x8b, 0xd5, 0xb8, 0x52, 0x97, 0x93, 0xe4,
	0x08, 0x89, 0x7f, 0x05, 0xd0, 0xb9, 0xfe, 0x81, 0xb9, 0x5d, 0x4e, 0xf8, 0xff, 0xd0, 0x5a, 0x64,
	0x76, 0x83, 0x82, 0x45, 0x6c, 0xd4, 0x4b, 0x7c, 0xc0, 0x39, 0x34, 0x17, 0xf8, 0xd3, 0x8a, 0x06,
	0x89, 0xf4, 0xcd, 0x63, 0x18, 0x7c, 0x92, 0x16, 0x67, 0xb2, 0x28, 0x30, 0xc7, 0x07, 0x11, 0x44,
	0x6c, 0x14, 0x24, 0x35, 0x8d, 0xbf, 0x80, 0xee, 0x4c, 0x19, 0x9b, 0xcb, 0x2d, 0x8a, 0x26, 0xd5,
	0x56, 0x31, 0xbf, 0x80, 0xe1, 0x55, 0x9a, 0x6a, 0x4c, 0xa5, 0xcd, 0x54, 0xfe, 0x05, 0xf7, 0xa2,
	0x45, 0xc4, 0x13, 0xd5, 0x71, 0x73, 0xb5, 0xd3, 0x6b, 0x5c, 0xec, 0x0b, 0xbc, 0x73, 0x9d, 0xda,
	0x9e, 0xab, 0xab, 0x34, 0xa3, 0x4c, 0x8d, 0xe8, 0x44, 0x01, 0xcd, 0x28, 0x53, 0xe3, 0xce, 0xf7,
	0xd4, 0xcd, 0xbd, 0xe8, 0xfa, 0xf3, 0x0f, 0x31, 0x7f, 0x03, 0xdd, 0x7b, 0x9d, 0x29, 0x9d, 0xd9,
	0xbd, 0xe8, 0x45, 0x6c, 0x34, 0x9c, 0x3e, 0x77, 0x26, 0x95, 0x46, 0xf8, 0xdf, 0x03, 0x90, 0x54,
	0x28, 0x7f, 0x09, 0x4d, 0x77, 0xa4, 0x00, 0x2a, 0x79, 0x76, 0x5c, 0x72, 0xb5, 0x41, 0x6d, 0x5d,
	0x32, 0x21, 0x24, 0x3e, 0x87, 0xd3, 0x5a, 0x17, 0x0e, 0xd0, 0xbe, 0x53, 0x7a, 0x2b, 0x37, 0xe1,
	0x09, 0xef, 0x40, 0xf0, 0x55, 0x3d, 0x86, 0x2c, 0x7e, 0x07, 0xbd, 0xaa, 0x90, 0x77, 0xa1, 0x79,
	0x93, 0x7f, 0x57, 0xe1, 0x09, 0xef, 0x43, 0xe7, 0x9b, 0xd4, 0x79, 0x96, 0xa7, 0x21, 0xe3, 0x3d,
	0x68, 0x5d, 0x6b, 0xad, 0x74, 0xd8, 0x70, 0xfa, 0x7c, 0xb7, 0x5e, 0xa3, 0x31, 0x61, 0x10, 0xff,
	0x66, 0xd0, 0x3f, 0xda, 0xaa, 0xf3, 0x80, 0x1c, 0xf2, 0xcb, 0x6b, 0xd6, 0x7c, 0x69, 0xd4, 0x7d,
	0xa9, 0xf6, 0x12, 0x3c, 0xd9, 0xcb, 0x2b, 0xe8, 0xac, 0xd5, 0x2e, 0xb7, 0xa8, 0xc5, 0x20, 0x62,
	0xa3, 0xfe, 0x34, 0x2c, 0xef, 0xce, 0x47, 0xaf, 0x2e, 0x27, 0xb3, 0x93, 0xe4, 0x80, 0xf0, 0x0b,
	0x68, 0xa5, 0x72, 0x97, 0xa2, 0x38, 0x25, 0x76, 0x58, 0xb2, 0x9f, 0x9d, 0x46, 0xa4, 0x4f, 0xf3,
	0x08, 0x02, 0x83, 0x56, 0x0c, 0x89, 0x1a, 0x94, 0xd4, 0x1c, 0x2d, 0x31, 0x2e, 0xe5, 0x3a, 0xd9,
	0x6c, 0x8b, 0x5a, 0x9c, 0xd5, 0x3a, 0x2d, 0x9c, 0xe6, 0x3b, 0x51, 0xfa, 0x43, 0x00, 0xec, 0x36,
	0x3e, 0xa7, 0x3b, 0x5f, 0x4d, 0xe4, 0xae, 0xed, 0x52, 0x6e, 0x76, 0xfe, 0x9f, 0xb3, 0xc4, 0x07,
	0x71, 0x4c, 0x2f, 0xa1, 0x9c, 0xe5, 0x1f, 0x4c, 0x04, 0xdd, 0xc3, 0x24, 0x75, 0xa2, 0x77, 0x20,
	0xde, 0x52, 0x97, 0x72, 0x8e, 0xbf, 0x77, 0x71, 0x26, 0x27, 0xd2, 0x22, 0x3d, 0x10, 0x96, 0xd0,
	0xf7, 0xaa, 0x4d, 0x4f, 0xf4, 0xf5, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xae, 0xfb, 0x54, 0x5c,
	0xb7, 0x03, 0x00, 0x00,
}
