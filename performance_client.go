package main

import (
	"context"
	"crypto/tls"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type PerfClient struct {
	client  *http.Client // reference to exist HTTP Client
	session string       // session
	//lastRequest    time.Time    // last request
	requests       uint64 // success created request
	responses      uint64 // success obtains responses
	responseErrors uint64 // error obtain responses
}

func NewPerfClient() *PerfClient {
	var cp PerfClient
	if config.IgnoreCertificate {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		cp = PerfClient{
			client:         &http.Client{Transport: tr},
			requests:       0,
			responses:      0,
			responseErrors: 0,
			session:        "",
		}
	} else {
		cp = PerfClient{
			client:         &http.Client{},
			requests:       0,
			responses:      0,
			responseErrors: 0,
			session:        "",
		}
	}
	return &cp
}

func (p *PerfClient) processRequest(name string, inner string) (body string, err error) {
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
		p.responseErrors++
		return fmt.Sprintf("%d", resp.StatusCode), fmt.Errorf("response status is %s", resp.Status)
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

func (p *PerfClient) isSessionOpen() bool {
	return len(p.session) > 0
}

func (p *PerfClient) logFields(operation ...string) log.Fields {
	var f log.Fields
	if len(operation) == 1 {
		f = log.Fields{
			"session": p.session,
			Routine:   operation,
		}
	} else {
		f = log.Fields{
			"session": p.session,
		}
	}
	return f
}
