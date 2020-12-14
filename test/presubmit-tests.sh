#!/usr/bin/env bash

export GO111MODULE=on

source $(dirname $0)/../vendor/knative.dev/hack/presubmit-tests.sh

# DockerHubSource use the default build, unit and integration test runners.

main $@
