#!/usr/bin/env bash

# Build script for Lambda artifacts.
# Usage: ./scripts/build.sh [VERSION] [OUTPUT_DIR] [ARCH]

set -euo pipefail

VERSION=${1:-dev}
OUTPUT_DIR=${2:-./dist}
ARCH=${3:-arm64}
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
GO_VERSION=$(go version | awk '{print $3}')

LD_FLAGS="-X github.com/PlatformStackPulse/go-lambda-template/pkg/version.Version=$VERSION \
          -X github.com/PlatformStackPulse/go-lambda-template/pkg/version.Commit=$COMMIT \
          -X github.com/PlatformStackPulse/go-lambda-template/pkg/version.BuildTime=$BUILD_TIME \
          -X github.com/PlatformStackPulse/go-lambda-template/pkg/version.GoVersion=$GO_VERSION"

mkdir -p "$OUTPUT_DIR"
rm -f "$OUTPUT_DIR/bootstrap" "$OUTPUT_DIR/lambda.zip"

echo "Building go-lambda-template v$VERSION"
echo "   Commit: $COMMIT"
echo "   Build Time: $BUILD_TIME"
echo "   Go Version: $GO_VERSION"
echo "   Architecture: $ARCH"
echo ""

CGO_ENABLED=0 GOOS=linux GOARCH="$ARCH" go build \
  -ldflags "$LD_FLAGS" \
  -o "$OUTPUT_DIR/bootstrap" \
  ./cmd/lambda

(
  cd "$OUTPUT_DIR"
  zip -q lambda.zip bootstrap
)

echo "Artifacts created in $OUTPUT_DIR"
