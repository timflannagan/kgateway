// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/solo-io/solo-projects/projects/glooshot/api/v1/glooshot.proto

package v1 // import "github.com/solo-io/solo-projects/projects/glooshot/pkg/api/v1"

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import _ "github.com/gogo/protobuf/types"
import v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
import faultinjection "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/faultinjection"
import core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

import time "time"

import bytes "bytes"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// TODO: this should probably go in the status. consult with ilackarms
type ExperimentSpec_State int32

const (
	ExperimentSpec_Init    ExperimentSpec_State = 0
	ExperimentSpec_Running ExperimentSpec_State = 1
	ExperimentSpec_Done    ExperimentSpec_State = 2
)

var ExperimentSpec_State_name = map[int32]string{
	0: "Init",
	1: "Running",
	2: "Done",
}
var ExperimentSpec_State_value = map[string]int32{
	"Init":    0,
	"Running": 1,
	"Done":    2,
}

func (x ExperimentSpec_State) String() string {
	return proto.EnumName(ExperimentSpec_State_name, int32(x))
}
func (ExperimentSpec_State) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_glooshot_c04baf15aaed035e, []int{1, 0}
}

//
// @solo-kit:resource.short_name=exp
// @solo-kit:resource.plural_name=experiments
// @solo-kit:resource.resource_groups=api.glooshot.solo.io
type Experiment struct {
	// Metadata contains the object metadata for this resource
	Metadata core.Metadata   `protobuf:"bytes,1,opt,name=metadata" json:"metadata"`
	Spec     *ExperimentSpec `protobuf:"bytes,2,opt,name=spec" json:"spec,omitempty"`
	// Status indicates the validation status of the resource. Status is read-only by clients, and set by gloo during validation
	Status               core.Status `protobuf:"bytes,3,opt,name=status" json:"status" testdiff:"ignore"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Experiment) Reset()         { *m = Experiment{} }
func (m *Experiment) String() string { return proto.CompactTextString(m) }
func (*Experiment) ProtoMessage()    {}
func (*Experiment) Descriptor() ([]byte, []int) {
	return fileDescriptor_glooshot_c04baf15aaed035e, []int{0}
}
func (m *Experiment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Experiment.Unmarshal(m, b)
}
func (m *Experiment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Experiment.Marshal(b, m, deterministic)
}
func (dst *Experiment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Experiment.Merge(dst, src)
}
func (m *Experiment) XXX_Size() int {
	return xxx_messageInfo_Experiment.Size(m)
}
func (m *Experiment) XXX_DiscardUnknown() {
	xxx_messageInfo_Experiment.DiscardUnknown(m)
}

var xxx_messageInfo_Experiment proto.InternalMessageInfo

func (m *Experiment) GetMetadata() core.Metadata {
	if m != nil {
		return m.Metadata
	}
	return core.Metadata{}
}

func (m *Experiment) GetSpec() *ExperimentSpec {
	if m != nil {
		return m.Spec
	}
	return nil
}

func (m *Experiment) GetStatus() core.Status {
	if m != nil {
		return m.Status
	}
	return core.Status{}
}

type ExperimentSpec struct {
	Faults        []*ExperimentSpec_InjectedFault `protobuf:"bytes,1,rep,name=faults" json:"faults,omitempty"`
	StopCondition *StopCondition                  `protobuf:"bytes,2,opt,name=stop_condition,json=stopCondition" json:"stop_condition,omitempty"`
	// State is the enum indicating the state of the resource
	State                ExperimentSpec_State `protobuf:"varint,4,opt,name=state,proto3,enum=glooshot.solo.io.ExperimentSpec_State" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ExperimentSpec) Reset()         { *m = ExperimentSpec{} }
func (m *ExperimentSpec) String() string { return proto.CompactTextString(m) }
func (*ExperimentSpec) ProtoMessage()    {}
func (*ExperimentSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_glooshot_c04baf15aaed035e, []int{1}
}
func (m *ExperimentSpec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExperimentSpec.Unmarshal(m, b)
}
func (m *ExperimentSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExperimentSpec.Marshal(b, m, deterministic)
}
func (dst *ExperimentSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExperimentSpec.Merge(dst, src)
}
func (m *ExperimentSpec) XXX_Size() int {
	return xxx_messageInfo_ExperimentSpec.Size(m)
}
func (m *ExperimentSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_ExperimentSpec.DiscardUnknown(m)
}

var xxx_messageInfo_ExperimentSpec proto.InternalMessageInfo

func (m *ExperimentSpec) GetFaults() []*ExperimentSpec_InjectedFault {
	if m != nil {
		return m.Faults
	}
	return nil
}

func (m *ExperimentSpec) GetStopCondition() *StopCondition {
	if m != nil {
		return m.StopCondition
	}
	return nil
}

func (m *ExperimentSpec) GetState() ExperimentSpec_State {
	if m != nil {
		return m.State
	}
	return ExperimentSpec_Init
}

type ExperimentSpec_InjectedFault struct {
	// TODO(yuval-k) should this be an upstream ref?
	Service              *v1.Destination             `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
	Fault                *faultinjection.RouteFaults `protobuf:"bytes,2,opt,name=fault" json:"fault,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *ExperimentSpec_InjectedFault) Reset()         { *m = ExperimentSpec_InjectedFault{} }
func (m *ExperimentSpec_InjectedFault) String() string { return proto.CompactTextString(m) }
func (*ExperimentSpec_InjectedFault) ProtoMessage()    {}
func (*ExperimentSpec_InjectedFault) Descriptor() ([]byte, []int) {
	return fileDescriptor_glooshot_c04baf15aaed035e, []int{1, 0}
}
func (m *ExperimentSpec_InjectedFault) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExperimentSpec_InjectedFault.Unmarshal(m, b)
}
func (m *ExperimentSpec_InjectedFault) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExperimentSpec_InjectedFault.Marshal(b, m, deterministic)
}
func (dst *ExperimentSpec_InjectedFault) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExperimentSpec_InjectedFault.Merge(dst, src)
}
func (m *ExperimentSpec_InjectedFault) XXX_Size() int {
	return xxx_messageInfo_ExperimentSpec_InjectedFault.Size(m)
}
func (m *ExperimentSpec_InjectedFault) XXX_DiscardUnknown() {
	xxx_messageInfo_ExperimentSpec_InjectedFault.DiscardUnknown(m)
}

var xxx_messageInfo_ExperimentSpec_InjectedFault proto.InternalMessageInfo

func (m *ExperimentSpec_InjectedFault) GetService() *v1.Destination {
	if m != nil {
		return m.Service
	}
	return nil
}

func (m *ExperimentSpec_InjectedFault) GetFault() *faultinjection.RouteFaults {
	if m != nil {
		return m.Fault
	}
	return nil
}

type StopCondition struct {
	Duration             *time.Duration     `protobuf:"bytes,1,opt,name=duration,stdduration" json:"duration,omitempty"`
	Metric               []*MetricThreshold `protobuf:"bytes,2,rep,name=metric" json:"metric,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *StopCondition) Reset()         { *m = StopCondition{} }
func (m *StopCondition) String() string { return proto.CompactTextString(m) }
func (*StopCondition) ProtoMessage()    {}
func (*StopCondition) Descriptor() ([]byte, []int) {
	return fileDescriptor_glooshot_c04baf15aaed035e, []int{2}
}
func (m *StopCondition) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopCondition.Unmarshal(m, b)
}
func (m *StopCondition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopCondition.Marshal(b, m, deterministic)
}
func (dst *StopCondition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopCondition.Merge(dst, src)
}
func (m *StopCondition) XXX_Size() int {
	return xxx_messageInfo_StopCondition.Size(m)
}
func (m *StopCondition) XXX_DiscardUnknown() {
	xxx_messageInfo_StopCondition.DiscardUnknown(m)
}

var xxx_messageInfo_StopCondition proto.InternalMessageInfo

func (m *StopCondition) GetDuration() *time.Duration {
	if m != nil {
		return m.Duration
	}
	return nil
}

func (m *StopCondition) GetMetric() []*MetricThreshold {
	if m != nil {
		return m.Metric
	}
	return nil
}

type MetricThreshold struct {
	MetricName           string   `protobuf:"bytes,1,opt,name=metric_name,json=metricName,proto3" json:"metric_name,omitempty"`
	Value                float64  `protobuf:"fixed64,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MetricThreshold) Reset()         { *m = MetricThreshold{} }
func (m *MetricThreshold) String() string { return proto.CompactTextString(m) }
func (*MetricThreshold) ProtoMessage()    {}
func (*MetricThreshold) Descriptor() ([]byte, []int) {
	return fileDescriptor_glooshot_c04baf15aaed035e, []int{3}
}
func (m *MetricThreshold) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MetricThreshold.Unmarshal(m, b)
}
func (m *MetricThreshold) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MetricThreshold.Marshal(b, m, deterministic)
}
func (dst *MetricThreshold) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MetricThreshold.Merge(dst, src)
}
func (m *MetricThreshold) XXX_Size() int {
	return xxx_messageInfo_MetricThreshold.Size(m)
}
func (m *MetricThreshold) XXX_DiscardUnknown() {
	xxx_messageInfo_MetricThreshold.DiscardUnknown(m)
}

var xxx_messageInfo_MetricThreshold proto.InternalMessageInfo

func (m *MetricThreshold) GetMetricName() string {
	if m != nil {
		return m.MetricName
	}
	return ""
}

func (m *MetricThreshold) GetValue() float64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type ChaosSettings struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChaosSettings) Reset()         { *m = ChaosSettings{} }
func (m *ChaosSettings) String() string { return proto.CompactTextString(m) }
func (*ChaosSettings) ProtoMessage()    {}
func (*ChaosSettings) Descriptor() ([]byte, []int) {
	return fileDescriptor_glooshot_c04baf15aaed035e, []int{4}
}
func (m *ChaosSettings) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaosSettings.Unmarshal(m, b)
}
func (m *ChaosSettings) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaosSettings.Marshal(b, m, deterministic)
}
func (dst *ChaosSettings) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaosSettings.Merge(dst, src)
}
func (m *ChaosSettings) XXX_Size() int {
	return xxx_messageInfo_ChaosSettings.Size(m)
}
func (m *ChaosSettings) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaosSettings.DiscardUnknown(m)
}

var xxx_messageInfo_ChaosSettings proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Experiment)(nil), "glooshot.solo.io.Experiment")
	proto.RegisterType((*ExperimentSpec)(nil), "glooshot.solo.io.ExperimentSpec")
	proto.RegisterType((*ExperimentSpec_InjectedFault)(nil), "glooshot.solo.io.ExperimentSpec.InjectedFault")
	proto.RegisterType((*StopCondition)(nil), "glooshot.solo.io.StopCondition")
	proto.RegisterType((*MetricThreshold)(nil), "glooshot.solo.io.MetricThreshold")
	proto.RegisterType((*ChaosSettings)(nil), "glooshot.solo.io.ChaosSettings")
	proto.RegisterEnum("glooshot.solo.io.ExperimentSpec_State", ExperimentSpec_State_name, ExperimentSpec_State_value)
}
func (this *Experiment) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Experiment)
	if !ok {
		that2, ok := that.(Experiment)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Metadata.Equal(&that1.Metadata) {
		return false
	}
	if !this.Spec.Equal(that1.Spec) {
		return false
	}
	if !this.Status.Equal(&that1.Status) {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *ExperimentSpec) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ExperimentSpec)
	if !ok {
		that2, ok := that.(ExperimentSpec)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.Faults) != len(that1.Faults) {
		return false
	}
	for i := range this.Faults {
		if !this.Faults[i].Equal(that1.Faults[i]) {
			return false
		}
	}
	if !this.StopCondition.Equal(that1.StopCondition) {
		return false
	}
	if this.State != that1.State {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *ExperimentSpec_InjectedFault) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ExperimentSpec_InjectedFault)
	if !ok {
		that2, ok := that.(ExperimentSpec_InjectedFault)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Service.Equal(that1.Service) {
		return false
	}
	if !this.Fault.Equal(that1.Fault) {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *StopCondition) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*StopCondition)
	if !ok {
		that2, ok := that.(StopCondition)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.Duration != nil && that1.Duration != nil {
		if *this.Duration != *that1.Duration {
			return false
		}
	} else if this.Duration != nil {
		return false
	} else if that1.Duration != nil {
		return false
	}
	if len(this.Metric) != len(that1.Metric) {
		return false
	}
	for i := range this.Metric {
		if !this.Metric[i].Equal(that1.Metric[i]) {
			return false
		}
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *MetricThreshold) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MetricThreshold)
	if !ok {
		that2, ok := that.(MetricThreshold)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.MetricName != that1.MetricName {
		return false
	}
	if this.Value != that1.Value {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *ChaosSettings) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ChaosSettings)
	if !ok {
		that2, ok := that.(ChaosSettings)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}

func init() {
	proto.RegisterFile("github.com/solo-io/solo-projects/projects/glooshot/api/v1/glooshot.proto", fileDescriptor_glooshot_c04baf15aaed035e)
}

var fileDescriptor_glooshot_c04baf15aaed035e = []byte{
	// 616 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0xcf, 0x52, 0x13, 0x4f,
	0x10, 0xc7, 0x59, 0x48, 0x02, 0xbf, 0x4e, 0x05, 0xf2, 0x9b, 0xa2, 0xac, 0xc0, 0x81, 0xc4, 0x3d,
	0x68, 0x2e, 0xee, 0x16, 0xe0, 0x01, 0xff, 0x70, 0x09, 0x88, 0x70, 0x40, 0xab, 0x26, 0x9e, 0xbc,
	0x50, 0xcb, 0x66, 0xb2, 0x19, 0xd9, 0x4c, 0x6f, 0xed, 0xcc, 0x52, 0x78, 0xf6, 0xe0, 0xcd, 0xb3,
	0x8f, 0xe0, 0x8b, 0x58, 0xe5, 0x53, 0x60, 0x95, 0x8f, 0xe0, 0x13, 0x58, 0xf3, 0x67, 0x17, 0x03,
	0x5a, 0xa8, 0xa7, 0x6c, 0x4f, 0xf7, 0xa7, 0xfb, 0xdb, 0x3d, 0xd3, 0x81, 0xc3, 0x84, 0xab, 0x49,
	0x71, 0x1a, 0xc4, 0x38, 0x0d, 0x25, 0xa6, 0xf8, 0x80, 0xa3, 0xfd, 0xcd, 0x72, 0x7c, 0xc3, 0x62,
	0x25, 0xc3, 0xea, 0x23, 0x49, 0x11, 0xe5, 0x04, 0x55, 0x18, 0x65, 0x3c, 0x3c, 0xdf, 0xac, 0xec,
	0x20, 0xcb, 0x51, 0x21, 0x69, 0x57, 0xb6, 0xe6, 0x03, 0x8e, 0xeb, 0xab, 0x09, 0x26, 0x68, 0x9c,
	0xa1, 0xfe, 0xb2, 0x71, 0xeb, 0x1b, 0x09, 0x62, 0x92, 0xb2, 0xd0, 0x58, 0xa7, 0xc5, 0x38, 0x1c,
	0x15, 0x79, 0xa4, 0x38, 0x0a, 0xe7, 0x0f, 0x7f, 0xa7, 0xe8, 0x8c, 0x57, 0xa5, 0xa5, 0x8a, 0x54,
	0x21, 0x1d, 0xb0, 0xf9, 0x07, 0xc0, 0x94, 0xa9, 0x68, 0x14, 0xa9, 0xc8, 0x21, 0x3b, 0xbf, 0x40,
	0xb4, 0xfc, 0xd9, 0x66, 0x4b, 0x38, 0xcb, 0xf1, 0xe2, 0xad, 0x23, 0x5f, 0xfe, 0x1d, 0x99, 0x16,
	0x09, 0x17, 0x32, 0x1c, 0x47, 0x45, 0xaa, 0xb8, 0xd0, 0x01, 0x1c, 0x85, 0x35, 0x6d, 0x42, 0xff,
	0xb3, 0x07, 0xf0, 0xec, 0x22, 0x63, 0x39, 0x9f, 0x32, 0xa1, 0xc8, 0x0e, 0x2c, 0x95, 0x5a, 0x3b,
	0x5e, 0xcf, 0xeb, 0x37, 0xb7, 0xee, 0x04, 0x31, 0xe6, 0xac, 0x1c, 0x6a, 0x70, 0xec, 0xbc, 0x83,
	0xda, 0x97, 0xcb, 0xee, 0x1c, 0xad, 0xa2, 0xc9, 0x43, 0xa8, 0xc9, 0x8c, 0xc5, 0x9d, 0x79, 0x43,
	0xf5, 0x82, 0xeb, 0xd7, 0x11, 0x5c, 0x55, 0x19, 0x66, 0x2c, 0xa6, 0x26, 0x9a, 0x3c, 0x87, 0x86,
	0x1d, 0x66, 0x67, 0xc1, 0x70, 0xab, 0xb3, 0xd5, 0x86, 0xc6, 0x37, 0x58, 0xd3, 0xb5, 0xbe, 0x5f,
	0x76, 0xff, 0x57, 0x4c, 0xaa, 0x11, 0x1f, 0x8f, 0x1f, 0xfb, 0x3c, 0x11, 0x98, 0x33, 0x9f, 0x3a,
	0xdc, 0xff, 0xb0, 0x00, 0xcb, 0xb3, 0x15, 0xc8, 0x01, 0x34, 0x4c, 0xa7, 0xb2, 0xe3, 0xf5, 0x16,
	0xfa, 0xcd, 0xad, 0xe0, 0x36, 0x4d, 0xc1, 0x91, 0x19, 0x11, 0x1b, 0x1d, 0x68, 0x8c, 0x3a, 0x9a,
	0x1c, 0xc0, 0xb2, 0x54, 0x98, 0x9d, 0xc4, 0x28, 0x46, 0x5c, 0x0f, 0xd0, 0xf5, 0xd8, 0xbd, 0x99,
	0x6f, 0xa8, 0x30, 0xdb, 0x2b, 0xc3, 0x68, 0x4b, 0xfe, 0x6c, 0x92, 0xa7, 0x50, 0xd7, 0x62, 0x59,
	0xa7, 0xd6, 0xf3, 0xfa, 0xcb, 0x5b, 0xf7, 0x6e, 0x95, 0xa3, 0xbb, 0x67, 0xd4, 0x42, 0xeb, 0xef,
	0x3c, 0x68, 0xcd, 0xe8, 0x23, 0xdb, 0xb0, 0x28, 0x59, 0x7e, 0xce, 0x63, 0xe6, 0xae, 0x6a, 0xcd,
	0x64, 0xac, 0xb2, 0xed, 0x33, 0xa9, 0xb8, 0x30, 0x6f, 0x9b, 0x96, 0x91, 0x64, 0x17, 0xea, 0xa6,
	0x2d, 0xd7, 0xc3, 0xfd, 0xc0, 0x3d, 0x06, 0xfb, 0x52, 0x66, 0x13, 0x50, 0x2c, 0x14, 0x33, 0xb5,
	0x24, 0xb5, 0x94, 0xdf, 0x87, 0xba, 0x51, 0x45, 0x96, 0xa0, 0x76, 0x24, 0xb8, 0x6a, 0xcf, 0x91,
	0x26, 0x2c, 0xd2, 0x42, 0x08, 0x2e, 0x92, 0xb6, 0xa7, 0x8f, 0xf7, 0x51, 0xb0, 0xf6, 0xbc, 0xff,
	0xde, 0x83, 0xd6, 0xcc, 0x38, 0xc8, 0x13, 0x58, 0x2a, 0x77, 0xed, 0x4a, 0xb0, 0x59, 0xc6, 0xa0,
	0x5c, 0xc6, 0x60, 0xdf, 0x05, 0x0c, 0x6a, 0x1f, 0xbf, 0x76, 0x3d, 0x5a, 0x01, 0xe4, 0x11, 0x34,
	0xa6, 0x4c, 0xe5, 0x5c, 0x3f, 0x30, 0x7d, 0x99, 0x77, 0x6f, 0x4e, 0xef, 0xd8, 0xf8, 0x5f, 0x4d,
	0x72, 0x26, 0x27, 0x98, 0x8e, 0xa8, 0x03, 0xfc, 0x43, 0x58, 0xb9, 0xe6, 0x22, 0x5d, 0x68, 0x5a,
	0xe7, 0x89, 0x88, 0xa6, 0x76, 0x7c, 0xff, 0x51, 0xb0, 0x47, 0x2f, 0xa2, 0x29, 0x23, 0xab, 0x50,
	0x3f, 0x8f, 0xd2, 0x82, 0x99, 0x31, 0x79, 0xd4, 0x1a, 0xfe, 0x0a, 0xb4, 0xf6, 0x26, 0x11, 0xca,
	0x21, 0x53, 0x8a, 0x8b, 0x44, 0x0e, 0xf6, 0x3e, 0x7d, 0xdb, 0xf0, 0x5e, 0xef, 0xfe, 0xc3, 0x9f,
	0x58, 0x76, 0x96, 0xb8, 0x2d, 0x3d, 0x6d, 0x98, 0xee, 0xb7, 0x7f, 0x04, 0x00, 0x00, 0xff, 0xff,
	0xd4, 0x1c, 0x58, 0xf7, 0x0c, 0x05, 0x00, 0x00,
}
