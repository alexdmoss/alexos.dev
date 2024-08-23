#!/usr/bin/env bash
set -euoE pipefail

hugo server -p 1314 -wDs src/
