package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	VERSION      string
	BUILD_DATE   string
	GIT_REVISION string
	GO_VERSION   string
)

var (
	listenAddress      = "0.0.0.0:9800"
	libratoEmail       = ""
	libratoToken       = ""
	libratoInterval    = int64(60)
	metricsFilter      = ""
	metricsResolution  = int64(1)
	metricsOffset      = int64(120)
	metricsGCPeriod    = int64(86400)
	metricsSummaryFunc = ""
	maxRequestTries    = int64(3)
)

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, "usage: %s [options]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func parseFlags() {
	flag.StringVar(&listenAddress, "address", "0.0.0.0:9800", "librato-exporter listen address")
	flag.Int64Var(&maxRequestTries, "max_tries", 3, "number of tries to do a request")

	flag.StringVar(&libratoEmail, "librato.email", "", "Librato Email account owner of token.")
	flag.StringVar(&libratoToken, "librato.token", "", "Librato API token created by email.")
	flag.Int64Var(&libratoInterval, "librato.interval", 60, "Interval in seconds to retrieve metrics from API")

	flag.StringVar(&metricsFilter, "metrics.filter", "", "List of metrics sepparated by comma.")
	flag.Int64Var(&metricsResolution, "metrics.resolution", 1, "Metrics resolution in seconds.")
	flag.Int64Var(&metricsOffset, "metrics.offset", 120, "How long wait to remove `nodata` metrics from output")
	flag.Int64Var(&metricsGCPeriod, "metrics.gc_period", 86400, "Time offset in seconds to define the start timestamp.")
	flag.StringVar(&metricsSummaryFunc, "metrics.summary_func", "", "Summary function.")

	flag.Usage = usage
	flag.Parse()
}
