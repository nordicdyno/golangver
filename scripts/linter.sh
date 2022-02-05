#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

set -x
golangci-lint run ./...
