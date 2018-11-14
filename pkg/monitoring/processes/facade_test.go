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
	"errors"
	"testing"
	"time"
)

type FailingMonitor struct{}

func (fm FailingMonitor) Update() error {
	return errors.New("fake error")
}

type FakeMonitor struct{}

func (m FakeMonitor) Update() error {
	return nil
}

func TestFakeMonitor(t *testing.T) {
	for _, flag := range []bool{false, true} {
		err := CollectStats(FakeMonitor{}, time.Now(), flag)
		if err != nil {
			t.Errorf("flag=%v unexpected error %v", flag, err)
		}
	}
}

func TestFailingMonitor(t *testing.T) {
	for _, flag := range []bool{false, true} {
		err := CollectStats(FailingMonitor{}, time.Now(), flag)
		if err == nil {
			t.Errorf("flag=%v unexpected success", flag)
		}
	}
}
