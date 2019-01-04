// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/solo-io/solo-projects/projects/vcs/api/v1/changeset.proto

package v1 // import "github.com/solo-io/solo-projects/projects/vcs/pkg/api/v1"

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import types "github.com/gogo/protobuf/types"
import v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
import v11 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
import core "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
import v12 "github.com/solo-io/solo-projects/projects/sqoop/pkg/api/v1"

import bytes "bytes"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Action int32

const (
	Action_NONE      Action = 0
	Action_COMMIT    Action = 1
	Action_CHECK_OUT Action = 2
)

var Action_name = map[int32]string{
	0: "NONE",
	1: "COMMIT",
	2: "CHECK_OUT",
}
var Action_value = map[string]int32{
	"NONE":      0,
	"COMMIT":    1,
	"CHECK_OUT": 2,
}

func (x Action) String() string {
	return proto.EnumName(Action_name, int32(x))
}
func (Action) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_changeset_e6ab33656361388b, []int{0}
}

//
// @solo-kit:resource.short_name=chg
// @solo-kit:resource.plural_name=changesets
// @solo-kit:resource.resource_groups=api.vcs.solo.io
//
// The ChangeSet object represents the current status of a Gloo user's working directory. Each element in the "data"
// element represents the complete snapshot of a resource.
//
type ChangeSet struct {
	// Status indicates the validation status of this resource
	Status core.Status `protobuf:"bytes,1,opt,name=status" json:"status" testdiff:"ignore"`
	// Metadata for this resource
	Metadata core.Metadata `protobuf:"bytes,2,opt,name=metadata" json:"metadata"`
	// The name of the git branch the changes will be applied to
	Branch types.StringValue `protobuf:"bytes,3,opt,name=branch" json:"branch"`
	// Represents an action that this changeset is waiting to be performed. The default value is NONE.
	PendingAction Action `protobuf:"varint,4,opt,name=pending_action,json=pendingAction,proto3,enum=vcs.solo.io.Action" json:"pending_action,omitempty"`
	// Description of the changeset. This will be the git commit message
	Description types.StringValue `protobuf:"bytes,5,opt,name=description" json:"description"`
	// The number of edits that the user applied to the previous commit.
	// A value greater than zero represents a dirty work tree.
	EditCount types.UInt32Value `protobuf:"bytes,6,opt,name=edit_count,json=editCount" json:"edit_count"`
	// The user who owns this changeset
	// TODO use dedicated message? Also, determine how to handle secrets?
	UserId types.StringValue `protobuf:"bytes,7,opt,name=user_id,json=userId" json:"user_id"`
	// The hash of the commit that the changeset represents an increment upon
	RootCommit types.StringValue `protobuf:"bytes,8,opt,name=root_commit,json=rootCommit" json:"root_commit"`
	// The git commit message for the root commit
	RootDescription types.StringValue `protobuf:"bytes,9,opt,name=root_description,json=rootDescription" json:"root_description"`
	// If a git commit attempt fails, this field will be populated with a user-friendly error message
	// No further git commit attempts will be possible until the user clears this field
	ErrorMsg *types.StringValue `protobuf:"bytes,10,opt,name=error_msg,json=errorMsg" json:"error_msg,omitempty"`
	// A collection of Gloo resources
	Data                 Data     `protobuf:"bytes,11,opt,name=data" json:"data"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChangeSet) Reset()         { *m = ChangeSet{} }
func (m *ChangeSet) String() string { return proto.CompactTextString(m) }
func (*ChangeSet) ProtoMessage()    {}
func (*ChangeSet) Descriptor() ([]byte, []int) {
	return fileDescriptor_changeset_e6ab33656361388b, []int{0}
}
func (m *ChangeSet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChangeSet.Unmarshal(m, b)
}
func (m *ChangeSet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChangeSet.Marshal(b, m, deterministic)
}
func (dst *ChangeSet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangeSet.Merge(dst, src)
}
func (m *ChangeSet) XXX_Size() int {
	return xxx_messageInfo_ChangeSet.Size(m)
}
func (m *ChangeSet) XXX_DiscardUnknown() {
	xxx_messageInfo_ChangeSet.DiscardUnknown(m)
}

var xxx_messageInfo_ChangeSet proto.InternalMessageInfo

func (m *ChangeSet) GetStatus() core.Status {
	if m != nil {
		return m.Status
	}
	return core.Status{}
}

func (m *ChangeSet) GetMetadata() core.Metadata {
	if m != nil {
		return m.Metadata
	}
	return core.Metadata{}
}

func (m *ChangeSet) GetBranch() types.StringValue {
	if m != nil {
		return m.Branch
	}
	return types.StringValue{}
}

func (m *ChangeSet) GetPendingAction() Action {
	if m != nil {
		return m.PendingAction
	}
	return Action_NONE
}

func (m *ChangeSet) GetDescription() types.StringValue {
	if m != nil {
		return m.Description
	}
	return types.StringValue{}
}

func (m *ChangeSet) GetEditCount() types.UInt32Value {
	if m != nil {
		return m.EditCount
	}
	return types.UInt32Value{}
}

func (m *ChangeSet) GetUserId() types.StringValue {
	if m != nil {
		return m.UserId
	}
	return types.StringValue{}
}

func (m *ChangeSet) GetRootCommit() types.StringValue {
	if m != nil {
		return m.RootCommit
	}
	return types.StringValue{}
}

func (m *ChangeSet) GetRootDescription() types.StringValue {
	if m != nil {
		return m.RootDescription
	}
	return types.StringValue{}
}

func (m *ChangeSet) GetErrorMsg() *types.StringValue {
	if m != nil {
		return m.ErrorMsg
	}
	return nil
}

func (m *ChangeSet) GetData() Data {
	if m != nil {
		return m.Data
	}
	return Data{}
}

// A user-specific snapshot of all gloo resources at a given commit plus any non-committed changes made by the user
type Data struct {
	Gateways             []*v1.Gateway        `protobuf:"bytes,1,rep,name=gateways" json:"gateways,omitempty"`
	VirtualServices      []*v1.VirtualService `protobuf:"bytes,2,rep,name=virtual_services,json=virtualServices" json:"virtual_services,omitempty"`
	Proxies              []*v11.Proxy         `protobuf:"bytes,3,rep,name=proxies" json:"proxies,omitempty"`
	Settings             []*v11.Settings      `protobuf:"bytes,4,rep,name=settings" json:"settings,omitempty"`
	Upstreams            []*v11.Upstream      `protobuf:"bytes,5,rep,name=upstreams" json:"upstreams,omitempty"`
	ResolverMaps         []*v12.ResolverMap   `protobuf:"bytes,6,rep,name=resolver_maps,json=resolverMaps" json:"resolver_maps,omitempty"`
	Schemas              []*v12.Schema        `protobuf:"bytes,7,rep,name=schemas" json:"schemas,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Data) Reset()         { *m = Data{} }
func (m *Data) String() string { return proto.CompactTextString(m) }
func (*Data) ProtoMessage()    {}
func (*Data) Descriptor() ([]byte, []int) {
	return fileDescriptor_changeset_e6ab33656361388b, []int{1}
}
func (m *Data) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Data.Unmarshal(m, b)
}
func (m *Data) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Data.Marshal(b, m, deterministic)
}
func (dst *Data) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Data.Merge(dst, src)
}
func (m *Data) XXX_Size() int {
	return xxx_messageInfo_Data.Size(m)
}
func (m *Data) XXX_DiscardUnknown() {
	xxx_messageInfo_Data.DiscardUnknown(m)
}

var xxx_messageInfo_Data proto.InternalMessageInfo

func (m *Data) GetGateways() []*v1.Gateway {
	if m != nil {
		return m.Gateways
	}
	return nil
}

func (m *Data) GetVirtualServices() []*v1.VirtualService {
	if m != nil {
		return m.VirtualServices
	}
	return nil
}

func (m *Data) GetProxies() []*v11.Proxy {
	if m != nil {
		return m.Proxies
	}
	return nil
}

func (m *Data) GetSettings() []*v11.Settings {
	if m != nil {
		return m.Settings
	}
	return nil
}

func (m *Data) GetUpstreams() []*v11.Upstream {
	if m != nil {
		return m.Upstreams
	}
	return nil
}

func (m *Data) GetResolverMaps() []*v12.ResolverMap {
	if m != nil {
		return m.ResolverMaps
	}
	return nil
}

func (m *Data) GetSchemas() []*v12.Schema {
	if m != nil {
		return m.Schemas
	}
	return nil
}

func init() {
	proto.RegisterType((*ChangeSet)(nil), "vcs.solo.io.ChangeSet")
	proto.RegisterType((*Data)(nil), "vcs.solo.io.Data")
	proto.RegisterEnum("vcs.solo.io.Action", Action_name, Action_value)
}
func (this *ChangeSet) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ChangeSet)
	if !ok {
		that2, ok := that.(ChangeSet)
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
	if !this.Status.Equal(&that1.Status) {
		return false
	}
	if !this.Metadata.Equal(&that1.Metadata) {
		return false
	}
	if !this.Branch.Equal(&that1.Branch) {
		return false
	}
	if this.PendingAction != that1.PendingAction {
		return false
	}
	if !this.Description.Equal(&that1.Description) {
		return false
	}
	if !this.EditCount.Equal(&that1.EditCount) {
		return false
	}
	if !this.UserId.Equal(&that1.UserId) {
		return false
	}
	if !this.RootCommit.Equal(&that1.RootCommit) {
		return false
	}
	if !this.RootDescription.Equal(&that1.RootDescription) {
		return false
	}
	if !this.ErrorMsg.Equal(that1.ErrorMsg) {
		return false
	}
	if !this.Data.Equal(&that1.Data) {
		return false
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
func (this *Data) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Data)
	if !ok {
		that2, ok := that.(Data)
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
	if len(this.Gateways) != len(that1.Gateways) {
		return false
	}
	for i := range this.Gateways {
		if !this.Gateways[i].Equal(that1.Gateways[i]) {
			return false
		}
	}
	if len(this.VirtualServices) != len(that1.VirtualServices) {
		return false
	}
	for i := range this.VirtualServices {
		if !this.VirtualServices[i].Equal(that1.VirtualServices[i]) {
			return false
		}
	}
	if len(this.Proxies) != len(that1.Proxies) {
		return false
	}
	for i := range this.Proxies {
		if !this.Proxies[i].Equal(that1.Proxies[i]) {
			return false
		}
	}
	if len(this.Settings) != len(that1.Settings) {
		return false
	}
	for i := range this.Settings {
		if !this.Settings[i].Equal(that1.Settings[i]) {
			return false
		}
	}
	if len(this.Upstreams) != len(that1.Upstreams) {
		return false
	}
	for i := range this.Upstreams {
		if !this.Upstreams[i].Equal(that1.Upstreams[i]) {
			return false
		}
	}
	if len(this.ResolverMaps) != len(that1.ResolverMaps) {
		return false
	}
	for i := range this.ResolverMaps {
		if !this.ResolverMaps[i].Equal(that1.ResolverMaps[i]) {
			return false
		}
	}
	if len(this.Schemas) != len(that1.Schemas) {
		return false
	}
	for i := range this.Schemas {
		if !this.Schemas[i].Equal(that1.Schemas[i]) {
			return false
		}
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}

func init() {
	proto.RegisterFile("github.com/solo-io/solo-projects/projects/vcs/api/v1/changeset.proto", fileDescriptor_changeset_e6ab33656361388b)
}

var fileDescriptor_changeset_e6ab33656361388b = []byte{
	// 759 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x95, 0xdf, 0x52, 0xdb, 0x38,
	0x14, 0xc6, 0x09, 0x09, 0xf9, 0xa3, 0x2c, 0x10, 0x04, 0xcb, 0x78, 0x99, 0x1d, 0x60, 0x72, 0xc5,
	0xec, 0x0e, 0xf6, 0x10, 0xb8, 0x60, 0x61, 0x76, 0x77, 0xc0, 0x61, 0x68, 0xda, 0x09, 0x74, 0x1c,
	0xe0, 0xa2, 0x37, 0x19, 0xc5, 0x11, 0x8a, 0x4a, 0x6c, 0xb9, 0x92, 0x1c, 0xda, 0xe7, 0xe8, 0x4b,
	0xf4, 0x51, 0x78, 0x0a, 0x2e, 0xfa, 0x08, 0x7d, 0x82, 0x8e, 0x64, 0x39, 0x7f, 0x28, 0x0c, 0x81,
	0xab, 0xf8, 0xf8, 0x7c, 0xbf, 0x4f, 0x8e, 0x74, 0xfc, 0x19, 0xd4, 0x09, 0x95, 0xbd, 0xb8, 0x63,
	0xfb, 0x2c, 0x70, 0x04, 0xeb, 0xb3, 0x6d, 0xca, 0x92, 0xdf, 0x88, 0xb3, 0x8f, 0xd8, 0x97, 0xc2,
	0x19, 0x5e, 0x0c, 0x7c, 0xe1, 0xa0, 0x88, 0x3a, 0x83, 0x1d, 0xc7, 0xef, 0xa1, 0x90, 0x60, 0x81,
	0xa5, 0x1d, 0x71, 0x26, 0x19, 0x2c, 0x0f, 0x7c, 0x61, 0x2b, 0xcc, 0xa6, 0x6c, 0x6d, 0x85, 0x30,
	0xc2, 0xf4, 0x7d, 0x47, 0x5d, 0x25, 0x92, 0xb5, 0x75, 0xc2, 0x18, 0xe9, 0x63, 0x47, 0x57, 0x9d,
	0xf8, 0xda, 0xb9, 0xe5, 0x28, 0x8a, 0x30, 0x17, 0xa6, 0xbf, 0xf3, 0xd4, 0x83, 0xdc, 0x50, 0x99,
	0x2e, 0x1b, 0x60, 0x89, 0xba, 0x48, 0x22, 0x83, 0x38, 0x53, 0x20, 0x42, 0x22, 0x19, 0xa7, 0x6b,
	0xfc, 0xfb, 0x08, 0x40, 0xfa, 0x8c, 0x8d, 0xfe, 0x23, 0x41, 0x12, 0xdf, 0xa2, 0x2f, 0x29, 0x6d,
	0x4a, 0x83, 0x9f, 0xbc, 0x14, 0x1f, 0x50, 0x2e, 0x63, 0xd4, 0x6f, 0x0b, 0xcc, 0x07, 0xd4, 0xc7,
	0xc6, 0x66, 0xff, 0x79, 0x1b, 0x55, 0x19, 0x8f, 0x88, 0xb3, 0xcf, 0xe9, 0x03, 0x1c, 0xbe, 0x88,
	0x14, 0x58, 0x4a, 0x1a, 0x12, 0xf1, 0x2a, 0x38, 0x8e, 0x84, 0xe4, 0x18, 0x05, 0x06, 0x6e, 0x4c,
	0x3f, 0x26, 0xe2, 0x13, 0x63, 0x51, 0x6a, 0xc3, 0xb1, 0x60, 0xfd, 0x01, 0xe6, 0xed, 0x00, 0x45,
	0xc6, 0xca, 0x7d, 0xa5, 0x95, 0xf0, 0x7b, 0x38, 0x30, 0x47, 0x5f, 0xbd, 0x9b, 0x03, 0x25, 0x57,
	0x0f, 0x61, 0x0b, 0x4b, 0x78, 0x0a, 0xf2, 0xc9, 0x39, 0x5b, 0x99, 0xcd, 0xcc, 0x56, 0xb9, 0xb6,
	0x62, 0xfb, 0x8c, 0xe3, 0x74, 0x20, 0xed, 0x96, 0xee, 0x1d, 0xff, 0x71, 0x77, 0xbf, 0x31, 0xf3,
	0xe3, 0x7e, 0x63, 0x49, 0x62, 0x21, 0xbb, 0xf4, 0xfa, 0xfa, 0xa0, 0x4a, 0x49, 0xc8, 0x38, 0xae,
	0x7a, 0x06, 0x87, 0xfb, 0xa0, 0x98, 0xce, 0x98, 0x35, 0xab, 0xad, 0x56, 0x27, 0xad, 0x9a, 0xa6,
	0x7b, 0x9c, 0x53, 0x66, 0xde, 0x50, 0x0d, 0x0f, 0x40, 0xbe, 0xc3, 0x51, 0xe8, 0xf7, 0xac, 0xac,
	0xe6, 0xfe, 0xb4, 0x93, 0x79, 0xb7, 0xd3, 0x79, 0xb7, 0x5b, 0x92, 0xd3, 0x90, 0x5c, 0xa1, 0x7e,
	0x8c, 0x0d, 0x6d, 0x08, 0x78, 0x00, 0x16, 0x22, 0x1c, 0x76, 0x69, 0x48, 0xda, 0xc8, 0x97, 0x94,
	0x85, 0x56, 0x6e, 0x33, 0xb3, 0xb5, 0x50, 0x5b, 0xb6, 0xc7, 0x5e, 0x2b, 0xfb, 0x48, 0xb7, 0xbc,
	0x79, 0x23, 0x4d, 0x4a, 0x58, 0x07, 0xe5, 0x2e, 0x16, 0x3e, 0xa7, 0x91, 0x06, 0xe7, 0xa6, 0x5e,
	0x7c, 0x1c, 0x83, 0x47, 0x00, 0xe0, 0x2e, 0x95, 0x6d, 0x9f, 0xc5, 0xa1, 0xb4, 0xf2, 0x4f, 0x98,
	0x5c, 0x36, 0x42, 0xb9, 0x5b, 0x1b, 0x37, 0x29, 0x29, 0xca, 0x55, 0x10, 0x3c, 0x04, 0x85, 0x58,
	0x60, 0xde, 0xa6, 0x5d, 0xab, 0x30, 0xfd, 0x0e, 0x28, 0xa4, 0xd1, 0x85, 0x2e, 0x28, 0x73, 0xc6,
	0xd4, 0xfa, 0x41, 0x40, 0xa5, 0x55, 0x9c, 0xda, 0x00, 0x28, 0xcc, 0xd5, 0x14, 0x6c, 0x82, 0x8a,
	0x36, 0x19, 0xdf, 0x8f, 0xd2, 0xd4, 0x4e, 0x8b, 0x8a, 0xad, 0x8f, 0xed, 0xc9, 0x3f, 0xa0, 0x84,
	0x39, 0x67, 0xbc, 0x1d, 0x08, 0x62, 0x81, 0xe7, 0x7d, 0xbc, 0xa2, 0x96, 0x37, 0x05, 0x81, 0x7f,
	0x83, 0x9c, 0x1e, 0xa1, 0xb2, 0xa6, 0x96, 0x26, 0x8e, 0xb1, 0x3e, 0x9a, 0x1e, 0x2d, 0xaa, 0x7e,
	0xcd, 0x82, 0x9c, 0xba, 0x09, 0xf7, 0x40, 0xd1, 0xe4, 0x87, 0x9a, 0xe3, 0xec, 0x56, 0xb9, 0x66,
	0xd9, 0x69, 0x00, 0xa5, 0xf4, 0x69, 0x52, 0x7b, 0x43, 0x25, 0x7c, 0x0b, 0x2a, 0x0f, 0x62, 0x46,
	0x58, 0xb3, 0x9a, 0xde, 0xf8, 0x85, 0xbe, 0x4a, 0x84, 0xad, 0x44, 0xe7, 0x2d, 0x0e, 0x26, 0x6a,
	0x01, 0xb7, 0x41, 0x41, 0xc5, 0x0d, 0xc5, 0xc2, 0xca, 0x6a, 0x8b, 0x65, 0x5b, 0x65, 0xc2, 0x90,
	0x7f, 0xaf, 0xb2, 0xc8, 0x4b, 0x35, 0xb0, 0x06, 0x8a, 0x69, 0xc6, 0x58, 0x39, 0xad, 0x5f, 0x9d,
	0xd4, 0xb7, 0x4c, 0xd7, 0x1b, 0xea, 0xe0, 0x1e, 0x28, 0xa5, 0xd1, 0x22, 0xac, 0xb9, 0xc7, 0xa0,
	0x4b, 0xd3, 0xf6, 0x46, 0x42, 0xf8, 0x3f, 0x98, 0x1f, 0x4f, 0x12, 0x61, 0xe5, 0x35, 0xb9, 0x66,
	0xeb, 0x84, 0x18, 0xa2, 0x9e, 0xd1, 0x34, 0x51, 0xe4, 0xfd, 0xc6, 0x47, 0x85, 0x80, 0x0e, 0x28,
	0x24, 0xf9, 0x21, 0xac, 0x82, 0x46, 0x7f, 0x7f, 0x80, 0xb6, 0x74, 0xd7, 0x4b, 0x55, 0x7f, 0x6d,
	0x83, 0xbc, 0x79, 0xc3, 0x8a, 0x20, 0x77, 0x76, 0x7e, 0x76, 0x52, 0x99, 0x81, 0x00, 0xe4, 0xdd,
	0xf3, 0x66, 0xb3, 0x71, 0x51, 0xc9, 0xc0, 0x79, 0x50, 0x72, 0xdf, 0x9c, 0xb8, 0xef, 0xda, 0xe7,
	0x97, 0x17, 0x95, 0xd9, 0xe3, 0xff, 0xbe, 0x7d, 0x5f, 0xcf, 0x7c, 0xd8, 0x7f, 0xd9, 0xc7, 0x34,
	0xba, 0x21, 0x26, 0xdc, 0x3a, 0x79, 0x3d, 0x51, 0xbb, 0x3f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xe0,
	0xef, 0x00, 0x92, 0x8f, 0x07, 0x00, 0x00,
}
