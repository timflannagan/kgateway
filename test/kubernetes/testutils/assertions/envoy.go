//go:build ignore

package assertions

import (
	"context"
	"io"
	"net"
	"time"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kgateway-dev/kgateway/v2/internal/gloo/pkg/defaults"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/envoyutils/admincli"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils/portforward"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/requestutils/curl"
)

func (p *Provider) AssertEnvoyAdminApi(
	ctx context.Context,
	envoyDeployment metav1.ObjectMeta,
	adminAssertions ...func(ctx context.Context, adminClient *admincli.Client),
) {
	// Before opening a port-forward, we assert that there is at least one Pod that is ready
	p.EventuallyReadyReplicas(ctx, envoyDeployment, BeNumerically(">=", 1))

	portForwarder, err := p.clusterContext.Cli.StartPortForward(ctx,
		portforward.WithDeployment(envoyDeployment.GetName(), envoyDeployment.GetNamespace()),
		portforward.WithPorts(int(defaults.EnvoyAdminPort), int(defaults.EnvoyAdminPort)),
	)
	p.Require.NoError(err, "can open port-forward")
	defer func() {
		portForwarder.Close()
		portForwarder.WaitForStop()
	}()

	// the port-forward returns before it completely starts up (https://github.com/kgateway-dev/kgateway/issues/9353),
	// so as a workaround we try to keep dialing the address until it succeeds
	p.Gomega.Eventually(func(g Gomega) {
		_, err = net.Dial("tcp", portForwarder.Address())
		g.Expect(err).NotTo(HaveOccurred())
	}).
		WithContext(ctx).
		WithTimeout(time.Second * 15).
		WithPolling(time.Second).
		Should(Succeed())

	adminClient := admincli.NewClient().
		WithReceiver(io.Discard). // adminAssertion can overwrite this
		WithCurlOptions(
			curl.WithRetries(3, 0, 10),
			curl.WithPort(int(defaults.EnvoyAdminPort)),
		)

	for _, adminAssertion := range adminAssertions {
		adminAssertion(ctx, adminClient)
	}
}
