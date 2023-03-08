package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"
)

var (
	callMetrics    map[string]*prometheus.GaugeVec
	counterMetrics map[string]*prometheus.CounterVec
	counterActual  map[string]float64
)

const (
	mainPage = "<html><head><title>%s</title></head><body>%s</body></html>"
	aHref    = "<a href=\"%s\">%s</a><br>"
)

func newWebServer(quit chan<- os.Signal) *http.Server {
	toStopChannel = make(chan bool)
	router := http.NewServeMux()

	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"metricsUri": "/", Routine: "newWebServer"}).Debug("request /")
		w.WriteHeader(http.StatusOK)
		body := fmt.Sprintf(aHref, "/metrics", "Prometheus Export Data")
		body += fmt.Sprintf(aHref, "/config", "Actual configuration")
		body += fmt.Sprintf(aHref, "/version", "Program version")
		if config.AllowStop {
			body += fmt.Sprintf(aHref, "/stop", "Stop program")
		}
		_, _ = w.Write([]byte(fmt.Sprintf(mainPage, applicationName, body)))
	}))
	router.Handle("/version", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"metricsUri": "/version", Routine: "newWebServer"}).Debug("request /version")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(version.Print(applicationName)))
	}))
	router.Handle("/err", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"metricsUri": "/err", Routine: "newWebServer"}).Debug("request /err")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("Error"))
	}))
	router.Handle("/config", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"metricsUri": "/config", Routine: "newWebServer"}).Debug("request /config")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(config.print()))
	}))
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/stop", func(writer http.ResponseWriter, request *http.Request) {
		log.WithFields(log.Fields{"metricsUri": "/stop", Routine: "newWebServer"}).Infof("request from %s", request.URL.Path)
		if config.AllowStop {
			writer.WriteHeader(http.StatusOK)
			_, _ = writer.Write([]byte("Stop processing"))
			select {
			case toStopChannel <- true:
				quit <- syscall.SIGINT
				break
			case <-time.After(time.Millisecond * 100):
				break
			default:
				break
			}
		} else {
			_, _ = writer.Write([]byte("Not allowed stop"))
		}
	})
	port := fmt.Sprintf(":%d", config.Port)
	server := &http.Server{Handler: router, Addr: port}

	log.WithFields(log.Fields{"port": port, "metricsUri": "/metrics", Routine: "newWebServer"}).Infof("listener start on 0.0.0.0:%d/metrics", config.Port)
	return server
}

func prometheusCreateMetrics() {
	log.WithFields(log.Fields{Routine: "prometheusCreateMetrics"}).Infof("prepare new metrics")
	callMetrics = make(map[string]*prometheus.GaugeVec)
	counterMetrics = make(map[string]*prometheus.CounterVec)
	counterActual = make(map[string]float64)

	var counter *CounterDetails
	var err error
	if !config.Metrics.GoCollector {
		prometheus.Unregister(collectors.NewGoCollector())
	}
	if !config.Metrics.ProcessStatus {
		prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	}

	for _, cnt := range SupportedCounters {
		if !config.Metrics.enablePrometheusCounter(cnt.allowedCounterName) {
			log.WithFields(log.Fields{Routine: "newWebServer", MetricsName: cnt.prometheusName}).Infof("metrics %s not enabled", cnt.allowedCounterName)
			continue
		}
		counter, err = monitors.GetCounterDetails(cnt.allowedCounterName)
		if err != nil {
			log.WithFields(log.Fields{Routine: "newWebServer", MetricsName: cnt.prometheusName}).Errorf("not defined description for %s", cnt.allowedCounterName)
		}
		if counter == nil {
			counter = &CounterDetails{name: cnt.allowedCounterName, description: fmt.Sprintf("Description for %s not exists", cnt.allowedCounterName)}
		}
		if strings.HasSuffix(strings.ToLower(cnt.allowedCounterName), "failed") {
			counterMetrics[cnt.allowedCounterName] = prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: cnt.prometheusName,
					Help: counter.description,
				}, []string{"server"})
			prometheus.MustRegister(counterMetrics[cnt.allowedCounterName])
			counterActual[cnt.allowedCounterName] = float64(0)
		} else {
			callMetrics[cnt.allowedCounterName] = prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: cnt.prometheusName,
					Help: counter.description,
				}, []string{"server"})
			prometheus.MustRegister(callMetrics[cnt.allowedCounterName])
			for _, srv := range monitors.monitors {
				callMetrics[cnt.allowedCounterName].WithLabelValues(srv.server).Set(0)
			}
		}
	}
}

func prometheusRemoveMetrics() {
	log.WithFields(log.Fields{Routine: "prometheusCreateMetrics"}).Infof("prepare remove all metrics")
	for _, cnt := range SupportedCounters {
		if strings.HasSuffix(strings.ToLower(cnt.allowedCounterName), "failed") {
			prometheus.Unregister(counterMetrics[cnt.allowedCounterName])
			counterActual[cnt.allowedCounterName] = float64(0)
		} else {
			prometheus.Unregister(callMetrics[cnt.allowedCounterName])
			for _, srv := range monitors.monitors {
				callMetrics[cnt.allowedCounterName].WithLabelValues(srv.server).Set(0)
			}
		}
	}
}

func gracefullyShutdown(server *http.Server, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	log.WithFields(log.Fields{Routine: "gracefullyShutdown"}).Infof("server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.WithFields(log.Fields{Routine: "gracefullyShutdown"}).Infof("could not gracefully shutdown the server: %v", err)
	}
	close(done)
}

func runHttpServer(srv *http.Server, done <-chan bool) {
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.WithFields(log.Fields{"port": config.Port, "error": err, Routine: "runHttpServer"}).Errorf("listener didn't start port %d. Error: %s", config.Port, err)
		panic(fmt.Sprintf("Server not start on port %d. Error: %s", config.Port, err))
	}
	toStopChannel <- true
	log.WithFields(log.Fields{"status": "stop", Routine: "runHttpServer"}).Info("HTTP server stopped")
}
