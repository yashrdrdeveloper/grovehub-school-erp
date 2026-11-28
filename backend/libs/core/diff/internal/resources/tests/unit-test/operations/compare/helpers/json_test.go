package helpers_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/diff/internal/src/operations/compare/helpers"
)

var _ = Describe("JSONValue", func() {
	It("uses JSON names, ignores hidden fields and preserves number precision", func() {
		input := struct {
			Count uint64 `json:"count"`
			Skip  string `json:"-"`
			Empty string `json:"empty,omitempty"`
		}{
			Count: 9007199254740993,
			Skip:  "hidden",
		}

		value, err := helpers.JSONValue(input)

		Expect(err).NotTo(HaveOccurred())
		Expect(value).To(Equal(map[string]any{"count": json.Number("9007199254740993")}))
	})

	It("rejects unsupported values", func() {
		value, err := helpers.JSONValue(make(chan int))

		Expect(err).To(HaveOccurred())
		Expect(value).To(BeNil())
	})
})
