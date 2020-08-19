#!/usr/bin/env bash

if [[ ! -v TAG ]]; then
  TAG="nightly"
  readonly TAG
fi

function build_release() {
  # config/ contains the manifests
  ko resolve --tags "$1" -f config/ > release.yaml
  ARTIFACTS_TO_PUBLISH="release.yaml"
}

build_release ${TAG}
