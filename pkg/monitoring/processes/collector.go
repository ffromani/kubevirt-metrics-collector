/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2018 Red Hat, Inc.
 *
 */

package processes

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/process"

	"github.com/fromanirh/kubevirt-metrics-collector/pkg/procscanner"

	verinfo "github.com/fromanirh/kubevirt-metrics-collector/internal/pkg/version"
)

var labels = []string{
	"host",    // On which host is the domain running?
	"domain",  // Which domain the process belongs to?
	"process", // What's the process?
	"type",    // what are we measuring?
}

var (
	// see https://www.robustperception.io/exposing-the-software-version-to-prometheus
	version = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "kubevirt",
			Name:      "info",
			Help:      "Version information",
			ConstLabels: prometheus.Labels{
				"branch":      verinfo.BRANCH,
				"goversion":   runtime.Version(),
				"revision":    verinfo.REVISION,
				"kubeversion": verinfo.KUBEVERSION,
				"version":     "1",
			},
		},
	)
	cpuTimesDesc = prometheus.NewDesc(
		"kubevirt_pod_infra_cpu_seconds_total",
		"CPU time spent, seconds.",
		labels,
		nil,
	)
	memoryAmountDesc = prometheus.NewDesc(
		"kubevirt_pod_infra_memory_amount_bytes",
		"Memory amount, bytes.",
		labels,
		nil,
	)
)

type Collector struct {
	conf *Config
	mon  Monitor
}

func NewSelfCollector() (*Collector, error) {
	mon, err := NewSelfMonitor()
	if err != nil {
		// how so?!
		return nil, err
	}

	return &Collector{
		conf: NewConfig(),
		mon:  mon,
	}, nil
}

func NewCollectorFromConf(conf *Config) (*Collector, error) {
	scanner := procscanner.ProcScanner{
		Targets: conf.Targets,
	}

	finder, err := NewCRIPodFinder(conf.CRIEndPoint, DefaultTimeout, scanner)
	if err != nil {
		return nil, err
	}

	mon, err := NewDomainMonitor(finder)
	if err != nil {
		return nil, err
	}

	return &Collector{
		conf: conf,
		mon:  mon,
	}, nil
}

func (co Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(co, ch)
}

// Note that Collect could be called concurrently
func (co Collector) Collect(ch chan<- prometheus.Metric) {
	pods, err := co.mon.Update()
	if err != nil {
		// TODO: log
		return
	}

	updated := 0
	for podName, podInfo := range pods {
		for _, proc := range podInfo.Procs {
			err = co.collectCPU(ch, podName, proc)
			if err != nil {
				log.Printf("failed to update CPU for pod %v: %v", podName, err)
				continue
			}

			err = co.collectMemory(ch, podName, proc)
			if err != nil {
				log.Printf("failed to update Memory for pod %v: %v", podName, err)
				continue
			}
			updated++
		}
	}
	log.Printf("updated metrics for %v pods", updated)
}

func (co *Collector) collectCPU(ch chan<- prometheus.Metric, domain string, proc *process.Process) error {
	process, err := extractProcName(proc)
	if err != nil {
		return err
	}

	times, err := proc.Times()
	if err != nil {
		return err
	}

	m, err := prometheus.NewConstMetric(
		cpuTimesDesc, prometheus.GaugeValue,
		times.User,
		co.conf.Hostname, domain, process, "user",
	)
	if err != nil {
		return err
	}
	ch <- m

	m, err = prometheus.NewConstMetric(
		cpuTimesDesc, prometheus.GaugeValue,
		times.System,
		co.conf.Hostname, domain, process, "system",
	)
	if err != nil {
		return err
	}
	ch <- m

	return nil
}

func (co *Collector) collectMemory(ch chan<- prometheus.Metric, domain string, proc *process.Process) error {
	process, err := extractProcName(proc)
	if err != nil {
		return err
	}

	memInfo, err := proc.MemoryInfoEx()
	if err != nil {
		return err
	}

	m, err := prometheus.NewConstMetric(
		memoryAmountDesc, prometheus.GaugeValue,
		float64(memInfo.VMS),
		co.conf.Hostname, domain, process, "virtual",
	)
	if err != nil {
		return err
	}
	ch <- m

	m, err = prometheus.NewConstMetric(
		memoryAmountDesc, prometheus.GaugeValue,
		float64(memInfo.RSS),
		co.conf.Hostname, domain, process, "resident",
	)
	if err != nil {
		return err
	}
	ch <- m

	m, err = prometheus.NewConstMetric(
		memoryAmountDesc, prometheus.GaugeValue,
		float64(memInfo.Shared),
		co.conf.Hostname, domain, process, "shared",
	)
	if err != nil {
		return err
	}
	ch <- m

	m, err = prometheus.NewConstMetric(
		memoryAmountDesc, prometheus.GaugeValue,
		float64(memInfo.Dirty),
		co.conf.Hostname, domain, process, "dirty",
	)
	if err != nil {
		return err
	}
	ch <- m

	return nil
}

func extractProcName(proc *process.Process) (string, error) {
	cmdline, err := proc.CmdlineSlice()
	if err != nil || len(cmdline) < 1 {
		return "", err
	}
	return filepath.Base(cmdline[0]), nil
}

func init() {
	prometheus.MustRegister(version)

	version.Set(1)
}
