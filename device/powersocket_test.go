package device_test

import (
	"encoding/json"
	"net/http"
	"time"

	. "github.com/mrbuk/sop112_exporter/device"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

func convertToJson(s string, target interface{}) {
	r := []byte(s)
	err := json.Unmarshal(r, target)
	Expect(err).NotTo(HaveOccurred())
}

var _ = Describe("Measurement", func() {

	var server *ghttp.Server
	var statusCode int
	var powersocket *Powersocket
	var result map[string]interface{}
	var rawResult string

	BeforeEach(func() {
		server = ghttp.NewServer()
		powersocket = NewPowersocket("SWP1234", server.Addr())
		powersocket.SetTimeout(500 * time.Millisecond)

		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET", "/"),
			ghttp.VerifyFormKV("cmd", "511"),
			ghttp.RespondWithJSONEncodedPtr(&statusCode, &result),
		))
	})

	AfterEach(func() {
		result = make(map[string]interface{})
		server.Close()
	})

	Context("Well formed message", func() {
		BeforeEach(func() {
			statusCode = http.StatusOK
			rawResult = `{"response":511,"code":200,"data":{"watt":["18.33"],"amp":["0.1"],"switch":[1]}}`
			convertToJson(rawResult, &result)
		})

		It("results in a correct measurement", func() {
			measurement, err := powersocket.Get()
			Expect(err).NotTo(HaveOccurred())
			Expect(measurement).To(Equal(18.33))
		})
	})

	Context("No data field returned", func() {
		BeforeEach(func() {
			statusCode = http.StatusOK
			rawResult = `{"response":511,"code":500}`
			convertToJson(rawResult, &result)
		})

		It("results in an error", func() {
			_, err := powersocket.Get()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Empty data field returned", func() {
		BeforeEach(func() {
			statusCode = http.StatusInternalServerError
			rawResult = `{"response":511,"code":500,"data":{}}`
			convertToJson(rawResult, &result)
		})

		It("results in an error", func() {
			_, err := powersocket.Get()
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Endpoint not reachable", func() {
		BeforeEach(func() {
			powersocket = NewPowersocket("SWP1234", "192.0.2.1")
			powersocket.SetTimeout(500 * time.Millisecond)
		})

		It("result in an error", func() {
			_, err := powersocket.Get()
			Expect(err).To(HaveOccurred())
		})
	})
})
