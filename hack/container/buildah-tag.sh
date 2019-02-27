#!/bin/bash

TAG="${1:-devel}"
XTAGS="${@:2}"  # see https://stackoverflow.com/questions/9057387/process-all-arguments-except-the-first-one-in-a-bash-script

buildah bud -t fromani/kubevirt-metrics-collector:$TAG .
for XTAG in $XTAGS; do
	buildah tag fromani/kubevirt-metrics-collector:$TAG fromani/kubevirt-metrics-collector:$XTAG
done
