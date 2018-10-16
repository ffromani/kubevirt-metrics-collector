#!/usr/bin/bash

set -xe

/usr/sbin/collectd -C /etc/collectd/collectd.conf
/usr/sbin/kube-metrics-collector -U /var/run/collectd.sock /etc/kube-metrics-collector.json
