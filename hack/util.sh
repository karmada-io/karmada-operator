#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# This script holds common bash variables and utility functions.

KARMADA_GO_PACKAGE="github.com/daocloud/karmada-operator"

# util::create_gopath_tree create the GOPATH tree
# Parameters:
#  - $1: the root path of repo
#  - $2: go path
function util:create_gopath_tree() {
  local repo_root=$1
  local go_path=$2

  local go_pkg_dir="${go_path}/src/${KARMADA_GO_PACKAGE}"
  go_pkg_dir=$(dirname "${go_pkg_dir}")

  mkdir -p "${go_pkg_dir}"

  if [[ ! -e "${go_pkg_dir}" || "$(readlink "${go_pkg_dir}")" != "${repo_root}" ]]; then
    ln -snf "${repo_root}" "${go_pkg_dir}"
  fi
}
