package plugins

import (
	"context"

	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyauthv2 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v2"
	"github.com/gogo/protobuf/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/solo-projects/projects/extauth/pkg/plugins/mocks"
)

//go:generate mockgen -destination mocks/auth_service_mock.go -package mocks github.com/solo-io/ext-auth-plugins/api AuthService

var _ = Describe("Plugin Chain", func() {

	Describe("mergeHeaders function", func() {

		var (
			buildHeader = func(key, value string, append bool) *envoycore.HeaderValueOption {
				return &envoycore.HeaderValueOption{
					Header: &envoycore.HeaderValue{
						Key:   key,
						Value: value,
					},
					Append: &types.BoolValue{
						Value: append,
					},
				}
			}

			getFirst = func() *envoyauthv2.OkHttpResponse {
				return &envoyauthv2.OkHttpResponse{
					Headers: []*envoycore.HeaderValueOption{
						buildHeader("a", "foo", true),
						buildHeader("b", "bar", false),
						buildHeader("c", "baz", true),
					},
				}
			}

			getSecond = func() *envoyauthv2.OkHttpResponse {
				return &envoyauthv2.OkHttpResponse{
					Headers: []*envoycore.HeaderValueOption{
						buildHeader("b", "new-b", true),
						buildHeader("c", "new-c", false),
					},
				}
			}
		)

		It("works as expected", func() {
			result := mergeHeaders(getFirst(), getSecond())
			Expect(result.Headers).To(ConsistOf(
				buildHeader("a", "foo", true),
				buildHeader("b", "bar, new-b", true),
				buildHeader("c", "new-c", false),
			))

			result = mergeHeaders(getSecond(), getFirst())
			Expect(result.Headers).To(ConsistOf(
				buildHeader("a", "foo", true),
				buildHeader("b", "bar", false),
				buildHeader("c", "new-c, baz", true),
			))
		})

		It("covers edge cases", func() {
			result := mergeHeaders(getFirst(), nil)
			Expect(result).To(BeEquivalentTo(getFirst()))

			result = mergeHeaders(nil, getSecond())
			Expect(result).To(BeEquivalentTo(getSecond()))

			result = mergeHeaders(nil, nil)
			Expect(result).To(BeNil())
		})
	})

	Describe("plugin chain execution", func() {

		var (
			ctrl          *gomock.Controller
			pluginWrapper pluginChain
			mockSvc1,
			mockSvc2,
			mockSvc3 *mocks.MockAuthService
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(T)

			mockSvc1 = mocks.NewMockAuthService(ctrl)
			mockSvc2 = mocks.NewMockAuthService(ctrl)
			mockSvc3 = mocks.NewMockAuthService(ctrl)

			pluginWrapper = pluginChain{}
			err := pluginWrapper.AddPlugin("One", mockSvc1)
			Expect(err).NotTo(HaveOccurred())
			err = pluginWrapper.AddPlugin("Two", mockSvc2)
			Expect(err).NotTo(HaveOccurred())
			err = pluginWrapper.AddPlugin("Three", mockSvc3)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("fails when adding plugins with the same name", func() {
			err := pluginWrapper.AddPlugin("One", mockSvc1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(DuplicatePluginNameError("One").Error()))
		})

		It("panics when adding plugins after Start has been called", func() {
			mockSvc1.EXPECT().Start(gomock.Any()).Return(nil).Times(1)
			mockSvc2.EXPECT().Start(gomock.Any()).Return(nil).Times(1)
			mockSvc3.EXPECT().Start(gomock.Any()).Return(nil).Times(1)

			err := pluginWrapper.Start(context.Background())
			Expect(err).NotTo(HaveOccurred())
			Expect(func() { _ = pluginWrapper.AddPlugin("Four", mockSvc1) }).To(Panic())
		})

		Describe("start functions", func() {

			Context("all functions succeed", func() {

				BeforeEach(func() {
					mockSvc1.EXPECT().Start(gomock.Any()).Return(nil).Times(1)
					mockSvc2.EXPECT().Start(gomock.Any()).Return(nil).Times(1)
					mockSvc3.EXPECT().Start(gomock.Any()).Return(nil).Times(1)
				})

				It("runs all start functions", func() {
					err := pluginWrapper.Start(context.Background())
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("a function fails", func() {

				BeforeEach(func() {
					mockSvc1.EXPECT().Start(gomock.Any()).Return(nil).Times(1)
					mockSvc2.EXPECT().Start(gomock.Any()).Return(errors.New("start failed")).Times(1)
					mockSvc3.EXPECT().Start(gomock.Any()).Return(nil).Times(0)
				})

				It("does not run start function after failure", func() {
					err := pluginWrapper.Start(context.Background())
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe("authorize functions", func() {

			Context("all functions succeed", func() {

				BeforeEach(func() {
					mockSvc1.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(api.AuthorizedResponse(), nil).Times(1)
					mockSvc2.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(api.AuthorizedResponse(), nil).Times(1)
					mockSvc3.EXPECT().Authorize(gomock.Any(), gomock.Any()).Return(api.AuthorizedResponse(), nil).Times(1)
				})

				It("runs all authorize functions", func() {
					response, err := pluginWrapper.Authorize(context.Background(), &envoyauthv2.CheckRequest{})
					Expect(err).NotTo(HaveOccurred())
					Expect(response).To(BeEquivalentTo(api.AuthorizedResponse()))
				})
			})

			Context("a function fails", func() {

				BeforeEach(func() {
					mockSvc1.EXPECT().Start(gomock.Any()).Return(nil).Times(1)
					mockSvc2.EXPECT().Start(gomock.Any()).Return(errors.New("start failed")).Times(1)
					mockSvc3.EXPECT().Start(gomock.Any()).Return(nil).Times(0)
				})

				It("does not run start function after failure", func() {
					err := pluginWrapper.Start(context.Background())
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
