#!/usr/bin/env bash
# Download and summarize the latest failing CI run on the current branch,
# printing just the bits worth reading: failed test names, error messages,
# and the visible error text from the DOM snapshots.
#
# Usage: scripts/inspect-ci-failure.sh [run-id]
#   run-id defaults to the latest failing run on the current branch.

set -euo pipefail

BRANCH="$(git rev-parse --abbrev-ref HEAD)"
RUN_ID="${1:-$(gh run list --branch "$BRANCH" --status failure --limit 1 --json databaseId --jq '.[0].databaseId')}"

if [[ -z "$RUN_ID" || "$RUN_ID" == "null" ]]; then
  echo "No failed runs found on branch $BRANCH."
  exit 0
fi

OUT="$(mktemp -d)"
trap 'rm -rf "$OUT"' EXIT

echo "=== Run $RUN_ID on $BRANCH ==="
echo

# Pull just the failed-step log lines (test names + playwright errors)
gh run view "$RUN_ID" --log-failed 2>/dev/null \
  | grep -E '✘|Error|> [0-9]+ \|' \
  | head -40
echo

# Download playwright-report artifact
if ! gh run download "$RUN_ID" --name playwright-report --dir "$OUT" 2>/dev/null; then
  echo "(no playwright-report artifact uploaded for this run)"
  exit 0
fi

echo "=== Visible error text per failed scenario ==="
for ctx in "$OUT"/test-results/*/error-context.md; do
  [[ -f "$ctx" ]] || continue
  name=$(basename "$(dirname "$ctx")")
  echo
  echo "--- $name ---"
  # The DOM snapshot's form-error paragraph often has the real error code
  grep -E 'paragraph.*ref=e[0-9]+.*:' "$ctx" | head -5
  echo
  # Error details block
  sed -n '/# Error details/,/```$/p' "$ctx" | head -10
done

echo
echo "Artifact dir (removed on exit): $OUT"
