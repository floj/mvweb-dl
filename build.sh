#!/usr/bin/env bash
set -eu
export GOARCH=amd64
export GOOS=linux
export CGO_ENABLED=0
go build -v -o "mvweb-dl-${GOARCH}"
