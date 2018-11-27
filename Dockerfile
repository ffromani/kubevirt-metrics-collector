FROM fedora:28

MAINTAINER "Francesco Romani" <fromani@redhat.com>
ENV container docker

RUN \
  dnf install -y \
    procps-ng curl less && \
  dnf clean all

RUN mkdir -p /etc/kubevirt-metrics-collector
COPY cluster/kubevirt-metrics-collector.json /etc/kubevirt-metrics-collector/config.json
COPY cmd/kubevirt-metrics-collector/kubevirt-metrics-collector /usr/sbin/kubevirt-metrics-collector
COPY cluster/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
