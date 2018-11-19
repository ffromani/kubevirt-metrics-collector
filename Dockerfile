FROM fedora:28

MAINTAINER "Francesco Romani" <fromani@redhat.com>
ENV container docker

RUN \
  dnf install -y \
    procps-ng curl less && \
  dnf clean all

RUN mkdir -p /etc/kube-metrics-collector
COPY cluster/kube-metrics-collector.json /etc/kube-metrics-collector/config.json
COPY cmd/kube-metrics-collector/kube-metrics-collector /usr/sbin/kube-metrics-collector
COPY cluster/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
