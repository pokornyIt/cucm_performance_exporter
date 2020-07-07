package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type perfClient struct {
	client         *http.Client // reference to exist HTTP Client
	session        string       // session
	lastRequest    time.Time    // last request
	requests       uint64       // success created request
	responses      uint64       // success obtains responses
	responseErrors uint64       // error obtain responses
}

func NewPerfClient() *perfClient {
	var cp perfClient
	if config.IgnoreCertificate {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		cp = perfClient{
			client:         &http.Client{Transport: tr},
			requests:       0,
			responses:      0,
			responseErrors: 0,
			session:        "",
		}
	} else {
		cp = perfClient{
			client:         &http.Client{},
			requests:       0,
			responses:      0,
			responseErrors: 0,
			session:        "",
		}
	}
	return &cp
}

func (p *perfClient) processRequest(name string, inner string) (body string, err error) {
	s := fmt.Sprintf(Envelope, inner)
	requestId := RandomString()
	req, err := perfRequestCreate(requestId, s)
	p.requests++
	if err != nil {
		log.WithFields(p.logFields(name)).Errorf("problem prepare %s request. Error: %s", name, err)
		return "", err
	}
	body, resp, err := perfRequestResponse(requestId, p.client, req)
	if resp != nil && resp.StatusCode > 299 {
		log.WithFields(p.logFields(name)).Errorf("problem read %s response. Status code: %s", name, resp.Status)
		p.responseErrors++
		return fmt.Sprintf("%d", resp.StatusCode), errors.New(fmt.Sprintf("response status is %s", resp.Status))
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

func (p *perfClient) isSessionOpen() bool {
	return len(p.session) > 0
}

func (p *perfClient) logFields(operation ...string) log.Fields {
	var f log.Fields
	if len(operation) == 1 {
		f = log.Fields{
			"session":   p.session,
			"operation": operation,
		}
	} else {
		f = log.Fields{
			"session": p.session,
		}
	}
	return f
}
