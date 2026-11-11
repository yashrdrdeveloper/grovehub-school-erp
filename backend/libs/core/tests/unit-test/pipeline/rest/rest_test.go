package rest_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	httprequest "grovehub/libs/core/http/request"
	rest "grovehub/libs/core/pipeline/rest"
)

type input struct {
	Name string `json:"name"`
}

func request(body string) *http.Request {
	return httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
}

var _ = Describe("REST request pipeline", func() {
	It("decodes first, normalizes second, validates third", func() {
		order := []string{}
		value, err := rest.Run(request(`{"name":"  Anu  "}`), rest.Config[input]{
			Decode: func(r *http.Request, dst any) error {
				order = append(order, "decode")
				return httprequest.DecodeStrictJSON(r, dst)
			},
			Normalize: func(v input) (input, error) {
				order = append(order, "normalize")
				v.Name = strings.TrimSpace(v.Name)
				return v, nil
			},
			Validate: func(v input) error {
				order = append(order, "validate")
				Expect(v.Name).To(Equal("Anu"))
				return nil
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(value.Name).To(Equal("Anu"))
		Expect(order).To(Equal([]string{"decode", "normalize", "validate"}))
	})

	It("supports defaults and optional stages", func() {
		value, err := rest.Run(request(`{"name":"Sam"}`), rest.Config[input]{})
		Expect(err).NotTo(HaveOccurred())
		Expect(value.Name).To(Equal("Sam"))
	})

	It("stops on decode errors without running later stages", func() {
		called := false
		_, err := rest.Run(request(`{"name":12}`), rest.Config[input]{
			Normalize: func(v input) (input, error) { called = true; return v, nil },
		})
		var decodeError *rest.DecodeError
		Expect(errors.As(err, &decodeError)).To(BeTrue())
		Expect(called).To(BeFalse())
	})

	It("keeps strict JSON decoding configurable", func() {
		_, err := rest.Run(request(`{"name":"Sam","unknown":true}`), rest.Config[input]{Decode: httprequest.DecodeStrictJSON})
		var decodeError *rest.DecodeError
		Expect(errors.As(err, &decodeError)).To(BeTrue())
	})

	It("stops at normalization and returns the original error", func() {
		sentinel := errors.New("invalid name")
		validated := false
		_, err := rest.Run(request(`{"name":"Sam"}`), rest.Config[input]{
			Normalize: func(v input) (input, error) { return v, sentinel },
			Validate: func(v input) error { validated = true; return nil },
		})
		Expect(err).To(MatchError(sentinel))
		Expect(validated).To(BeFalse())
	})

	It("stops at validation without invoking any service", func() {
		sentinel := errors.New("validation failed")
		_, err := rest.Run(request(`{"name":"Sam"}`), rest.Config[input]{
			Validate: func(v input) error { return sentinel },
		})
		Expect(err).To(MatchError(sentinel))
	})
})
