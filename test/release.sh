source vendor/knative.dev/test-infra/scripts/release.sh

function build_release() {
  # config/ contains the manifests
  ko resolve ${KO_FLAGS} -f config/ > release.yaml
  ARTIFACTS_TO_PUBLISH="release.yaml"
}

main "$@"