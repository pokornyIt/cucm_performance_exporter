package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"regexp"
	"strings"
)

type Config struct {
	MonitorNames      []string       `yaml:"monitor_names" json:"monitor_names"`
	Metrics           MetricsEnabled `yaml:"metrics" json:"metrics"`
	Log               ConfigLog      `yaml:"log" json:"log"`
	ApiAddress        string         `yaml:"apiAddress" json:"apiAddress"`
	ApiUser           string         `yaml:"apiUser" json:"apiUser"`
	ApiPassword       string         `yaml:"apiPwd" json:"apiPwd"`
	Port              int            `yaml:"port" json:"port"`
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
	GoCollector              bool `yaml:"goCollector" json:"goCollector"`
	ProcessStatus            bool `yaml:"processStatus" json:"processStatus"`
}

type ConfigLog struct {
	Level          string `json:"level" yaml:"level"`                   // Log level FATAL, ERROR, WARNING, INFO, DEBUG, TRACE. Default is INFO
	FileName       string `json:"fileName" yaml:"fileName"`             // Log filename
	JSONFormat     bool   `json:"jsonFormat" yaml:"jsonFormat"`         // enable log in JSON format
	LogProgramInfo bool   `json:"logProgramInfo" yaml:"logProgramInfo"` // enable log program details (line, file name)
	MaxSize        int    `json:"maxSize" yaml:"maxSize"`               // Maximal log file size in MB
	MaxBackups     int    `json:"maxBackups" yaml:"maxBackups"`         // Maximal Number of backups
	MaxAge         int    `json:"maxAge" yaml:"maxAge"`                 // Maximal backup in days
	Quiet          bool   `json:"quiet" yaml:"quiet"`                   // Logging quiet - output only to file or only panic
}

type Intervals struct {
	Default int
	Min     int
	Max     int
}

var (
	showConfig    = kingpin.Flag("config.show", "Show actual configuration and ends").Default("false").Bool()
	configFile    = kingpin.Flag("config.file", "Configuration file default is \"server.yml\".").PlaceHolder("cfg.yml").Default("server.yml").String()
	LogMaxSize    = Intervals{Default: 50, Min: 1, Max: 5000}       // Limits and defaults for Log MaxSize
	LogMaxBackups = Intervals{Default: 5, Min: 0, Max: 100}         // Limits and defaults for Log MaxBackups
	LogMaxAge     = Intervals{Default: 30, Min: 1, Max: 365}        // Limits and defaults for Log MaxAge
	portLimits    = Intervals{Default: 9717, Min: 1024, Max: 65535} // Limits and defaults for ports

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
			GoCollector:              true,
			ProcessStatus:            true,
		},
		Log: ConfigLog{
			Level:          "Info",
			FileName:       "",
			JSONFormat:     false,
			LogProgramInfo: false,
			MaxSize:        LogMaxSize.Default,
			MaxBackups:     LogMaxBackups.Default,
			MaxAge:         LogMaxAge.Default,
			Quiet:          false,
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
	logFile   = kingpin.Flag("log.file", "Path and file name for store log. Default is disabled.").PlaceHolder("file.log").Default("").String()
)

// LoadFile Load configuration file form filename
func (c *Config) LoadFile(filename string) (err error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return c.ProcessLoadFile(content)
}

// ProcessLoadFile Process config content/*
func (c *Config) ProcessLoadFile(content []byte) (err error) {
	err = yaml.UnmarshalStrict(content, c)
	if err != nil {
		err1 := json.Unmarshal(content, c)
		if err1 != nil {
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
	if len(*logFile) > 0 {
		c.Log.FileName = *logFile
	}

	return c.Validate()
}

func (c *Config) Validate() (err error) {
	// validate master part
	if !validServer(c.ApiAddress) {
		return errors.New("API Address isn't valid FQDN or IP address")
	}
	if len(c.ApiUser) < 1 {
		return errors.New("API User must be defined")
	}
	if len(c.ApiPassword) < 1 {
		return errors.New("API User password must be defined")
	}
	if !portLimits.Validate(c.Port) {
		return errors.New("defined port not valid")
	}

	// validate child
	if err = c.Log.Validate(); err != nil {
		return err
	}
	return nil
}

func (m *MetricsEnabled) Validate() bool {
	return true
}

func (m *MetricsEnabled) Print() string {
	a := "Metrics:\r\n"
	lenTxt := len("ProcessStatus")
	const fmtFormat = "%s\t- %s:%s [%t]\r\n"
	for _, name := range SupportedCounters {
		if len(name.allowedCounterName) > lenTxt {
			lenTxt = len(name.allowedCounterName)
		}
	}
	var reqSpaces int
	for _, name := range SupportedCounters {
		reqSpaces = lenTxt - len(name.allowedCounterName)
		a = fmt.Sprintf(fmtFormat, a, name.allowedCounterName, strings.Repeat(" ", reqSpaces), m.enablePrometheusCounter(name.allowedCounterName))
	}
	reqSpaces = lenTxt - len("GoCollector")
	a = fmt.Sprintf(fmtFormat, a, "GoCollector", strings.Repeat(" ", reqSpaces), m.GoCollector)
	reqSpaces = lenTxt - len("ProcessStatus")
	a = fmt.Sprintf(fmtFormat, a, "ProcessStatus", strings.Repeat(" ", reqSpaces), m.ProcessStatus)

	return a
}

func (c *Config) print() string {
	a := fmt.Sprintf("API:                  [https://%s:8443/perfmonservice2/services/PerfmonService?wsdl]\r\n", c.ApiAddress)
	a = fmt.Sprintf("%sIgnore Certificate:   [%t]\r\n", a, c.IgnoreCertificate)
	a = fmt.Sprintf("%sUser:                 [%s]\r\n", a, c.ApiUser)
	a = fmt.Sprintf("%sServers:              [%s]\r\n", a, strings.Join(c.MonitorNames, ", "))
	a = fmt.Sprintf("%sPort:                 [:%d]\r\n", a, c.Port)

	a = fmt.Sprintf("%s%s", a, c.Metrics.Print())
	a = fmt.Sprintf("%s%s", a, c.Log.Print())
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

func (a *ConfigLog) Validate() (err error) {
	lvl := validLogLevel(a.Level)
	a.Level = strings.ToUpper(lvl.String())
	a.FileName = FixFileName(a.FileName)
	a.MaxSize = LogMaxSize.ValidOrDefault(a.MaxSize)
	a.MaxAge = LogMaxSize.ValidOrDefault(a.MaxAge)
	a.MaxBackups = LogMaxSize.ValidOrDefault(a.MaxBackups)

	return nil
}

func (a *ConfigLog) Print() string {
	o := "Logging\r\n"
	o = fmt.Sprintf("%s\t- Level                     [%s]\r\n", o, a.Level)
	o = fmt.Sprintf("%s\t- Use JSON format           [%t]\r\n", o, a.JSONFormat)
	if len(a.FileName) > 0 {
		o = fmt.Sprintf("%s\t- Logging file              [%s]\r\n", o, a.FileName)
		o = fmt.Sprintf("%s\t- Logging program details   [%t]\r\n", o, a.LogProgramInfo)
		o = fmt.Sprintf("%s\t- Maximal file size in MB   [%d]\r\n", o, a.MaxSize)
		o = fmt.Sprintf("%s\t- Number of backups         [%d]\r\n", o, a.MaxBackups)
		o = fmt.Sprintf("%s\t- Maximal age in days       [%d]\r\n", o, a.MaxAge)
		o = fmt.Sprintf("%s\t- Backup compress           [%t]\r\n", o, true)
	} else {
		o = fmt.Sprintf("%s\t- Don't logging to file\r\n", o)
	}
	return o
}

func (a *ConfigLog) LogToFile() bool {
	return len(a.FileName) > 0
}

func (i *Intervals) Validate(actual int) bool {
	return actual >= i.Min && i.Max >= actual
}

func (i *Intervals) ValidOrDefault(actual int) int {
	if i.Validate(actual) {
		return actual
	}
	return i.Default
}

func (i *Intervals) Print() string {
	return fmt.Sprintf("%d (min: %d / max:%d)", i.Default, i.Min, i.Max)
}

func IsValidFileName(file string) bool {
	file = strings.Trim(file, " ")
	if len(file) < 1 {
		return false
	}
	file = strings.ReplaceAll(file, "\\", "/")
	_, f := path.Split(file)
	if f == "" || f == "." || f == ".." {
		return false
	}
	return true
}

func FixFileName(file string) string {
	if !IsValidFileName(file) {
		return ""
	}
	file = strings.ReplaceAll(file, "\\", "/")
	dir, file1 := path.Split(file)
	dir = path.Clean(dir)
	file1 = path.Join(dir, file1)
	return strings.ReplaceAll(file1, "/", string(os.PathSeparator))
}

func validServer(srv string) bool {
	ipAddress := regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	invalidAddress := regexp.MustCompile(`^((\d+)\.){3}(\d+)$`)
	dnsName := regexp.MustCompile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	if ipAddress.MatchString(srv) {
		return true
	}
	if invalidAddress.MatchString(srv) {
		return false
	}
	return dnsName.MatchString(srv)
}
