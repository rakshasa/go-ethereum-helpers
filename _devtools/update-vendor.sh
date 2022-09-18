#!/bin/bash

set -eux

project_root="$(cd "$(git -C "$( dirname "${BASH_SOURCE[0]}" )" rev-parse --show-toplevel)" && pwd)"; readonly project_root
git_revision="$(git -C "${project_root}" rev-parse --short HEAD)"; readonly git_revision


readonly build_image="go-ethereum-helpers--update-vendor-${git_revision}"

readonly dependencies=(
  "github.com/ethereum/go-ethereum@v1.10.25"
)

build_dir=$(mktemp -d); readonly build_dir

cleanup() {
  local -r retval="$?"
  set +eu

  docker rm "${container}"
  docker rmi "${build_image}"

  rm -rf "${build_dir}"

  set +x

  if [[ "${success:-no}" == "yes" ]]; then
    echo
    echo "*******************************"
    echo "*** Vendor Update Succeeded ***"
    echo "*******************************"
    echo
  else
    echo
    echo "****************************"
    echo "*** Vendor Update Failed ***"
    echo "****************************"
    echo
  fi

  exit "${retval}"
}
trap cleanup EXIT ERR

cd "${project_root}"

git clone --depth 1 file://"${project_root}" "${build_dir}"

docker build \
  --tag "${build_image}" \
  --target "build-env" \
  - < "${project_root}"/dockerfile

docker run \
  --rm \
  --interactive \
  --volume "${build_dir}":/build/ \
  "${build_image}" \
  /bin/sh - <<EOF
#!/bin/sh
set -eux

cd /build/

rm -rf ./go.mod ./go.sum ./vendor/

go clean -cache
go mod init github.com/rakshasa/go-ethereum-helpers

for dep in ${dependencies[@]}; do
  go get -u -v "\${dep}"
done

go mod tidy -v -compat=1.17
go mod vendor -v

set +x
echo
echo "+----------------------+"
echo "| Vendor Files Created |"
echo "+----------------------+"
echo

EOF

ls "${build_dir}"

rm -rf ./{go.mod,go.sum,vendor}
cp -r "${build_dir}"/{go.mod,go.sum,vendor} ./

success="yes"
