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
	"testing"

	"github.com/fromanirh/kubevirt-metrics-collector/pkg/procscanner"
)

func TestConfigInvalidByDefault(t *testing.T) {
	checkInvalid(t, NewConfig())
	checkInvalid(t, &Config{})
}

func TestConfigFillValidate(t *testing.T) {
	conf := NewConfig()
	conf.Targets = []procscanner.ProcTarget{
		{
			Name: "init",
			Argv: []string{"/sbin/init"},
		},
	}
	conf.ListenAddress = ":9999"
	conf.CRIEndPoint = "/var/run/cri.sock"
	conf.Hostname = "foobar.test.lan"
	conf.DebugMode = false

	checkValid(t, conf)
}

func TestConfigMinimalFillValidate(t *testing.T) {
	conf := NewConfig()
	conf.Targets = []procscanner.ProcTarget{
		{
			Name: "init",
			Argv: []string{"/sbin/init"},
		},
	}
	conf.ListenAddress = ":9999"
	conf.CRIEndPoint = "/var/run/cri.sock"

	checkValid(t, conf)
}

func TestConfigInvalidWithoutTargets(t *testing.T) {
	conf := NewConfig()
	conf.ListenAddress = ":9999"
	conf.CRIEndPoint = "/var/run/cri.sock"

	checkInvalid(t, conf)
}

func TestConfigInvalidWithoutListenAddress(t *testing.T) {
	conf := NewConfig()
	conf.Targets = []procscanner.ProcTarget{
		{
			Name: "init",
			Argv: []string{"/sbin/init"},
		},
	}
	conf.CRIEndPoint = "/var/run/cri.sock"

	checkInvalid(t, conf)
}

func TestConfigInvalidWithoutCRIEndPoint(t *testing.T) {
	conf := NewConfig()
	conf.Targets = []procscanner.ProcTarget{
		{
			Name: "init",
			Argv: []string{"/sbin/init"},
		},
	}
	conf.ListenAddress = ":9999"

	checkInvalid(t, conf)
}

func TestConfigReadFromFile(t *testing.T) {
	conf, err := NewConfigFromFile("testconf.json")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	checkValid(t, conf)
}

func TestConfigReadFromFileInexistent(t *testing.T) {
	conf, err := NewConfigFromFile("missing.json")
	if err == nil {
		t.Errorf("unexpected success")
		return
	}
	if conf != nil {
		t.Errorf("unexpected conf: %v", conf)
		return
	}
}

func checkValid(t *testing.T, conf *Config) {
	err := conf.Validate()
	if err != nil {
		t.Errorf("conf unexpectedly invalid: %#v: %v", conf, err)
	}
}

func checkInvalid(t *testing.T, conf *Config) {
	err := conf.Validate()
	if err == nil {
		t.Errorf("conf unexpectedly valid: %#v", conf)
	}
}
