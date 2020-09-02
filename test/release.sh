#!/usr/bin/env bash

set -o errexit

# (tom24d) donotpush to org repo

if [[ ! -v TAG ]]; then
  TAG="nightly"
  readonly TAG
elif [[ ${TAG} =~ refs/tags/([v|0-9|.]+) ]]; then
  TAG=${BASH_REMATCH[1]}
  readonly TAG
fi

function build_release() {
  # Update release labels if this is a tagged release
  if [[ -n "${TAG}" ]]; then
    echo "Tagged release, updating release labels to contrib.eventing.knative.dev/release: \"${TAG}\""
    LABEL_YAML_CMD=(sed -e "s|contrib.eventing.knative.dev/release: devel|contrib.eventing.knative.dev/release: \"${TAG}\"|")
  else
    echo "Untagged release, will NOT update release labels"
    LABEL_YAML_CMD=(cat)
  fi

  # config/ contains the manifests
  ko resolve --tags "$1" -f config/ | "${LABEL_YAML_CMD[@]}" > release.yaml
  ARTIFACTS_TO_PUBLISH="release.yaml"
}

build_release ${TAG}
