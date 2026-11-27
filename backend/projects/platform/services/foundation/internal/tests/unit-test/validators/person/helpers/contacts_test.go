package helpers_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    dto "grovehub/services/foundation/internal/app/dto/person"
    helpers "grovehub/services/foundation/internal/app/validators/person/helpers"
    ids "grovehub/services/foundation/internal/tests/constants/person"
)

// @Testing suite=FoundationPersonValidation
func ContactsSpecs() {
    Describe("NormalizeContacts",func(){
        It("normalizes email and phone while preserving primary flags",func(){
            in:=dto.Contacts{Emails:[]dto.ContactValue{{Value:"  PERSON@EXAMPLE.COM  ",IsPrimary:true}},
                Phones:[]dto.ContactValue{{Value:ids.Phone,IsPrimary:true}}}
            result,err:=helpers.NormalizeContacts(in)
            Expect(err).NotTo(HaveOccurred())
            Expect(result.Emails).To(Equal([]dto.ContactValue{{Value:ids.Email,IsPrimary:true}}))
            Expect(result.Phones[0].IsPrimary).To(BeTrue())
            Expect(in.Emails[0].Value).To(Equal("  PERSON@EXAMPLE.COM  "))
        })
        It("rejects malformed email",func(){
            _,err:=helpers.NormalizeContacts(dto.Contacts{Emails:[]dto.ContactValue{{Value:"invalid"}}})
            Expect(err).To(HaveOccurred())
        })
        It("rejects malformed phone",func(){
            _,err:=helpers.NormalizeContacts(dto.Contacts{Phones:[]dto.ContactValue{{Value:"abc"}}})
            Expect(err).To(HaveOccurred())
        })
        It("handles empty groups",func(){
            normalized,err:=helpers.NormalizeContacts(dto.Contacts{})
            Expect(err).NotTo(HaveOccurred())
            Expect(normalized.Emails).To(BeEmpty())
            Expect(normalized.Phones).To(BeEmpty())
        })
    })
    Describe("ValidateContacts",func(){
        It("accepts one primary contact per type",func(){
            Expect(helpers.ValidateContacts(dto.Contacts{
                Emails:[]dto.ContactValue{{Value:ids.Email,IsPrimary:true}},
                Phones:[]dto.ContactValue{{Value:ids.Phone,IsPrimary:true}},
            })).To(Succeed())
        })
        It("rejects duplicate contacts",func(){
            Expect(helpers.ValidateContacts(dto.Contacts{Emails:[]dto.ContactValue{{Value:ids.Email},{Value:ids.Email}}})).To(HaveOccurred())
        })
        It("rejects two primary values of one type",func(){
            Expect(helpers.ValidateContacts(dto.Contacts{Emails:[]dto.ContactValue{{Value:ids.Email,IsPrimary:true},{Value:ids.OtherEmail,IsPrimary:true}}})).To(HaveOccurred())
        })
        It("rejects more than ten emails",func(){
            emails:=make([]dto.ContactValue,11)
            for i:=range emails{ emails[i].Value=string(rune('a'+i))+"@example.com" }
            Expect(helpers.ValidateContacts(dto.Contacts{Emails:emails})).To(HaveOccurred())
        })
        It("rejects more than ten phones",func(){
            phones:=make([]dto.ContactValue,11)
            Expect(helpers.ValidateContacts(dto.Contacts{Phones:phones})).To(HaveOccurred())
        })
        It("treats email and phone with same text as different types",func(){
            Expect(helpers.ValidateContacts(dto.Contacts{Emails:[]dto.ContactValue{{Value:"123"}},
                Phones:[]dto.ContactValue{{Value:"123"}}})).To(Succeed())
        })
    })
}
