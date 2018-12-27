package metrics

import (
	"github.com/mrbuk/sop112_exporter/device"
	"github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"
)

// Collector describes how metrics are collected
//go:generate counterfeiter . Collector
type Collector interface {
	Collect(device device.Powersocket)
}

// MetricCollector allows collections of Powersocket measures into
// Prometheus metrics
type MetricCollector struct {
	metric *prometheus.GaugeVec
	errs   *prometheus.CounterVec
}

// NewMetricCollector creates a new MetricCollector
func NewMetricCollector(metric *prometheus.GaugeVec, errs *prometheus.CounterVec) *MetricCollector {
	return &MetricCollector{
		metric: metric,
		errs:   errs,
	}
}

// Collect the metrics of a Powersocket into the provided GaugeVec
func (c *MetricCollector) Collect(deviceLabel string, device device.Measureable) {
	measurement, err := device.Get()
	if err != nil {
		log.Errorf("Error fetching metrics: %v", err)
		c.errs.With(prometheus.Labels{"device": deviceLabel}).Inc()
		return
	}

	c.metric.With(prometheus.Labels{"device": deviceLabel}).Set(measurement)
}
