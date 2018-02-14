// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: virtualhost.proto

package v1

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf1 "github.com/gogo/protobuf/types"
import _ "github.com/gogo/protobuf/gogoproto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type VirtualHost struct {
	Name      string     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Domains   []string   `protobuf:"bytes,2,rep,name=domains" json:"domains,omitempty"`
	Routes    []*Route   `protobuf:"bytes,3,rep,name=routes" json:"routes,omitempty"`
	SslConfig *SSLConfig `protobuf:"bytes,4,opt,name=ssl_config,json=sslConfig" json:"ssl_config,omitempty"`
}

func (m *VirtualHost) Reset()                    { *m = VirtualHost{} }
func (m *VirtualHost) String() string            { return proto.CompactTextString(m) }
func (*VirtualHost) ProtoMessage()               {}
func (*VirtualHost) Descriptor() ([]byte, []int) { return fileDescriptorVirtualhost, []int{0} }

func (m *VirtualHost) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *VirtualHost) GetDomains() []string {
	if m != nil {
		return m.Domains
	}
	return nil
}

func (m *VirtualHost) GetRoutes() []*Route {
	if m != nil {
		return m.Routes
	}
	return nil
}

func (m *VirtualHost) GetSslConfig() *SSLConfig {
	if m != nil {
		return m.SslConfig
	}
	return nil
}

type Route struct {
	Matcher *Matcher `protobuf:"bytes,1,opt,name=matcher" json:"matcher,omitempty"`
	// Types that are valid to be assigned to Destination:
	//	*Route_MultipleDestinations
	//	*Route_SingleDestination
	Destination   isRoute_Destination      `protobuf_oneof:"destination"`
	PrefixRewrite string                   `protobuf:"bytes,4,opt,name=prefix_rewrite,json=prefixRewrite,proto3" json:"prefix_rewrite,omitempty"`
	Extensions    *google_protobuf1.Struct `protobuf:"bytes,5,opt,name=extensions" json:"extensions,omitempty"`
}

func (m *Route) Reset()                    { *m = Route{} }
func (m *Route) String() string            { return proto.CompactTextString(m) }
func (*Route) ProtoMessage()               {}
func (*Route) Descriptor() ([]byte, []int) { return fileDescriptorVirtualhost, []int{1} }

type isRoute_Destination interface {
	isRoute_Destination()
	Equal(interface{}) bool
}

type Route_MultipleDestinations struct {
	MultipleDestinations *MultipleDestinations `protobuf:"bytes,2,opt,name=multiple_destinations,json=multipleDestinations,oneof"`
}
type Route_SingleDestination struct {
	SingleDestination *SingleDestination `protobuf:"bytes,3,opt,name=single_destination,json=singleDestination,oneof"`
}

func (*Route_MultipleDestinations) isRoute_Destination() {}
func (*Route_SingleDestination) isRoute_Destination()    {}

func (m *Route) GetDestination() isRoute_Destination {
	if m != nil {
		return m.Destination
	}
	return nil
}

func (m *Route) GetMatcher() *Matcher {
	if m != nil {
		return m.Matcher
	}
	return nil
}

func (m *Route) GetMultipleDestinations() *MultipleDestinations {
	if x, ok := m.GetDestination().(*Route_MultipleDestinations); ok {
		return x.MultipleDestinations
	}
	return nil
}

func (m *Route) GetSingleDestination() *SingleDestination {
	if x, ok := m.GetDestination().(*Route_SingleDestination); ok {
		return x.SingleDestination
	}
	return nil
}

func (m *Route) GetPrefixRewrite() string {
	if m != nil {
		return m.PrefixRewrite
	}
	return ""
}

func (m *Route) GetExtensions() *google_protobuf1.Struct {
	if m != nil {
		return m.Extensions
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Route) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Route_OneofMarshaler, _Route_OneofUnmarshaler, _Route_OneofSizer, []interface{}{
		(*Route_MultipleDestinations)(nil),
		(*Route_SingleDestination)(nil),
	}
}

func _Route_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Route)
	// destination
	switch x := m.Destination.(type) {
	case *Route_MultipleDestinations:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.MultipleDestinations); err != nil {
			return err
		}
	case *Route_SingleDestination:
		_ = b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.SingleDestination); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("Route.Destination has unexpected type %T", x)
	}
	return nil
}

func _Route_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Route)
	switch tag {
	case 2: // destination.multiple_destinations
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(MultipleDestinations)
		err := b.DecodeMessage(msg)
		m.Destination = &Route_MultipleDestinations{msg}
		return true, err
	case 3: // destination.single_destination
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(SingleDestination)
		err := b.DecodeMessage(msg)
		m.Destination = &Route_SingleDestination{msg}
		return true, err
	default:
		return false, nil
	}
}

func _Route_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Route)
	// destination
	switch x := m.Destination.(type) {
	case *Route_MultipleDestinations:
		s := proto.Size(x.MultipleDestinations)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *Route_SingleDestination:
		s := proto.Size(x.SingleDestination)
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Matcher struct {
	// Types that are valid to be assigned to Path:
	//	*Matcher_PathPrefix
	//	*Matcher_PathRegex
	//	*Matcher_PathExact
	Path        isMatcher_Path    `protobuf_oneof:"path"`
	Headers     map[string]string `protobuf:"bytes,4,rep,name=headers" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	QueryParams map[string]string `protobuf:"bytes,5,rep,name=query_params,json=queryParams" json:"query_params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Verbs       []string          `protobuf:"bytes,6,rep,name=verbs" json:"verbs,omitempty"`
}

func (m *Matcher) Reset()                    { *m = Matcher{} }
func (m *Matcher) String() string            { return proto.CompactTextString(m) }
func (*Matcher) ProtoMessage()               {}
func (*Matcher) Descriptor() ([]byte, []int) { return fileDescriptorVirtualhost, []int{2} }

type isMatcher_Path interface {
	isMatcher_Path()
	Equal(interface{}) bool
}

type Matcher_PathPrefix struct {
	PathPrefix string `protobuf:"bytes,1,opt,name=path_prefix,json=pathPrefix,proto3,oneof"`
}
type Matcher_PathRegex struct {
	PathRegex string `protobuf:"bytes,2,opt,name=path_regex,json=pathRegex,proto3,oneof"`
}
type Matcher_PathExact struct {
	PathExact string `protobuf:"bytes,3,opt,name=path_exact,json=pathExact,proto3,oneof"`
}

func (*Matcher_PathPrefix) isMatcher_Path() {}
func (*Matcher_PathRegex) isMatcher_Path()  {}
func (*Matcher_PathExact) isMatcher_Path()  {}

func (m *Matcher) GetPath() isMatcher_Path {
	if m != nil {
		return m.Path
	}
	return nil
}

func (m *Matcher) GetPathPrefix() string {
	if x, ok := m.GetPath().(*Matcher_PathPrefix); ok {
		return x.PathPrefix
	}
	return ""
}

func (m *Matcher) GetPathRegex() string {
	if x, ok := m.GetPath().(*Matcher_PathRegex); ok {
		return x.PathRegex
	}
	return ""
}

func (m *Matcher) GetPathExact() string {
	if x, ok := m.GetPath().(*Matcher_PathExact); ok {
		return x.PathExact
	}
	return ""
}

func (m *Matcher) GetHeaders() map[string]string {
	if m != nil {
		return m.Headers
	}
	return nil
}

func (m *Matcher) GetQueryParams() map[string]string {
	if m != nil {
		return m.QueryParams
	}
	return nil
}

func (m *Matcher) GetVerbs() []string {
	if m != nil {
		return m.Verbs
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*Matcher) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _Matcher_OneofMarshaler, _Matcher_OneofUnmarshaler, _Matcher_OneofSizer, []interface{}{
		(*Matcher_PathPrefix)(nil),
		(*Matcher_PathRegex)(nil),
		(*Matcher_PathExact)(nil),
	}
}

func _Matcher_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*Matcher)
	// path
	switch x := m.Path.(type) {
	case *Matcher_PathPrefix:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		_ = b.EncodeStringBytes(x.PathPrefix)
	case *Matcher_PathRegex:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		_ = b.EncodeStringBytes(x.PathRegex)
	case *Matcher_PathExact:
		_ = b.EncodeVarint(3<<3 | proto.WireBytes)
		_ = b.EncodeStringBytes(x.PathExact)
	case nil:
	default:
		return fmt.Errorf("Matcher.Path has unexpected type %T", x)
	}
	return nil
}

func _Matcher_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*Matcher)
	switch tag {
	case 1: // path.path_prefix
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Path = &Matcher_PathPrefix{x}
		return true, err
	case 2: // path.path_regex
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Path = &Matcher_PathRegex{x}
		return true, err
	case 3: // path.path_exact
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Path = &Matcher_PathExact{x}
		return true, err
	default:
		return false, nil
	}
}

func _Matcher_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*Matcher)
	// path
	switch x := m.Path.(type) {
	case *Matcher_PathPrefix:
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.PathPrefix)))
		n += len(x.PathPrefix)
	case *Matcher_PathRegex:
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.PathRegex)))
		n += len(x.PathRegex)
	case *Matcher_PathExact:
		n += proto.SizeVarint(3<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(len(x.PathExact)))
		n += len(x.PathExact)
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type MultipleDestinations struct {
	WeightedDestinations []*WeightedSingleDestination `protobuf:"bytes,1,rep,name=weighted_destinations,json=weightedDestinations" json:"weighted_destinations,omitempty"`
}

func (m *MultipleDestinations) Reset()                    { *m = MultipleDestinations{} }
func (m *MultipleDestinations) String() string            { return proto.CompactTextString(m) }
func (*MultipleDestinations) ProtoMessage()               {}
func (*MultipleDestinations) Descriptor() ([]byte, []int) { return fileDescriptorVirtualhost, []int{3} }

func (m *MultipleDestinations) GetWeightedDestinations() []*WeightedSingleDestination {
	if m != nil {
		return m.WeightedDestinations
	}
	return nil
}

type WeightedSingleDestination struct {
	Destination *SingleDestination `protobuf:"bytes,1,opt,name=destination" json:"destination,omitempty"`
	Weight      uint32             `protobuf:"varint,2,opt,name=weight,proto3" json:"weight,omitempty"`
}

func (m *WeightedSingleDestination) Reset()         { *m = WeightedSingleDestination{} }
func (m *WeightedSingleDestination) String() string { return proto.CompactTextString(m) }
func (*WeightedSingleDestination) ProtoMessage()    {}
func (*WeightedSingleDestination) Descriptor() ([]byte, []int) {
	return fileDescriptorVirtualhost, []int{4}
}

func (m *WeightedSingleDestination) GetDestination() *SingleDestination {
	if m != nil {
		return m.Destination
	}
	return nil
}

func (m *WeightedSingleDestination) GetWeight() uint32 {
	if m != nil {
		return m.Weight
	}
	return 0
}

type SingleDestination struct {
	// Types that are valid to be assigned to Destination:
	//	*SingleDestination_Function
	//	*SingleDestination_Upstream
	Destination isSingleDestination_Destination `protobuf_oneof:"destination"`
}

func (m *SingleDestination) Reset()                    { *m = SingleDestination{} }
func (m *SingleDestination) String() string            { return proto.CompactTextString(m) }
func (*SingleDestination) ProtoMessage()               {}
func (*SingleDestination) Descriptor() ([]byte, []int) { return fileDescriptorVirtualhost, []int{5} }

type isSingleDestination_Destination interface {
	isSingleDestination_Destination()
	Equal(interface{}) bool
}

type SingleDestination_Function struct {
	Function *FunctionDestination `protobuf:"bytes,1,opt,name=function,oneof"`
}
type SingleDestination_Upstream struct {
	Upstream *UpstreamDestination `protobuf:"bytes,2,opt,name=upstream,oneof"`
}

func (*SingleDestination_Function) isSingleDestination_Destination() {}
func (*SingleDestination_Upstream) isSingleDestination_Destination() {}

func (m *SingleDestination) GetDestination() isSingleDestination_Destination {
	if m != nil {
		return m.Destination
	}
	return nil
}

func (m *SingleDestination) GetFunction() *FunctionDestination {
	if x, ok := m.GetDestination().(*SingleDestination_Function); ok {
		return x.Function
	}
	return nil
}

func (m *SingleDestination) GetUpstream() *UpstreamDestination {
	if x, ok := m.GetDestination().(*SingleDestination_Upstream); ok {
		return x.Upstream
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*SingleDestination) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _SingleDestination_OneofMarshaler, _SingleDestination_OneofUnmarshaler, _SingleDestination_OneofSizer, []interface{}{
		(*SingleDestination_Function)(nil),
		(*SingleDestination_Upstream)(nil),
	}
}

func _SingleDestination_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*SingleDestination)
	// destination
	switch x := m.Destination.(type) {
	case *SingleDestination_Function:
		_ = b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Function); err != nil {
			return err
		}
	case *SingleDestination_Upstream:
		_ = b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Upstream); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("SingleDestination.Destination has unexpected type %T", x)
	}
	return nil
}

func _SingleDestination_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*SingleDestination)
	switch tag {
	case 1: // destination.function
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(FunctionDestination)
		err := b.DecodeMessage(msg)
		m.Destination = &SingleDestination_Function{msg}
		return true, err
	case 2: // destination.upstream
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(UpstreamDestination)
		err := b.DecodeMessage(msg)
		m.Destination = &SingleDestination_Upstream{msg}
		return true, err
	default:
		return false, nil
	}
}

func _SingleDestination_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*SingleDestination)
	// destination
	switch x := m.Destination.(type) {
	case *SingleDestination_Function:
		s := proto.Size(x.Function)
		n += proto.SizeVarint(1<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case *SingleDestination_Upstream:
		s := proto.Size(x.Upstream)
		n += proto.SizeVarint(2<<3 | proto.WireBytes)
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type FunctionDestination struct {
	UpstreamName string `protobuf:"bytes,1,opt,name=upstream_name,json=upstreamName,proto3" json:"upstream_name,omitempty"`
	FunctionName string `protobuf:"bytes,2,opt,name=function_name,json=functionName,proto3" json:"function_name,omitempty"`
}

func (m *FunctionDestination) Reset()                    { *m = FunctionDestination{} }
func (m *FunctionDestination) String() string            { return proto.CompactTextString(m) }
func (*FunctionDestination) ProtoMessage()               {}
func (*FunctionDestination) Descriptor() ([]byte, []int) { return fileDescriptorVirtualhost, []int{6} }

func (m *FunctionDestination) GetUpstreamName() string {
	if m != nil {
		return m.UpstreamName
	}
	return ""
}

func (m *FunctionDestination) GetFunctionName() string {
	if m != nil {
		return m.FunctionName
	}
	return ""
}

type UpstreamDestination struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (m *UpstreamDestination) Reset()                    { *m = UpstreamDestination{} }
func (m *UpstreamDestination) String() string            { return proto.CompactTextString(m) }
func (*UpstreamDestination) ProtoMessage()               {}
func (*UpstreamDestination) Descriptor() ([]byte, []int) { return fileDescriptorVirtualhost, []int{7} }

func (m *UpstreamDestination) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type SSLConfig struct {
	SecretRef string `protobuf:"bytes,1,opt,name=secret_ref,json=secretRef,proto3" json:"secret_ref,omitempty"`
}

func (m *SSLConfig) Reset()                    { *m = SSLConfig{} }
func (m *SSLConfig) String() string            { return proto.CompactTextString(m) }
func (*SSLConfig) ProtoMessage()               {}
func (*SSLConfig) Descriptor() ([]byte, []int) { return fileDescriptorVirtualhost, []int{8} }

func (m *SSLConfig) GetSecretRef() string {
	if m != nil {
		return m.SecretRef
	}
	return ""
}

func init() {
	proto.RegisterType((*VirtualHost)(nil), "v1.VirtualHost")
	proto.RegisterType((*Route)(nil), "v1.Route")
	proto.RegisterType((*Matcher)(nil), "v1.Matcher")
	proto.RegisterType((*MultipleDestinations)(nil), "v1.MultipleDestinations")
	proto.RegisterType((*WeightedSingleDestination)(nil), "v1.WeightedSingleDestination")
	proto.RegisterType((*SingleDestination)(nil), "v1.SingleDestination")
	proto.RegisterType((*FunctionDestination)(nil), "v1.FunctionDestination")
	proto.RegisterType((*UpstreamDestination)(nil), "v1.UpstreamDestination")
	proto.RegisterType((*SSLConfig)(nil), "v1.SSLConfig")
}
func (this *VirtualHost) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*VirtualHost)
	if !ok {
		that2, ok := that.(VirtualHost)
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
	if this.Name != that1.Name {
		return false
	}
	if len(this.Domains) != len(that1.Domains) {
		return false
	}
	for i := range this.Domains {
		if this.Domains[i] != that1.Domains[i] {
			return false
		}
	}
	if len(this.Routes) != len(that1.Routes) {
		return false
	}
	for i := range this.Routes {
		if !this.Routes[i].Equal(that1.Routes[i]) {
			return false
		}
	}
	if !this.SslConfig.Equal(that1.SslConfig) {
		return false
	}
	return true
}
func (this *Route) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Route)
	if !ok {
		that2, ok := that.(Route)
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
	if !this.Matcher.Equal(that1.Matcher) {
		return false
	}
	if that1.Destination == nil {
		if this.Destination != nil {
			return false
		}
	} else if this.Destination == nil {
		return false
	} else if !this.Destination.Equal(that1.Destination) {
		return false
	}
	if this.PrefixRewrite != that1.PrefixRewrite {
		return false
	}
	if !this.Extensions.Equal(that1.Extensions) {
		return false
	}
	return true
}
func (this *Route_MultipleDestinations) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Route_MultipleDestinations)
	if !ok {
		that2, ok := that.(Route_MultipleDestinations)
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
	if !this.MultipleDestinations.Equal(that1.MultipleDestinations) {
		return false
	}
	return true
}
func (this *Route_SingleDestination) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Route_SingleDestination)
	if !ok {
		that2, ok := that.(Route_SingleDestination)
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
	if !this.SingleDestination.Equal(that1.SingleDestination) {
		return false
	}
	return true
}
func (this *Matcher) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Matcher)
	if !ok {
		that2, ok := that.(Matcher)
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
	if that1.Path == nil {
		if this.Path != nil {
			return false
		}
	} else if this.Path == nil {
		return false
	} else if !this.Path.Equal(that1.Path) {
		return false
	}
	if len(this.Headers) != len(that1.Headers) {
		return false
	}
	for i := range this.Headers {
		if this.Headers[i] != that1.Headers[i] {
			return false
		}
	}
	if len(this.QueryParams) != len(that1.QueryParams) {
		return false
	}
	for i := range this.QueryParams {
		if this.QueryParams[i] != that1.QueryParams[i] {
			return false
		}
	}
	if len(this.Verbs) != len(that1.Verbs) {
		return false
	}
	for i := range this.Verbs {
		if this.Verbs[i] != that1.Verbs[i] {
			return false
		}
	}
	return true
}
func (this *Matcher_PathPrefix) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Matcher_PathPrefix)
	if !ok {
		that2, ok := that.(Matcher_PathPrefix)
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
	if this.PathPrefix != that1.PathPrefix {
		return false
	}
	return true
}
func (this *Matcher_PathRegex) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Matcher_PathRegex)
	if !ok {
		that2, ok := that.(Matcher_PathRegex)
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
	if this.PathRegex != that1.PathRegex {
		return false
	}
	return true
}
func (this *Matcher_PathExact) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Matcher_PathExact)
	if !ok {
		that2, ok := that.(Matcher_PathExact)
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
	if this.PathExact != that1.PathExact {
		return false
	}
	return true
}
func (this *MultipleDestinations) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MultipleDestinations)
	if !ok {
		that2, ok := that.(MultipleDestinations)
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
	if len(this.WeightedDestinations) != len(that1.WeightedDestinations) {
		return false
	}
	for i := range this.WeightedDestinations {
		if !this.WeightedDestinations[i].Equal(that1.WeightedDestinations[i]) {
			return false
		}
	}
	return true
}
func (this *WeightedSingleDestination) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*WeightedSingleDestination)
	if !ok {
		that2, ok := that.(WeightedSingleDestination)
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
	if !this.Destination.Equal(that1.Destination) {
		return false
	}
	if this.Weight != that1.Weight {
		return false
	}
	return true
}
func (this *SingleDestination) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SingleDestination)
	if !ok {
		that2, ok := that.(SingleDestination)
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
	if that1.Destination == nil {
		if this.Destination != nil {
			return false
		}
	} else if this.Destination == nil {
		return false
	} else if !this.Destination.Equal(that1.Destination) {
		return false
	}
	return true
}
func (this *SingleDestination_Function) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SingleDestination_Function)
	if !ok {
		that2, ok := that.(SingleDestination_Function)
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
	if !this.Function.Equal(that1.Function) {
		return false
	}
	return true
}
func (this *SingleDestination_Upstream) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SingleDestination_Upstream)
	if !ok {
		that2, ok := that.(SingleDestination_Upstream)
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
	if !this.Upstream.Equal(that1.Upstream) {
		return false
	}
	return true
}
func (this *FunctionDestination) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*FunctionDestination)
	if !ok {
		that2, ok := that.(FunctionDestination)
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
	if this.UpstreamName != that1.UpstreamName {
		return false
	}
	if this.FunctionName != that1.FunctionName {
		return false
	}
	return true
}
func (this *UpstreamDestination) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*UpstreamDestination)
	if !ok {
		that2, ok := that.(UpstreamDestination)
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
	if this.Name != that1.Name {
		return false
	}
	return true
}
func (this *SSLConfig) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*SSLConfig)
	if !ok {
		that2, ok := that.(SSLConfig)
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
	if this.SecretRef != that1.SecretRef {
		return false
	}
	return true
}

func init() { proto.RegisterFile("virtualhost.proto", fileDescriptorVirtualhost) }

var fileDescriptorVirtualhost = []byte{
	// 690 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xcd, 0x6e, 0xd3, 0x40,
	0x10, 0xae, 0x93, 0x34, 0xc1, 0xe3, 0x04, 0xb5, 0xdb, 0x94, 0x9a, 0xa8, 0x85, 0xd4, 0x55, 0xa5,
	0x80, 0x50, 0xaa, 0x06, 0xa1, 0xa2, 0x1e, 0x40, 0x2a, 0xb4, 0xca, 0x81, 0x9f, 0xb2, 0x15, 0x70,
	0xb4, 0xdc, 0x64, 0xe3, 0x18, 0x1c, 0x3b, 0xdd, 0x5d, 0xa7, 0xe9, 0x3b, 0x70, 0xe0, 0x31, 0x78,
	0x02, 0x1e, 0x88, 0xf7, 0x40, 0x42, 0x3b, 0x6b, 0x07, 0x27, 0x0d, 0x07, 0x6e, 0x33, 0xdf, 0x7c,
	0xdf, 0xce, 0xcc, 0xce, 0xec, 0xc2, 0xfa, 0x24, 0xe0, 0x32, 0xf1, 0xc2, 0x61, 0x2c, 0x64, 0x7b,
	0xcc, 0x63, 0x19, 0x93, 0xc2, 0xe4, 0xb0, 0xb1, 0xed, 0xc7, 0xb1, 0x1f, 0xb2, 0x03, 0x44, 0x2e,
	0x93, 0xc1, 0x81, 0x90, 0x3c, 0xe9, 0xa5, 0x8c, 0x46, 0xdd, 0x8f, 0xfd, 0x18, 0xcd, 0x03, 0x65,
	0x69, 0xd4, 0xf9, 0x66, 0x80, 0xf5, 0x49, 0x9f, 0xd6, 0x8d, 0x85, 0x24, 0x04, 0x4a, 0x91, 0x37,
	0x62, 0xb6, 0xd1, 0x34, 0x5a, 0x26, 0x45, 0x9b, 0xd8, 0x50, 0xe9, 0xc7, 0x23, 0x2f, 0x88, 0x84,
	0x5d, 0x68, 0x16, 0x5b, 0x26, 0xcd, 0x5c, 0xb2, 0x0b, 0x65, 0x1e, 0x27, 0x92, 0x09, 0xbb, 0xd8,
	0x2c, 0xb6, 0xac, 0x8e, 0xd9, 0x9e, 0x1c, 0xb6, 0xa9, 0x42, 0x68, 0x1a, 0x20, 0x4f, 0x00, 0x84,
	0x08, 0xdd, 0x5e, 0x1c, 0x0d, 0x02, 0xdf, 0x2e, 0x35, 0x8d, 0x96, 0xd5, 0xa9, 0x29, 0xda, 0xc5,
	0xc5, 0x9b, 0x57, 0x08, 0x52, 0x53, 0x88, 0x50, 0x9b, 0xce, 0xcf, 0x02, 0xac, 0xa2, 0x9e, 0xec,
	0x43, 0x65, 0xe4, 0xc9, 0xde, 0x90, 0x71, 0xac, 0xc5, 0xea, 0x58, 0x4a, 0xf4, 0x56, 0x43, 0x34,
	0x8b, 0x91, 0xf7, 0xb0, 0x39, 0x4a, 0x42, 0x19, 0x8c, 0x43, 0xe6, 0xf6, 0x99, 0x90, 0x41, 0xe4,
	0xc9, 0x20, 0xc6, 0x4a, 0x95, 0xc8, 0x46, 0x51, 0x4a, 0x78, 0x9d, 0x8b, 0x77, 0x57, 0x68, 0x7d,
	0xb4, 0x04, 0x27, 0x67, 0x40, 0x44, 0x10, 0xf9, 0xf3, 0xc7, 0xd9, 0x45, 0x3c, 0x6d, 0x13, 0xeb,
	0xc6, 0x68, 0x4e, 0xd3, 0x5d, 0xa1, 0xeb, 0x62, 0x11, 0x24, 0xfb, 0x70, 0x77, 0xcc, 0xd9, 0x20,
	0x98, 0xba, 0x9c, 0x5d, 0xf3, 0x40, 0x32, 0xec, 0xdd, 0xa4, 0x35, 0x8d, 0x52, 0x0d, 0x92, 0x23,
	0x00, 0x36, 0x95, 0x2c, 0x12, 0x58, 0xf4, 0x2a, 0xa6, 0xd9, 0x6a, 0xeb, 0x41, 0xb6, 0xb3, 0x41,
	0xb6, 0x2f, 0x70, 0x90, 0x34, 0x47, 0x3d, 0xa9, 0x81, 0x95, 0x2b, 0xd0, 0xf9, 0x5d, 0x80, 0x4a,
	0x7a, 0x39, 0x64, 0x17, 0xac, 0xb1, 0x27, 0x87, 0xae, 0xce, 0xa4, 0x47, 0xd9, 0x5d, 0xa1, 0xa0,
	0xc0, 0x73, 0xc4, 0xc8, 0x43, 0x40, 0xcf, 0xe5, 0xcc, 0x67, 0x53, 0xbc, 0x2b, 0xc5, 0x30, 0x15,
	0x46, 0x15, 0x34, 0x23, 0xb0, 0xa9, 0xd7, 0x93, 0xd8, 0xfe, 0x8c, 0x70, 0xaa, 0x20, 0xd2, 0x81,
	0xca, 0x90, 0x79, 0x7d, 0xc6, 0x85, 0x5d, 0xc2, 0xd9, 0xdb, 0xb9, 0xf9, 0xb4, 0xbb, 0x3a, 0x74,
	0x1a, 0x49, 0x7e, 0x43, 0x33, 0x22, 0x79, 0x09, 0xd5, 0xab, 0x84, 0xf1, 0x1b, 0x77, 0xec, 0x71,
	0x6f, 0xa4, 0xda, 0x55, 0xc2, 0xed, 0xbc, 0xf0, 0x83, 0x8a, 0x9f, 0x63, 0x58, 0x8b, 0xad, 0xab,
	0xbf, 0x08, 0xa9, 0xc3, 0xea, 0x84, 0xf1, 0x4b, 0x61, 0x97, 0x71, 0x0f, 0xb5, 0xd3, 0x38, 0x86,
	0x6a, 0x3e, 0x1f, 0x59, 0x83, 0xe2, 0x57, 0x76, 0x93, 0xae, 0xb0, 0x32, 0x51, 0xe7, 0x85, 0x09,
	0xd3, 0x9d, 0x52, 0xed, 0x1c, 0x17, 0x9e, 0x1b, 0x8d, 0x17, 0xb0, 0xb6, 0x98, 0xf2, 0x7f, 0xf4,
	0x27, 0x65, 0x28, 0xa9, 0x3b, 0x71, 0xbe, 0x40, 0x7d, 0xd9, 0x9a, 0x11, 0x0a, 0x9b, 0xd7, 0x2c,
	0xf0, 0x87, 0x92, 0xf5, 0xe7, 0xf7, 0xd3, 0xc0, 0xde, 0x77, 0x54, 0xef, 0x9f, 0x53, 0xc2, 0xad,
	0xcd, 0xa2, 0xf5, 0x4c, 0x9b, 0x3f, 0xd3, 0x09, 0xe1, 0xfe, 0x3f, 0x25, 0xe4, 0x68, 0x6e, 0x2f,
	0xd2, 0xb7, 0xb3, 0x7c, 0x71, 0x69, 0x9e, 0x49, 0xee, 0x41, 0x59, 0x67, 0xc3, 0x26, 0x6b, 0x34,
	0xf5, 0x9c, 0xef, 0x06, 0xac, 0xdf, 0x4e, 0xf3, 0x0c, 0xee, 0x0c, 0x92, 0xa8, 0x97, 0xcb, 0xb1,
	0xa5, 0x72, 0x9c, 0xa5, 0xd8, 0xfc, 0xf3, 0x98, 0x51, 0x95, 0x2c, 0x19, 0x0b, 0xc9, 0x99, 0x37,
	0x4a, 0x5f, 0x28, 0xca, 0x3e, 0xa6, 0xd8, 0x82, 0x2c, 0xa3, 0x2e, 0x2e, 0xbb, 0x0b, 0x1b, 0x4b,
	0x12, 0x91, 0x3d, 0xa8, 0x65, 0x0a, 0x37, 0xf7, 0x89, 0x55, 0x33, 0xf0, 0x9d, 0xfa, 0xcc, 0xf6,
	0xa0, 0x96, 0x55, 0xa3, 0x49, 0x7a, 0xa4, 0xd5, 0x0c, 0x54, 0x24, 0xe7, 0x11, 0x6c, 0x2c, 0x29,
	0x69, 0xd9, 0xe7, 0xe8, 0x3c, 0x06, 0x73, 0xf6, 0x93, 0x91, 0x1d, 0x00, 0xc1, 0x7a, 0x9c, 0x49,
	0x97, 0xb3, 0x41, 0x4a, 0x33, 0x35, 0x42, 0xd9, 0xe0, 0xa4, 0xf4, 0xe3, 0xd7, 0x03, 0xe3, 0xb2,
	0x8c, 0xcf, 0xfa, 0xe9, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb5, 0x37, 0x49, 0x83, 0xc6, 0x05,
	0x00, 0x00,
}
