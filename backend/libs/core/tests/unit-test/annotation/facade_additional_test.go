package annotation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation"
)

var _ = Describe("Annotation public metadata and configuration",func(){
	It("provides catalog entries, use cases and patterns",func(){
		Expect(annotation.All()).NotTo(BeEmpty())
		Expect(annotation.Uses()).NotTo(BeEmpty())
		Expect(annotation.Patterns()).NotTo(BeEmpty())
		directive,found,err:=annotation.ParseComment("// @Entity sample.records")
		Expect(err).NotTo(HaveOccurred())
		Expect(found).To(BeTrue())
		Expect(annotation.ValidateDirective(directive)).To(Succeed())
	})
	It("decodes explicit and environment-based configuration",func(){
		type settings struct {
			Host string `env:"CORE_COVERAGE_TEST_HOST" required:"true"`
			Port int `env:"CORE_COVERAGE_TEST_PORT" default:"8080"`
		}
		GinkgoT().Setenv("CORE_COVERAGE_TEST_HOST"," localhost ")
		GinkgoT().Setenv("CORE_COVERAGE_TEST_PORT","9090")
		env,err:=annotation.ConfigurationFromEnv[settings]()
		Expect(err).NotTo(HaveOccurred())
		Expect(env).To(Equal(settings{Host:"localhost",Port:9090}))
		value,err:=annotation.DecodeConfiguration[settings](func(key string)(string,bool){
			if key=="CORE_COVERAGE_TEST_HOST" {return "example",true}
			return "",false
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal(settings{Host:"example",Port:8080}))
	})
	It("normalizes values through the public entry point",func(){
		type input struct { Name string `normalize:"trim" validate:"notblank"` }
		value,err:=annotation.Normalize(input{Name:" Ada "})
		Expect(err).NotTo(HaveOccurred())
		Expect(value.Name).To(Equal("Ada"))
	})
})
