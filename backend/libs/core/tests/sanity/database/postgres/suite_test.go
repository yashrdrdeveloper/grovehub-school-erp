package postgres_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestProjectDatabaseIsolation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ProjectDatabaseIsolation Suite")
}
