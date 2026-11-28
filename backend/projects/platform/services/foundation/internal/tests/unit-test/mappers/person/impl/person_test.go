package impl_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    "grovehub/services/foundation/internal/app/entity"
    mapper "grovehub/services/foundation/internal/app/mappers/person/impl"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
    ids "grovehub/services/foundation/internal/tests/constants/person"
    dto "grovehub/services/foundation/internal/app/dto/person"
)

// @Testing suite=FoundationPersonMapper
func MapperSpecs() {
    Describe("PersonMapper",func(){
        It("extracts GID from request params",func(){
            m:=mapper.NewPersonMapper()
            Expect(m.ToPersonGID(dto.PersonByGID{GID:ids.GID})).To(Equal(ids.GID))
        })
        It("returns strongly typed error metadata",func(){
            m:=mapper.NewPersonMapper()
            Expect(m.ToPersonErrorMetadata(ids.GID)).To(Equal(dto.PersonErrorMetadata{GID:ids.GID}))
        })
        It("projects all Person fields and typed contacts without mutating entities",func(){
            p:=fixtures.Entity()
            middle:="Kumar"
            p.MiddleName=&middle
            email:=fixtures.Contact()
            email.IsVerified=true
            phone:=fixtures.Contact()
            phone.Type=entity.ContactTypePhone
            phone.Value=ids.Phone
            phone.NormalizedValue=ids.Phone
            r:=mapper.NewPersonMapper().ToPersonResponse(p,[]entity.Contact{email,phone})
            Expect(r.GID).To(Equal(ids.GID))
            Expect(r.MiddleName).To(Equal(middle))
            Expect(r.Status).To(Equal("ACTIVE"))
            Expect(r.Version).To(Equal(int64(1)))
            Expect(r.CreatedAt).To(Equal(p.CreatedAt))
            Expect(r.UpdatedAt).To(Equal(p.UpdatedAt))
            Expect(r.Contacts.Emails).To(HaveLen(1))
            Expect(r.Contacts.Emails[0].IsVerified).To(BeTrue())
            Expect(r.Contacts.Phones).To(HaveLen(1))
            Expect(r.Contacts.Phones[0].Value).To(Equal(ids.Phone))
        })
        It("ignores unsupported contact types",func(){
            contact:=fixtures.Contact()
            contact.Type=entity.ContactType("UNKNOWN")
            result:=mapper.NewPersonMapper().ToPersonResponse(fixtures.Entity(),[]entity.Contact{contact})
            Expect(result.MiddleName).To(BeEmpty())
            Expect(result.Contacts.Emails).To(BeEmpty())
            Expect(result.Contacts.Phones).To(BeEmpty())
        })
    })
}
