#!/bin/bash -ex

# 0. Assign default values to some of our environment variables
# Get directory this script is located in to access script local files
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
# The name of the kind cluster to deploy to
CLUSTER_NAME="${CLUSTER_NAME:-kind}"
# The version of the Node Docker image to use for booting the cluster
CLUSTER_NODE_VERSION="${CLUSTER_NODE_VERSION:-v1.31.0}"
# The version used to tag images
VERSION="${VERSION:-1.0.0-ci1}"
# Skip building docker images if we are testing a released version
SKIP_DOCKER="${SKIP_DOCKER:-false}"
# Stop after creating the kind cluster
JUST_KIND="${JUST_KIND:-false}"
# Set the default image variant to standard
IMAGE_VARIANT="${IMAGE_VARIANT:-standard}"
# If true, run extra steps to set up k8s gateway api conformance test environment
CONFORMANCE="${CONFORMANCE:-false}"
# The version of the k8s gateway api conformance tests to run. Requires CONFORMANCE=true
CONFORMANCE_VERSION="${CONFORMANCE_VERSION:-v1.2.0}"
# The channel of the k8s gateway api conformance tests to run. Requires CONFORMANCE=true
CONFORMANCE_CHANNEL="${CONFORMANCE_CHANNEL:-"experimental"}"
# The kind CLI to use. Defaults to the latest version from the kind repo.
KIND="${KIND:-go tool kind}"
# If true, use localstack for lambda functions
LOCALSTACK="${LOCALSTACK:-false}"
# If true, install the helm charts
INSTALL_KGATEWAY="${INSTALL_KGATEWAY:-false}"

function create_kind_cluster_or_skip() {
  activeClusters=$(kind get clusters)

  # if the kind cluster exists already, return
  if [[ "$activeClusters" =~ .*"$CLUSTER_NAME".* ]]; then
    echo "cluster exists, skipping cluster creation"
    return
  fi

  echo "creating cluster ${CLUSTER_NAME}"
  $KIND create cluster \
    --name "$CLUSTER_NAME" \
    --image "kindest/node:$CLUSTER_NODE_VERSION" \
    --config="$SCRIPT_DIR/cluster.yaml"
  echo "Finished setting up cluster $CLUSTER_NAME"

  # so that you can just build the kind image alone if needed
  if [[ $JUST_KIND == 'true' ]]; then
    echo "JUST_KIND=true, not building images"
    exit
  fi
}

# 1. Create a kind cluster (or skip creation if a cluster with name=CLUSTER_NAME already exists)
# This config is roughly based on: https://kind.sigs.k8s.io/docs/user/ingress/
create_kind_cluster_or_skip

# 2. Apply the Kubernetes Gateway API CRDs
# Note, we're using kustomize to apply the CRDs from the k8s gateway api repo as
# kustomize supports remote GH URLs and provides more flexibility compared to
# alternatives like running a series of `kubectl apply -f <url>` commands. This
# approach is largely necessary since upstream hasn't adopted a helm chart for
# the CRDs yet, or won't be for the foreseeable future.
kubectl apply --kustomize "https://github.com/kubernetes-sigs/gateway-api/config/crd/$CONFORMANCE_CHANNEL?ref=$CONFORMANCE_VERSION"

if [[ $SKIP_DOCKER == 'true' ]]; then
  # TODO(tim): refactor the Makefile & CI scripts so we're loading local
  # charts to real helm repos, and then we can remove this block.
  echo "SKIP_DOCKER=true, not building images or chart"
  helm repo add gloo https://storage.googleapis.com/solo-public-helm
  helm repo update
else
  # 3. Make all the docker images and load them to the kind cluster
  VERSION=$VERSION CLUSTER_NAME=$CLUSTER_NAME IMAGE_VARIANT=$IMAGE_VARIANT make kind-build-and-load

  # 4. Build the test helm chart, ensuring we have a chart in the `_test` folder
  VERSION=$VERSION make package-kgateway-charts

  # 5. Install the locally packagedhelm charts if the INSTALL_KGATEWAY flag is set
  if [[ $INSTALL_KGATEWAY == "true" ]]; then
    echo "Installing the helm charts"
    helm upgrade -i kgateway-crds _test/kgateway-crds-$VERSION.tgz
    helm upgrade -i kgateway _test/kgateway-$VERSION.tgz --create-namespace -n kgateway-system --set image.registry=ghcr.io/kgateway-dev
  fi
fi

# 5. Conformance test setup
if [[ $CONFORMANCE == "true" ]]; then
  echo "Running conformance test setup"

  . $SCRIPT_DIR/setup-metalllb-on-kind.sh
fi

# 6. Setup localstack
if [[ $LOCALSTACK == "true" ]]; then
  echo "Setting up localstack"
  . $SCRIPT_DIR/setup-localstack.sh
fi
