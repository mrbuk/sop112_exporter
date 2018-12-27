package device_test

import (
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDevice(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Device Suite")
}

var _ = Describe("Ensure GinkgoWriter is used for logging", func() {
	BeforeSuite(func() {
		log.SetOutput(GinkgoWriter)
	})
})
