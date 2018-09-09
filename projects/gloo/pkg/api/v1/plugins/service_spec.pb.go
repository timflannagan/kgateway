// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: service_spec.proto

/*
Package plugins is a generated protocol buffer package.

It is generated from these files:
	service_spec.proto

It has these top-level messages:
	ServiceSpec
*/
package plugins

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import rest_plugins_gloo_solo_io "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/rest"
import grpc_plugins_gloo_solo_io "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/grpc"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// Describes APIs and application-level information for services
// Gloo routes to. ServiceSpec is contained within the UpstreamSpec for certain types
// of upstreams, including Kubernetes, Consul, and Static.
// ServiceSpec configuration is opaque to Gloo and handled by Service Plugins.
type ServiceSpec struct {
	// Note to developers: new Service Plugins must be added to this oneof field
	// to be usable by Gloo.
	//
	// Types that are valid to be assigned to PluginType:
	//	*ServiceSpec_Rest
	//	*ServiceSpec_Grpc
	PluginType isServiceSpec_PluginType `protobuf_oneof:"plugin_type"`
}

func (m *ServiceSpec) Reset()                    { *m = ServiceSpec{} }
func (m *ServiceSpec) String() string            { return proto.CompactTextString(m) }
func (*ServiceSpec) ProtoMessage()               {}
func (*ServiceSpec) Descriptor() ([]byte, []int) { return fileDescriptorServiceSpec, []int{0} }

type isServiceSpec_PluginType interface {
	isServiceSpec_PluginType()
	Equal(interface{}) bool
}

type ServiceSpec_Rest struct {
	Rest *rest_plugins_gloo_solo_io.ServiceSpec `protobuf:"bytes,1,opt,name=rest,oneof"`
}
type ServiceSpec_Grpc struct {
	Grpc *grpc_plugins_gloo_solo_io.ServiceSpec `protobuf:"bytes,2,opt,name=grpc,oneof"`
}

func (*ServiceSpec_Rest) isServiceSpec_PluginType() {}
func (*ServiceSpec_Grpc) isServiceSpec_PluginType() {}

func (m *ServiceSpec) GetPluginType() isServiceSpec_PluginType {
	if m != nil {
		return m.PluginType
	}
	return nil
}

func (m *ServiceSpec) GetRest() *rest_plugins_gloo_solo_io.ServiceSpec {
	if x, ok := m.GetPluginType().(*ServiceSpec_Rest); ok {
		return x.Rest
	}
	return nil
}

func (m *ServiceSpec) GetGrpc() *grpc_plugins_gloo_solo_io.ServiceSpec {
	if x, ok := m.GetPluginType().(*ServiceSpec_Grpc); ok {
		return x.Grpc
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*ServiceSpec) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _ServiceSpec_OneofMarshaler, _ServiceSpec_OneofUnmarshaler, _ServiceSpec_OneofSizer, []interface{}{
		(*ServiceSpec_Rest)(nil),
		(*ServiceSpec_Grpc)(nil),
	}
}

func _ServiceSpec_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*ServiceSpec)
	// plugin_type
	switch x := m.PluginType.(type) {
	case *ServiceSpec_Rest:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Rest); err != nil {
			return err
		}
	case *ServiceSpec_Grpc:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Grpc); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("ServiceSpec.PluginType has unexpected type %T", x)
	}
	return nil
}

func _ServiceSpec_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*ServiceSpec)
	switch tag {
	case 1: // plugin_type.rest
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(rest_plugins_gloo_solo_io.ServiceSpec)
		err := b.DecodeMessage(msg)
		m.PluginType = &ServiceSpec_Rest{msg}
		return true, err
	case 2: // plugin_type.grpc
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(grpc_plugins_gloo_solo_io.ServiceSpec)
		err := b.DecodeMessage(msg)
		m.PluginType = &ServiceSpec_Grpc{msg}
		return true, err
	default:
		return false, nil
	}
}

func _ServiceSpec_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*ServiceSpec)
	// plugin_type
	switch x := m.PluginType.(type) {
	case *ServiceSpec_Rest:
		s := proto.Size(x.Rest)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *ServiceSpec_Grpc:
		s := proto.Size(x.Grpc)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

func init() {
	proto.RegisterType((*ServiceSpec)(nil), "plugins.gloo.solo.io.ServiceSpec")
}
func (this *ServiceSpec) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ServiceSpec)
	if !ok {
		that2, ok := that.(ServiceSpec)
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
	if that1.PluginType == nil {
		if this.PluginType != nil {
			return false
		}
	} else if this.PluginType == nil {
		return false
	} else if !this.PluginType.Equal(that1.PluginType) {
		return false
	}
	return true
}
func (this *ServiceSpec_Rest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ServiceSpec_Rest)
	if !ok {
		that2, ok := that.(ServiceSpec_Rest)
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
	if !this.Rest.Equal(that1.Rest) {
		return false
	}
	return true
}
func (this *ServiceSpec_Grpc) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ServiceSpec_Grpc)
	if !ok {
		that2, ok := that.(ServiceSpec_Grpc)
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
	if !this.Grpc.Equal(that1.Grpc) {
		return false
	}
	return true
}

func init() { proto.RegisterFile("service_spec.proto", fileDescriptorServiceSpec) }

var fileDescriptorServiceSpec = []byte{
	// 228 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2a, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0x8d, 0x2f, 0x2e, 0x48, 0x4d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x12,
	0x29, 0xc8, 0x29, 0x4d, 0xcf, 0xcc, 0x2b, 0xd6, 0x4b, 0xcf, 0xc9, 0xcf, 0xd7, 0x2b, 0xce, 0xcf,
	0xc9, 0xd7, 0xcb, 0xcc, 0x97, 0x12, 0x49, 0xcf, 0x4f, 0xcf, 0x07, 0x2b, 0xd0, 0x07, 0xb1, 0x20,
	0x6a, 0xa5, 0x3c, 0xd2, 0x33, 0x4b, 0x32, 0x4a, 0x93, 0xf4, 0x92, 0xf3, 0x73, 0xf5, 0x41, 0x2a,
	0x75, 0x33, 0xf3, 0x21, 0x74, 0x76, 0x66, 0x89, 0x7e, 0x41, 0x51, 0x7e, 0x56, 0x6a, 0x72, 0x49,
	0xb1, 0x3e, 0xc8, 0x20, 0xfd, 0xc4, 0x82, 0x4c, 0xfd, 0x32, 0x43, 0x7d, 0xa8, 0xe1, 0xfa, 0x45,
	0xa9, 0xc5, 0x25, 0x60, 0x82, 0x62, 0x93, 0xd2, 0x8b, 0x0a, 0x92, 0xc1, 0x04, 0xc4, 0x24, 0xa5,
	0x19, 0x8c, 0x5c, 0xdc, 0xc1, 0x10, 0x6f, 0x05, 0x17, 0xa4, 0x26, 0x0b, 0xd9, 0x70, 0xb1, 0x80,
	0xec, 0x91, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x36, 0x52, 0xd3, 0x83, 0x58, 0x8a, 0xc5, 0x8f, 0x7a,
	0x48, 0xba, 0x3c, 0x18, 0x82, 0xc0, 0xba, 0x40, 0xba, 0x41, 0x66, 0x4b, 0x30, 0x41, 0x75, 0x43,
	0x2c, 0x22, 0x42, 0x37, 0x48, 0xa1, 0x13, 0x2f, 0x17, 0x37, 0x44, 0x6d, 0x7c, 0x49, 0x65, 0x41,
	0xaa, 0x93, 0xd3, 0x8a, 0x47, 0x72, 0x8c, 0x51, 0x36, 0xc4, 0x7b, 0xb5, 0x20, 0x3b, 0x1d, 0xcd,
	0xbb, 0x49, 0x6c, 0x60, 0x5f, 0x1a, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x95, 0x73, 0x75, 0x56,
	0xbb, 0x01, 0x00, 0x00,
}
