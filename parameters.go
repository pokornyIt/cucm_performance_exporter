package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"regexp"
	"strings"
)

type Config struct {
	MonitorNames        []string       `yaml:"monitor_names" json:"monitor_names"`
	Metrics             MetricsEnabled `yaml:"metrics" json:"metrics"`
	Log                 ConfigLog      `yaml:"log" json:"log"`
	ApiAddress          string         `yaml:"apiAddress" json:"apiAddress"`
	ApiUser             string         `yaml:"apiUser" json:"apiUser"`
	ApiPassword         string         `yaml:"apiPwd" json:"apiPwd"`
	Port                int            `yaml:"port" json:"port"`
	IgnoreCertificate   bool           `yaml:"ignoreCertificate" json:"ignoreCertificate"`
	ApiTimeout          int            `yaml:"apiTimeout" json:"apiTimeout"`
	AllowStop           bool           `yaml:"allowStop" json:"allowStop"`
	SleepBetweenRequest int            `yaml:"sleepBetweenRequest" json:"sleepBetweenRequest"`
}

type MetricsEnabled struct {
	CallsActive                               bool `yaml:"callsActive" json:"callsActive"`
	CallsAttempted                            bool `yaml:"callsAttempted" json:"callsAttempted"`
	CallsInProgress                           bool `yaml:"callsInProgress" json:"callsInProgress"`
	CallsCompleted                            bool `yaml:"callsCompleted" json:"callsCompleted"`
	PartiallyRegisteredPhone                  bool `yaml:"partiallyRegisteredPhone" json:"partiallyRegisteredPhone"`
	RegisteredHardwarePhones                  bool `yaml:"registeredHardwarePhones" json:"registeredHardwarePhones"`
	GatewayRegistrationFailures               bool `yaml:"gatewayRegistrationFailures" json:"gatewayRegistrationFailures"`
	GatewaysInService                         bool `yaml:"gatewaysInService" json:"gatewaysInService"`
	GatewaysOutOfService                      bool `yaml:"gatewaysOutOfService" json:"gatewaysOutOfService"`
	GatewaysSessionsActive                    bool `yaml:"gatewaysSessionsActive" json:"gatewaysSessionsActive"`
	GatewaysSessionsFailed                    bool `yaml:"gatewaysSessionsFailed" json:"gatewaysSessionsFailed"`
	PhoneSessionsActive                       bool `yaml:"phoneSessionsActive" json:"phoneSessionsActive"`
	PhoneSessionsFailed                       bool `yaml:"phoneSessionsFailed" json:"phoneSessionsFailed"`
	GoCollector                               bool `yaml:"goCollector" json:"goCollector"`
	ProcessStatus                             bool `yaml:"processStatus" json:"processStatus"`
	AnnunciatorOutOfResources                 bool `yaml:"annunciatorOutOfResources" json:"annunciatorOutOfResources"`
	AnnunciatorResourceActive                 bool `yaml:"annunciatorResourceActive" json:"annunciatorResourceActive"`
	AnnunciatorResourceAvailable              bool `yaml:"annunciatorResourceAvailable" json:"annunciatorResourceAvailable"`
	AnnunciatorResourceTotal                  bool `yaml:"annunciatorResourceTotal" json:"annunciatorResourceTotal"`
	AuthenticatedCallsActive                  bool `yaml:"authenticatedCallsActive" json:"authenticatedCallsActive"`
	AuthenticatedCallsCompleted               bool `yaml:"authenticatedCallsCompleted" json:"authenticatedCallsCompleted"`
	AuthenticatedPartiallyRegisteredPhone     bool `yaml:"authenticatedPartiallyRegisteredPhone" json:"authenticatedPartiallyRegisteredPhone"`
	AuthenticatedRegisteredPhones             bool `yaml:"authenticatedRegisteredPhones" json:"authenticatedRegisteredPhones"`
	CallManagerHeartBeat                      bool `yaml:"callManagerHeartBeat" json:"callManagerHeartBeat"`
	CumulativeAllocatedResourceCannotOpenPort bool `yaml:"cumulativeAllocatedResourceCannotOpenPort" json:"cumulativeAllocatedResourceCannotOpenPort"`
	EncryptedCallsActive                      bool `yaml:"encryptedCallsActive" json:"encryptedCallsActive"`
	EncryptedCallsCompleted                   bool `yaml:"encryptedCallsCompleted" json:"encryptedCallsCompleted"`
	EncryptedPartiallyRegisteredPhones        bool `yaml:"encryptedPartiallyRegisteredPhones" json:"encryptedPartiallyRegisteredPhones"`
	EncryptedRegisteredPhones                 bool `yaml:"encryptedRegisteredPhones" json:"encryptedRegisteredPhones"`
	HWConferenceActive                        bool `yaml:"hwConferenceActive" json:"hwConferenceActive"`
	HWConferenceCompleted                     bool `yaml:"hwConferenceCompleted" json:"hwConferenceCompleted"`
	HWConferenceOutOfResources                bool `yaml:"hwConferenceOutOfResources" json:"hwConferenceOutOfResources"`
	HWConferenceResourceActive                bool `yaml:"hwConferenceResourceActive" json:"hwConferenceResourceActive"`
	HWConferenceResourceAvailable             bool `yaml:"hwConferenceResourceAvailable" json:"hwConferenceResourceAvailable"`
	HWConferenceResourceTotal                 bool `yaml:"hwConferenceResourceTotal" json:"hwConferenceResourceTotal"`
	MTPOutOfResources                         bool `yaml:"mtpOutOfResources" json:"mtpOutOfResources"`
	MTPRequestsThrottled                      bool `yaml:"mtpRequestsThrottled" json:"mtpRequestsThrottled"`
	MTPResourceActive                         bool `yaml:"mtpResourceActive" json:"mtpResourceActive"`
	MTPResourceAvailable                      bool `yaml:"mtpResourceAvailable" json:"mtpResourceAvailable"`
	MTPResourceTotal                          bool `yaml:"mtpResourceTotal" json:"mtpResourceTotal"`
	RegisteredAnalogAccess                    bool `yaml:"registeredAnalogAccess" json:"registeredAnalogAccess"`
	RegisteredMGCPGateway                     bool `yaml:"registeredMGCPGateway" json:"registeredMGCPGateway"`
	RegisteredOtherStationDevices             bool `yaml:"registeredOtherStationDevices" json:"registeredOtherStationDevices"`
	SIPLineServerAuthorizationChallenges      bool `yaml:"sipLineServerAuthorizationChallenges" json:"sipLineServerAuthorizationChallenges"`
	SIPLineServerAuthorizationFailures        bool `yaml:"sipLineServerAuthorizationFailures" json:"sipLineServerAuthorizationFailures"`
	SIPTrunkApplicationAuthorizationFailures  bool `yaml:"sipTrunkApplicationAuthorizationFailures" json:"sipTrunkApplicationAuthorizationFailures"`
	SIPTrunkApplicationAuthorizations         bool `yaml:"sipTrunkApplicationAuthorizations" json:"sipTrunkApplicationAuthorizations"`
	SIPTrunkAuthorizationFailures             bool `yaml:"sipTrunkAuthorizationFailures" json:"sipTrunkAuthorizationFailures"`
	SIPTrunkAuthorizations                    bool `yaml:"sipTrunkAuthorizations" json:"sipTrunkAuthorizations"`
	SIPTrunkServerAuthenticationChallenges    bool `yaml:"sipTrunkServerAuthenticationChallenges" json:"sipTrunkServerAuthenticationChallenges"`
	SWConferenceActive                        bool `yaml:"swConferenceActive" json:"swConferenceActive"`
	SWConferenceCompleted                     bool `yaml:"swConferenceCompleted" json:"swConferenceCompleted"`
	SWConferenceOutOfResources                bool `yaml:"swConferenceOutOfResources" json:"swConferenceOutOfResources"`
	SWConferenceResourceActive                bool `yaml:"swConferenceResourceActive" json:"swConferenceResourceActive"`
	SWConferenceResourceAvailable             bool `yaml:"swConferenceResourceAvailable" json:"swConferenceResourceAvailable"`
	SWConferenceResourceTotal                 bool `yaml:"swConferenceResourceTotal" json:"swConferenceResourceTotal"`
	SystemCallsAttempted                      bool `yaml:"systemCallsAttempted" json:"systemCallsAttempted"`
	TranscoderOutOfResources                  bool `yaml:"transcoderOutOfResources" json:"transcoderOutOfResources"`
	TranscoderRequestsThrottled               bool `yaml:"transcoderRequestsThrottled" json:"transcoderRequestsThrottled"`
	TranscoderResourceActive                  bool `yaml:"transcoderResourceActive" json:"transcoderResourceActive"`
	TranscoderResourceAvailable               bool `yaml:"transcoderResourceAvailable" json:"transcoderResourceAvailable"`
	TranscoderResourceTotal                   bool `yaml:"transcoderResourceTotal" json:"transcoderResourceTotal"`
	UnEncryptedCallFailures                   bool `yaml:"unEncryptedCallFailures" json:"unEncryptedCallFailures"`
	VideoCallsActive                          bool `yaml:"videoCallsActive" json:"videoCallsActive"`
	VideoCallsCompleted                       bool `yaml:"videoCallsCompleted" json:"videoCallsCompleted"`
	VideoOnHoldOutOfResources                 bool `yaml:"videoOnHoldOutOfResources" json:"videoOnHoldOutOfResources"`
	VideoOnHoldResourceActive                 bool `yaml:"videoOnHoldResourceActive" json:"videoOnHoldResourceActive"`
	VideoOutOfResources                       bool `yaml:"videoOutOfResources" json:"videoOutOfResources"`
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
	showConfig               = kingpin.Flag("config.show", "Show actual configuration and ends").Default("false").Bool()
	configFile               = kingpin.Flag("config.file", "Configuration file default is \"server.yml\".").PlaceHolder("cfg.yml").Default("server.yml").String()
	LogMaxSize               = Intervals{Default: 50, Min: 1, Max: 5000}       // Limits and defaults for Log MaxSize
	LogMaxBackups            = Intervals{Default: 5, Min: 0, Max: 100}         // Limits and defaults for Log MaxBackups
	LogMaxAge                = Intervals{Default: 30, Min: 1, Max: 365}        // Limits and defaults for Log MaxAge
	PortLimits               = Intervals{Default: 9717, Min: 1024, Max: 65535} // Limits and defaults for ports
	ApiTimeoutLimit          = Intervals{Default: 5, Min: 1, Max: 30}          // Limits and defaults for API Timeouts in sec
	SleepBetweenRequestLimit = Intervals{Default: 30, Min: 5, Max: 120}        // Limits and defaults for sleep between API requests in  sec

	config = &Config{
		Metrics: MetricsEnabled{
			CallsActive:                               true,
			CallsInProgress:                           true,
			CallsCompleted:                            true,
			PartiallyRegisteredPhone:                  true,
			RegisteredHardwarePhones:                  true,
			GatewaysSessionsActive:                    true,
			GatewaysSessionsFailed:                    true,
			PhoneSessionsActive:                       true,
			PhoneSessionsFailed:                       true,
			GoCollector:                               true,
			ProcessStatus:                             true,
			CallsAttempted:                            false,
			GatewayRegistrationFailures:               false,
			GatewaysInService:                         false,
			GatewaysOutOfService:                      false,
			AnnunciatorOutOfResources:                 false,
			AnnunciatorResourceActive:                 false,
			AnnunciatorResourceAvailable:              false,
			AnnunciatorResourceTotal:                  false,
			AuthenticatedCallsActive:                  false,
			AuthenticatedCallsCompleted:               false,
			AuthenticatedPartiallyRegisteredPhone:     false,
			AuthenticatedRegisteredPhones:             false,
			CallManagerHeartBeat:                      false,
			CumulativeAllocatedResourceCannotOpenPort: false,
			EncryptedCallsActive:                      false,
			EncryptedCallsCompleted:                   false,
			EncryptedPartiallyRegisteredPhones:        false,
			EncryptedRegisteredPhones:                 false,
			HWConferenceActive:                        false,
			HWConferenceCompleted:                     false,
			HWConferenceOutOfResources:                false,
			HWConferenceResourceActive:                false,
			HWConferenceResourceAvailable:             false,
			HWConferenceResourceTotal:                 false,
			MTPOutOfResources:                         false,
			MTPRequestsThrottled:                      false,
			MTPResourceActive:                         false,
			MTPResourceAvailable:                      false,
			MTPResourceTotal:                          true,
			RegisteredAnalogAccess:                    false,
			RegisteredMGCPGateway:                     false,
			RegisteredOtherStationDevices:             false,
			SIPLineServerAuthorizationChallenges:      false,
			SIPLineServerAuthorizationFailures:        false,
			SIPTrunkApplicationAuthorizationFailures:  false,
			SIPTrunkApplicationAuthorizations:         false,
			SIPTrunkAuthorizationFailures:             false,
			SIPTrunkAuthorizations:                    false,
			SIPTrunkServerAuthenticationChallenges:    false,
			SWConferenceActive:                        false,
			SWConferenceCompleted:                     false,
			SWConferenceOutOfResources:                false,
			SWConferenceResourceActive:                false,
			SWConferenceResourceAvailable:             false,
			SWConferenceResourceTotal:                 false,
			SystemCallsAttempted:                      false,
			TranscoderOutOfResources:                  false,
			TranscoderRequestsThrottled:               false,
			TranscoderResourceActive:                  false,
			TranscoderResourceAvailable:               false,
			TranscoderResourceTotal:                   false,
			UnEncryptedCallFailures:                   false,
			VideoCallsActive:                          false,
			VideoCallsCompleted:                       false,
			VideoOnHoldOutOfResources:                 false,
			VideoOnHoldResourceActive:                 false,
			VideoOutOfResources:                       false,
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
		MonitorNames:        []string{},
		ApiAddress:          "",
		ApiUser:             "",
		ApiPassword:         "",
		Port:                9717,
		IgnoreCertificate:   false,
		ApiTimeout:          15,
		AllowStop:           false,
		SleepBetweenRequest: 30,
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
	if !PortLimits.Validate(c.Port) {
		return errors.New("defined port not valid")
	}
	if !ApiTimeoutLimit.Validate(c.ApiTimeout) {
		return errors.New("defined API timeouts not valid")
	}
	if !SleepBetweenRequestLimit.Validate(c.SleepBetweenRequest) {
		return errors.New("defined sleep between request is not valid")
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
	a = fmt.Sprintf("%sTimeout:              [%d]\r\n", a, c.ApiTimeout)
	a = fmt.Sprintf("%sSleep time:           [%d]\r\n", a, c.SleepBetweenRequest)
	a = fmt.Sprintf("%sAllow stop:           [%t]\r\n", a, c.AllowStop)

	a = fmt.Sprintf("%s%s", a, c.Metrics.Print())
	a = fmt.Sprintf("%s%s", a, c.Log.Print())
	return a
}

func (c *Config) logFields(operation ...string) log.Fields {
	f := log.Fields{
		"monitorNames":      strings.Join(c.MonitorNames, ";"),
		"apiUser":           c.ApiUser,
		"apiAddress":        c.ApiAddress,
		"telemetryPort":     c.Port,
		"ignoreCertificate": c.IgnoreCertificate,
	}
	if len(operation) > 0 {
		for i, s := range operation {
			f[fmt.Sprintf("option_%d", i)] = s
		}
	}
	return f
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
	if name == AnnunciatorOutOfResources {
		return m.AnnunciatorOutOfResources
	}
	if name == CallsAttempted {
		return m.CallsAttempted
	}
	if name == GatewayRegistrationFailures {
		return m.GatewayRegistrationFailures
	}
	if name == GatewaysInService {
		return m.GatewaysInService
	}
	if name == GatewaysOutOfService {
		return m.GatewaysOutOfService
	}
	if name == AnnunciatorResourceActive {
		return m.AnnunciatorResourceActive
	}
	if name == AnnunciatorResourceAvailable {
		return m.AnnunciatorResourceAvailable
	}
	if name == AnnunciatorResourceTotal {
		return m.AnnunciatorResourceTotal
	}
	if name == AuthenticatedCallsActive {
		return m.AuthenticatedCallsActive
	}
	if name == AuthenticatedCallsCompleted {
		return m.AuthenticatedCallsCompleted
	}
	if name == AuthenticatedPartiallyRegisteredPhone {
		return m.AuthenticatedPartiallyRegisteredPhone
	}
	if name == AuthenticatedRegisteredPhones {
		return m.AuthenticatedRegisteredPhones
	}
	if name == CallManagerHeartBeat {
		return m.CallManagerHeartBeat
	}
	if name == CumulativeAllocatedResourceCannotOpenPort {
		return m.CumulativeAllocatedResourceCannotOpenPort
	}
	if name == EncryptedCallsActive {
		return m.EncryptedCallsActive
	}
	if name == EncryptedCallsCompleted {
		return m.EncryptedCallsCompleted
	}
	if name == EncryptedPartiallyRegisteredPhones {
		return m.EncryptedPartiallyRegisteredPhones
	}
	if name == EncryptedRegisteredPhones {
		return m.EncryptedRegisteredPhones
	}
	if name == MTPOutOfResources {
		return m.MTPOutOfResources
	}
	if name == MTPRequestsThrottled {
		return m.MTPRequestsThrottled
	}
	if name == MTPResourceActive {
		return m.MTPResourceActive
	}
	if name == MTPResourceAvailable {
		return m.MTPResourceAvailable
	}
	if name == MTPResourceTotal {
		return m.MTPResourceTotal
	}
	if name == SIPLineServerAuthorizationChallenges {
		return m.SIPLineServerAuthorizationChallenges
	}
	if name == SIPLineServerAuthorizationFailures {
		return m.SIPLineServerAuthorizationFailures
	}
	if name == SIPTrunkApplicationAuthorizationFailures {
		return m.SIPTrunkApplicationAuthorizationFailures
	}
	if name == SIPTrunkApplicationAuthorizations {
		return m.SIPTrunkApplicationAuthorizations
	}
	if name == SIPTrunkAuthorizationFailures {
		return m.SIPTrunkAuthorizationFailures
	}
	if name == SIPTrunkAuthorizations {
		return m.SIPTrunkAuthorizations
	}
	if name == SIPTrunkServerAuthenticationChallenges {
		return m.SIPTrunkServerAuthenticationChallenges
	}
	if name == SystemCallsAttempted {
		return m.SystemCallsAttempted
	}
	if name == TranscoderOutOfResources {
		return m.TranscoderOutOfResources
	}
	if name == TranscoderRequestsThrottled {
		return m.TranscoderRequestsThrottled
	}
	if name == TranscoderResourceActive {
		return m.TranscoderResourceActive
	}
	if name == TranscoderResourceAvailable {
		return m.TranscoderResourceAvailable
	}
	if name == TranscoderResourceTotal {
		return m.TranscoderResourceTotal
	}
	if name == UnEncryptedCallFailures {
		return m.UnEncryptedCallFailures
	}
	if name == VideoCallsActive {
		return m.VideoCallsActive
	}
	if name == VideoCallsCompleted {
		return m.VideoCallsCompleted
	}
	if name == VideoOnHoldOutOfResources {
		return m.VideoOnHoldOutOfResources
	}
	if name == VideoOnHoldResourceActive {
		return m.VideoOnHoldResourceActive
	}
	if name == VideoOutOfResources {
		return m.VideoOutOfResources
	}
	if name == HWConferenceActive {
		return m.HWConferenceActive
	}
	if name == HWConferenceCompleted {
		return m.HWConferenceCompleted
	}
	if name == HWConferenceOutOfResources {
		return m.HWConferenceOutOfResources
	}
	if name == HWConferenceResourceActive {
		return m.HWConferenceResourceActive
	}
	if name == HWConferenceResourceAvailable {
		return m.HWConferenceResourceAvailable
	}
	if name == HWConferenceResourceTotal {
		return m.HWConferenceResourceTotal
	}
	if name == SWConferenceActive {
		return m.SWConferenceActive
	}
	if name == SWConferenceCompleted {
		return m.SWConferenceCompleted
	}
	if name == SWConferenceOutOfResources {
		return m.SWConferenceOutOfResources
	}
	if name == SWConferenceResourceActive {
		return m.SWConferenceResourceActive
	}
	if name == SWConferenceResourceAvailable {
		return m.SWConferenceResourceAvailable
	}
	if name == SWConferenceResourceTotal {
		return m.SWConferenceResourceTotal
	}
	if name == RegisteredAnalogAccess {
		return m.RegisteredAnalogAccess
	}
	if name == RegisteredMGCPGateway {
		return m.RegisteredMGCPGateway
	}
	if name == RegisteredOtherStationDevices {
		return m.RegisteredOtherStationDevices
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
