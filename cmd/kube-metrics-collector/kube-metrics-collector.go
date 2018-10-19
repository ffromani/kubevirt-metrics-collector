package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"

	"github.com/fromanirh/kube-metrics-collector/pkg/monitoring/processes"
	_ "github.com/fromanirh/kube-metrics-collector/pkg/monitoring/processes/prometheus"
	"github.com/fromanirh/kube-metrics-collector/pkg/procscanner"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func (c *Config) Validate() *Config {
	if len(c.Targets) == 0 {
		log.Fatalf("missing process(es) to track")
	}
	if c.CRIEndPoint == "" {
		log.Fatalf("missing CRI endpoint")
	}
	if c.ListenAddress == "" {
		log.Fatalf("missing listen address")
	}
	return c
}

func (c *Config) Setup(confFile, intervalString string, debugMode bool) *Config {
	err := readFile(c, confFile)
	if err != nil {
		log.Fatalf("error reading the configuration on '%s': %s", confFile, err)
	}

	if c.Hostname == "" {
		c.Hostname, err = os.Hostname()
		if err != nil {
			log.Fatalf("error getting the host name: %s", err)
		}
	}
	if c.Interval == "" {
		c.Interval = intervalString
	}
	if debugMode {
		c.DebugMode = true
	}
	return c
}

func readFile(conf *Config, path string) error {
	log.Printf("trying configuration: %s", path)

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

	log.Printf("read from file: %s", path)
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s /path/to/kube-metrics-collector.json\n", os.Args[0])
		flag.PrintDefaults()
	}
	intervalString := flag.StringP("interval", "I", DefaultInterval, "metrics collection interval")
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

	conf := &Config{Interval: DefaultInterval}
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

	finder, err := processes.NewCRIPodFinder(conf.CRIEndPoint, processes.DefaultTimeout, scanner)
	if err != nil {
		log.Fatalf("%v", err)
	}

	mon, err := processes.NewDomainMonitor(conf.Hostname, finder)
	if err != nil {
		log.Fatalf("%v", err)
	}

	go collect(mon, interval, conf.DebugMode)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(conf.ListenAddress, nil))
}

func collect(mon *processes.DomainMonitor, interval time.Duration, debugMode bool) {
	ticker := time.NewTicker(interval)
	for {
		t := <-ticker.C
		if debugMode {
			log.Printf("updating at %v", t)
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
	}

}
