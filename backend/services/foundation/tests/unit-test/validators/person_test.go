package validators_test

import (
	"errors"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	annotationvalidation "grovehub/libs/core/annotation/validation"
	coreerrors "grovehub/libs/core/errors"
	coresearch "grovehub/libs/core/search"
	personerror "grovehub/services/foundation/constants/error"
	persondto "grovehub/services/foundation/dto/person"
	personvalidator "grovehub/services/foundation/validators/person/impl"
)

var _ = Describe("Person validation rules", func() {
	var validator = personvalidator.NewPersonValidator()

	It("trims and validates names using DTO tags and normalizes contacts before custom rules", func() {
		request := persondto.CreatePerson{
			PersonNames: persondto.PersonNames{FirstName: "  Arun ", MiddleName: "  K  ", LastName: "  Kumar  "},
			Contacts: persondto.Contacts{Emails: []persondto.ContactValue{{Value: "  Arun@Example.COM ", IsPrimary: true}}},
		}
		validated, err := annotationvalidation.Process(persondto.CreatePersonRequest{RequestBody: request}, validator.CreateRules())
		Expect(err).NotTo(HaveOccurred())
		Expect(validated.RequestBody.FirstName).To(Equal("Arun"))
		Expect(validated.RequestBody.MiddleName).To(Equal("K"))
		Expect(validated.RequestBody.LastName).To(Equal("Kumar"))
		Expect(validated.RequestBody.Contacts.Emails[0].Value).To(Equal("arun@example.com"))
		Expect(request.Contacts.Emails[0].Value).To(Equal("  Arun@Example.COM "))
	})

	It("does not repeat name validation in custom rules", func() {
		request := persondto.CreatePerson{PersonNames: persondto.PersonNames{FirstName: "  ", LastName: "Kumar"}}
		_, err := annotationvalidation.Process(persondto.CreatePersonRequest{RequestBody: request}, validator.CreateRules())
		var invalid *coreerrors.Error
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Definition().Code()).To(Equal(personerror.InvalidPersonRequest.Code()))
	})

	It("keeps duplicate and multiple-primary checks as custom domain rules", func() {
		request := persondto.CreatePerson{
			PersonNames: persondto.PersonNames{FirstName: "Arun", LastName: "Kumar"},
			Contacts: persondto.Contacts{Emails: []persondto.ContactValue{
				{Value: "ARUN@example.com"}, {Value: "arun@EXAMPLE.com"},
			}},
		}
		_, err := annotationvalidation.Process(persondto.CreatePersonRequest{RequestBody: request}, validator.CreateRules())
		var invalid *coreerrors.Error
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Definition().Code()).To(Equal(personerror.InvalidPersonContact.Code()))

		request.Contacts.Emails = []persondto.ContactValue{
			{Value: "a@example.com", IsPrimary: true}, {Value: "b@example.com", IsPrimary: true},
		}
		_, err = annotationvalidation.Process(persondto.CreatePersonRequest{RequestBody: request}, validator.CreateRules())
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Definition().Code()).To(Equal(personerror.InvalidPersonContact.Code()))
	})

	It("preserves explicit PATCH null while trimming non-null supplied fields", func() {
		firstName := "  Arun "
		request := persondto.UpdatePerson{}
		request.FirstName.Set, request.FirstName.Value = true, &firstName
		request.MiddleName.Set = true
		validated, err := annotationvalidation.Process(persondto.UpdatePersonRequest{Params:persondto.PersonByGID{GID:"018f5d04-315c-75a9-b890-df0d02e787c1"}, RequestBody:request}, validator.UpdateRules())
		Expect(err).NotTo(HaveOccurred())
		Expect(*validated.RequestBody.FirstName.Value).To(Equal("Arun"))
		Expect(validated.RequestBody.MiddleName.Set).To(BeTrue())
		Expect(validated.RequestBody.MiddleName.Value).To(BeNil())
		Expect(*request.FirstName.Value).To(Equal("  Arun "))
	})

	It("validates generic search requests using SearchRules", func() {
		normalized, err := annotationvalidation.Process(persondto.SearchPersonsRequest{RequestBody:coresearch.Request{}}, validator.SearchRules())
		Expect(err).NotTo(HaveOccurred())
		Expect(normalized).To(Equal(persondto.SearchPersonsRequest{}))

		_, err = annotationvalidation.Process(persondto.SearchPersonsRequest{RequestBody:coresearch.Request{
			SearchText: strings.Repeat("x", 257),
		}}, validator.SearchRules())
		var applicationError *coreerrors.Error
		Expect(errors.As(err, &applicationError)).To(BeTrue())
		Expect(applicationError.Definition().Code()).To(Equal(personerror.InvalidPersonSearch.Code()))
	})

	It("enforces PATCH required names and max lengths through tags", func() {
		tooLong := strings.Repeat("a", 101)
		request := persondto.UpdatePerson{}
		request.FirstName.Set, request.FirstName.Value = true, &tooLong
		_, err := annotationvalidation.Process(persondto.UpdatePersonRequest{Params:persondto.PersonByGID{GID:"018f5d04-315c-75a9-b890-df0d02e787c1"}, RequestBody:request}, validator.UpdateRules())
		var invalid *coreerrors.Error
		Expect(errors.As(err, &invalid)).To(BeTrue())
		Expect(invalid.Definition().Code()).To(Equal(personerror.InvalidPersonRequest.Code()))
	})
})
