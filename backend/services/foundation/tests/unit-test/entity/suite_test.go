package entity_test

import (
 "testing"

 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
)

func TestFoundationEntity(t *testing.T) {
 RegisterFailHandler(Fail)
 RunSpecs(t, "Foundation Entity Generation Suite")
}
