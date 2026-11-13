package services_test

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	annotationvalidation "grovehub/libs/core/annotation/validation"
	foundationrequest "grovehub/services/foundation/dto/request"
	"grovehub/services/foundation/entity"
	personservice "grovehub/services/foundation/services/impl"
	personvalidator "grovehub/services/foundation/validators/impl"
)

const patchPersonGID = "01900000-0000-7000-8000-000000000000"

func patchServiceFixture() (*fakePersonRepository, *fakeContactRepository, *fakeUnitOfWork) {
	middleName := "Kumar"
	log := new([]string)
	person := &fakePersonRepository{
		log: log,
		person: entity.Person{GID:patchPersonGID,FirstName:"Arun",MiddleName:&middleName,LastName:"Roy",Version:1,UpdatedAt:time.Now().UTC()},
	}
	contacts := &fakeContactRepository{
		log: log,
		byPerson: map[string][]entity.Contact{
			patchPersonGID: {{PersonGID:patchPersonGID,Type:entity.ContactTypeEmail,Value:"arun@example.com",NormalizedValue:"arun@example.com"}},
		},
	}
	return person, contacts, &fakeUnitOfWork{}
}

var _ = Describe("Person PATCH application", func() {
	It("clears middleName and all contacts when explicitly null", func() {
		person, contacts, work := patchServiceFixture()
		service := personservice.NewPersonService(person, contacts, work, fakeMapper{}, personvalidator.NewPersonValidator())
		var payload foundationrequest.UpdatePerson
		Expect(json.Unmarshal([]byte(`{"middleName":null,"contacts":null}`), &payload)).To(Succeed())
		_, err := service.Update(context.Background(), foundationrequest.UpdatePersonRequest{
			Params: foundationrequest.PersonByGID{GID: patchPersonGID},
			RequestBody: payload,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(person.person.MiddleName).To(BeNil())
		Expect(contacts.byPerson[patchPersonGID]).To(BeEmpty())
		Expect(contacts.replacements).To(Equal(1))
		Expect(work.runs).To(Equal(1))
		Expect(person.person.Version).To(Equal(int64(2)))
	})

	It("keeps names and contacts unchanged for omitted fields", func() {
		person, contacts, work := patchServiceFixture()
		service := personservice.NewPersonService(person, contacts, work, fakeMapper{}, personvalidator.NewPersonValidator())
		var payload foundationrequest.UpdatePerson
		Expect(json.Unmarshal([]byte(`{}`), &payload)).To(Succeed())
		_, err := service.Update(context.Background(), foundationrequest.UpdatePersonRequest{
			Params: foundationrequest.PersonByGID{GID: patchPersonGID},
			RequestBody: payload,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(*person.person.MiddleName).To(Equal("Kumar"))
		Expect(contacts.byPerson[patchPersonGID]).To(HaveLen(1))
		Expect(work.runs).To(Equal(0))
		Expect(person.person.Version).To(Equal(int64(1)))
	})

	It("rejects explicit null or blank required names", func() {
		validator := personvalidator.NewPersonValidator()
		for _, payload := range []string{`{"firstName":null}`,`{"lastName":null}`,`{"firstName":"   "}`,`{"lastName":""}`} {
			var request foundationrequest.UpdatePerson
			Expect(json.Unmarshal([]byte(payload), &request)).To(Succeed())
			_, err := annotationvalidation.Process(foundationrequest.UpdatePersonRequest{Params:foundationrequest.PersonByGID{GID:patchPersonGID}, RequestBody:request}, validator.UpdateRules())
			Expect(err).To(HaveOccurred(), payload)
		}
	})

	It("accepts an empty middle name as a removal and validates supplied names", func() {
		person, _, _ := patchServiceFixture()
		service := personservice.NewPersonService(person, &fakeContactRepository{log:person.log,byPerson:map[string][]entity.Contact{}}, &fakeUnitOfWork{}, fakeMapper{}, personvalidator.NewPersonValidator())
		var payload foundationrequest.UpdatePerson
		Expect(json.Unmarshal([]byte(`{"middleName":"  ","firstName":"  Arun  "}`), &payload)).To(Succeed())
		_, err := service.Update(context.Background(), foundationrequest.UpdatePersonRequest{
			Params: foundationrequest.PersonByGID{GID: patchPersonGID},
			RequestBody: payload,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(person.person.MiddleName).To(BeNil())
		Expect(person.person.FirstName).To(Equal("Arun"))
	})

	It("clears contact groups when a supplied object contains empty arrays", func() {
		person, contacts, _ := patchServiceFixture()
		service := personservice.NewPersonService(person, contacts, &fakeUnitOfWork{}, fakeMapper{}, personvalidator.NewPersonValidator())
		var payload foundationrequest.UpdatePerson
		Expect(json.Unmarshal([]byte(`{"contacts":{"emails":[],"phones":[]}}`), &payload)).To(Succeed())
		_, err := service.Update(context.Background(), foundationrequest.UpdatePersonRequest{
			Params: foundationrequest.PersonByGID{GID: patchPersonGID},
			RequestBody: payload,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(contacts.byPerson[patchPersonGID]).To(BeEmpty())
	})

	It("rejects overlong supplied names without mutating repository values", func() {
		person, contacts, work := patchServiceFixture()
		service := personservice.NewPersonService(person, contacts, work, fakeMapper{}, personvalidator.NewPersonValidator())
		var payload foundationrequest.UpdatePerson
		Expect(json.Unmarshal([]byte(`{"middleName":"`+strings.Repeat("x",101)+`"}`), &payload)).To(Succeed())
		_, err := service.Update(context.Background(), foundationrequest.UpdatePersonRequest{
			Params: foundationrequest.PersonByGID{GID: patchPersonGID},
			RequestBody: payload,
		})
		Expect(err).To(HaveOccurred())
		Expect(*person.person.MiddleName).To(Equal("Kumar"))
		Expect(work.runs).To(Equal(0))
	})
})
