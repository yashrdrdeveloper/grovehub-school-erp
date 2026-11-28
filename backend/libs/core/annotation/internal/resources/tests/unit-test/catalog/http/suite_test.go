package cataloghttp_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHttpCatalog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HttpCatalog Suite")
}
