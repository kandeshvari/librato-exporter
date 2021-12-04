package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"sort"
	"strings"
	"sync"
)

// Metric represents a Librato Metric.
type Metric struct {
	Name       string                 `json:"name"`
	Series     []Serie                `json:"series"`
	Attributes map[string]interface{} `json:"attributes"`
	Resolution uint                   `json:"resolution"`
	Period     uint                   `json:"period"`
}

type Serie struct {
	Tags         map[string]string `json:"tags"`
	Measurements []TimeVal         `json:"measurements"`
}

type TimeVal struct {
	Time  uint    `json:"time"`
	Value float64 `json:"value"`
}

type StringSlice []string

var metricsHolder = map[string]TimeVal{}
var lock = sync.RWMutex{}

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }

func collectLabels(m map[string]string) string {
	tags := make([]string, 0, len(m))
	for k, v := range m {
		tags = append(tags, fmt.Sprintf("%s=%q", k, v))
	}
	sort.Sort(StringSlice(tags))
	return strings.Join(tags, ",")
}

func collectLabelsInterface(m map[string]interface{}, filterAttrs []string) string {
	tags := make([]string, 0, len(m))

OUTER:
	for k, v := range m {
		for _, f := range filterAttrs {
			if strings.Compare(k, f) == 0 {
				continue OUTER
			}
		}
		tags = append(tags, fmt.Sprintf("%s=%q", k, fmt.Sprintf("%v", v)))
	}
	sort.Sort(StringSlice(tags))
	return strings.Join(tags, ",")
}

func updateMetricsHolder(ms *Metric) error {
	var tagsLabels string
	attrLabels := collectLabelsInterface(ms.Attributes, []string{"created_by_ua"})
	for _, s := range ms.Series {
		tagsLabels = collectLabels(s.Tags)

		mReplaced := strings.ReplaceAll(ms.Name, ".", ":")
		mReplaced = strings.ReplaceAll(mReplaced, "-", "_")
		metricName := fmt.Sprintf("%s{%s,%s}", mReplaced, attrLabels, tagsLabels)
		if len(s.Measurements) < 1 {
			return fmt.Errorf("empty measurements on %s", metricName)
		}
		lock.Lock()
		metricsHolder[metricName] = s.Measurements[len(s.Measurements)-1]
		log.Debugf("added metric: %s %f", metricName, metricsHolder[metricName].Value)
		lock.Unlock()
	}
	return nil
}

func writeMetrics(w io.Writer) (int, error) {
	lock.RLock()
	numMetrics := 0
	for k, v := range metricsHolder {
		_, err := fmt.Fprintf(w, "%s %f\n", k, v.Value)
		if err != nil {
			return numMetrics, fmt.Errorf("can't write metric %s : %w", k, err)
		}
		numMetrics += 1
	}
	lock.RUnlock()
	return numMetrics, nil
}
