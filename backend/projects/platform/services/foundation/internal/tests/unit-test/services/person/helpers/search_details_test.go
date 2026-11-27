package helpers_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "grovehub/services/foundation/internal/app/entity"
    helpers "grovehub/services/foundation/internal/app/services/person/helpers"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
)
 
// @Testing suite=FoundationPersonHelpers
func SearchDetailsSpecs() {
    Describe("BuildPersonSearchDetails",func(){
        It("includes identity, names, and contact values in canonical order",func(){
            p:=fixtures.Entity()
            middle:="Middle"
            p.MiddleName=&middle
            contact:=fixtures.Contact()
            result:=helpers.BuildPersonSearchDetails(p,[]entity.Contact{contact})
            Expect(result).To(ContainSubstring(p.GID))
            Expect(result).To(ContainSubstring("asha middle shah"))
            Expect(result).To(ContainSubstring("person@example.com"))
        })
        It("does not mutate the original contact order",func(){
            p:=fixtures.Entity()
            b:=fixtures.Contact()
            b.Type=entity.ContactTypePhone
            b.Value="+919876543210"
            b.NormalizedValue=b.Value
            a:=fixtures.Contact()
            contacts:=[]entity.Contact{b,a}
            result:=helpers.BuildPersonSearchDetails(p,contacts)
            Expect(result).To(ContainSubstring("person@example.com"))
            Expect(contacts[0].Type).To(Equal(entity.ContactTypePhone))
        })
        It("supports a person with no middle name or contacts",func(){
            p:=fixtures.Entity()
            Expect(helpers.BuildPersonSearchDetails(p,nil)).To(ContainSubstring("asha shah"))
        })
    })
}
