package id

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const expectedUUIDv7 = "01a1697d-9400-7001-8203-040506070809"

var fixedUUIDv7Time = time.Date(2026, time.October, 23, 16, 0, 0, 0, time.UTC)

func TestCoreIDUnit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core ID Unit Tests")
}
