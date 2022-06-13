package syncer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/aws"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"

	"github.com/solo-io/solo-projects/projects/observability/pkg/grafana/template"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	v1 "github.com/solo-io/solo-projects/projects/observability/pkg/api/v1"
	"github.com/solo-io/solo-projects/projects/observability/pkg/grafana"
	"github.com/solo-io/solo-projects/projects/observability/pkg/grafana/mocks"
	"go.uber.org/zap"
)

var (
	mockCtrl                *gomock.Controller
	dashboardClient         *mocks.MockDashboardClient
	snapshotClient          *mocks.MockSnapshotClient
	testErr                 = errors.New("test err")
	dashboardsSnapshot      *v1.DashboardsSnapshot
	testGrafanaState        *grafanaState
	dashboardSyncer         *GrafanaDashboardsSyncer
	logger                  *zap.SugaredLogger
	upstreamOne             *gloov1.Upstream
	upstreamTwo             *gloov1.Upstream
	upstreamList            gloov1.UpstreamList
	snapshotResponse        *grafana.SnapshotListResponse
	dashboardSearchResponse []grafana.FoundBoard
	dashboardJsonTemplate   string
	defaultDashboardUids    map[string]struct{}
)

var _ = Describe("Grafana Syncer", func() {
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		dashboardClient = mocks.NewMockDashboardClient(mockCtrl)
		snapshotClient = mocks.NewMockSnapshotClient(mockCtrl)
		defaultDashboardUids = make(map[string]struct{})
		dashboardsSnapshot = &v1.DashboardsSnapshot{
			Upstreams: []*gloov1.Upstream{
				{
					Metadata: &core.Metadata{Name: "us1", Namespace: "ns"},
				},
				{
					Metadata: &core.Metadata{Name: "us2", Namespace: "ns"},
				},
			},
		}
		testGrafanaState = &grafanaState{
			boards:    nil,
			snapshots: nil,
		}
		dashboardJsonTemplate = "{\"test-content\": \"content\"}"
		dashboardSyncer = NewGrafanaDashboardSyncer(dashboardClient, snapshotClient, dashboardJsonTemplate, generalFolderId, defaultDashboardUids)
		logger = contextutils.LoggerFrom(context.TODO())
		upstreamOne = &gloov1.Upstream{
			UpstreamType: &gloov1.Upstream_Aws{
				Aws: &aws.UpstreamSpec{Region: "test"},
			},
			NamespacedStatuses: &core.NamespacedStatuses{},
			Metadata:           &core.Metadata{Name: "test-upstream-one", Namespace: "ns"},
		}
		upstreamTwo = &gloov1.Upstream{
			UpstreamType: &gloov1.Upstream_Kube{
				Kube: &kubernetes.UpstreamSpec{
					ServiceName:      "upstream-service",
					ServiceNamespace: "ns",
					ServicePort:      80,
				},
			},
			NamespacedStatuses: &core.NamespacedStatuses{},
			Metadata:           &core.Metadata{Name: "test-upstream-two", Namespace: "ns"},
		}
		upstreamList = []*gloov1.Upstream{upstreamOne, upstreamTwo}

		snapshotResponse = &grafana.SnapshotListResponse{}
		dashboardSearchResponse = []grafana.FoundBoard{}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Dashboard Syncer", func() {
		var addIdToDashboardJson = func(bytes []byte, id int) []byte {
			jsonDashboard := map[string]interface{}{}
			err := json.Unmarshal(bytes, &jsonDashboard)
			Expect(err).NotTo(HaveOccurred())
			jsonDashboard["id"] = id
			jsonWithId, err := json.Marshal(jsonDashboard)
			Expect(err).NotTo(HaveOccurred())

			return jsonWithId
		}

		It("builds dashboards from scratch", func() {
			for _, upstream := range upstreamList {
				templateGenerator := template.NewUpstreamTemplateGenerator(upstream, dashboardJsonTemplate)
				uid := templateGenerator.GenerateUid()

				dashboardClient.EXPECT().
					GetRawDashboard(uid).
					Return(nil, grafana.BoardProperties{}, grafana.DashboardNotFound(uid))

				dashPost, err := templateGenerator.GenerateDashboardPost(generalFolderId) // should just return "test-json"
				Expect(err).NotTo(HaveOccurred())
				snapshotBytes, err := templateGenerator.GenerateSnapshot() // should just return "test-json"
				Expect(err).NotTo(HaveOccurred())

				dashboardClient.EXPECT().
					PostDashboard(dashPost).
					Return(nil)
				snapshotClient.EXPECT().
					SetRawSnapshot(snapshotBytes).
					Return(nil, nil) // we don't consume the snapshot response apart from checking the error
			}

			snapshotClient.EXPECT().
				GetSnapshots().
				Return([]grafana.SnapshotListResponse{}, nil)
			dashboardClient.EXPECT().
				SearchDashboards("", false, tags[0], tags[1]).
				Return([]grafana.FoundBoard{}, nil)

			err := dashboardSyncer.Sync(context.TODO(), &v1.DashboardsSnapshot{
				Upstreams: upstreamList, // we just init'd two new upstreams
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("rebuilds dashboards when they exist already", func() {
			for index, upstream := range upstreamList {
				templateGenerator := template.NewUpstreamTemplateGenerator(upstream, dashboardJsonTemplate)
				renderedDashboardPost, err := templateGenerator.GenerateDashboardPost(generalFolderId)
				Expect(err).NotTo(HaveOccurred())

				uid := templateGenerator.GenerateUid()

				// have to set the key "id" in the json that gets returned by GetRawDashboard
				bytes, err := json.Marshal(renderedDashboardPost)
				Expect(err).NotTo(HaveOccurred())
				jsonWithId := addIdToDashboardJson(bytes, index)

				dashboardClient.EXPECT().
					GetRawDashboard(uid).
					Return(jsonWithId, grafana.BoardProperties{}, nil)
				dashboardClient.EXPECT().
					GetDashboardVersions(float64(index)).
					Return([]*grafana.Version{
						{
							VersionId: 1,
							Message:   template.DefaultCommitMessage,
						},
					}, nil)

				snapshotBytes, err := templateGenerator.GenerateSnapshot()
				Expect(err).NotTo(HaveOccurred())

				dashboardClient.EXPECT().
					PostDashboard(renderedDashboardPost).
					Return(nil)
				snapshotClient.EXPECT().
					SetRawSnapshot(snapshotBytes).
					Return(nil, nil) // we don't consume the snapshot response apart from checking the error
			}

			snapshotClient.EXPECT().
				GetSnapshots().
				Return([]grafana.SnapshotListResponse{}, nil)
			dashboardClient.EXPECT().
				SearchDashboards("", false, tags[0], tags[1]).
				Return([]grafana.FoundBoard{}, nil)

			err := dashboardSyncer.Sync(context.TODO(), &v1.DashboardsSnapshot{
				Upstreams: upstreamList, // we just init'd two new upstreams
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not regenerate the dashboard when it has been edited by a human", func() {
			for index, upstream := range upstreamList {
				templateGenerator := template.NewUpstreamTemplateGenerator(upstream, dashboardJsonTemplate)
				renderedDashboardPost, err := templateGenerator.GenerateDashboardPost(generalFolderId)
				Expect(err).NotTo(HaveOccurred())

				uid := templateGenerator.GenerateUid()

				// have to set the key "id" in the json that gets returned by GetRawDashboard
				bytes, err := json.Marshal(renderedDashboardPost)
				Expect(err).NotTo(HaveOccurred())
				jsonWithId := addIdToDashboardJson(bytes, index)

				dashboardClient.EXPECT().
					GetRawDashboard(uid).
					Return(jsonWithId, grafana.BoardProperties{}, nil)
				dashboardClient.EXPECT().
					GetDashboardVersions(float64(index)).
					Return([]*grafana.Version{
						{
							VersionId: 1,
							Message:   "This is a message written by a human",
						},
					}, nil)
			}

			snapshotClient.EXPECT().
				GetSnapshots().
				Return([]grafana.SnapshotListResponse{}, nil)
			dashboardClient.EXPECT().
				SearchDashboards("", false, tags[0], tags[1]).
				Return([]grafana.FoundBoard{}, nil)

			err := dashboardSyncer.Sync(context.TODO(), &v1.DashboardsSnapshot{
				Upstreams: upstreamList, // we just init'd two new upstreams
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("deletes the dashboard for a deleted upstream", func() {
			var snapshots []grafana.SnapshotListResponse
			var dashboards []grafana.FoundBoard
			for _, upstream := range upstreamList {
				generator := template.NewUpstreamTemplateGenerator(upstream, dashboardJsonTemplate)
				uid := generator.GenerateUid()
				snapshots = append(snapshots, grafana.SnapshotListResponse{
					Key: uid,
				})
				dashboards = append(dashboards, grafana.FoundBoard{
					UID: uid,
				})

				dashboardClient.EXPECT().
					DeleteDashboard(uid).
					Return(grafana.StatusMessage{}, nil)
				snapshotClient.EXPECT().
					DeleteSnapshot(uid).
					Return(nil)
			}

			snapshotClient.EXPECT().
				GetSnapshots().
				Return(snapshots, nil)
			dashboardClient.EXPECT().
				SearchDashboards("", false, tags[0], tags[1]).
				Return(dashboards, nil)

			err := dashboardSyncer.Sync(context.TODO(), &v1.DashboardsSnapshot{
				Upstreams: []*gloov1.Upstream{}, // deleted both the upstreams
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	It("Ignores invalid default folder ids to use the general id", func() {

		invalidFolderId := uint(1)
		dashboardSyncer = NewGrafanaDashboardSyncer(dashboardClient, snapshotClient, dashboardJsonTemplate, invalidFolderId, defaultDashboardUids)
		for _, upstream := range upstreamList {
			templateGenerator := template.NewUpstreamTemplateGenerator(upstream, dashboardJsonTemplate)
			uid := templateGenerator.GenerateUid()

			dashboardClient.EXPECT().
				GetRawDashboard(uid).
				Return(nil, grafana.BoardProperties{}, grafana.DashboardNotFound(uid))

			dashPost, err := templateGenerator.GenerateDashboardPost(generalFolderId) // only 2 upstreams, and 0 and 1 are valid
			Expect(err).NotTo(HaveOccurred())
			snapshotBytes, err := templateGenerator.GenerateSnapshot() // should just return "test-json"
			Expect(err).NotTo(HaveOccurred())

			dashboardClient.EXPECT().
				PostDashboard(dashPost).
				Return(nil)
			snapshotClient.EXPECT().
				SetRawSnapshot(snapshotBytes).
				Return(nil, nil) // we don't consume the snapshot response apart from checking the error
		}
		snapshotClient.EXPECT().
			GetSnapshots().
			Return([]grafana.SnapshotListResponse{}, nil)

		// unlike other functions, this is only called once because the default id gets corrected after the first loop
		// which prompts the second loop to not bother retrieving a valid id list until it encounters an annotated folder
		// id, which isn't present in this test
		dashboardClient.EXPECT().
			GetAllFolderIds().
			Return(nil, nil)

		dashboardClient.EXPECT().
			SearchDashboards("", false, tags[0], tags[1]).
			Return([]grafana.FoundBoard{}, nil)

		err := dashboardSyncer.Sync(context.TODO(), &v1.DashboardsSnapshot{
			Upstreams: upstreamList, // we just init'd two new upstreams
		})
		Expect(err).NotTo(HaveOccurred())
	})

	It("Ignores invalid annotated folder ids and places dashboards in the general folder", func() {
		// Although errors aren't propagated out of the sync loop, we can still tell that it noticed the invalid id error
		// by the fact that it loops twice and calls most functions twice as a result.
		upstreamOne.Metadata.Annotations = map[string]string{upstreamFolderIdAnnotationKey: "1"}
		upstreamTwo.Metadata.Annotations = map[string]string{upstreamFolderIdAnnotationKey: "test"}
		dashboardSyncer = NewGrafanaDashboardSyncer(dashboardClient, snapshotClient, dashboardJsonTemplate, generalFolderId, defaultDashboardUids)

		for _, upstream := range upstreamList {
			templateGenerator := template.NewUpstreamTemplateGenerator(upstream, dashboardJsonTemplate)
			uid := templateGenerator.GenerateUid()

			dashboardClient.EXPECT().
				GetRawDashboard(uid).
				Return(nil, grafana.BoardProperties{}, grafana.DashboardNotFound(uid))

			dashPost, err := templateGenerator.GenerateDashboardPost(generalFolderId) // only 2 upstreams, and 0 and 1 are valid
			Expect(err).NotTo(HaveOccurred())
			snapshotBytes, err := templateGenerator.GenerateSnapshot() // should just return "test-json"
			Expect(err).NotTo(HaveOccurred())

			dashboardClient.EXPECT().
				PostDashboard(dashPost).
				Return(nil)
			snapshotClient.EXPECT().
				SetRawSnapshot(snapshotBytes).
				Return(nil, nil) // we don't consume the snapshot response apart from checking the error
		}
		snapshotClient.EXPECT().
			GetSnapshots().
			Return([]grafana.SnapshotListResponse{}, nil)

		dashboardClient.EXPECT().
			GetAllFolderIds().
			Return(nil, nil)

		dashboardClient.EXPECT().
			SearchDashboards("", false, tags[0], tags[1]).
			Return([]grafana.FoundBoard{}, nil)

		err := dashboardSyncer.Sync(context.TODO(), &v1.DashboardsSnapshot{
			Upstreams: upstreamList, // we just init'd two new upstreams
		})
		Expect(err).NotTo(HaveOccurred())
	})
})

var _ = Describe("Build Rest Client", func() {
	It("builds a rest client for an HTTP URL when given grafana username and password", func() {
		err := os.Setenv(grafanaUsername, "username")
		Expect(err).NotTo(HaveOccurred())
		err = os.Setenv(grafanaPassword, "password")
		Expect(err).NotTo(HaveOccurred())
		_, err = buildRestClient(context.Background(), "http://test.com")
		Expect(err).NotTo(HaveOccurred())

		// fails without it
		err = os.Unsetenv(grafanaUsername)
		Expect(err).NotTo(HaveOccurred())
		err = os.Unsetenv(grafanaPassword)
		Expect(err).NotTo(HaveOccurred())
		_, err = buildRestClient(context.Background(), "http://test.com")
		Expect(err).To(Equal(grafana.IncompleteGrafanaCredentials))
	})

	It("builds a rest client for an HTTP URL when given grafana api key", func() {
		err := os.Setenv(grafanaApiKey, "apiKey")
		Expect(err).NotTo(HaveOccurred())
		_, err = buildRestClient(context.Background(), "http://test.com")
		Expect(err).NotTo(HaveOccurred())

		// fails without it
		err = os.Unsetenv(grafanaApiKey)
		Expect(err).NotTo(HaveOccurred())
		_, err = buildRestClient(context.Background(), "http://test.com")
		Expect(err).To(Equal(grafana.IncompleteGrafanaCredentials))
	})

	It("Do not require a CA cert if an HTTPS URL is given", func() {
		err := os.Setenv(grafanaUsername, "username")
		Expect(err).NotTo(HaveOccurred())
		err = os.Setenv(grafanaPassword, "password")
		Expect(err).NotTo(HaveOccurred())
		err = os.Unsetenv(grafanaCaCrt)
		Expect(err).NotTo(HaveOccurred())
		_, err = buildRestClient(context.Background(), "https://test.com")
		Expect(err).NotTo(HaveOccurred())
	})
})

var _ = Describe("Load Default Dashboards", func() {
	mockCtrl = gomock.NewController(GinkgoT())
	dashboardClient = mocks.NewMockDashboardClient(mockCtrl)
	templateGenerator := template.NewDefaultJsonGenerator("default", `{"foo":"bar"}`)
	uid := templateGenerator.GenerateUid()

	It("returns before generating when dashboard already exists", func() {
		dashboardClient.EXPECT().GetRawDashboard(uid).
			Return(nil, grafana.BoardProperties{}, nil)

		loadDefaultDashboard(context.Background(), templateGenerator, generalFolderId, dashboardClient)
	})

	It("returns before generating when getting the dashboard results in an unexpected error", func() {
		dashboardClient.EXPECT().GetRawDashboard(uid).
			Return(nil, grafana.BoardProperties{}, fmt.Errorf("fake error"))

		loadDefaultDashboard(context.Background(), templateGenerator, generalFolderId, dashboardClient)
	})

	It("returns attempts to save dashboard if it does not already exist", func() {
		dashboardClient.EXPECT().GetRawDashboard(uid).
			Return(nil, grafana.BoardProperties{}, grafana.DashboardNotFound(uid))
		dashPost, err := templateGenerator.GenerateDashboardPost(generalFolderId)
		Expect(err).NotTo(HaveOccurred())
		dashboardClient.EXPECT().PostDashboard(dashPost).Return(nil)

		loadDefaultDashboard(context.Background(), templateGenerator, generalFolderId, dashboardClient)
	})
})
