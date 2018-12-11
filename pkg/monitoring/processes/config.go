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
	"github.com/fromanirh/kubevirt-metrics-collector/pkg/procscanner"

	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// Config encodes the configuration of the monitoring package
type Config struct {
	Targets       []procscanner.ProcTarget `json:"targets"`
	ListenAddress string                   `json:"listenaddress"`
	CRIEndPoint   string                   `json:"criendpoint"`
	Hostname      string                   `json:"hostname"`
	DebugMode     bool                     `json:"debugmode"`
}

// NewConfig creates a new Config object with the current defaults
func NewConfig() *Config {
	return &Config{}
}

// NewConfigFromFile creates a new Config object with the settings taken from the given file.
// Should the file not specify a setting, the value is the default one (see NewCOnfig)
func NewConfigFromFile(confFile string) (*Config, error) {
	conf := NewConfig()

	err := readFile(conf, confFile)
	if err != nil {
		return nil, fmt.Errorf("error reading the configuration on '%s': %s", confFile, err)
	}

	return conf, nil
}

// Validate returns nil if the current configuration is consistent and legal, an error otherwise
func (c *Config) Validate() error {
	// mandatory
	if len(c.Targets) == 0 {
		return errors.New("missing process(es) to track")
	}
	if c.ListenAddress == "" {
		return errors.New("missing listen address")
	}
	if c.CRIEndPoint == "" {
		return errors.New("missing CRI endpoint")
	}
	if c.Hostname == "" {
		var err error
		c.Hostname = os.Getenv("KUBE_NODE_NAME")
		if c.Hostname == "" {
			c.Hostname, err = os.Hostname()
			if err != nil {
				return fmt.Errorf("error getting the host name: %s", err)
			}
		}
	}
	// noone really cares about DebugMode
	return nil
}

func readFile(conf *Config, path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if len(content) > 0 {
		err = json.Unmarshal(content, conf)
		if err != nil {
			return err
		}
	}

	return nil
}
