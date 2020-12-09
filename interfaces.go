package main

import log "github.com/sirupsen/logrus"

type status interface {
	print() string
	logFields(operation ...string) log.Fields
}

type configValid interface {
	Validate() (err error)
	Print() string
}
