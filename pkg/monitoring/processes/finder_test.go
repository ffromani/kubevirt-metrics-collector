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
)

func TestPodMapBasic(t *testing.T) {
	myPid := int32(os.Getpid())
	procs := make(map[string][]int32)
	procs["self"] = []int32{myPid}
	sc := &SelfScanner{}

	pods := make(PodMap)
	podMap, err := pods.MapProcsToPods(sc, procs)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	info, ok := podMap["selfPod"]
	if !ok {
		t.Errorf("missing selfPod in %#v", podMap)
		return
	}
	if len(info.Procs) != 1 {
		t.Errorf("unexpected procs: %#v", info.Procs)
		return
	}
	if info.Procs[0].Pid != myPid {
		t.Errorf("unexpected pid: found %v expected %v", info.Procs[0].Pid, myPid)
	}
}
