#!/bin/bash
set -e
mkdir -p build
export CGO_ENABLED=0
go tool dist list | grep -v wasm | while IFS=/ read -r goos goarch; do
  echo "Building for $goos/$goarch..."
  GOOS="$goos" GOARCH="$goarch" go build -v -o "build/simple-dashboardd-$goos-$goarch" ./cmd/simple-dashboardd || echo "Failed to build for $goos/$goarch"
done
