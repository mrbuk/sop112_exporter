package metrics_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	log "github.com/sirupsen/logrus"
)

func TestMetrics(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Metrics Suite")
}

var _ = Describe("Ensure GinkgoWriter is used for logging", func() {
	BeforeSuite(func() {
		log.SetOutput(GinkgoWriter)
	})
})
