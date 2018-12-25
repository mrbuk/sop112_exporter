package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSop112Exporter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sop112Exporter Suite")
}
