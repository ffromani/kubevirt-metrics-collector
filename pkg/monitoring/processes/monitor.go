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
	"github.com/shirou/gopsutil/process"
)

type MetricsUpdater interface {
	UpdateCPU(domain string, proc *process.Process) error
	UpdateMemory(domain string, proc *process.Process) error
}

type Monitor interface {
	Update() error
}

type DomainMonitor struct {
	podFinder PodFinder
	pods      map[string]*PodInfo
	metrics   MetricsUpdater
}

func NewDomainMonitor(podFinder PodFinder, metricsUpdater MetricsUpdater) (Monitor, error) {
	dm := DomainMonitor{
		podFinder: podFinder,
		pods:      make(map[string]*PodInfo),
		metrics:   metricsUpdater,
	}
	return &dm, nil
}

func (dm *DomainMonitor) Update() error {
	var err error
	if err = dm.refreshPods(); err != nil {
		return err
	}
	if err = dm.updateMetrics(); err != nil {
		return err
	}
	return nil
}

func (dm *DomainMonitor) updateMetrics() error {
	var err error
	for podName, podInfo := range dm.pods {
		for _, proc := range podInfo.Procs {
			err = dm.metrics.UpdateCPU(podName, proc)
			if err != nil {
				continue // TODO
			}

			err = dm.metrics.UpdateMemory(podName, proc)
			if err != nil {
				continue // TODO
			}
		}
	}

	return nil
}

func (dm *DomainMonitor) refreshPods() error {
	pods, err := dm.podFinder.FindPods()
	if err != nil {
		return err
	}

	podsToRemove := []string{}
	for name, _ := range dm.pods {
		_, ok := pods[name]
		if !ok {
			// pod is gone
			podsToRemove = append(podsToRemove, name)
		}
	}

	for _, name := range podsToRemove {
		delete(dm.pods, name)
	}

	for name, podInfo := range pods {
		// since pod content is immutable, we can just overwrite the old data
		// TODO: log diffs
		dm.pods[name] = podInfo
	}

	return nil
}
