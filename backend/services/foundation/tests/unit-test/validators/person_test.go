package validators_test

import (
	"errors"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	annotationvalidation "grovehub/libs/core/annotation/validation"
	coreerrors "grovehub/libs/core/errors"
	personerror "grovehub/services/foundation/constants/error"
	foundationrequest "grovehub/services/foundation/dto/request"
	personvalidator "grovehub/services/foundation/validators/impl"
)

var _ = Describe("Person validation rules", func() {
	var validator = personvalidator.NewPersonValidator()

	It("trims and validates names using DTO tags and normalizes contacts before custom rules", func() {
		request := foundationrequest.CreatePerson{
			Person: foundationrequest.Person{FirstName: "  Arun ", MiddleName: "  K  ", LastName: "  Kumar  "},
			Contacts: foundationrequest.Contacts{Emails: []foundationrequest.ContactValue{{Value: "  Arun@Example.COM ", IsPrimary: true}}},
		}
		validated, err := annotationvalidation.Process(request, validator.CreateRules())
		Expect(err).NotTo(HaveOccurred())
		Expect(validated.FirstName).To(Equal("Arun"))
		Expect(validated.MiddleName).To(Equal("K"))
		Expect(validated.LastName).To(Equal("Kumar"))
		Expect(validated.Contacts.Emails[0].Value).To(Equal("arun@example.com"))
		Expect(request.Contacts.Emails[0].Value).To(Equal("  Arun@Example.COM "))
	})

	It("does not repeat name validation in custom rules", func() {
		request := foundationrequest.CreatePerson{Person: foundationrequest.Person{FirstName: "  ", LastName: "Kumar"}}
		_, err := annotationvalidation.Process(request, validator.CreateRules())
		var invalid *coreerrors.Error
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Definition().Code()).To(Equal(personerror.InvalidPersonRequest.Code()))
	})

	It("keeps duplicate and multiple-primary checks as custom domain rules", func() {
		request := foundationrequest.CreatePerson{
			Person: foundationrequest.Person{FirstName: "Arun", LastName: "Kumar"},
			Contacts: foundationrequest.Contacts{Emails: []foundationrequest.ContactValue{
				{Value: "ARUN@example.com"}, {Value: "arun@EXAMPLE.com"},
			}},
		}
		_, err := annotationvalidation.Process(request, validator.CreateRules())
		var invalid *coreerrors.Error
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Definition().Code()).To(Equal(personerror.InvalidPersonContact.Code()))

		request.Contacts.Emails = []foundationrequest.ContactValue{
			{Value: "a@example.com", IsPrimary: true}, {Value: "b@example.com", IsPrimary: true},
		}
		_, err = annotationvalidation.Process(request, validator.CreateRules())
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Definition().Code()).To(Equal(personerror.InvalidPersonContact.Code()))
	})

	It("preserves explicit PATCH null while trimming non-null supplied fields", func() {
		firstName := "  Arun "
		request := foundationrequest.UpdatePerson{}
		request.FirstName.Set, request.FirstName.Value = true, &firstName
		request.MiddleName.Set = true
		validated, err := annotationvalidation.Process(request, validator.UpdateRules())
		Expect(err).NotTo(HaveOccurred())
		Expect(*validated.FirstName.Value).To(Equal("Arun"))
		Expect(validated.MiddleName.Set).To(BeTrue())
		Expect(validated.MiddleName.Value).To(BeNil())
		Expect(*request.FirstName.Value).To(Equal("  Arun "))
	})

	It("enforces PATCH required names and max lengths through tags", func() {
		tooLong := strings.Repeat("a", 101)
		request := foundationrequest.UpdatePerson{}
		request.FirstName.Set, request.FirstName.Value = true, &tooLong
		_, err := annotationvalidation.Process(request, validator.UpdateRules())
		var invalid *coreerrors.Error
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Definition().Code()).To(Equal(personerror.InvalidPersonRequest.Code()))
	})
})
