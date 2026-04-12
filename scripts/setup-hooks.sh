#!/usr/bin/env bash

# Setup Git hooks for development
# Run this script during initial setup: ./scripts/setup-hooks.sh

set -e

HOOKS_DIR=".git/hooks"
SCRIPT_DIR="./scripts"

# Create hooks directory if it doesn't exist
mkdir -p "$HOOKS_DIR"

# Pre-commit hook
if [ -f "$SCRIPT_DIR/pre-commit" ]; then
    cp "$SCRIPT_DIR/pre-commit" "$HOOKS_DIR/pre-commit"
    chmod +x "$HOOKS_DIR/pre-commit"
    echo "✅ Pre-commit hook installed"
else
    echo "⚠️  Pre-commit hook not found"
fi

echo "✅ Git hooks setup complete"
echo ""
echo "Configured hooks:"
echo "  - pre-commit: Validates conventional commit format"
