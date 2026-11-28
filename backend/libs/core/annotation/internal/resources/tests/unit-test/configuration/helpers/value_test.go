package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/annotation/internal/src/configuration/helpers"
	"reflect"
	"time"
)

var _ = Describe("DecodeField", func() {
	DescribeTable("rejects invalid typed values without exposing input",
		func(value any, text string) {
			target := reflect.ValueOf(value).Elem()
			field := reflect.StructField{Name: "Secret", Tag: reflect.StructTag(`env:"SECRET"`)}
			err := helpers.DecodeField(target, field, func(string) (string, bool) { return text, true })
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).NotTo(ContainSubstring(text))
		},
		Entry("duration", new(time.Duration), "secret-duration"),
		Entry("boolean", new(bool), "secret-bool"),
		Entry("integer", new(int8), "999999"),
		Entry("unsupported", new(float64), "secret-float"),
	)
})
