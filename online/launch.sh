#!/usr/bin/env bash

set -o errexit
set -o errtrace
set -o nounset
set -o pipefail

NAMESPACE="${1:-default}"

# set permissions 
kubectl -n $NAMESPACE create sa kubed-sh 
kubectl -n $NAMESPACE create rolebinding givekdsuperpower \
        --clusterrole=cluster-admin \
        --serviceaccount=$NAMESPACE:kubed-sh \
        --namespace=$NAMESPACE

# launch kubed-sh
kubectl -n $NAMESPACE run kubed-sh \
        --image=quay.io/mhausenblas/kubed-sh:0.5.2 \
        --serviceaccount=kubed-sh

sleep 20

# make accessible in local environment
kubectl  -n $NAMESPACE port-forward deploy/kubed-sh 8888:8080