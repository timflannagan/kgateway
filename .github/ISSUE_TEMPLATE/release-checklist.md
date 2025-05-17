---
name: Release Checklist
about: Checklist for creating a new kgateway release.
title: 'Release Checklist: vx.y.z'
labels: release
assignees: ''

---

# Release Checklist for Version: `vx.y.z`

This checklist outlines the steps for creating a new release of the kgateway project.
Replace `vx.y.z` in the issue title and below with the target release version.

## Pre-Release Preparations

-   [ ] **Confirm Stability:** Ensure the `main` branch is stable and all planned features/fixes for this release are merged.
-   [ ] **Update Dependencies (if needed):** Review and update Go module dependencies (`go.mod`) and any other project dependencies (e.g., Python requirements, base Docker images).
-   [ ] **Documentation Review:** Ensure all relevant documentation (READMEs, `docs/` folder, API docs, user guides) is up-to-date with the changes in this release.
-   [ ] **Changelog Preparation (Manual Review):**
    *   [ ] Gather significant changes, new features, bug fixes, and breaking changes since the last release.
    *   [ ] Draft or review the changelog/release notes. GoReleaser will generate notes based on commit messages between the custom header and footer, but manual review and additions are crucial for clarity.
-   [ ] **Test Thoroughly:**
    *   [ ] Ensure all unit tests pass: `make test`
    *   [ ] Ensure all integration/E2E tests pass locally: `make conformance` (or equivalent comprehensive test suite).
    *   [ ] Perform any necessary manual Quality Assurance (QA) or exploratory testing for critical paths.
-   [ ] **Check `MAIN_VERSION` in `.github/workflows/release.yaml` (if applicable):** If the `MAIN_VERSION` environment variable (used for rolling releases from the `main` branch) needs to be updated for the next development cycle *after* this tag is cut, prepare a separate PR for that change.

## Release Execution (Triggering the Workflow)

-   [ ] **Verify Permissions:** Ensure you have maintainer/admin permissions to trigger the "Release" workflow and that the workflow has permissions to push to GitHub Packages/`ghcr.io` and create/update GitHub Releases.
-   [ ] **Navigate to GitHub Actions:** Go to the "Actions" tab of the repository.
-   [ ] **Select "Release" Workflow:** Find the "Release" workflow in the list.
-   [ ] **Click "Run workflow" dropdown.**
-   [ ] **Input Parameters for the Workflow Dispatch:**
    *   **`version`**: Enter the exact semantic version for this release (e.g., `v1.2.3`, `v2.0.0-rc.1`). **This Git tag must NOT exist yet.** GoReleaser will create it.
    *   **`validate`**: (Optional) Set to `true` to explicitly run the `validate` job. Note: The `validate` job runs by default for new tags anyway due to the `if: startsWith(github.ref, 'refs/tags/')` condition, but this input provides an override.
-   [ ] **Initiate Workflow:** Click the "Run workflow" button.

## Monitor Workflow Execution

-   [ ] **Monitor `setup` job:** Verify the correct `version` and `goreleaser_args` are determined.
-   [ ] **Monitor `helm` job:** Verify Helm charts (`kgateway` and `kgateway-crds`) are packaged and pushed successfully to `oci://ghcr.io/.../charts`.
-   [ ] **Monitor `goreleaser` job:**
    *   [ ] Verify Go binaries (`kgateway`, `sds`, `envoyinit`) are built for all target architectures.
    *   [ ] Verify Docker images (including `ai_extension`) are built and pushed to `ghcr.io` for all architectures.
    *   [ ] Verify multi-arch Docker manifests are created and pushed.
    *   [ ] Verify GitHub Release is created with the correct tag (`vx.y.z`).
    *   [ ] Review the generated release notes in the GitHub Release (header/footer from `.goreleaser.yaml` + auto-generated content). Ensure they are accurate.
-   [ ] **Monitor `validate` job:** (If triggered/applicable for tags)
    *   [ ] Verify KinD cluster setup is successful.
    *   [ ] Verify Helm charts are installed correctly from the OCI registry.
    *   [ ] Verify `kgateway` deployment becomes available.
    *   [ ] Verify conformance tests (`make conformance`) pass against the deployed release.

## Post-Release Verification & Tasks

-   [ ] **Verify GitHub Release:**
    *   [ ] Double-check the newly created release on the GitHub "Releases" page.
    *   [ ] Ensure the tag (`vx.y.z`) is correct and points to the intended commit.
    *   [ ] Confirm all expected artifacts (e.g., source code archives, checksums if configured by GoReleaser) are present.
    *   [ ] Perform a final review of the release notes. Edit directly on GitHub if minor tweaks are needed.
-   [ ] **Verify Container Images:**
    *   [ ] Check `ghcr.io` (e.g., `ghcr.io/OWNER/kgateway:vx.y.z`) for the new image tags. Confirm accessibility via the vanity URL `cr.kgateway.dev`.
    *   [ ] Verify multi-arch manifests are correct (e.g., using `docker manifest inspect ghcr.io/OWNER/kgateway:vx.y.z`).
-   [ ] **Verify Helm Charts:**
    *   [ ] Check the OCI registry (`oci://ghcr.io/OWNER/charts/kgateway:vx.y.z` and `oci://ghcr.io/OWNER/charts/kgateway-crds:vx.y.z`) for the new chart versions. Confirm accessibility via the vanity URL `oci://cr.kgateway.dev`.
    *   [ ] (Recommended) Perform a test installation of the new Helm chart from the OCI registry into a test environment.
-   [ ] **Announce Release:**
    *   [ ] Communicate the new release to users and contributors (e.g., project Slack, mailing list, social media, community forums).
    *   [ ] Include links to the GitHub release, changelog, and relevant documentation.
-   [ ] **Update `MAIN_VERSION` (if a PR was prepared):** Merge the PR to update the `MAIN_VERSION` environment variable in `.github/workflows/release.yaml` for the next development cycle on the `main` branch.
-   [ ] **Housekeeping:**
    *   [ ] Close this release issue.
    *   [ ] If your project uses release branches, create or update branches for the next development cycle as per your branching strategy (e.g., create `release/vx.y+1.z` or merge `release/vx.y.z` back to `main` if applicable).

## Troubleshooting Notes

*   If the workflow fails, examine the GitHub Actions logs for the specific job and step that errored.
*   Common failure points:
    *   Incorrect `GITHUB_TOKEN` permissions.
    *   Errors in `.goreleaser.yaml` configuration.
    *   Docker build issues (e.g., Dockerfile errors, base image not found, architecture-specific build problems).
    *   Helm packaging or OCI push failures.
    *   KinD cluster setup issues or test failures in the `validate` job.
    *   Git tag already exists (ensure the `version` provided to `workflow_dispatch` is new).
*   If GoReleaser fails to create the GitHub release or tag, check its logs for details. It might be due to authentication issues or conflicts.