#!/bin/sh
set -ex
cd "$(dirname "$0")"
GOBIN=$HOME/.terraform.d/plugins GO111MODULE=on go get github.com/yi-jiayu/terraform-provider-telegram
terraform init
