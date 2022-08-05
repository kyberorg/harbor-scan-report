#!/bin/sh -l

# Harbor instance
export HARBOR_HOST=${1}
export HARBOR_PROTO=${7}
export HARBOR_PORT=${8}

# Harbor Access
export HARBOR_ROBOT=${2}
export HARBOR_TOKEN=${3}

# Harbor image
export IMAGE=${4}

# GitHub Settings (if set comment will be written)
export GITHUB_TOKEN=${6}
export GITHUB_EVENT_TYPE=${5}
export GITHUB_ISSUE_COMMENT_URL=${9}
export GITHUB_PR_COMMENT_URL=${10}

/hsr