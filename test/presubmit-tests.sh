#!/usr/bin/env bash

export GO111MODULE=on

source $(dirname $0)/../vendor/knative.dev/test-infra/scripts/presubmit-tests.sh

# DockerHubSource use the default build, unit and integration test runners.

main $@
