#!/usr/bin/env bash
set -euo pipefail

GCP_PROJECT_NAME=jl-digital-docs 
TFSTATE_BUCKET_NAME=${GCP_PROJECT_NAME}-tfstate
DOMAIN=docs.jl-digital.net 
REGION=europe-west2

source functions.sh 

function help() {
  echo -e "Usage: go <command>"
  echo -e
  echo -e "    help                     Print this help"
  echo -e "    run                      Runs site locally on for fast-feedback / exploratory testing"
  echo -e "    page <chapter/page.md>   Add a new page called <chapter>/<page.md>"
  echo -e "    chapter <chapter>        Add a new chapter called <chapter>/"
  echo -e "    build                    Builds the site (creates static HTML) - either locally or as part of CI"
  echo -e "    deploy                   Deploys site into Google AppEngine. Assumes pre-requisites are done"
  echo -e "    test                     Runs local tests to ensure build created required files"
  echo -e 
  exit 0
}

function run() {
    pushd $(dirname $BASH_SOURCE[0]) >/dev/null
    run_hugo server -p 1314 -wDs src/
    popd >/dev/null
}

function page() {
    if [[ ! ${1:-} ]]; then console_msg "You must specify the new page name - <chapter>/<page.md>" ERROR; exit 1; fi
    pushd "src/" >/dev/null
    run_hugo new ${1}
    popd >/dev/null
}

function chapter() {
    if [[ ! ${1:-} ]]; then console_msg "You must specify the new chapter name" ERROR; exit 1; fi
    pushd "src/" >/dev/null
    run_hugo new -k chapter ${1}/_index.md
    popd >/dev/null
}

function build() {

    console_msg "Building site ..."

    mkdir -p "www/"

    pushd "www/" > /dev/null
    rm -rf ./*
    popd >/dev/null 
    
    pushd "src/" >/dev/null
    run_hugo -s .
    popd >/dev/null

    test

    console_msg "Build complete" INFO true 

}

function deploy() {

    assert_variables_set GCP_PROJECT_NAME SITE DOMAIN

    pushd "app/" >/dev/null

    echo "${GOOGLE_CREDENTIALS}" | gcloud auth activate-service-account --key-file -
    trap "gcloud auth revoke" EXIT

    console_msg "Deploying docs site ..." INFO true

    cat app-template.yaml | sed 's/${SITE}/'${SITE}'/g' > app.yaml
    gcloud app deploy --project ${GCP_PROJECT_NAME} --quiet

    console_msg "Deployment complete" INFO true

    popd >/dev/null

    # [AM] This test is not going to be so great for brand new sites - seems to take about 10-15 mins for new SSL certs to be issued and usable
    # If that's the case, maybe query deploy status with `gcloud app` instead - version check perhaps?

    console_msg "Checking HTTP status code for https://${SITE}.${DOMAIN}/ ..."

    response_code=$(curl -k -L -o /dev/null -w "%{http_code}" https://${SITE}.${DOMAIN}/)

    if [[ ${response_code:0:1} == "4" ]] || [[ ${response_code:0:1} == "5" ]]; then
        console_msg "Test FAILED - HTTP response code was ${response_code}" ERROR
        exit 1
    else 
        console_msg "Test PASSED - HTTP response code was ${response_code}"
    fi

}

function test() {

    local error=0

    console_msg "Running unit tests ..."

    markdown_files=$(find content -type f -name '*.md')

    for md_file in ${markdown_files}; do
        html_path=$(dirname ${md_file} | sed 's#^content#www#')

        if [[ $(basename ${md_file}) == "_index.md" ]]; then
            html_file="index.html"
        else
            html_file=$(basename ${md_file} | sed 's#\.md$#/index.html#')
        fi

        test_file=$(echo ${html_path}/${html_file} | awk '{print tolower($0)}')
        if [[ ! -f "${test_file}" ]]; then
            error=1
            console_msg "Expected HTML file was missing. Markdown: ${md_file} should be assembled into: ${test_file}"
        fi

    done

    if [[ "${error}" != "0" ]]; then
        console_msg "Tests FAILED - see messages above for for detail" ERROR true
        exit 1
    else
        console_msg "All tests passed!"
    fi

}

# We assume on a JLP laptop where we can't set path, it is installed in C:/hugo/ as per instructions. Otherwise, we assume its in the user's path
function run_hugo() {

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

function ctrl_c() {
    if [ ! -z ${PID:-} ]; then
        kill ${PID}
    fi
    exit 1
}

trap ctrl_c INT

if [[ ${1:-} =~ ^(help|run|page|chapter|build|deploy|test)$ ]]; then
  COMMAND=${1}
  shift
  $COMMAND "$@"
else
  help
  exit 1
fi
