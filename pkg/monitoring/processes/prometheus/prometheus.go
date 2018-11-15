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

package prometheus

import (
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
	"github.com/shirou/gopsutil/process"

	verinfo "github.com/fromanirh/kube-metrics-collector/internal/pkg/version"
)

var labels = []string{
	"host",   //  On which host is the domain running?
	"domain", // Which domain the process belongs to?
	"name",   // What's the process?
	"type",   // what are we measuring?
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
	cpuTimes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kubevirt",
			Subsystem: "process_cpu",
			Name:      "seconds_total",
			Help:      "Amount spent, seconds.",
		},
		labels,
	)
	memoryAmount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kubevirt",
			Subsystem: "process_memory",
			Name:      "amount_bytes",
			Help:      "Amount of memory, bytes.",
		},
		labels,
	)
)

func init() {
	prometheus.MustRegister(version)
	prometheus.MustRegister(cpuTimes)
	prometheus.MustRegister(memoryAmount)

	version.Set(1)
}

// fill the metrics with data about itself
func autoFillMetrics() error {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return err
	}
	err = UpdateMetricsCPU("localhost", "init", proc)
	if err != nil {
		return err
	}
	err = UpdateMetricsMemory("localhost", "init", proc)
	if err != nil {
		return err
	}
	return nil
}

func DumpMetrics(w io.Writer) error {
	err := autoFillMetrics()
	if err != nil {
		return err
	}

	mfs, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return err
	}

	for _, mf := range mfs {
		if _, err := expfmt.MetricFamilyToText(w, mf); err != nil {
			return err
		}
	}
	return nil
}

func UpdateMetricsCPU(host, domain string, proc *process.Process) error {
	name, err := extractProcName(proc)
	if err != nil {
		return err
	}

	times, err := proc.Times()
	if err != nil {
		return err
	}

	cpuTimes.With(prometheus.Labels{"host": host, "domain": domain, "name": name, "type": "user"}).Set(times.User)
	cpuTimes.With(prometheus.Labels{"host": host, "domain": domain, "name": name, "type": "system"}).Set(times.System)

	return nil
}

func UpdateMetricsMemory(host, domain string, proc *process.Process) error {
	name, err := extractProcName(proc)
	if err != nil {
		return err
	}

	memInfo, err := proc.MemoryInfoEx()
	if err != nil {
		return err
	}

	memoryAmount.With(prometheus.Labels{"host": host, "domain": domain, "name": name, "type": "virtual"}).Set(float64(memInfo.VMS))
	memoryAmount.With(prometheus.Labels{"host": host, "domain": domain, "name": name, "type": "resident"}).Set(float64(memInfo.RSS))
	memoryAmount.With(prometheus.Labels{"host": host, "domain": domain, "name": name, "type": "shared"}).Set(float64(memInfo.Shared))
	memoryAmount.With(prometheus.Labels{"host": host, "domain": domain, "name": name, "type": "dirty"}).Set(float64(memInfo.Dirty))
	return nil
}

func extractProcName(proc *process.Process) (string, error) {
	exe, err := proc.Exe()
	if err != nil {
		return "", err
	}
	return filepath.Base(exe), nil
}
