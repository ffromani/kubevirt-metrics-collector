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

	"log"
	"time"
)

// TODO: this should be const
var NOTIME = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)

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
	return CollectLoop(mon, ticker.C, conf.DebugMode)
}

func CollectLoop(mon Monitor, clock <-chan time.Time, debugMode bool) error {
	for {
		t := <-clock
		if t == NOTIME {
			break
		}
		CollectStats(mon, t, debugMode)
	}
	return nil
}

func CollectStats(mon Monitor, now time.Time, debugMode bool) error {
	if debugMode {
		log.Printf("updating at %v", now)
	}

	start := time.Now()
	err := mon.Update()
	stop := time.Now()

	if debugMode {
		log.Printf("update took %v", stop.Sub(start))
	}
	if err != nil {
		log.Printf("error while updating: %v", err)
	}
	return err
}
