package middleware_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMiddlewareCatalog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MiddlewareCatalog Suite")
}
