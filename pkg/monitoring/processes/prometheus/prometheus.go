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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/process"

	"path/filepath"
)

var labels = []string{
	"host",   //  On which host is the domain running?
	"domain", // Which domain the process belongs to?
	"name",   // What's the process?
	"type",   // what are we measuring?
}

var (
	cpuTimes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kubevirt",
			Subsystem: "cpu_times",
			Name:      "time_seconds",
			Help:      "Amount spent, seconds.",
		},
		labels,
	)
	memoryAmount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "kubevirt",
			Subsystem: "process_memory",
			Name:      "amount_kib",
			Help:      "Amount of memory, KiB.",
		},
		labels,
	)
)

func init() {
	prometheus.MustRegister(cpuTimes)
	prometheus.MustRegister(memoryAmount)
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

	memInfo, err := proc.MemoryInfo()
	if err != nil {
		return err
	}

	memoryAmount.With(prometheus.Labels{"host": host, "domain": domain, "name": name, "type": "virtual"}).Set(float64(memInfo.VMS / 1024.))
	memoryAmount.With(prometheus.Labels{"host": host, "domain": domain, "name": name, "type": "resident"}).Set(float64(memInfo.RSS / 1024.))
	return nil
}

func extractProcName(proc *process.Process) (string, error) {
	exe, err := proc.Exe()
	if err != nil {
		return "", err
	}
	return filepath.Base(exe), nil
}
