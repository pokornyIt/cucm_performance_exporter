package main

import (
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"time"
)

const (
	Envelope = "<soapenv:Envelope xmlns:soapenv=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:soap=\"http://schemas.cisco.com/ast/soap\">\r\n" +
		"<soapenv:Header/>\r\n" +
		"<soapenv:Body>\r\n" +
		"%s\r\n" +
		"</soapenv:Body>\r\n" +
		"</soapenv:Envelope>"
	EnvelopeList            = "<soap:perfmonListCounter>\r\n<soap:Host>%s</soap:Host>\r\n</soap:perfmonListCounter>"
	QueryCounterDescription = "<soap:perfmonQueryCounterDescription>\r\n<soap:Counter>%s</soap:Counter>\r\n</soap:perfmonQueryCounterDescription>"
)

type ClusterHostMonitorData struct {
	server      string           // server name
	counterList counterGroupList // counterList list of available counters
}

type counterGroupList struct {
	group []counterGroup // counter group array
}

type counterGroup struct {
	groupName     string           // name of group same as used in counter group list
	multiInstance bool             // is multi instance of counter
	counterName   []CounterDetails // list of counter name
}
type CounterDetails struct {
	name        string
	description string
}

type XmlListCounterResponse struct {
	XMLName           xml.Name `xml:"perfmonListCounterResponse"`
	Text              string   `xml:",chardata"`
	Ns1               string   `xml:"ns1,attr"`
	ListCounterReturn []struct {
		Text           string `xml:",chardata"`
		Name           string `xml:"Name"`
		MultiInstance  bool   `xml:"MultiInstance"`
		ArrayOfCounter struct {
			Text string `xml:",chardata"`
			Item []struct {
				Text string `xml:",chardata"`
				Name string `xml:"Name"`
			} `xml:"item"`
		} `xml:"ArrayOfCounter"`
	} `xml:"perfmonListCounterReturn"`
}

type XmlDescriptionCounterResponse struct {
	XMLName                       xml.Name `xml:"perfmonQueryCounterDescriptionResponse"`
	Text                          string   `xml:",chardata"`
	Ns1                           string   `xml:"ns1,attr"`
	QueryCounterDescriptionReturn string   `xml:"perfmonQueryCounterDescriptionReturn"`
}

func (c *counterGroup) counterPathBase(server string, counter string) string {
	return fmt.Sprintf("\\\\%s\\%s\\%s", server, c.groupName, counter)
}

//func (c *counterGroup) counterPathWithInstanceBase(server string, instance string, counter string) string {
//	return fmt.Sprintf("\\\\%s\\%s(%s)\\%s", server, c.groupName, instance, counter)
//}
//
//func (g *counterGroupList) counterPathBase(server string, counter string) string {
//	for _, gr := range g.group {
//		for _, name := range gr.counterName {
//			if name.name == counter {
//				return gr.counterPathBase(server, counter)
//			}
//		}
//	}
//	return ""
//}
//
//func (g *counterGroupList) counterPathWithInstanceBase(server string, instance string, counter string) string {
//	for _, gr := range g.group {
//		for _, name := range gr.counterName {
//			if name.name == counter {
//				return gr.counterPathWithInstanceBase(server, instance, counter)
//			}
//		}
//	}
//	return ""
//}

// NewClusterHostMonitorData create new monitored hosts (CUCM servers) with empty counter group
func NewClusterHostMonitorData(srv string) *ClusterHostMonitorData {
	grp := make([]counterGroup, 0)
	h := ClusterHostMonitorData{
		server:      srv,
		counterList: counterGroupList{group: grp},
	}
	log.WithFields(h.logFields("NewClusterHostMonitorData")).Trace("create session id holder for host")
	return &h

}

func (h *ClusterHostMonitorData) createCounterList(data XmlListCounterResponse) {
	log.WithFields(h.logFields("createCounterList")).Tracef("create counter list from response")
	defer duration(track(log.Fields{FieldRoutine: "createCounterList"}, "procedure ends"))
	for _, listReturn := range data.ListCounterReturn {
		if !inSlice(listReturn.Name, AllowedGroupNames) {
			continue
		}
		m := make([]CounterDetails, 0)
		for _, cnt := range listReturn.ArrayOfCounter.Item {
			if isNameInAllowedCounter(cnt.Name) && config.Metrics.enablePrometheusCounter(cnt.Name) {
				m = append(m, CounterDetails{
					name:        cnt.Name,
					description: "",
				})
			}
		}
		if len(m) > 0 {
			h.counterList.group = append(h.counterList.group, counterGroup{
				groupName:     listReturn.Name,
				multiInstance: listReturn.MultiInstance,
				counterName:   m,
			})
		}
	}
}

func inSlice(name string, list []string) bool {
	for _, v := range list {
		if v == name {
			return true
		}
	}
	return false
}

// AddCounters request PerfMon API for add new counters into session
func (h *ClusterHostMonitorData) AddCounters(client *ApiMonitorClient) (err error) {
	log.WithFields(h.logFields("AddCounter")).Trace("add counters to session")
	defer duration(track(h.logFields("AddCounter"), "procedure ends"))
	cnt := ""
	for _, group := range h.counterList.group {
		for _, counter := range group.counterName {
			cnt = fmt.Sprintf("%s<soap:Counter><soap:Name>%s</soap:Name></soap:Counter>", cnt, group.counterPathBase(h.server, counter.name))
		}
	}

	req := fmt.Sprintf("<soap:perfmonAddCounter><soap:SessionHandle>%s</soap:SessionHandle><soap:ArrayOfCounter>%s</soap:ArrayOfCounter></soap:perfmonAddCounter>", client.session, cnt)
	body, err := client.processRequest("AddCounters", req)

	if body == "401" {
		log.WithFields(h.logFields("AddCounter")).Fatal("user not authorize for use performance API")
	}

	var fault FaultResponse
	e := xml.Unmarshal([]byte(body), &fault)
	if e == nil {
		log.WithFields(log.Fields{"error": fault.FaultCode, "message": fault.FaultString}).WithFields(h.logFields("AddCounters")).Error("problem add counters")
		return err
	}
	log.WithFields(h.logFields("AddCounter")).Trace("success add counters to server")
	return nil
}

// ListCounters collect all counters from API server for specific CUCM host
func (h *ClusterHostMonitorData) ListCounters(client *ApiMonitorClient) (err error) {
	log.WithFields(h.logFields("ListCounters")).Trace("collect counters from server")
	defer duration(track(h.logFields("ListCounters"), "procedure ends"))
	if len(h.counterList.group) > 0 {
		log.WithFields(h.logFields("ListCounters")).Trace("collect counters are read from list")
		return nil
	}
	s := fmt.Sprintf(EnvelopeList, h.server)
	body, err := client.processRequest("ListCounters", s)
	if err != nil && body == "401" {
		log.WithFields(h.logFields("ListCounters")).Fatal("user not authorize for use performance API")
	}

	if err != nil {
		return err
	}

	var list XmlListCounterResponse
	err = xml.Unmarshal([]byte(body), &list)
	if err != nil {
		log.WithFields(h.logFields("ListCounters")).Errorf("problem convert XML body to struct. Error: %s", err)
		return err
	}
	h.createCounterList(list)
	return nil
}

func (h *ClusterHostMonitorData) ReadCounterDescription(client *ApiMonitorClient) (err error) {
	log.WithFields(h.logFields("ReadCounterDescription")).Trace("collect counters descriptions from server")
	defer duration(track(h.logFields("ReadCounterDescription"), "procedure ends"))
	var s string
	var base string
	errCounter := 0

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	for g, group := range h.counterList.group {
		if group.multiInstance {
			log.WithFields(h.logFields("ReadCounterDescription")).Info("multi-instance not implement yet")
			continue
		}

		for c, counter := range group.counterName {
			base = group.counterPathBase(h.server, counter.name)
			select {
			case <-time.After(time.Millisecond * 2):
				break
			case a := <-quit:
				log.WithFields(log.Fields{FieldRoutine: "ReadCounterDescription"}).
					Fatalf("function canceled by user request with signal %s", a)
			}
			log.WithFields(h.logFields("ReadCounterDescription")).WithField(FieldMetricsName, base).
				Tracef("collect counters descriptions for %s", counter.name)
			s = fmt.Sprintf(QueryCounterDescription, base)
			body, errRequest := client.processRequest("ReadCounterDescription", s)
			if body == "401" {
				log.WithFields(h.logFields("ReadCounterDescription")).WithField(FieldMetricsName, base).
					Fatal("user not authorize for use performance API")
			}
			if strings.Contains(body, "RateControl") {
				log.WithFields(h.logFields("ReadCounterDescription")).WithField(FieldMetricsName, base).
					Fatal("exceeded allowed rate for Perfmon information")
			}
			if errRequest != nil {
				errCounter++
				continue
			}

			var description XmlDescriptionCounterResponse
			err = xml.Unmarshal([]byte(body), &description)
			if err != nil {
				log.WithFields(h.logFields("ReadCounterDescription")).WithField(FieldMetricsName, base).
					Errorf("problem convert XML body to struct. Error: %s", err)
				errCounter++
				continue
			}
			h.counterList.group[g].counterName[c].description = description.QueryCounterDescriptionReturn
		}
	}
	if errCounter > 0 {
		return fmt.Errorf("in read description for counters get %d errors", errCounter)
	}
	return nil

}

func (h *ClusterHostMonitorData) print() string {

	return ""
}

func (h *ClusterHostMonitorData) string() string {
	return fmt.Sprintf("Errors in %s : %d", applicationName, monitors.client.responseErrors)
}

func (h *ClusterHostMonitorData) logFields(operation ...string) log.Fields {
	var f log.Fields
	if len(operation) == 2 {
		f = log.Fields{
			FieldMonitorName: h.server,
			FieldRoutine:     operation[0],
			FieldSessionId:   operation[1],
		}
	} else if len(operation) == 1 {
		f = log.Fields{
			FieldMonitorName: h.server,
			FieldRoutine:     operation[0],
		}
	} else {
		f = log.Fields{
			FieldMonitorName: h.server,
		}
	}
	return f
}
