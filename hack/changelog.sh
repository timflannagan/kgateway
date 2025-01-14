#!/usr/bin/env bash
set -o pipefail
set -o errexit
set -o nounset
set -x

LAST_TAG=$(git describe --tags --abbrev=0)
LAST_TAG_DATE=$(git log -1 --format=%aI "$LAST_TAG")
echo "Last tag: $LAST_TAG"
echo "Last tag date: $LAST_TAG_DATE"

# FIXME. Only search for PRs merged after the last tag date.
# FIXME: REPO_NAME argument.
# FIXME: --search "merged:>$LAST_TAG_DATE" \
PR_JSON=$(gh pr list \
  --label release-note \
  --repo github.com/timflannagan/k8sgateway \
  --state open \
  --json number,title,labels,url,mergedAt,body \
  --limit 100)

# If PR_JSON is empty or equals "[]", skip changelog generation
if [[ "$PR_JSON" == "[]" ]]; then
  echo "No PRs found. Skipping changelog generation."
  exit 0
fi

# Initialize empty variables for each section
new_features=""
bug_fixes=""
deprecations=""
breaking_changes=""

# Determine number of PRs
PR_COUNT=$(echo "$PR_JSON" | jq 'length')

# Loop over each PR to categorize and extract release notes
for (( i=0; i<$PR_COUNT; i++ )); do
  title=$(echo "$PR_JSON" | jq -r ".[$i].title")
  url=$(echo "$PR_JSON" | jq -r ".[$i].url")
  labels=$(echo "$PR_JSON" | jq -r ".[$i].labels[].name")
  body=$(echo "$PR_JSON" | jq -r ".[$i].body")

  # Extract release note content from the PR body
  note=$(echo "$body" | awk '/```release-note/{flag=1;next}/```/{flag=0}flag' | xargs)
  if [[ -n "$note" ]]; then
    entry="- $note ([$title]($url))"$'\n'
  else
    entry="- [$title]($url)"$'\n'
  fi

  # Categorize the entry based on kind labels
  for label in $labels; do
    case "$label" in
      "kind/new_feature"|"kind/enhancement")
        new_features+="$entry"
        ;;
      "kind/bug")
        bug_fixes+="$entry"
        ;;
      "kind/deprecation")
        deprecations+="$entry"
        ;;
      "kind/breaking-change")
        breaking_changes+="$entry"
        ;;
    esac
  done
done

# Build the changelog content conditionally using proper newlines
changelog_content=""
if [ -n "$new_features" ]; then
  changelog_content+="## New Features"$'\n\n'"$new_features"$'\n'
fi
if [ -n "$bug_fixes" ]; then
  changelog_content+="## Bug Fixes"$'\n\n'"$bug_fixes"$'\n'
fi
if [ -n "$deprecations" ]; then
  changelog_content+="## Deprecations"$'\n\n'"$deprecations"$'\n'
fi
if [ -n "$breaking_changes" ]; then
  changelog_content+="## Breaking Changes"$'\n\n'"$breaking_changes"$'\n'
fi

# Only create CHANGELOG.md if there's content
if [ -n "$changelog_content" ]; then
  echo "$changelog_content" > CHANGELOG.md
  echo "Generated CHANGELOG.md:"
  cat CHANGELOG.md
else
  echo "No content for CHANGELOG.md. Skipping file creation."
fi
