package rest_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	annotationvalidation "grovehub/libs/core/annotation/validation"
	httprequest "grovehub/libs/core/http/request"
	rest "grovehub/libs/core/pipeline/rest"
)

type input struct {
	Name string `json:"name" normalize:"trim" validate:"notblank,maxrunes=10"`
}

func request(body string) *http.Request {
	return httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(body))
}

var _ = Describe("REST request pipeline", func() {
	It("runs decode -> tagged normalization -> resource normalization -> DTO rules -> custom validation", func() {
		order := []string{}
		value, err := rest.Run(request(`{"name":"  Anu  "}`), rest.Config[input]{
			Decode: func(r *http.Request, dst any) error {
				order = append(order, "decode")
				return httprequest.DecodeStrictJSON(r, dst)
			},
			Rules: annotationvalidation.Rules[input]{
				Normalize: func(v input) (input, error) {
					order = append(order, "resource-normalize")
					Expect(v.Name).To(Equal("Anu")) // tag normalization happened first
					v.Name = "Canonical"
					return v, nil
				},
				Validate: func(v input) error {
					order = append(order, "custom-validate")
					Expect(v.Name).To(Equal("Canonical")) // DTO validation passed
					return nil
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(value.Name).To(Equal("Canonical"))
		Expect(order).To(Equal([]string{"decode", "resource-normalize", "custom-validate"}))
	})

	It("validates the trimmed DTO, rejecting whitespace-only values before custom rules", func() {
		called := false
		_, err := rest.Run(request(`{"name":"   "}`), rest.Config[input]{
			Rules: annotationvalidation.Rules[input]{
				Validate: func(v input) error { called = true; return nil },
			},
		})
		var invalid annotationvalidation.Error
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Field).To(Equal("Name"))
		Expect(invalid.Rule).To(Equal("notblank"))
		Expect(called).To(BeFalse())
	})

	It("uses normalized values when evaluating max length", func() {
		value, err := rest.Run(request(`{"name":" 1234567890 "}`), rest.Config[input]{})
		Expect(err).NotTo(HaveOccurred())
		Expect(value.Name).To(Equal("1234567890"))
	})

	It("runs custom validation only after DTO validation succeeds", func() {
		sentinel := errors.New("resource rule")
		_, err := rest.Run(request(`{"name":"Anu"}`), rest.Config[input]{
			Rules: annotationvalidation.Rules[input]{
				Validate: func(v input) error { return sentinel },
			},
		})
		Expect(err).To(MatchError(sentinel))
	})

	It("keeps decoding errors separate from domain validation errors", func() {
		called := false
		_, err := rest.Run(request(`{"name":12}`), rest.Config[input]{
			Rules: annotationvalidation.Rules[input]{
				Normalize: func(v input) (input, error) { called = true; return v, nil },
			},
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

	It("stops on resource normalization errors before DTO and custom checks", func() {
		sentinel := errors.New("invalid canonical value")
		validated := false
		_, err := rest.Run(request(`{"name":"Sam"}`), rest.Config[input]{
			Rules: annotationvalidation.Rules[input]{
				Normalize: func(v input) (input, error) { return v, sentinel },
				Validate: func(v input) error { validated = true; return nil },
			},
		})
		Expect(err).To(MatchError(sentinel))
		Expect(validated).To(BeFalse())
	})

	It("maps only DTO rule failures through the configured mapper", func() {
		sentinel := errors.New("mapped DTO failure")
		_, err := rest.Run(request(`{"name":"  "}`), rest.Config[input]{
			Rules: annotationvalidation.Rules[input]{
				MapDTOError: func(err error) error {
					var invalid annotationvalidation.Error
					Expect(errors.As(err, &invalid)).To(BeTrue())
					return sentinel
				},
			},
		})
		Expect(err).To(MatchError(sentinel))
	})
})
