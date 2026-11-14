package helpers_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    dto "grovehub/services/foundation/internal/app/dto/person"
    "grovehub/services/foundation/internal/app/entity"
    helpers "grovehub/services/foundation/internal/app/services/person/helpers"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
    ids "grovehub/services/foundation/internal/tests/constants/person"
)

// @Testing suite=FoundationPersonHelpers
func ContactSpecs() {
    Describe("BuildContacts",func(){
        It("builds typed emails and phones and preserves verification of an existing value",func(){
            prior:=fixtures.Contact()
            prior.IsVerified=true
            request:=dto.Contacts{
                Emails:[]dto.ContactValue{{Value:ids.Email,IsPrimary:true},{Value:ids.OtherEmail}},
                Phones:[]dto.ContactValue{{Value:ids.Phone,IsPrimary:true}},
            }
            contacts:=helpers.BuildContacts(request,[]entity.Contact{prior})
            Expect(contacts).To(HaveLen(3))
            Expect(contacts[0]).To(And(HaveField("Type",entity.ContactTypeEmail),HaveField("NormalizedValue",ids.Email),HaveField("IsVerified",true)))
            Expect(contacts[1].IsVerified).To(BeFalse())
            Expect(contacts[2].Type).To(Equal(entity.ContactTypePhone))
            Expect(contacts[2].IsPrimary).To(BeTrue())
        })
        It("handles empty input without claiming that previous contacts still exist",func(){
            Expect(helpers.BuildContacts(dto.Contacts{},[]entity.Contact{fixtures.Contact()})).To(BeEmpty())
        })
    })
    Describe("AssociateContacts",func(){
        It("adds ownership and audit timestamps without mutating caller-owned input",func(){
            person:=fixtures.Entity()
            contacts:=[]entity.Contact{{Type:entity.ContactTypeEmail,Value:ids.Email,NormalizedValue:ids.Email}}
            result:=helpers.AssociateContacts(person,contacts)
            Expect(result).To(HaveLen(1))
            Expect(result[0].PersonGID).To(Equal(person.GID))
            Expect(result[0].CreatedAt).To(Equal(person.UpdatedAt))
            Expect(result[0].UpdatedAt).To(Equal(person.UpdatedAt))
            Expect(contacts[0].PersonGID).To(BeEmpty())
            result[0].Value="changed"
            Expect(contacts[0].Value).To(Equal(ids.Email))
        })
        It("retains the original creation time for an existing contact",func(){
            person:=fixtures.Entity()
            old:=fixtures.Contact()
            old.CreatedAt=person.CreatedAt.Add(-24*60*60*1e9)
            result:=helpers.AssociateContacts(person,[]entity.Contact{old})
            Expect(result[0].CreatedAt).To(Equal(old.CreatedAt))
            Expect(result[0].UpdatedAt).To(Equal(person.UpdatedAt))
        })
    })
    Describe("CopyContacts",func(){
        It("returns nil for empty input",func(){Expect(helpers.CopyContacts(nil)).To(BeNil())})
        It("performs a defensive copy",func(){
            contacts:=[]entity.Contact{fixtures.Contact()}
            copied:=helpers.CopyContacts(contacts)
            copied[0].Value="edited"
            Expect(contacts[0].Value).To(Equal(ids.Email))
        })
    })
}
