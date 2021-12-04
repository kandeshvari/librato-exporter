package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"net"
	"strings"
	"time"
)

const (
	urlLibratoMeasurements = "https://metrics.librato.com/metrics-api/v1/measurements/"
)

func putInternalMetric(key, metric string, value int64) {
	err := internalMetrics.Put(fmt.Sprintf("%s{metric=%q}", key, metric), value)
	if err != nil {
		log.WithError(err).
			WithField("key", key).
			WithField("metric", metric).
			Error("unable to update internal counter")
	}

}

func goRequestMetricsLoop() {
	metrics := strings.Split(metricsFilter, ",")
	for {
		for _, metricName := range metrics {
			tmNow := time.Now().Unix()
			req := makeRequest(metricName, tmNow-metricsOffset, tmNow, metricsResolution)
			body, err := doRequest(req)
			if err != nil {
				log.WithError(err).WithField("metric", metricName).Error("can't make request")
				continue
			}

			// unmarshal metric
			m := &Metric{}
			err = json.Unmarshal(body, m)
			if err != nil {
				log.WithError(err).WithField("metric", metricName).Error("can't parse response")
				continue
			}

			// update metrics holder
			err = updateMetricsHolder(m)
			if err != nil {
				log.WithError(err).WithField("metric", metricName).Error("can't error update metrics holder")
				continue
			}

			// get some stats
			numMeasurements := 0
			numSeries := len(m.Series)
			if numSeries > 0 {
				numMeasurements = len(m.Series[0].Measurements)
			}

			// update internal metrics
			putInternalMetric("librato_exporter_num_series", metricName, int64(numSeries))
			putInternalMetric("librato_exporter_num_measurements", metricName, int64(numMeasurements))

			log.WithField("metric", metricName).
				WithField("num_series", numSeries).
				WithField("num_measurements", numMeasurements).
				Info("got from librato")
		}

		// sleep for interval after gathering
		time.Sleep(time.Second * time.Duration(libratoInterval))
	}
}

func doRequest(req *fasthttp.Request) ([]byte, error) {
	triesCount := int64(0)
DO_RETRY:
	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
			return fasthttp.DialTimeout(addr, time.Second*5)
		},
	}
	err := client.Do(req, resp)
	if err != nil {
		if triesCount <= maxRequestTries {
			log.WithError(err).WithField("try", triesCount).Warn("can't do request. Retry")
			triesCount++
			goto DO_RETRY
		}
		log.WithError(err).Error("can't do request")
		return nil, err
	}
	return resp.Body(), nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func makeRequest(metric string, startTime, stopTime, resolution int64) *fasthttp.Request {
	url := fmt.Sprintf("%s%s?start_time=%d&end_time=%d&resolution=%d", urlLibratoMeasurements, metric, startTime, stopTime, resolution)
	if metricsSummaryFunc != "" {
		url = fmt.Sprintf("%s&summary_function=%s", url, metricsSummaryFunc)
	}
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.Add("Authorization", "Basic "+basicAuth(libratoEmail, libratoToken))
	req.Header.Add("User-Agent", "curl/7.38.0")
	return req
}
