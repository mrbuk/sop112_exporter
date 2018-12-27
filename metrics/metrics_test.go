package metrics_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/mrbuk/sop112_exporter/device/devicefakes"
	. "github.com/mrbuk/sop112_exporter/metrics"

	dto "github.com/prometheus/client_model/go"
)

var _ = Describe("Metrics", func() {

	var collector *MetricCollector
	var metrics *prometheus.GaugeVec
	var errs *prometheus.CounterVec

	var fakeDevice *devicefakes.FakeMeasureable

	var fakeDeviceName string

	BeforeEach(func() {
		metrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "test_metric",
		},
			[]string{"device"},
		)
		errs = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "test_errs",
		},
			[]string{"device"},
		)
		collector = NewMetricCollector(metrics, errs)
		fakeDevice = new(devicefakes.FakeMeasureable)

		fakeDeviceName = "Test-Device"
	})

	Context("Metrics collected", func() {
		It("set a new value for powerconsumption metric", func() {
			fakeDevice.GetReturns(1.0, nil)

			collector.Collect(fakeDeviceName, fakeDevice)

			m := dto.Metric{}
			metrics.With(prometheus.Labels{"device": fakeDeviceName}).Write(&m)
			Expect(m.GetGauge().GetValue()).To(Equal(1.0))

			m = dto.Metric{}
			errs.With(prometheus.Labels{"device": fakeDeviceName}).Write(&m)
			Expect(m.GetCounter().GetValue()).To(Equal(0.0))
		})
	})

	Context("Error collecting metrics", func() {
		It("increases the error metric", func() {
			fakeDevice.GetReturns(0.0, errors.New("Something went wrong"))

			collector.Collect(fakeDeviceName, fakeDevice)

			m := dto.Metric{}
			metrics.With(prometheus.Labels{"device": fakeDeviceName}).Write(&m)
			Expect(m.GetGauge().GetValue()).To(Equal(0.0))

			m = dto.Metric{}
			errs.With(prometheus.Labels{"device": fakeDeviceName}).Write(&m)
			Expect(m.GetCounter().GetValue()).To(Equal(1.0))

		})
	})

})
