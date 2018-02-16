#!/usr/bin/env bash

kubectl delete ds,deploy,svc,po --selector=gen=kubed-sh
