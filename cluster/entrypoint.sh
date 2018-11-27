#!/usr/bin/bash

set -xe

/usr/sbin/kubevirt-metrics-collector /etc/kubevirt-metrics-collector/config.json
