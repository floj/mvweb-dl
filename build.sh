#!/usr/bin/env bash
set -eu
export GOARCH=arm64 GOOS=linux CGO_ENABLED=0
go build -v -o "mvweb-dl-${GOARCH}"
