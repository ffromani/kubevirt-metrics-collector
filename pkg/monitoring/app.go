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
 * Copyright 2019 Red Hat, Inc.
 */

package monitoring

import (
	"fmt"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"

	"github.com/fromanirh/kubevirt-metrics-collector/internal/pkg/k8sutils"
	"github.com/fromanirh/kubevirt-metrics-collector/internal/pkg/log"
	"github.com/fromanirh/kubevirt-metrics-collector/internal/pkg/service"

	"github.com/fromanirh/kubevirt-metrics-collector/pkg/monitoring/processes"
)

const (
	defaultPort = 8443
	defaultHost = "0.0.0.0"
)

type App struct {
	service.ServiceListen
	TLSInfo   *k8sutils.TLSInfo
	dumpMode  bool
	fakeMode  bool
	debugMode bool
	checkMode bool
}

var _ service.Service = &App{}

func (app *App) AddFlags() {
	app.InitFlags()

	app.BindAddress = defaultHost
	app.Port = defaultPort

	app.AddCommonFlags()

	flag.StringVarP(&app.TLSInfo.CertFilePath, "cert-file", "c", "", "override path to TLS certificate - you need also the key to enable TLS")
	flag.StringVarP(&app.TLSInfo.KeyFilePath, "key-file", "k", "", "override path to TLS key - you need also the cert to enable TLS")
	flag.BoolVarP(&app.fakeMode, "fake", "F", false, "run even connection to CRI runtime fails")
	flag.BoolVarP(&app.debugMode, "debug", "D", false, "enable pod resolution debug mode")
	flag.BoolVarP(&app.dumpMode, "dump-metrics", "M", false, "dump the available metrics and exit")
	flag.BoolVarP(&app.checkMode, "check-config", "C", false, "validate (and dump) configuration and exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s /path/to/kubevirt-metrics-collector.json\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func (app *App) Run() {
	app.TLSInfo.UpdateFromK8S()
	defer app.TLSInfo.Clean()

	args := flag.Args()

	var err error

	if app.dumpMode {
		co, err := processes.NewSelfCollector()
		if err != nil {
			log.Log.Infof("error creating the collector: %v", err)
			os.Exit(1)
		}
		prometheus.MustRegister(co)

		processes.DumpMetrics(os.Stderr)
		os.Exit(1)
	}

	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	conf, err := processes.NewConfigFromFile(args[0])
	if err != nil {
		log.Log.Infof("error reading the configuration file %s: %v", args[0], err)
		os.Exit(1)
	}

	conf.DebugMode = app.debugMode
	conf.Validate()

	if app.debugMode || app.checkMode {
		spew.Fdump(os.Stderr, conf)
	}

	if app.checkMode {
		return
	}

	log.Log.Infof("kubevirt-metrics-collector started")
	defer log.Log.Infof("kubevirt-metrics-collector stopped")

	co, err := processes.NewCollectorFromConf(conf)
	if err == nil {
		prometheus.MustRegister(co)
	} else {
		log.Log.Warningf("error creating the collector: %v", err)
		if !app.fakeMode {
			os.Exit(2)
		}
	}

	http.Handle("/metrics", promhttp.Handler())
	if app.TLSInfo.IsEnabled() {
		log.Log.Infof("TLS configured, serving over HTTPS")
		log.Log.Infof("%s", http.ListenAndServeTLS(conf.ListenAddress, app.TLSInfo.CertFilePath, app.TLSInfo.KeyFilePath, nil))
	} else {
		log.Log.Infof("TLS *NOT* configured, serving over HTTP")
		log.Log.Infof("%s", http.ListenAndServe(conf.ListenAddress, nil))
	}
}
