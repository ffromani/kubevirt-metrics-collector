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

import "github.com/shirou/gopsutil/process"

type PodInfo struct {
	Procs []*process.Process
}

type PodFinder interface {
	FindPods() (map[string]*PodInfo, error)
	FindPodByPID(pid int32) (string, error)
}

type PodMap map[string]*PodInfo

func (pods PodMap) MapProcsToPods(pf PodFinder, procs map[string][]int32) (PodMap, error) {
	for _, pids := range procs {
		for _, pid := range pids {
			podName, err := pf.FindPodByPID(pid)
			if err != nil {
				continue // TODO: log
			}

			podInfo, ok := pods[podName]
			if !ok {
				podInfo = &PodInfo{}
				pods[podName] = podInfo
			}

			proc, err := process.NewProcess(pid)
			if err != nil {
				continue // TODO: log
			}

			podInfo.Procs = append(podInfo.Procs, proc)
		}
	}
	return pods, nil
}
