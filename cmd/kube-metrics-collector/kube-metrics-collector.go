package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"

	"github.com/fromanirh/kube-metrics-collector/pkg/monitoring/processes"
	_ "github.com/fromanirh/kube-metrics-collector/pkg/monitoring/processes/prometheus"
	"github.com/fromanirh/kube-metrics-collector/pkg/procscanner"

	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s /path/to/kube-metrics-collector.json\n", os.Args[0])
		flag.PrintDefaults()
	}
	intervalString := flag.StringP("interval", "I", processes.DefaultInterval, "metrics collection interval")
	debugMode := flag.BoolP("debug", "D", false, "enable pod resolution debug mode")
	checkMode := flag.BoolP("check-config", "C", false, "validate (and dump) configuration and exit")
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
		return
	}

	log.Printf("kube-metrics-collectorer started")
	defer log.Printf("kube-metrics-collectorer stopped")

	var err error

	conf := processes.NewConfig()
	conf.Setup(args[0], *intervalString, *debugMode)
	conf.Validate()

	interval, err := time.ParseDuration(conf.Interval)
	if err != nil {
		log.Fatalf("error getting the polling interval: %s", err)
	}

	scanner := procscanner.ProcScanner{}
	scanner.Targets = conf.Targets
	if conf.DebugMode {
		spew.Fdump(os.Stderr, scanner)
	}

	// here because this way the debug mode can emit both conf and scanner content
	if *checkMode {
		spew.Fdump(os.Stderr, conf)
		return
	}

	go processes.Collect(conf, scanner, interval)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(conf.ListenAddress, nil))
}
