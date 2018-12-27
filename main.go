package main

import (
	"net/http"
	"os"

	"github.com/mrbuk/sop112_exporter/device"
	"github.com/mrbuk/sop112_exporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"time"

	log "github.com/sirupsen/logrus"
)

var (
	devices = []*device.Powersocket{
		device.NewPowersocket("SWP1040003000954", "192.168.178.108"),
		device.NewPowersocket("SWP1040003000536", "192.168.178.107"),
	}

	powerConsumption = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "powersocket_consumption_watts",
		Help: "Current power consumption of the socket",
	},
		[]string{"device"},
	)

	errs = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "powersocket_errors",
		Help: "Errors occured talking to powersockets",
	},
		[]string{"device"},
	)
)

func init() {
	initLogging()

	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(powerConsumption, errs)
}

func main() {
	// fetch measurements from SOP112 devices
	go collectMetrics()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func collectMetrics() {
	c := metrics.NewMetricCollector(powerConsumption, errs)

	for {
		for _, device := range devices {
			go c.Collect(device.Name, device)
		}
		time.Sleep(time.Duration(5 * time.Second))
	}
}

// InitLogging initilizes the log subsystem with the default
// log level or one provided via environment variable LOG_LEVEL
func initLogging() {
	var logLevel log.Level

	defaultLogLevel := log.InfoLevel
	log.SetLevel(defaultLogLevel)

	customLoglevel := os.Getenv("LOG_LEVEL")
	if customLoglevel != "" {
		var err error
		logLevel, err = log.ParseLevel(customLoglevel)
		if err != nil {
			log.Warnf("Couldn't parse '%s'", customLoglevel)
		} else {
			log.SetLevel(logLevel)
		}
	}
}
