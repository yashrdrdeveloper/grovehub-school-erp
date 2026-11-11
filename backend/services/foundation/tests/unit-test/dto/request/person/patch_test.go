package person_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	foundationrequest "grovehub/services/foundation/dto/request/person"
)

var _ = Describe("Person PATCH JSON presence", func() {
	It("distinguishes omitted, explicit null, and provided values", func() {
		var omitted foundationrequest.UpdatePerson
		Expect(json.Unmarshal([]byte(`{}`), &omitted)).To(Succeed())
		Expect(omitted.MiddleName.Set).To(BeFalse())
		Expect(omitted.Contacts.Set).To(BeFalse())

		var cleared foundationrequest.UpdatePerson
		Expect(json.Unmarshal([]byte(`{"middleName":null,"contacts":null}`), &cleared)).To(Succeed())
		Expect(cleared.MiddleName.Set).To(BeTrue())
		Expect(cleared.MiddleName.Value).To(BeNil())
		Expect(cleared.Contacts.Set).To(BeTrue())
		Expect(cleared.Contacts.Value).To(BeNil())

		var supplied foundationrequest.UpdatePerson
		Expect(json.Unmarshal([]byte(`{"middleName":" Kumar ","contacts":{"emails":[]}}`), &supplied)).To(Succeed())
		Expect(supplied.MiddleName.Set).To(BeTrue())
		Expect(*supplied.MiddleName.Value).To(Equal(" Kumar "))
		Expect(supplied.Contacts.Set).To(BeTrue())
		Expect(supplied.Contacts.Value).NotTo(BeNil())
		Expect(supplied.Contacts.Value.Emails).To(BeEmpty())
	})

	It("rejects JSON types inconsistent with the declared PATCH field", func() {
		for _, payload := range []string{
			`{"firstName":123}`,
			`{"middleName":true}`,
			`{"contacts":"wrong"}`,
			`{"contacts":[]}`,
		} {
			var request foundationrequest.UpdatePerson
			Expect(json.Unmarshal([]byte(payload), &request)).NotTo(Succeed(), payload)
		}
	})

	It("round-trips omission and explicit null without inventing fields", func() {
		var request foundationrequest.UpdatePerson
		Expect(json.Unmarshal([]byte(`{"middleName":null,"firstName":"Arun"}`), &request)).To(Succeed())
		bytes, err := json.Marshal(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(bytes)).To(MatchJSON(`{"firstName":"Arun","middleName":null}`))
	})
})
