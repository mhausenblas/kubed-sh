#!/usr/bin/env bash

set -e

if ! command -v jq >/dev/null 2>&1; then
    echo "Please install jq before continuing"
    exit 1
fi

if ! command -v eksctl >/dev/null 2>&1; then
    echo "Please install eksctl before continuing"
    exit 1
fi

echo "Provisioning EKS cluster using Fargate"

# set region

# eksctl create cluster -f fg-cluster-spec.yaml

# check if cluster if available

# kubectl create namespace serverless

# kubed-sh
