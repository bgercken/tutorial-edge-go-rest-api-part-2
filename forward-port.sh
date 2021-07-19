#!/usr/bin/env bash
kubectl port-forward --address 0.0.0.0 service/comments-api 8080:8080
