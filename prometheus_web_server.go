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
	// callMetrics list of call gauge metrics (i.e. number of devices)
	callMetrics map[string]*prometheus.GaugeVec
	// counterMetrics list of counter metrics (i.e. failure recorded calls)
	counterMetrics map[string]*prometheus.CounterVec
	// counterActual actual presented value in counterMetrics
	counterActual map[string]float64
)

const (
	mainPage = "<html><head><title>%s</title></head><body>%s</body></html>"
	aHref    = "<a href=\"%s\">%s</a><br>"
)

// newWebServer create web server structure
func newWebServer(quit chan<- os.Signal) *http.Server {
	defer duration(track(log.Fields{FieldRoutine: "newWebServer"}, "procedure ends"))
	toStopChannel = make(chan bool)
	router := http.NewServeMux()

	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"metricsUri": "/", FieldRoutine: "newWebServer"}).Debug("request /")
		w.WriteHeader(http.StatusOK)
		body := fmt.Sprintf(aHref, "/metrics", "Prometheus Export Data")
		body += fmt.Sprintf(aHref, "/status", "Program status")
		body += fmt.Sprintf(aHref, "/config", "Program configuration")
		body += fmt.Sprintf(aHref, "/version", "Program version")
		if config.AllowStop {
			body += fmt.Sprintf(aHref, "/stop", "Stop program")
		}
		_, _ = w.Write([]byte(fmt.Sprintf(mainPage, applicationName, body)))
	}))
	router.Handle("/version", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"metricsUri": "/version", FieldRoutine: "newWebServer"}).Debug("request /version")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(version.Print(applicationName)))
	}))
	router.Handle("/status", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"metricsUri": "/err", FieldRoutine: "newWebServer"}).Debug("request /status")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(monitors.client.print()))
	}))
	router.Handle("/config", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"metricsUri": "/config", FieldRoutine: "newWebServer"}).Debug("request /config")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(config.print()))
	}))
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/stop", func(writer http.ResponseWriter, request *http.Request) {
		log.WithFields(log.Fields{"metricsUri": "/stop", FieldRoutine: "newWebServer"}).Infof("request from %s", request.URL.Path)
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

	log.WithFields(log.Fields{"port": port, "metricsUri": "/metrics", FieldRoutine: "newWebServer"}).Infof("listener start on 0.0.0.0:%d/metrics", config.Port)
	return server
}

// prometheusCreateMetrics create all necessary metrics for prometheus
func prometheusCreateMetrics() {
	log.WithFields(log.Fields{FieldRoutine: "prometheusCreateMetrics"}).Infof("prepare new metrics")
	defer duration(track(log.Fields{FieldRoutine: "prometheusCreateMetrics"}, "procedure ends"))
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

	for _, supportedCounter := range SupportedCounters {
		if !config.Metrics.enablePrometheusCounter(supportedCounter.allowedCounterName) {
			log.WithFields(log.Fields{FieldRoutine: "newWebServer", FieldMetricsName: supportedCounter.prometheusName}).Debugf("metrics %s not enabled", supportedCounter.allowedCounterName)
			continue
		}
		counter, err = monitors.GetCounterDetails(supportedCounter.allowedCounterName)
		if err != nil {
			log.WithFields(log.Fields{FieldRoutine: "newWebServer", FieldMetricsName: supportedCounter.prometheusName}).Errorf("not defined description for %s", supportedCounter.allowedCounterName)
		}
		if counter == nil {
			counter = &CounterDetails{name: supportedCounter.allowedCounterName, description: fmt.Sprintf("Description for %s not exists", supportedCounter.allowedCounterName)}
		}
		if strings.HasSuffix(strings.ToLower(supportedCounter.allowedCounterName), "failed") {
			counterMetrics[supportedCounter.allowedCounterName] = prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: supportedCounter.prometheusName,
					Help: counter.description,
				}, []string{"server"})
			prometheus.MustRegister(counterMetrics[supportedCounter.allowedCounterName])
			counterActual[supportedCounter.allowedCounterName] = float64(0)
		} else {
			callMetrics[supportedCounter.allowedCounterName] = prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: supportedCounter.prometheusName,
					Help: counter.description,
				}, []string{"server"})
			prometheus.MustRegister(callMetrics[supportedCounter.allowedCounterName])
			for _, srv := range monitors.monitors {
				callMetrics[supportedCounter.allowedCounterName].WithLabelValues(srv.server).Set(0)
			}
		}
	}
}

// prometheusRemoveMetrics remove all CUCM metrics from prometheus
func prometheusRemoveMetrics() {
	log.WithFields(log.Fields{FieldRoutine: "prometheusCreateMetrics"}).Infof("prepare remove all metrics")
	defer duration(track(log.Fields{FieldRoutine: "prometheusCreateMetrics"}, "procedure ends"))
	for _, cnt := range SupportedCounters {
		if config.Metrics.enablePrometheusCounter(cnt.allowedCounterName) {
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
}

// gracefullyShutdown shutdown all services, web servers and GO routines
func gracefullyShutdown(server *http.Server, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	log.WithFields(log.Fields{FieldRoutine: "gracefullyShutdown"}).Infof("server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.WithFields(log.Fields{FieldRoutine: "gracefullyShutdown"}).Infof("could not gracefully shutdown the server: %v", err)
	}
	close(done)
}

// runHttpServer run web server
func runHttpServer(srv *http.Server) {
	defer duration(track(log.Fields{FieldRoutine: "runHttpServer"}, "procedure ends"))
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.WithFields(log.Fields{"port": config.Port, "error": err, FieldRoutine: "runHttpServer"}).Errorf("listener didn't start port %d. Error: %s", config.Port, err)
		panic(fmt.Sprintf("server not start on port %d. Error: %s", config.Port, err))
	}
	toStopChannel <- true
	log.WithFields(log.Fields{"status": "stop", FieldRoutine: "runHttpServer"}).Info("HTTP server stopped")
}
