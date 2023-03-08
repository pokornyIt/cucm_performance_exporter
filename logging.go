package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	Routine     = "routine"   // flag for routine name
	RequestId   = "requestId" // unique request ID
	MetricsName = "metricsName"
)

func validLogLevel(level string) log.Level {
	switch strings.ToUpper(level) {
	case "FAT", "F", "FATAL":
		return log.FatalLevel
	case "ERR", "E", "ERROR":
		return log.ErrorLevel
	case "WAR", "W", "WARNING":
		return log.WarnLevel
	case "INF", "I", "INFO":
		return log.InfoLevel
	case "TRC", "T", "TRACE":
		return log.TraceLevel
	case "DEB", "D", "DEBUG":
		return log.DebugLevel
	default:
		return log.InfoLevel
	}
}

func initLog() {
	logToFile := config.Log.LogToFile()
	if help {
		logToFile = false
	}

	if config.Log.JSONFormat {
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

	log.SetReportCaller(config.Log.LogProgramInfo)
	log.SetLevel(validLogLevel(config.Log.Level))
	if logToFile {
		lJack := &lumberjack.Logger{
			Filename:   config.Log.FileName,
			MaxBackups: config.Log.MaxBackups,
			MaxAge:     config.Log.MaxAge,
			MaxSize:    config.Log.MaxSize,
			Compress:   true,
		}
		if config.Log.Quiet {
			log.SetOutput(lJack)
		} else {
			mWriter := io.MultiWriter(os.Stdout, lJack)
			log.SetOutput(mWriter)
		}
	} else {
		if config.Log.Quiet {
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
