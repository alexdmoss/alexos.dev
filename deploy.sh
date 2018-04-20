#!/usr/bin/env bash
#
# Deploys latest image to GKE by applying Kubernetes deployment manifest,
# after automatically updating with latest image tag in GCR.

# source global variables
. ./config/vars.sh

# check required variables set
if [[ -z ${GCP_PROJECT_NAME} ]];  then echo "[ERROR] GCP_PROJECT_NAME not set, aborting.";  exit 1; fi
if [[ -z ${APP_NAME} ]];  then echo "[ERROR] APP_NAME not set, aborting.";  exit 1; fi
if [[ -z ${NAMESPACE} ]];  then echo "[ERROR] NAMESPACE not set, and required when installing. Aborting.";  exit 1; fi


NGINX_BUILD_IMAGE=eu.gcr.io/${GCP_PROJECT_NAME}/${APP_NAME}-nginx
PHP_BUILD_IMAGE=eu.gcr.io/${GCP_PROJECT_NAME}/${APP_NAME}-app

# get latest build info pushed to GCR (assumes NGINX & PHP version are linked)
LATEST_TAG=$(gcloud container images list-tags ${PHP_BUILD_IMAGE} --sort-by="~timestamp" --limit=1 --format='value(tags)')
if [[ $(echo $LATEST_TAG | grep -c ",") -gt 0 ]]; then LATEST_TAG=$(echo $LATEST_TAG | awk -F, '{print $2}'); fi

set -x

if [[ $1 == "--install" ]]; then
  cat ./k8s/00-create-namespace.yml | sed 's#${NAMESPACE}#'${NAMESPACE}'#g' | kubectl apply -f -
fi

cat ./k8s/${APP_NAME}-nginx.yml | sed 's#${NAMESPACE}#'${NAMESPACE}'#g' | sed 's#${APP_NAME}#'${APP_NAME}'#g' | sed 's#${NGINX_BUILD_IMAGE}#'${NGINX_BUILD_IMAGE}'#g' | sed 's#${IMAGE_VERSION}#'${LATEST_TAG}'#g' | kubectl apply -f -

cat ./k8s/${APP_NAME}-app.yml | sed 's#${NAMESPACE}#'${NAMESPACE}'#g' | sed 's#${APP_NAME}#'${APP_NAME}'#g' | sed 's#${PHP_BUILD_IMAGE}#'${PHP_BUILD_IMAGE}'#g' | sed 's#${IMAGE_VERSION}#'${LATEST_TAG}'#g' | kubectl apply -f -

if [[ $1 == "--install" ]]; then
  cat ./k8s/99-create-ingress.yml | sed 's#${NAMESPACE}#'${NAMESPACE}'#g' | sed 's#${APP_NAME}#'${APP_NAME}'#g' | kubectl apply -f -
fi
