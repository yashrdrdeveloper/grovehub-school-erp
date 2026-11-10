package entity_test

import (
 "time"

 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"

 foundationentity "grovehub/services/foundation/entity"
)

var _ = Describe("Generated entity boilerplate",func(){
 It("creates Contact builders and typed accessors with the exact GORM table",func(){
  now:=time.Now().UTC()
  contact:=foundationentity.NewContactBuilder().
   WithPersonGID(foundationentity.GID("a")).
   WithType(foundationentity.ContactTypeEmail).
   WithValue("test@example.com").
   WithNormalizedValue("test@example.com").
   WithIsPrimary(true).
   WithCreatedAt(now).
   Build()
  Expect(contact.TableName()).To(Equal("foundation.person_contacts"))
  Expect(contact.GetValue()).To(Equal("test@example.com"))
  contact.SetIsVerified(true)
  Expect(contact.GetIsVerified()).To(BeTrue())
 })
 It("preserves Person's domain-specific getters while providing a builder",func(){
  person:=foundationentity.NewPersonBuilder().
   WithGID(foundationentity.GID("p")).
   WithFirstName("Ada").
   WithStatus(foundationentity.PersonStatusActive).
   Build()
  Expect(person.TableName()).To(Equal("foundation.persons"))
  Expect(person.GetFirstName()).To(Equal("Ada"))
  Expect(person.GetMiddleName()).To(BeEmpty())
 })
})
