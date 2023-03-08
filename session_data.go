package main

import (
	"encoding/xml"
	"errors"
	log "github.com/sirupsen/logrus"
	"strings"
)

type SessionData struct {
	XMLName     xml.Name         `xml:"perfmonCollectSessionDataResponse"`
	Text        string           `xml:",chardata"`
	Ns1         string           `xml:"ns1,attr"`
	CollectData []OneCollectData `xml:"perfmonCollectSessionDataReturn"`
}

type OneCollectData struct {
	Text    string  `xml:",chardata"`
	Name    string  `xml:"Name"`
	Value   float64 `xml:"Value"`
	CStatus string  `xml:"CStatus"`
}

// processData base on collected data update Prometheus metrics
func (s *SessionData) processData() {
	var server, counter string
	var err error
	for _, data := range s.CollectData {
		server, _, counter, err = data.splitName()
		if err != nil {
			continue
		}
		if config.Metrics.enablePrometheusCounter(counter) {
			if strings.HasSuffix(strings.ToLower(counter), "failed") {
				newVal := data.Value - counterActual[counter]
				if newVal < 0 {
					//log.WithFields(log.Fields{FieldRoutine: "processData", FieldMetricsName: counter, "counter": counterActual[counter], "counter new": data.Value}).
					//	Infof("reset metrics for counter %s, because difference is %f", counter, newVal)
					//var counterDetail *CounterDetails
					//counterDetail, err = monitors.GetCounterDetails(counter)
					//if err != nil {
					//	log.WithFields(log.Fields{FieldRoutine: "processData", FieldMetricsName: counter}).Errorf("not defined description for %s", counter)
					//	continue
					//}
					//prometheus.Unregister(counterMetrics[counter])
					//counterMetrics[counter] = prometheus.NewCounterVec(
					//	prometheus.CounterOpts{
					//		Name: getPrometheusName(counter),
					//		Help: counterDetail.description,
					//	}, []string{"server"})
					//prometheus.MustRegister(counterMetrics[counter])
					//counterActual[counter] = float64(0)
					//newVal = data.Value
					continue
				}
				counterMetrics[counter].WithLabelValues(server).Add(newVal)
				counterActual[counter] = data.Value
			} else {
				callMetrics[counter].WithLabelValues(server).Set(data.Value)
			}
		}
	}
}

// splitName split data path to parts include group
func (o *OneCollectData) splitName() (server string, group string, counter string, err error) {
	v := strings.Trim(o.Name, "\\")
	subst := strings.Split(v, "\\")
	if len(subst) != 3 {
		log.WithFields(log.Fields{FieldRoutine: "splitName", "name": o.Name}).Error("problem split counter name")
		return "", "", "", errors.New("problem split name")
	}
	return subst[0], subst[1], subst[2], nil
}
