package entity_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const validGID = "01890f57-7b5c-7cc6-98b2-1f3d4e5a6b7c"

var fixedTime = time.Date(2026, time.October, 23, 16, 0, 0, 0, time.UTC)

func TestFoundationEntityUnit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Foundation Entity Unit Tests")
}
