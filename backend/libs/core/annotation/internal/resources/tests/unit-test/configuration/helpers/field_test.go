package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"reflect"
	"grovehub/libs/core/annotation/internal/src/configuration/helpers"
)

var _ = Describe("DecodeField", func() {
	It("ignores untagged fields and handles missing optional values", func() {
		value := struct { Untagged string; Optional string `env:"OPTIONAL"` }{Untagged:"keep"}
		target := reflect.ValueOf(&value).Elem()
		lookup := func(string)(string,bool){return "",false}
		for i:=0; i<target.NumField(); i++ {
			Expect(helpers.DecodeField(target.Field(i), target.Type().Field(i), lookup)).To(Succeed())
		}
		Expect(value.Untagged).To(Equal("keep"))
		Expect(value.Optional).To(BeEmpty())
	})
	It("rejects unexported fields and blank environment keys", func() {
		value := struct { private string `env:"PRIVATE"`; Empty string `env:" "` }{}
		target := reflect.ValueOf(&value).Elem()
		for i:=0; i<target.NumField(); i++ {
			Expect(helpers.DecodeField(target.Field(i), target.Type().Field(i), func(string)(string,bool){return "x",true})).To(HaveOccurred())
		}
	})
})
