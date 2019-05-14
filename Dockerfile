FROM registry.access.redhat.com/ubi8/ubi-minimal

MAINTAINER "Francesco Romani" <fromani@redhat.com>
ENV container docker

RUN mkdir -p /etc/kubevirt-metrics-collector
COPY cluster/kubevirt-metrics-collector.json /etc/kubevirt-metrics-collector/config.json
COPY cmd/kubevirt-metrics-collector/kubevirt-metrics-collector /usr/sbin/kubevirt-metrics-collector

ENTRYPOINT [ "/usr/sbin/kubevirt-metrics-collector", "/etc/kubevirt-metrics-collector/config.json" ]
