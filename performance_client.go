package main

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// ApiMonitorClient API client
type ApiMonitorClient struct {
	client         *http.Client // client reference to exist HTTP Client
	session        string       // session actual id
	requests       uint64       // requests success created request
	responses      uint64       // responses success obtains response
	responseErrors uint64       // responseErrors error obtain response
}

// NewApiMonitorClient create new API client with prepared http.Client
func NewApiMonitorClient() *ApiMonitorClient {
	var cp ApiMonitorClient
	if config.IgnoreCertificate {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		cp = ApiMonitorClient{
			client:         &http.Client{Transport: tr},
			requests:       0,
			responses:      0,
			responseErrors: 0,
			session:        "",
		}
	} else {
		cp = ApiMonitorClient{
			client:         &http.Client{},
			requests:       0,
			responses:      0,
			responseErrors: 0,
			session:        "",
		}
	}
	return &cp
}

// processRequest process one request to API with predefined timeout
// program returns collected body or error if here any problem
func (p *ApiMonitorClient) processRequest(name string, inner string) (body string, err error) {
	if LogRequestDuration {
		defer duration(track(log.Fields{FieldRoutine: "processRequest"}, "procedure ends"))
	}
	var req *http.Request
	var resp *http.Response
	s := fmt.Sprintf(Envelope, inner)
	requestId := RandomString()
	req, err = perfRequestCreate(requestId, s)
	p.requests++
	if err != nil {
		log.WithFields(p.logFields(name)).Errorf("problem prepare %s request. Error: %s", name, err)
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.ApiTimeout)*time.Second)
	defer cancel()
	req = req.WithContext(ctx)
	body, resp, err = perfRequestResponse(requestId, p.client, req)
	if resp != nil && resp.StatusCode > 299 {
		log.WithFields(p.logFields(name)).Errorf("problem read %s response. Status code: %s", name, resp.Status)
		var f FaultResponse
		err = xml.Unmarshal([]byte(body), &f)
		p.responseErrors++
		if resp.StatusCode == 401 || err != nil {
			return fmt.Sprintf("%d", resp.StatusCode), fmt.Errorf("response status is %s - %s %s", resp.Status, f.FaultCode, f.FaultString)
		}
		return f.FaultCode, fmt.Errorf("response status is %s - %s %s", resp.Status, f.FaultCode, f.FaultString)
	}

	if err != nil {
		log.WithFields(p.logFields(name)).Errorf("problem read %s response. Error: %s", name, err)
		p.responseErrors++
		return "", err
	}
	p.responses++
	body, err = perfRequestBodyRelevant(body)
	if err != nil {
		log.WithFields(p.logFields(name)).Errorf("problem analyze %s response. Error: %s", name, err)
		return "", err
	}
	return body, nil
}

// isSessionOpen Define if connection is UP
func (p *ApiMonitorClient) isSessionOpen() bool {
	return len(p.session) > 0
}

// logFields create valid list of log fields depend on server
func (p *ApiMonitorClient) logFields(operation ...string) log.Fields {
	var f log.Fields
	if len(operation) == 1 {
		f = log.Fields{
			FieldSession: p.session,
			FieldIsUp:    p.isSessionOpen(),
			FieldRoutine: operation,
		}
	} else {
		f = log.Fields{
			FieldIsUp:    p.isSessionOpen(),
			FieldSession: p.session,
		}
	}
	return f
}

// print List actual error status of client
func (p *ApiMonitorClient) print() string {
	msg := "Client status"
	msg = fmt.Sprintf("%s\r\nRequests     %d", msg, p.requests)
	msg = fmt.Sprintf("%s\r\nResponses    %d", msg, p.responses)
	msg = fmt.Sprintf("%s\r\nError        %d", msg, p.responseErrors)
	msg = fmt.Sprintf("%s\r\nConnected    %t", msg, p.isSessionOpen())
	return msg
}
