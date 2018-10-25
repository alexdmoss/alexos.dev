#!/usr/bin/env bash
#
# Builds docker images and publishes to Google Container Registry
#
# Notes:
# - As part of CD pipeline, version should be passed in by CI tool, but
#   for portability/simplicity for now, we set some export variables

# source global variables
. ./config/vars.sh

if [[ -z ${GCP_PROJECT_NAME} ]];  then echo "[ERROR] GCP_PROJECT_NAME not set, aborting."; exit 1; fi
if [[ -z ${APP_NAME} ]];  then echo "[ERROR] NGINX_IMAGE_NAME not set, aborting."; exit 1; fi

# get latest build info pushed to GCR (assumes NGINX & PHP version are linked)
LATEST_TAG=$(gcloud container images list-tags eu.gcr.io/${GCP_PROJECT_NAME}/${APP_NAME}-nginx --sort-by="~timestamp" --limit=1 --format='value(tags)')
if [[ $(echo $LATEST_TAG | grep -c ",") -gt 0 ]]; then
  LATEST_TAG=$(echo $LATEST_TAG | awk -F, '{print $2}');
fi

# might be first build
if [[ -z $LATEST_TAG ]];then
  NEW_TAG=0.1
else
  NEW_TAG=$(echo $LATEST_TAG| awk -F. -v OFS=. 'NF==1{print ++$NF}; NF>1{if(length($NF+1)>length($NF))$(NF-1)++; $NF=sprintf("%0*d", length($NF), ($NF+1)%(10^length($NF))); print}')
fi

set -x

NGINX_BUILD_IMAGE=eu.gcr.io/${GCP_PROJECT_NAME}/${APP_NAME}-nginx:${NEW_TAG}

find . -name '*.DS_Store' -exec rm {} \;

# build NGINX image
docker build -t ${NGINX_BUILD_IMAGE} -f ./Dockerfile.nginx .

# push to GCR - assumes command line already authenticated
docker push ${NGINX_BUILD_IMAGE}
