package integrationtest

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/mrbuk/sop112_exporter/device"
)

var _ = Describe("Detect", func() {

	var broadcast = "192.168.178.255"

	Context("Devices present", func() {
		It("will detect 2 devices", func() {
			detectedSockets, err := device.Detect(broadcast)
			Expect(err).NotTo(HaveOccurred())
			Expect(detectedSockets).ShouldNot(BeEmpty())
		})
	})

})
