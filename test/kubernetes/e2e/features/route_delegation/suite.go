package route_delegation

import (
	"context"
	"net/http"
	"time"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/pkg/utils/requestutils/curl"
	testmatchers "github.com/kgateway-dev/kgateway/v2/test/gomega/matchers"
	"github.com/kgateway-dev/kgateway/v2/test/kubernetes/e2e"
	"github.com/kgateway-dev/kgateway/v2/test/kubernetes/e2e/defaults"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

type tsuite struct {
	suite.Suite

	ctx context.Context
	// ti contains all the metadata/utilities necessary to execute a series of tests
	// against an installation of Gloo Gateway
	ti *e2e.TestInstallation

	// maps test name to a list of manifests to apply before the test
	manifests map[string][]string

	manifestObjects map[string][]client.Object
}

func NewTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &tsuite{
		ctx: ctx,
		ti:  testInst,
	}
}

func (s *tsuite) SetupSuite() {
	s.manifests = map[string][]string{
		"TestBasic":                       {commonManifest, basicRoutesManifest},
		"TestRecursive":                   {commonManifest, recursiveRoutesManifest},
		"TestCyclic":                      {commonManifest, cyclicRoutesManifest},
		"TestInvalidChild":                {commonManifest, invalidChildRoutesManifest},
		"TestHeaderQueryMatch":            {commonManifest, headerQueryMatchRoutesManifest},
		"TestMultipleParents":             {commonManifest, multipleParentsManifest},
		"TestInvalidChildValidStandalone": {commonManifest, invalidChildValidStandaloneManifest},
		"TestUnresolvedChild":             {commonManifest, unresolvedChildManifest},
		"TestMatcherInheritance":          {commonManifest, matcherInheritanceManifest},
	}
	// Not every resource that is applied needs to be verified. We are not testing `kubectl apply`,
	// but the below code demonstrates how it can be done if necessary
	s.manifestObjects = map[string][]client.Object{
		commonManifest:                      {proxyService, proxyDeployment, defaults.CurlPod, httpbinTeam1, httpbinTeam2, gateway},
		basicRoutesManifest:                 {routeRoot, routeTeam1, routeTeam2},
		cyclicRoutesManifest:                {routeRoot, routeTeam1, routeTeam2},
		recursiveRoutesManifest:             {routeRoot, routeTeam1, routeTeam2},
		invalidChildRoutesManifest:          {routeRoot, routeTeam1, routeTeam2},
		headerQueryMatchRoutesManifest:      {routeRoot, routeTeam1, routeTeam2},
		multipleParentsManifest:             {routeParent1, routeParent2, routeTeam1, routeTeam2},
		invalidChildValidStandaloneManifest: {proxyTestService, proxyTestDeployment, routeRoot, routeTeam1, routeTeam2},
		unresolvedChildManifest:             {routeRoot},
		matcherInheritanceManifest:          {routeParent1, routeParent2, routeTeam1},
	}
}

func (s *tsuite) TearDownSuite() {
	// nothing at the moment
}

func (s *tsuite) BeforeTest(suiteName, testName string) {
	manifests := s.manifests[testName]
	for _, manifest := range manifests {
		err := s.ti.Actions.Kubectl().ApplyFile(s.ctx, manifest)
		s.Require().NoError(err)
		s.ti.Assertions.EventuallyObjectsExist(s.ctx, s.manifestObjects[manifest]...)
	}
}

func (s *tsuite) AfterTest(suiteName, testName string) {
	manifests := s.manifests[testName]
	for _, manifest := range manifests {
		err := s.ti.Actions.Kubectl().DeleteFileSafe(s.ctx, manifest)
		s.Require().NoError(err)
		s.ti.Assertions.EventuallyObjectsNotExist(s.ctx, s.manifestObjects[manifest]...)
	}
}

func (s *tsuite) TestBasic() {
	// Assert traffic to team1 route
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam1)},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to team2 route
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam2)},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam2)})
}

func (s *tsuite) TestRecursive() {
	// Assert traffic to team1 route
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam1)},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to team2 route
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam2)},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam2)})
}

func (s *tsuite) TestCyclic() {
	// Assert traffic to team1 route
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam1)},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to team2 route fails with HTTP 404 as it is a cyclic route
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam2)},
		&testmatchers.HttpResponse{StatusCode: http.StatusNotFound})

	s.ti.Assertions.EventuallyHTTPRouteStatusContainsMessage(s.ctx, routeTeam2.Name, routeTeam2.Namespace,
		"cyclic reference detected", 10*time.Second, 1*time.Second)
}

func (s *tsuite) TestInvalidChild() {
	// Assert traffic to team1 route
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam1)},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to team2 route fails with HTTP 404 as the route is invalid due to specifying a hostname on the child route
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam2)},
		&testmatchers.HttpResponse{StatusCode: http.StatusNotFound})

	s.ti.Assertions.EventuallyHTTPRouteStatusContainsMessage(s.ctx, routeTeam2.Name, routeTeam2.Namespace,
		"spec.hostnames must be unset", 10*time.Second, 1*time.Second)
}

func (s *tsuite) TestHeaderQueryMatch() {
	// Assert traffic to team1 route with matching header and query parameters
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam1),
			curl.WithHeader("header1", "val1"),
			curl.WithHeader("headerX", "valX"),
			curl.WithQueryParameters(map[string]string{"query1": "val1", "queryX": "valX"}),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to team2 route fails with HTTP 404 as it does not match the parent's header and query parameters
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam2),
			curl.WithHeader("headerX", "valX"),
			curl.WithQueryParameters(map[string]string{"queryX": "valX"}),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusNotFound})
}

func (s *tsuite) TestMultipleParents() {
	// Assert traffic to parent1.com/anything/team1
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyHostPort),
			curl.WithPath(pathTeam1),
			curl.WithHostHeader(routeParent1Host),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to parent1.com/anything/team2
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyHostPort),
			curl.WithPath(pathTeam2),
			curl.WithHostHeader(routeParent1Host),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam2)})

	// Assert traffic to parent2.com/anything/team1
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyHostPort),
			curl.WithPath(pathTeam1),
			curl.WithHostHeader(routeParent2Host),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to parent2.com/anything/team2 fails as it is not selected by parent2 route
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyHostPort),
			curl.WithPath(pathTeam2),
			curl.WithHostHeader(routeParent2Host),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusNotFound})
}

func (s *tsuite) TestInvalidChildValidStandalone() {
	// Assert traffic to team1 route
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyTestHostPort),
			curl.WithPath(pathTeam1),
			curl.WithHostHeader(routeParentHost),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to team2 route on parent hostname fails with HTTP 404 as the route is invalid due to specifying a hostname on the child route
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyTestHostPort),
			curl.WithPath(pathTeam2),
			curl.WithHostHeader(routeParentHost),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusNotFound})

	// Assert traffic to team2 route on standalone host succeeds
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyTestHostPort),
			curl.WithPath(pathTeam2),
			curl.WithHostHeader(routeTeam2Host),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam2)})

	s.ti.Assertions.EventuallyHTTPRouteStatusContainsMessage(s.ctx, routeTeam2.Name, routeTeam2.Namespace,
		"spec.hostnames must be unset", 10*time.Second, 1*time.Second)
}

func (s *tsuite) TestUnresolvedChild() {
	s.ti.Assertions.EventuallyHTTPRouteStatusContainsReason(s.ctx, routeRoot.Name, routeRoot.Namespace,
		string(gwv1.RouteReasonBackendNotFound), 10*time.Second, 1*time.Second)
}

func (s *tsuite) TestMatcherInheritance() {
	// Assert traffic on parent1's prefix
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath("/anything/foo/child")},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring("/anything/foo/child")})

	// Assert traffic on parent2's prefix
	s.ti.Assertions.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath("/anything/baz/child")},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring("/anything/baz/child")})
}
