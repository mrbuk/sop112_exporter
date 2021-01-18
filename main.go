package main

import (
	"flag"
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
	listenAddrPtr := flag.String("listen", ":9132", "address to expose metrics endpoint")
	broadcastPtr := flag.String("broadcast", "", "broadcast address to be used for device detected (required)")
	flag.Parse()

	if *broadcastPtr == "" {
		flag.Usage()
		os.Exit(1)
	}

	// detect devices in network
	devices, err := device.Detect(*broadcastPtr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Detected devices %s", devices)

	// fetch measurements from SOP112 devices
	go collectMetrics(devices)

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	log.Printf("Listening on %s", *listenAddrPtr)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddrPtr, nil))
}

func collectMetrics(devices []*device.Powersocket) {
	c := metrics.NewMetricCollector(powerConsumption, errs)

	// set timeout of device to be shorter than the sleep
	// to avoid calls pilling up in case of connectivity issues
	for _, device := range devices {
		device.SetTimeout(3 * time.Second)
	}

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
