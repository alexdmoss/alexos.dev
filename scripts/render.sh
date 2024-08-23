#!/usr/bin/env bash
set -euoE pipefail

echo "-> [INFO] Building site ..."

pushd "$(dirname "${BASH_SOURCE[0]}")/../" >/dev/null

mkdir -p "www/"

pushd "www/" > /dev/null
rm -rf ./*
popd >/dev/null 

pushd "src/" >/dev/null
hugo --source .
popd >/dev/null

echo "-> [INFO] Build complete"
