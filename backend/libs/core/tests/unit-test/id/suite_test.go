package id_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCoreIDUnit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core ID Unit Tests")
}
