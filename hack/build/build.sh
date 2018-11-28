#!/bin/sh

VERSIONDIR="internal/pkg/version"
VERSIONFILE="${VERSIONDIR}/version.go"

mkdir -p ${VERSIONDIR} && ./hack/build/genver.sh > ${VERSIONFILE}
cd cmd/kubevirt-metrics-collector && go build -v .
