#!/usr/bin/env bash

# Build script for production builds
# Usage: ./scripts/build.sh [VERSION] [OUTPUT_DIR]

set -e

VERSION=${1:-dev}
OUTPUT_DIR=${2:-./dist}
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION=$(go version | awk '{print $3}')

LD_FLAGS="-X github.com/PlatformStackPulse/go-template/pkg/version.Version=$VERSION \
          -X github.com/PlatformStackPulse/go-template/pkg/version.Commit=$COMMIT \
          -X github.com/PlatformStackPulse/go-template/pkg/version.BuildTime=$BUILD_TIME \
          -X github.com/PlatformStackPulse/go-template/pkg/version.GoVersion=$GO_VERSION"

mkdir -p "$OUTPUT_DIR"

echo "🔨 Building go-template v$VERSION"
echo "   Commit: $COMMIT"
echo "   Build Time: $BUILD_TIME"
echo "   Go Version: $GO_VERSION"
echo ""

# Build for multiple platforms
PLATFORMS=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
)

for platform in "${PLATFORMS[@]}"; do
  OS=$(echo $platform | cut -d'/' -f1)
  ARCH=$(echo $platform | cut -d'/' -f2)
  OUTPUT_NAME="go-template-$VERSION-$OS-$ARCH"

  if [ "$OS" == "windows" ]; then
    OUTPUT_NAME="$OUTPUT_NAME.exe"
  fi

  echo "Building $OUTPUT_NAME..."

  GOOS=$OS GOARCH=$ARCH go build \
    -ldflags "$LD_FLAGS" \
    -o "$OUTPUT_DIR/$OUTPUT_NAME" \
    cmd/app/main.go

  echo "✅ Built: $OUTPUT_DIR/$OUTPUT_NAME"
done

echo ""
echo "✅ Build complete! Artifacts in $OUTPUT_DIR"
