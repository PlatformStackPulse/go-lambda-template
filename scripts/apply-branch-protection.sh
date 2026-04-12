#!/usr/bin/env bash

# Apply branch protection to main branch via GitHub REST API.
# Requirements:
#   - curl
#   - jq
#   - GITHUB_TOKEN env var with repo admin permissions
# Optional:
#   - GITHUB_OWNER (default: PlatformStackPulse)
#   - GITHUB_REPO (default: go-template)
#   - BRANCH (default: main)

set -euo pipefail

if ! command -v curl >/dev/null 2>&1; then
  echo "curl is required"
  exit 1
fi

if ! command -v jq >/dev/null 2>&1; then
  echo "jq is required"
  exit 1
fi

if [[ -z "${GITHUB_TOKEN:-}" ]]; then
  echo "GITHUB_TOKEN is required"
  echo "Example: export GITHUB_TOKEN=ghp_xxx"
  exit 1
fi

GITHUB_OWNER="${GITHUB_OWNER:-PlatformStackPulse}"
GITHUB_REPO="${GITHUB_REPO:-go-template}"
BRANCH="${BRANCH:-main}"

API_URL="https://api.github.com/repos/${GITHUB_OWNER}/${GITHUB_REPO}/branches/${BRANCH}/protection"

echo "Applying branch protection to ${GITHUB_OWNER}/${GITHUB_REPO}:${BRANCH}"

PAYLOAD='{
  "required_status_checks": {
    "strict": true,
    "contexts": [
      "CI Pipeline / Lint & Format Check",
      "CI Pipeline / Test (1.21)",
      "CI Pipeline / Test (1.22)",
      "CI Pipeline / Security Scans",
      "CI Pipeline / Commit Lint",
      "CI Pipeline / Build",
      "CodeQL Analysis / Analyze"
    ]
  },
  "enforce_admins": true,
  "required_pull_request_reviews": {
    "dismiss_stale_reviews": true,
    "required_approving_review_count": 1,
    "require_code_owner_reviews": false,
    "require_last_push_approval": false
  },
  "restrictions": null,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "required_conversation_resolution": true,
  "lock_branch": false,
  "allow_fork_syncing": true
}'

HTTP_CODE=$(curl -sS -o /tmp/branch-protection-response.json -w "%{http_code}" \
  -X PUT \
  -H "Accept: application/vnd.github+json" \
  -H "Authorization: Bearer ${GITHUB_TOKEN}" \
  -H "X-GitHub-Api-Version: 2022-11-28" \
  "${API_URL}" \
  -d "${PAYLOAD}")

if [[ "${HTTP_CODE}" != "200" ]]; then
  echo "Failed to apply branch protection (HTTP ${HTTP_CODE})"
  cat /tmp/branch-protection-response.json
  exit 1
fi

echo "Branch protection applied successfully."
echo "Configured required checks:"
jq -r '.required_status_checks.contexts[]' /tmp/branch-protection-response.json
