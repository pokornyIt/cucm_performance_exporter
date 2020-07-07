package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type perfMonService struct {
	monitors []perfMonHost // host monitor parts
	client   *perfClient   // http client
}

type openSessionResponse struct {
	XMLName       xml.Name `xml:"perfmonOpenSessionResponse"`
	Text          string   `xml:",chardata"`
	Ns1           string   `xml:"ns1,attr"`
	OpenSessionId string   `xml:"perfmonOpenSessionReturn"`
}

func NewPerfMonServers() *perfMonService {
	p := perfMonService{
		monitors: make([]perfMonHost, 0),
		client:   NewPerfClient(),
	}
	for _, r := range config.MonitorNames {
		p.monitors = append(p.monitors, *NewPerMonHost(r))
	}
	log.WithFields(p.logFields("NewPerfMonServers")).Trace("create monitor service")
	return &p
}

func (s *perfMonService) OpenSession() (err error) {
	log.WithFields(s.logFields("OpenSession")).Trace("open new session")
	req := " <soap:perfmonOpenSession/>"
	body, err := s.client.processRequest("OpenSession", req)

	var data openSessionResponse
	err = xml.Unmarshal([]byte(body), &data)
	if err != nil {
		log.WithFields(s.logFields("OpenSession")).Errorf("problem convert XML body to struct. Error: %s", err)
		return err
	}
	s.client.session = data.OpenSessionId
	log.WithFields(s.logFields("OpenSession", data.OpenSessionId)).Infof("open new monitoring session")
	return nil
}

func (s *perfMonService) AddCounters() {
	log.WithFields(s.logFields("AddCounters", s.client.session)).Trace("register counters for session")
	if !s.client.isSessionOpen() {
		log.WithFields(s.logFields("AddCounters")).Debug("session not open")
		return
	}
	cnt := 0
	for _, mon := range s.monitors {
		if mon.AddCounters(s.client) != nil {
			log.WithFields(s.logFields("AddCounters", s.client.session)).Errorf("problem register counter toi session for monitor %s", mon.server)
		} else {
			cnt++
		}
	}
	if cnt == len(s.monitors) {
		log.WithFields(s.logFields("AddCounters", s.client.session)).Info("success register counters for session")
	}
}

func (s *perfMonService) CloseSession() {
	log.WithFields(s.logFields("CloseSession")).Trace("close existing session")
	if !s.client.isSessionOpen() {
		log.WithFields(s.logFields("CloseSession")).Debug("not any open session")
	}
	req := fmt.Sprintf("<soap:perfmonCloseSession><soap:SessionHandle>%s</soap:SessionHandle></soap:perfmonCloseSession>", s.client.session)
	_, _ = s.client.processRequest("CloseSession", req)
	log.WithFields(s.logFields("CloseSession", s.client.session)).Debug("current session is closed")
	s.client.session = ""
}

func (s *perfMonService) ExistSession() bool {
	return s.client.isSessionOpen()
}

func (s *perfMonService) CollectSessionData() (err error) {
	log.WithFields(s.logFields("CollectSessionData", s.client.session)).Trace("collect session data")
	if !s.client.isSessionOpen() {
		log.WithFields(s.logFields("CollectSessionData")).Debug("session not open")
		return errors.New("session not exist for open data")
	}
	req := fmt.Sprintf("<soap:perfmonCollectSessionData><soap:SessionHandle>%s</soap:SessionHandle></soap:perfmonCollectSessionData>", s.client.session)
	body, err := s.client.processRequest("CollectSessionData", req)

	var data SessionData
	err = xml.Unmarshal([]byte(body), &data)
	if err != nil {
		log.WithFields(s.logFields("CollectSessionData", s.client.session)).Errorf("problem convert XML body to required struct. Error: %s", err)
		return err
	}
	data.processData()
	return nil
}

func (s *perfMonService) ListAllCounters() (err error) {
	log.WithFields(s.logFields("ListAllCounters")).Trace("collect all counters")
	for r, _ := range s.monitors {
		e := s.monitors[r].ListCounters(s.client)
		if e != nil {
			err = e
		}
		e = s.monitors[r].ReadCounterDescription(s.client)
		if e != nil {
			err = e
		}
	}
	return err
}

func (s *perfMonService) GetCounterDetails(name string) (details *counterDetails, err error) {
	for srv, server := range s.monitors {
		for g, group := range server.counterList.group {
			for c, counter := range group.counterName {
				if counter.name == name {
					return &s.monitors[srv].counterList.group[g].counterName[c], nil
				}
			}
		}
	}
	log.WithFields(s.logFields("GetCounterDetails")).Errorf("not found details for counter %s", name)
	details = &counterDetails{name: name, description: fmt.Sprintf("Description for %s not exists", name)}
	return details, errors.New(fmt.Sprintf("problem found required counter [%s] on any server", name))
}

func (s *perfMonService) print() string {

	return ""
}
func (s *perfMonService) logFields(operation ...string) log.Fields {
	names := strings.Builder{}
	delimiter := ""
	for _, r := range s.monitors {
		names.WriteString(delimiter)
		names.WriteString(r.server)
		delimiter = ";"
	}
	var f log.Fields
	if len(operation) == 2 {
		f = log.Fields{
			"monitorNames": names.String(),
			"operation":    operation[0],
			"sessionId":    operation[1],
		}
	} else if len(operation) == 1 {
		f = log.Fields{
			"monitorNames": names.String(),
			"operation":    operation,
		}
	} else {
		f = log.Fields{
			"monitorNames": names.String(),
		}
	}

	return f
}
