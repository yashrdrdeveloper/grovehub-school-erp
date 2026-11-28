package gormtable_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGormTable(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GormTable Suite")
}
