package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"regexp"
	"time"
)

type FaultResponse struct {
	XMLName     xml.Name `xml:"Fault"`
	Text        string   `xml:",chardata"`
	FaultCode   string   `xml:"faultcode"`
	FaultString string   `xml:"faultstring"`
	Detail      string   `xml:"detail"`
}

// perfRequestCreate generate http request wit request ID and body
//   - request to https://<API server>:8443/perfmonservice2/services/PerfmonService?wsdl
func perfRequestCreate(requestId string, body string) (req *http.Request, err error) {
	log.WithFields(log.Fields{FieldRoutine: "perfRequestCreate", FieldRequestId: requestId}).Trace("prepare request")
	if LogRequestDuration {
		defer duration(track(log.Fields{FieldRoutine: "perfRequestCreate", FieldRequestId: requestId}, "procedure ends"))
	}
	server := fmt.Sprintf("https://%s:8443/perfmonservice2/services/PerfmonService?wsdl", config.ApiAddress)
	log.WithFields(log.Fields{FieldRoutine: "perfRequestCreate", FieldRequestId: requestId}).Tracef("prepare server API name: %s", server)
	req, err = http.NewRequest("POST", server, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.WithField(FieldRoutine, "perfRequestCreate").Errorf("problem create request. Error: %s", err)
		return nil, err
	}
	req.Header.Add("User-Agent", httpApplicationName())
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("Accept", "text/xml")
	req.Header.Add("Cache-Control", "no-cache")
	req.SetBasicAuth(config.ApiUser, config.ApiPassword)

	return req, nil
}

func perfRequestResponse(requestId string, client *http.Client, req *http.Request) (body string, resp *http.Response, err error) {
	log.WithFields(log.Fields{FieldRoutine: "perfRequestResponse", FieldRequestId: requestId}).Trace("get response")
	if LogRequestDuration {
		defer duration(track(log.Fields{FieldRoutine: "perfRequestResponse", FieldRequestId: requestId}, "procedure ends"))
	}
	requestsCount := rateRequest.requests
	waitTime := rateRequest.delay()
	if waitTime > time.Millisecond {
		log.WithFields(log.Fields{FieldRoutine: "perfRequestResponse", FieldRequestId: requestId}).
			Warnf("wait after %d requests for %s", requestsCount, waitTime.String())
		time.Sleep(waitTime)
	}
	resp, err = client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{FieldRoutine: "perfRequestResponse", FieldRequestId: requestId}).Errorf("problem process request. Error: %s", err)
		return "", resp, err
	}
	s, err := io.ReadAll(resp.Body)
	return string(s), resp, err
}

func perfRequestBodyRelevant(body string) (data string, err error) {
	var rex = regexp.MustCompile(`(?m)<soapenv:Body>((.|\n)*?)</soapenv:Body>`)

	if !rex.Match([]byte(body)) {
		return "", errors.New("response body not contains \"<soapenv:Body>\"")
	}
	x := rex.FindStringSubmatch(body)
	if len(x) < 2 {
		return "", errors.New("response body not contains \"<soapenv:Body>\"")
	}
	data = x[1]
	return data, nil
}
