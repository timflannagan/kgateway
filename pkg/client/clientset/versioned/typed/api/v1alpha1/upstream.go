// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	context "context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"

	applyconfigurationapiv1alpha1 "github.com/kgateway-dev/kgateway/v2/api/applyconfiguration/api/v1alpha1"
	apiv1alpha1 "github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	scheme "github.com/kgateway-dev/kgateway/v2/pkg/client/clientset/versioned/scheme"
)

// UpstreamsGetter has a method to return a UpstreamInterface.
// A group's client should implement this interface.
type UpstreamsGetter interface {
	Upstreams(namespace string) UpstreamInterface
}

// UpstreamInterface has methods to work with Upstream resources.
type UpstreamInterface interface {
	Create(ctx context.Context, upstream *apiv1alpha1.Upstream, opts v1.CreateOptions) (*apiv1alpha1.Upstream, error)
	Update(ctx context.Context, upstream *apiv1alpha1.Upstream, opts v1.UpdateOptions) (*apiv1alpha1.Upstream, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, upstream *apiv1alpha1.Upstream, opts v1.UpdateOptions) (*apiv1alpha1.Upstream, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*apiv1alpha1.Upstream, error)
	List(ctx context.Context, opts v1.ListOptions) (*apiv1alpha1.UpstreamList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *apiv1alpha1.Upstream, err error)
	Apply(ctx context.Context, upstream *applyconfigurationapiv1alpha1.UpstreamApplyConfiguration, opts v1.ApplyOptions) (result *apiv1alpha1.Upstream, err error)
	// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
	ApplyStatus(ctx context.Context, upstream *applyconfigurationapiv1alpha1.UpstreamApplyConfiguration, opts v1.ApplyOptions) (result *apiv1alpha1.Upstream, err error)
	UpstreamExpansion
}

// upstreams implements UpstreamInterface
type upstreams struct {
	*gentype.ClientWithListAndApply[*apiv1alpha1.Upstream, *apiv1alpha1.UpstreamList, *applyconfigurationapiv1alpha1.UpstreamApplyConfiguration]
}

// newUpstreams returns a Upstreams
func newUpstreams(c *GatewayV1alpha1Client, namespace string) *upstreams {
	return &upstreams{
		gentype.NewClientWithListAndApply[*apiv1alpha1.Upstream, *apiv1alpha1.UpstreamList, *applyconfigurationapiv1alpha1.UpstreamApplyConfiguration](
			"upstreams",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *apiv1alpha1.Upstream { return &apiv1alpha1.Upstream{} },
			func() *apiv1alpha1.UpstreamList { return &apiv1alpha1.UpstreamList{} },
		),
	}
}
