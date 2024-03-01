package main

import (
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"
)

//goland:noinspection SpellCheckingInspection
const (
	applicationName      = "cucm-perfmon-exporter"                                // application name
	letterBytes          = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // map for random string
	letterIdxBits        = 6                                                      // 6 bits to represent a letter index
	letterIdxMask        = 1<<letterIdxBits - 1                                   // All 1-bits, as many as letterIdxBits
	letterIdxMax         = 63 / letterIdxBits                                     // # of letter indices fitting in 63 bits
	maxRandomSize        = 10                                                     // required size of random string
	sleepBetweenSessions = 10                                                     // sleep second between open new session or reconnect to server
)

var src = rand.NewSource(time.Now().UnixNano())
var (
	help          bool           // show help?
	toStopChannel chan bool      // used for setup stop
	monitors      PerfMonService // registered service
	Version       string         // for build data
	Revision      string         // for build data
	Branch        string         // for build data
	BuildUser     string         // for build data
	BuildDate     string         // for build data
	rateRequest   RateControl    // for limit requests
)

// httpApplicationName define user agent name for API requests
func httpApplicationName() string {
	return fmt.Sprintf("%s/%s", applicationName, version.Version)
}

// RandomString generate unique string
func RandomString() string {
	sb := strings.Builder{}
	sb.Grow(maxRandomSize)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := maxRandomSize-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

func monitoringOpenSession() (ret bool, err error) {
	log.WithFields(log.Fields{FieldRoutine: "monitoringOpenSession"}).Trace("try open new monitor session")
	defer duration(track(log.Fields{FieldRoutine: "monitoringOpenSession"}, "procedure ends"))
	if err = monitors.OpenSession(); err != nil {
		log.WithFields(log.Fields{FieldRoutine: "monitoringOpenSession"}).Error("problem open monitor session to target server")
		return true, err
	}
	monitors.AddCounters()
	prometheusCreateMetrics()
	return false, nil
}

func track(fields log.Fields, msg string) (log.Fields, string, time.Time) {
	return fields, msg, time.Now()
}
func duration(fields log.Fields, msg string, timeStart time.Time) {
	timeEnd := time.Now()
	log.WithFields(fields).WithField("duration", timeEnd.Sub(timeStart).Round(time.Millisecond).String()).Trace(msg)
}

// monitoringProcess start and run monitoring process
func monitoringProcess() {
	var roundStartTime time.Time
	var durationWait time.Duration

	defer duration(track(log.Fields{FieldRoutine: "monitoringProcess"}, "procedure ends"))
	log.WithFields(log.Fields{FieldRoutine: "monitoringProcess"}).Debug("start with configuration")
	rateRequest.reset()
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	log.WithFields(log.Fields{FieldRoutine: "monitoringProcess"}).Trace("read performance counters and description")
	monitors = *NewPerfMonServers()
	errMonitor := monitors.ListAllCounters()
	if errMonitor != nil {
		time.Sleep(2 * time.Second) // for finalize write to log
		log.WithFields(log.Fields{FieldRoutine: "monitoringProcess"}).Fatal("problem collect counters from server")
	}

	log.WithFields(log.Fields{FieldRoutine: "monitoringProcess"}).Trace("start web server and gracefully shutdown GO routines")
	srv := newWebServer(quit)
	go gracefullyShutdown(srv, quit, done)
	go runHttpServer(srv)

	if toStopChannel == nil {
		log.WithFields(log.Fields{FieldRoutine: "monitoringProcess"}).Fatal("problem start http listener")
		return
	}

	requiredStop, err := monitoringOpenSession()
	if err != nil && err.Error() == "" {
		return
	}

	// processing cycle
	for {
		roundStartTime = time.Now()
		if requiredStop {
			log.WithFields(log.Fields{FieldRoutine: "monitoringProcess"}).Debug("close existing open routines")
			monitors.CloseSession()
			<-done
			break
		}
		if !requiredStop {
			if !monitors.ExistSession() {
				log.WithFields(log.Fields{FieldRoutine: "monitoringProcess"}).Infof("session is closed wait %ds and try open new one", sleepBetweenSessions)
				time.Sleep(time.Second * sleepBetweenSessions)
				_, err = monitoringOpenSession()
			}
			if err == nil {
				err = monitors.CollectSessionData()
			}
			if err != nil {
				log.WithFields(log.Fields{FieldRoutine: "monitoringProcess"}).Info("problem read data close session")
				monitors.CloseSession()
			} else {
				log.WithFields(log.Fields{FieldRoutine: "monitoringProcess"}).Trace("collect session data")
			}
			durationWait = time.Second*time.Duration(config.SleepBetweenRequest) - time.Now().Sub(roundStartTime)
			if !monitors.ExistSession() {
				durationWait = time.Second * 60 // wait for the next try to connect to the server
			} else if durationWait < 1*time.Millisecond {
				durationWait = 1 * time.Second // Wait time is too shor wait 1 second
			}
			select {
			case <-toStopChannel:
				requiredStop = true
				log.WithFields(log.Fields{FieldRoutine: "monitoringProcess"}).Trace("request stop channel for monitoring")
				break
			case <-time.After(durationWait):
				break
			}
		}
	}
	log.WithFields(log.Fields{FieldRoutine: "monitoringProcess"}).Debug("procedure ends")
}

func main() {
	defer duration(track(log.Fields{FieldRoutine: "main"}, "program ends"))
	exitCode := 0
	version.Branch = Branch
	if len(Branch) == 0 {
		version.Branch = "Undefined"
	}
	version.Revision = Revision
	if len(Revision) == 0 {
		version.Revision = "Undefined"
	}
	version.BuildUser = BuildUser
	if len(BuildUser) == 0 {
		version.BuildUser = "Undefined"
	}
	version.BuildDate = BuildDate
	if len(BuildDate) == 0 {
		t := time.Now()
		version.BuildDate = t.Format("20060102-15:04:05")
	}
	version.Version = Version
	if len(Version) == 0 {
		version.Version = "Undefined"
	}

	kingpin.Version(version.Print(applicationName))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	err := config.LoadFile(*configFile)
	initLog()

	if *showConfig {
		fmt.Println(config.print())
		log.WithFields(log.Fields{"ApplicationName": applicationName}).Info("show only configuration ane exit")
		os.Exit(0)
	}
	if err == nil {
		monitoringProcess()
	} else {
		log.Errorf("problem with configuration. Error: %s", err)
		fmt.Printf("Program did not start due to configuration error! \r\n\tError: %s", err)
		exitCode = 1
	}

	time.Sleep(time.Second * 2)
	os.Exit(exitCode)
}
