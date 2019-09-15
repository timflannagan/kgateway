package virtualservice_test

import (
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	glooCmdutils "github.com/solo-io/gloo/projects/gloo/cli/pkg/cmdutils"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	ratelimitpb "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/ratelimit"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/testutils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/ratelimit"
)

var _ = Describe("CustomEnvoyConfig", func() {

	var (
		vsvc     *gatewayv1.VirtualService
		vsClient gatewayv1.VirtualServiceClient
	)
	BeforeEach(func() {
		helpers.UseMemoryClients()
		// create a settings object
		vsClient = helpers.MustVirtualServiceClient()
		vsvc = &gatewayv1.VirtualService{
			Metadata: core.Metadata{
				Name:      "vs",
				Namespace: "gloo-system",
			},
			VirtualHost: &gatewayv1.VirtualHost{
				Routes: []*gatewayv1.Route{{
					Matcher: &v1.Matcher{
						PathSpecifier: &v1.Matcher_Prefix{Prefix: "/"},
					}}},
			},
		}

		var err error
		vsvc, err = vsClient.Write(vsvc, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
	})

	rateLimitExtension := func() *ratelimitpb.RateLimitVhostExtension {
		var routeExt ratelimitpb.RateLimitVhostExtension
		var err error
		vsvc, err = vsClient.Read(vsvc.Metadata.Namespace, vsvc.Metadata.Name, clients.ReadOpts{})
		Expect(err).NotTo(HaveOccurred())

		err = utils.UnmarshalExtension(vsvc.VirtualHost.VirtualHostPlugins, ratelimit.EnvoyExtensionName, &routeExt)
		Expect(err).NotTo(HaveOccurred())

		return &routeExt
	}

	It("should edit virtual service", func() {

		glooCmdutils.EditFileForTest = func(prefix, suffix string, r io.Reader) ([]byte, string, error) {
			b := `
rate_limits:
- actions:
  - source_cluster: {}`
			return []byte(b), "", nil
		}

		err := testutils.GlooctlEE("edit virtualservice --name vs --namespace gloo-system ratelimit custom-envoy-config")
		Expect(err).NotTo(HaveOccurred())

		ext := rateLimitExtension()
		Expect(ext.RateLimits).To(HaveLen(1))
		Expect(ext.RateLimits[0].Actions).To(HaveLen(1))
	})

})
