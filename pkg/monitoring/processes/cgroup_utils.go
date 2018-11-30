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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Classifiers for CGroup found
const (
	MissingCGroup = iota
	MalformedCGroup
	UnknownCGroup
	DockerCGroup
)

// FindContainerIDByCGroup fetches cgroup information for the given PID
// Returns the cgroup name and the CGroup classifier
// Should the pid belong to more than one cgroup: returns the first one, in kernel order.
func FindContainerIDByCGroup(pid int32) (string, int) {
	return parseProcCGroupEntry(fmt.Sprintf("/proc/%d/cgroup", pid))
}

func parseProcCGroupEntry(entry string) (string, int) {
	file, err := os.Open(entry)
	defer file.Close()

	if err != nil {
		return "", MissingCGroup
	}

	reader := bufio.NewReader(file)
	line, isPrefix, err := reader.ReadLine()
	if isPrefix || err != nil {
		return "", MissingCGroup
	}

	return parseProcCGroupLine(string(line))
}

func parseProcCGroupLine(line string) (string, int) {
	fields := strings.Split(line, ":")
	// per cgroups(7):
	// hierarchy-ID:controller-list:cgroup-path
	if len(fields) != 3 {
		return "", MalformedCGroup
	}
	name := filepath.Base(fields[2])
	if strings.HasPrefix(name, "docker-") {
		name = name[len("docker-"):]
		if strings.HasSuffix(name, ".scope") {
			name = name[:len(name)-len(".scope")]
		}
		return name, DockerCGroup
	}
	return name, UnknownCGroup
}
