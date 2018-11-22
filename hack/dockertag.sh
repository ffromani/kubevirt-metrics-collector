#!/bin/bash

TAG="${1:-devel}"

docker build -t fromanirh/kube-metrics-collector:$TAG .
