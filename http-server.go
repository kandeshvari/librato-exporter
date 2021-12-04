package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func handlerMetrics(w http.ResponseWriter, r *http.Request) {
	//r.Response.Header.Add("Content-Type", "text/plain")
	numMetrics, err := writeMetrics(w)
	if err != nil {
		log.WithError(err).Error("can't handle /metrics request")
	}

	internalMetrics.Range(func(key, value interface{}) bool {
		_, err := fmt.Fprintf(w, "%s %d\n", key, *value.(*int64))
		if err != nil {
			log.WithError(err).Error("can't write internal metric")
		}
		return true
	})

	log.WithField("num_metrics", numMetrics).Info("metrics requested")
}

func runServer() {
	http.HandleFunc("/metrics", handlerMetrics)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`librato-exporter. Use /metrics for scraping`))
	})

	log.WithField("address", listenAddress).Info("listen on")
	log.Fatal(http.ListenAndServe(listenAddress, nil))

}
