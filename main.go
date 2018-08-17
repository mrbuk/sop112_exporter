package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"time"

	log "github.com/sirupsen/logrus"
)

var (
	netClient = &http.Client{
		Timeout: time.Second * 10,
	}

	devices = []Powersocket{
		Powersocket{Name: "SWP1040003000954", IP: "192.168.178.108"},
		Powersocket{Name: "SWP1040003000536", IP: "192.168.178.107"},
	}

	powerConsumption = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "powersocket_consumption_watts",
		Help: "Current power consumption of the socket",
	},
		[]string{"device"},
	)
)

func init() {
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
	for {
		for _, device := range devices {
			go collectMetric(device)
		}
		time.Sleep(time.Duration(1 * time.Second))
	}
}

func collectMetric(device Powersocket) {
	measurement, err := fetchMeasurement(device.IP)
	if err != nil {
		log.Errorf("Error fetching metrics: %v", err)
		return
	}

	powerConsumption.With(prometheus.Labels{"device": device.Name}).Set(measurement)
}

func fetchMeasurement(ip string) (float64, error) {
	// return rand.Float64(), nil
	//response, _ := netClient.Get(url)

	var consumption PowerConsumption

	log.Debugf("Checking '%s' for consumption data", ip)
	response, err := netClient.Get(fmt.Sprintf("http://%s/?cmd=511", ip))
	if err != nil {
		return 0.0, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != err {
		return 0.0, err
	}

	err = json.Unmarshal(body, &consumption)
	if err != nil {
		return 0.0, err
	}

	rawWattage := consumption.Data.Watt
	if len(rawWattage) < 1 {
		return 0.0, errors.New("No measures")
	}

	wattage, err := strconv.ParseFloat(rawWattage[0], 64)
	if err != nil {
		return 0.0, err
	}

	return wattage, nil
}
