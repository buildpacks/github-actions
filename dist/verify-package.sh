#!/usr/bin/env bash

set -euo pipefail

expected_namespace="${INPUT_NAMESPACE}"
expected_name="${INPUT_NAME}"
expected_version="${INPUT_VERSION}"
expected_addr="${INPUT_ADDRESS}"
expected_id="${expected_namespace}/${expected_name}"
expected_escaped_id="${expected_namespace}_${expected_name}"

crane export "${expected_addr}" - \
  | tar xOf - "/cnb/buildpacks/${expected_escaped_id}/${expected_version}/buildpack.toml" \
  | yj -tj \
  > /tmp/buildpack.json

actual_id=$(cat /tmp/buildpack.json | jq -r .buildpack.id)
actual_version=$(cat /tmp/buildpack.json | jq -r .buildpack.version)

if [ "${expected_id}" != "${actual_id}" ]; then
  echo "invalid id in buildpackage: expected '${expected_id}' found '${actual_id}'"
  exit 1
elif [ "${expected_version}" != "${actual_version}" ]; then
  echo "invalid version in buildpackage: expected '${expected_version}' found '${actual_version}'"
  exit 1
fi

echo "successfully verified ${expected_addr}"
exit 0
