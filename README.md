# kubevirt-metrics-collector

`kubevirt-metrics-collector` is a [kubevirt](http://kubevirt.io) addon that watches the processes used by `virt-launcher` to run VMs, and report their usage consumption (CPU, memory).

`kubevirt-metrics-collector` *augments* the metrics reported by kubevirt itself, does not replace it.

Use `kubevirt-metrics-collector` if you want to have *extra* insights about how your VMs are performing, for troubleshooting, monitoring, tuning.

[![Go Report Card](https://goreportcard.com/badge/github.com/fromanirh/kubevirt-metrics-collector)](https://goreportcard.com/report/github.com/fromanirh/kubevirt-metrics-collector)

## License

Apache v2

## Dependencies

* [gopsutil](https://github.com/shirou/gopsutil)
* [kubernetes APIs](https://github.com/kubernetes/kubernetes)


## Installation

This project can be deployed in a kubevirt cluster to report metrics about processes running inside PODs.
You may use this to monitor the resource consumption of these infrastructure processes for VM-based workloads.
We assume that the cluster is running [prometheus-operator](https://github.com/coreos/prometheus-operator/blob/master/Documentation/user-guides/getting-started.md) to manage the monitoring using prometheus,
and [kubevirt](https://github.com/kubevirt/kubevirt/releases/tag/v0.12.0) >= 0.12.0, which includes itself better integration with prometheus operator.

### Deploy on KubeVirt running on [Kubernetes](https://kubernetes.io/)

We expect kubernetes >= 1.12
Move into the manifests directory:
```bash
cd cluster/manifests
```

TODO

#### Deploy on KubeVirt running on [OKD](https://www.okd.io/)

We expect okd >= 3.11.
Move into the manifests directory:
```bash
cd cluster/manifests
```

Create the service using the provided manifest:
```bash
oc create -n openshift-monitoring -f vmi-service.yaml
```

Create the config map:
```bash
oc create -n openshift-monitoring -f config-map.yaml
```

You need to set up a new accounts and a new securityContextConstraints, both achieved doing:
```bash
oc create -n openshift-monitoring -f okd-account-scc.yaml
```

Now add the permissions to the `securityContextConstraints` created on the step right above.
First add permissions for the 'hostPath' volume:
```bash
oc patch scc scc-hostpath -p '{"allowHostDirVolumePlugin": true}'
```

Then make sure the `hostPID` setting is allowed:
```bash
oc patch scc scc-hostpath -p '{"allowHostPID": true}'
```

Now you can deploy the collector
```
oc create -n openshift-monitoring -f okd-daemonset.yaml
```

#### HACK: Deploy on KubeVirt running on [OKD](https://www.okd.io/) faking node-exporter

We expect okd >= 3.11.
the `kubevirt-metrics-collector` can plug in the existing `node-exporter` service and can
be used transparently with the vanilla infrastructure, with no extra configuration.
If you are unsure if this applies to your case, please ignore this and just use the previous,
supported installatin method.

Move into the manifests directory:
```bash
cd cluster/manifests
```

Create the config map:
```bash
oc create -n openshift-monitoring -f fake-node-exporter/config-map.yaml
```

You need to set up a new accounts and a new securityContextConstraints, both achieved doing:
```bash
oc create -n openshift-monitoring -f okd-account-scc.yaml
```

Now add the permissions to the `securityContextConstraints` created on the step right above.
First add permissions for the 'hostPath' volume:
```bash
oc patch scc scc-hostpath -p '{"allowHostDirVolumePlugin": true}'
```

Then make sure the `hostPID` setting is allowed:
```bash
oc patch scc scc-hostpath -p '{"allowHostPID": true}'
```

Now you can deploy the collector
```
oc create -n openshift-monitoring -f fake-node-exporter/okd-daemonset.yaml
```

This is it. You should soon see the metrics in your prometheus servers.

### Fix namespace mismatch (optional)
`kubevirt-metrics-collector` uses a deployment in the `kube-system` namespace. VM pods usually run in the `default` namespace.
This may make the prometheus server unable to scrape the metrics endpoint that `kubevirt-metrics-collector` added.
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

## Exposed metrics

`kubevirt-metrics-collector` exposes metrics about the resource consumption of the infrastructural processes which make it possible
to run VMs inside PODs. The intended usage of those metrics is for development, troubleshooting, optimization. In the day-to-day
scenario, you should not worry about those metrics.

### Metrics naming conventions

`kubevirt-metrics-collector` use a naming scheme as similar as possible to the [node_exporter](https://github.com/prometheus/node_exporter);
if the purpose of a metric allows that, we try to use an identical name; otherwise, we aim to be as close as possible.

The server exposes the metrics in the `kubevirt` namespace, in the `pod_infra` subsystem. See the [full definition here](https://github.com/fromanirh/kubevirt-metrics-collector/blob/master/pkg/monitoring/processes/prometheus/prometheus.go#L42)

The metrics are versioned following [these recommendations](https://www.robustperception.io/exposing-the-software-version-to-prometheus).

### Metrics listing

You can learn about all the metrics exposed by `kubevirt-metrics-collector` without deploying in your cluster, using the `-M` flag of the server.
Example:
```bash
$ ./cmd/kubevirt-metrics-collector/kubevirt-metrics-collector -M 2>&1 | grep -v '^#' | grep kube
kubevirt_info{branch="master",goversion="go1.10.5",kubeversion="0.9.1",revision="566d93d",version="1"} 1
kubevirt_pod_infra_cpu_seconds_total{domain="init",host="localhost",process="kubevirt-metrics-collector",type="system"} 0
kubevirt_pod_infra_cpu_seconds_total{domain="init",host="localhost",process="kubevirt-metrics-collector",type="user"} 0
kubevirt_pod_infra_memory_amount_bytes{domain="init",host="localhost",process="kubevirt-metrics-collector",type="dirty"} 5.6410112e+07
kubevirt_pod_infra_memory_amount_bytes{domain="init",host="localhost",process="kubevirt-metrics-collector",type="resident"} 1.2054528e+07
kubevirt_pod_infra_memory_amount_bytes{domain="init",host="localhost",process="kubevirt-metrics-collector",type="shared"} 1.009664e+07
kubevirt_pod_infra_memory_amount_bytes{domain="init",host="localhost",process="kubevirt-metrics-collector",type="virtual"} 4.80759808e+08
```

## Notes about integration with kubernetes/kubevirt

Please be aware that in order to resolve the PIDs to meaningful VM domain names, procwatch **needs to access the CRI socket on the host**.
This is equivalent of exposing the docker socket inside the container. This may or may not be an issue on your cluster setup.

If you disable the CRI socket access, procwatch will just report the PIDs of the monitored processes.


## Caveats & Gotchas

content pending
