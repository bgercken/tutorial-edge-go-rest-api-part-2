#!/usr/bin/env bash
TAG="bgercken/comments-api:latest"

docker build -t $TAG .

[[ $? -eq 0 ]] && docker push $TAG
