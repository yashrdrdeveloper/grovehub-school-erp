package postgres_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCorePostgresUnit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Core PostgreSQL Unit Tests")
}
