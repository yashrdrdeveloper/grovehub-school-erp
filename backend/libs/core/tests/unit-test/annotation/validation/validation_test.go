package validation_test

import (
	"encoding/json"
	"errors"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	annotation "grovehub/libs/core/annotation"

)


type SampleNames struct {
	FirstName  string `json:"firstName" normalize:"trim" validate:"notblank,maxrunes=100"`
	MiddleName string `json:"middleName,omitempty" normalize:"trim" validate:"maxrunes=100"`
	LastName   string `json:"lastName" normalize:"trim" validate:"notblank,maxrunes=100"`
}

type sampleValue struct {
	Value string `json:"value" normalize:"trim" validate:"notblank"`
}

type sampleContacts struct {
	Emails []sampleValue `json:"emails,omitempty" validate:"valid"`
	Phones []sampleValue `json:"phones,omitempty" validate:"valid"`
}

type sampleCreate struct {
	SampleNames
	Contacts sampleContacts `json:"contacts,omitempty" validate:"valid"`
}

type samplePatch struct {
	FirstName  annotation.Patch[string]         `json:"firstName" normalize:"trim" validate:"nonnull,notblank,maxrunes=100"`
	MiddleName annotation.Patch[string]         `json:"middleName" normalize:"trim" validate:"maxrunes=100"`
	LastName   annotation.Patch[string]         `json:"lastName" normalize:"trim" validate:"nonnull,notblank,maxrunes=100"`
	Contacts   annotation.Patch[sampleContacts] `json:"contacts" validate:"valid"`
}

var _ = Describe("DTO tag stages", func() {
	It("normalizes embedded names before enforcing notblank/length rules", func() {
		input := sampleCreate{
			SampleNames: SampleNames{
				FirstName: "  Amit  ",
				MiddleName: "  Kumar  ",
				LastName: "  Singh  ",
			},
		}
		normalized, err := annotation.NormalizeFields(input)
		Expect(err).NotTo(HaveOccurred())
		Expect(normalized.FirstName).To(Equal("Amit"))
		Expect(normalized.MiddleName).To(Equal("Kumar"))
		Expect(normalized.LastName).To(Equal("Singh"))
		Expect(input.FirstName).To(Equal("  Amit  "))
		Expect(annotation.Validate(normalized)).To(Succeed())
	})

	It("rejects trimmed blank strings with a field-specific rule", func() {
		in := sampleCreate{SampleNames: SampleNames{FirstName: "   ", LastName: "Smith"}}
		normalized, err := annotation.NormalizeFields(in)
		Expect(err).NotTo(HaveOccurred())
		err = annotation.Validate(normalized)
		var invalid annotation.ValidationError
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Field).To(Equal("SampleNames.FirstName"))
		Expect(invalid.Rule).To(Equal("notblank"))
	})

	It("uses rune counts after trimming and rejects oversized names", func() {
		allowed := sampleCreate{SampleNames: SampleNames{FirstName: " " + strings.Repeat("ऋ", 100) + " ", LastName: "Singh"}}
		normalized, err := annotation.NormalizeFields(allowed)
		Expect(err).NotTo(HaveOccurred())
		Expect(annotation.Validate(normalized)).To(Succeed())
		tooLong := allowed
		tooLong.FirstName = strings.Repeat("ऋ", 101)
		_, err = annotation.Process(tooLong, annotation.ValidationRules[sampleCreate]{})
		var invalid annotation.ValidationError
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Rule).To(Equal("maxrunes"))
	})

	It("recursively validates nested contact slices after normalization", func() {
		input := sampleCreate{
			SampleNames: SampleNames{FirstName: "Riya", LastName: "Shah"},
			Contacts: sampleContacts{Emails: []sampleValue{{Value: "   "}}},
		}
		_, err := annotation.Process(input, annotation.ValidationRules[sampleCreate]{})
		var invalid annotation.ValidationError
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Field).To(Equal("Contacts.Emails[0].Value"))
		Expect(input.Contacts.Emails[0].Value).To(Equal("   ")) // no in-place slice mutation
	})

	It("preserves PATCH omission and null semantics while checking supplied fields", func() {
		cases := []struct {
			payload string
			field string
			rule string
		}{
			{`{}`, "", ""},
			{`{"middleName":null,"contacts":null}`, "", ""},
			{`{"firstName":null}`, "FirstName", "nonnull"},
			{`{"lastName":null}`, "LastName", "nonnull"},
			{`{"firstName":"   "}`, "FirstName", "notblank"},
			{`{"lastName":" ` + strings.Repeat("a", 101) + ` "}`, "LastName", "maxrunes"},
			{`{"contacts":{"phones":[{"value":"  "}]}}`, "Contacts.Phones[0].Value", "notblank"},
		}
		for _, testcase := range cases {
			var request samplePatch
			Expect(json.Unmarshal([]byte(testcase.payload), &request)).To(Succeed())
			_, err := annotation.Process(request, annotation.ValidationRules[samplePatch]{})
			if testcase.field == "" {
				Expect(err).NotTo(HaveOccurred(), testcase.payload)
				continue
			}
			var invalid annotation.ValidationError
			Expect(errors.As(err, &invalid)).To(BeTrue(), testcase.payload)
			Expect(invalid.Field).To(Equal(testcase.field), testcase.payload)
			Expect(invalid.Rule).To(Equal(testcase.rule), testcase.payload)
		}
	})

	It("trims a supplied PATCH name without modifying the original pointer", func() {
		var original samplePatch
		Expect(json.Unmarshal([]byte(`{"firstName":"  Asha  "}`), &original)).To(Succeed())
		normalized, err := annotation.Process(original, annotation.ValidationRules[samplePatch]{})
		Expect(err).NotTo(HaveOccurred())
		Expect(*normalized.FirstName.Value).To(Equal("Asha"))
		Expect(*original.FirstName.Value).To(Equal("  Asha  "))
	})
})
