kube-metrics-collector
======================

kube-metrics-collector is a [kubevirt](http://kubevirt.io) addon that watches the processes used by `virt-launcher` to run VMs, and report their usage consumption (CPU, memory).

License
=======

Apache v2

Dependencies
============

* [gopsutil](https://github.com/shirou/gopsutil)
* [kubernetes APIs](https://github.com/kubernetes/kubernetes)


Installation: kubernetes/kubevirt cluster
=========================================

This project can be deployed in a kubevirt cluster to report metrics about processes running inside PODs.
You may use this to monitor the resource consumption of these infrastructure processes for VM-based workloads.
We assume that the cluster is running [prometheus-operator](https://github.com/coreos/prometheus-operator/blob/master/Documentation/user-guides/getting-started.md) to manage the monitoring using prometheus,
and [kubevirt](https://github.com/kubevirt/kubevirt/releases/tag/v0.9.1) >= 0.9.1, which includes itself better integration with prometheus operator.

0. First deploy a new service to plug in the configuration KubeVirt uses to interact with prometheus-operator:
```
kubectl create -f cluster/kube-service-vmi.yaml
```

1. Now deploy the tool itself into the cluster:
Set PLATFORM to either "k8s" or "ocp" and
```
kubectl create -f cluster/kube-metrics-collector-config-map.yaml
kubectl create -f cluster/kube-metrics-collector-$PLATFORM.yaml
```

TODO: add template

2. procwatch installs a new deployment in the `kube-system` namespace. VM pods usually run in the `default` namespace.
This may make the prometheus server unable to scrape the metrics endpoint.
procwatch added. To fix this, deploy your prometheus server in your cluster like this:
```yaml
apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: prometheus
spec:
  serviceAccountName: prometheus
  serviceMonitorNamespaceSelector:
    matchLabels:
      prometheus.kubevirt.io: ""
  serviceMonitorSelector:
    matchLabels:
      prometheus.kubevirt.io: ""
  resources:
    requests:
      memory: 400Mi

```
Note the usage of `serviceMonitorNamespaceSelector`. [See here for more details](https://github.com/coreos/prometheus-operator/issues/1331)

3. Now, you may need to add labels to the namespaces, like kube-system. Here's an example of how it could look like:
```yaml
apiVersion: v1
kind: Namespace
metadata:
  ...
  creationTimestamp: 2018-09-21T13:53:16Z
  labels:
    prometheus.kubevirt.io: ""
...
```

3. The last step: now you need to deploy a `Service Monitor` to let `prometheus-operator` pickup and add rules for this endpoint:
```
kubectl create -f cluster/kube-service-monitor-vmi.yaml
```

Please check the next sections for Caveats.

Supported metrics
=================

You can learn the metrics exposed by `kube-metrics-collector` without deploying in your cluster, using the `-M` flag of the server. Example:
```bash
$ ./cmd/kube-metrics-collector/kube-metrics-collector -M 2>&1 | grep -v '^#' | grep kube
kubevirt_info{branch="master",goversion="go1.10.4",kubeversion="0.9.1",revision="f2a62fa",version="1"} 1
kubevirt_process_cpu_seconds_total{domain="init",host="localhost",name="kube-metrics-collector",type="system"} 0
kubevirt_process_cpu_seconds_total{domain="init",host="localhost",name="kube-metrics-collector",type="user"} 0
kubevirt_process_memory_amount_bytes{domain="init",host="localhost",name="kube-metrics-collector",type="dirty"} 6.4933888e+07
kubevirt_process_memory_amount_bytes{domain="init",host="localhost",name="kube-metrics-collector",type="resident"} 1.1931648e+07
kubevirt_process_memory_amount_bytes{domain="init",host="localhost",name="kube-metrics-collector",type="shared"} 1.015808e+07
kubevirt_process_memory_amount_bytes{domain="init",host="localhost",name="kube-metrics-collector",type="virtual"} 5.5625728e+08
```


Notes about integration with kubernetes/kubevirt
================================================

Please be aware that in order to resolve the PIDs to meaningful VM domain names, procwatch **needs to access the CRI socket on the host**.
This is equivalent of exposing the docker socket inside the container. This may or may not be an issue on your cluster setup.

If you disable the CRI socket access, procwatch will just report the PIDs of the monitored processes.

