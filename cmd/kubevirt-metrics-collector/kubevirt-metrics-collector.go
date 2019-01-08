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
	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"

	"github.com/fromanirh/kubevirt-metrics-collector/pkg/monitoring/processes"

	"fmt"
	"log"
	"net/http"
	"os"
)

func Main() int {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s /path/to/kubevirt-metrics-collector.json\n", os.Args[0])
		flag.PrintDefaults()
	}
	certFilePath := flag.StringP("cert-file", "c", "", "path to TLS certificate - you need also the key to enable TLS")
	keyFilePath := flag.StringP("key-file", "k", "", "path to TLS key - you need also the cert to enable TLS")
	fakeMode := flag.BoolP("fake", "F", false, "run even connection to CRI runtime fails")
	debugMode := flag.BoolP("debug", "D", false, "enable pod resolution debug mode")
	dumpMode := flag.BoolP("dump-metrics", "M", false, "dump the available metrics and exit")
	checkMode := flag.BoolP("check-config", "C", false, "validate (and dump) configuration and exit")
	flag.Parse()

	args := flag.Args()

	var err error

	if *dumpMode {
		co, err := processes.NewSelfCollector()
		if err != nil {
			log.Printf("error creating the collector: %v", err)
			return 1
		}
		prometheus.MustRegister(co)

		processes.DumpMetrics(os.Stderr)
		return 1
	}

	if len(args) < 1 {
		flag.Usage()
		return 1
	}

	conf, err := processes.NewConfigFromFile(args[0])
	if err != nil {
		log.Printf("error reading the configuration file %s: %v", args[0], err)
		return 1
	}

	conf.DebugMode = *debugMode
	conf.Validate()

	if *debugMode || *checkMode {
		spew.Fdump(os.Stderr, conf)
	}

	if *checkMode {
		return 0
	}

	log.Printf("kubevirt-metrics-collector started")
	defer log.Printf("kubevirt-metrics-collector stopped")

	co, err := processes.NewCollectorFromConf(conf)
	if err == nil {
		prometheus.MustRegister(co)
	} else {
		log.Printf("error creating the collector: %v", err)
		if !*fakeMode {
			return 2
		}
	}

	http.Handle("/metrics", promhttp.Handler())
	if *certFilePath == "" || *keyFilePath == "" {
		log.Printf("TLS NOT fully configured (cert AND key), serving over HTTP")
		log.Printf("%s", http.ListenAndServe(conf.ListenAddress, nil))
	} else {
		log.Printf("TLS configured: cert='%s' key='%s', serving over HTTPS", *certFilePath, *keyFilePath)
		log.Printf("%s", http.ListenAndServeTLS(conf.ListenAddress, *certFilePath, *keyFilePath, nil))
	}
	return 0
}

func main() {
	os.Exit(Main())
}
