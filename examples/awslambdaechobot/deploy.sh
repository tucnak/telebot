#!/bin/sh
set -ex
cd "$(dirname "$0")"
CGO_ENABLED=0 go build -ldflags="-s -w" -v
terraform apply "$@"
