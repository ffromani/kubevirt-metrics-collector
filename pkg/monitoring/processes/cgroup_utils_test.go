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

import "testing"

func TestNoCgroupFile(t *testing.T) {
	containerID, cgroupStyle := parseProcCGroupEntry("/this/path/does/not/exist")
	if containerID != "" {
		t.Errorf("unexpected containerID: %v", containerID)
	}
	if cgroupStyle != MissingCGroup {
		t.Errorf("unexpected cgroupStyle: %v", cgroupStyle)
	}
}

func TestNoCgroupData(t *testing.T) {
	containerID, cgroupStyle := parseProcCGroupEntry("/dev/null")
	if containerID != "" {
		t.Errorf("unexpected containerID: %v", containerID)
	}
	if cgroupStyle != MissingCGroup {
		t.Errorf("unexpected cgroupStyle: %v", cgroupStyle)
	}
}

func TestEmptyCgroupData(t *testing.T) {
	containerID, cgroupStyle := parseProcCGroupEntry("testdata/cgroup-empty")
	if containerID != "/" {
		t.Errorf("unexpected containerID: %v", containerID)
	}
	if cgroupStyle != UnknownCGroup {
		t.Errorf("unexpected cgroupStyle: %v", cgroupStyle)
	}
}

func TestEmptyDockerData(t *testing.T) {
	containerID, cgroupStyle := parseProcCGroupEntry("testdata/cgroup-docker")
	if containerID != "0fca315d4f86002639693ca3dbf16e4376a1951ea18ab537a38bb12478de161c" {
		t.Errorf("unexpected containerID: %v", containerID)
	}
	if cgroupStyle != DockerCGroup {
		t.Errorf("unexpected cgroupStyle: %v", cgroupStyle)
	}
}

func TestMalformedContent(t *testing.T) {
	containerID, cgroupStyle := parseProcCGroupEntry("testdata/malformed")
	if containerID != "" {
		t.Errorf("unexpected containerID: %v", containerID)
	}
	if cgroupStyle != MalformedCGroup {
		t.Errorf("unexpected cgroupStyle: %v", cgroupStyle)
	}
}
