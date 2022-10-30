#!/bin/bash


set -Eeuxo pipefail

project_root="$(cd "$(git -C "$( dirname "${BASH_SOURCE[0]}" )" rev-parse --show-toplevel)" && pwd)"; readonly project_root

git_short_rev="$(git -C "${project_root}" rev-parse --short HEAD)"; readonly git_short_rev

export DOCKER_BUILDKIT=1

exit_cleanup() {
  local -r retval="$?"
  set +Eeu

  if (( retval == 0 )); then
    echo -e "\n*** Success (Build Candidate Images) ***\n"
  else
    echo -e "\n*** Failure (Build Candidate Images) ***\n"
  fi

  exit "${retval}"
}
trap "exit_cleanup" EXIT

build_targets=(
  build
)

build_args=(
  --file "${project_root}"/dockerfile
  --progress plain
)

if [[ -n "${EXTRA_TEST_ARGS:-}" ]]; then
  build_args+=(--build-arg "EXTRA_TEST_ARGS=${EXTRA_TEST_ARGS}")
fi

for target in "${build_targets[@]}"; do
  docker build \
    --target "${target}" \
    --tag "go-ethereum-helpers-${target}:${git_short_rev}" \
    "${build_args[@]}" \
    "${project_root}"
done
