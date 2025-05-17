# Repository Guide

This repository contains the `kgateway` project, an Envoy-powered API gateway for Kubernetes. The codebase is primarily written in **Go**, with some **Python** packages (AI extension) and a small **Rust** module under `internal/envoyinit/rustformations`.

## Key directories
- `api/` – Kubernetes API types and CRDs
- `cmd/` – Go entry points for binaries (`kgateway`, `modelschema`, `sds`)
- `examples/` – Sample YAML manifests and plugin examples
- `hack/` – Development scripts and utilities
- `install/helm/` – Helm charts used to deploy kgateway
- `internal/` – Private Go packages for the controller, envoy init container, Helm assets, etc.
- `pkg/` – Reusable Go libraries and plugins
- `python/` – AI extension, docs helpers, and Python unit tests
- `docs/` – Contains a short note that documentation has moved elsewhere
- `test/kubernetes/` – Current end-to-end test framework and utilities

The `devel/` directory is stale and should be ignored. Other subdirectories under `test/` are also deprecated; only `test/kubernetes/` is actively maintained.

## Building and testing
- Run `make help` from the repository root to list Make targets.
- Go unit tests can be run with `make run-tests` (skips e2e tests) or `make go-test` for all tests. Add `TEST_PKG=./path/to/pkg` to limit packages.
- Coverage reports are generated with `make go-test-with-coverage`.
- Python tests for the AI extension live in `python/ai_extension/test/`. Use the provided Makefile in `python/` (`make unit-tests`) or invoke `python -m pytest` directly after creating a virtualenv and installing `requirements-dev.txt`.
- Kubernetes e2e tests reside in `test/kubernetes/`. The [README](test/kubernetes/README.md) explains the philosophy and provides debugging tips. Typical workflow:
  1. Create a kind cluster via `./hack/kind/setup-kind.sh`.
  2. Run the tests with `go test -v -timeout 600s ./test/kubernetes/e2e/tests` or use the IDE configurations described in the README.
- A `Tiltfile` is available for iterative local development using Tilt.

## Documentation
The [`docs/README.md`](docs/README.md) notes that all maintained documentation lives at [kgateway.dev](https://kgateway.dev/) and in [github.com/kgateway-dev/kgateway.dev](https://github.com/kgateway-dev/kgateway.dev/). The files in `docs/` are kept for historical reference only.

