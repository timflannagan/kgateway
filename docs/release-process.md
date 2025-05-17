# Release Process

This document outlines the process for releasing new versions of the kgateway project. The process is largely automated via GitHub Actions and GoReleaser.

## Overview

Releases involve building Go binaries, creating and publishing multi-architecture Docker images, packaging and publishing Helm charts, and creating a GitHub Release with release notes and associated artifacts.

## Release Triggers

The release workflow (`.github/workflows/release.yaml`) can be triggered in several ways:

1.  **Manual Dispatch (`workflow_dispatch`):**
    *   Maintainers can manually trigger the workflow from the GitHub Actions UI.
    *   **Inputs:**
        *   `version` (optional string): An override for the release version (e.g., `v1.2.3`). If not provided, a development version like `v0.0.0-manual-<git-sha>` is generated. The provided version must be a valid semantic version starting with 'v'.
        *   `validate` (optional boolean, default: `false`): If set to `true`, the release validation job will be run.
    *   This is the primary method for creating official tagged releases. Maintainers supply the desired tag (e.g., `v1.2.3`), and GoReleaser creates this tag and the corresponding GitHub Release.

2.  **Tag Push (`push: tags: v*`):**
    *   Pushing a Git tag that matches the pattern `v*` (e.g., `v1.2.3`, `v2.0.0-alpha.1`) to the repository automatically triggers a full release for that tag.
    *   The tag `v2.1.0-main` is explicitly excluded from this trigger to avoid conflict with the main branch rolling release.

3.  **Main Branch Push (`push: branches: main`):**
    *   Every push to the `main` branch triggers a "rolling" release.
    *   This release is tagged with a consistent version defined by the `MAIN_VERSION` environment variable in the workflow (e.g., `v2.1.0-main`).
    *   Existing releases and tags for this `MAIN_VERSION` are deleted before the new one is created to ensure it always points to the latest `main` commit.

4.  **Pull Request to Main (`pull_request: branches: main`):**
    *   When a pull request is opened or updated against the `main` branch, a snapshot release is triggered.
    *   This is for testing and validation purposes.
    *   The version for these snapshot releases is typically formatted as `<latest-git-tag>-pr.<PR_NUMBER>-<commit-sha>` (e.g., `v1.2.2-pr.123-abcdef1`).
    *   Snapshot releases use the `--snapshot` flag in GoReleaser and skip certain steps like publishing final release artifacts.

## Workflow Jobs

The `.github/workflows/release.yaml` workflow consists of the following main jobs:

1.  **`setup`**:
    *   Checks out the repository code (including full history for version determination).
    *   Determines the `version` string based on the event that triggered the workflow (manual dispatch, tag push, main branch push, or pull request). It validates the format of manually provided versions.
    *   Sets appropriate `goreleaser_args` based on the trigger type (e.g., `--clean`, `--skip=validate` for manual/main, `--snapshot` for PRs). These arguments are passed to the `make release` command.

2.  **`helm`**:
    *   Packages the Helm charts for `kgateway` and `kgateway-crds` using `make package-kgateway-charts`. The version from the `setup` job is used.
    *   For pull request events, it only lints the Helm charts (`make lint-kgateway-charts`).
    *   For other events (manual, tag, main push), it logs into the GitHub Container Registry (`ghcr.io`) and pushes the packaged Helm chart OCI artifacts to `oci://${{ env.IMAGE_REGISTRY }}/charts`.

3.  **`goreleaser`**:
    *   This is the core job for building and releasing artifacts. It depends on the `setup` and `helm` jobs.
    *   Prepares the Go environment and sets up Docker Buildx and QEMU for multi-arch builds.
    *   If triggered by a push to `main`, it attempts to delete any pre-existing GitHub release and tag for the `MAIN_VERSION` to ensure the rolling release is always fresh.
    *   For non-PR events, it logs into `ghcr.io` using a Docker login action.
    *   Executes `make release`. This Makefile target invokes GoReleaser, passing necessary environment variables like `GITHUB_TOKEN`, `VERSION`, `IMAGE_REGISTRY`, `GORELEASER_ARGS`, and `GORELEASER_CURRENT_TAG`.
    *   **GoReleaser's Role (`.goreleaser.yaml`):**
        *   Runs pre-hooks: `go mod tidy`, `go mod download`.
        *   Builds binaries for `controller` (kgateway), `sds`, and `envoyinit` for Linux (amd64 and arm64).
        *   Builds Docker images for each binary and architecture using specified Dockerfiles (e.g., `cmd/kgateway/Dockerfile`). It also builds a Python `ai_extension` image.
        *   Creates and pushes multi-architecture Docker manifest lists to `${IMAGE_REGISTRY}` (which resolves to `ghcr.io/${{ github.repository_owner }}`). The final images are available under this registry.
        *   Manages GitHub Releases:
            *   Automatically determines if a release is a pre-release based on the version string (e.g., `v1.0.0-alpha.1`).
            *   Replaces existing releases and artifacts if a release with the same tag is re-run.
            *   Generates release notes using a custom `header` and `footer`. The header distinguishes between rolling `main` builds and tagged releases. The footer includes installation instructions for Helm charts and Docker images, referencing a vanity URL (`cr.kgateway.dev/kgateway-dev`).
            *   Changelog generation by GoReleaser itself is disabled, implying the content between the header and footer is auto-generated based on commit messages unless a custom release note is provided.

4.  **`validate`**:
    *   This job runs if the workflow was triggered by a tag push or if the `validate: true` input was provided during a manual dispatch.
    *   It performs end-to-end validation of the released artifacts.
    *   Sets up a KinD (Kubernetes in Docker) cluster.
    *   Installs the just-released `kgateway-crds` and `kgateway` Helm charts from the OCI registry (`oci://${{ env.IMAGE_REGISTRY }}/charts`) using the current release version.
    *   Waits for the `kgateway` deployment to become available in the KinD cluster.
    *   Runs conformance tests using `make conformance` against the deployed version.

## How to Create a New Tagged Release (Manual Process for Maintainers)

1.  **Ensure Pre-Release Checks are Complete:**
    *   The `main` branch (or a release branch if used) is stable and contains all changes intended for the release.
    *   All tests are passing.
    *   Relevant documentation has been updated.
    *   A changelog has been prepared or reviewed.

2.  **Trigger the "Release" Workflow:**
    *   Navigate to the "Actions" tab of the kgateway GitHub repository.
    *   Find and select the "Release" workflow from the list on the left.
    *   Click the "Run workflow" dropdown button.

3.  **Provide Inputs:**
    *   **`version`**: Enter the desired semantic version for the new release (e.g., `v1.2.3`). **Important:** This Git tag should *not* already exist. GoReleaser will create the tag as part of its process.
    *   **`validate`**: You can leave this as `false` (default) or set it to `true`. The validation job will run automatically for new tags anyway as per the `if: startsWith(github.ref, 'refs/tags/')` condition on the `validate` job, but setting it `true` ensures it runs even if there's a logic change there.

4.  **Run the Workflow:**
    *   Click the "Run workflow" button.

5.  **Monitor Execution:**
    *   Track the progress of the workflow through the GitHub Actions UI.
    *   Ensure all jobs (`setup`, `helm`, `goreleaser`, `validate`) complete successfully.

6.  **Verify Release:**
    *   Once the workflow is complete, navigate to the "Releases" page of the repository.
    *   Confirm that a new release with the specified tag (e.g., `v1.2.3`) has been created.
    *   Check the release notes and attached artifacts.
    *   Verify that the Docker images and Helm charts are available in `ghcr.io` (and accessible via the vanity URL `cr.kgateway.dev`).

## Artifact Locations

*   **Docker Images:** `ghcr.io/${{ github.repository_owner }}/<image-name>:<version>` (also available via `cr.kgateway.dev/kgateway-dev/<image-name>:<version>`)
    *   `kgateway`
    *   `sds`
    *   `envoyinit`
    *   `ai_extension`
*   **Helm Charts:** `oci://ghcr.io/${{ github.repository_owner }}/charts/<chart-name>:<version>` (also available via `oci://cr.kgateway.dev/kgateway-dev/charts/<chart-name>:<version>`)
    *   `kgateway`
    *   `kgateway-crds`
*   **GitHub Release:** Contains source code archives, and potentially other artifacts if configured in GoReleaser (e.g., compiled binaries, checksums).

This process ensures that releases are consistent, automated, and include validation steps to maintain quality.