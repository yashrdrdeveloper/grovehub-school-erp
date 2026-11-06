package services_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestIntegrationServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Service Boundaries")
}
