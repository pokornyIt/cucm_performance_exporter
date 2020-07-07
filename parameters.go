package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"strings"
)

type Config struct {
	MonitorNames      []string       `yaml:"monitor_names" json:"monitor_names"`
	Metrics           MetricsEnabled `yaml:"metrics" json:"metrics"`
	ApiAddress        string         `yaml:"apiAddress" json:"apiAddress"`
	ApiUser           string         `yaml:"apiUser" json:"apiUser"`
	ApiPassword       string         `yaml:"apiPwd" json:"apiPwd"`
	Port              uint16         `yaml:"port" json:"port"`
	IgnoreCertificate bool           `yaml:"ignoreCertificate" json:"ignoreCertificate"`
}

type MetricsEnabled struct {
	CallsActive              bool `yaml:"callsActive" json:"callsActive"`
	CallsInProgress          bool `yaml:"callsInProgress" json:"callsInProgress"`
	CallsCompleted           bool `yaml:"callsCompleted" json:"callsCompleted"`
	PartiallyRegisteredPhone bool `yaml:"partiallyRegisteredPhone" json:"partiallyRegisteredPhone"`
	RegisteredHardwarePhones bool `yaml:"registeredHardwarePhones" json:"registeredHardwarePhones"`
	GatewaysSessionsActive   bool `yaml:"gatewaysSessionsActive" json:"gatewaysSessionsActive"`
	GatewaysSessionsFailed   bool `yaml:"gatewaysSessionsFailed" json:"gatewaysSessionsFailed"`
	PhoneSessionsActive      bool `yaml:"phoneSessionsActive" json:"phoneSessionsActive"`
	PhoneSessionsFailed      bool `yaml:"phoneSessionsFailed" json:"phoneSessionsFailed"`
}

var (
	showConfig = kingpin.Flag("config.show", "Show actual configuration and ends").Default("false").Bool()
	configFile = kingpin.Flag("config.file", "Configuration file default is \"server.yml\".").PlaceHolder("cfg.yml").Default("server.yml").String()
	//listenAddress = kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").Default(":9717").String()
	config = &Config{
		Metrics: MetricsEnabled{
			CallsActive:              true,
			CallsInProgress:          true,
			CallsCompleted:           true,
			PartiallyRegisteredPhone: true,
			RegisteredHardwarePhones: true,
			GatewaysSessionsActive:   true,
			GatewaysSessionsFailed:   true,
			PhoneSessionsActive:      true,
			PhoneSessionsFailed:      true,
		},
		MonitorNames:      []string{},
		ApiAddress:        "",
		ApiUser:           "",
		ApiPassword:       "",
		Port:              9717,
		IgnoreCertificate: false,
	}
	apiServer = kingpin.Flag("api.address", "CUCM Server FQDN or IP address.").PlaceHolder("server").Default("").String()
	apiUser   = kingpin.Flag("api.user", "CUCM user with access to PerfMON data.").PlaceHolder("User").Default("").String()
	apiPwd    = kingpin.Flag("api.pwd", "CUCM user password.").PlaceHolder("pwd").Default("").String()
)

func (c *Config) LoadFile(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.UnmarshalStrict(content, c)
	if err != nil {
		err = json.Unmarshal(content, c)
		if err != nil {
			return err
		}
	}
	if len(*apiServer) > 0 {
		c.ApiAddress = *apiServer
	}
	if len(*apiUser) > 0 {
		c.ApiUser = *apiUser
	}
	if len(*apiPwd) > 0 {
		c.ApiPassword = *apiPwd
	}

	match, err := regexp.MatchString("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$", c.ApiAddress)

	if !match || err != nil {
		match, err = regexp.MatchString("^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$", c.ApiAddress)
		if !match || err != nil {
			return errors.New("API Address isn't valid FQDN or IP address")
		}
	}
	if len(c.ApiUser) < 1 {
		return errors.New("API User must be defined")
	}
	if len(c.ApiPassword) < 1 {
		return errors.New("API User password must be defined")
	}
	if c.Port < 1024 || c.Port > 65535 {
		return errors.New("defined port not valid")
	}
	return nil
}

func (c *Config) print() string {
	a := fmt.Sprintf("API:               [https://%s:8443/perfmonservice2/services/PerfmonService?wsdl]\r\n", c.ApiAddress)
	a = fmt.Sprintf("%sIgnoreCertificate: [%t]\r\n", a, c.IgnoreCertificate)
	a = fmt.Sprintf("%sUser:              [%s]\r\n", a, c.ApiUser)
	a = fmt.Sprintf("%sServers:           [%s]\r\n", a, strings.Join(c.MonitorNames, ", "))
	a = fmt.Sprintf("%sPort:              [:%d]\r\n", a, c.Port)
	a = fmt.Sprintf("%sMetrics:\r\n", a)
	lenTxt := 0
	for _, m := range AllowedCounterNames {
		if len(m.allowedCounterName) > lenTxt {
			lenTxt = len(m.allowedCounterName)
		}
	}
	var reqSpaces int
	for _, m := range AllowedCounterNames {
		reqSpaces = lenTxt - len(m.allowedCounterName)
		a = fmt.Sprintf("%s\t- %s:%s [%t]\r\n", a, m.allowedCounterName, strings.Repeat(" ", reqSpaces), c.Metrics.enablePrometheusCounter(m.allowedCounterName))
	}
	return a
}

func (c *Config) logFields(operation ...string) log.Fields {
	return log.Fields{
		"monitorNames":      strings.Join(c.MonitorNames, ";"),
		"apiUser":           c.ApiUser,
		"apiAddress":        c.ApiAddress,
		"telemetryPort":     c.Port,
		"ignoreCertificate": c.IgnoreCertificate,
	}
}

func (m *MetricsEnabled) enablePrometheusCounter(name string) bool {
	if name == CallsActive {
		return m.CallsActive
	}
	if name == CallsInProgress {
		return m.CallsInProgress
	}
	if name == CallsCompleted {
		return m.CallsCompleted
	}
	if name == GatewaysSessionsActive {
		return m.GatewaysSessionsActive
	}
	if name == GatewaysSessionsFailed {
		return m.GatewaysSessionsFailed
	}
	if name == RegisteredHardwarePhones {
		return m.RegisteredHardwarePhones
	}
	if name == PartiallyRegisteredPhone {
		return m.PartiallyRegisteredPhone
	}
	if name == PhoneSessionsActive {
		return m.PhoneSessionsActive
	}
	if name == PhoneSessionsFailed {
		return m.PhoneSessionsFailed
	}

	return false
}
