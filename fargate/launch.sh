#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

if ! command -v jq >/dev/null 2>&1; then
    echo "Please install jq before continuing"
    exit 1
fi

if ! command -v eksctl >/dev/null 2>&1; then
    echo "Please install eksctl before continuing"
    exit 1
fi

TARGET_REGION=${1:-eu-west-1} 

echo "Provisioning EKS on Fargate cluster in $TARGET_REGION"

tmpdir=$(mktemp -d)
cat <<EOF >> ${tmpdir}/fg-cluster-spec.yaml
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: kubed-sh
  region: $TARGET_REGION

iam:
  withOIDC: true

fargateProfiles:
  - name: defaultfp
    selectors:
      - namespace: serverless

cloudWatch:
  clusterLogging:
    enableTypes: ["*"]
EOF
eksctl create cluster -f ${tmpdir}/fg-cluster-spec.yaml

# check if cluster if available
# eksctl get cluster kubed-sh -o json | jq -r '.[0].Status'

# kubectl create namespace serverless

# kubed-sh
