package person_test

import (
	"encoding/json"
	"errors"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreannotation "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation"
	persondto "grovehub/services/foundation/internal/app/dto/person"
)

var _ = Describe("DTO tag stages", func() {
	It("normalizes embedded names before enforcing notblank/length rules", func() {
		input := persondto.CreatePerson{
			PersonNames: persondto.PersonNames{
				FirstName: "  Amit  ",
				MiddleName: "  Kumar  ",
				LastName: "  Singh  ",
			},
		}
		normalized, err := coreannotation.NormalizeFields(input)
		Expect(err).NotTo(HaveOccurred())
		Expect(normalized.FirstName).To(Equal("Amit"))
		Expect(normalized.MiddleName).To(Equal("Kumar"))
		Expect(normalized.LastName).To(Equal("Singh"))
		Expect(input.FirstName).To(Equal("  Amit  "))
		Expect(coreannotation.Validate(normalized)).To(Succeed())
	})

	It("rejects trimmed blank strings with a field-specific rule", func() {
		in := persondto.CreatePerson{PersonNames: persondto.PersonNames{FirstName: "   ", LastName: "Smith"}}
		normalized, err := coreannotation.NormalizeFields(in)
		Expect(err).NotTo(HaveOccurred())
		err = coreannotation.Validate(normalized)
		var invalid coreannotation.ValidationError
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Field).To(Equal("PersonNames.FirstName"))
		Expect(invalid.Rule).To(Equal("notblank"))
	})

	It("uses rune counts after trimming and rejects oversized names", func() {
		allowed := persondto.CreatePerson{PersonNames: persondto.PersonNames{FirstName: " " + strings.Repeat("ऋ", 100) + " ", LastName: "Singh"}}
		normalized, err := coreannotation.NormalizeFields(allowed)
		Expect(err).NotTo(HaveOccurred())
		Expect(coreannotation.Validate(normalized)).To(Succeed())
		tooLong := allowed
		tooLong.FirstName = strings.Repeat("ऋ", 101)
		_, err = coreannotation.Process(tooLong, coreannotation.ValidationRules[persondto.CreatePerson]{})
		var invalid coreannotation.ValidationError
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Rule).To(Equal("maxrunes"))
	})

	It("recursively validates nested contact slices after normalization", func() {
		input := persondto.CreatePerson{
			PersonNames: persondto.PersonNames{FirstName: "Riya", LastName: "Shah"},
			Contacts: persondto.Contacts{Emails: []persondto.ContactValue{{Value: "   "}}},
		}
		_, err := coreannotation.Process(input, coreannotation.ValidationRules[persondto.CreatePerson]{})
		var invalid coreannotation.ValidationError
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
			var request persondto.UpdatePerson
			Expect(json.Unmarshal([]byte(testcase.payload), &request)).To(Succeed())
			_, err := coreannotation.Process(request, coreannotation.ValidationRules[persondto.UpdatePerson]{})
			if testcase.field == "" {
				Expect(err).NotTo(HaveOccurred(), testcase.payload)
				continue
			}
			var invalid coreannotation.ValidationError
			Expect(errors.As(err, &invalid)).To(BeTrue(), testcase.payload)
			Expect(invalid.Field).To(Equal(testcase.field), testcase.payload)
			Expect(invalid.Rule).To(Equal(testcase.rule), testcase.payload)
		}
	})

	It("trims a supplied PATCH name without modifying the original pointer", func() {
		var original persondto.UpdatePerson
		Expect(json.Unmarshal([]byte(`{"firstName":"  Asha  "}`), &original)).To(Succeed())
		normalized, err := coreannotation.Process(original, coreannotation.ValidationRules[persondto.UpdatePerson]{})
		Expect(err).NotTo(HaveOccurred())
		Expect(*normalized.FirstName.Value).To(Equal("Asha"))
		Expect(*original.FirstName.Value).To(Equal("  Asha  "))
	})
})
