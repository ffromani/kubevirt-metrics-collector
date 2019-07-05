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
 */

package main

import (
	"os"

	"github.com/fromanirh/kubevirt-metrics-collector/internal/pkg/k8sutils"
	"github.com/fromanirh/kubevirt-metrics-collector/internal/pkg/log"
	"github.com/fromanirh/kubevirt-metrics-collector/internal/pkg/service"

	"github.com/fromanirh/kubevirt-metrics-collector/pkg/monitoring"
)

func Main() int {
	tlsInfo := &k8sutils.TLSInfo{}
	app := &monitoring.App{
		TLSInfo: tlsInfo,
	}
	service.Setup(app)
	log.InitializeLogging("kubevirt-metrics-collector")
	app.Run()
	return 0
}

func main() {
	os.Exit(Main())
}
