package entity_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestEntityCatalog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EntityCatalog Suite")
}
