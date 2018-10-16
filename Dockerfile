FROM fedora:28

MAINTAINER "Francesco Romani" <fromani@redhat.com>
ENV container docker

RUN \
  dnf install -y \
    collectd collectd-virt collectd-write_prometheus && \
  dnf clean all

RUN \
  dnf install -y \
    procps-ng curl less && \
  dnf clean all

RUN mkdir -p /etc/collectd
# safety fallback
COPY cluster/collectd.conf /etc/collectd/collectd.conf

COPY cluster/kube-metrics-collector.json /etc/kube-metrics-collector.json
COPY cmd/kube-metrics-collector/kube-metrics-collector /usr/sbin/kube-metrics-collector
COPY cluster/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
