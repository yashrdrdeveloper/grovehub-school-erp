package helpers_test

import (
    "encoding/json"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    dto "grovehub/services/foundation/internal/app/dto/person"
    helpers "grovehub/services/foundation/internal/app/validators/person/helpers"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
)

// @Testing suite=FoundationPersonValidation
func RequestSpecs() {
    Describe("NormalizeCreatePersonRequest",func(){
        It("copies normalized contacts without changing the original",func(){
            source:=fixtures.CreateRequest()
            source.RequestBody.Contacts.Emails[0].Value="  PERSON@EXAMPLE.COM "
            normalized,err:=helpers.NormalizeCreatePersonRequest(source)
            Expect(err).NotTo(HaveOccurred())
            Expect(normalized.RequestBody.Contacts.Emails[0].Value).To(Equal("person@example.com"))
            Expect(source.RequestBody.Contacts.Emails[0].Value).To(Equal("  PERSON@EXAMPLE.COM "))
            Expect(helpers.ValidateCreatePersonRequest(normalized)).To(Succeed())
        })
        It("rejects invalid emails",func(){
            source:=fixtures.CreateRequest()
            source.RequestBody.Contacts.Emails[0].Value="invalid"
            _,err:=helpers.NormalizeCreatePersonRequest(source)
            Expect(err).To(HaveOccurred())
        })
        It("rejects duplicate create contacts",func(){
            source:=fixtures.CreateRequest()
            source.RequestBody.Contacts.Emails=append(source.RequestBody.Contacts.Emails,source.RequestBody.Contacts.Emails[0])
            Expect(helpers.ValidateCreatePersonRequest(source)).To(HaveOccurred())
        })
    })
    Describe("NormalizeUpdatePersonRequest",func(){
        It("preserves omitted and explicitly cleared contacts",func(){
            omitted:=fixtures.UpdateRequest()
            normalized,err:=helpers.NormalizeUpdatePersonRequest(omitted)
            Expect(err).NotTo(HaveOccurred())
            Expect(normalized.RequestBody.Contacts.Set).To(BeFalse())
            Expect(helpers.ValidateUpdatePersonRequest(normalized)).To(Succeed())
            var cleared dto.UpdatePerson
            Expect(json.Unmarshal([]byte(`{"contacts":null}`),&cleared)).To(Succeed())
            update:=dto.UpdatePersonRequest{RequestBody:cleared}
            normalized,err=helpers.NormalizeUpdatePersonRequest(update)
            Expect(err).NotTo(HaveOccurred())
            Expect(helpers.ValidateUpdatePersonRequest(normalized)).To(Succeed())
        })
        It("normalizes supplied contacts",func(){
            var patch dto.UpdatePerson
            Expect(json.Unmarshal([]byte(`{"contacts":{"emails":[{"value":" PERSON@EXAMPLE.COM "}]}}`),&patch)).To(Succeed())
            normalized,err:=helpers.NormalizeUpdatePersonRequest(dto.UpdatePersonRequest{RequestBody:patch})
            Expect(err).NotTo(HaveOccurred())
            Expect(normalized.RequestBody.Contacts.Value.Emails[0].Value).To(Equal("person@example.com"))
            Expect(helpers.ValidateUpdatePersonRequest(normalized)).To(Succeed())
        })
        It("rejects invalid updated contacts",func(){
            var patch dto.UpdatePerson
            Expect(json.Unmarshal([]byte(`{"contacts":{"emails":[{"value":"invalid"}]}}`),&patch)).To(Succeed())
            _,err:=helpers.NormalizeUpdatePersonRequest(dto.UpdatePersonRequest{RequestBody:patch})
            Expect(err).To(HaveOccurred())
        })
        It("rejects duplicated updated contacts",func(){
            var patch dto.UpdatePerson
            Expect(json.Unmarshal([]byte(`{"contacts":{"emails":[{"value":"person@example.com"},{"value":"person@example.com"}]}}`),&patch)).To(Succeed())
            Expect(helpers.ValidateUpdatePersonRequest(dto.UpdatePersonRequest{RequestBody:patch})).To(HaveOccurred())
        })
    })
    Describe("ValidateSearchPersonRequest",func(){
        It("accepts a minimal search request",func(){
            Expect(helpers.ValidateSearchPersonRequest(dto.SearchPersonsRequest{})).To(Succeed())
        })
    })
}
