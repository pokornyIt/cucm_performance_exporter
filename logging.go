package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
)

const (
	identificationId      = "id"
	identificationIdOrder = "00__id__"
)

var (
	logLevel         string
	logPath          string
	logToFile        bool
	logMaxSize       int
	logMaxAge        int
	logMaxBackups    int
	logReportCallers bool
	logJsonFormat    bool
	logSilent        bool
)

func validLogLevel() log.Level {
	switch strings.ToUpper(logLevel) {
	case "FAT", "F", "FATAL":
		logLevel = "Fatal"
		return log.FatalLevel
	case "ERR", "E", "ERROR":
		logLevel = "ERROR"
		return log.ErrorLevel
	case "WAR", "W", "WARNING":
		logLevel = "WARNING"
		return log.WarnLevel
	case "INF", "I", "INFO":
		logLevel = "INFO"
		return log.InfoLevel
	case "TRC", "T", "TRACE":
		logLevel = "TRACE"
		return log.TraceLevel
	default:
		logLevel = "DEBUG"
		return log.DebugLevel
	}
}

func initLog() {
	if help {
		logToFile = false
	}

	if logJsonFormat {
		jsonFormatter := new(log.JSONFormatter)
		jsonFormatter.TimestampFormat = "2006-01-02 15:04:05.000"
		jsonFormatter.CallerPrettyfier = func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			format := " %s:%d"
			if logToFile {
				format = "%s:%d"
			}
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf(format, filename, f.Line)
		}
		log.SetFormatter(jsonFormatter)
	} else {
		Formatter := new(log.TextFormatter)
		Formatter.TimestampFormat = "2006-01-02 15:04:05.000"
		Formatter.FullTimestamp = true
		Formatter.DisableLevelTruncation = false
		Formatter.SortingFunc = func(i []string) {
			if len(i) < 2 {
				return
			}
			idx := -1
			for j, s := range i {
				if s == identificationId {
					idx = j
				}
			}
			if idx > -1 && idx < len(i) {
				i[idx] = identificationIdOrder
			}
			sort.Strings(i)
			idx = -1
			for j, s := range i {
				if s == identificationIdOrder {
					idx = j
				}
			}
			if idx > -1 && idx < len(i) {
				i[idx] = identificationId
			}

		}
		Formatter.ForceColors = !logToFile
		Formatter.CallerPrettyfier = func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			format := " %s:%d"
			if logToFile {
				format = "%s:%d"
			}
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf(format, filename, f.Line)
		}
		log.SetFormatter(Formatter)
	}

	log.SetReportCaller(logReportCallers)
	lvl := validLogLevel()
	log.SetLevel(lvl)
	if logToFile {
		lJack := &lumberjack.Logger{
			Filename:   logPath,
			MaxBackups: logMaxBackups,
			MaxAge:     logMaxAge,
			MaxSize:    logMaxSize,
			Compress:   true,
		}
		if logSilent {
			log.SetOutput(lJack)
		} else {
			mWriter := io.MultiWriter(os.Stdout, lJack)
			log.SetOutput(mWriter)
		}
	} else {
		if logSilent {
			log.SetLevel(log.PanicLevel)
		}
	}

	if !help {
		log.WithFields(log.Fields{
			"ApplicationName": applicationName,
			"RuntimeVersion":  runtime.Version(),
			"CPUs":            runtime.NumCPU(),
			"Arch":            runtime.GOARCH,
		}).Info("application Initializing")
	}
}

func VersionDetail() string {
	return fmt.Sprintf("Version details\r\n\tApplication Name: %s\r\n\tRuntime Version: %s\r\n\tCPUs: %d\r\n\tArchitectire: %s\r\n\tBuild Time: %s\n\r\tCommit Hash: %s",
		applicationName, runtime.Version(), runtime.NumCPU(), runtime.GOARCH, buildTime, commitHash)
}
