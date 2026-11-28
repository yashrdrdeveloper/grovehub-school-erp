package pack_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestModulePack(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Module Pack Suite")
}
