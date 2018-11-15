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

type procInfo struct {
	domain string
	exe    string
}

type NullUpdater struct {
	CPU    []procInfo
	Memory []procInfo
}

func (nu *NullUpdater) UpdateCPU(domain string, proc *process.Process) error {
	exe, err := proc.Exe()
	if err != nil {
		return err
	}
	nu.CPU = append(nu.CPU, procInfo{
		domain: domain,
		exe:    exe,
	})
	return nil
}

func (nu *NullUpdater) UpdateMemory(domain string, proc *process.Process) error {
	exe, err := proc.Exe()
	if err != nil {
		return err
	}
	nu.Memory = append(nu.Memory, procInfo{
		domain: domain,
		exe:    exe,
	})
	return nil
}

type SelfScanner struct{}

func (sc SelfScanner) FindPods() (map[string]*PodInfo, error) {
	ret := make(map[string]*PodInfo)
	pi := PodInfo{}
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return ret, err
	}

	pi.Procs = append(pi.Procs, proc)
	ret["self"] = &pi
	return ret, nil
}

func TestUpdateHappyPath(t *testing.T) {
	nu := &NullUpdater{}
	mon, err := NewDomainMonitor(SelfScanner{}, nu)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	err = mon.Update()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if len(nu.CPU) != 1 {
		t.Errorf("unexpected CPU updates: %#v", nu.CPU)
	}
	if len(nu.Memory) != 1 {
		t.Errorf("unexpected Memory updates: %#v", nu.Memory)
	}
}
