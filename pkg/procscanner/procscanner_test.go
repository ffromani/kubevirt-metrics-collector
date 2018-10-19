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

package procscanner

import (
	"testing"
)

func TestReadProcCmdline(t *testing.T) {
	argv := readProcCmdline("/proc/1/cmdline")
	if len(argv) == 0 {
		t.Errorf("failed to read cmdline of pid 1")
	}
}

func TestReadProcCmdlineInexistent(t *testing.T) {
	argv := readProcCmdline("/proc/0/cmdline")
	if len(argv) > 0 {
		t.Errorf("Unexpected data for pid 0: %#v", argv)
	}
}

func TestScanWithoutGlob(t *testing.T) {
	target := ProcTarget{
		Name: "journald",
		Argv: []string{"/usr/lib/systemd/systemd-journald"},
	}
	ps := ProcScanner{
		Targets: []ProcTarget{target},
	}
	ret, err := ps.Scan("testdata/proc")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if len(ret) != 1 {
		t.Errorf("Unexpected results: %#v", ret)
	}
	pids, ok := ret["journald"]
	if !ok {
		t.Errorf("journald not found: %#v", ret)
	}
	if len(pids) != 1 || pids[0] != 2159 {
		t.Errorf("Unexpected pids: %#v", pids)
	}
}
