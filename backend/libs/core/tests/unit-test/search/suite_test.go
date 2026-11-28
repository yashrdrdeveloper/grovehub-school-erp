package search_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSearchFacade(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t,"Search Facade Suite")
}
