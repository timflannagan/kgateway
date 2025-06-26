# kgateway

Note, all commands should be run from the root of the kgateway repo.

## Quickstart

### Setup

<!-- TODO: Deflate this README.md? -->

To create the local test environment in kind, run:

```shell
make run
```

To create a gateway, use the Gateway resource:

```shell
kubectl apply -f - <<EOF
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: http
spec:
  gatewayClassName: kgateway
  listeners:
  - allowedRoutes:
      namespaces:
        from: All
    name: http
    port: 8080
    protocol: HTTP
EOF
```

Apply a test application such as bookinfo:

```shell
kubectl create namespace bookinfo
kubectl apply -f https://raw.githubusercontent.com/istio/istio/release-1.20/samples/bookinfo/platform/kube/bookinfo.yaml -n bookinfo
```

Then create a corresponding HTTPRoute:

```shell
kubectl apply -f- <<EOF
apiVersion: gateway.networking.k8s.io/v1beta1
kind: HTTPRoute
metadata:
  name: productpage
  namespace: bookinfo
  labels:
    example: productpage-route
spec:
  parentRefs:
    - name: http
      namespace: default
  hostnames:
    - "www.example.com"
  rules:
    - backendRefs:
        - name: productpage
          port: 9080
EOF
```

### Testing

Expose the gateway that gets created via the Gateway resource:

```shell
kubectl port-forward deployment/http 8080:8080
```

Send some traffic through the gateway:

```shell
curl -I localhost:8080/productpage -H "host: www.example.com" -v
```

## Istio Integration

Let's bootstrap the test environment with the Istio auto mTLS feature enabled:

```shell
HELM_ADDITIONAL_VALUES=<(cat <<EOF
controller:
  logLevel: debug
  extraEnv:
    KGW_ENABLE_ISTIO_AUTO_MTLS: "true"
EOF
) make run
```

Next, we need to install Istio in the cluster along with the bookinfo test application in the mesh:

```shell
./hack/istio.sh
```

To create a gateway, use the Gateway resource:

```shell
kubectl apply -f - <<EOF
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: http
spec:
  gatewayClassName: kgateway
  listeners:
  - allowedRoutes:
      namespaces:
        from: All
    name: http
    port: 8080
    protocol: HTTP
EOF
```

Then create a corresponding HTTPRoute:

```shell
kubectl apply -f- <<EOF
apiVersion: gateway.networking.k8s.io/v1beta1
kind: HTTPRoute
metadata:
  name: productpage
  namespace: bookinfo
  labels:
    example: productpage-route
spec:
  parentRefs:
    - name: http
      namespace: default
  hostnames:
    - "www.example.com"
  rules:
    - backendRefs:
        - name: productpage
          port: 9080
EOF
```

Then expose the gateway that gets created via the Gateway resource:

```shell
kubectl port-forward deployment/http 8080:8080
```

Send some traffic through the gateway:

```shell
curl localhost:8080/productpage -H "host: www.example.com" -v | grep Book
```

Now turn on strict mode for bookinfo to ensure the kgateway proxy can still connect:

```shell
kubectl apply -f- <<EOF
apiVersion: security.istio.io/v1
kind: PeerAuthentication
metadata:
  name: default
  namespace: bookinfo
spec:
  mtls:
    mode: STRICT
EOF
```

Then send another request:

```shell
curl localhost:8080/productpage -H "host: www.example.com" -v | grep Book
```

Test sending traffic to an application not in the mesh:

```shell
# Create non-mesh app (helloworld namespace is not labeled for istio injection)
kubectl create namespace helloworld
kubectl apply -f https://raw.githubusercontent.com/istio/istio/master/samples/helloworld/helloworld.yaml -n helloworld
```

Apply an HTTPRoute for helloworld:

```shell
kubectl apply -f- <<EOF
apiVersion: gateway.networking.k8s.io/v1beta1
kind: HTTPRoute
metadata:
  name: helloworld
  namespace: helloworld
  labels:
    example: helloworld-route
spec:
  parentRefs:
    - name: http
      namespace: default
  hostnames:
    - "helloworld"
  rules:
    - backendRefs:
        - name: helloworld
          port: 5000
EOF
```

Send traffic to the non-mesh app:

```shell
curl localhost:8080/hello -H "host: helloworld" -v
```
