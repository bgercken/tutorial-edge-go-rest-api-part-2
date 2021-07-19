#!/usr/bin/env bash

. ./ENV.sh
envsubst < config/deployment.yaml | kubectl apply -f -
