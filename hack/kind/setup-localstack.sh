#! /bin/bash

set -o errexit
set -o pipefail
set -o nounset

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

function install_localstack() {
  # Install localstack
  helm repo add localstack-repo https://helm.localstack.cloud
  helm repo update

  helm upgrade -i --create-namespace localstack localstack-repo/localstack --namespace localstack -f ${ROOT_DIR}/localstack-values.yaml
  kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=localstack -n localstack --timeout=120s
}

function get_localstack_endpoint() {
  localstack_port=$(kubectl get -n localstack svc/localstack -o jsonpath="{.spec.ports[0].nodePort}")
  localstack_host=$(kubectl get nodes -o jsonpath="{.items[0].status.addresses[0].address}")
  echo "http://${localstack_host}:${localstack_port}"
}

function create_test_user() {
  local localstack_endpoint=$1
  aws --endpoint-url=${localstack_endpoint} iam create-user --user-name test
}

function add_lambda_function() {
    local localstack_endpoint=$1

    # Use dummy credentials for localstack as it doesn't support IAM
    # for free tier.
    echo "Setting AWS credentials..."
    export AWS_ACCESS_KEY_ID=test
    export AWS_SECRET_ACCESS_KEY=test
    export AWS_DEFAULT_REGION=us-east-1

    echo "Deleting existing function if it exists..."
    aws --endpoint-url=${localstack_endpoint} lambda delete-function --function-name hello-function --no-cli-pager 2>/dev/null || true

    echo "Creating function..."
    aws --endpoint-url=${localstack_endpoint} lambda create-function \
        --no-cli-pager \
        --region us-east-1 \
        --function-name hello-function \
        --zip-file fileb://${ROOT_DIR}/lambda/hello-function.zip \
        --handler hello-function.handler \
        --runtime nodejs18.x \
        --role arn:aws:iam::000000000000:role/my-lambda-role
}

install_localstack
localstack_endpoint=$(get_localstack_endpoint)
add_lambda_function ${localstack_endpoint}

echo "Localstack endpoint: ${localstack_endpoint}"
echo "export ENDPOINT=${localstack_endpoint}"
