set -o errexit

export GO111MODULE=on

function upload_test_images() {
  echo ">> Publishing test images"
  # Script needs to be executed from the root directory
  # to pickup .ko.yaml
  cd "$( dirname "$0")/.."
  local image_dir="test/test_images"
  local docker_tag=$1
  local tag_option=""
  if [ -n "${docker_tag}" ]; then
    tag_option="--tags $docker_tag,latest"
  fi

  # ko resolve is being used for the side-effect of publishing images,
  # so the resulting yaml produced is ignored.
  ko resolve ${tag_option} -RBf "${image_dir}" > /dev/null
}

: ${KO_DOCKER_REPO:?"You must set 'KO_DOCKER_REPO', see DEVELOPMENT.md"}

upload_test_images $@