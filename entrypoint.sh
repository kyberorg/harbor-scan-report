#!/bin/sh -l

# Harbor instance
export HARBOR_HOST=${1}
export HARBOR_PROTO=${8}
export HARBOR_PORT=${9}

# Harbor Access
export HARBOR_ROBOT=${2}
export HARBOR_TOKEN=${3}

# Harbor image
export IMAGE=${4}
export DIGEST=${14}
export MAX_ALLOWED_SEVERITY=${5}

# GitHub Settings (if set comment will be written)
export GITHUB_TOKEN=${7}
export GITHUB_URL=${6}

# Comment customization
export COMMENT_TITLE=${10}
export COMMENT_MODE=${11}

# Timing options
export TIMEOUT=${12}
export CHECK_INTERVAL=${13}

# Report
export REPORT_SORT_BY=${15}

# Run it!
/hsr