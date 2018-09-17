// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: bootstrap.proto

package v1

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"
import google_protobuf1 "github.com/gogo/protobuf/types"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

//
// @solo-kit:resource.short_name=st
// @solo-kit:resource.plural_name=settings
// @solo-kit:resource.resource_groups=setup.gloo.solo.io
type Settings struct {
	// namespace to write discovered data
	DiscoveryNamespace string `protobuf:"bytes,1,opt,name=discovery_namespace,json=discoveryNamespace,proto3" json:"discovery_namespace,omitempty"`
	// namespaces to watch for user config as well as
	WatchNamespaces string `protobuf:"bytes,2,opt,name=watch_namespaces,json=watchNamespaces,proto3" json:"watch_namespaces,omitempty"`
	// where to read user config (upstream, proxy) from
	// if nil, use only in memory config
	//
	// Types that are valid to be assigned to ConfigSource:
	//	*Settings_KubernetesConfigSource
	//	*Settings_DirectoryConfigSource
	ConfigSource isSettings_ConfigSource `protobuf_oneof:"config_source"`
	// where to read secrets from (vault, k8s)
	//
	// Types that are valid to be assigned to SecretSource:
	//	*Settings_KubernetesSecretSource
	//	*Settings_VaultSecretSource
	//	*Settings_DirectorySecretSource
	SecretSource isSettings_SecretSource `protobuf_oneof:"secret_source"`
	// where to read artifacts from (configmap, file)
	//
	// Types that are valid to be assigned to ArtifactSource:
	//	*Settings_KubernetesArtifactSource
	//	*Settings_DirectoryArtifactSource
	ArtifactSource isSettings_ArtifactSource `protobuf_oneof:"artifact_source"`
	// where the gloo xds server should bind (should not need configuration by user)
	BindAddr string `protobuf:"bytes,11,opt,name=bind_addr,json=bindAddr,proto3" json:"bind_addr,omitempty"`
	// how frequently to resync watches, etc
	RefreshRate *google_protobuf1.Duration `protobuf:"bytes,12,opt,name=refresh_rate,json=refreshRate" json:"refresh_rate,omitempty"`
	// enable serving debug data on port 9090
	DevMode bool `protobuf:"varint,13,opt,name=DevMode,proto3" json:"DevMode,omitempty"`
}

func (m *Settings) Reset()                    { *m = Settings{} }
func (m *Settings) String() string            { return proto.CompactTextString(m) }
func (*Settings) ProtoMessage()               {}
func (*Settings) Descriptor() ([]byte, []int) { return fileDescriptorBootstrap, []int{0} }

type isSettings_ConfigSource interface {
	isSettings_ConfigSource()
	Equal(interface{}) bool
}
type isSettings_SecretSource interface {
	isSettings_SecretSource()
	Equal(interface{}) bool
}
type isSettings_ArtifactSource interface {
	isSettings_ArtifactSource()
	Equal(interface{}) bool
}

type Settings_KubernetesConfigSource struct {
	KubernetesConfigSource *Settings_KubernetesCrds `protobuf:"bytes,4,opt,name=kubernetes_config_source,json=kubernetesConfigSource,oneof"`
}
type Settings_DirectoryConfigSource struct {
	DirectoryConfigSource *Settings_Directory `protobuf:"bytes,5,opt,name=directory_config_source,json=directoryConfigSource,oneof"`
}
type Settings_KubernetesSecretSource struct {
	KubernetesSecretSource *Settings_KubernetesSecrets `protobuf:"bytes,6,opt,name=kubernetes_secret_source,json=kubernetesSecretSource,oneof"`
}
type Settings_VaultSecretSource struct {
	VaultSecretSource *Settings_VaultSecrets `protobuf:"bytes,7,opt,name=vault_secret_source,json=vaultSecretSource,oneof"`
}
type Settings_DirectorySecretSource struct {
	DirectorySecretSource *Settings_Directory `protobuf:"bytes,8,opt,name=directory_secret_source,json=directorySecretSource,oneof"`
}
type Settings_KubernetesArtifactSource struct {
	KubernetesArtifactSource *Settings_KubernetesConfigmaps `protobuf:"bytes,9,opt,name=kubernetes_artifact_source,json=kubernetesArtifactSource,oneof"`
}
type Settings_DirectoryArtifactSource struct {
	DirectoryArtifactSource *Settings_Directory `protobuf:"bytes,10,opt,name=directory_artifact_source,json=directoryArtifactSource,oneof"`
}

func (*Settings_KubernetesConfigSource) isSettings_ConfigSource()     {}
func (*Settings_DirectoryConfigSource) isSettings_ConfigSource()      {}
func (*Settings_KubernetesSecretSource) isSettings_SecretSource()     {}
func (*Settings_VaultSecretSource) isSettings_SecretSource()          {}
func (*Settings_DirectorySecretSource) isSettings_SecretSource()      {}
func (*Settings_KubernetesArtifactSource) isSettings_ArtifactSource() {}
func (*Settings_DirectoryArtifactSource) isSettings_ArtifactSource()  {}

func (m *Settings) GetConfigSource() isSettings_ConfigSource {
	if m != nil {
		return m.ConfigSource
	}
	return nil
}
func (m *Settings) GetSecretSource() isSettings_SecretSource {
	if m != nil {
		return m.SecretSource
	}
	return nil
}
func (m *Settings) GetArtifactSource() isSettings_ArtifactSource {
	if m != nil {
		return m.ArtifactSource
	}
	return nil
}

func (m *Settings) GetDiscoveryNamespace() string {
	if m != nil {
		return m.DiscoveryNamespace
	}
	return ""
}

func (m *Settings) GetWatchNamespaces() string {
	if m != nil {
		return m.WatchNamespaces
	}
	return ""
}

func (m *Settings) GetKubernetesConfigSource() *Settings_KubernetesCrds {
	if x, ok := m.GetConfigSource().(*Settings_KubernetesConfigSource); ok {
		return x.KubernetesConfigSource
	}
	return nil
}

func (m *Settings) GetDirectoryConfigSource() *Settings_Directory {
	if x, ok := m.GetConfigSource().(*Settings_DirectoryConfigSource); ok {
		return x.DirectoryConfigSource
	}
	return nil
}

func (m *Settings) GetKubernetesSecretSource() *Settings_KubernetesSecrets {
	if x, ok := m.GetSecretSource().(*Settings_KubernetesSecretSource); ok {
		return x.KubernetesSecretSource
	}
	return nil
}

func (m *Settings) GetVaultSecretSource() *Settings_VaultSecrets {
	if x, ok := m.GetSecretSource().(*Settings_VaultSecretSource); ok {
		return x.VaultSecretSource
	}
	return nil
}

func (m *Settings) GetDirectorySecretSource() *Settings_Directory {
	if x, ok := m.GetSecretSource().(*Settings_DirectorySecretSource); ok {
		return x.DirectorySecretSource
	}
	return nil
}

func (m *Settings) GetKubernetesArtifactSource() *Settings_KubernetesConfigmaps {
	if x, ok := m.GetArtifactSource().(*Settings_KubernetesArtifactSource); ok {
		return x.KubernetesArtifactSource
	}
	return nil
}

func (m *Settings) GetDirectoryArtifactSource() *Settings_Directory {
	if x, ok := m.GetArtifactSource().(*Settings_DirectoryArtifactSource); ok {
		return x.DirectoryArtifactSource
	}
	return nil
}

func (m *Settings) GetBindAddr() string {
	if m != nil {
		return m.BindAddr
	}
	return ""
}

func (m *Settings) GetRefreshRate() *google_protobuf1.Duration {
	if m != nil {
		return m.RefreshRate
	}
	return nil
}

func (m *Settings) GetDevMode() bool {
	if m != nil {
		return m.DevMode
	}
	return false
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Settings) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Settings_OneofMarshaler, _Settings_OneofUnmarshaler, _Settings_OneofSizer, []interface{}{
		(*Settings_KubernetesConfigSource)(nil),
		(*Settings_DirectoryConfigSource)(nil),
		(*Settings_KubernetesSecretSource)(nil),
		(*Settings_VaultSecretSource)(nil),
		(*Settings_DirectorySecretSource)(nil),
		(*Settings_KubernetesArtifactSource)(nil),
		(*Settings_DirectoryArtifactSource)(nil),
	}
}

func _Settings_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Settings)
	// config_source
	switch x := m.ConfigSource.(type) {
	case *Settings_KubernetesConfigSource:
		_ = b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.KubernetesConfigSource); err != nil {
			return err
		}
	case *Settings_DirectoryConfigSource:
		_ = b.EncodeVarint(5<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.DirectoryConfigSource); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Settings.ConfigSource has unexpected type %T", x)
	}
	// secret_source
	switch x := m.SecretSource.(type) {
	case *Settings_KubernetesSecretSource:
		_ = b.EncodeVarint(6<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.KubernetesSecretSource); err != nil {
			return err
		}
	case *Settings_VaultSecretSource:
		_ = b.EncodeVarint(7<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.VaultSecretSource); err != nil {
			return err
		}
	case *Settings_DirectorySecretSource:
		_ = b.EncodeVarint(8<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.DirectorySecretSource); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Settings.SecretSource has unexpected type %T", x)
	}
	// artifact_source
	switch x := m.ArtifactSource.(type) {
	case *Settings_KubernetesArtifactSource:
		_ = b.EncodeVarint(9<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.KubernetesArtifactSource); err != nil {
			return err
		}
	case *Settings_DirectoryArtifactSource:
		_ = b.EncodeVarint(10<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.DirectoryArtifactSource); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Settings.ArtifactSource has unexpected type %T", x)
	}
	return nil
}

func _Settings_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Settings)
	switch tag {
	case 4: // config_source.kubernetes_config_source
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Settings_KubernetesCrds)
		err := b.DecodeMessage(msg)
		m.ConfigSource = &Settings_KubernetesConfigSource{msg}
		return true, err
	case 5: // config_source.directory_config_source
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Settings_Directory)
		err := b.DecodeMessage(msg)
		m.ConfigSource = &Settings_DirectoryConfigSource{msg}
		return true, err
	case 6: // secret_source.kubernetes_secret_source
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Settings_KubernetesSecrets)
		err := b.DecodeMessage(msg)
		m.SecretSource = &Settings_KubernetesSecretSource{msg}
		return true, err
	case 7: // secret_source.vault_secret_source
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Settings_VaultSecrets)
		err := b.DecodeMessage(msg)
		m.SecretSource = &Settings_VaultSecretSource{msg}
		return true, err
	case 8: // secret_source.directory_secret_source
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Settings_Directory)
		err := b.DecodeMessage(msg)
		m.SecretSource = &Settings_DirectorySecretSource{msg}
		return true, err
	case 9: // artifact_source.kubernetes_artifact_source
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Settings_KubernetesConfigmaps)
		err := b.DecodeMessage(msg)
		m.ArtifactSource = &Settings_KubernetesArtifactSource{msg}
		return true, err
	case 10: // artifact_source.directory_artifact_source
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Settings_Directory)
		err := b.DecodeMessage(msg)
		m.ArtifactSource = &Settings_DirectoryArtifactSource{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Settings_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Settings)
	// config_source
	switch x := m.ConfigSource.(type) {
	case *Settings_KubernetesConfigSource:
		s := proto.Size(x.KubernetesConfigSource)
		n += proto.SizeVarint(4<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Settings_DirectoryConfigSource:
		s := proto.Size(x.DirectoryConfigSource)
		n += proto.SizeVarint(5<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	// secret_source
	switch x := m.SecretSource.(type) {
	case *Settings_KubernetesSecretSource:
		s := proto.Size(x.KubernetesSecretSource)
		n += proto.SizeVarint(6<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Settings_VaultSecretSource:
		s := proto.Size(x.VaultSecretSource)
		n += proto.SizeVarint(7<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Settings_DirectorySecretSource:
		s := proto.Size(x.DirectorySecretSource)
		n += proto.SizeVarint(8<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	// artifact_source
	switch x := m.ArtifactSource.(type) {
	case *Settings_KubernetesArtifactSource:
		s := proto.Size(x.KubernetesArtifactSource)
		n += proto.SizeVarint(9<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Settings_DirectoryArtifactSource:
		s := proto.Size(x.DirectoryArtifactSource)
		n += proto.SizeVarint(10<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

// ilackarms(todo: make sure these are configurable)
type Settings_KubernetesCrds struct {
}

func (m *Settings_KubernetesCrds) Reset()         { *m = Settings_KubernetesCrds{} }
func (m *Settings_KubernetesCrds) String() string { return proto.CompactTextString(m) }
func (*Settings_KubernetesCrds) ProtoMessage()    {}
func (*Settings_KubernetesCrds) Descriptor() ([]byte, []int) {
	return fileDescriptorBootstrap, []int{0, 0}
}

type Settings_KubernetesSecrets struct {
}

func (m *Settings_KubernetesSecrets) Reset()         { *m = Settings_KubernetesSecrets{} }
func (m *Settings_KubernetesSecrets) String() string { return proto.CompactTextString(m) }
func (*Settings_KubernetesSecrets) ProtoMessage()    {}
func (*Settings_KubernetesSecrets) Descriptor() ([]byte, []int) {
	return fileDescriptorBootstrap, []int{0, 1}
}

type Settings_VaultSecrets struct {
}

func (m *Settings_VaultSecrets) Reset()         { *m = Settings_VaultSecrets{} }
func (m *Settings_VaultSecrets) String() string { return proto.CompactTextString(m) }
func (*Settings_VaultSecrets) ProtoMessage()    {}
func (*Settings_VaultSecrets) Descriptor() ([]byte, []int) {
	return fileDescriptorBootstrap, []int{0, 2}
}

type Settings_KubernetesConfigmaps struct {
}

func (m *Settings_KubernetesConfigmaps) Reset()         { *m = Settings_KubernetesConfigmaps{} }
func (m *Settings_KubernetesConfigmaps) String() string { return proto.CompactTextString(m) }
func (*Settings_KubernetesConfigmaps) ProtoMessage()    {}
func (*Settings_KubernetesConfigmaps) Descriptor() ([]byte, []int) {
	return fileDescriptorBootstrap, []int{0, 3}
}

type Settings_Directory struct {
}

func (m *Settings_Directory) Reset()                    { *m = Settings_Directory{} }
func (m *Settings_Directory) String() string            { return proto.CompactTextString(m) }
func (*Settings_Directory) ProtoMessage()               {}
func (*Settings_Directory) Descriptor() ([]byte, []int) { return fileDescriptorBootstrap, []int{0, 4} }

func init() {
	proto.RegisterType((*Settings)(nil), "gloo.solo.io.Settings")
	proto.RegisterType((*Settings_KubernetesCrds)(nil), "gloo.solo.io.Settings.KubernetesCrds")
	proto.RegisterType((*Settings_KubernetesSecrets)(nil), "gloo.solo.io.Settings.KubernetesSecrets")
	proto.RegisterType((*Settings_VaultSecrets)(nil), "gloo.solo.io.Settings.VaultSecrets")
	proto.RegisterType((*Settings_KubernetesConfigmaps)(nil), "gloo.solo.io.Settings.KubernetesConfigmaps")
	proto.RegisterType((*Settings_Directory)(nil), "gloo.solo.io.Settings.Directory")
}
func (this *Settings) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Settings)
	if !ok {
		that2, ok := that.(Settings)
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
	if this.DiscoveryNamespace != that1.DiscoveryNamespace {
		return false
	}
	if this.WatchNamespaces != that1.WatchNamespaces {
		return false
	}
	if that1.ConfigSource == nil {
		if this.ConfigSource != nil {
			return false
		}
	} else if this.ConfigSource == nil {
		return false
	} else if !this.ConfigSource.Equal(that1.ConfigSource) {
		return false
	}
	if that1.SecretSource == nil {
		if this.SecretSource != nil {
			return false
		}
	} else if this.SecretSource == nil {
		return false
	} else if !this.SecretSource.Equal(that1.SecretSource) {
		return false
	}
	if that1.ArtifactSource == nil {
		if this.ArtifactSource != nil {
			return false
		}
	} else if this.ArtifactSource == nil {
		return false
	} else if !this.ArtifactSource.Equal(that1.ArtifactSource) {
		return false
	}
	if this.BindAddr != that1.BindAddr {
		return false
	}
	if !this.RefreshRate.Equal(that1.RefreshRate) {
		return false
	}
	if this.DevMode != that1.DevMode {
		return false
	}
	return true
}
func (this *Settings_KubernetesConfigSource) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Settings_KubernetesConfigSource)
	if !ok {
		that2, ok := that.(Settings_KubernetesConfigSource)
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
	if !this.KubernetesConfigSource.Equal(that1.KubernetesConfigSource) {
		return false
	}
	return true
}
func (this *Settings_DirectoryConfigSource) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Settings_DirectoryConfigSource)
	if !ok {
		that2, ok := that.(Settings_DirectoryConfigSource)
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
	if !this.DirectoryConfigSource.Equal(that1.DirectoryConfigSource) {
		return false
	}
	return true
}
func (this *Settings_KubernetesSecretSource) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Settings_KubernetesSecretSource)
	if !ok {
		that2, ok := that.(Settings_KubernetesSecretSource)
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
	if !this.KubernetesSecretSource.Equal(that1.KubernetesSecretSource) {
		return false
	}
	return true
}
func (this *Settings_VaultSecretSource) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Settings_VaultSecretSource)
	if !ok {
		that2, ok := that.(Settings_VaultSecretSource)
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
	if !this.VaultSecretSource.Equal(that1.VaultSecretSource) {
		return false
	}
	return true
}
func (this *Settings_DirectorySecretSource) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Settings_DirectorySecretSource)
	if !ok {
		that2, ok := that.(Settings_DirectorySecretSource)
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
	if !this.DirectorySecretSource.Equal(that1.DirectorySecretSource) {
		return false
	}
	return true
}
func (this *Settings_KubernetesArtifactSource) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Settings_KubernetesArtifactSource)
	if !ok {
		that2, ok := that.(Settings_KubernetesArtifactSource)
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
	if !this.KubernetesArtifactSource.Equal(that1.KubernetesArtifactSource) {
		return false
	}
	return true
}
func (this *Settings_DirectoryArtifactSource) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Settings_DirectoryArtifactSource)
	if !ok {
		that2, ok := that.(Settings_DirectoryArtifactSource)
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
	if !this.DirectoryArtifactSource.Equal(that1.DirectoryArtifactSource) {
		return false
	}
	return true
}
func (this *Settings_KubernetesCrds) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Settings_KubernetesCrds)
	if !ok {
		that2, ok := that.(Settings_KubernetesCrds)
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
	return true
}
func (this *Settings_KubernetesSecrets) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Settings_KubernetesSecrets)
	if !ok {
		that2, ok := that.(Settings_KubernetesSecrets)
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
	return true
}
func (this *Settings_VaultSecrets) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Settings_VaultSecrets)
	if !ok {
		that2, ok := that.(Settings_VaultSecrets)
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
	return true
}
func (this *Settings_KubernetesConfigmaps) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Settings_KubernetesConfigmaps)
	if !ok {
		that2, ok := that.(Settings_KubernetesConfigmaps)
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
	return true
}
func (this *Settings_Directory) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Settings_Directory)
	if !ok {
		that2, ok := that.(Settings_Directory)
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
	return true
}

func init() { proto.RegisterFile("bootstrap.proto", fileDescriptorBootstrap) }

var fileDescriptorBootstrap = []byte{
	// 536 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x94, 0xdd, 0x6e, 0xd3, 0x30,
	0x14, 0xc7, 0x97, 0x09, 0xb6, 0xf6, 0xb4, 0x5b, 0xd7, 0x74, 0x8c, 0x2c, 0x48, 0x53, 0x05, 0x42,
	0x2a, 0x42, 0x24, 0xe2, 0xe3, 0x0a, 0x71, 0xb3, 0xae, 0x17, 0x95, 0x10, 0x5c, 0xa4, 0x82, 0x8b,
	0x5d, 0x50, 0x39, 0xb6, 0x9b, 0x9a, 0xb6, 0x39, 0x91, 0xed, 0x14, 0xed, 0x31, 0x78, 0x0b, 0x9e,
	0x8b, 0x27, 0x41, 0x71, 0x9b, 0xaf, 0xb2, 0x69, 0xbb, 0x6a, 0x6d, 0x9f, 0xf3, 0xfb, 0xfb, 0xfc,
	0x7d, 0x4e, 0xa0, 0x13, 0x22, 0x6a, 0xa5, 0x25, 0x49, 0xbc, 0x44, 0xa2, 0x46, 0xbb, 0x1d, 0x2d,
	0x11, 0x3d, 0x85, 0x4b, 0xf4, 0x04, 0xba, 0xa7, 0x11, 0x46, 0x68, 0x0e, 0xfc, 0xec, 0xdf, 0x26,
	0xc6, 0xbd, 0x88, 0x10, 0xa3, 0x25, 0xf7, 0xcd, 0x2a, 0x4c, 0x67, 0x3e, 0x4b, 0x25, 0xd1, 0x02,
	0xe3, 0xcd, 0xf9, 0xf3, 0xdf, 0x0d, 0x68, 0x4c, 0xb8, 0xd6, 0x22, 0x8e, 0x94, 0xed, 0x43, 0x8f,
	0x09, 0x45, 0x71, 0xcd, 0xe5, 0xcd, 0x34, 0x26, 0x2b, 0xae, 0x12, 0x42, 0xb9, 0x63, 0xf5, 0xad,
	0x41, 0x33, 0xb0, 0x8b, 0xa3, 0xaf, 0xf9, 0x89, 0xfd, 0x0a, 0x4e, 0x7e, 0x11, 0x4d, 0xe7, 0x65,
	0xb0, 0x72, 0xf6, 0x4d, 0x74, 0xc7, 0xec, 0x17, 0x91, 0xca, 0x26, 0xe0, 0x2c, 0xd2, 0x90, 0xcb,
	0x98, 0x6b, 0xae, 0xa6, 0x14, 0xe3, 0x99, 0x88, 0xa6, 0x0a, 0x53, 0x49, 0xb9, 0xf3, 0xa8, 0x6f,
	0x0d, 0x5a, 0xef, 0x5e, 0x7a, 0xd5, 0x7a, 0xbc, 0xfc, 0x56, 0xde, 0xe7, 0x22, 0xed, 0x4a, 0x32,
	0x35, 0xde, 0x0b, 0xce, 0x4a, 0xd0, 0x95, 0xe1, 0x4c, 0x0c, 0xc6, 0xbe, 0x86, 0xa7, 0x4c, 0x48,
	0x4e, 0x35, 0xca, 0x9b, 0x1d, 0x85, 0xc7, 0x46, 0xa1, 0x7f, 0x87, 0xc2, 0x28, 0xcf, 0x1a, 0xef,
	0x05, 0x4f, 0x0a, 0x44, 0x8d, 0xcd, 0x6a, 0xd7, 0x57, 0x9c, 0x4a, 0xae, 0x73, 0xf8, 0x81, 0x81,
	0x0f, 0xee, 0xbd, 0xfe, 0xc4, 0x64, 0xa9, 0xb1, 0x55, 0xad, 0x60, 0xb3, 0xb9, 0x55, 0xf9, 0x06,
	0xbd, 0x35, 0x49, 0x97, 0x7a, 0x47, 0xe0, 0xd0, 0x08, 0xbc, 0xb8, 0x43, 0xe0, 0x7b, 0x96, 0x51,
	0xb2, 0xbb, 0xeb, 0x72, 0x7d, 0x9b, 0x31, 0x75, 0x74, 0xe3, 0x81, 0xc6, 0x58, 0x15, 0x63, 0x6a,
	0xec, 0x05, 0xb8, 0x15, 0x63, 0x88, 0xd4, 0x62, 0x46, 0x68, 0x81, 0x6f, 0x1a, 0xfc, 0xeb, 0xfb,
	0x5f, 0xd6, 0x78, 0xbd, 0x22, 0x89, 0x1a, 0xef, 0x07, 0x15, 0xa7, 0x2f, 0xb7, 0xbc, 0xad, 0xd8,
	0x0f, 0x38, 0x2f, 0x0b, 0xd9, 0xd5, 0x82, 0x07, 0x96, 0xb2, 0x1f, 0x94, 0x6e, 0xec, 0xf0, 0x9f,
	0x41, 0x33, 0x14, 0x31, 0x9b, 0x12, 0xc6, 0xa4, 0xd3, 0x32, 0x8d, 0xdc, 0xc8, 0x36, 0x2e, 0x19,
	0x93, 0xf6, 0x27, 0x68, 0x4b, 0x3e, 0x93, 0x5c, 0xcd, 0xa7, 0x92, 0x68, 0xee, 0xb4, 0x8d, 0xde,
	0xb9, 0xb7, 0x99, 0x30, 0x2f, 0x9f, 0x30, 0x6f, 0xb4, 0x9d, 0xb0, 0xa0, 0xb5, 0x0d, 0x0f, 0x88,
	0xe6, 0xb6, 0x03, 0x87, 0x23, 0xbe, 0xfe, 0x82, 0x8c, 0x3b, 0x47, 0x7d, 0x6b, 0xd0, 0x08, 0xf2,
	0xa5, 0x7b, 0x02, 0xc7, 0xf5, 0x16, 0x77, 0x7b, 0xd0, 0xfd, 0xaf, 0x6b, 0xdc, 0x63, 0x68, 0x57,
	0x5f, 0xda, 0x3d, 0x83, 0xd3, 0xdb, 0xfc, 0x73, 0x5b, 0xd0, 0x2c, 0x6a, 0x1d, 0x76, 0xe0, 0xa8,
	0x36, 0x08, 0xd9, 0x46, 0xad, 0x01, 0x86, 0x5d, 0xe8, 0xec, 0x18, 0x39, 0xfc, 0xf8, 0xe7, 0xef,
	0x85, 0x75, 0xfd, 0x21, 0x12, 0x7a, 0x9e, 0x86, 0x1e, 0xc5, 0x95, 0x9f, 0xb9, 0xf9, 0x46, 0xe0,
	0xe6, 0x77, 0x21, 0x74, 0xf6, 0x39, 0xf9, 0xc9, 0xa9, 0x56, 0x7e, 0x66, 0xb6, 0x9f, 0x2c, 0x22,
	0x9f, 0x24, 0xc2, 0x5f, 0xbf, 0x0d, 0x0f, 0x8c, 0x0d, 0xef, 0xff, 0x05, 0x00, 0x00, 0xff, 0xff,
	0x94, 0x86, 0x43, 0x15, 0xad, 0x04, 0x00, 0x00,
}
