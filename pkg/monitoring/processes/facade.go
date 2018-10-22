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
	"github.com/fromanirh/kube-metrics-collector/pkg/procscanner"

	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	DefaultInterval = "5s"
)

type Config struct {
	Targets       []procscanner.ProcTarget `json:"targets"`
	Interval      string                   `json:"interval"`
	Hostname      string                   `json:"hostname"`
	ListenAddress string                   `json:"listenaddress"`
	CRIEndPoint   string                   `json:"criendpoint"`
	DebugMode     bool                     `json:"debugmode"`
}

func NewConfig() *Config {
	return &Config{
		Interval: DefaultInterval,
	}
}

func (c *Config) Validate() error {
	if len(c.Targets) == 0 {
		errors.New("missing process(es) to track")
	}
	if c.CRIEndPoint == "" {
		errors.New("missing CRI endpoint")
	}
	if c.ListenAddress == "" {
		errors.New("missing listen address")
	}
	return nil
}

func (c *Config) Setup(confFile, intervalString string, debugMode bool) error {
	err := readFile(c, confFile)
	if err != nil {
		errors.New(fmt.Sprintf("error reading the configuration on '%s': %s", confFile, err))
	}

	if c.Hostname == "" {
		c.Hostname, err = os.Hostname()
		if err != nil {
			errors.New(fmt.Sprintf("error getting the host name: %s", err))
		}
	}
	if c.Interval == "" {
		c.Interval = intervalString
	}
	if debugMode {
		c.DebugMode = true
	}
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

func Collect(conf *Config, scanner procscanner.ProcScanner, interval time.Duration) error {
	finder, err := NewCRIPodFinder(conf.CRIEndPoint, DefaultTimeout, scanner)
	if err != nil {
		return err
	}

	mon, err := NewDomainMonitor(conf.Hostname, finder)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	for {
		t := <-ticker.C
		if conf.DebugMode {
			log.Printf("updating at %v", t)
		}

		start := time.Now()
		err := mon.Update()
		stop := time.Now()

		if conf.DebugMode {
			log.Printf("update took %v", stop.Sub(start))
		}
		if err != nil {
			log.Printf("error while updating: %v", err)
		}
	}
}
