package services_test

import (
 "testing"
 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
)

func TestIntegrationServices(t *testing.T){
 RegisterFailHandler(Fail)
 RunSpecs(t,"Integration Service Boundaries")
}
