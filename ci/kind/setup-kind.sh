#!/bin/bash -ex

# 0. Assign default values to some of our environment variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
CLUSTER_NAME="${CLUSTER_NAME:-kind}"
CLUSTER_NODE_VERSION="${CLUSTER_NODE_VERSION:-v1.31.0}"
VERSION="${VERSION:-1.0.0-ci1}"
SKIP_DOCKER="${SKIP_DOCKER:-false}"
JUST_KIND="${JUST_KIND:-false}"
IMAGE_VARIANT="${IMAGE_VARIANT:-standard}"
CONFORMANCE="${CONFORMANCE:-false}"
CONFORMANCE_VERSION="${CONFORMANCE_VERSION:-v1.2.0}"
CONFORMANCE_CHANNEL="${CONFORMANCE_CHANNEL:-experimental}"
CILIUM_VERSION="${CILIUM_VERSION:-1.15.5}"

function create_kind_cluster_or_skip() {
  activeClusters=$(kind get clusters)
  if [[ "$activeClusters" =~ .*"$CLUSTER_NAME".* ]]; then
    echo "cluster exists, skipping cluster creation"
    return
  fi

  echo "creating cluster ${CLUSTER_NAME}"
  kind create cluster \
    --name "$CLUSTER_NAME" \
    --image "kindest/node:$CLUSTER_NODE_VERSION" \
    --config="$SCRIPT_DIR/cluster.yaml"

  helm repo add cilium-setup-kind https://helm.cilium.io/
  helm repo update
  helm install cilium cilium-setup-kind/cilium --version $CILIUM_VERSION \
   --namespace kube-system \
   --set image.pullPolicy=IfNotPresent \
   --set ipam.mode=kubernetes \
   --set operator.replicas=1
  helm repo remove cilium-setup-kind

  echo "Finished setting up cluster $CLUSTER_NAME"
  if [[ $JUST_KIND == 'true' ]]; then
    echo "JUST_KIND=true, not building images"
    exit
  fi
}

# 1. Create or skip cluster
create_kind_cluster_or_skip

if [[ $SKIP_DOCKER == 'true' ]]; then
  echo "SKIP_DOCKER=true, not building images or chart"
  exit 0
else
  VERSION=$VERSION CLUSTER_NAME=$CLUSTER_NAME IMAGE_VARIANT=$IMAGE_VARIANT make kind-build-and-load

  # Build the test helm chart (if needed)
  VERSION=$VERSION make build-test-chart
fi

# Build the gloo CLI for local usage
make -s build-cli-local

# Apply the Kubernetes Gateway API CRDs
kubectl apply --kustomize "https://github.com/kubernetes-sigs/gateway-api/config/crd/$CONFORMANCE_CHANNEL?ref=$CONFORMANCE_VERSION"

# 6. If we want the conformance environment, set up MetalLB
if [[ $CONFORMANCE == "true" ]]; then
  . $SCRIPT_DIR/setup-metalllb-on-kind.sh
fi

echo "Ready to install the new kgateway chart if needed. Example usage:"
echo "  helm install my-kgateway install/helm/kgateway --namespace my-namespace --create-namespace"