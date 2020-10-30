#!/usr/bin/env bash

set -euo pipefail

pack_version="${1}"
jq_version="1.6"
crane_version="0.1.4"
yj_version="5.0.0"

mkdir -p "$HOME/bin"

echo "---> Installing jq ${jq_version}"
curl \
  --retry 3 \
  --output "${HOME}/bin/jq" \
  --location \
  --show-error \
  --silent \
  "https://github.com/stedolan/jq/releases/download/jq-${jq_version}/jq-linux64"
chmod +x "${HOME}/bin/jq"

echo "---> Installing crane ${crane_version}"
curl \
  --retry 3 \
  --location \
  --show-error \
  --silent \
  "https://github.com/google/go-containerregistry/releases/download/v${crane_version}/go-containerregistry_Linux_x86_64.tar.gz" \
  | tar -C "${HOME}/bin/" -xzv crane

echo "---> Installing yj ${yj_version}"
curl \
  --retry 3 \
  --output "${HOME}/bin/yj" \
  --location \
  --show-error \
  --silent \
  "https://github.com/sclevine/yj/releases/download/v${yj_version}/yj-linux"
chmod +x "${HOME}/bin/yj"

echo "---> Installing pack ${pack_version}"
curl \
  --retry 3 \
  --location \
  --show-error \
  --silent \
  "https://github.com/buildpacks/pack/releases/download/v${pack_version}/pack-v${pack_version}-linux.tgz" \
  | tar -C "${HOME}/bin/" -xzv pack

echo "PATH=${HOME}/bin:${PATH}" >> $GITHUB_ENV
