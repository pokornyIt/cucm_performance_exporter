package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"
)

const (
	applicationName = "CUCM PerfMon Exporter 2020.03.10"                     // application name
	httpApplication = "cucm-permon-exporter/2020.03.10"                      //  http user agent
	letterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // map for random string
	letterIdxBits   = 6                                                      // 6 bits to represent a letter index
	letterIdxMask   = 1<<letterIdxBits - 1                                   // All 1-bits, as many as letterIdxBits
	letterIdxMax    = 63 / letterIdxBits                                     // # of letter indices fitting in 63 bits
	maxRandomSize   = 10                                                     // required size of random string
	shortBodyChars  = 120                                                    // Max length print from string
	TimeFormat      = "15:04:05.0000"                                        // time format
)

var src = rand.NewSource(time.Now().UnixNano())
var (
	help          bool           // show help?
	toStopChannel chan bool      // used for setup stop
	monitors      perfMonService // registered service
	buildTime     string         // update in build
	commitHash    string         // commit hash
)

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
	log.WithFields(log.Fields{"operation": "monitoringOpenSession"}).Trace("try open new monitor session")
	if err = monitors.OpenSession(); err != nil {
		log.WithFields(log.Fields{"operation": "monitoringOpenSession"}).Fatal("problem open monitor session to target server")
		return true, err
	}
	monitors.AddCounters()
	return false, nil
}

func monitoringProcess() {
	var requiredStop bool
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	log.WithFields(log.Fields{"operation": "monitoringProcess"}).Debug("start with configuration")
	requiredStop = false
	log.WithFields(log.Fields{"operation": "monitoringProcess"}).Trace("read performance counters and description")
	monitors = *NewPerfMonServers()
	_ = monitors.ListAllCounters()

	srv := newWebServer(quit)
	go gracefullyShutdown(srv, quit, done)
	go runHttpServer(srv, done)

	if toStopChannel == nil {
		log.WithFields(log.Fields{"operation": "monitoringProcess"}).Fatal("problem start http listener")
		requiredStop = true
		return
	}

	requiredStop, err := monitoringOpenSession()
	if err != nil && err.Error() == "" {
		requiredStop = true
		return
	}

	// processing cycle
	for {
		if requiredStop {
			log.WithFields(log.Fields{"operation": "monitoringProcess"}).Debug("close existing open routines")
			monitors.CloseSession()
			<-done
			break
		}
		if !requiredStop {
			if !monitors.ExistSession() {
				_, err = monitoringOpenSession()
			}
			if err == nil {
				err = monitors.CollectSessionData()
			}
			if err != nil {
				log.WithFields(log.Fields{"operation": "monitoringProcess"}).Info("problem read data")
				monitors.CloseSession()
			} else {
				log.WithFields(log.Fields{"operation": "monitoringProcess"}).Trace("collect one data")
			}
			select {
			case requiredStop = <-toStopChannel:
				requiredStop = true
				break
			case <-time.After(time.Second * 5):
				break
			}
		}
	}
	log.WithFields(log.Fields{"operation": "monitoringProcess"}).Debug("procedure ends")
}

func main() {
	timeStart := time.Now()
	exitCode := 0
	kingpin.Version(VersionDetail())
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	initLog()

	err := config.LoadFile(*configFile)
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

	timeEnd := time.Now()
	log.WithFields(log.Fields{"duration": timeEnd.Sub(timeStart).String()}).Infof("program end at %s", time.Now().Format(TimeFormat))
	time.Sleep(time.Second * 2)
	os.Exit(exitCode)
}
