package lambda

import (
	"context"
	"fmt"
	"net/http"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kgateway-dev/kgateway/v2/pkg/utils/kubeutils"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/requestutils/curl"
	testmatchers "github.com/kgateway-dev/kgateway/v2/test/gomega/matchers"
	"github.com/kgateway-dev/kgateway/v2/test/kubernetes/e2e"
	testdefaults "github.com/kgateway-dev/kgateway/v2/test/kubernetes/e2e/defaults"
	"github.com/kgateway-dev/kgateway/v2/test/testutils"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

// testingSuite is a suite of Lambda backend routing tests
type testingSuite struct {
	suite.Suite
	ctx context.Context
	ti  *e2e.TestInstallation
}

func NewTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &testingSuite{
		ctx: ctx,
		ti:  testInst,
	}
}

func (s *testingSuite) SetupSuite() {
	err := s.ti.Actions.Kubectl().ApplyFile(s.ctx, setupManifest)
	s.NoError(err, "can apply "+setupManifest)
	err = s.ti.Actions.Kubectl().ApplyFile(s.ctx, testdefaults.CurlPodManifest)
	s.NoError(err, "can apply curl pod manifest")

	s.ti.Assertions.EventuallyObjectsExist(s.ctx, testdefaults.CurlPod)
	s.ti.Assertions.EventuallyPodsRunning(s.ctx, testdefaults.CurlPod.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=curl",
	})
}

func (s *testingSuite) TearDownSuite() {
	if testutils.ShouldTearDown() {
		fmt.Println("TEAR_DOWN is set, skipping cleanup")
		return
	}
	err := s.ti.Actions.Kubectl().DeleteFileSafe(s.ctx, setupManifest)
	s.NoError(err, "can delete setup manifest")
	err = s.ti.Actions.Kubectl().DeleteFileSafe(s.ctx, testdefaults.CurlPodManifest)
	s.NoError(err, "can delete curl pod manifest")
}

func (s *testingSuite) TestLambdaBackendRouting() {
	// create backend lambda
	err := s.ti.Actions.Kubectl().ApplyFile(s.ctx, lambdaBackendManifest)
	s.NoError(err, "can apply lambda backend manifest")
	s.ti.Assertions.EventuallyObjectsExist(s.ctx, lambdaBackend)

	// Test Lambda backend with custom endpoint
	s.ti.Assertions.AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(gatewayObjectMeta)),
			curl.WithHostHeader("www.example.com"),
			curl.WithPort(8080),
			curl.WithPath("/lambda"),
		},
		&testmatchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body:       gomega.ContainSubstring(`Hello from Lambda`),
		},
	)
}

func (s *testingSuite) TestLambdaBackendAsyncRouting() {
	// create backend lambda
	err := s.ti.Actions.Kubectl().ApplyFile(s.ctx, lambdaAsyncManifest)
	s.NoError(err, "can apply lambda backend manifest")
	s.ti.Assertions.EventuallyObjectsExist(s.ctx, lambdaBackend)

	// Test Lambda backend with custom endpoint
	s.ti.Assertions.AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(gatewayObjectMeta)),
			curl.WithHostHeader("www.example.com"),
			curl.WithPort(8080),
			curl.WithPath("/lambda"),
		},
		&testmatchers.HttpResponse{
			StatusCode: http.StatusAccepted,
			Body:       gomega.BeEmpty(),
		},
	)
}
