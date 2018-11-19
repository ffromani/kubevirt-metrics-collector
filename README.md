# kube-metrics-collector

kube-metrics-collector is a [kubevirt](http://kubevirt.io) addon that watches the processes used by `virt-launcher` to run VMs, and report their usage consumption (CPU, memory).

## License

Apache v2

## Dependencies

* [gopsutil](https://github.com/shirou/gopsutil)
* [kubernetes APIs](https://github.com/kubernetes/kubernetes)


## Installation: kubernetes/kubevirt cluster

This project can be deployed in a kubevirt cluster to report metrics about processes running inside PODs.
You may use this to monitor the resource consumption of these infrastructure processes for VM-based workloads.
We assume that the cluster is running [prometheus-operator](https://github.com/coreos/prometheus-operator/blob/master/Documentation/user-guides/getting-started.md) to manage the monitoring using prometheus,
and [kubevirt](https://github.com/kubevirt/kubevirt/releases/tag/v0.9.1) >= 0.9.1, which includes itself better integration with prometheus operator.

### Deploy a new service to plug in the configuration KubeVirt uses to interact with prometheus-operator:
Just use the provided manifest:
```
kubectl create -f cluster/kube-service-vmi.yaml
```

### Deploy the server itself into the cluster:
Set PLATFORM to either "k8s" or "ocp" and
```
kubectl create -f cluster/kube-metrics-collector-config-map.yaml
kubectl create -f cluster/kube-metrics-collector-$PLATFORM.yaml
```

TODO: add template

### Fix namespace mismatch (optional)
`kube-metrics-collector` uses a deployment in the `kube-system` namespace. VM pods usually run in the `default` namespace.
This may make the prometheus server unable to scrape the metrics endpoint that `kube-metrics-collector` added.
To fix this, deploy your prometheus server in your cluster like this:
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

### Add labels to the namespaces
This step depends on how you are running your cluster. As per kubevirt 0.10.z, if you using the `kube-system` namespace (by default it is so), you need to perform this step.
Here's an example of how it could look like:
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

### Deploy a `Service Monitor`
This is the last step. You need this  to let `prometheus-operator` pickup and add rules for this endpoint.
Just use the provided manifest:
```
kubectl create -f cluster/kube-service-monitor-vmi.yaml
```

## Exposed metrics

`kube-metrics-collector` exposes metrics about the resource consumption of the infrastructural processes which make it possible
to run VMs inside PODs. The intended usage of those metrics is for development, troubleshooting, optimization. In the day-to-day
scenario, you should not worry about those metrics.

### Metrics naming conventions

`kube-metrics collector` use a naming scheme as similar as possible to the [node_exporter](https://github.com/prometheus/node_exporter);
if the purpose of a metric allows that, we try to use an identical name; otherwise, we aim to be as close as possible.

The server exposes the metrics in the `kubevirt` namespace, in the `pod_infra` subsystem. See the [full definition here](https://github.com/fromanirh/kube-metrics-collector/blob/master/pkg/monitoring/processes/prometheus/prometheus.go#L42)

The metrics are versioned following [these recommendations](https://www.robustperception.io/exposing-the-software-version-to-prometheus).

### Metrics listing

You can learn about all the metrics exposed by `kube-metrics-collector` without deploying in your cluster, using the `-M` flag of the server.
Example:
```bash
$ ./cmd/kube-metrics-collector/kube-metrics-collector -M 2>&1 | grep -v '^#' | grep kube
kubevirt_info{branch="master",goversion="go1.10.5",kubeversion="0.9.1",revision="566d93d",version="1"} 1
kubevirt_pod_infra_cpu_seconds_total{domain="init",host="localhost",process="kube-metrics-collector",type="system"} 0
kubevirt_pod_infra_cpu_seconds_total{domain="init",host="localhost",process="kube-metrics-collector",type="user"} 0
kubevirt_pod_infra_memory_amount_bytes{domain="init",host="localhost",process="kube-metrics-collector",type="dirty"} 5.6410112e+07
kubevirt_pod_infra_memory_amount_bytes{domain="init",host="localhost",process="kube-metrics-collector",type="resident"} 1.2054528e+07
kubevirt_pod_infra_memory_amount_bytes{domain="init",host="localhost",process="kube-metrics-collector",type="shared"} 1.009664e+07
kubevirt_pod_infra_memory_amount_bytes{domain="init",host="localhost",process="kube-metrics-collector",type="virtual"} 4.80759808e+08
```

## Notes about integration with kubernetes/kubevirt

Please be aware that in order to resolve the PIDs to meaningful VM domain names, procwatch **needs to access the CRI socket on the host**.
This is equivalent of exposing the docker socket inside the container. This may or may not be an issue on your cluster setup.

If you disable the CRI socket access, procwatch will just report the PIDs of the monitored processes.


## Caveats & Gotchas

content pending
