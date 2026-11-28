package persistence_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPersistenceCatalog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PersistenceCatalog Suite")
}
