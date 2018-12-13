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
	"os"
	"sync"
	"time"

	"github.com/shirou/gopsutil/process"
)

type PodInfoMap map[string]*PodInfo

type Monitor interface {
	Update() (PodInfoMap, error)
}

const FreshnessThreshold = 1 * time.Second

type DomainMonitor struct {
	podFinder PodFinder
	lock      sync.RWMutex
	pods      PodInfoMap
	timestamp time.Time
}

type SelfMonitor struct {
	pods PodInfoMap
}

func NewSelfMonitor() (Monitor, error) {
	mon := &SelfMonitor{
		pods: make(PodInfoMap),
	}
	self, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return mon, err
	}
	info := &PodInfo{}
	info.Procs = append(info.Procs, self)
	mon.pods["self"] = info
	return mon, nil
}

func (sm *SelfMonitor) Update() (PodInfoMap, error) {
	return sm.pods, nil
}

func NewDomainMonitor(podFinder PodFinder) (Monitor, error) {
	return &DomainMonitor{
		podFinder: podFinder,
		pods:      make(PodInfoMap),
	}, nil
}

func (dm *DomainMonitor) Update() (PodInfoMap, error) {
	// this is racy, and we don't care
	age := time.Now().Sub(dm.timestamp)
	if age <= FreshnessThreshold {
		return dm.currentPodInfo()
	}
	return dm.updatePodInfo()

}

func (dm *DomainMonitor) currentPodInfo() (PodInfoMap, error) {
	dm.lock.RLock()
	defer dm.lock.RUnlock()
	return dm.pods, nil
}

func (dm *DomainMonitor) updatePodInfo() (PodInfoMap, error) {
	dm.lock.Lock()
	defer dm.lock.Unlock()
	pods, err := dm.podFinder.FindPods()
	if err != nil {
		log.Printf("error finding available pods: %v", err)
		return make(PodInfoMap), err
	}

	podsToRemove := []string{}
	for name := range dm.pods {
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

	log.Printf("refreshed %v pods", len(dm.pods))
	return dm.pods, nil
}
