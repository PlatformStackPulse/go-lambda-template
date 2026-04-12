#!/usr/bin/env bash

# Regenerates CHANGELOG.md from git history and Conventional Commits.

set -euo pipefail

if ! command -v git-chglog >/dev/null 2>&1; then
  echo "git-chglog is not installed."
  echo "Install with: go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest"
  exit 1
fi

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

if git tag --list | grep -q .; then
  git-chglog --config .chglog/config.yml --template .chglog/CHANGELOG.tpl.md --output CHANGELOG.md
else
  git-chglog --config .chglog/config.yml --template .chglog/CHANGELOG.tpl.md --next-tag v0.1.0 --output CHANGELOG.md
fi

echo "CHANGELOG.md updated"
