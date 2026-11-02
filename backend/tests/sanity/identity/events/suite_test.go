package events_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	sanityTimeout  = 10 * time.Second
	sanityInterval = 100 * time.Millisecond
)

func TestIdentityEventsSanity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Identity Events Sanity Tests")
}
