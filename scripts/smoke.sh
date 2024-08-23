#!/usr/bin/env bash
set -euoE pipefail

function smoke() {

    local error=0

    _console_msg "Checking HTTP status codes for https://${DOMAIN}/ ..."
    
    _smoke_test "${DOMAIN}" https://"${DOMAIN}"/ "Recent Posts" "Recent"
    _smoke_test "${DOMAIN}" https://"${DOMAIN}"/about/ "A little bit of info about me" "About"
    _smoke_test "${DOMAIN}" https://"${DOMAIN}"/contact/ "Send Message" "Contact"
    _smoke_test "${DOMAIN}" https://"${DOMAIN}"/posts/ "Previous Page" "Posts"
    _smoke_test "${DOMAIN}" https://"${DOMAIN}"/tags/ "/tags/google" "Tags"
    _smoke_test "${DOMAIN}" https://"${DOMAIN}"/2019/02/23/a-year-in-google-cloud/ "This time last year" "GCP-Blog"

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
    local explanation=$4

    output=$(curl -H "Host: ${domain}" -s -k -L "${url}" || true)

    if [[ $(echo "${output}" | grep -c "${match}") -eq 0 ]]; then 
        _console_msg "Test $explanation FAILED - ${url} - missing phrase" ERROR
        error=1
    else
        _console_msg "Test $explanation PASSED - ${url}" INFO
    fi

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

smoke "${@:-}"
