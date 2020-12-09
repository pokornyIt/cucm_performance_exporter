package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
)

type FaultResponse struct {
	XMLName     xml.Name `xml:"Fault"`
	Text        string   `xml:",chardata"`
	FaultCode   string   `xml:"faultcode"`
	FaultString string   `xml:"faultstring"`
	Detail      string   `xml:"detail"`
}

func perfRequestCreate(requestId string, body string) (req *http.Request, err error) {
	log.WithFields(log.Fields{"routine": "perfRequestCreate", "requestId": requestId}).Trace("prepare request")
	server := fmt.Sprintf("https://%s:8443/perfmonservice2/services/PerfmonService?wsdl", config.ApiAddress)
	log.WithFields(log.Fields{"routine": "perfRequestCreate", "requestId": requestId}).Tracef("prepare server API name: %s", server)
	req, err = http.NewRequest("POST", server, bytes.NewBuffer([]byte(body)))
	if err != nil {
		log.WithField("routine", "perfRequestCreate").Errorf("problem create request. Error: %s", err)
		return nil, err
	}
	req.Header.Add("User-Agent", httpApplication)
	req.Header.Add("Content-Type", "text/xml")
	req.Header.Add("Accept", "text/xml")
	req.Header.Add("Cache-Control", "no-cache")
	req.SetBasicAuth(config.ApiUser, config.ApiPassword)

	return req, nil
}

func perfRequestResponse(requestId string, client *http.Client, req *http.Request) (body string, resp *http.Response, err error) {
	log.WithFields(log.Fields{"routine": "perfRequestResponse", "requestId": requestId}).Trace("get response")
	resp, err = client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{"routine": "perfRequestResponse", "requestId": requestId}).Errorf("problem process request. Error: %s", err)
		return "", resp, err
	}
	s, err := ioutil.ReadAll(resp.Body)
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
