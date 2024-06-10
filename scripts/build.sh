#!/bin/sh
set -eu

IMAGE_TAG=${IMAGE_NAME}:${CI_COMMIT_SHA}-$(echo "${CI_COMMIT_TIMESTAMP}" | sed 's/[:+]/./g')

echo "Building image [${IMAGE_TAG}]"

set -x

# we cannot wrap ECP_IMAGE_LABELS or it sets as one giant label
#shellcheck disable=SC2086
/kaniko/executor \
    --context="dir://${CI_PROJECT_DIR}" \
    --destination="${IMAGE_TAG}" \
    --dockerfile="${CI_PROJECT_DIR}/Dockerfile" \
    --label=alexos-dev.source.pipeline.url=${CI_PIPELINE_URL} \
    --compressed-caching=false \
    --cache=true \
    --use-new-run \
    --cleanup
