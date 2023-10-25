// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.6.1
// source: github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/rt_selector.proto

package v1

import (
	context "context"
	reflect "reflect"
	sync "sync"

	_ "github.com/solo-io/protoc-gen-ext/extproto"
	v12 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	v11 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	matchers "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/core/matchers"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SubRouteTableRow struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Action:
	//
	//	*SubRouteTableRow_RouteAction
	//	*SubRouteTableRow_RedirectAction
	//	*SubRouteTableRow_DirectResponseAction
	//	*SubRouteTableRow_DelegateAction
	Action          isSubRouteTableRow_Action         `protobuf_oneof:"action"`
	Matcher         string                            `protobuf:"bytes,5,opt,name=matcher,proto3" json:"matcher,omitempty"`
	MatchType       string                            `protobuf:"bytes,6,opt,name=match_type,json=matchType,proto3" json:"match_type,omitempty"`
	Methods         []string                          `protobuf:"bytes,7,rep,name=methods,proto3" json:"methods,omitempty"`
	Headers         []*matchers.HeaderMatcher         `protobuf:"bytes,8,rep,name=headers,proto3" json:"headers,omitempty"`
	QueryParameters []*matchers.QueryParameterMatcher `protobuf:"bytes,9,rep,name=query_parameters,json=queryParameters,proto3" json:"query_parameters,omitempty"`
	RtRoutes        []*SubRouteTableRow               `protobuf:"bytes,10,rep,name=rt_routes,json=rtRoutes,proto3" json:"rt_routes,omitempty"`
}

func (x *SubRouteTableRow) Reset() {
	*x = SubRouteTableRow{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubRouteTableRow) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubRouteTableRow) ProtoMessage() {}

func (x *SubRouteTableRow) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubRouteTableRow.ProtoReflect.Descriptor instead.
func (*SubRouteTableRow) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_rawDescGZIP(), []int{0}
}

func (m *SubRouteTableRow) GetAction() isSubRouteTableRow_Action {
	if m != nil {
		return m.Action
	}
	return nil
}

func (x *SubRouteTableRow) GetRouteAction() *v1.RouteAction {
	if x, ok := x.GetAction().(*SubRouteTableRow_RouteAction); ok {
		return x.RouteAction
	}
	return nil
}

func (x *SubRouteTableRow) GetRedirectAction() *v1.RedirectAction {
	if x, ok := x.GetAction().(*SubRouteTableRow_RedirectAction); ok {
		return x.RedirectAction
	}
	return nil
}

func (x *SubRouteTableRow) GetDirectResponseAction() *v1.DirectResponseAction {
	if x, ok := x.GetAction().(*SubRouteTableRow_DirectResponseAction); ok {
		return x.DirectResponseAction
	}
	return nil
}

func (x *SubRouteTableRow) GetDelegateAction() *v11.DelegateAction {
	if x, ok := x.GetAction().(*SubRouteTableRow_DelegateAction); ok {
		return x.DelegateAction
	}
	return nil
}

func (x *SubRouteTableRow) GetMatcher() string {
	if x != nil {
		return x.Matcher
	}
	return ""
}

func (x *SubRouteTableRow) GetMatchType() string {
	if x != nil {
		return x.MatchType
	}
	return ""
}

func (x *SubRouteTableRow) GetMethods() []string {
	if x != nil {
		return x.Methods
	}
	return nil
}

func (x *SubRouteTableRow) GetHeaders() []*matchers.HeaderMatcher {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *SubRouteTableRow) GetQueryParameters() []*matchers.QueryParameterMatcher {
	if x != nil {
		return x.QueryParameters
	}
	return nil
}

func (x *SubRouteTableRow) GetRtRoutes() []*SubRouteTableRow {
	if x != nil {
		return x.RtRoutes
	}
	return nil
}

type isSubRouteTableRow_Action interface {
	isSubRouteTableRow_Action()
}

type SubRouteTableRow_RouteAction struct {
	RouteAction *v1.RouteAction `protobuf:"bytes,1,opt,name=route_action,json=routeAction,proto3,oneof"`
}

type SubRouteTableRow_RedirectAction struct {
	RedirectAction *v1.RedirectAction `protobuf:"bytes,2,opt,name=redirect_action,json=redirectAction,proto3,oneof"`
}

type SubRouteTableRow_DirectResponseAction struct {
	DirectResponseAction *v1.DirectResponseAction `protobuf:"bytes,3,opt,name=direct_response_action,json=directResponseAction,proto3,oneof"`
}

type SubRouteTableRow_DelegateAction struct {
	DelegateAction *v11.DelegateAction `protobuf:"bytes,4,opt,name=delegate_action,json=delegateAction,proto3,oneof"`
}

func (*SubRouteTableRow_RouteAction) isSubRouteTableRow_Action() {}

func (*SubRouteTableRow_RedirectAction) isSubRouteTableRow_Action() {}

func (*SubRouteTableRow_DirectResponseAction) isSubRouteTableRow_Action() {}

func (*SubRouteTableRow_DelegateAction) isSubRouteTableRow_Action() {}

type GetVirtualServiceRoutesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VirtualServiceRef *v12.ClusterObjectRef `protobuf:"bytes,1,opt,name=virtual_service_ref,json=virtualServiceRef,proto3" json:"virtual_service_ref,omitempty"`
}

func (x *GetVirtualServiceRoutesRequest) Reset() {
	*x = GetVirtualServiceRoutesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVirtualServiceRoutesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVirtualServiceRoutesRequest) ProtoMessage() {}

func (x *GetVirtualServiceRoutesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVirtualServiceRoutesRequest.ProtoReflect.Descriptor instead.
func (*GetVirtualServiceRoutesRequest) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_rawDescGZIP(), []int{1}
}

func (x *GetVirtualServiceRoutesRequest) GetVirtualServiceRef() *v12.ClusterObjectRef {
	if x != nil {
		return x.VirtualServiceRef
	}
	return nil
}

type GetVirtualServiceRoutesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VsRoutes []*SubRouteTableRow `protobuf:"bytes,1,rep,name=vs_routes,json=vsRoutes,proto3" json:"vs_routes,omitempty"`
}

func (x *GetVirtualServiceRoutesResponse) Reset() {
	*x = GetVirtualServiceRoutesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVirtualServiceRoutesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVirtualServiceRoutesResponse) ProtoMessage() {}

func (x *GetVirtualServiceRoutesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVirtualServiceRoutesResponse.ProtoReflect.Descriptor instead.
func (*GetVirtualServiceRoutesResponse) Descriptor() ([]byte, []int) {
	return file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_rawDescGZIP(), []int{2}
}

func (x *GetVirtualServiceRoutesResponse) GetVsRoutes() []*SubRouteTableRow {
	if x != nil {
		return x.VsRoutes
	}
	return nil
}

var File_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto protoreflect.FileDescriptor

var file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_rawDesc = []byte{
	0x0a, 0x5a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c,
	0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x73,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x64,
	0x67, 0x65, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x74, 0x5f, 0x73, 0x65,
	0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x15, 0x72, 0x70,
	0x63, 0x2e, 0x65, 0x64, 0x67, 0x65, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f,
	0x2e, 0x69, 0x6f, 0x1a, 0x12, 0x65, 0x78, 0x74, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x4a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x73, 0x6f, 0x6c, 0x6f,
	0x2d, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x67,
	0x6c, 0x6f, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x6d, 0x61, 0x74, 0x63,
	0x68, 0x65, 0x72, 0x73, 0x2f, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x39, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x61, 0x70, 0x69,
	0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f,
	0x76, 0x31, 0x2f, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x46,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d,
	0x69, 0x6f, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x61, 0x70, 0x69, 0x73, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x67, 0x6c, 0x6f, 0x6f, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x2f, 0x76, 0x31,
	0x2f, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x73, 0x6b, 0x76, 0x32, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x89, 0x05, 0x0a, 0x10, 0x53, 0x75, 0x62, 0x52, 0x6f,
	0x75, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x6f, 0x77, 0x12, 0x3e, 0x0a, 0x0c, 0x72,
	0x6f, 0x75, 0x74, 0x65, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f,
	0x2e, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x0b,
	0x72, 0x6f, 0x75, 0x74, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x47, 0x0a, 0x0f, 0x72,
	0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f,
	0x2e, 0x69, 0x6f, 0x2e, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x41, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x48, 0x00, 0x52, 0x0e, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x5a, 0x0a, 0x16, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x72,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f,
	0x2e, 0x69, 0x6f, 0x2e, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x14, 0x64, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x4a, 0x0a, 0x0f, 0x64, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x5f, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x67, 0x61, 0x74, 0x65,
	0x77, 0x61, 0x79, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x67, 0x61, 0x74, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x0e, 0x64, 0x65,
	0x6c, 0x65, 0x67, 0x61, 0x74, 0x65, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x61, 0x74, 0x63,
	0x68, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73,
	0x18, 0x07, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x73, 0x12,
	0x43, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x29, 0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x73, 0x2e, 0x63, 0x6f, 0x72, 0x65,
	0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x48, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x52, 0x07, 0x68, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x73, 0x12, 0x5c, 0x0a, 0x10, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x70, 0x61,
	0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x31,
	0x2e, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x73, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x67,
	0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x51, 0x75, 0x65, 0x72,
	0x79, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65, 0x72, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x65,
	0x72, 0x52, 0x0f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65, 0x74, 0x65,
	0x72, 0x73, 0x12, 0x44, 0x0a, 0x09, 0x72, 0x74, 0x5f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x18,
	0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x64, 0x67, 0x65,
	0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x53, 0x75,
	0x62, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x6f, 0x77, 0x52, 0x08,
	0x72, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x42, 0x08, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x22, 0x75, 0x0a, 0x1e, 0x47, 0x65, 0x74, 0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x53, 0x0a, 0x13, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x5f,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f, 0x72, 0x65, 0x66, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x23, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x73, 0x6b, 0x76, 0x32, 0x2e, 0x73, 0x6f,
	0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x43, 0x6c, 0x75, 0x73, 0x74, 0x65, 0x72, 0x4f, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x52, 0x65, 0x66, 0x52, 0x11, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x66, 0x22, 0x67, 0x0a, 0x1f, 0x47, 0x65, 0x74,
	0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x6f,
	0x75, 0x74, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x44, 0x0a, 0x09,
	0x76, 0x73, 0x5f, 0x72, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x27, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x64, 0x67, 0x65, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e,
	0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x53, 0x75, 0x62, 0x52, 0x6f, 0x75, 0x74, 0x65,
	0x54, 0x61, 0x62, 0x6c, 0x65, 0x52, 0x6f, 0x77, 0x52, 0x08, 0x76, 0x73, 0x52, 0x6f, 0x75, 0x74,
	0x65, 0x73, 0x32, 0xa6, 0x01, 0x0a, 0x17, 0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x41, 0x70, 0x69, 0x12, 0x8a,
	0x01, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x12, 0x35, 0x2e, 0x72, 0x70, 0x63,
	0x2e, 0x65, 0x64, 0x67, 0x65, 0x2e, 0x67, 0x6c, 0x6f, 0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e,
	0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x52, 0x6f, 0x75, 0x74, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x36, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x64, 0x67, 0x65, 0x2e, 0x67, 0x6c, 0x6f,
	0x6f, 0x2e, 0x73, 0x6f, 0x6c, 0x6f, 0x2e, 0x69, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x69, 0x72,
	0x74, 0x75, 0x61, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x6f, 0x75, 0x74, 0x65,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x56, 0xb8, 0xf5, 0x04,
	0x01, 0xc0, 0xf5, 0x04, 0x01, 0x5a, 0x4c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x69, 0x6f, 0x2f, 0x73, 0x6f, 0x6c, 0x6f, 0x2d, 0x70,
	0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x73,
	0x2f, 0x61, 0x70, 0x69, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x72, 0x70, 0x63, 0x2e, 0x65, 0x64, 0x67, 0x65, 0x2e, 0x67, 0x6c, 0x6f, 0x6f,
	0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_rawDescOnce sync.Once
	file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_rawDescData = file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_rawDesc
)

func file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_rawDescGZIP() []byte {
	file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_rawDescOnce.Do(func() {
		file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_rawDescData)
	})
	return file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_rawDescData
}

var file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_goTypes = []interface{}{
	(*SubRouteTableRow)(nil),                // 0: rpc.edge.gloo.solo.io.SubRouteTableRow
	(*GetVirtualServiceRoutesRequest)(nil),  // 1: rpc.edge.gloo.solo.io.GetVirtualServiceRoutesRequest
	(*GetVirtualServiceRoutesResponse)(nil), // 2: rpc.edge.gloo.solo.io.GetVirtualServiceRoutesResponse
	(*v1.RouteAction)(nil),                  // 3: gloo.solo.io.RouteAction
	(*v1.RedirectAction)(nil),               // 4: gloo.solo.io.RedirectAction
	(*v1.DirectResponseAction)(nil),         // 5: gloo.solo.io.DirectResponseAction
	(*v11.DelegateAction)(nil),              // 6: gateway.solo.io.DelegateAction
	(*matchers.HeaderMatcher)(nil),          // 7: matchers.core.gloo.solo.io.HeaderMatcher
	(*matchers.QueryParameterMatcher)(nil),  // 8: matchers.core.gloo.solo.io.QueryParameterMatcher
	(*v12.ClusterObjectRef)(nil),            // 9: core.skv2.solo.io.ClusterObjectRef
}
var file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_depIdxs = []int32{
	3,  // 0: rpc.edge.gloo.solo.io.SubRouteTableRow.route_action:type_name -> gloo.solo.io.RouteAction
	4,  // 1: rpc.edge.gloo.solo.io.SubRouteTableRow.redirect_action:type_name -> gloo.solo.io.RedirectAction
	5,  // 2: rpc.edge.gloo.solo.io.SubRouteTableRow.direct_response_action:type_name -> gloo.solo.io.DirectResponseAction
	6,  // 3: rpc.edge.gloo.solo.io.SubRouteTableRow.delegate_action:type_name -> gateway.solo.io.DelegateAction
	7,  // 4: rpc.edge.gloo.solo.io.SubRouteTableRow.headers:type_name -> matchers.core.gloo.solo.io.HeaderMatcher
	8,  // 5: rpc.edge.gloo.solo.io.SubRouteTableRow.query_parameters:type_name -> matchers.core.gloo.solo.io.QueryParameterMatcher
	0,  // 6: rpc.edge.gloo.solo.io.SubRouteTableRow.rt_routes:type_name -> rpc.edge.gloo.solo.io.SubRouteTableRow
	9,  // 7: rpc.edge.gloo.solo.io.GetVirtualServiceRoutesRequest.virtual_service_ref:type_name -> core.skv2.solo.io.ClusterObjectRef
	0,  // 8: rpc.edge.gloo.solo.io.GetVirtualServiceRoutesResponse.vs_routes:type_name -> rpc.edge.gloo.solo.io.SubRouteTableRow
	1,  // 9: rpc.edge.gloo.solo.io.VirtualServiceRoutesApi.GetVirtualServiceRoutes:input_type -> rpc.edge.gloo.solo.io.GetVirtualServiceRoutesRequest
	2,  // 10: rpc.edge.gloo.solo.io.VirtualServiceRoutesApi.GetVirtualServiceRoutes:output_type -> rpc.edge.gloo.solo.io.GetVirtualServiceRoutesResponse
	10, // [10:11] is the sub-list for method output_type
	9,  // [9:10] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() {
	file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_init()
}
func file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_init() {
	if File_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubRouteTableRow); i {
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
		file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVirtualServiceRoutesRequest); i {
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
		file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVirtualServiceRoutesResponse); i {
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
	file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*SubRouteTableRow_RouteAction)(nil),
		(*SubRouteTableRow_RedirectAction)(nil),
		(*SubRouteTableRow_DirectResponseAction)(nil),
		(*SubRouteTableRow_DelegateAction)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_goTypes,
		DependencyIndexes: file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_depIdxs,
		MessageInfos:      file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_msgTypes,
	}.Build()
	File_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto = out.File
	file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_rawDesc = nil
	file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_goTypes = nil
	file_github_com_solo_io_solo_projects_projects_apiserver_api_rpc_edge_gloo_v1_rt_selector_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// VirtualServiceRoutesApiClient is the client API for VirtualServiceRoutesApi service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type VirtualServiceRoutesApiClient interface {
	GetVirtualServiceRoutes(ctx context.Context, in *GetVirtualServiceRoutesRequest, opts ...grpc.CallOption) (*GetVirtualServiceRoutesResponse, error)
}

type virtualServiceRoutesApiClient struct {
	cc grpc.ClientConnInterface
}

func NewVirtualServiceRoutesApiClient(cc grpc.ClientConnInterface) VirtualServiceRoutesApiClient {
	return &virtualServiceRoutesApiClient{cc}
}

func (c *virtualServiceRoutesApiClient) GetVirtualServiceRoutes(ctx context.Context, in *GetVirtualServiceRoutesRequest, opts ...grpc.CallOption) (*GetVirtualServiceRoutesResponse, error) {
	out := new(GetVirtualServiceRoutesResponse)
	err := c.cc.Invoke(ctx, "/rpc.edge.gloo.solo.io.VirtualServiceRoutesApi/GetVirtualServiceRoutes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VirtualServiceRoutesApiServer is the server API for VirtualServiceRoutesApi service.
type VirtualServiceRoutesApiServer interface {
	GetVirtualServiceRoutes(context.Context, *GetVirtualServiceRoutesRequest) (*GetVirtualServiceRoutesResponse, error)
}

// UnimplementedVirtualServiceRoutesApiServer can be embedded to have forward compatible implementations.
type UnimplementedVirtualServiceRoutesApiServer struct {
}

func (*UnimplementedVirtualServiceRoutesApiServer) GetVirtualServiceRoutes(context.Context, *GetVirtualServiceRoutesRequest) (*GetVirtualServiceRoutesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetVirtualServiceRoutes not implemented")
}

func RegisterVirtualServiceRoutesApiServer(s *grpc.Server, srv VirtualServiceRoutesApiServer) {
	s.RegisterService(&_VirtualServiceRoutesApi_serviceDesc, srv)
}

func _VirtualServiceRoutesApi_GetVirtualServiceRoutes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetVirtualServiceRoutesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VirtualServiceRoutesApiServer).GetVirtualServiceRoutes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpc.edge.gloo.solo.io.VirtualServiceRoutesApi/GetVirtualServiceRoutes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VirtualServiceRoutesApiServer).GetVirtualServiceRoutes(ctx, req.(*GetVirtualServiceRoutesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _VirtualServiceRoutesApi_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpc.edge.gloo.solo.io.VirtualServiceRoutesApi",
	HandlerType: (*VirtualServiceRoutesApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetVirtualServiceRoutes",
			Handler:    _VirtualServiceRoutesApi_GetVirtualServiceRoutes_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/solo-io/solo-projects/projects/apiserver/api/rpc.edge.gloo/v1/rt_selector.proto",
}
