#!/usr/bin/env bash

set -euo pipefail

# create the package name with version

if [ ! -f buildpack.toml ]; then
  echo "missing buildpack.toml!"
  exit 1
fi

version="$(cat buildpack.toml | yj -t | jq -r .buildpack.version)"
echo "::set-output name=version::${version}"
echo "Selected version ${version} from buildpack.toml"

package="${INPUT_PACKAGE}:${version}"

if [ ! -f package.toml ]; then
  echo "[buildpack]\nuri = \".\"" > package.toml
fi

if [[ -n "${INPUT_PUBLISH+x}" ]]; then
  pack package-buildpack \
    "${package}" \
    --config package.toml \
    --publish
else
  pack package-buildpack \
    "${package}" \
    --config package.toml
fi

echo "::set-output name=digest::$(crane digest "${package}")"
