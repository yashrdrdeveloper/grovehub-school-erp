package messaging_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMessagingCatalog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MessagingCatalog Suite")
}
