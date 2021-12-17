package main

import log "github.com/sirupsen/logrus"

func main() {
	parseFlags()

	log.WithField("version", VERSION).
		WithField("go", GO_VERSION).
		WithField("commit", GIT_REVISION).
		WithField("build_date", BUILD_DATE).
		Info("starting librato-exporter")

	go goRequestMetricsLoop()
	go goMetricsGC()

	runServer()
}
