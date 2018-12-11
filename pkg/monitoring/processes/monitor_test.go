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
	"os"
	"testing"

	"github.com/shirou/gopsutil/process"
)

type SelfScanner struct {
	Skip bool
}

func (sc *SelfScanner) FindPods() (map[string]*PodInfo, error) {
	ret := make(map[string]*PodInfo)
	if !sc.Skip {
		pi := PodInfo{}
		proc, err := process.NewProcess(int32(os.Getpid()))
		if err != nil {
			return ret, err
		}

		pi.Procs = append(pi.Procs, proc)
		ret["self"] = &pi
	}
	return ret, nil
}

func (sc *SelfScanner) FindPodByPID(pid int32) (string, error) {
	return "selfPod", nil
}

func TestUpdateHappyPath(t *testing.T) {
	mon, err := NewDomainMonitor(&SelfScanner{})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	pods, err := mon.Update()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if len(pods) != 1 {
		t.Errorf("unexpected pods: %#v", pods)
	}
}

func TestUpdatePodDisappears(t *testing.T) {
	sc := &SelfScanner{}
	mon, err := NewDomainMonitor(sc)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	pods, err := mon.Update()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if len(pods) != 1 {
		t.Errorf("unexpected pods: %#v", pods)
	}

	sc.Skip = true
	pods, err = mon.Update()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if len(pods) != 0 {
		t.Errorf("unexpected pods: %#v", pods)
	}
}
