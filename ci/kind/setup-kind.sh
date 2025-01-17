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
# The version of Cilium to install.
CILIUM_VERSION="${CILIUM_VERSION:-1.15.5}"

function create_kind_cluster_or_skip() {
  activeClusters=$(kind get clusters)

  # if the kind cluster exists already, return
  if [[ "$activeClusters" =~ .*"$CLUSTER_NAME".* ]]; then
    echo "cluster exists, skipping cluster creation"
    return
  fi

  echo "creating cluster ${CLUSTER_NAME}"
  kind create cluster \
    --name "$CLUSTER_NAME" \
    --image "kindest/node:$CLUSTER_NODE_VERSION" \
    --config="$SCRIPT_DIR/cluster.yaml"

  # Install cilium as we need to define custom network policies to simulate kube api server unavailability
  # in some of our kube2e tests
  helm repo add cilium-setup-kind https://helm.cilium.io/
  helm repo update
  helm install cilium cilium-setup-kind/cilium --version $CILIUM_VERSION \
   --namespace kube-system \
   --set image.pullPolicy=IfNotPresent \
   --set ipam.mode=kubernetes \
   --set operator.replicas=1
  helm repo remove cilium-setup-kind
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

if [[ $SKIP_DOCKER == 'true' ]]; then
  # TODO(tim): refactor the Makefile & CI scripts so we're loading local
  # charts to real helm repos, and then we can remove this block.
  echo "SKIP_DOCKER=true, not building images or chart"
  helm repo add gloo https://storage.googleapis.com/solo-public-helm
  helm repo update
else
  # 2. Make all the docker images and load them to the kind cluster
  VERSION=$VERSION CLUSTER_NAME=$CLUSTER_NAME IMAGE_VARIANT=$IMAGE_VARIANT make kind-build-and-load

  # 3. Build the test helm chart, ensuring we have a chart in the `_test` folder
  VERSION=$VERSION make package-kgateway-chart
fi

# 4. Build the gloo command line tool, ensuring we have one in the `_output` folder
make -s build-cli-local

# 5. Apply the Kubernetes Gateway API CRDs
# Note, we're using kustomize to apply the CRDs from the k8s gateway api repo as
# kustomize supports remote GH URLs and provides more flexibility compared to
# alternatives like running a series of `kubectl apply -f <url>` commands. This
# approach is largely necessary since upstream hasn't adopted a helm chart for
# the CRDs yet, or won't be for the foreseeable future.
kubectl apply --kustomize "https://github.com/kubernetes-sigs/gateway-api/config/crd/$CONFORMANCE_CHANNEL?ref=$CONFORMANCE_VERSION"

# 6. Conformance test setup
if [[ $CONFORMANCE == "true" ]]; then
  echo "Running conformance test setup"

  if [[ "$(uname -s)" == "Darwin" ]]; then
    echo "Detected macOS. Installing cloud-provider-kind@v0.4.0."
    go install sigs.k8s.io/cloud-provider-kind@v0.4.0

    echo "Starting cloud-provider-kind in the background."
    # We redirect stdout/stderr to a file to avoid cluttering logs, then
    # store the process PID so we can terminate it if desired.
    nohup sudo "$(go env GOPATH)/bin/cloud-provider-kind" &> cloud-provider-kind.log &
    CLOUD_PROVIDER_KIND_PID=$!

    cat <<EOF
    ==================================
    cloud-provider-kind is now installed and running in the background.
    Background process PID: ${CLOUD_PROVIDER_KIND_PID}
    Log file: cloud-provider-kind.log

    Next steps:
      1) If you need to debug, check the cloud-provider-kind.log file.
      2) When this job finishes, you may optionally kill the background process with:
          kill ${CLOUD_PROVIDER_KIND_PID}
      3) Your control plane node is labeled to exclude LB traffic by default.
        To allow the LB to function on macOS, remove the label with:
          kubectl label node \${CLUSTER_NAME}-control-plane node.kubernetes.io/exclude-from-external-load-balancers-
    ==================================
EOF
  elif [[ "$(uname -s)" == "Linux" ]]; then
    . $SCRIPT_DIR/setup-metalllb-on-kind.sh
  fi
fi
