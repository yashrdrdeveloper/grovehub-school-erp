package messaging_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	sanityTimeout  = 10 * time.Second
	sanityInterval = 250 * time.Millisecond
)

func TestInfrastructureMessagingSanity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Infrastructure Messaging Sanity Tests")
}
