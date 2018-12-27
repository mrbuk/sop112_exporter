package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"time"

	log "github.com/sirupsen/logrus"
)

var (
	devices = []*Powersocket{
		NewPowersocket("SWP1040003000954", "192.168.178.108"),
		NewPowersocket("SWP1040003000536", "192.168.178.107"),
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
	InitLogging()

	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(powerConsumption)
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
	c := NewMetricCollector(powerConsumption, errs)

	for {
		for _, device := range devices {
			go c.Collect(device)
		}
		time.Sleep(time.Duration(1 * time.Second))
	}
}
