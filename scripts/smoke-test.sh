#!/usr/bin/env bash
#
# Post-deploy smoke test. Each check: URL, expected status, a phrase that only
# appears when the page really rendered, and a description. Redirect checks
# (30x) match the phrase against the Location header instead of the body.
#
# TLS is checked without -k: a load balancer serving the wrong certificate
# must fail the smoke test.

set -uo pipefail

HOST="${SMOKE_HOST:-https://${DOMAIN:?DOMAIN must be set}}"

error=0

function main() {

  echo "-> [INFO] Running smoke tests against ${HOST} ..."

  _check "${HOST}/"                                   "200" "Recent Posts"                   "Home page"
  _check "${HOST}/about/"                             "200" "A little bit of info about me"  "About"
  _check "${HOST}/contact/"                           "200" "Send Message"                   "Contact"
  _check "${HOST}/posts/"                             "200" "Previous Page"                  "Posts"
  _check "${HOST}/tags/"                              "200" "/tags/google"                   "Tags"
  _check "${HOST}/2019/02/23/a-year-in-google-cloud/" "200" "This time last year"            "Blog post"
  _check "${HOST}/index.xml"                          "200" "<item>"                         "RSS feed"
  _check "${HOST}/sitemap.xml"                        "200" "<loc>https://alexos.dev/</loc>" "Sitemap"
  _check "${HOST}/index.json"                         "200" '"uri": "https://alexos.dev/"'   "Search index"
  _check "${HOST}/health"                             "200" "OK"                             "Health endpoint"
  _check "${HOST}/woofwoof/"                          "404" "Four-Oh-Four"                   "404 page"
  _check "${HOST}/about"                              "301" "https://alexos.dev/about/"      "Trailing slash redirect"
  _check "${HOST}/categories/google"                  "301" "https://alexos.dev/tags/google" "Categories redirect"
  _check "http://alexos.dev/"                         "302" "https://alexos.dev"             "Redirect to https"

  echo

  if [[ "${error}" != "0" ]]; then
    echo "-> [ERROR] Tests FAILED - see messages above for detail"
    exit 1
  fi

  echo "-> [INFO] All tests passed!"

}

function _check() {

  local url=$1
  local expected_rc=$2
  local match=$3
  local description=$4
  local output actual_rc rc

  echo

  if [[ ${expected_rc} == 30? ]]; then
    # GET rather than -I: HEAD isn't answered the same way everywhere
    output=$(curl -s -o /dev/null -D - --max-time 10 "${url}")
    rc=$?
    actual_rc=$(echo "${output}" | grep 'HTTP/' | head -n1 | awk '{print $2}')
    match_against=$(echo "${output}" | grep -i '^location:')
  else
    output=$(curl -s -L --max-time 10 -w "\n%{http_code}" "${url}")
    rc=$?
    actual_rc=$(echo "${output}" | tail -n1)
    match_against=${output}
  fi

  if [[ ${rc} -eq 60 ]]; then
    echo "-> [ERROR] Test FAILED - ${description} ::: ${url} - TLS certificate not valid for this host"
    error=1
    return
  fi

  if [[ "${actual_rc}" != "${expected_rc}" ]]; then
    echo "-> [ERROR] Test FAILED - ${description} ::: ${url} - expected ${expected_rc}, received ${actual_rc:-nothing (curl exit ${rc})}"
    error=1
    return
  fi

  if ! grep -qF -- "${match}" <<< "${match_against}"; then
    echo "-> [ERROR] Test FAILED - ${description} ::: ${url} - missing phrase (expected ${match})"
    error=1
    return
  fi

  echo "-> [INFO] Test PASSED - ${description} ::: ${url}"

}

main "$@"
