#!/usr/bin/env bash
set -euo pipefail

function help() {
  echo -e "Usage: go <command>"
  echo -e
  echo -e "    help                     Print this help"
  echo -e "    run                      Runs site locally on for fast-feedback / exploratory testing"
  echo -e "    local_build              Builds the site (HTML + docker image) locally only - no push"
  echo -e "    build                    Builds the site (HTML + docker image) and pushes to registry"
  echo -e "    deploy                   Deploys site onto hosts. Assumes pre-requisites are done"
  echo -e "    smoke                    Runs smoke tests"
  echo -e 
  exit 0
}

function run() {
    pushd $(dirname $BASH_SOURCE[0]) >/dev/null
    _run_hugo server -p 1314 -wDs src/
    popd >/dev/null
}


function local_build() {

    _console_msg "Building site locally ..."

    _assert_variables_set GCP_PROJECT_NAME APP_NAME DOMAIN

    mkdir -p "www/"

    pushd "www/" > /dev/null
    rm -rf ./*
    popd >/dev/null 
    
    pushd "src/" >/dev/null
    _run_hugo --source .
    popd >/dev/null

    _build-test

    pushd $(dirname $BASH_SOURCE[0]) >/dev/null

    _console_msg "Baking docker image ..."

    IMAGE_NAME=eu.gcr.io/${GCP_PROJECT_NAME}/${APP_NAME}

    docker build --tag ${APP_NAME}:latest .

    _local-test ${APP_NAME}:latest

    _console_msg "Built locally:
          - run with:  docker run -p 32080:32080 ${APP_NAME}:latest
          - test with: curl -H \"Host: ${DOMAIN}\" -s http://localhost:32080/" INFO true

    popd >/dev/null 

}

function build() {

    _console_msg "Building site ..."

    pushd $(dirname $BASH_SOURCE[0]) > /dev/null

    _assert_variables_set GCP_PROJECT_NAME APP_NAME CI_COMMIT_SHA

    if [[ ${CI_SERVER:-} == "yes" ]]; then
        _assert_variables_set GOOGLE_CREDENTIALS
        echo "${GOOGLE_CREDENTIALS}" | gcloud auth activate-service-account --key-file -
        trap "gcloud auth revoke --verbosity=error" EXIT

        wget --no-verbose -O hugo.tar.gz https://github.com/gohugoio/hugo/releases/download/v0.79.0/hugo_extended_0.79.0_Linux-64bit.tar.gz && \
        tar zxf hugo.tar.gz && \
        mv ./hugo /usr/local/bin/ && \
        rm hugo.tar.gz

    fi
    
    mkdir -p "www/"

    pushd "www/" > /dev/null
    rm -rf ./*
    popd >/dev/null 
    
    pushd "src/" >/dev/null
    _run_hugo --source .
    popd >/dev/null

    _build-test

    pushd $(dirname $BASH_SOURCE[0]) >/dev/null

    _console_msg "Baking docker image ..."

    IMAGE_NAME=eu.gcr.io/${GCP_PROJECT_NAME}/${APP_NAME}

    gcloud auth configure-docker --quiet
    docker pull ${IMAGE_NAME}:latest || true
    docker build --cache-from ${IMAGE_NAME}:latest --tag ${APP_NAME}:latest .

    _local-test ${APP_NAME}:latest

    _console_msg "Pushing image to registry ..."

    docker tag ${APP_NAME}:latest ${IMAGE_NAME}:${CI_COMMIT_SHA}
    docker tag ${APP_NAME}:latest ${IMAGE_NAME}:latest

    docker push ${IMAGE_NAME}:${CI_COMMIT_SHA}
    docker push ${IMAGE_NAME}:latest

    popd >/dev/null 
    
    _console_msg "Build complete" INFO true 

}

function deploy() {

    _assert_variables_set GCP_PROJECT_NAME GCP_REGION CLUSTER_NAME APP_NAME DOMAIN NAMESPACE CI_COMMIT_SHA

    _console_msg "Deploying app ..." INFO true

    export IMAGE_NAME=eu.gcr.io/${GCP_PROJECT_NAME}/${APP_NAME}

    pushd "k8s/" >/dev/null

    if [[ ${CI_SERVER:-} == "yes" ]]; then

        echo "${GOOGLE_CREDENTIALS}" | gcloud auth activate-service-account --key-file -
        trap "gcloud auth revoke --verbosity=error" EXIT

        gcloud config set project ${GCP_PROJECT_NAME}
        gcloud config set compute/region ${GCP_REGION}
        gcloud config set container/cluster ${CLUSTER_NAME}
        gcloud container clusters get-credentials ${CLUSTER_NAME} --region ${GCP_REGION} --project ${GCP_PROJECT_NAME}

    fi

    cat *.yaml | envsubst | kubectl apply -n ${NAMESPACE} -f -

    _console_msg "Deployment complete" INFO true

    popd >/dev/null

}

function _build-test() {

    local error=0

    _console_msg "Running local build tests ..."

    markdown_files=$(find content -type f -name '*.md')

    for md_file in ${markdown_files}; do
    
        html_file="index.html"
        html_path=$(dirname ${md_file} | sed 's#^content#www#')

        if [[ $(basename ${md_file}) == "_index.md" ]]; then
            html_file="index.html"
        elif [[ $(echo ${md_file} | grep -c '/posts/') -gt 0 ]]; then
            if [[ $(grep -c 'draft: true' ${md_file}) -gt 0 ]]; then
                _console_msg "${md_file} in draft - SKIPPING"
            else
                # permalinks mean we need to extract the date to know its destination
                publish_date=$(grep 'date: ' ${md_file})
                publish_year=$(echo ${publish_date} | awk -F- '{print $1}' | awk -F': ' '{print $2}')
                publish_month=$(echo ${publish_date} | awk -F- '{print $2}')
                publish_day=$(echo ${publish_date} | awk -F- '{print $3}' | awk -FT '{print $1}')
                html_file=${publish_year}/${publish_month}/${publish_day}/$(basename ${md_file} | sed 's#\.md$#/index.html#')
            fi
        else
            html_file=$(basename ${md_file} | sed 's#\.md$#/index.html#')
        fi

        test_file=$(echo ${html_path}/${html_file} | sed 's#/posts##' | awk '{print tolower($0)}')
        if [[ ! -f "${test_file}" ]]; then
            error=1
            _console_msg "Expected HTML file was missing. Markdown: ${md_file} should be assembled into: ${test_file}"
        fi

    done

    if [[ ! -f "www/tags/index.html" ]]; then
        error=1
        _console_msg "Tags index is missing"
    fi
    if [[ ! -f "www/categories/index.html" ]]; then
        error=1
        _console_msg "Categories index is missing"
    fi
    if [[ ! -f "www/index.json" ]]; then
        error=1
        _console_msg "index.json (for Search) is missing"
    fi
    if [[ ! -f "www/sitemap.xml" ]]; then
        error=1
        _console_msg "sitemap.xml missing"
    fi
    if [[ ! -f "www/robots.txt" ]]; then
        error=1
        _console_msg "robots.txt missing"
    fi
    if [[ ! -f "www/404.html" ]]; then
        error=1
        _console_msg "404.html file missing"
    fi

    if [[ "${error:-}" != "0" ]]; then
        _console_msg "Tests FAILED - see messages above for for detail" ERROR
        exit 1
    else
        _console_msg "All build tests passed!"
    fi

}

function _local-test() {

    local error=0
    local image=${1:-}

    if [[ ${CI_SERVER:-} == "yes" ]]; then
        local_hostname=docker
    else
        local_hostname=localhost
    fi

    _console_msg "Running local docker image tests ..."

    _assert_variables_set APP_NAME DOMAIN

    docker run -d --name ${APP_NAME} -p 32080:32080 ${image}
    trap "docker rm -f ${APP_NAME} >/dev/null 2>&1 || true" EXIT

    # wow really, does it actually need this? /sigh
    sleep 2

    _smoke_test ${DOMAIN} http://${local_hostname}:32080/ "Recent Posts"
    _smoke_test ${DOMAIN} http://${local_hostname}:32080/about/ "A little bit of info about me"
    _smoke_test ${DOMAIN} http://${local_hostname}:32080/contact/ "Send Message"
    _smoke_test ${DOMAIN} http://${local_hostname}:32080/posts/ "Previous Page"
    _smoke_test ${DOMAIN} http://${local_hostname}:32080/categories/ "/categories/cloud"
    _smoke_test ${DOMAIN} http://${local_hostname}:32080/tags/ "/tags/google"
    _smoke_test ${DOMAIN} http://${local_hostname}:32080/2019/02/23/a-year-in-google-cloud/ "This time last year"

    if [[ "${error:-}" != "0" ]]; then
        _console_msg "Tests FAILED - see messages above for for detail" ERROR
        exit 1
    else
        _console_msg "All local tests passed!"
    fi

}

function smoke() {

    local error=0

    _assert_variables_set DOMAIN

    _console_msg "Checking HTTP status code for https://${DOMAIN}/ ..."
    
    _smoke_test ${DOMAIN} http://${DOMAIN}/ "Recent Posts"
    _smoke_test ${DOMAIN} https://${DOMAIN}/about/ "A little bit of info about me"
    _smoke_test ${DOMAIN} https://${DOMAIN}/contact/ "Send Message"
    _smoke_test ${DOMAIN} https://${DOMAIN}/posts/ "Previous Page"
    _smoke_test ${DOMAIN} https://${DOMAIN}/categories/ "/categories/cloud"
    _smoke_test ${DOMAIN} https://${DOMAIN}/tags/ "/tags/google"
    _smoke_test ${DOMAIN} https://${DOMAIN}/2019/02/23/a-year-in-google-cloud/ "This time last year"

    if [[ "${error:-}" != "0" ]]; then
        _console_msg "Tests FAILED - see messages above for for detail" ERROR
        exit 1
    else
        _console_msg "All local tests passed!"
    fi

}

function _smoke_test() {
    local domain=$1
    local url=$2
    local match=$3
    output=$(curl -H "Host: ${domain}" -s -k -L -w "\nHTTP-%{http_code}" ${url} || true)
    if [[ $(echo ${output} | grep -c "HTTP-200") -eq 0 ]]; then
        _console_msg "Test FAILED - ${url} - non-200 return code" ERROR
        error=1
    fi
    if [[ $(echo ${output} | grep -c "${match}") -eq 0 ]]; then 
        _console_msg "Test FAILED - ${url} - missing phrase" ERROR
        error=1
    else
        _console_msg "Test PASSED - ${url}" INFO
    fi
}

function _run_hugo() {
    case "$OSTYPE" in
        # darwin*)  HUGO_BIN='hugo';;
        # linux*)   HUGO_BIN='hugo';;
        win*)     HUGO_BIN='C:/hugo/hugo';;
        cygwin*)  HUGO_BIN='C:/hugo/hugo';;
        msys*)    HUGO_BIN='C:/hugo/hugo';;
        *)        HUGO_BIN='hugo';;
    esac
    (${HUGO_BIN} "$@") 
}

function _assert_variables_set() {
  local error=0
  local varname
  for varname in "$@"; do
    if [[ -z "${!varname-}" ]]; then
      echo "${varname} must be set" >&2
      error=1
    fi
  done
  if [[ ${error} = 1 ]]; then
    exit 1
  fi
}

function _console_msg() {
  local msg=${1}
  local level=${2:-}
  local ts=${3:-}
  if [[ -z ${level} ]]; then level=INFO; fi
  if [[ -n ${ts} ]]; then ts=" [$(date +"%Y-%m-%d %H:%M")]"; fi

  if [[ ${level} == "ERROR" ]] || [[ ${level} == "CRIT" ]] || [[ ${level} == "FATAL" ]]; then
    (echo >&2 "-> [${level}]${ts} ${msg}")
  else 
    (echo "-> [${level}]${ts} ${msg}")
  fi
}

function ctrl_c() {
    if [ ! -z ${PID:-} ]; then
        kill ${PID}
    fi
    exit 1
}

trap ctrl_c INT

if [[ ${1:-} =~ ^(help|run|local_build|build|deploy|test|smoke)$ ]]; then
  COMMAND=${1}
  shift
  $COMMAND "$@"
else
  help
  exit 1
fi
