# imports should be after the set up flags so are lower

# https://www.gnu.org/software/make/manual/html_node/Special-Variables.html#Special-Variables
.DEFAULT_GOAL := help

#----------------------------------------------------------------------------------
# Help
#----------------------------------------------------------------------------------
# Our Makefile is quite large, and hard to reason through
# `make help` can be used to self-document targets
# To update a target to be self-documenting (and appear with the `help` command),
# place a comment after the target that is prefixed by `##`. For example:
#	custom-target: ## comment that will appear in the documentation when running `make help`
#
# **NOTE TO DEVELOPERS**
# As you encounter make targets that are frequently used, please make them self-documenting
.PHONY: help
help: NAME_COLUMN_WIDTH=35
help: LINE_COLUMN_WIDTH=5
help: ## Output the self-documenting make targets
	@grep -hnE '^[%a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "[:]|(## )"}; {printf "\033[36mL%-$(LINE_COLUMN_WIDTH)s%-$(NAME_COLUMN_WIDTH)s\033[0m %s\n", $$1, $$2, $$4}'

#----------------------------------------------------------------------------------
# Base
#----------------------------------------------------------------------------------

ROOTDIR := $(shell pwd)
OUTPUT_DIR ?= $(ROOTDIR)/_output

export IMAGE_REGISTRY ?= cr.kgateway.dev/kgateway-dev

# Kind of a hack to make sure _output exists
z := $(shell mkdir -p $(OUTPUT_DIR))

# A semver resembling 1.0.1-dev. Most calling GHA jobs customize this. Exported for use in goreleaser.yaml.
VERSION ?= 1.0.1-dev
export VERSION

SOURCES := $(shell find . -name "*.go" | grep -v test.go)

export ENVOY_IMAGE ?= quay.io/solo-io/envoy-gloo:1.34.1-patch3
export LDFLAGS := -X 'github.com/kgateway-dev/kgateway/v2/internal/version.Version=$(VERSION)'
export GCFLAGS ?=

UNAME_M := $(shell uname -m)
# if `GO_ARCH` is set, then it will keep its value. Else, it will be changed based off the machine's host architecture.
# if the machines architecture is set to arm64 then we want to set the appropriate values, else we only support amd64
IS_ARM_MACHINE := $(or	$(filter $(UNAME_M), arm64), $(filter $(UNAME_M), aarch64))
ifneq ($(IS_ARM_MACHINE), )
	ifneq ($(GOARCH), amd64)
		GOARCH := arm64
	endif
else
	# currently we only support arm64 and amd64 as a GOARCH option.
	ifneq ($(GOARCH), arm64)
		GOARCH := amd64
	endif
endif

PLATFORM := --platform=linux/$(GOARCH)
PLATFORM_MULTIARCH := $(PLATFORM)
LOAD_OR_PUSH := --load
ifeq ($(MULTIARCH), true)
	PLATFORM_MULTIARCH := --platform=linux/amd64,linux/arm64
	LOAD_OR_PUSH :=

	ifeq ($(MULTIARCH_PUSH), true)
		LOAD_OR_PUSH := --push
	endif
endif

GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
GO_BUILD_FLAGS := CGO_ENABLED=0 GOARCH=$(GOARCH)

TEST_ASSET_DIR ?= $(ROOTDIR)/_test

# This is the location where assets are placed after a test failure
# This is used by our e2e tests to emit information about the running instance of kgateway
BUG_REPORT_DIR := $(TEST_ASSET_DIR)/bug_report
$(BUG_REPORT_DIR):
	mkdir -p $(BUG_REPORT_DIR)

# This is the location where logs are stored for future processing.
# This is used to generate summaries of test outcomes and may be used in the future to automate
# processing of data based on test outcomes.
TEST_LOG_DIR := $(TEST_ASSET_DIR)/test_log
$(TEST_LOG_DIR):
	mkdir -p $(TEST_LOG_DIR)

# Base Alpine image used for all containers. Exported for use in goreleaser.yaml.
# TODO(tim): Why does SDS need an alphine-based base image? Can we remove this and
# use distroless, scratch, chaingurad, etc.
export ALPINE_BASE_IMAGE ?= alpine:3.17.6

#----------------------------------------------------------------------------------
# Macros
#----------------------------------------------------------------------------------

# This macro takes a relative path as its only argument and returns all the files
# in the tree rooted at that directory that match the given criteria.
get_sources = $(shell find $(1) -name "*.go" | grep -v test | grep -v generated.go | grep -v mock_)

#----------------------------------------------------------------------------------
# Repo setup
#----------------------------------------------------------------------------------

.PHONY: init-git-hooks
init-git-hooks:  ## Use the tracked version of Git hooks from this repo
	git config core.hooksPath .githooks

GOIMPORTS ?= go tool goimports
.PHONY: fmt
fmt:  ## Format the code with goimports
	$(GOIMPORTS) -local "github.com/kgateway-dev/kgateway/v2/"  -w $(shell ls -d */ | grep -v vendor)

.PHONY: fmt-changed
fmt-changed:  ## Format the code with goimports
	git diff --name-only | grep '.*.go$$' | xargs -- $(GOIMPORTS) -w

# must be a separate target so that make waits for it to complete before moving on
.PHONY: mod-download
mod-download:  ## Download the dependencies
	go mod download all

.PHONY: mod-tidy
mod-tidy: mod-download  ## Tidy the go mod file
	go mod tidy

#----------------------------------------------------------------------------
# Analyze
#----------------------------------------------------------------------------

# TODO: No need for a dedicated section or an analyze esq target name. Move
# to the lint section.

GO_VERSION := $(shell cat go.mod | grep -E '^go' | awk '{print $$2}')
GOTOOLCHAIN ?= go$(GO_VERSION)

GOLANGCI_LINT ?= go tool golangci-lint
ANALYZE_ARGS ?= --fix --verbose
.PHONY: analyze
analyze:  ## Run golangci-lint. Override options with ANALYZE_ARGS.
	GOTOOLCHAIN=$(GOTOOLCHAIN) $(GOLANGCI_LINT) run $(ANALYZE_ARGS) ./...

#----------------------------------------------------------------------------
# Info
#----------------------------------------------------------------------------

.PHONY: envoyversion
envoyversion: ENVOY_VERSION_TAG ?= $(shell echo $(ENVOY_IMAGE) | cut -d':' -f2)
envoyversion:
	echo "Version is $(ENVOY_VERSION_TAG)"
	echo "Commit for envoyproxy is $(shell curl -s https://raw.githubusercontent.com/solo-io/envoy-gloo/refs/tags/v$(ENVOY_VERSION_TAG)/bazel/repository_locations.bzl | grep "envoy =" -A 4 | grep commit | cut -d'"' -f2)"
	echo "Current ABI in envoyinit can be found in the cargo.toml's envoy-proxy-dynamic-modules-rust-sdk"

#----------------------------------------------------------------------------------
# E2E Tests
#----------------------------------------------------------------------------------

GO_E2E_TEST_PKGS ?= ./test/kubernetes/e2e/tests
GO_E2E_TEST_ARGS ?= -v -timeout=600s -ldflags='$(LDFLAGS)' $(GO_TEST_ARGS)
GO_E2E_USER_TEST_ARGS ?=

# TODO(tim): I think we'll need to either explicitly skip the tests that need additional setup, or
# add build tags that allow us to determine which tests to run. I think I'm leaning towards the latter.
# Note: this isn't a problem with CI as we're explicitly configuring which suites are run, but a problem
# if you run this locally.
.PHONY: e2e
e2e: setup test-e2e ## Set up development environment and run E2E tests

.PHONY: test-e2e
test-e2e: ## Run E2E tests
	go test $(GO_E2E_TEST_ARGS) $(GO_E2E_USER_TEST_ARGS) $(GO_E2E_TEST_PKGS)

deploy-localstack: ## Deploy LocalStack
	./hack/kind/setup-localstack.sh

e2e-lambda: setup deploy-localstack test-e2e-lambda

# TODO(tim): Move the lambda to a dedicated AWS suite? The regex works well, and we don't have a need to set
# any custom helm chart values for lambda functionality, so likely the wrong approach.
test-e2e-lambda: ## Run E2E tests with Lambda support
	go test $(GO_E2E_TEST_ARGS) -run "^TestKgateway$$/Lambda$$" $(GO_E2E_TEST_PKGS)

e2e-ai: setup kind-build-and-load-test-ai-provider kind-build-and-load-kgateway-ai-extension test-e2e-ai

test-e2e-ai: ## Run E2E tests with AI support
	go test $(GO_E2E_TEST_ARGS) -run '^TestAIExtensions$$' $(GO_E2E_TEST_PKGS)

e2e-agent-gateway: setup kind-build-and-load-test-a2a-agent test-e2e-agent-gateway

test-e2e-agent-gateway: ## Run E2E tests with Agent Gateway support
	go test $(GO_E2E_TEST_ARGS) -run '^TestAgentGatewayIntegration$$' $(GO_E2E_TEST_PKGS)

#----------------------------------------------------------------------------------
# Env test
#----------------------------------------------------------------------------------

# TODO(tim): Dynamic envtest version based on client-go dependency
ENVTEST_K8S_VERSION = 1.23
ENVTEST ?= go tool setup-envtest
# The setup suite is the envtest-based suite right now.
ENVTEST_PKGS ?= ./internal/kgateway/setup
# Provide a way to override the test to run, e.g. ENVTEST_ARGS='-run TestWithStandardSettings'
ENVTEST_ARGS ?=

.PHONY: envtest-path
envtest-path: ## Set the envtest path
	@$(ENVTEST) use $(ENVTEST_K8S_VERSION) -p path --arch=amd64

.PHONY: envtest
envtest: envtest-path ## Run envtest suite (controller-runtime integration tests)
	go test -v -ldflags='$(LDFLAGS)' $(GO_TEST_ARGS) $(ENVTEST_PKGS) $(ENVTEST_ARGS)

#----------------------------------------------------------------------------------
# Unit Tests
#----------------------------------------------------------------------------------

GO_TEST_ENV ?=
# Testings flags: https://pkg.go.dev/cmd/go#hdr-Testing_flags
# The default timeout for a suite is 10 minutes, but this can be overridden by setting the -timeout flag. Currently set
# to 25 minutes based on the time it takes to run the longest test setup (kgateway_test).
GO_TEST_ARGS ?= -timeout=25m -cpu=4 -race -outputdir=$(OUTPUT_DIR)
GO_UNIT_TEST_PKGS ?= $(shell go list ./... | grep -v hack | grep -v test | grep -v internal/kgateway/setup)

.PHONY: unit-go
unit-go: ## Run Go unit tests with coverage and validation
unit-go: mod-download kgateway
	go test -ldflags='$(LDFLAGS)' $(GO_TEST_ARGS) $(GO_UNIT_TEST_PKGS)

# TODO(tim): Needs more thought, but it's passing CI. Cannot run this target locally.
.PHONY: unit-python
unit-python: ## Run Python unit tests
unit-python:
	cd python && python3 -m venv .venv
	cd python && .venv/bin/pip install -r requirements-dev.txt --no-cache-dir
	cd python && .venv/bin/pip install -r requirements.txt --no-cache-dir
	cd python/ai_extension && ../.venv/bin/python -m pytest test

.PHONY: unit
unit: unit-go unit-python ## Run all unit tests

#----------------------------------------------------------------------------------
# Clean
#----------------------------------------------------------------------------------

# Important to clean before pushing new releases. Dockerfiles and binaries may not update properly
# FIXME(tim): Why doesn't this clean _everything_ and run all the sub-clean targets?
.PHONY: clean
clean:
	rm -rf _output
	rm -rf _test
	git clean -f -X install

# Clean generated code
# see hack/generate.sh for source of truth of dirs to clean
.PHONY: clean-gen
clean-gen:
	rm -rf api/applyconfiguration
	rm -rf pkg/generated/openapi
	rm -rf pkg/client
	rm -f install/helm/kgateway-crds/templates/gateway.kgateway.dev_*.yaml

.PHONY: clean-tests
clean-tests:
	find * -type f -name '*.test' -exec rm {} \;
	find * -type f -name '*.cov' -exec rm {} \;
	find * -type f -name 'junit*.xml' -exec rm {} \;

.PHONY: clean-bug-report
clean-bug-report:
	rm -rf $(BUG_REPORT_DIR)

.PHONY: clean-test-logs
clean-test-logs:
	rm -rf $(TEST_LOG_DIR)

#----------------------------------------------------------------------------------
# Generated Code
#----------------------------------------------------------------------------------

.PHONY: verify
verify: generate-all  ## Verify that generated code is up to date
	git diff -U3 --exit-code

.PHONY: generate-all
generate-all: generated-code

.PHONY: generated-code
generated-code: clean-gen go-generate-all mod-tidy generate-licenses fmt ## Generate all required code, cleaning and formatting as well; this target is executed in CI

.PHONY: go-generate-all
go-generate-all: go-generate-apis go-generate-mocks

.PHONY: go-generate-apis
go-generate-apis: ## Run all go generate directives in the repo, including codegen for protos, mockgen, and more
	GO111MODULE=on go generate ./hack/...

.PHONY: go-generate-mocks
go-generate-mocks: ## Runs all generate directives for mockgen in the repo
	GO111MODULE=on go generate -run="mockgen" ./...

.PHONY: generate-licenses
generate-licenses: ## Update the licenses for the project
	GO111MODULE=on go run hack/utils/oss_compliance/oss_compliance.go osagen -c "GNU General Public License v2.0,GNU General Public License v3.0,GNU Lesser General Public License v2.1,GNU Lesser General Public License v3.0,GNU Affero General Public License v3.0"
	GO111MODULE=on go run hack/utils/oss_compliance/oss_compliance.go osagen -s "Mozilla Public License 2.0,GNU General Public License v2.0,GNU General Public License v3.0,GNU Lesser General Public License v2.1,GNU Lesser General Public License v3.0,GNU Affero General Public License v3.0"> hack/utils/oss_compliance/osa_provided.md
	GO111MODULE=on go run hack/utils/oss_compliance/oss_compliance.go osagen -i "Mozilla Public License 2.0"> hack/utils/oss_compliance/osa_included.md

#----------------------------------------------------------------------------------
# AI Extensions ExtProc Server
#----------------------------------------------------------------------------------

PYTHON_DIR := $(ROOTDIR)/python

export AI_EXTENSION_IMAGE_REPO ?= kgateway-ai-extension
.PHONY: kgateway-ai-extension-docker
kgateway-ai-extension-docker:
	docker buildx build $(LOAD_OR_PUSH) $(PLATFORM_MULTIARCH) -f $(PYTHON_DIR)/Dockerfile $(ROOTDIR) \
		--build-arg PYTHON_DIR=python \
		-t  $(IMAGE_REGISTRY)/kgateway-ai-extension:$(VERSION)

#----------------------------------------------------------------------------------
# Controller
#----------------------------------------------------------------------------------

K8S_GATEWAY_DIR=internal/kgateway
K8S_GATEWAY_SOURCES=$(call get_sources,$(K8S_GATEWAY_DIR))
CONTROLLER_OUTPUT_DIR=$(OUTPUT_DIR)/$(K8S_GATEWAY_DIR)
export CONTROLLER_IMAGE_REPO ?= kgateway

# We include the files in K8S_GATEWAY_DIR as dependencies to the kgateway build
# so changes in those directories cause the make target to rebuild
$(CONTROLLER_OUTPUT_DIR)/kgateway-linux-$(GOARCH): $(K8S_GATEWAY_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags='$(LDFLAGS)' -gcflags='$(GCFLAGS)' -o $@ ./cmd/kgateway/...

.PHONY: kgateway
kgateway: $(CONTROLLER_OUTPUT_DIR)/kgateway-linux-$(GOARCH)

$(CONTROLLER_OUTPUT_DIR)/Dockerfile: cmd/kgateway/Dockerfile
	cp $< $@

.PHONY: kgateway-docker
kgateway-docker: $(CONTROLLER_OUTPUT_DIR)/kgateway-linux-$(GOARCH) $(CONTROLLER_OUTPUT_DIR)/Dockerfile
	docker buildx build --load $(PLATFORM) $(CONTROLLER_OUTPUT_DIR) -f $(CONTROLLER_OUTPUT_DIR)/Dockerfile \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg ENVOY_IMAGE=$(ENVOY_IMAGE) \
		-t $(IMAGE_REGISTRY)/$(CONTROLLER_IMAGE_REPO):$(VERSION)

#----------------------------------------------------------------------------------
# SDS Server - gRPC server for serving Secret Discovery Service config
#----------------------------------------------------------------------------------

SDS_DIR=internal/sds
SDS_SOURCES=$(call get_sources,$(SDS_DIR))
SDS_OUTPUT_DIR=$(OUTPUT_DIR)/$(SDS_DIR)
export SDS_IMAGE_REPO ?= sds

$(SDS_OUTPUT_DIR)/sds-linux-$(GOARCH): $(SDS_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags='$(LDFLAGS)' -gcflags='$(GCFLAGS)' -o $@ ./cmd/sds/...

.PHONY: sds
sds: $(SDS_OUTPUT_DIR)/sds-linux-$(GOARCH)

$(SDS_OUTPUT_DIR)/Dockerfile.sds: cmd/sds/Dockerfile
	cp $< $@

.PHONY: sds-docker
sds-docker: $(SDS_OUTPUT_DIR)/sds-linux-$(GOARCH) $(SDS_OUTPUT_DIR)/Dockerfile.sds
	docker buildx build --load $(PLATFORM) $(SDS_OUTPUT_DIR) -f $(SDS_OUTPUT_DIR)/Dockerfile.sds \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg BASE_IMAGE=$(ALPINE_BASE_IMAGE) \
		-t $(IMAGE_REGISTRY)/$(SDS_IMAGE_REPO):$(VERSION)

#----------------------------------------------------------------------------------
# Envoy init (BASE/SIDECAR)
#----------------------------------------------------------------------------------

ENVOYINIT_DIR=internal/envoyinit
ENVOYINIT_SOURCES=$(call get_sources,$(ENVOYINIT_DIR))
ENVOYINIT_OUTPUT_DIR=$(OUTPUT_DIR)/$(ENVOYINIT_DIR)
export ENVOYINIT_IMAGE_REPO ?= envoy-wrapper

$(ENVOYINIT_OUTPUT_DIR)/envoyinit-linux-$(GOARCH): $(ENVOYINIT_SOURCES)
	$(GO_BUILD_FLAGS) GOOS=linux go build -ldflags='$(LDFLAGS)' -gcflags='$(GCFLAGS)' -o $@ ./internal/envoyinit/cmd/...

.PHONY: envoyinit
envoyinit: $(ENVOYINIT_OUTPUT_DIR)/envoyinit-linux-$(GOARCH)

# TODO(nfuden) cheat the process for now with -r but try to find a cleaner method
$(ENVOYINIT_OUTPUT_DIR)/Dockerfile.envoyinit: internal/envoyinit/Dockerfile.envoyinit
	cp -r ${ENVOYINIT_DIR}/rustformations $(ENVOYINIT_OUTPUT_DIR)
	cp $< $@

$(ENVOYINIT_OUTPUT_DIR)/docker-entrypoint.sh: internal/envoyinit/cmd/docker-entrypoint.sh
	cp $< $@

.PHONY: envoy-wrapper-docker
envoy-wrapper-docker: $(ENVOYINIT_OUTPUT_DIR)/envoyinit-linux-$(GOARCH) $(ENVOYINIT_OUTPUT_DIR)/Dockerfile.envoyinit $(ENVOYINIT_OUTPUT_DIR)/docker-entrypoint.sh
	docker buildx build --load $(PLATFORM) $(ENVOYINIT_OUTPUT_DIR) -f $(ENVOYINIT_OUTPUT_DIR)/Dockerfile.envoyinit \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg ENVOY_IMAGE=$(ENVOY_IMAGE) \
		-t $(IMAGE_REGISTRY)/$(ENVOYINIT_IMAGE_REPO):$(VERSION)

#----------------------------------------------------------------------------------
# Helm
#----------------------------------------------------------------------------------

HELM ?= go tool helm
HELM_PACKAGE_ARGS ?= --version $(VERSION)
HELM_CHART_DIR=install/helm/kgateway
HELM_CHART_DIR_CRD=install/helm/kgateway-crds

.PHONY: package-kgateway-charts
package-kgateway-charts: package-kgateway-chart package-kgateway-crd-chart ## Package the kgateway charts

.PHONY: package-kgateway-chart
package-kgateway-chart: ## Package the kgateway charts
	mkdir -p $(TEST_ASSET_DIR); \
	$(HELM) package $(HELM_PACKAGE_ARGS) --destination $(TEST_ASSET_DIR) $(HELM_CHART_DIR); \
	$(HELM) repo index $(TEST_ASSET_DIR);

.PHONY: package-kgateway-crd-chart
package-kgateway-crd-chart: ## Package the kgateway crd chart
	mkdir -p $(TEST_ASSET_DIR); \
	$(HELM) package $(HELM_PACKAGE_ARGS) --destination $(TEST_ASSET_DIR) $(HELM_CHART_DIR_CRD); \
	$(HELM) repo index $(TEST_ASSET_DIR);

.PHONY: deploy-kgateway-crd-chart
deploy-kgateway-crd-chart: ## Deploy the kgateway crd chart
	$(HELM) upgrade --install kgateway-crds $(TEST_ASSET_DIR)/kgateway-crds-$(VERSION).tgz --namespace kgateway-system --create-namespace

HELM_ADDITIONAL_VALUES ?= hack/helm/dev.yaml
.PHONY: deploy-kgateway-chart
deploy-kgateway-chart: ## Deploy the kgateway chart
	$(HELM) upgrade --install kgateway $(TEST_ASSET_DIR)/kgateway-$(VERSION).tgz \
	--namespace kgateway-system --create-namespace \
	--set image.tag=$(VERSION) \
	-f $(HELM_ADDITIONAL_VALUES)

.PHONY: lint-kgateway-charts
lint-kgateway-charts: ## Lint the kgateway charts
	$(HELM) lint $(HELM_CHART_DIR)
	$(HELM) lint $(HELM_CHART_DIR_CRD)

#----------------------------------------------------------------------------------
# Release
#----------------------------------------------------------------------------------

GORELEASER ?= go tool github.com/goreleaser/goreleaser/v2
GORELEASER_ARGS ?= --snapshot --clean
GORELEASER_TIMEOUT ?= 60m
GORELEASER_CURRENT_TAG ?= $(VERSION)

.PHONY: release
release: ## Create a release using goreleaser
	GORELEASER_CURRENT_TAG=$(GORELEASER_CURRENT_TAG) $(GORELEASER) release $(GORELEASER_ARGS) --timeout $(GORELEASER_TIMEOUT)

#----------------------------------------------------------------------------------
# Docker
#----------------------------------------------------------------------------------

.PHONY: docker
docker: kgateway-docker ## Build docker images
docker: envoy-wrapper-docker
docker: sds-docker
docker: kgateway-ai-extension-docker

.PHONY: docker-push
docker-push: docker-push-kgateway
docker-push: docker-push-envoy-wrapper
docker-push: docker-push-sds
docker-push: docker-push-kgateway-ai-extension

.PHONY: docker-retag
docker-retag: docker-retag-kgateway
docker-retag: docker-retag-envoy-wrapper
docker-retag: docker-retag-sds
docker-retag: docker-retag-kgateway-ai-extension

docker-retag-%:
	docker tag $(ORIGINAL_IMAGE_REGISTRY)/$*:$(VERSION) $(IMAGE_REGISTRY)/$*:$(VERSION)

docker-push-%:
	docker push $(IMAGE_REGISTRY)/$*:$(VERSION)

#----------------------------------------------------------------------------------
# Development
#----------------------------------------------------------------------------------

KIND ?= go tool kind
CLUSTER_NAME ?= kind
INSTALL_NAMESPACE ?= kgateway-system

.PHONY: kind-create
kind-create: ## Create a KinD cluster
	$(KIND) get clusters | grep $(CLUSTER_NAME) || $(KIND) create cluster --name $(CLUSTER_NAME)

CONFORMANCE_CHANNEL ?= experimental
CONFORMANCE_VERSION ?= v1.3.0
.PHONY: gw-api-crds
gw-api-crds: ## Install the Gateway API CRDs
	kubectl apply --kustomize "https://github.com/kubernetes-sigs/gateway-api/config/crd/$(CONFORMANCE_CHANNEL)?ref=$(CONFORMANCE_VERSION)"

.PHONY: metallb
metallb: ## Install the MetalLB load balancer
	./hack/kind/setup-metalllb-on-kind.sh

.PHONY: deploy-kgateway
deploy-kgateway: package-kgateway-charts deploy-kgateway-crd-chart deploy-kgateway-chart ## Deploy the kgateway chart and CRDs

.PHONY: setup
setup: kind-create kind-build-and-load gw-api-crds metallb package-kgateway-charts ## Set up basic infrastructure (kind cluster, images, CRDs, MetalLB)

.PHONY: run
run: setup deploy-kgateway  ## Set up complete development environment

#----------------------------------------------------------------------------------
# Build assets for kubernetes e2e tests
#----------------------------------------------------------------------------------

kind-load-%:
	$(KIND) load docker-image $(IMAGE_REGISTRY)/$*:$(VERSION) --name $(CLUSTER_NAME)

# Build an image and load it into the KinD cluster
# Depends on: IMAGE_REGISTRY, VERSION, CLUSTER_NAME
# Envoy image may be specified via ENVOY_IMAGE on the command line or at the top of this file
kind-build-and-load-%: %-docker kind-load-% ; ## Use to build specified image and load it into kind

# Update the docker image used by a deployment
# This works for most of our deployments because the deployment name and container name both match
# NOTE TO DEVS:
#	I explored using a special format of the wildcard to pass deployment:image,
# 	but ran into some challenges with that pattern, while calling this target from another one.
#	It could be a cool extension to support, but didn't feel pressing so I stopped
kind-set-image-%:
	kubectl rollout pause deployment $* -n $(INSTALL_NAMESPACE) || true
	kubectl set image deployment/$* $*=$(IMAGE_REGISTRY)/$*:$(VERSION) -n $(INSTALL_NAMESPACE)
	kubectl patch deployment $* -n $(INSTALL_NAMESPACE) -p '{"spec": {"template":{"metadata":{"annotations":{"kgateway-kind-last-update":"$(shell date)"}}}} }'
	kubectl rollout resume deployment $* -n $(INSTALL_NAMESPACE)

# Reload an image in KinD
# This is useful to developers when changing a single component
# You can reload an image, which means it will be rebuilt and reloaded into the kind cluster, and the deployment
# will be updated to reference it
# Depends on: IMAGE_REGISTRY, VERSION, INSTALL_NAMESPACE , CLUSTER_NAME
# Envoy image may be specified via ENVOY_IMAGE on the command line or at the top of this file
kind-reload-%: kind-build-and-load-% kind-set-image-% ; ## Use to build specified image, load it into kind, and restart its deployment

.PHONY: kind-build-and-load ## Use to build all images and load them into kind
kind-build-and-load: kind-build-and-load-kgateway
kind-build-and-load: kind-build-and-load-envoy-wrapper
kind-build-and-load: kind-build-and-load-sds

.PHONY: kind-load ## Use to load all images into kind
kind-load: kind-load-kgateway
kind-load: kind-load-envoy-wrapper
kind-load: kind-load-sds

define kind_reload_msg
The kind-reload-% targets exist in order to assist developers with the work cycle of
build->test->change->build->test. To that end, rebuilding/reloading every image, then
restarting every deployment is seldom necessary. Consider using kind-reload-% to do so
for a specific component, or kind-build-and-load to push new images for every component.
endef
export kind_reload_msg
.PHONY: kind-reload
kind-reload:
	@echo "$$kind_reload_msg"

# Useful utility for listing images loaded into the kind cluster
.PHONY: kind-list-images
kind-list-images: ## List solo-io images in the kind cluster named {CLUSTER_NAME}
	docker exec -ti $(CLUSTER_NAME)-control-plane crictl images | grep "solo-io"

# Useful utility for pruning images that were previously loaded into the kind cluster
.PHONY: kind-prune-images
kind-prune-images: ## Remove images in the kind cluster named {CLUSTER_NAME}
	docker exec -ti $(CLUSTER_NAME)-control-plane crictl rmi --prune

#----------------------------------------------------------------------------------
# A2A Test Server (for agentgateway a2a integration in e2e tests)
#----------------------------------------------------------------------------------

TEST_A2A_AGENT_SERVER_DIR := $(ROOTDIR)/test/kubernetes/e2e/features/agentgateway/a2a-example
.PHONY: test-a2a-agent-docker
test-a2a-agent-docker:
	docker buildx build $(LOAD_OR_PUSH) $(PLATFORM_MULTIARCH) -f $(TEST_A2A_AGENT_SERVER_DIR)/Dockerfile $(TEST_A2A_AGENT_SERVER_DIR) \
		-t $(IMAGE_REGISTRY)/test-a2a-agent:$(VERSION)

#----------------------------------------------------------------------------------
# AI Extensions Test Server (for mocking AI Providers in e2e tests)
#----------------------------------------------------------------------------------

TEST_AI_PROVIDER_SERVER_DIR := $(ROOTDIR)/test/mocks/mock-ai-provider-server
.PHONY: test-ai-provider-docker
test-ai-provider-docker:
	docker buildx build $(LOAD_OR_PUSH) $(PLATFORM_MULTIARCH) -f $(TEST_AI_PROVIDER_SERVER_DIR)/Dockerfile $(TEST_AI_PROVIDER_SERVER_DIR) \
		-t $(IMAGE_REGISTRY)/test-ai-provider:$(VERSION)

#----------------------------------------------------------------------------------
# Targets for running Kubernetes Gateway API conformance tests
#----------------------------------------------------------------------------------

# Pull the conformance test suite from the k8s gateway api repo and copy it into the test dir.
$(TEST_ASSET_DIR)/conformance/conformance_test.go:
	mkdir -p $(TEST_ASSET_DIR)/conformance
	echo "//go:build conformance" > $@
	cat $(shell go list -json -m sigs.k8s.io/gateway-api | jq -r '.Dir')/conformance/conformance_test.go >> $@
	go fmt $@

CONFORMANCE_UNSUPPORTED_FEATURES ?= -exempt-features=GatewayAddressEmpty,GatewayHTTPListenerIsolation,GatewayInfrastructurePropagation,GatewayPort8080,GatewayStaticAddresses,HTTPRouteBackendRequestHeaderModification,HTTPRouteDestinationPortMatching,HTTPRouteParentRefPort,HTTPRouteRequestMultipleMirrors,HTTPRouteRequestPercentageMirror
CONFORMANCE_SUPPORTED_PROFILES ?= -conformance-profiles=GATEWAY-HTTP
CONFORMANCE_GATEWAY_CLASS ?= kgateway
CONFORMANCE_REPORT_ARGS ?= -report-output=$(TEST_ASSET_DIR)/conformance/$(VERSION)-report.yaml -organization=kgateway-dev -project=kgateway -version=$(VERSION) -url=github.com/kgateway-dev/kgateway -contact=github.com/kgateway-dev/kgateway/issues/new/choose
CONFORMANCE_ARGS := -gateway-class=$(CONFORMANCE_GATEWAY_CLASS) $(CONFORMANCE_UNSUPPORTED_FEATURES) $(CONFORMANCE_SUPPORTED_PROFILES) $(CONFORMANCE_REPORT_ARGS)

.PHONY: conformance ## Run the conformance test suite
test/conformance: $(TEST_ASSET_DIR)/conformance/conformance_test.go
	go test -mod=mod -ldflags='$(LDFLAGS)' -tags conformance -test.v $(TEST_ASSET_DIR)/conformance/... -args $(CONFORMANCE_ARGS)

# Run only the specified conformance test. The name must correspond to the ShortName of one of the k8s gateway api
# conformance tests.
test/conformance-%: $(TEST_ASSET_DIR)/conformance/conformance_test.go
	go test -mod=mod -ldflags='$(LDFLAGS)' -tags conformance -test.v $(TEST_ASSET_DIR)/conformance/... -args $(CONFORMANCE_ARGS) \
	-run-test=$*

conformance: run test/conformance

#----------------------------------------------------------------------------------
# Printing makefile variables utility
#----------------------------------------------------------------------------------

# use `make print-MAKEFILE_VAR` to print the value of MAKEFILE_VAR

print-%  : ; @echo $($*)
