package strategy_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStrategyCatalog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StrategyCatalog Suite")
}
