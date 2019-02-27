#!/bin/bash
TAGS="$*"
for TAG in $TAGS; do
	buildah push fromani/kubevirt-metrics-collector:$TAG docker://quay.io/fromani/kubevirt-metrics-collector:$TAG
done
