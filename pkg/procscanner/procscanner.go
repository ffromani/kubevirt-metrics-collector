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
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// ProcTarget defines a target for a ProcScanner
type ProcTarget struct {
	Name string   // user-visible name of the target. If not specified, Argv[0] is used
	Argv []string // command line arguments to match to identify the target
}

// ProcScanner scans the Linux procfs to find the PID of the given targets
type ProcScanner struct {
	Targets []ProcTarget
}

// Scan performs the scan of the procfs and returns a mapping between ProcTargets - by name, and the PIDs
// basePath is the path where procfs is mounted (default: /proc)
// returns a map whose keys are ProcTargets Names, and whose valuesa is an unordered list of all the PIDs
// of the live instances of the ProcTargets.
// For example, you should expect only one entry for '/sbin/init' (or equivalent), but you should expect
// more than one entry for '/bin/sh' or '/bin/getty' (or equivalent)
func (p *ProcScanner) Scan(basePath string) (map[string][]int32, error) {
	res := make(map[string][]int32)

	procEntries, err := filepath.Glob(path.Join(basePath, "*", "cmdline"))
	if err != nil {
		return res, err
	}

	for _, procEntry := range procEntries {
		argv := readProcCmdline(procEntry)
		if argv == nil || len(argv) == 0 {
			continue
		}

		name, ok := p.findTarget(argv)
		if !ok {
			continue
		}

		items := strings.Split(procEntry, string(os.PathSeparator))
		// "", "proc", "$PID", "cmdline"
		pid, err := strconv.Atoi(items[2])
		if err != nil {
			return res, err
		}

		_, ok = res[name]
		if !ok {
			res[name] = []int32{}
		}
		res[name] = append(res[name], int32(pid))
	}

	return res, nil
}

func (p *ProcScanner) findTarget(argv []string) (string, bool) {
	for _, target := range p.Targets {
		match, err := MatchArgv(argv, target.Argv)
		if err != nil {
			break
		} else if match {
			return target.Name, true
		}

	}
	return "", false
}

func readProcCmdline(pathname string) []string {
	argv := make([]string, 0)
	content, err := ioutil.ReadFile(pathname)
	if err == nil && len(content) > 0 {
		for _, chunk := range bytes.Split(content, []byte{0}) {
			arg := string(chunk)
			if len(arg) > 0 {
				argv = append(argv, arg)
			}
		}
	}
	return argv
}

// MatchArgv matches two commanddlines.
// Return error iff the normalization of the commandline fails.
func MatchArgv(argv, model []string) (bool, error) {
	ref := model
	oth := argv
	if len(argv) < len(model) {
		ref = argv
		oth = model
	}
	for idx, elem := range ref {
		matched, err := filepath.Match(elem, oth[idx])
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}
	}
	return true, nil
}
