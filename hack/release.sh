#!/bin/bash

set -e

if [ -z "$1" ]; then
	echo "usage: $0 <tag>"
	exit 1
fi
if [ -z "${GITHUB_TOKEN}" ] || [ -z "${GITHUB_USER}" ]; then
	echo "make sure to set GITHUB_TOKEN and GITHUB_USER env vars"
	exit 2
fi

TAG="$1"  #TODO: validate tag is vX.Y.Z
BRANCH=$(git rev-parse --abbrev-ref HEAD)

./hack/build/build.sh ${TAG}
if [ -d _out ]; then
	rm -rf _out;
fi
mkdir -p _out
cp cmd/kubevirt-metrics-collector/kubevirt-metrics-collector _out/kubevirt-metrics-collector-${TAG}-linux-amd64
git add cmd/kubevirt-metrics-collector/kubevirt-metrics-collector && git ci -s -m "binaries: rebuild for tag ${TAG}"
git tag -a -m "kubevirt-metrics-collector ${TAG}" ${TAG}
git push origin --tags ${BRANCH}
if  which github-release 2> /dev/null; then
	github-release release -t ${TAG} -r kubevirt-metrics-collector
	github-release upload -t ${TAG} -r kubevirt-metrics-collector \
		-n kubevirt-metrics-collector-${TAG}-linux-amd64 \
		-f _out/kubevirt-metrics-collector-${TAG}-linux-amd64
fi
